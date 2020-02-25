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
