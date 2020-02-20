package gameslg

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/gameutils"
	"github.com/llr104/LiFrame/utils"
)

var game mainLogic

type mainLogic struct {

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

var Enter enterGame

func init() {
	Enter = enterGame{}
	game = mainLogic{}
	gameutils.STS.SetGame(&game)
}


type enterGame struct {
	liNet.BaseRouter
}

func (s *enterGame) NameSpace() string {
	return "*.*"
}

func (s *enterGame) EveryThingHandle(req liFace.IRequest) {
	if req.GetMsgName() == proto.GameEnterGameReq{
		reqInfo := proto.EnterGameReq{}
		ackInfo := proto.EnterGameAck{}
		if err := json.Unmarshal(req.GetData(), &reqInfo); err != nil {
			ackInfo.Code = proto.Code_Illegal
			utils.Log.Info("GameEnterGameReq error:%s", err.Error())
		} else {
			ackInfo.Code = proto.Code_Success
			req.GetConnection().SetProperty("userId", reqInfo.UserId)
		}
		data, _ := json.Marshal(ackInfo)
		req.GetConnection().SendMsg(proto.GameEnterGameAck, data)
	}
}

