package slgproto

import (
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/gameslg/data"
)

type QryNpcSceneReq struct{

}

type QryNpcSceneAck struct{
	proto.BaseAck
	NpcScenes []data.NpcScene  `json:"scenes"`
}
