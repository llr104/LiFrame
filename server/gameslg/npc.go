package gameslg

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/server/gameslg/data"
	"github.com/llr104/LiFrame/server/gameslg/slgproto"
	"github.com/llr104/LiFrame/utils"
)

var NPC npc

func init() {
	NPC = npc{}
	NPC.scenes = make(map[int]*data.NpcScene)

}

type npc struct {
	liNet.BaseRouter
	scenes map[int]*data.NpcScene
}

func (s *npc) After()  {
	/*
		先简单点默认随机几个npc场景，后面会改成配表，先做大概功能
	*/

	for i:=0; i<3; i++ {
		NPC.scenes[i] = data.RandomNPCScene(uint16(i))
	}
}

func (s *npc) NameSpace() string {
	return "npc"
}

func (s *npc) PreHandle(req liFace.IRequest) bool{
	_, err := req.GetConnection().GetProperty("roleId")
	if err == nil {
		return true
	}else{
		utils.Log.Warning("%s not has roleId", req.GetMsgName())
		return false
	}
}


func (s *npc) QrySceneReq(req liFace.IRequest)  {
	reqInfo := slgproto.QryNpcSceneReq{}
	ackInfo := slgproto.QryNpcSceneAck{}
	json.Unmarshal(req.GetData(), &reqInfo)
	ackInfo.Code = slgproto.CodeSlgSuccess

	n := len(s.scenes)
	ackInfo.NpcScenes = make([]data.NpcScene, n)
	i := 0
	for _, v := range s.scenes{
		ackInfo.NpcScenes[i] = *v
		i++
	}

	data, _ := json.Marshal(ackInfo)
	req.GetConnection().SendMsg(slgproto.NpcQrySceneAck, data)

}