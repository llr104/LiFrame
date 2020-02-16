package gate

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/utils"
	"io"
	"strconv"
	"strings"
	"sync"
	"time"
)

var MyGate gate

type offline struct {
	offlineProxyMap 	map[string] *liNet.Client
	offlineTime         int64
}

type gate struct {
	onlineProxyMap map[string] map[string] *liNet.Client //key:handshakeId key:proxyId
	offlineMap     map[string] offline
	wsUserHMap 	   map[uint32] *liNet.WsConnection //key:userId
	lock           sync.RWMutex
}

func init() {

	MyGate = gate{onlineProxyMap: make(map[string]map[string] *liNet.Client),
		offlineMap:make(map[string] offline),
		wsUserHMap:make( map[uint32] *liNet.WsConnection)}

	utils.Scheduler.NewTimerInterval(10*time.Second, utils.IntervalForever, checkOffLine, []interface{}{})
}

func (g*gate) ProxyClient(wsConn* liNet.WsConnection, msgProxyId string, router liFace.IRouter) (*liNet.Client, bool ){
	g.lock.Lock()
	defer g.lock.Unlock()

	msgProxy, err := app.ServerMgr.GetProxy(msgProxyId)
	if err != nil{
		utils.Log.Warn("%s", err.Error())
		wsConn.WriteMessage(msgProxyId, proto.ProxyError, []byte(err.Error()))
		return nil, false
	}

	id, err1 := wsConn.GetProperty("handshakeId")
	if err1 != nil{
		return nil, false
	}
	handshakeId := id.(string)
	isNeedCreate := false

	if _, ok := g.onlineProxyMap[handshakeId]; ok == false{
		g.onlineProxyMap[handshakeId] = make(map[string] *liNet.Client)
		isNeedCreate = true
	}

	proxyMap := g.onlineProxyMap[handshakeId]
	if clientProxy, ok := proxyMap[msgProxyId]; ok == false {
		isNeedCreate = true
	}else{
		if clientProxy.GetConn() == nil || clientProxy.GetConn().IsClose() {
			isNeedCreate = true
		}
	}

	if isNeedCreate {
		//创建proxy
		delete(proxyMap, msgProxyId)

		clientId := fmt.Sprintf("id_%d_%s", wsConn.GetId(), msgProxy)
		arr := strings.Split(msgProxy,":")
		if len(arr) != 2 {
			return nil, false
		}
		name := fmt.Sprintf("name_%s", msgProxy)
		port, _ := strconv.Atoi(arr[1])
		c := liNet.NewClient(name, clientId, arr[0], port)
		c.AddRouter(router)

		c.Start()
		c.Running()
		if c.GetConn() != nil && c.GetConn().IsClose() == false{
			proxyMap[msgProxyId] = c
		}
	}

	proxyClient, b := g.onlineProxyMap[handshakeId][msgProxyId]
	return proxyClient, b
}

func (g*gate) CloseProxy(wsConn* liNet.WsConnection, msgProxyId string) {
	g.lock.Lock()
	defer g.lock.Unlock()

	id, err := wsConn.GetProperty("handshakeId")
	if err != nil{
		utils.Log.Info("closeProxy not handshakeId")
	}else{
		handshakeId := id.(string)
		if proxy, ok := g.onlineProxyMap[handshakeId][msgProxyId]; ok{
			proxy.Stop()
			delete(g.onlineProxyMap[handshakeId], msgProxyId)
		}
	}
}

func (g*gate) closeAllProxy(wsConn* liNet.WsConnection) {

	id, err := wsConn.GetProperty("handshakeId")
	if err == nil{
		handshakeId := id.(string)

		g.lock.Lock()
		proxyMap, ok1 := g.onlineProxyMap[handshakeId]
		off, ok2 := g.offlineMap[handshakeId]
		delete(g.onlineProxyMap, handshakeId)
		delete(g.offlineMap, handshakeId)
		g.lock.Unlock()

		if ok1 {
			for _,v := range proxyMap {
				v.Stop()
			}
		}

		if ok2 {
			for _, v := range off.offlineProxyMap{
				v.Stop()
			}
		}
	}
}

