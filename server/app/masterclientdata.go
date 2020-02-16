package app

import "sync"

var MClientData masterClientData

func init() {
	MClientData = masterClientData{}
}

type masterClientData struct {
	lock 			sync.RWMutex
	onlineCnt		uint32		//server的在线数
}


func (s *masterClientData) GetOnlineCnt() uint32 {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.onlineCnt
}

func (s *masterClientData) Inc() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.onlineCnt++
}

func (s *masterClientData) Dec() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.onlineCnt--
}


func (s *masterClientData) SetOnlineCnt(onlineCnt uint32){
	s.lock.Lock()
	defer s.lock.Unlock()
	s.onlineCnt = onlineCnt
}