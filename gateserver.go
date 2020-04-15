package main

import (
	"bytes"
	"compress/gzip"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db"
	"github.com/llr104/LiFrame/server/gate"
	"github.com/llr104/LiFrame/utils"
	"github.com/thinkoner/openssl"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var cid uint64 = 0

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
	gate.MyGate.ConnectExit(wsConn)
	utils.Log.Debug("handleOnClose wsCount:%d", app.MClientData.GetOnlineCnt())
}

func handleWsMessage(wsConn *liNet.WsConnection, req *liNet.WsMessageReq, rsp* liNet.WsMessageRsp) {
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
	if len(msgArr) == 4{

		msgName := msgArr[0]
		msgProxyId := msgArr[1]
		seq := msgArr[2]
		req.Seq, _ = strconv.Atoi(seq)
		body := msgArr[3]

		if msgName == proto.GateHandshake{
			if body == ""{
				utils.Log.Info("不是断线重连")
				handshakeId := gate.MyGate.ConnectEnter(wsConn)
				rsp.Seq = req.Seq
				rsp.Data = handshakeId
				rsp.FuncName = proto.GateHandshake
				rsp.ProxyName = ""
			}else{
				utils.Log.Info("是断线重连")
				handshakeId := gate.MyGate.Reconnect(wsConn, body)
				rsp.Seq = req.Seq
				rsp.Data = handshakeId
				rsp.FuncName = proto.GateHandshake
				rsp.ProxyName = ""
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
				utils.Log.Info("gateServer.LoginServerReq error:%s", err.Error())
			}else{
				ackInfo.Code = proto.Code_Success
				ackInfo.ServerInfo = serverInfo
			}
			rsp.FuncName = proto.GateLoginServerAck
			rsp.ProxyName = ""
			rsp.Data = ackInfo
			rsp.Seq = req.Seq

		}else if msgName == proto.GateExitProxy{
			gate.MyGate.CloseProxy(wsConn, msgProxyId)
		}else{
			rsp.Seq = req.Seq
			routerToTarget(wsConn, msgName, msgProxyId, body, rsp)
		}
	}

}

func routerToTarget(wsConn* liNet.WsConnection, msgName string, msgProxyId string, body string, rsp* liNet.WsMessageRsp){
	isAuth := false
	if proto.EnterLoginLoginReq == msgName || proto.EnterLoginRegisterReq == msgName{
		isAuth = true
	}else{
		_, err := wsConn.GetProperty("session")
		if err == nil {
			isAuth = true
		}
	}

	if isAuth {
		isLive := false
		proxyClient, err := gate.MyGate.ProxyClient(wsConn, msgProxyId, gate.Router)
		if err == nil {
			sendData  := []byte(body)
			if proxyClient.GetConn() != nil && proxyClient.GetConn().IsClose() == false{
				isLive = true
				proxyClient.GetConn().SetProperty("gateConn", wsConn)
				proxyClient.GetConn().SetProperty("proxy", msgProxyId)
				proxyClient.GetConn().SendMsg(msgName, sendData)
			}
		}else{

			utils.Log.Warn(err.Error())

			rsp.ProxyName = msgProxyId
			rsp.FuncName = proto.ProxyError
			rsp.Data = err.Error()
		}

		if isLive == false{
			msg := fmt.Sprintf("%s ProxyClient not live", msgProxyId)
			utils.Log.Warn(msg)

			rsp.ProxyName = msgProxyId
			rsp.FuncName = proto.ProxyError
			rsp.Data = msg
		}
	}else {
		rsp.ProxyName = msgProxyId
		rsp.FuncName = proto.AuthError
		rsp.Data = ""
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

	db.Init()

	go app.MasterClient(proto.ServerTypeGate)
	app.SetShutDownFunc(ShutDown)

	addr := fmt.Sprintf("%s:%d", utils.GlobalObject.AppConfig.Host, utils.GlobalObject.AppConfig.TcpPort)
	http.HandleFunc("/", wsHandler)
	http.ListenAndServe(addr, nil)

}