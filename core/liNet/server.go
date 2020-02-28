package liNet

import (
	"fmt"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/utils"
	"net"
	"os"
)

//iServer 接口实现，定义一个Server服务类
type Server struct {
	//服务器的名称
	name 		string
	ipVersion 	string
	ip 			string
	port     	int
	id       	string
	listener 	*net.TCPListener
	exit     	chan bool
	//当前Server的消息管理模块，用来绑定msgName和对应的处理方法
	msgHandler liFace.IMsgHandle
	//当前Server的链接管理器
	connMgr liFace.IConnManager
	//该Server的连接创建时Hook函数
	onConnStart func(conn liFace.IConnection)
	//该Server的连接断开时的Hook函数
	onConnStop func(conn liFace.IConnection)
}

/*
  创建一个服务器句柄
 */
func NewServer () *Server{

	s:= &Server {
		id:         utils.GlobalObject.AppConfig.ServerId,
		name:       utils.GlobalObject.AppConfig.ServerName,
		ipVersion:  "tcp4",
		ip:         utils.GlobalObject.AppConfig.Host,
		port:       utils.GlobalObject.AppConfig.TcpPort,
		msgHandler: NewMsgHandle(utils.GlobalObject.ServerWorkerSize),
		connMgr:    NewConnManager(),
		exit:       make(chan bool, 1),
	}
	return s
}
//============== 实现 liFace.INetWork 里的全部接口方法 ========

func (s *Server) GetName() string{
	return s.name
}

func (s *Server) GetId() string{
	return s.id
}

func (s *Server) GetHost()string{
	return s.ip
}

func (s *Server) GetPort() int{
	return s.port
}

//开启网络服务
func (s *Server) Start() {
	utils.Log.Info("[START] Server name: %s, listenner at ip: %s, port %d is starting", s.name, s.ip, s.port)

	//开启一个go去做服务端Linster业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.ipVersion, fmt.Sprintf("%s:%d", s.ip, s.port))
		if err != nil {
			utils.Log.Error("resolve tcp addr err: %s", err.Error())
			os.Exit(0)
		}

		//2 监听服务器地址
		listener, err:= net.ListenTCP(s.ipVersion, addr)
		if err != nil {
			utils.Log.Error("listen %s err %s", s.ipVersion,  err.Error())
			os.Exit(0)
		}else{
			s.listener = listener
		}

		//已经监听成功
		utils.Log.Info("start server %s succ, now listenning...",s.name)

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		//3 启动server网络连接业务
		for {

			//3.1 阻塞等待客户端建立连接请求
			conn, err := listener.AcceptTCP()
			if err != nil {
				utils.Log.Info("Accept err %s ", err.Error())
				return
			}
			utils.Log.Info("Get conn remote addr = %s", conn.RemoteAddr().String())

			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.connMgr.Len() >= utils.GlobalObject.MaxConn {

				utils.Log.Error("too much connect %d", s.connMgr.Len())
				conn.Close()
				continue
			}

			utils.Log.Info("Get conn cid %d", cid)
			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid ++

			//3.4 启动当前链接的处理业务
			go dealConn.Start()
		}
	}()
}

//停止服务
func (s *Server) Stop() {
	utils.Log.Info("[STOP] server name %s" , s.name)

	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.msgHandler.StopWorkerPool()
	s.connMgr.ClearConn()
	s.listener.Close()
	s.exit <- true
}

//运行服务
func (s *Server) Running() {
	s.Start()

	//TODO Server.Running() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞,否则主Go退出， listenner的go将会退出
	select{
		case exit:=<-s.exit:
			if exit {
				return
			}
	}

}

//路由功能：给当前服务注册一个路由业务方法，供客户端链接处理使用
func (s *Server) AddRouter(router liFace.IRouter) {
	s.msgHandler.AddRouter(router)
	router.After()
}

//得到链接管理
func (s *Server) GetConnMgr() liFace.IConnManager {
	return s.connMgr
}

//设置该Server的连接创建时Hook函数
func (s *Server) SetOnConnStart(hookFunc func (liFace.IConnection)) {
	s.onConnStart = hookFunc
}

//设置该Server的连接断开时的Hook函数
func (s *Server) SetOnConnStop(hookFunc func (liFace.IConnection)) {
	s.onConnStop = hookFunc
}

//调用连接OnConnStart Hook函数
func (s *Server) CallOnConnStart(conn liFace.IConnection) {
	if s.onConnStart != nil {
		utils.Log.Info("---> CallOnConnStart....")
		s.onConnStart(conn)
	}
}

//调用连接OnConnStop Hook函数
func (s *Server) CallOnConnStop(conn liFace.IConnection) {
	if s.onConnStop != nil {
		utils.Log.Info("---> CallOnConnStop....")
		s.onConnStop(conn)
	}
}





