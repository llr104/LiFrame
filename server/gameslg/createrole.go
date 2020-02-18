package gameslg

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/server/db/slgdb"
	"github.com/llr104/LiFrame/server/gameslg/proto"
	"github.com/llr104/LiFrame/utils"
)

var CreateRole createRole

func init() {
	CreateRole = createRole{}
}
type createRole struct {
	liNet.BaseRouter
}

func (s *createRole) NameSpace() string {
	return "birth"
}

func (s *createRole) QryRoleReq(req liFace.IRequest)  {
	ackInfo := proto.QryRoleAck{}
	p, err := req.GetConnection().GetProperty("userId")
	if err != nil{
		ackInfo.Code = proto.Code_Not_Auth
	}else{
		userId := p.(uint32)
		r := slgdb.NewDefaultRole()
		r.UserId = userId
		if err := slgdb.FindRoleByUserId(&r); err == nil{
			ackInfo.Role = r
			ackInfo.Code = proto.Code_SLG_Success
		}else{
			ackInfo.Code = proto.Code_Role_NoFound
		}
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.BirthQryRoleAck, data)
}

func (s *createRole) NewRoleReq(req liFace.IRequest) {
	reqInfo := proto.NewRoleReq{}
	ackInfo := proto.NewRoleAck{}
	json.Unmarshal(req.GetData(), &reqInfo)
	p, err := req.GetConnection().GetProperty("userId")

	if err != nil{
		ackInfo.Code = proto.Code_Not_Auth
	}else{
		//创建角色
		userId := p.(uint32)
		r := slgdb.NewDefaultRole()
		r.Name = reqInfo.Name
		r.UserId = userId
		r.Nation = reqInfo.Nation
		if err := slgdb.FindRoleByName(&r); err == nil{
			ackInfo.Code = proto.Code_Role_Exit
		}else{

			if id, err := slgdb.InsertRoleToDB(&r); err == nil {
				ackInfo.Role = r
				ackInfo.Code = proto.Code_SLG_Success
				utils.Log.Info("new role:%d", id)
			}else {
				utils.Log.Info("new role error: %s", err.Error())
			}

		}
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(proto.BirthNewRoleAck, data)
}
