package utils

import (
	"LiFrame/core/logs"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
)

/*
	存储一切有关Zinx框架的全局参数，供其他模块使用
	一些参数也可以通过 用户根据 LiFrame.json来配置
*/
type GlobalObj struct {

	AppConfig        Config
	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	ServerWorkerSize uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
	MaxMsgChanLen    uint32 //SendBuffMsg发送消息的缓冲最大长度
}

/*
	定义一个全局的对象
*/
var GlobalObject *GlobalObj
var Log *logs.LiLogger


//读取用户的配置文件
func (g *GlobalObj) Load(configFile string) {

	if confFileExists, _ := PathExists(configFile); confFileExists != true {
		text := fmt.Sprintf("Config File %s is not exist!!", configFile)
		Log.Error(text)
		panic(errors.New(text))
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}
	//将json数据解析到struct中
	err = json.Unmarshal(data, &g.AppConfig)
	if err != nil {
		Log.Error("load config error")
		panic(err)
	}

	Log.Info("Config:%v", g.AppConfig)

	if g.AppConfig.LogFile != ""{
		logCfg := fmt.Sprintf(`{"filename":"%s"}`,g.AppConfig.LogFile)
		Log.SetLogger(logs.AdapterFile, logCfg)
	}

}

/*
	提供init方法，默认加载
*/
func init() {

	Log = logs.GetLiLogger()
	Log.SetLogger(logs.AdapterConsole)
	Log.EnableFuncCallDepth(true)
	//初始化GlobalObject变量，设置一些默认值
	GlobalObject = &GlobalObj{
		MaxConn:          12000,
		MaxPacketSize:    40960,
		ServerWorkerSize: 2,
		MaxWorkerTaskLen: 1024,
		MaxMsgChanLen:    1024,
		AppConfig:        NewConfig(),
	}

}
