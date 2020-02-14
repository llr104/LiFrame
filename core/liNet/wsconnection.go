package liNet

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/llr104/LiFrame/utils"
	"github.com/thinkoner/openssl"
	"sync"
)

var GateMessageKey = []byte("liFrameVeryGood!")

// 客户端读写消息
type WsMessage struct {
	MsgType int
	Data    []byte
}

// 客户端连接
type WsConnection struct {
	wsSocket 	*websocket.Conn 			// 底层websocket
	outChan 	chan *WsMessage 			// 写队列
	mutex 		sync.Mutex					// 避免重复关闭管道
	isClosed 	bool
	id       	uint64                		// id
	onMessage	func(wsConn *WsConnection, req *WsMessage)
	onClose     func(wsConn* WsConnection)
	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

func NewWsConnection(wsSocket *websocket.Conn, cid uint64) *WsConnection {
	wsConn := &WsConnection{
		wsSocket: wsSocket,
		outChan: make(chan *WsMessage, 1000),
		isClosed:false,
		id:cid,
		property:make(map[string]interface{}),
	}

	return wsConn
}

func (wsConn *WsConnection) GetId() uint64 {
	return wsConn.id
}

func (wsConn *WsConnection) Running() {
	// 读协程
	go wsConn.wsReadLoop()
	// 写协程
	go wsConn.wsWriteLoop()
}

func (wsConn *WsConnection) wsReadLoop() {
	for {
		// 读一个message
		msgType, data, err := wsConn.wsSocket.ReadMessage()
		if err != nil {
			break
		}
		req := &WsMessage{
			msgType,
			data,
		}

		if wsConn.onMessage != nil{
			wsConn.onMessage(wsConn, req)
		}

	}

	wsConn.wsClose()
}

func (wsConn *WsConnection)SetOnMessage(hookFunc func (*WsConnection, *WsMessage))  {
	wsConn.onMessage = hookFunc
}

func (wsConn *WsConnection)SetOnClose(hookFunc func (*WsConnection))  {
	wsConn.onClose = hookFunc
}


func (wsConn *WsConnection) wsWriteLoop() {
	for {
		select {
			// 取一个消息
			case msg := <- wsConn.outChan:
				// 写给websocket
				if err := wsConn.wsSocket.WriteMessage(msg.MsgType, msg.Data); err != nil {
					utils.Log.Warn("wsWriteLoop error %s", err.Error())
					wsConn.wsClose()
					return
				}
			}
	}

}

func (wsConn *WsConnection) WriteProxyMessage(proxyName string, funcName string, body interface{})  {
	data, err := json.Marshal(body)
	if err != nil{
		return
	}

	wsConn.WriteMessage(proxyName, funcName, data)
}


func (wsConn *WsConnection) WriteObject(proxyName string, funcName string, body interface{})  {
	data, err := json.Marshal(body)
	if err != nil{
		return
	}

	wsConn.WriteMessage(proxyName, funcName, data)
}

func (wsConn *WsConnection) WriteMessage(proxyName string, funcName string, body[] byte){
	text := fmt.Sprintf("%s|%s|%s", funcName, proxyName, body)

	enData, err := openssl.AesCBCEncrypt([]byte(text), GateMessageKey, GateMessageKey, openssl.ZEROS_PADDING)

	if err != nil{
		return
	}
	data := hex.EncodeToString(enData)

	var b bytes.Buffer
	gz, _ := gzip.NewWriterLevel(&b, 9)
	if _, err := gz.Write([]byte(data)); err != nil {
		return
	}
	if err := gz.Flush(); err != nil {
		return
	}
	if err := gz.Close(); err != nil {
		return
	}

	wsConn.writeBytes(b.Bytes())
}

func (wsConn *WsConnection) writeBytes(bytes []byte)  {
	wsConn.outChan <- &WsMessage{websocket.BinaryMessage, bytes,}
}

func (wsConn *WsConnection) writeText(text string)  {
	data := []byte(text)
	wsConn.outChan <- &WsMessage{websocket.TextMessage, data,}
}


func (wsConn *WsConnection) wsClose() {
	wsConn.wsSocket.Close()
	wsConn.mutex.Lock()
	defer wsConn.mutex.Unlock()
	if !wsConn.isClosed {
		if wsConn.onClose != nil{
			wsConn.onClose(wsConn)
		}
		wsConn.isClosed = true
	}

}


//设置链接属性
func (c *WsConnection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取链接属性
func (c *WsConnection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *WsConnection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}