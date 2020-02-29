package main

import (
	"fmt"
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db"
	"github.com/llr104/LiFrame/server/master"
	"github.com/llr104/LiFrame/utils"
	"net/http"
	"os"
)


func main() {

	if len(os.Args) > 1 {
		cfgPath := os.Args[1]
		utils.GlobalObject.Load(cfgPath)
	}else{
		utils.GlobalObject.Load("conf/master.json")
	}
	db.Init()

	/*
	http 监测服务器状态
	*/
	httpCfg := utils.GlobalObject.AppConfig.Http
	addr := fmt.Sprintf("%s:%d", httpCfg.IP, httpCfg.Port)
	http.HandleFunc("/", master.HomeHandler)
	http.HandleFunc("/state", master.StateHandler)
	http.HandleFunc("/shutdown", master.ShutdownHandler)

	utils.Log.Info("start monitor:%s", addr)
	go http.ListenAndServe(addr, nil)

	s := liNet.NewServer()
	s.AddRouter(&master.STS)
	s.SetOnConnStop(master.ClientConnStop)
	s.SetOnConnStart(master.ClientConnStart)
	app.SetShutDownFunc(master.ShutDown)
	app.SetServer(s)
	s.Running()


}
