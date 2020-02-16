package world

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
	app.SessionMgr.SessionExitByConn(conn)

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

func (s *sts) CheckSessionAck(req liFace.IRequest) {

	reqInfo := proto.CheckSessionAck{}
	err := json.Unmarshal(req.GetData(), &reqInfo)
	utils.Log.Info("CheckSessionAck: %v", reqInfo)
	if err != nil{
		utils.Log.Info("CheckSessionAck req error:",err.Error())
		ackInfo := proto.JoinWorldAck{}
		ackInfo.Code = proto.Code_Illegal
		ackInfo.UserId = reqInfo.UserId
		ackInfo.Session = reqInfo.Session
		data, _ := json.Marshal(ackInfo)
		req.GetConnection().SendMsg(proto.EnterWorldJoinWorldAck, data)

	}else{

		ser := app.GetServer()
		n := ser.(liFace.INetWork)
		conn, err := n.GetConnMgr().Get(reqInfo.ConnId)
		if err != nil{
			utils.Log.Info("CheckSessionAck conn: %d, error:%s", reqInfo.ConnId, err.Error())
		}else{
			//绑定session 绑定userId
			if reqInfo.Code == proto.Code_Success {

				app.SessionMgr.SessionEnter(reqInfo.Session, conn)
				conn.SetProperty("userId", reqInfo.UserId)
				conn.SetProperty("lastKeepLive", time.Now().Unix())

				c := conn.(*liNet.Connection)
				OnlineInstance.Join(c)
			}

			ackInfo := proto.JoinWorldAck{}
			ackInfo.Code = reqInfo.Code
			ackInfo.UserId = reqInfo.UserId
			ackInfo.Session = reqInfo.Session
			data, _ := json.Marshal(ackInfo)
			conn.SendMsg(proto.EnterWorldJoinWorldAck, data)
		}
	}
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