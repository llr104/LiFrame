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

func (s *router) EveryThingHandle(req liFace.IRequest, rsp liFace.IMessage) {

	conn, err := req.GetConnection().GetProperty("gateConn")
	if err != nil{
		utils.Log.Warn("EveryThingHandle not found gateConn")
	}

	msg := req.GetMessage()
	msgName := msg.GetMsgName()
	gateConn := conn.(*liNet.WsConnection)

	proxy, e :=  req.GetConnection().GetProperty("proxy")
	if e != nil{
		gateConn.Push("", msgName,  msg.GetBody())
	}else{
		proxyName := proxy.(string)
		gateConn.Response(proxyName, msgName, msg.GetSeq(), msg.GetBody())
	}
}

func (s *router) Handle(rsp liFace.IRespond) {
	req := rsp.GetRequest()
	conn, err := req.GetConnection().GetProperty("gateConn")
	if err != nil{
		utils.Log.Warn("Handle not found gateConn")
	}

	msg := rsp.GetMessage()
	msgName := msg.GetMsgName()
	gateConn := conn.(*liNet.WsConnection)
	if msgName == proto.EnterLoginLoginReq{
		loginAck := proto.LoginAck{}
		err := json.Unmarshal(msg.GetBody(), &loginAck)
		if err == nil && loginAck.Code == proto.Code_Success{
			gateConn.SetProperty("session", loginAck.Session)
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
		gateConn.Push("", msgName,  rsp.GetData())
	}else{
		proxyName := proxy.(string)
		gateConn.Response(proxyName, msgName, msg.GetSeq(), rsp.GetData())
	}

}
