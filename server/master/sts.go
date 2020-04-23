package master

import (
	"encoding/json"
	"fmt"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/utils"
	"strings"
	"sync"
	"time"
)

func ClientConnStart(conn liFace.IConnection){
	utils.Log.Info("ClientConnStart:%s", conn.RemoteAddr().String())
}

func ClientConnStop(conn liFace.IConnection){
	utils.Log.Info("ClientConnStop:%s", conn.RemoteAddr().String())
}

func ShutDown(){
	utils.Log.Info("ShutDown")
}

type sts struct {
	liNet.BaseRouter
	nextProxyId int
	serverMap   map[string] proto.ServerInfo
	lock 		sync.Mutex
}

var STS sts

func init() {
	STS = sts{
		nextProxyId: 0,
		serverMap:make(map[string]proto.ServerInfo),
	}

	utils.Scheduler.NewTimerInterval(10*time.Second, utils.IntervalForever, checkClientLive, []interface{}{})
}


func checkClientLive(v ...interface{}){
	STS.liveCheck()
}


func (s *sts) NameSpace() string{
	return "System"
}

func (s *sts) getServerMap() map[string]proto.ServerInfo{
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.serverMap
}

func (s *sts) liveCheck() {

	s.lock.Lock()
	defer s.lock.Unlock()
	/*
		超过60s没有上报服务信息，认为断线
	*/
	var cur = time.Now().Unix()
	for k,v := range s.serverMap {
		if cur - v.LastTime > 60{
			v.State = proto.ServerStateDead
			s.serverMap[k] = v
		}
	}
}


func (s*sts) ServerInfoReport(req liFace.IRequest, rsp liFace.IMessage){
	msg := req.GetMessage()
	remote := req.GetConnection().GetTCPConnection().RemoteAddr().String()
	info := proto.ServerInfoReport{}
	sArr := strings.Split(remote, ":")
	if len(sArr) != 2{
		return
	}

	ip := sArr[0]
	err := json.Unmarshal(msg.GetBody(), &info)
	utils.Log.Info("ServerInfoReport %v ", info)

	if err != nil{
		utils.Log.Info("ServerInfoReport req error:",err.Error())
	}else{
		info.State = proto.ServerStateNormal
		info.LastTime = time.Now().Unix()
		info.IP = ip

		s.lock.Lock()
		defer s.lock.Unlock()

		last, ok := s.serverMap[info.Id]
		if ok {
			info.ProxyName = last.ProxyName
			s.serverMap[info.Id] = info.ServerInfo
		}else{
			info.ServerInfo.ProxyName = fmt.Sprintf("%d", s.nextProxyId)
			s.serverMap[info.Id] = info.ServerInfo
			s.nextProxyId++
		}
	}

}

func (s* sts) Ping(req liFace.IRequest, rsp liFace.IMessage){
	utils.Log.Info("Ping")
	info := proto.PingPong{}
	info.CurTime = time.Now().Unix()
	data, _ := json.Marshal(info)
	rsp.SetBody(data)
}


func (s*sts) ServerListReq(req liFace.IRequest, rsp liFace.IMessage){
	msg := req.GetMessage()
	utils.Log.Info("ServerListReq req : %s", req.GetConnection().GetTCPConnection().RemoteAddr())
	info := proto.ServerListReq{}
	json.Unmarshal(msg.GetBody(), &info)

	//发送服务器列表
	ack := proto.ServerListAck{}
	ack.ServerMap = s.serverMap
	data, _ := json.Marshal(ack)
	rsp.SetBody(data)

}