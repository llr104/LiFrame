package proto

import "LiFrame/dbobject"

type UserInfoReq struct {
	UserId		uint32
}

type UserInfoAck struct {
	BaseAck
	User dbobject.User
}

type UserLogoutReq struct {
	UserId		uint32
}

type UserLogoutAck struct {
	BaseAck
}

type GameServersInfo struct {
	Id        string
	Name      string
	ProxyName string
}

type GameServersReq struct {
	UserId uint32
}

type GameServersAck struct {
	BaseAck
	Servers map[string]GameServersInfo
}