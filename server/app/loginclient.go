package app

import (
	"encoding/json"
	"errors"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/utils"
	"time"
)

var (
	lPingTimerId uint32
)

func lClientConnStart(conn liFace.IConnection){
	utils.Log.Info("lClientConnStart:%s", conn.RemoteAddr().String())
	lPingTimerId, _ = utils.Scheduler.NewTimerInterval(15*time.Second,utils.IntervalForever, loginPingTimer, []interface{}{conn})
}

func lClientConnStop(conn liFace.IConnection){

	client := conn.GetTcpNetWork().(*liNet.Client)
	Name := client.GetName()
	Id := client.GetId()
	IP := client.GetHost()
	Port := client.GetPort()
	cType := client.GetClientType()

	utils.Log.Info("lClientConnStop:%s,%s,%s:%d,%d", Name, Id, IP, Port, cType)
	utils.Scheduler.NewTimerAfter(5*time.Second, restartLoginClient, []interface{}{Name, Id, IP, Port, cType})
	utils.Log.Info("lClientConnStop end")

	if lPingTimerId > 0{
		utils.Scheduler.CancelTimer(lPingTimerId)
		lPingTimerId = 0
	}
}

func LoginClient(clientName string, clientId string,
	remoteHost string, remotePort int, clientType proto.ServerType) (*liNet.Client, error){

	if remotePort > 0 && remoteHost != ""{
		var c *liNet.Client
		c = liNet.NewClient(clientName, clientId, remoteHost, remotePort)
		c.SetClientType(clientType)
		c.SetOnConnStart(lClientConnStart)
		c.SetOnConnStop(lClientConnStop)
		c.Running()
		return c, nil
	}
	return nil, errors.New("new LoginClient Error")
}


func restartLoginClient(v ...interface{}) {
	Name := v[0].(string)
	Id := v[1].(string)
	Ip := v[2].(string)
	Port := v[3].(int)
	cType := v[4].(proto.ServerType)

	LoginClient(Name, Id, Ip, Port, cType)
}

func loginPingTimer(v ...interface{})  {

	conn := v[0].(liFace.IConnection)
	info := proto.PingPong{}
	info.CurTime = time.Now().Unix()

	data ,err := json.Marshal(info)
	if err == nil{
		conn.RpcCall(proto.SystemPing, data, nil, nil)
	}else{
		utils.Log.Info("loginPingTimer error:%s", err.Error())
	}
}

