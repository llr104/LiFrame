package login

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/utils"
	"time"
)

func ClientConnStart(conn liFace.IConnection) {
	app.MClientData.Inc()
	utils.Log.Info("ClientConnStart:%s", conn.RemoteAddr().String())
}

func ClientConnStop(conn liFace.IConnection) {
	app.MClientData.Dec()
	SessLoginMgr.SessionExitByConn(conn)
	utils.Log.Info("ClientConnStop:%s", conn.RemoteAddr().String())
}

func ShutDown(){
	utils.Log.Info("ShutDown")
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

func (s* sts) Ping(req liFace.IRequest){
	utils.Log.Info("Ping")
	info := proto.PingPong{}
	info.CurTime = time.Now().Unix()
	data, _ := json.Marshal(info)
	req.GetConnection().SendMsg(proto.SystemPong, data)
}

func (s* sts) Pong(req liFace.IRequest){
	utils.Log.Info("Pong")
}

/*
校验session
*/
func (s *sts) CheckSessionReq(req liFace.IRequest) {

	reqInfo := proto.CheckSessionReq{}
	ackInfo := proto.CheckSessionAck{}
	err := json.Unmarshal(req.GetData(), &reqInfo)
	utils.Log.Info("CheckSessionReq: %v", reqInfo)
	if err != nil {
		ackInfo.Code = proto.Code_Illegal
		utils.Log.Info("CheckSessionReq error:", err.Error())
	} else {
		ok := SessLoginMgr.SessionIsLive(reqInfo.UserId, reqInfo.Session)
		if ok {
			ackInfo.Code = proto.Code_Success
		}else{
			ackInfo.Code = proto.Code_Session_Error
		}
	}
	ackInfo.Session = reqInfo.Session
	ackInfo.UserId = reqInfo.UserId
	ackInfo.ConnId = reqInfo.ConnId

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.SystemCheckSessionAck, data)
}

/*
更新session操作
*/
func (s *sts) SessionUpdateReq(req liFace.IRequest) {

	reqInfo := proto.SessionUpdateReq{}
	ackInfo := proto.SessionUpdateAck{}

	ackInfo.Session = reqInfo.Session
	ackInfo.UserId = reqInfo.UserId
	ackInfo.ConnId = reqInfo.ConnId
	ackInfo.OpType = reqInfo.OpType
	utils.Log.Info("SessionUpdateReq: %v", reqInfo)
	if err := json.Unmarshal(req.GetData(), &reqInfo); err != nil {
		ackInfo.Code = proto.Code_Illegal
		utils.Log.Info("SessionUpdateReq error:%s", err.Error())
	} else {
		if reqInfo.OpType == proto.SessionOpDelete {
			Enter.logout(reqInfo.UserId, reqInfo.Session)
		}else if reqInfo.OpType == proto.SessionOpKeepLive {
			SessLoginMgr.SessionKeepLive(reqInfo.UserId, reqInfo.Session)
		}
		ackInfo.Code = proto.Code_Success
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.SystemSessionUpdateAck, data)

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

}