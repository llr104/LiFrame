package gameslg

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/utils"
)

var Enter enterGame

func init() {
	Enter = enterGame{}
}


type enterGame struct {
	liNet.BaseRouter
}

func (s *enterGame) NameSpace() string {
	return "*.*"
}

func (s *enterGame) EveryThingHandle(req liFace.IRequest, rsp liFace.IMessage) {
	if req.GetMessage().GetMsgName() == proto.GameEnterGameReq{
		reqInfo := proto.EnterGameReq{}
		ackInfo := proto.EnterGameAck{}
		if err := json.Unmarshal(req.GetMessage().GetBody(), &reqInfo); err != nil {
			ackInfo.Code = proto.Code_Illegal
			utils.Log.Info("GameEnterGameReq error:%s", err.Error())
		} else {
			ackInfo.Code = proto.Code_Success
			req.GetConnection().SetProperty("userId", reqInfo.UserId)
		}
		data, _ := json.Marshal(ackInfo)
		rsp.SetBody(data)
	}
}

