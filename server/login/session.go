package login

import (
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/utils"
	"sync"
	"time"
)

var SessLoginMgr loginSessionMgr

func init() {
	SessLoginMgr = loginSessionMgr{
		sessionMap:make(map[uint32] string),
		liveTimeMap:make(map[uint32] int64),
	}

	utils.Scheduler.NewTimerInterval(60*time.Second, utils.IntervalForever, checkSessionLive, []interface{}{})
}

func checkSessionLive(v ...interface{}) {
	curTime := time.Now().Unix()
	for k,v := range SessLoginMgr.liveTimeMap{
		if curTime> v{
			SessLoginMgr.RemoveSessionForce(k)
		}
	}
}

type loginSessionMgr struct {
	loginLock 	   	sync.RWMutex
	sessionMap 		map[uint32] string      //key-value:userId-session
	liveTimeMap 	map[uint32] int64      //key-value:userId-time
}

func (s*loginSessionMgr) NewSession(serverId string, id uint32, conn liFace.IConnection) string {

	oldSession, ok := s.liveSession(id)
	if ok == false {
		utils.Log.Info("%s:被顶号了", oldSession)
		app.SessionMgr.SessionExit(oldSession)
	}

	session := app.SessionMgr.CreateSession(serverId, id)
	app.SessionMgr.SessionEnter(session, conn)

	s.loginLock.Lock()
	defer s.loginLock.Unlock()

	s.sessionMap[id] = session
	//有效期5分钟
	s.liveTimeMap[id] = time.Now().Unix() + 5*60

	return session
}

func (s*loginSessionMgr) RemoveSessionForce(userId uint32) {
	s.loginLock.Lock()
	defer s.loginLock.Unlock()

	if sess, ok := s.sessionMap[userId]; ok{
		app.SessionMgr.SessionExit(sess)
	}

	delete(s.sessionMap, userId)
	delete(s.liveTimeMap, userId)
}

func (s*loginSessionMgr) RemoveSession(userId uint32, session string) {
	s.loginLock.Lock()
	defer s.loginLock.Unlock()

	sess, ok := s.sessionMap[userId]
	if ok {
		if session == sess {
			delete(s.sessionMap, userId)
			delete(s.liveTimeMap, userId)
			app.SessionMgr.SessionExit(session)
		}
	}
}


func (s*loginSessionMgr) SessionIsLive(userId uint32, session string) bool {
	s.loginLock.Lock()
	defer s.loginLock.Unlock()

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

func (s *loginSessionMgr) liveSession(userId uint32) (string, bool){
	s.loginLock.Lock()
	defer s.loginLock.Unlock()

	sess, ok := s.sessionMap[userId]
	if ok {
		live := s.liveTimeMap[userId]
		if live > time.Now().Unix(){
			return sess, true
		}
	}
	return "", false
}

func (s*loginSessionMgr) SessionKeepLive(userId uint32, session string) {
	s.loginLock.Lock()
	defer s.loginLock.Unlock()

	sess, ok := s.sessionMap[userId]
	if ok {
		if session == sess {
			s.liveTimeMap[userId] = time.Now().Unix() + 5*60
		}
	}
}

func (s*loginSessionMgr) SessionExitByConn(conn liFace.IConnection){
	app.SessionMgr.SessionExitByConn(conn)
}


func (s *loginSessionMgr) OnlineCnt() int {
	s.loginLock.Lock()
	defer s.loginLock.Unlock()

	return len(s.sessionMap)
}




