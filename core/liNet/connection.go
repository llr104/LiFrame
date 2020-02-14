package liNet

import (
	"errors"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/utils"
	"github.com/thinkoner/openssl"
	"io"
	"net"
	"sync"
)

type Connection struct {
	//当前Conn属于哪个Server
	TcpNetWork liFace.INetWork
	//当前连接的socket TCP套接字
	Conn *net.TCPConn
	//当前连接的ID 也可以称作为SessionID，ID全局唯一
	ConnID uint32
	//当前连接的关闭状态
	isClosed bool
	//消息管理msgName和对应处理方法的消息管理模块
	MsgHandler liFace.IMsgHandle
	//告知该链接已经退出/停止的channel
	ExitBuffChan chan bool
	//无缓冲管道，用于读、写两个goroutine之间的消息通信
	msgChan chan []byte
	//有关冲管道，用于读、写两个goroutine之间的消息通信
	msgBuffChan chan []byte

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

//创建连接的方法
func NewConnection(netWork liFace.INetWork, conn *net.TCPConn, connID uint32, msgHandler liFace.IMsgHandle) *Connection {
	//初始化Conn属性
	c := &Connection{
		TcpNetWork:   netWork,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
		msgBuffChan:  make(chan []byte, utils.GlobalObject.MaxMsgChanLen),
		property:     make(map[string]interface{}),
	}

	if conn != nil{
		//将新创建的Conn添加到链接管理中
		c.TcpNetWork.GetConnMgr().Add(c)
	}else{
		c.isClosed = true
	}


	return c
}


/*
	写消息Goroutine， 用户将数据发送给客户端
*/
func (c *Connection) StartWriter() {
	utils.Log.Info ("%s [Writer Goroutine is running]", c.TcpNetWork.GetName())
	defer utils.Log.Info("%s %s [conn Writer exit!]",c.TcpNetWork.GetName(), c.RemoteAddr().String())

	for {
		select {
			case data := <-c.msgChan:
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					utils.Log.Error("%s Send Body error %s:, Conn Writer exit ",c.TcpNetWork.GetName(), err.Error())
					return
				}

			case data := <-c.msgBuffChan:
				//有数据要写给客户端
				if _, err := c.Conn.Write(data); err != nil {
					utils.Log.Error("%s Send Buff Body error: %s, Conn Writer exit",c.TcpNetWork.GetName(), err.Error())
					return
				}
			case <-c.ExitBuffChan:
				return
		}
	}
}

/*
	读消息Goroutine，用于从客户端中读取数据
*/
func (c *Connection) StartReader() {
	utils.Log.Info("%s [Reader Goroutine is running]",c.TcpNetWork.GetName())
	defer utils.Log.Info("%s %s [conn Reader exit!]", c.TcpNetWork.GetName(), c.RemoteAddr().String())
	//将链接从连接管理器中删除
	defer c.TcpNetWork.GetConnMgr().Remove(c)
	defer c.Stop()

	for {

		// 创建拆包解包的对象
		dp := NewDataPack()
		//读取客户端的Msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.Conn, headData); err != nil {
			utils.Log.Warning("%s read msg head error %s", c.TcpNetWork.GetName(), err.Error())
			break
		}
		//普通socket拆包，得到nameLen 和 dataLen 放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			utils.Log.Warning("%s unpack error %s", c.TcpNetWork.GetName(), err.Error())
			break
		}

		//根据 nameLen bodyLen 读取 Data
		var data []byte
		if msg.GetBodyLen() + msg.GetNameLen()> 0 {
			data = make([]byte, msg.GetBodyLen()+msg.GetNameLen())
			if _, err := io.ReadFull(c.Conn, data); err != nil {
				utils.Log.Warning( "%s read msg Data error %s", c.TcpNetWork.GetName(), err)
				break
			}
		}

		//解密body
		body := data[msg.GetNameLen():]
		if dst, err := openssl.AesECBDecrypt(body, DataPackKey, openssl.PKCS7_PADDING); err != nil {
			return
		}else{
			body = dst
		}

		msg.SetMsgNameByte(data[0:msg.GetNameLen()])
		msg.SetBody(body)

		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}

		if utils.GlobalObject.ServerWorkerSize > 0 {
			//已经启动工作池机制，将消息交给Worker处理
			c.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			//从绑定好的消息和对应的处理方法中执行对应的Handle方法
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}


//启动连接，让当前连接开始工作
func (c *Connection) Start() {
	if c.Conn != nil{
		//1 开启用户从客户端读取数据流程的Goroutine
		go c.StartReader()
		//2 开启用于写回客户端数据流程的Goroutine
		go c.StartWriter()
		//按照用户传递进来的创建连接时需要处理的业务，执行钩子方法
		c.TcpNetWork.CallOnConnStart(c)
	}
}

//停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	utils.Log.Info("%s Conn Stop()...ConnID = %d",c.TcpNetWork.GetName(), c.ConnID)
	//如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true

	//如果用户注册了该链接的关闭回调业务，那么在此刻应该显示调用
	c.TcpNetWork.CallOnConnStop(c)

	// 关闭socket链接
	c.Conn.Close()
	//关闭Writer
	c.ExitBuffChan <- true

	//关闭该链接全部管道
	close(c.ExitBuffChan)
	close(c.msgBuffChan)
}

func (c *Connection) GetTcpNetWork() liFace.INetWork {
	return c.TcpNetWork
}

//从当前连接获取原始的socket TCPConn
func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

//获取当前连接ID
func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) IsClose() bool {
	return c.isClosed
}

//获取远程客户端地址信息
func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}


//直接将Message数据发送数据给远程的TCP客户端
func (c *Connection) SendMsg(msgName string, data []byte) error {
	if c.isClosed == true {
		utils.Log.Warning("connection closed when send msg")
		return errors.New("connection closed when send msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgName, data))
	if err != nil {
		utils.Log.Warning("%s Pack error msg id = %s", c.TcpNetWork.GetName(), msgName)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgChan <- msg

	return nil
}

func (c *Connection) SendBuffMsg(msgName string, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send buff msg")
	}
	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgName, data))
	if err != nil {
		utils.Log.Warning("%s Pack error msg id = %s",c.TcpNetWork.GetName(),  msgName)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	c.msgBuffChan <- msg

	return nil
}

//设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	c.property[key] = value
}

//获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	if value, ok := c.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

//移除链接属性
func (c *Connection) RemoveProperty(key string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()

	delete(c.property, key)
}
