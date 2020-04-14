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

type WsMessageReq struct {
	MsgType int
	Seq     int
	Data    []byte
}

type WsMessageRsp struct {
	MsgType 	int
	Seq     	int
	FuncName 	string
	ProxyName   string
	Data    	interface{}
}


// 客户端连接
type WsConnection struct {
	wsSocket 	*websocket.Conn   	// 底层websocket
	outChan 	chan *WsMessageReq 	// 写队列
	mutex 		sync.Mutex       	// 避免重复关闭管道
	isClosed 	bool
	id       	uint64              // id
	onMessage	func(wsConn *WsConnection, req *WsMessageReq, rsp*WsMessageRsp)
	onClose     func(wsConn* WsConnection)
	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

func NewWsConnection(wsSocket *websocket.Conn, cid uint64) *WsConnection {
	wsConn := &WsConnection{
		wsSocket: wsSocket,
		outChan: make(chan *WsMessageReq, 1000),
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

		req := &WsMessageReq{
			msgType,
			1,
			data,
		}

		rsp := &WsMessageRsp{
			msgType,
			1,
			"",
			"",
			data,
		}

		if wsConn.onMessage != nil{
			wsConn.onMessage(wsConn, req, rsp)
			wsConn.WriteObject(rsp.ProxyName, rsp.FuncName, rsp.Seq, rsp.Data)
		}
	}

	wsConn.Close()
}

func (wsConn *WsConnection) SetOnMessage(hookFunc func (*WsConnection, *WsMessageReq, *WsMessageRsp))  {
	wsConn.onMessage = hookFunc
}

func (wsConn *WsConnection) SetOnClose(hookFunc func (*WsConnection))  {
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
					wsConn.Close()
					return
				}
			}
	}

}

func (wsConn *WsConnection) WriteProxyMessage(proxyName string, funcName string, seq int, body interface{})  {
	data, err := json.Marshal(body)
	if err != nil{
		return
	}
	wsConn.WriteMessage(proxyName, funcName, seq, data)
}


func (wsConn *WsConnection) WriteObject(proxyName string, funcName string, seq int, body interface{})  {
	data, err := json.Marshal(body)
	if err != nil{
		return
	}

	wsConn.WriteMessage(proxyName, funcName, seq, data)
}

func (wsConn *WsConnection) WriteMessage(proxyName string, funcName string, seq int, body[] byte){
	text := fmt.Sprintf("%s|%s|%d|%s", funcName, proxyName, seq, body)

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

	wsConn.writeBytes(b.Bytes(), seq)
}

func (wsConn *WsConnection) writeBytes(bytes []byte, seq int)  {
	wsConn.outChan <- &WsMessageReq{websocket.BinaryMessage,seq,bytes,}
}

func (wsConn *WsConnection) writeText(text string, seq int)  {
	data := []byte(text)
	wsConn.outChan <- &WsMessageReq{websocket.TextMessage,seq,data,}
}


func (wsConn *WsConnection) Close() {
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