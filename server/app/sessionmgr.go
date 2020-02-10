package app

import (
	"LiFrame/utils"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/thinkoner/openssl"
	"strings"
	"time"
)

var SessionMgr SessionManager

func init() {
	SessionMgr = SessionManager{sessionKey:"sessionKey123456"}
}

type SessionManager struct {
	sessionKey string //密钥的长度可以是16/24/32个字符（128/192/256位）
}

/*
根据用户id生成一个session
*/
func (s* SessionManager) CreateSession(appId string, userId uint32) string {
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

func  (s* SessionManager) CheckSessionFrom(session string) (string, error){

	key := []byte(s.sessionKey)
	bytes, err := base64.StdEncoding.DecodeString(session)
	if err != nil{
		return "", errors.New("session not found from server")
	}

	dSrc, _ := openssl.AesECBDecrypt(bytes, key, openssl.PKCS7_PADDING)
	//fmt.Println("dSrc:",dSrc,err)

	arr := strings.Split(string(dSrc), "_")
	if len(arr) != 3{
		return "", errors.New("session not found from server")
	}else{
		appId := arr[0]
		return appId, nil
	}
}


func  (s* SessionManager) CheckSessionValid(session string, userId uint32) bool{
	
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





