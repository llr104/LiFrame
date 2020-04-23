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

func (s* sts) Ping(req liFace.IRequest, rsp liFace.IRespond){
	utils.Log.Info("Ping")
	info := proto.PingPong{}
	info.CurTime = time.Now().Unix()
	data, _ := json.Marshal(info)
	rsp.GetMessage().SetBody(data)
}


func (s *sts) CheckSessionAck(rsp liFace.IRespond) {

	req := rsp.GetRequest()
	msg := req.GetMessage()
	reqInfo := proto.CheckSessionAck{}
	err := json.Unmarshal(msg.GetBody(), &reqInfo)
	utils.Log.Info("CheckSessionAck: %v", reqInfo)
	if err != nil{
		utils.Log.Info("CheckSessionAck req error:", err.Error())
		ackInfo := proto.SessionAck{}
		ackInfo.Code = proto.Code_Illegal
		ackInfo.UserId = reqInfo.UserId
		ackInfo.Session = reqInfo.Session
		data, _ := json.Marshal(ackInfo)
		req.GetConnection().RpcCall(proto.EnterWorldSession, data,nil)

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

			ackInfo := proto.SessionAck{}
			ackInfo.Code = reqInfo.Code
			ackInfo.UserId = reqInfo.UserId
			ackInfo.Session = reqInfo.Session
			data, _ := json.Marshal(ackInfo)
			conn.RpcCall(proto.EnterWorldSession, data, nil)
		}
	}
}

func (s* sts) UserOnOrOffReq(req liFace.IRequest, rsp liFace.IRespond) {
	msg := req.GetMessage()
	reqInfo := proto.UserOnlineOrOffLineReq{}
	json.Unmarshal(msg.GetBody(), &reqInfo)

	utils.Log.Info("UserOnOrOffReq: %v", reqInfo)

	ackInfo := proto.UserOnlineOrOffLineAck{}
	ackInfo.Type = reqInfo.Type
	ackInfo.UserId = reqInfo.UserId

	data, _ := json.Marshal(ackInfo)
	rsp.GetMessage().SetBody(data)

}