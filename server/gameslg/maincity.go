package gameslg

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/server/gameslg/slgproto"
	"github.com/llr104/LiFrame/utils"
)

var MainCity mainCity

func init() {
	MainCity = mainCity{}
}


type mainCity struct {
	liNet.BaseRouter
}

func (s *mainCity) NameSpace() string {
	return "mainCity"
}

func (s *mainCity) PreHandle(req liFace.IRequest) bool{
	_, err := req.GetConnection().GetProperty("roleId")
	if err == nil {
		return true
	}else{
		utils.Log.Warning("%s not has roleId", req.GetMsgName())
		return false
	}
}

func (s *mainCity) QryBuildingQeq(req liFace.IRequest) {
	reqInfo := slgproto.QryBuildingQeq{}
	ackInfo := slgproto.QryBuildingAck{}
	json.Unmarshal(req.GetData(), &reqInfo)
	ackInfo.Code = slgproto.CodeSlgSuccess

	p, _ := req.GetConnection().GetProperty("roleId")
	roleId := p.(uint32)
	buildings := playerMgr.getBuilding(roleId, reqInfo.BuildType)
	if buildings != nil{
		data, _ := json.Marshal(buildings)
		ackInfo.BuildType = reqInfo.BuildType
		ackInfo.Buildings = string(data)
		ackInfo.Yield = playerMgr.getYield(roleId, reqInfo.BuildType)
	}
	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(slgproto.MainCityQryBuildingAck, data)
}

func (s *mainCity) UpBuildingQeq(req liFace.IRequest) {
	reqInfo := slgproto.UpBuildingQeq{}
	ackInfo := slgproto.UpBuildingAck{}
	json.Unmarshal(req.GetData(), &reqInfo)

	p, _ := req.GetConnection().GetProperty("roleId")
	roleId := p.(uint32)
	b, ok := playerMgr.upBuilding(roleId, reqInfo.BuildId, reqInfo.BuildType)
	if ok {
		ackInfo.Code = slgproto.CodeSlgSuccess
		ackInfo.BuildType = reqInfo.BuildType
		ackInfo.Yield = playerMgr.getYield(roleId, reqInfo.BuildType)
		data, _ := json.Marshal(b)
		ackInfo.Build = string(data)
		if role := playerMgr.getRole(roleId);role != nil{
			ackInfo.Role = *role
		}
	}else{
		ackInfo.Code = slgproto.CodeBuildingUpError
	}
	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(slgproto.MainCityUpBuildingAck, data)
}


func (s *mainCity) QryGeneralReq(req liFace.IRequest) {
	reqInfo := slgproto.QryGeneralReq{}
	ackInfo := slgproto.QryGeneralAck{}
	json.Unmarshal(req.GetData(), &reqInfo)

	p, _ := req.GetConnection().GetProperty("roleId")
	roleId := p.(uint32)
	b, ok := playerMgr.getGenerals(roleId)
	if ok {
		ackInfo.Code = slgproto.CodeSlgSuccess
		ackInfo.Generals = b
	}else{
		ackInfo.Code = slgproto.CodeGeneralError
	}
	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(slgproto.MainCityQryGeneralAck, data)
}

