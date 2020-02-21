package slgproto

import (
	"github.com/llr104/LiFrame/proto"
)

const (
	Building_Dwelling = iota
	Building_Minefield
	Building_Farmland
	Building_Lumberyard
	Building_Barrack
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
	BuildType    int8	 `json:"type"`
	Build        string  `json:"build"`
	Yield        uint32  `json:"yield"`
}
