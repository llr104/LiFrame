package app

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/utils"
	"time"
)

var (
	mPingTimerId       uint32
	mReportTimerId     uint32
	mServerListTimerId uint32
	
	MClient *liNet.Client
)

func mClientConnStart(conn liFace.IConnection){
	utils.Log.Info("mClientConnStart:%s", conn.RemoteAddr().String())

	//启动的时候快速同步几次，确保数据实时
	utils.Scheduler.NewTimerInterval(2*time.Second, 15, mReportTimer, []interface{}{conn})
	utils.Scheduler.NewTimerInterval(2*time.Second, 15, mServerListTimer, []interface{}{conn})

	//后续同步次数减少
	mReportTimerId, _ = utils.Scheduler.NewTimerInterval(30*time.Second, utils.IntervalForever, mReportTimer, []interface{}{conn})
	mPingTimerId, _ = utils.Scheduler.NewTimerInterval(15*time.Second, utils.IntervalForever, mPingTimer, []interface{}{conn})
	mServerListTimerId, _ = utils.Scheduler.NewTimerInterval(35*time.Second, utils.IntervalForever, mServerListTimer, []interface{}{conn})

}

func mClientConnStop(conn liFace.IConnection){

	tcp := conn.GetTcpNetWork()
	Name := tcp.GetName()
	Id := tcp.GetId()
	IP := tcp.GetHost()
	Port := tcp.GetPort()

	utils.Log.Info("mClientConnStop:%s,%s,%s:%d", Name, Id, IP,Port)
	utils.Scheduler.NewTimerAfter(5*time.Second, restartClientMaster, []interface{}{tcp})


	if mReportTimerId > 0{
		utils.Scheduler.CancelTimer(mReportTimerId)
		mReportTimerId = 0
	}

	if mPingTimerId > 0{
		utils.Scheduler.CancelTimer(mPingTimerId)
		mPingTimerId = 0
	}

	if mServerListTimerId > 0{
		utils.Scheduler.CancelTimer(mServerListTimerId)
		mServerListTimerId = 0
	}

	utils.Log.Info("mClientConnStop end")
}


func MasterClient(sType proto.ServerType) {
	client := utils.GlobalObject.AppConfig.Master
	if client.RemoteHost != ""{
		MClient = liNet.NewClient(client.ClientName, client.ClientId, client.RemoteHost, client.RemoteTcpPort)
		MClient.SetClientType(sType)
		MClient.SetOnConnStart(mClientConnStart)
		MClient.SetOnConnStop(mClientConnStop)
		MClient.AddRouter(&MClientRouter)
		MClient.Running()
	}
}


func restartClientMaster(v ...interface{}) {
	c := v[0].(*liNet.Client)
	MasterClient(c.GetClientType())
}

func mPingTimer(v ...interface{})  {

	conn := v[0].(liFace.IConnection)
	info := proto.PingPong{}
	info.CurTime = time.Now().Unix()

	data ,err := json.Marshal(info)
	if err == nil{
		conn.SendMsg(proto.EnterMasterPing, data)
	}else{
		utils.Log.Info("mReportTimer error:%s", err.Error())
	}
}

func mReportTimer(v ...interface{})  {
	conn := v[0].(liFace.IConnection)
	tcp := conn.GetTcpNetWork()
	c := tcp.(*liNet.Client)
	info := proto.ServerInfoReport{}
	info.LastTime = time.Now().Unix()
	info.Type = c.GetClientType()
	info.State = proto.ServerStateNormal
	info.IP = utils.GlobalObject.AppConfig.Host
	info.Port = utils.GlobalObject.AppConfig.TcpPort
	info.Id = utils.GlobalObject.AppConfig.ServerId
	info.Name = utils.GlobalObject.AppConfig.ServerName
	info.OnlineCnt = MClientData.GetOnlineCnt()

	data ,err := json.Marshal(info)
	if err == nil{
		conn.SendMsg(proto.EnterMasterServerInfoReport, data)
	}else{
		utils.Log.Info("mReportTimer error:%s", err.Error())
	}
}

func mServerListTimer(v ...interface{})  {

	conn := MClient.GetConn()
	if conn != nil{
		info := proto.ServerListReq{}
		info.CurTime = time.Now().Unix()

		data ,err := json.Marshal(info)
		if err == nil{
			conn.SendMsg(proto.EnterMasterServerListReq, data)
		}else{
			utils.Log.Info("mServerListTimer error:%s", err.Error())
		}
	}
}