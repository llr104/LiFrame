package liNet

import "github.com/llr104/LiFrame/core/liFace"

type Request struct {
	conn liFace.IConnection //已经和客户端建立好的 链接
	msg  liFace.IMessage    //客户端请求的数据
}

//获取请求连接信息
func(r *Request) GetConnection() liFace.IConnection {
	return r.conn
}
//获取请求消息的数据
func(r *Request) GetData() []byte {
	return r.msg.GetBody()
}

//获取请求的消息的ID
func (r *Request) GetMsgName() string {
	return r.msg.GetMsgName()
}

