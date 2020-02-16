package proto

/*
服与服之间通讯基础消息
*/

type ServerType int
const (
	ServerTypeMaster = iota
	ServerTypeLogin
	ServerTypeWorld
	ServerTypeGate
	ServerTypeGame
)

const (
	ServerStateNormal= iota
	ServerStateDead
)

const (
	SessionOpDelete = iota
	SessionOpKeepLive
)

const (
	UserOnline = iota
	UserOffline
)

/*
上报服务器信息
 */
type ServerInfo struct{
	Name string
	Id   string
	IP   string
	Port  int
	Type ServerType
	OnlineCnt uint32
	State int16
	LastTime int64
	ProxyName string
}

type ServerInfoReport struct {
	ServerInfo
}

/*
服务器之间的心跳信息
*/
type PingPong struct {
	CurTime int64
}

type ServerListReq struct{
	CurTime int64
}

type ServerListAck struct{
	BaseAck
	ServerMap map[string] ServerInfo
}

type CheckSessionReq struct{
	UserId  uint32
	ConnId	uint32
	Session string
}

type CheckSessionAck struct{
	BaseAck
	UserId  uint32
	ConnId	uint32
	Session string
}

type SessionUpdateReq struct{
	UserId  uint32
	ConnId	uint32
	Session string
	OpType  int8
}

type SessionUpdateAck struct{
	BaseAck
	UserId  uint32
	ConnId	uint32
	Session string
	OpType  int8
}

type UserOnlineOrOffLineReq struct {
	Type    int
	UserId  uint32
}

type UserOnlineOrOffLineAck struct {
	Type    int
	UserId  uint32
}
/*
进入游戏
*/
type EnterGameReq struct{
	UserId		uint32
}

type EnterGameAck struct{
	BaseAck
}