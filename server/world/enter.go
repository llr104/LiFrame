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

var Enter EnterWorld

func init() {
	Enter = EnterWorld{}
}


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

type EnterWorld struct {
	liNet.BaseRouter
}

func (s *EnterWorld) NameSpace() string {
	return "EnterWorld"
}

func (s *EnterWorld) JoinWorldReq(req liFace.IRequest) {

	reqInfo := proto.JoinWorldReq{}
	ackInfo := proto.JoinWorldAck{}
	err := json.Unmarshal(req.GetData(), &reqInfo)
	utils.Log.Info("JoinWorldReq: %v", reqInfo)
	if err != nil{
		utils.Log.Info("JoinWorldReq req error:",err.Error())
		ackInfo.Code = proto.Code_Illegal
		data, _ := json.Marshal(ackInfo)
		req.GetConnection().SendMsg(proto.EnterWorldJoinWorldAck, data)
	}else{
		//向login校验session是否有效
		if serverId, err := app.SessionMgr.CheckSessionFrom(reqInfo.Session); err == nil{
			c, ok := W2Login.clientMap[serverId]
			if ok {
				sessionReq := proto.CheckSessionReq{}
				sessionReq.Session = reqInfo.Session
				sessionReq.UserId = reqInfo.UserId
				sessionReq.ConnId = req.GetConnection().GetConnID()
				data, _ := json.Marshal(sessionReq)

				conn := c.GetConn()
				if conn != nil{
					conn.SendMsg(proto.EnterLoginCheckSessionReq, data)
				}else{
					ackInfo.Code = proto.Code_Session_Error
					data, _ := json.Marshal(ackInfo)
					req.GetConnection().SendMsg(proto.EnterWorldJoinWorldAck, data)
				}
			}else{
				utils.Log.Info("session serverId: %s not found app connect server", serverId)

				ackInfo.Code = proto.Code_Session_Error
				data, _ := json.Marshal(ackInfo)
				req.GetConnection().SendMsg(proto.EnterWorldJoinWorldAck, data)
			}
		}else{
			utils.Log.Info("session serverId: %s not found from server", serverId)

			ackInfo.Code = proto.Code_Session_Error
			data, _ := json.Marshal(ackInfo)
			req.GetConnection().SendMsg(proto.EnterWorldJoinWorldAck, data)
		}
	}


}


func (s *EnterWorld) CheckSessionAck(req liFace.IRequest) {

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
