package app

import (
	"LiFrame/core/liFace"
	"LiFrame/core/liNet"
	"LiFrame/proto"
	"LiFrame/utils"
	"encoding/json"
)

var GRouter *GateRouter

func init() {
	GRouter = &GateRouter{}
}
type GateRouter struct {
	liNet.BaseRouter
}

func (s *GateRouter) NameSpace() string {
	return "*.*"
}

func (s *GateRouter) EveryThingHandle(req liFace.IRequest) {
	conn, err := req.GetConnection().GetProperty("gateConn")
	if err != nil{
		utils.Log.Warn("EveryThingHandle not found gateConn")
	}

	gateConn := conn.(*liNet.WsConnection)
	if req.GetMsgName() == proto.EnterLoginLoginAck{
		loginAck := proto.LoginAck{}
		err := json.Unmarshal(req.GetData(),&loginAck)
		if err == nil && loginAck.Code == proto.Code_Success{
			gateConn.SetProperty("isAuth", true)
		}else{
			gateConn.SetProperty("isAuth", false)
		}
	}

	name := req.GetMsgName()
	proxy, e :=  req.GetConnection().GetProperty("proxy")
	if e != nil{
		gateConn.WriteMessage("", name, req.GetData())
	}else{
		proxyName := proxy.(string)
		gateConn.WriteMessage(proxyName, name, req.GetData())
	}

}
