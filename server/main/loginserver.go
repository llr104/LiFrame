package main

import (
	"LiFrame/core/liNet"
	"LiFrame/proto"
	"LiFrame/server/app"
	"LiFrame/server/db"
	"LiFrame/server/login"
	"LiFrame/utils"
)


func main() {

	utils.GlobalObject.Load("conf/login.json")
	db.InitDataBase()

	s := liNet.NewServer()
	s.AddRouter(&login.Enter)

	s.SetOnConnStop(login.ClientConnStop)
	s.SetOnConnStart(login.ClientConnStart)
	app.SetShutDownFunc(login.ShutDown)
	app.SetServer(s)
	go app.MasterClient(proto.ServerTypeLogin)
	s.Running()
}
