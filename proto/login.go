package proto

/*
客户端与登录服之间消息
*/
type LoginReq struct {
	Name     string
	Password string
	Ip       string
}

type LoginAck struct {
	BaseAck
	Id			uint32
	Name     	string
	Password 	string
	Session     string
}

type RegisterReq struct {
	Name     	string
	Password 	string
	Ip       	string
}

type RegisterAck struct {
	BaseAck
	Id			  uint32
	Name     	  string
	Password 	  string
}

type DistributeServerReq struct{
	CurTime int64
}

type DistributeServerAck struct{
	BaseAck
	ServerInfo ServerInfo
}

type JoinWorldReq struct {
	Session      string
	UserId		 uint32	
}

type JoinWorldAck struct {
	BaseAck
	Session      string
	UserId		 uint32	
}