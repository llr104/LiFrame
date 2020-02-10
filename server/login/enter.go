package login

import (
	"LiFrame/core/liFace"
	"LiFrame/core/liNet"
	"LiFrame/dbobject"
	"LiFrame/proto"
	"LiFrame/server/app"
	"LiFrame/utils"
	"encoding/json"
	"time"
)

var Enter EnterLogin

func init() {
	Enter = EnterLogin{}
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

type EnterLogin struct {
	liNet.BaseRouter
}

func (s *EnterLogin) NameSpace() string {
	return "EnterLogin"
}

/*
登录
*/
func (s* EnterLogin) Ping(req liFace.IRequest){
	utils.Log.Info("Ping req: %s", req.GetMsgName())
	info := proto.PingPong{}
	info.CurTime = time.Now().Unix()
	data, _ := json.Marshal(info)
	req.GetConnection().SendMsg(proto.LoginClientPong, data)
}

func (s *EnterLogin) LoginReq(req liFace.IRequest) {
	beginTime := time.Now().Nanosecond()

	utils.Log.Info("LoginReq begin: %s", req.GetMsgName())
	reqInfo := proto.LoginReq{}
	ackInfo := proto.LoginAck{}

	err := json.Unmarshal(req.GetData(), &reqInfo)
	if err != nil {
		ackInfo.Code = proto.Code_Illegal
		utils.Log.Info("LoginReq error:%s", err.Error())
	} else {
		user := dbobject.User{}
		user.Name = reqInfo.Name
		user.Password = reqInfo.Password
		user.LastLoginIp = reqInfo.Ip

		if err := dbobject.FindUserByNP(&user); err != nil {
			ackInfo.Code = proto.Code_User_Error
			utils.Log.Info("LoginReq FindByNamePassword error:%s", err.Error())
		} else {

			if user.State != dbobject.UserStateNormal {
				ackInfo.Code = proto.Code_User_Forbid
			} else {
				user.LastLoginTime = time.Now().Unix()
				user.IsOnline = true
				user.LoginTimes += 1
				dbobject.UpdateUserToDB(&user)
				session := s.login(&user, req.GetConnection())

				ackInfo.Code = proto.Code_Success
				ackInfo.Password = user.Password
				ackInfo.Name = user.Name
				ackInfo.Id = user.Id
				ackInfo.Session = session
			}
		}
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.EnterLoginLoginAck, data)

	endTime := time.Now().Nanosecond()
	diff := endTime-beginTime
	u := uint64(diff) / uint64(time.Millisecond)


	utils.Log.Info("LoginReq end: %v,time:%d", reqInfo, u)
}

/*
注册
*/
func (s *EnterLogin) RegisterReq(req liFace.IRequest) {

	beginTime := time.Now().Nanosecond()

	utils.Log.Info("RegisterReq begin: %s", req.GetMsgName())
	reqInfo := proto.RegisterReq{}
	ackInfo := proto.RegisterAck{}
	err := json.Unmarshal(req.GetData(), &reqInfo)

	if err != nil {
		ackInfo.Code = proto.Code_Illegal
		utils.Log.Info("RegisterReq error:", err.Error())
	} else {
		user := dbobject.User{}
		user.Name = reqInfo.Name
		user.Password = reqInfo.Password
		user.LastLoginIp = reqInfo.Ip

		if err := dbobject.FindUserByName(&user); err == nil {
			ackInfo.Code = proto.Code_User_Exist
			utils.Log.Info("RegisterReq FindByName:%s Exist", ackInfo.Name)
		} else {
			user.LastLoginTime = time.Now().Unix()
			user.IsOnline = true
			user.LoginTimes = 1
			user.State = dbobject.UserStateNormal
			dbobject.InsertUserToDB(&user)

			ackInfo.Code = proto.Code_Success
			ackInfo.Password = user.Password
			ackInfo.Name = user.Name
			ackInfo.Id = user.Id
		}
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.EnterLoginRegisterAck, data)

	endTime := time.Now().Nanosecond()
	diff := endTime-beginTime
	u := uint64(diff) / uint64(time.Millisecond)

	utils.Log.Info("RegisterReq end: %v,%d", reqInfo,u)
}

/*
校验session
*/
func (s *EnterLogin) CheckSessionReq(req liFace.IRequest) {
	utils.Log.Info("CheckSessionReq begin: %s", req.GetMsgName())
	reqInfo := proto.CheckSessionReq{}
	ackInfo := proto.CheckSessionAck{}
	err := json.Unmarshal(req.GetData(), &reqInfo)

	if err != nil {
		ackInfo.Code = proto.Code_Illegal
		utils.Log.Info("CheckSessionReq error:", err.Error())
	} else {
		ok := LoginSessMgr.SessionIsLive(reqInfo.UserId, reqInfo.Session)
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
	req.GetConnection().SendMsg(proto.EnterWorldCheckSessionAck, data)
	utils.Log.Info("RegisterReq end: %v", reqInfo)
}

/*
根据负载分配world服务器
*/
func (s *EnterLogin) DistributeWorldReq(req liFace.IRequest) {
	utils.Log.Info("DistributeServerReq begin: %s", req.GetMsgName())
	reqInfo := proto.DistributeServerReq{}
	ackInfo := proto.DistributeServerAck{}

	if err := json.Unmarshal(req.GetData(), &reqInfo); err != nil {
		ackInfo.Code = proto.Code_Illegal
		utils.Log.Info("DistributeWorldReq error:%s", err.Error())
	} else {
		
		if serverInfo, err:= app.ServerMgr.Distribute(proto.ServerTypeWorld); err != nil {
			ackInfo.Code = proto.Code_Not_Server
			utils.Log.Info("DistributeWorldReq error:%s", err.Error())
		}else{
			ackInfo.Code = proto.Code_Success
			ackInfo.ServerInfo = serverInfo
		}
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.EnterLoginDistributeWorldAck, data)
	utils.Log.Info("DistributeWorldAck end: %v", reqInfo)
}

/*
更新session操作
*/
func (s *EnterLogin) SessionUpdateReq(req liFace.IRequest) {
	utils.Log.Info("SessionUpdateReq begin: %s", req.GetMsgName())
	reqInfo := proto.SessionUpdateReq{}
	ackInfo := proto.SessionUpdateAck{}

	ackInfo.Session = reqInfo.Session
	ackInfo.UserId = reqInfo.UserId
	ackInfo.ConnId = reqInfo.ConnId
	ackInfo.OpType = reqInfo.OpType

	if err := json.Unmarshal(req.GetData(), &reqInfo); err != nil {
		ackInfo.Code = proto.Code_Illegal
		utils.Log.Info("SessionUpdateReq error:%s", err.Error())
	} else {
		if reqInfo.OpType == proto.SessionOpDelete {
			s.logout(reqInfo.UserId, reqInfo.Session)
		}else if reqInfo.OpType == proto.SessionOpKeepLive {
			LoginSessMgr.SessionKeepLive(reqInfo.UserId, reqInfo.Session)
		}
		ackInfo.Code = proto.Code_Success
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.EnterLoginSessionUpdateAck, data)
	utils.Log.Info("SessionUpdateReq end: %v", reqInfo)
}


func (s *EnterLogin) login(user *dbobject.User, conn liFace.IConnection) string{

	ser := app.GetServer()
	n := ser.(liFace.INetWork)

	session := app.SessionMgr.CreateSession(n.GetId(), user.Id)
	LoginSessMgr.AddSession(user.Id, session)
	conn.SetProperty("session",session)
	conn.SetProperty("userId", user.Id)
	return session
}


func (s *EnterLogin) logout(userId uint32, session string) {

	LoginSessMgr.RemoveSession(userId, session)
	utils.Log.Info("logout userId:%d", userId)
}