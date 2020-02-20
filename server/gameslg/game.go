package gameslg

import (
	"github.com/llr104/LiFrame/server/gameutils"
	"github.com/llr104/LiFrame/utils"
)

var game mainLogic

func init() {
	game = mainLogic{}
	gameutils.STS.SetGame(&game)
}

type mainLogic struct {
	isShutDown bool
}

/*
返回true:用户离开了游戏，返回false:用户断线，保留用户的游戏状态
*/
func (s *mainLogic) UserOffLine(userId uint32) bool{
	playerMgr.ReleasePlayer(userId)
	return true
}

func (s *mainLogic) UserOnLine(userId uint32){
	utils.Log.Info("UserOnLine: %d", userId)
}

func (s *mainLogic) UserLogout(userId uint32) bool{
	return s.UserOffLine(userId)
}


func (s *mainLogic) ShutDown(){
	utils.Log.Info("ShutDown")
}