package world

import (
	"LiFrame/core/liNet"
	"LiFrame/proto"
	client2 "LiFrame/server/app"
	"LiFrame/utils"
	"encoding/json"
	"time"
)

var OnlineInstance Online

func init() {
	OnlineInstance = Online{connectMap: make(map[uint32] *liNet.Connection),}
	utils.Scheduler.NewTimerInterval(30*time.Second, utils.IntervalForever, checkOnline, []interface{}{})
}


func checkOnline(v ...interface{}) {
	curTime := time.Now().Unix()

	for k,v := range OnlineInstance.connectMap{
		if v.IsClose() {
			delete(OnlineInstance.connectMap, k)
		}else{
			t, err1 := v.GetProperty("lastKeepLive")
			u, err2 := v.GetProperty("userId")
			s, err3 := v.GetProperty("session")

			if err1 == nil && err2 == nil && err3 == nil {
				l := t.(int64)
				userId := u.(uint32)
				session := s.(string)

				if curTime - l >60{
					//60s没有同步session到login服，需要上报
					sessReq := proto.SessionUpdateReq{}
					sessReq.Session = session
					sessReq.UserId = userId
					sessReq.ConnId = 0
					sessReq.OpType = proto.SessionOpKeepLive

					appId, _ := client2.SessionMgr.CheckSessionFrom(session)
					client, ok := W2Login.GetLoginClient(appId)
					if ok {
						data, _ := json.Marshal(sessReq)
						client.GetConn().SendMsg(proto.EnterLoginSessionUpdateReq, data)
						v.SetProperty("lastKeepLive", curTime)
					}
				}

			}else{
				delete(OnlineInstance.connectMap, k)
			}
		}
	}
}

type Online struct {
	connectMap 		map[uint32] *liNet.Connection //key:userId
}

func (s *Online) Join(conn*liNet.Connection) {
	u, err := conn.GetProperty("userId")
	if err == nil{
		userId := u.(uint32)
		s.connectMap[userId] = conn
	}
}

func (s *Online) Exit(conn*liNet.Connection)  {
	u, err := conn.GetProperty("userId")
	if err == nil{
		userId := u.(uint32)
		delete(s.connectMap, userId)
	}
}