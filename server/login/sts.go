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

func (s* sts) Ping(req liFace.IRequest, rsp liFace.IMessage){
	utils.Log.Info("Ping")
	info := proto.PingPong{}
	info.CurTime = time.Now().Unix()
	data, _ := json.Marshal(info)
	rsp.SetBody(data)
}


/*
校验session
*/
func (s *sts) CheckSessionReq(req liFace.IRequest, rsp liFace.IMessage) {
	msg := req.GetMessage()
	reqInfo := proto.CheckSessionReq{}
	ackInfo := proto.CheckSessionAck{}
	err := json.Unmarshal(msg.GetBody(), &reqInfo)
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
	rsp.SetBody(data)

}

/*
更新session操作
*/
func (s *sts) SessionUpdateReq(req liFace.IRequest, rsp liFace.IMessage) {

	reqInfo := proto.SessionUpdateReq{}
	ackInfo := proto.SessionUpdateAck{}

	ackInfo.Session = reqInfo.Session
	ackInfo.UserId = reqInfo.UserId
	ackInfo.ConnId = reqInfo.ConnId
	ackInfo.OpType = reqInfo.OpType
	utils.Log.Info("SessionUpdateReq: %v", reqInfo)
	if err := json.Unmarshal(req.GetMessage().GetBody(), &reqInfo); err != nil {
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
	rsp.SetBody(data)

}

func (s* sts) UserOnOrOffReq(req liFace.IRequest, rsp liFace.IMessage) {

	reqInfo := proto.UserOnlineOrOffLineReq{}
	json.Unmarshal(req.GetMessage().GetBody(), &reqInfo)

	utils.Log.Info("UserOnOrOffReq: %v", reqInfo)

	ackInfo := proto.UserOnlineOrOffLineAck{}
	ackInfo.Type = reqInfo.Type
	ackInfo.UserId = reqInfo.UserId

	data, _ := json.Marshal(ackInfo)
	rsp.SetBody(data)

}