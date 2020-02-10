package main

import (
	"LiFrame/core/liNet"
	"LiFrame/server/app"
	"LiFrame/server/db"
	"LiFrame/server/master"
	"LiFrame/utils"
	"fmt"
	"net/http"
)


func main() {

	utils.GlobalObject.Load("conf/master.json")
	db.InitDataBase()

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
	s.AddRouter(&master.Enter)
	s.SetOnConnStop(master.ClientConnStop)
	s.SetOnConnStart(master.ClientConnStart)
	app.SetShutDownFunc(master.ShutDown)
	app.SetServer(s)
	s.Running()


}
