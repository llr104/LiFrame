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
	ackInfo.Code = slgproto.Code_SLG_Success

	p, _ := req.GetConnection().GetProperty("roleId")
	roleId := p.(uint32)
	buildings := buildingMgr.getBuilding(roleId, reqInfo.BuildType)
	if buildings != nil{
		data, _ := json.Marshal(buildings)
		ackInfo.BuildType = reqInfo.BuildType
		ackInfo.Buildings = string(data)
	}
	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(slgproto.MainCityQryBuildingAck, data)
}
