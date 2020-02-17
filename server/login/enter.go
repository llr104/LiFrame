package login

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db/dbobject"
	"github.com/llr104/LiFrame/utils"
	"time"
)

var Enter enterLogin

func init() {
	Enter = enterLogin{}
}



type enterLogin struct {
	liNet.BaseRouter
}

func (s *enterLogin) NameSpace() string {
	return "EnterLogin"
}

/*
登录
*/
func (s *enterLogin) LoginReq(req liFace.IRequest) {
	beginTime := time.Now().Nanosecond()

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


	utils.Log.Info("LoginReq: %v,time:%d", reqInfo, u)
}

/*
注册
*/
func (s *enterLogin) RegisterReq(req liFace.IRequest) {

	reqInfo := proto.RegisterReq{}
	ackInfo := proto.RegisterAck{}
	err := json.Unmarshal(req.GetData(), &reqInfo)
	utils.Log.Info("RegisterReq end: %v", reqInfo)
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

}


/*
根据负载分配world服务器
*/
func (s *enterLogin) DistributeWorldReq(req liFace.IRequest) {

	reqInfo := proto.DistributeServerReq{}
	ackInfo := proto.DistributeServerAck{}
	utils.Log.Info("DistributeWorldAck: %v", reqInfo)
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

}


func (s *enterLogin) login(user *dbobject.User, conn liFace.IConnection) string{

	ser := app.GetServer()
	n := ser.(liFace.INetWork)
	session := SessLoginMgr.NewSession(n.GetId(), user.Id, conn)
	conn.SetProperty("userId", user.Id)
	return session
}


func (s *enterLogin) logout(userId uint32, session string) {
	SessLoginMgr.RemoveSession(userId, session)
	utils.Log.Info("logout userId:%d", userId)
}