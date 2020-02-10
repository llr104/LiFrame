package world

import (
	"LiFrame/core/liNet"
	"LiFrame/proto"
	"LiFrame/server/app"
	"LiFrame/utils"
	"time"
)

type ToLogin struct {
	clientMap map[string] *liNet.Client //serverId,liNet.Client
}

var W2Login ToLogin
var clientTimer uint32


func init() {
	W2Login = ToLogin{clientMap: make(map[string] *liNet.Client),}
	clientTimer, _ = utils.Scheduler.NewTimerInterval(5*time.Second, utils.IntervalForever, checkLoginClient, []interface{}{})
}

func (s*ToLogin)GetLoginClient(appId string) (*liNet.Client, bool) {
	o, ok := W2Login.clientMap[appId]
	return o,ok
}

func checkLoginClient(v ...interface{}) {
	serverMap := app.ServerMgr.GetServerMap()
	for _, val := range serverMap {
		if val.Type == proto.ServerTypeLogin{
			o, ok := W2Login.clientMap[val.Id]
			if ok {
				if val.IP != o.GetHost() || val.Port != o.GetPort(){
					o.Stop()
					c, err := app.LoginClient(val.Name, val.Id, val.IP, val.Port, val.Type)
					if err == nil{
						c.AddRouter(&EnterWorld{})
						W2Login.clientMap[val.Id] = c
					}
				}
			}else{
				c, err := app.LoginClient(val.Name, val.Id, val.IP, val.Port, val.Type)
				if err == nil{
					c.AddRouter(&EnterWorld{})
					W2Login.clientMap[val.Id] = c
				}
			}
		}
	}
}
