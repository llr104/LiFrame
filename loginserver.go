package main

import (
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db"
	"github.com/llr104/LiFrame/server/login"
	"github.com/llr104/LiFrame/utils"
	"os"
)


func main() {
	if len(os.Args) > 1 {
		cfgPath := os.Args[1]
		utils.GlobalObject.Load(cfgPath)
	}else{
		utils.GlobalObject.Load("conf/login.json")
	}

	db.Init()

	s := liNet.NewServer()
	s.AddRouter(&login.STS)
	s.AddRouter(&login.Enter)

	s.SetOnConnStop(login.ClientConnStop)
	s.SetOnConnStart(login.ClientConnStart)
	app.SetShutDownFunc(login.ShutDown)
	app.SetServer(s)
	go app.MasterClient(proto.ServerTypeLogin)
	s.Running()
}
