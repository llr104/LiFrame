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
	m := Enter.getServerMap()
	data, _ := json.Marshal(m)
	resp.Write(data)
}

func ShutdownHandler(resp http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	if len(req.Form["isAll"]) > 0{
		utils.Log.Info("Shutdown all")
		ser := app.GetServer()
		n := ser.(liFace.INetWork)
		mgr := n.GetConnMgr()
		mgr.BroadcastMsg(proto.MasterClientShutDown, nil)
		resp.Write([]byte("ok"))

		utils.Scheduler.NewTimerAfter(5*time.Second, shutdown, []interface{}{})
		os.Exit(0)
	}else{
		resp.Write([]byte("error"))
	}
}

func shutdown(v ...interface{}) {
	os.Exit(0)
}