package main

import (
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db"
	"github.com/llr104/LiFrame/server/world"
	"github.com/llr104/LiFrame/utils"
	"os"
)


func main() {

	if len(os.Args) > 1 {
		cfgPath := os.Args[1]
		utils.GlobalObject.Load(cfgPath)
	}else{
		utils.GlobalObject.Load("conf/world.json")
	}

	db.Init()

	s := liNet.NewServer()
	s.AddRouter(&world.STS)
	s.AddRouter(&world.Enter)

	s.SetOnConnStart(world.ClientConnStart)
	s.SetOnConnStop(world.ClientConnStop)
	app.SetServer(s)
	app.SetShutDownFunc(world.ShutDown)

	go app.MasterClient(proto.ServerTypeWorld)

	s.Running()
}
