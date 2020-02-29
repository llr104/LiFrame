package slgproto

import (
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
)

const (
	BuildingDwelling = iota
	BuildingMinefield
	BuildingFarmland
	BuildingLumberyard
	BuildingBarrack
)

type QryBuildingQeq struct {
	BuildType    int8	 `json:"type"`
}

type QryBuildingAck struct {
	proto.BaseAck
	BuildType    int8	 `json:"type"`
	Buildings    string  `json:"buildings"`
	Yield        uint32  `json:"yield"`
}

type UpBuildingQeq struct {
	BuildType    int8	 `json:"type"`
	BuildId      int     `json:"Id"`
}

type UpBuildingAck struct {
	proto.BaseAck
	BuildType    int8	 	`json:"type"`
	Build        string  	`json:"build"`
	Yield        uint32  	`json:"yield"`
	Role		 slgdb.Role	`json:"role"`
}

type QryGeneralReq struct {

}

type QryGeneralAck struct {
	proto.BaseAck
	Generals[] *slgdb.General `json:"generals"`
}

