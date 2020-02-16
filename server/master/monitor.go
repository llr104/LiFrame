package master

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/utils"
	"net/http"
	"os"
	"time"
)

func HomeHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("well come to monitor!"))
}

func StateHandler(resp http.ResponseWriter, req *http.Request) {
	m := STS.getServerMap()
	data, _ := json.Marshal(m)
	resp.Write(data)
}

func ShutdownHandler(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(req.Form["isAll"]) > 0{
		utils.Log.Info("Shutdown all")
		ser := app.GetServer()
		resp.Write([]byte("ok"))

		utils.Scheduler.NewTimerInterval(1*time.Second, 5, shutdownTimer, []interface{}{ser})
		utils.Scheduler.NewTimerAfter(6*time.Second, shutdown, []interface{}{})

	}else{
		resp.Write([]byte("error"))
	}
}

func shutdownTimer(v ...interface{}) {
	ser := v[0].(liFace.INetWork)
	mgr := ser.GetConnMgr()
	mgr.BroadcastMsg(proto.SystemShutDown, nil)
}

func shutdown(v ...interface{}) {
	os.Exit(0)
}