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


func (r *Request) GetMessage() liFace.IMessage {
	return r.msg
}

func (r *Request) SetMessage(msg liFace.IMessage)  {
	r.msg = msg
}

type Respond struct {
	msg  	liFace.IMessage
	req 	liFace.IRequest
}


//获取响应消息的数据
func(r *Respond) GetData() []byte {
	return r.msg.GetBody()
}



func (r *Respond) GetRequest() liFace.IRequest {
	return r.req
}

func (r *Respond) SetRequest(req liFace.IRequest)  {
	r.req = req
}

func (r *Respond) GetMessage() liFace.IMessage {
	return r.msg
}

func (r *Respond) SetMessage(msg liFace.IMessage)  {
	r.msg = msg
}



