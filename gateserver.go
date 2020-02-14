package main

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db"
	"github.com/llr104/LiFrame/utils"
	"github.com/thinkoner/openssl"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var myGate gate
var cid uint64 = 0

type offline struct {
	offlineProxyMap 	map[string] *liNet.Client
	offlineTime         int64
}

type gate struct {
	onlineProxyMap 		map[string]map[string] *liNet.Client
	offlineMap          map[string] offline
	lock 		sync.RWMutex
}

func init() {
	myGate = gate{onlineProxyMap:make(map[string]map[string] *liNet.Client), offlineMap:make(map[string] offline)}
	utils.Scheduler.NewTimerInterval(10*time.Second, utils.IntervalForever, checkOffLine, []interface{}{})
}

func (g* gate) proxyClient(wsConn* liNet.WsConnection, msgProxyId string, router liFace.IRouter) (*liNet.Client, bool ){
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

		proxyKey := fmt.Sprintf("%s_%s", handshakeId, msgProxyId)
		arr := strings.Split(msgProxy,":")
		if len(arr) != 2 {
			return nil, false
		}
		name := fmt.Sprintf("name_%s", msgProxy)
		port, _ := strconv.Atoi(arr[1])
		c := liNet.NewClient(name, proxyKey, arr[0], port)
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

func (g* gate) closeProxy(wsConn* liNet.WsConnection, msgProxyId string) {
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

func (g* gate) reconnect(wsConn* liNet.WsConnection){
	g.lock.Lock()
	defer g.lock.Unlock()

	if id, err := wsConn.GetProperty("handshakeId"); err == nil{
		handshakeId := id.(string)
		if m, ok := g.offlineMap[handshakeId]; ok{
			utils.Log.Info("断线回来，代理归属到新的连接")
			g.onlineProxyMap[handshakeId] = m.offlineProxyMap
			delete(g.offlineMap, handshakeId)
		}
	}
}

func (g* gate) connectExit(wsConn* liNet.WsConnection){
	g.lock.Lock()
	defer g.lock.Unlock()

	id, err := wsConn.GetProperty("handshakeId")
	if err  != nil {
		return
	}
	handshakeId := id.(string)
	if proxyMap, ok := g.onlineProxyMap[handshakeId]; ok {
		off := offline{}
		off.offlineProxyMap = proxyMap
		off.offlineTime = time.Now().Unix()
		g.offlineMap[handshakeId] = off
		delete(g.onlineProxyMap, handshakeId)
	}
}

func (g* gate) check(){
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

func checkOffLine(v ...interface{})   {
	myGate.check()
}

/*
gate生成uuid用户handshake
*/
func newHandshakeId(cid uint64) string {
	t := time.Now().UnixNano()
	str := fmt.Sprintf("%d_%d", t, cid)
	w := md5.New()
	io.WriteString(w, str)
	uuid := fmt.Sprintf("%x", w.Sum(nil))
	return uuid
}

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}



func wsHandler(resp http.ResponseWriter, req *http.Request) {
	// 应答客户端告知升级连接为websocket
	wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
	if err != nil {
		return
	}

	cid++
	app.MClientData.Inc()
	wsConn := liNet.NewWsConnection(wsSocket, cid)
	wsConn.SetOnMessage(handleWsMessage)
	wsConn.SetOnClose(handleOnClose)
	wsConn.Running()
}

func handleOnClose(wsConn *liNet.WsConnection)  {
	app.MClientData.Dec()
	myGate.connectExit(wsConn)
	utils.Log.Debug("handleOnClose wsCount:%d", app.MClientData.GetOnlineCnt())
}

func handleWsMessage(wsConn *liNet.WsConnection, req *liNet.WsMessage) {
	if req.MsgType == websocket.TextMessage{
		return
	}

	//解压包
	b := new(bytes.Buffer)
	binary.Write(b, binary.LittleEndian, req.Data)
	r, err := gzip.NewReader(b)
	if err != nil{
		return
	}

	defer r.Close()
	unzipData, err := ioutil.ReadAll(r)
	if err!=nil{
		return
	}

	encode, err := hex.DecodeString(string(unzipData))
	if err != nil{
		return
	}

	decode, _ := openssl.AesCBCDecrypt(encode, liNet.GateMessageKey, liNet.GateMessageKey, openssl.ZEROS_PADDING)
	data := string(decode)

	msgArr := strings.Split(data,"|")
	if len(msgArr) == 3{

		msgName := msgArr[0]
		msgProxyId := msgArr[1]
		body := msgArr[2]

		if msgName == proto.GateHandshake{
			if body == ""{
				utils.Log.Info("不是断线重连")
				handshakeId := newHandshakeId(wsConn.GetId())
				wsConn.SetProperty("handshakeId", handshakeId)
				wsConn.WriteMessage("", proto.GateHandshake, []byte(handshakeId))
			}else{
				utils.Log.Info("是断线重连")
				wsConn.SetProperty("handshakeId", body)
				myGate.reconnect(wsConn)
				wsConn.WriteMessage("", proto.GateHandshake, []byte(body))
			}
			return
		}

		_, err := wsConn.GetProperty("handshakeId")
		if err != nil{
			return
		}

		/*检测授权是否成功才转发消息*/
		if msgName == proto.GateLoginServerReq{
			ackInfo := proto.DistributeServerAck{}
			if serverInfo, err:= app.ServerMgr.Distribute(proto.ServerTypeLogin); err != nil {
				ackInfo.Code = proto.Code_Not_Server
				utils.Log.Info("gate.LoginServerReq error:%s", err.Error())
			}else{
				ackInfo.Code = proto.Code_Success
				ackInfo.ServerInfo = serverInfo
			}
			wsConn.WriteObject("",proto.GateLoginServerAck, ackInfo)

		}else if msgName == proto.GateExitProxy{
			myGate.closeProxy(wsConn, msgProxyId)
		}else{
			routerToTarget(wsConn, msgName, msgProxyId, body)
		}
	}

}

func routerToTarget(wsConn* liNet.WsConnection, msgName string, msgProxyId string, body string){
	isAuth := false
	if proto.EnterLoginLoginReq == msgName || proto.EnterLoginRegisterReq == msgName{
		isAuth = true
	}else{
		r, err := wsConn.GetProperty("isAuth")
		if err == nil && r == true{
			isAuth = true
		}
	}

	if isAuth {
		isLive := false
		proxyClient, ok := myGate.proxyClient(wsConn, msgProxyId, app.GRouter)
		if ok {
			sendData  := []byte(body)
			if proxyClient.GetConn() != nil && proxyClient.GetConn().IsClose() == false{
				isLive = true
				proxyClient.GetConn().SetProperty("gateConn", wsConn)
				proxyClient.GetConn().SetProperty("proxy", msgProxyId)
				proxyClient.GetConn().SendMsg(msgName, sendData)
			}
		}

		if isLive == false{
			msg := fmt.Sprintf("%s ProxyClient not live", msgProxyId)
			utils.Log.Warn(msg)
			wsConn.WriteMessage(msgProxyId, proto.ProxyError, []byte(msg))
		}
	}else {
		wsConn.WriteMessage(msgProxyId, proto.AuthError, []byte(""))
	}

}


func ShutDown(){
	utils.Log.Info("ShutDown")
}

func main() {

	if len(os.Args) > 1 {
		cfgPath := os.Args[1]
		utils.GlobalObject.Load(cfgPath)
	}else{
		utils.GlobalObject.Load("conf/gate.json")
	}

	db.InitDataBase()

	go app.MasterClient(proto.ServerTypeGate)
	app.SetShutDownFunc(ShutDown)

	addr := fmt.Sprintf("%s:%d", utils.GlobalObject.AppConfig.Host, utils.GlobalObject.AppConfig.TcpPort)
	http.HandleFunc("/", wsHandler)
	http.ListenAndServe(addr, nil)

}