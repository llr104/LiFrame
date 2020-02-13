package master

import (
	"encoding/json"
	"fmt"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/utils"
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

type EnterMaster struct {
	liNet.BaseRouter
	nextProxyId int
	serverMap   map[string] proto.ServerInfo
	lock 		sync.Mutex
}


var Enter EnterMaster
func init() {
	Enter = EnterMaster{
		nextProxyId: 0,
		serverMap:make(map[string]proto.ServerInfo),
	}

	utils.Scheduler.NewTimerInterval(10*time.Second, utils.IntervalForever, checkClientLive, []interface{}{})
}


func checkClientLive(v ...interface{}){
	Enter.liveCheck()
}


func (s *EnterMaster) NameSpace() string{
	return "EnterMaster"
}

func (s *EnterMaster) getServerMap() map[string]proto.ServerInfo{
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.serverMap
}

func (s *EnterMaster) liveCheck() {


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


func (s* EnterMaster) ServerInfoReport(req liFace.IRequest){

	utils.Log.Info("ServerInfoReport req: %s", req.GetMsgName())
	info := proto.ServerInfoReport{}
	err := json.Unmarshal(req.GetData(), &info)
	if err != nil{
		utils.Log.Info("ServerInfoReport req error:",err.Error())
	}else{
		info.State = proto.ServerStateNormal
		info.LastTime = time.Now().Unix()

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

	utils.Log.Info("ServerInfoReport req: %v", info)
}

func (s* EnterMaster) Ping(req liFace.IRequest){
	utils.Log.Info("Ping req: %s", req.GetMsgName())
	info := proto.PingPong{}
	info.CurTime = time.Now().Unix()
	data, _ := json.Marshal(info)
	req.GetConnection().SendMsg(proto.MasterClientPong, data)
}

func (s* EnterMaster) ServerListReq(req liFace.IRequest){
	utils.Log.Info("ServerListReq req begin: %s", req.GetMsgName())
	info := proto.ServerListReq{}
	json.Unmarshal(req.GetData(), &info)

	//发送服务器列表
	ack := proto.ServerListAck{}
	ack.ServerMap = s.serverMap
	data, _ := json.Marshal(ack)
	req.GetConnection().SendMsg(proto.MasterClientServerListAck, data)
	utils.Log.Info("ServerListReq req end: %v", info)
}