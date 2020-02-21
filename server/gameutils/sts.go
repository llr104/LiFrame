package gameutils

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/utils"
)

func ClientConnStart(conn liFace.IConnection) {
	app.MClientData.Inc()
}

func ClientConnStop(conn liFace.IConnection) {
	app.MClientData.Dec()
}

func ShutDown(){
	utils.Log.Info("ShutDown")
	if STS.isShutDown{
		return
	}
	//关闭前处理
	STS.isShutDown = true
	if STS.game != nil{
		STS.game.ShutDown()
	}
}

var STS sts

func init() {
	STS = sts{}
}

type sts struct {
	liNet.BaseRouter
	game IGame
	isShutDown bool
}

func (s *sts) NameSpace() string {
	return "System"
}

func  (s *sts) SetGame(game IGame)  {
	s.game = game
}

func (s *sts) IsShutDown() bool {
	return s.isShutDown
}

func (s* sts) UserOnOrOffReq(req liFace.IRequest) {

	reqInfo := proto.UserOnlineOrOffLineReq{}
	json.Unmarshal(req.GetData(), &reqInfo)

	utils.Log.Info("UserOnOrOffReq: %v", reqInfo)

	ackInfo := proto.UserOnlineOrOffLineAck{}
	ackInfo.Type = reqInfo.Type
	ackInfo.UserId = reqInfo.UserId

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.SystemUserOnOrOffAck, data)

	if reqInfo.Type == proto.UserOffline {
		s.userOffline(reqInfo.UserId)
	}else{
		s.userOnline(reqInfo.UserId)
	}
}

func (s* sts) userOffline(userId uint32) {
	if s.game != nil{
		s.game.UserOffLine(userId)
	}
}

func (s* sts) userOnline(userId uint32){
	if s.game != nil{
		s.game.UserOnLine(userId)
	}
}