package app

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/utils"
	"github.com/thinkoner/openssl"
	"strings"
	"sync"
	"time"
)

var SessionMgr sessionMgr

func init() {
	SessionMgr = sessionMgr{sessionKey: "sessionKey123456", connMap:make(map[string] liFace.IConnection)}
}

type sessionMgr struct {
	sessionKey string
	connMap    map[string] liFace.IConnection
	lock 	   sync.RWMutex
}

/*
根据用户id生成一个session
*/
func (s*sessionMgr) CreateSession(appId string, userId uint32) string {
	unix := time.Now().UnixNano()

	src := fmt.Sprintf("%s_%d_%d", appId, unix, userId)
	key := []byte(s.sessionKey) 
	dst, err := openssl.AesECBEncrypt([]byte(src), key, openssl.PKCS7_PADDING)
	if err != nil{
		fmt.Println("CreateSession:",err.Error())
	}

	session := base64.StdEncoding.EncodeToString(dst)
	dSrc , _:= openssl.AesECBDecrypt(dst, key, openssl.PKCS7_PADDING)

	utils.Log.Info("CreateSession src: %s",string(dSrc))
	utils.Log.Info("CreateSession session: %s", session)

	return session
}

func  (s*sessionMgr) CheckSessionFrom(session string) (string, error){

	key := []byte(s.sessionKey)
	bytes, err := base64.StdEncoding.DecodeString(session)
	if err != nil{
		return "", errors.New("session not found from server")
	}

	dSrc, _ := openssl.AesECBDecrypt(bytes, key, openssl.PKCS7_PADDING)
	arr := strings.Split(string(dSrc), "_")
	if len(arr) != 3{
		return "", errors.New("session not found from server")
	}else{
		appId := arr[0]
		return appId, nil
	}
}


func  (s*sessionMgr) CheckSessionValid(session string, userId uint32) bool{
	
	key := []byte(s.sessionKey)
	bytes, err := base64.StdEncoding.DecodeString(session)
	if err != nil{
		return false
	}

	dSrc, _:= openssl.AesECBDecrypt(bytes, key, openssl.PKCS7_PADDING)
	//fmt.Println("dSrc:",dSrc,err)

	arr := strings.Split(string(dSrc), "_")
	if len(arr) != 3{
		return false
	}else{
		appId := arr[0]
		if ServerMgr.HasServerById(appId) == false{
			return false
		}

		userIdStr := fmt.Sprintf("%d",userId)
		if arr[2] == userIdStr{
			return true
		}else{
			return false
		}
	}

}

func (s *sessionMgr) SessionEnter(session string, conn liFace.IConnection)  {
	conn.SetProperty("session", session)

	//stop会触发SessionExitByConn，所以加锁要注意
	s.lock.Lock()
	oldConn, ok := s.connMap[session]
	s.lock.Unlock()

	if ok {
		if oldConn != conn{
			if oldConn.IsClose() == false {
				//关闭前可以发消息通知 todo
				utils.Log.Info("session:%s被新的连接顶替了", session)
				oldConn.Stop()
			}
		}
	}

	s.lock.Lock()
	s.connMap[session] = conn
	s.lock.Unlock()

}

func (s *sessionMgr) SessionExitByConn(conn liFace.IConnection)  {

	if v, err := conn.GetProperty("session"); err == nil{

		session := v.(string)
		utils.Log.Info("session:%s的连接断开了", session)

		s.lock.Lock()
		defer s.lock.Unlock()
		delete(s.connMap, session)
	}
}

func (s *sessionMgr) SessionExit(session string)  {
	s.lock.Lock()
	conn ,ok := s.connMap[session]
	s.lock.Unlock()

	//stop会触发SessionExitByConn，所以加锁要注意
	if ok {
		if conn.IsClose() == false {
			//关闭前可以发消息通知 todo
			conn.Stop()
			utils.Log.Info("session:%s的连接断开了", session)
		}
	}

	s.lock.Lock()
	delete(s.connMap, session)
	s.lock.Unlock()

}





