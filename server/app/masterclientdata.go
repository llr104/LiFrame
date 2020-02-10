package app

import "sync"

var MClientData MasterClientData

func init() {
	MClientData = MasterClientData{}
}

type MasterClientData struct {
	lock 			sync.RWMutex
	onlineCnt		uint32		//server的在线数
}


func (s* MasterClientData) GetOnlineCnt() uint32 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.onlineCnt
}

func (s* MasterClientData) Inc() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.onlineCnt++
}

func (s* MasterClientData) Dec() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.onlineCnt--
}


func (s* MasterClientData) SetOnlineCnt(onlineCnt uint32){
	s.lock.Lock()
	defer s.lock.Unlock()
	s.onlineCnt = onlineCnt
}