package utils

type ClientConfig struct {
	RemoteHost      string          //远程服务器主机ip
	RemoteTcpPort   int             //远程服务器主机端口号
	ClientName      string
	ClientId        string
}

type DBConfig struct {
	User			string
	Password        string
	Name            string
	Port            int
	IP              string
}

type HttpConfig struct {
	Port            int
	IP              string
}


type Config struct {
	Host      		string          //当前服务器主机IP
	TcpPort   		int             //当前服务器主机监听端口号
	ServerName      string          //当前服务器名称
	ServerId		string			//服务器id
	LogFile       	string 			//日志文件名称

	Master    		ClientConfig
	DataBase        DBConfig
	Http      		HttpConfig

	MaxPacketSize    uint32 //都需数据包的最大值
	MaxConn          int    //当前服务器主机允许的最大链接个数
	ServerWorkerSize uint32 //业务工作Worker池的数量
	MaxWorkerTaskLen uint32 //业务工作Worker对应负责的任务队列最大任务存储数量
}

func NewConfig() Config {
	c := Config{}
	c.Host = "0.0.0.0"
	c.TcpPort = 8000
	c.ServerName = "Default Server"
	c.ServerId = "Server1"
	c.LogFile = "./logout/run.log"


    c.MaxConn = 12000
	c.MaxPacketSize = 40960
	c.ServerWorkerSize = 2
	c.MaxWorkerTaskLen = 1024
	return c
}