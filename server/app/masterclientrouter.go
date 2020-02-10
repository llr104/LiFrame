
package app

import (
	"LiFrame/core/liFace"
	"LiFrame/core/liNet"
	"LiFrame/proto"
	"LiFrame/utils"
	"encoding/json"
	"os"
	"time"
)

var MClientRouter MasterClientRouter

type MasterClientRouter struct {
	liNet.BaseRouter
	isShutDown bool
}

func (s *MasterClientRouter) NameSpace() string {
	return "MasterClient"
}

func (s *MasterClientRouter) Pong(req liFace.IRequest) {
	utils.Log.Info("Pong:%s", req.GetMsgName())
}

func (s *MasterClientRouter) ServerListAck(req liFace.IRequest) {
	utils.Log.Info("ServerListAck begin: %s", req.GetMsgName())

	ackInfo := proto.ServerListAck{}
	err := json.Unmarshal(req.GetData(), &ackInfo)
	if err != nil{
		utils.Log.Info("ServerListAck error:%s",err.Error())
	}else{
		ServerMgr.Update(ackInfo.ServerMap)
	}

	utils.Log.Info("ServerListAck end: %v", ackInfo)
}

func (s *MasterClientRouter) ShutDown(req liFace.IRequest) {
	utils.Log.Info("ShutDown:%s", req.GetMsgName())

	if s.isShutDown == false {
		//是否需要做一些退出操作
		s.isShutDown = true
		f := GetShutDownFunc()
		if f != nil{
			f()
		}
		time.Sleep(5*time.Second)
		os.Exit(0)
	}

}
