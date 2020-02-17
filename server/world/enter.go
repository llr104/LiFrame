package world

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db/dbobject"
	"github.com/llr104/LiFrame/utils"
)


var Enter enterWorld

type enterWorld struct {
	liNet.BaseRouter
}

func init() {
	Enter = enterWorld{}
}

func (s *enterWorld) PreHandle(req liFace.IRequest) bool{
	name := req.GetMsgName()
	if name == proto.EnterWorldJoinWorldReq{
		return true
	}

	_, err := req.GetConnection().GetProperty("session")
	if err == nil {
		return true
	}else{

		//session检验不通过，跳过后面的逻辑
		ackInfo := proto.JoinWorldAck{}
		ackInfo.Code = proto.Code_Session_Error
		data, _ := json.Marshal(ackInfo)
		req.GetConnection().SendMsg(proto.EnterWorldJoinWorldAck, data)
		return false
	}
}

func (s *enterWorld) NameSpace() string {
	return "EnterWorld"
}

func (s *enterWorld) JoinWorldReq(req liFace.IRequest) {

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
					conn.SendMsg(proto.SystemCheckSessionReq, data)
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

func (s *enterWorld) UserInfoReq(req liFace.IRequest) {
	utils.Log.Info("UserInfoReq begin: %s", req.GetMsgName())
	reqInfo := proto.UserInfoReq{}
	ackInfo := proto.UserInfoAck{}

	err := json.Unmarshal(req.GetData(), &reqInfo)
	if err != nil{
		utils.Log.Info("UserInfoReq req error:",err.Error())
		ackInfo.Code = proto.Code_Illegal
		data, _ := json.Marshal(ackInfo)
		req.GetConnection().SendMsg(proto.EnterWorldUserInfoAck, data)
	}else{
		if u, e := req.GetConnection().GetProperty("userId"); e != nil{
			ackInfo.Code = proto.Code_Illegal
			data, _ := json.Marshal(ackInfo)
			req.GetConnection().SendMsg(proto.EnterWorldUserInfoAck, data)
		}else{
			reqInfo.UserId = u.(uint32)
			user := dbobject.User{}
			user.Id = reqInfo.UserId
			if err:= dbobject.FindUserById(&user); err != nil{
				ackInfo.Code = proto.Code_User_Error
			}else{
				ackInfo.User = user
				ackInfo.Code = proto.Code_Success
			}

			data, _ := json.Marshal(ackInfo)
			req.GetConnection().SendMsg(proto.EnterWorldUserInfoAck, data)
		}
	}

	utils.Log.Info("UserInfoReq end: %v", reqInfo)
}

func (s *enterWorld) UserLogoutReq(req liFace.IRequest) {

	reqInfo := proto.UserLogoutReq{}
	ackInfo := proto.UserLogoutAck{}

	ackInfo.Code = proto.Code_Success
	data, _ := json.Marshal(ackInfo)
	utils.Log.Info("UserLogoutReq end: %v", reqInfo)
	req.GetConnection().SendMsg(proto.EnterWorldUserLogoutAck, data)


	//上报到登录服begin
	v1, _ := req.GetConnection().GetProperty("session")
	v2, _ := req.GetConnection().GetProperty("userId")
	session := v1.(string)
	userId := v2.(uint32)

	if session != "" &&  userId > 0{

		sessReq := proto.SessionUpdateReq{}
		sessReq.Session = session
		sessReq.UserId = userId
		sessReq.ConnId = req.GetConnection().GetConnID()
		sessReq.OpType = proto.SessionOpDelete

		if appId, err := app.SessionMgr.CheckSessionFrom(session); err == nil {
			client, ok := W2Login.GetLoginClient(appId)
			if ok {
				data, _ := json.Marshal(sessReq)
				conn := client.GetConn()
				if conn != nil{
					conn.SendMsg(proto.SystemSessionUpdateReq, data)
				}
			}
		}
	}
	//上报到登录服end

	//退出online并关闭当前连接
	conn := req.GetConnection()
	c := conn.(*liNet.Connection)
	OnlineInstance.Exit(c)

	req.GetConnection().Stop()

}

func (s *enterWorld) GameServersReq(req liFace.IRequest) {
	reqInfo := proto.GameServersReq{}
	ackInfo := proto.GameServersAck{}

	m := app.ServerMgr.GetGameScenesMap()
	ackInfo.Servers = m
	ackInfo.Code = proto.Code_Success
	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.EnterWorldGameServersAck, data)
	utils.Log.Info("GameServersReq: %v", reqInfo)
}