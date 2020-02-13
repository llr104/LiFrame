package world

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/dbobject"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/utils"
)

var CommonWorld Common

type Common struct {
	liNet.BaseRouter
}

func init() {
	CommonWorld = Common{}
}

func (s *Common) PreHandle(req liFace.IRequest) bool{
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

func (s *Common) NameSpace() string {
	return "CommonWorld"
}

func (s *Common) UserInfoReq(req liFace.IRequest) {
	utils.Log.Info("UserInfoReq begin: %s", req.GetMsgName())
	reqInfo := proto.UserInfoReq{}
	ackInfo := proto.UserInfoAck{}

	err := json.Unmarshal(req.GetData(), &reqInfo)
	if err != nil{
		utils.Log.Info("UserInfoReq req error:",err.Error())
		ackInfo.Code = proto.Code_Illegal
		data, _ := json.Marshal(ackInfo)
		req.GetConnection().SendMsg(proto.CommonWorldUserInfoAck, data)
	}else{
		if u, e := req.GetConnection().GetProperty("userId"); e != nil{
			ackInfo.Code = proto.Code_Illegal
			data, _ := json.Marshal(ackInfo)
			req.GetConnection().SendMsg(proto.CommonWorldUserInfoAck, data)
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
			req.GetConnection().SendMsg(proto.CommonWorldUserInfoAck, data)
		}
	}

	utils.Log.Info("UserInfoReq end: %v", reqInfo)
}

func (s *Common) UserLogoutReq(req liFace.IRequest) {
	utils.Log.Info("UserLogoutReq begin: %s", req.GetMsgName())
	reqInfo := proto.UserLogoutReq{}
	ackInfo := proto.UserLogoutAck{}


	ackInfo.Code = proto.Code_Success
	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.CommonWorldUserLogoutAck, data)


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
				client.GetConn().SendMsg(proto.EnterLoginSessionUpdateReq, data)
			}
		}
	}
	//上报到登录服end

	//退出online并关闭当前连接
	conn := req.GetConnection()
	c := conn.(*liNet.Connection)
	OnlineInstance.Exit(c)

	req.GetConnection().Stop()
	utils.Log.Info("UserLogoutReq end: %v", reqInfo)
}

func (s *Common) GameServersReq(req liFace.IRequest) {
	utils.Log.Info("GameServersReq begin: %s", req.GetMsgName())
	reqInfo := proto.GameServersReq{}
	ackInfo := proto.GameServersAck{}

	m := app.ServerMgr.GetGameScenesMap()
	ackInfo.Servers = m
	ackInfo.Code = proto.Code_Success
	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.CommonWorldGameServersAck, data)
	utils.Log.Info("GameServersReq end: %v", reqInfo)
}