package gate

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/utils"
)

var Router *router

func init() {
	Router = &router{}
}

type router struct {
	liNet.BaseRouter
}

func (s *router) NameSpace() string {
	return "*.*"
}

func (s *router) EveryThingHandle(req liFace.IRequest) {
	conn, err := req.GetConnection().GetProperty("gateConn")
	if err != nil{
		utils.Log.Warn("EveryThingHandle not found gateConn")
	}

	msgName := req.GetMsgName()
	gateConn := conn.(*liNet.WsConnection)
	if msgName == proto.EnterLoginLoginAck{
		loginAck := proto.LoginAck{}
		err := json.Unmarshal(req.GetData(),&loginAck)
		if err == nil && loginAck.Code == proto.Code_Success{
			gateConn.SetProperty("session",loginAck.Session)
			gateConn.SetProperty("userId", loginAck.Id)
			//处理踢号
			MyGate.userEnter(gateConn)
		}else{
			MyGate.userExit(gateConn)
			gateConn.RemoveProperty("session")
			gateConn.RemoveProperty("userId")
		}
	}

	proxy, e :=  req.GetConnection().GetProperty("proxy")
	if e != nil{
		gateConn.WriteMessage("", msgName, req.GetData())
	}else{
		proxyName := proxy.(string)
		gateConn.WriteMessage(proxyName, msgName, req.GetData())
	}

}
