package app

import (
	"errors"
	"fmt"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/utils"
	"sync"
)

var ServerMgr serverMgr

func init()  {
	ServerMgr = serverMgr{}
}

type serverMgr struct {
	serverMap map[string] proto.ServerInfo
	mutex sync.RWMutex
}

func (s *serverMgr) Update(serverMap map[string] proto.ServerInfo) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	//utils.Log.Info("Update map: %v",serverMap)
	s.serverMap = serverMap
}

func (s *serverMgr) GetServerMap() map[string] proto.ServerInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.serverMap
}

func (s *serverMgr) GetGameScenesMap() map[string] proto.GameServersInfo {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	m := make(map[string] proto.GameServersInfo)

	for k, v:= range s.serverMap{
		if v.State == proto.ServerStateNormal && v.Type == proto.ServerTypeGame{
			s := proto.GameServersInfo{}
			s.Id = v.Id
			s.Name = v.Name
			s.ProxyName = v.ProxyName
			m[k] = s
		}
	}
	return m
}

func (s *serverMgr) HasServerById(id string) bool{
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, server := range s.serverMap {
		if server.Id == id {
			return true
		}
	}
	return false
}

/*
分配负载最低的服务
*/
func (s *serverMgr) Distribute(stype proto.ServerType) (proto.ServerInfo, error){
	s.mutex.Lock()
	defer s.mutex.Unlock()

	utils.Log.Info("Distribute type: %d, map: %v", stype, s.serverMap)
	var count uint32 = 1000000000
	var retServer proto.ServerInfo
	err := errors.New("not found server")
	for _, server := range s.serverMap {
		if server.Type == stype && server.OnlineCnt < count && server.State == proto.ServerStateNormal{
			count = server.OnlineCnt
			retServer = server
			err = nil
		}
	}
	return retServer, err
}

func (s *serverMgr) GetProxy(proxyName string) (string, error){
	s.mutex.Lock()
	defer s.mutex.Unlock()

	for _, server := range s.serverMap {
		if server.ProxyName == proxyName{
			return fmt.Sprintf("%s:%d", server.IP, server.Port), nil
		}
	}

	return "", errors.New(fmt.Sprintf("not found proxy:%s",proxyName))
}

