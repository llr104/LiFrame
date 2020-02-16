package world

import (
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/app"
	"github.com/llr104/LiFrame/utils"
	"time"
)

type toLogin struct {
	clientMap map[string] *liNet.Client //serverId,liNet.Client
}

var W2Login toLogin
var clientTimer uint32


func init() {
	W2Login = toLogin{clientMap: make(map[string] *liNet.Client),}
	clientTimer, _ = utils.Scheduler.NewTimerInterval(5*time.Second, utils.IntervalForever, checkLoginClient, []interface{}{})
}

func (s*toLogin)GetLoginClient(appId string) (*liNet.Client, bool) {
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
						c.AddRouter(&sts{})
						W2Login.clientMap[val.Id] = c
					}
				}
			}else{
				c, err := app.LoginClient(val.Name, val.Id, val.IP, val.Port, val.Type)
				if err == nil{
					c.AddRouter(&sts{})
					W2Login.clientMap[val.Id] = c
				}
			}
		}
	}
}
