package game

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
	utils.Log.Info("ClientConnStart:%s", conn.RemoteAddr().String())
}

func ClientConnStop(conn liFace.IConnection) {
	app.MClientData.Dec()

	//修改离线用户
	if p, err := conn.GetProperty("userId");err == nil{
		userId := p.(uint32)
		STS.userOffline(userId)
	}

	utils.Log.Info("ClientConnStop:%s", conn.RemoteAddr().String())
}

func ShutDown(){
	utils.Log.Info("ShutDown")
	game.shutDown()
}

var STS sts

func init() {
	STS = sts{}
}

type sts struct {
	liNet.BaseRouter
}

func (s *sts) NameSpace() string {
	return "System"
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
	ok, state := GUserMgr.UserIsIn(userId)
	if ok {
		GUserMgr.UserChangeState(userId, GUserStateOffLine, state.SceneId, nil)
		r := game.userOffLine(userId)
		if r {
			GUserMgr.UserChangeState(userId, GUserStateLeave, -1,nil)
		}
	}
}

func (s* sts) userOnline(userId uint32){

}