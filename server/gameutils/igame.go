package gameutils

type IGame interface {

	/*
	返回true:用户离开了游戏，返回false:用户断线，保留用户的游戏状态
	*/
	UserOffLine(userId uint32) bool
	UserOnLine(userId uint32)
	UserLogout(userId uint32) bool
	ShutDown()
}



