package main

import (
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db"
	"github.com/llr104/LiFrame/utils"
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/thinkoner/openssl"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)


// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var cid uint64 = 0
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

		/*
		检测授权是否成功才转发消息
		*/
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
			proxy,err := app.ServerMgr.GetProxy(msgProxyId)
			if err == nil {
				wsConn.CloseProxy(proxy)
			}else{
				utils.Log.Warn("%s", err.Error())
			}
		}else{
			isAuth := false
			if proto.EnterLoginLoginReq == msgName || proto.EnterLoginRegisterReq == msgName{
				isAuth = true
			}else{
				r, err := wsConn.GetProperty("isAuth")
				if err == nil && r == true{
					isAuth = true
				}
			}

			proxy,err := app.ServerMgr.GetProxy(msgProxyId)
			if err != nil{
				utils.Log.Warn("%s", err.Error())
				wsConn.WriteMessage(msgProxyId, proto.ProxyError, []byte(err.Error()))
				return
			}

			if isAuth {
				isLive := false
				proxyClient, ok := wsConn.ProxyClient(proxy, app.GRouter)
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