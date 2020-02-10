package world

import (
	"LiFrame/core/liFace"
	"LiFrame/core/liNet"
	"LiFrame/proto"
	"LiFrame/server/app"
	"LiFrame/utils"
	"encoding/json"
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
	utils.Log.Info("JoinWorldReq begin: %s", req.GetMsgName())

	reqInfo := proto.JoinWorldReq{}
	ackInfo := proto.JoinWorldAck{}
	err := json.Unmarshal(req.GetData(), &reqInfo)
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
				c.GetConn().SendMsg(proto.EnterLoginCheckSessionReq, data)
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

	utils.Log.Info("JoinWorldReq end: %v", reqInfo)
}


func (s *EnterWorld) CheckSessionAck(req liFace.IRequest) {
	utils.Log.Info("CheckSessionAck begin: %s", req.GetMsgName())

	reqInfo := proto.CheckSessionAck{}
	err := json.Unmarshal(req.GetData(), &reqInfo)
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
				conn.SetProperty("session", reqInfo.Session)
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

	utils.Log.Info("CheckSessionAck end: %v", reqInfo)
}
