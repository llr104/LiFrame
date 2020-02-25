package slgproto

import (
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
)

type QryWorldMapReq struct{

}

type QryWorldMapAck struct{
	proto.BaseAck
	Citys [] slgdb.City  `json:"citys"`
}

type GarrisonCityReq struct {
	GeneralId      	uint32   `json:"generalId"`
	CityId          int      `json:"cityId"`
}

type GarrisonCityAck struct {
	proto.BaseAck
	GeneralId      	uint32   `json:"generalId"`
	CityId          int      `json:"cityId"`
}


type AttackCityReq struct {
	GeneralId      	uint32   `json:"generalId"`
	CityId          int      `json:"cityId"`
}

type AttackCityAck struct {
	proto.BaseAck
	CityId          int      		`json:"cityId"`
	General       	slgdb.General	`json:"general"`
}