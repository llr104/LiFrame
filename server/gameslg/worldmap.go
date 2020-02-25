package gameslg

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/server/gameslg/data"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
	"github.com/llr104/LiFrame/server/gameslg/slgproto"
	"github.com/llr104/LiFrame/utils"
)

/*
世界地图
*/
var WorldMap worldMap


type worldMap struct {
	liNet.BaseRouter

}

func (s *worldMap) NameSpace() string {
	return "worldMap"
}

func (s *worldMap) PreHandle(req liFace.IRequest) bool{
	_, err := req.GetConnection().GetProperty("roleId")
	if err == nil {
		return true
	}else{
		utils.Log.Warning("%s not has roleId", req.GetMsgName())
		return false
	}
}


func (s *worldMap) QryWorldMap(req liFace.IRequest)  {
	reqInfo := slgproto.QryWorldMapReq{}
	ackInfo := slgproto.QryWorldMapAck{}
	json.Unmarshal(req.GetData(), &reqInfo)
	ackInfo.Code = slgproto.CodeSlgSuccess

	n := data.CityMgr.Count()
	cm := data.CityMgr.CityMap()
	ackInfo.Citys = make([]slgdb.City, n)
	i := 0
	for _, v := range cm{
		ackInfo.Citys[i] = *v
		i++
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(slgproto.WorldMapQryWorldMapAck, data)
}

/*
驻守城池
*/
func (s *worldMap) GarrisonCity(req liFace.IRequest)  {
	reqInfo := slgproto.GarrisonCityReq{}
	ackInfo := slgproto.GarrisonCityAck{}
	json.Unmarshal(req.GetData(), &reqInfo)

	p, _ := req.GetConnection().GetProperty("roleId")
	roleId := p.(uint32)

	role := playerMgr.getRole(roleId)
	if role == nil{
		utils.Log.Error("playerMgr not found role")
		return
	}

	general, ok := playerMgr.getGeneral(roleId, reqInfo.GeneralId)
	ackInfo.CityId = reqInfo.CityId
	ackInfo.GeneralId = reqInfo.GeneralId

	if ok {
		if reqInfo.CityId == 0{
			//取消驻守
			ackInfo.Code = slgproto.CodeSlgSuccess
			general.CityId = reqInfo.CityId
		}else{
			//驻守
			cm := data.CityMgr.CityMap()
			city, found := cm[reqInfo.CityId]
			if found{
				if city.Nation == role.Nation{
					general.CityId = reqInfo.CityId
					ackInfo.Code = slgproto.CodeSlgSuccess
				}else{
					ackInfo.Code = slgproto.CodeNotLocalCity
				}
			}else{
				ackInfo.Code = slgproto.CodeCityError
			}
		}
	}else{
		ackInfo.Code = slgproto.CodeGeneralError
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(slgproto.WorldMapGarrisonCityAck, data)
}

/*
攻击城池
*/
func (s *worldMap) AttackCityReq(req liFace.IRequest)  {
	reqInfo := slgproto.AttackCityReq{}
	ackInfo := slgproto.AttackCityAck{}
	json.Unmarshal(req.GetData(), &reqInfo)

	p, _ := req.GetConnection().GetProperty("roleId")
	roleId := p.(uint32)

	role := playerMgr.getRole(roleId)
	if role == nil{
		utils.Log.Error("playerMgr not found role")
		return
	}

	general, ok := playerMgr.getGeneral(roleId, reqInfo.GeneralId)
	ackInfo.CityId = reqInfo.CityId

	if ok {
		cm := data.CityMgr.CityMap()
		city, found := cm[reqInfo.CityId]
		if found{
			if city.Nation != role.Nation{
				general.CityId = reqInfo.CityId
				ackInfo.Code = slgproto.CodeSlgSuccess

				//这里需要完成战斗，得出城池是否被攻占，后面加

			}else{
				ackInfo.Code = slgproto.CodeAttackLocalCity
			}
		}else{
			ackInfo.Code = slgproto.CodeCityError
		}
	}else{
		ackInfo.Code = slgproto.CodeGeneralError
	}

	ackInfo.General = *general
	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(slgproto.WorldMapAttackCityAck, data)
}