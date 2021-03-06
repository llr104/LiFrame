package liFace

import "net"

const (
	RpcPush = iota //0
	RpcReq
	RpcAck
)

//定义连接接口
type IConnection interface {
	//启动连接，让当前连接开始工作
	Start()
	//停止连接，结束当前连接状态M
	Stop()

	GetConnID() uint32
	IsClose() bool
	//从当前连接获取原始的socket TCPConn
	GetTCPConnection() *net.TCPConn

	GetTcpNetWork() INetWork

	//获取远程客户端地址信息
	RemoteAddr() net.Addr

	//直接将Message数据发送数据给远程的TCP客户端
	RpcCall(msgName string, data []byte, success func(rsp IRespond), fail func(rsp IRespond)) error
	RpcReply(msgName string, seq uint32, data []byte) error
	RpcPush(msgName string, data []byte) error
	CheckRpc(seq uint32, rsp IMessage) bool


	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string)(interface{}, error)
	//移除链接属性
	RemoveProperty(key string)

}


