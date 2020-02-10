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

type ScenesInfo struct {
	Id        string
	Name      string
	ProxyName string
}

type GameScenesReq struct {
	UserId uint32
}

type GameScenesAck struct {
	BaseAck
	Scenes map[string]ScenesInfo
}