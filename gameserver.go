package main

import (
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db"
	"github.com/llr104/LiFrame/server/game"
	"github.com/llr104/LiFrame/server/gameutils"
	"github.com/llr104/LiFrame/utils"
	"os"
)


func main() {

	if len(os.Args) > 1 {
		cfgPath := os.Args[1]
		utils.GlobalObject.Load(cfgPath)
	}else{
		utils.GlobalObject.Load("conf/game.json")
	}

	db.Init()

	s := liNet.NewServer()
	s.AddRouter(&gameutils.STS)
	s.AddRouter(&game.Enter)

	s.SetOnConnStart(gameutils.ClientConnStart)
	s.SetOnConnStop(gameutils.ClientConnStop)
	app.SetShutDownFunc(gameutils.ShutDown)
	app.SetServer(s)

	go app.MasterClient(proto.ServerTypeGame)

	s.Running()
}
