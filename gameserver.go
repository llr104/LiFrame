package main

import (
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/db"
	"github.com/llr104/LiFrame/server/game"
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

	db.InitDataBase()

	s := liNet.NewServer()
	s.AddRouter(&game.Enter)

	s.SetOnConnStart(game.ClientConnStart)
	s.SetOnConnStop(game.ClientConnStop)
	app.SetShutDownFunc(game.ShutDown)
	app.SetServer(s)

	go app.MasterClient(proto.ServerTypeGame)

	s.Running()
}
