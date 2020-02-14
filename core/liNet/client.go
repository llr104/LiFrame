package liNet

import (
	"fmt"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/utils"
	"net"
)

type Client struct {
	name        string
	id          string
	ipVersion   string
	remoteIP    string
	remotePort  int
	clientType  proto.ServerType
	Server      *Server
	msgHandler  liFace.IMsgHandle
	connMgr     liFace.IConnManager
	onConnStart func(conn liFace.IConnection)
	onConnStop  func(conn liFace.IConnection)
}

func NewClient (clientName string, clientId string, remoteIP string, remotePort int) *Client {

	c:= Client {
		name:       clientName,
		id:         clientId,
		ipVersion:  "tcp4",
		remoteIP:   remoteIP,
		remotePort: remotePort,
		msgHandler: NewMsgHandle(1),
		connMgr:    NewConnManager(),
	}
	return &c
}

func (c *Client) SetClientType(sType proto.ServerType){
	c.clientType = sType
}

func (c *Client) GetClientType() proto.ServerType{
	return c.clientType
}


func (c *Client) GetName() string{
	return c.name
}

func (c *Client) GetId() string{
	return c.id
}

func (c *Client) Start(){
	c.msgHandler.StartWorkerPool()
}

func (c *Client) GetHost()string{
	return c.remoteIP
}

func (c *Client) GetPort() int{
	return c.remotePort
}

func (c *Client) Stop(){

	utils.Log.Info("[STOP] Client")
	c.msgHandler.StopWorkerPool()
	c.connMgr.ClearConn()

}


func (c *Client) Running(){
	c.Start()

	addr, err := net.ResolveTCPAddr(c.ipVersion, fmt.Sprintf("%s:%d", c.GetHost(), c.GetPort()))
	if err != nil {
		utils.Log.Warning("resolve tcp addr err:%s", err.Error())
		return
	}

	connTcp, err := net.DialTCP(c.ipVersion,nil,addr)
	if err != nil {
		utils.Log.Warning("app start exit err:%s",err.Error())
		conn := NewConnection(c, nil, 0, c.msgHandler)
		c.CallOnConnStop(conn)
		return
	}

	//保证client的时候只有一个conn
	c.connMgr.ClearConn()
	conn := NewConnection(c, connTcp, 0, c.msgHandler)
	conn.Start()
}

func (c *Client) GetConnMgr() liFace.IConnManager {
	return c.connMgr
}


func (c *Client) GetConn() liFace.IConnection {
	conn, err := c.connMgr.Get(0)
	if err == nil {
		return conn
	}else{
		return  nil
	}
}


func (c *Client)AddRouter(router liFace.IRouter){
	c.msgHandler.AddRouter(router)
}


func (c *Client)SetOnConnStart(hookFunc func (liFace.IConnection)){
	c.onConnStart = hookFunc
}

func (c *Client)SetOnConnStop(hookFunc func (liFace.IConnection)){
	c.onConnStop = hookFunc
}

func (c *Client)CallOnConnStart(conn liFace.IConnection){
	if c.onConnStart != nil {
		utils.Log.Info("---> CallOnConnStart....")
		c.onConnStart(conn)
	}
}

func (c *Client)CallOnConnStop(conn liFace.IConnection){
	if c.onConnStop != nil {
		utils.Log.Info("---> CallOnConnStop....")
		c.onConnStop(conn)
	}
}


