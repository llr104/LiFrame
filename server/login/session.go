package login

import (
	"LiFrame/utils"
	"time"
)

var LoginSessMgr LoginSessionMgr

func init() {
	LoginSessMgr = LoginSessionMgr{
		sessionMap:make(map[uint32] string),
		liveTimeMap:make(map[uint32] int64),
	}

	utils.Scheduler.NewTimerInterval(60*time.Second, utils.IntervalForever, checkSessionLive, []interface{}{})
}

func checkSessionLive(v ...interface{}) {
	curTime := time.Now().Unix()
	for k,v := range LoginSessMgr.liveTimeMap{
		if curTime> v{
			LoginSessMgr.RemoveSessionForce(k)
		}
	}
}

type LoginSessionMgr struct {
	sessionMap 		map[uint32] string      //key-value:userId-session
	liveTimeMap 	map[uint32] int64      //key-value:userId-time
}

func (s*LoginSessionMgr) AddSession(id uint32,session string)  {
	s.sessionMap[id] = session
	//有效期5分钟
	s.liveTimeMap[id] = time.Now().Unix() + 5*60
}

func (s*LoginSessionMgr) RemoveSessionForce(userId uint32) {
	delete(s.sessionMap, userId)
	delete(s.liveTimeMap, userId)
}

func (s*LoginSessionMgr) RemoveSession(userId uint32, session string) {
	sess, ok := s.sessionMap[userId]
	if ok {
		if session == sess {
			delete(s.sessionMap, userId)
			delete(s.liveTimeMap, userId)
		}
	}
}

func (s*LoginSessionMgr) SessionIsLive(userId uint32, session string) bool {
	sess, ok := s.sessionMap[userId]
	if ok {
		if session == sess {
			live := s.liveTimeMap[userId]
			if live > time.Now().Unix(){
				return true
			}
		}
	}
	return false
}

func (s*LoginSessionMgr) SessionKeepLive(userId uint32, session string) {
	sess, ok := s.sessionMap[userId]
	if ok {
		if session == sess {
			s.liveTimeMap[userId] = time.Now().Unix() + 5*60
		}
	}
}

func (s *LoginSessionMgr) OnlineCnt() int {
	return len(s.sessionMap)
}