func (g*gate) Reconnect(wsConn* liNet.WsConnection, handshakeId string) string{
	g.lock.Lock()
	defer g.lock.Unlock()

	newHandshakeId := g.NewHandshakeId(wsConn.GetId())
	if m, ok := g.offlineMap[handshakeId]; ok{
		utils.Log.Info("断线回来，代理归属到新的连接")
		g.onlineProxyMap[newHandshakeId] = m.offlineProxyMap
		delete(g.offlineMap, handshakeId)

		//通知其所在的代理，用户在线了
		if p, err := wsConn.GetProperty("userId"); err == nil {
			userId := p.(uint32)
			for _, v := range m.offlineProxyMap {
				c := v.GetConn()
				if c != nil{
					pack := proto.UserOnlineOrOffLineReq{}
					pack.UserId = userId
					pack.Type = proto.UserOnline
					data, _ := json.Marshal(pack)
					c.SendMsg(proto.SystemUserOnOrOffReq, data)
				}
			}
		}
	}else{
		utils.Log.Info("断线回来，代理已经不存在了")
	}
	wsConn.SetProperty("handshakeId", newHandshakeId)
	return newHandshakeId
}

func (g*gate) ConnectEnter(wsConn* liNet.WsConnection) string {
	handshakeId := g.NewHandshakeId(wsConn.GetId())
	wsConn.SetProperty("handshakeId", handshakeId)
	return  handshakeId
}

func (g*gate) ConnectExit(wsConn* liNet.WsConnection){

	id, err := wsConn.GetProperty("handshakeId")
	if err  != nil {
		return
	}
	handshakeId := id.(string)
	g.lock.Lock()
	if proxyMap, ok := g.onlineProxyMap[handshakeId]; ok {
		off := offline{}
		off.offlineProxyMap = proxyMap
		off.offlineTime = time.Now().Unix()
		g.offlineMap[handshakeId] = off
		delete(g.onlineProxyMap, handshakeId)

		//通知其所在的代理，用户断线了
		if p, err := wsConn.GetProperty("userId"); err == nil {
			userId := p.(uint32)
			for _, v := range proxyMap {
				c := v.GetConn()
				if c != nil{
					pack := proto.UserOnlineOrOffLineReq{}
					pack.UserId = userId
					pack.Type = proto.UserOffline
					data, _ := json.Marshal(pack)
					c.SendMsg(proto.SystemUserOnOrOffReq, data)
				}
			}
		}
	}
	g.lock.Unlock()

	g.userExit(wsConn)
}

func (g *gate) userEnter(wsConn* liNet.WsConnection) {

	if p, err := wsConn.GetProperty("userId"); err == nil{
		userId := p.(uint32)
		g.lock.Lock()
		oldConn, ok := g.wsUserHMap[userId]
		g.lock.Unlock()

		if ok && oldConn != wsConn{
			utils.Log.Info("账号在其他地方登录")
			//发送消息 todo

			g.closeAllProxy(oldConn)
			oldConn.Close()
		}

		g.lock.Lock()
		g.wsUserHMap[userId] = wsConn
		g.lock.Unlock()
	}

}

func (g *gate) userExit(wsConn* liNet.WsConnection) {

	if p, err := wsConn.GetProperty("userId"); err == nil{
		userId := p.(uint32)
		g.lock.Lock()
		defer g.lock.Unlock()
		delete(g.wsUserHMap, userId)
	}
}

func (g*gate) check(){
	g.lock.Lock()
	defer g.lock.Unlock()

	/*
		代理保留2分钟
	*/
	t := time.Now().Unix()
	for key, v:= range g.offlineMap {
		if t - v.offlineTime >120 {
			for _, v1 := range v.offlineProxyMap{
				v1.Stop()
			}
			utils.Log.Info("%s代理过期被删除了", key)
			delete(g.offlineMap, key)
		}
	}
}

/*
gate生成uuid用户handshake
*/
func (g*gate) NewHandshakeId(cid uint64) string {
	t := time.Now().UnixNano()
	str := fmt.Sprintf("%d_%d", t, cid)
	w := md5.New()
	io.WriteString(w, str)
	uuid := fmt.Sprintf("%x", w.Sum(nil))
	return uuid
}


func checkOffLine(v ...interface{})   {
	MyGate.check()
}

