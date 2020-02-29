package main

import (
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/server/gameslg"
	"github.com/llr104/LiFrame/server/gameutils"
	"github.com/llr104/LiFrame/utils"
	"os"
)


func main() {

	defaultXlsxDir := "conf/xlsx"
	if len(os.Args) == 3 {
		utils.GlobalObject.Load(os.Args[1])
		defaultXlsxDir = os.Args[2]
	}else if len(os.Args) == 2{
		utils.GlobalObject.Load(os.Args[1])
	}else{
		utils.GlobalObject.Load("conf/gameslg.json")
	}

	gameslg.Init(defaultXlsxDir)

	s := liNet.NewServer()
	s.AddRouter(&gameutils.STS)
	s.AddRouter(&gameslg.Enter)
	s.AddRouter(&gameslg.CreateRole)
	s.AddRouter(&gameslg.MainCity)
	s.AddRouter(&gameslg.NPC)
	s.AddRouter(&gameslg.WorldMap)



	s.SetOnConnStart(gameutils.ClientConnStart)
	s.SetOnConnStop(gameutils.ClientConnStop)
	app.SetShutDownFunc(gameutils.ShutDown)
	app.SetServer(s)

	go app.MasterClient(proto.ServerTypeGame)

	s.Running()
}
