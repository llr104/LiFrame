package liFace

/*
	IRequest 接口：
	实际上是把客户端请求的链接信息 和 请求的数据 包装到了 Request里
*/
type IRequest interface{
	GetConnection() IConnection //获取请求连接信息
	GetMessage() IMessage
	SetMessage(msg IMessage)
}

/*
	IRespond 接口：
*/

type IRespond interface{
	GetData() []byte            //获取请求消息的数据
	GetRequest() IRequest
	SetRequest(req IRequest)
	GetMessage() IMessage
	SetMessage(msg IMessage)
}
