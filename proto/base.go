package proto

const ProxyError  = "proxyError"
const AuthError  = "authError"
const GameEnterGameReq  = "enterGameReq"
const GameEnterGameAck  = "enterGameAck"

const SystemShutDown = "System.ShutDown"
const SystemSessionOnlineOrOffLine = "System.SessionOnlineOrOffLine"
const SystemPing = "System.Ping"
const SystemPong = "System.Pong"
const SystemServerInfoReport = "System.ServerInfoReport"
const SystemServerListReq = "System.ServerListReq"
const SystemServerListAck = "System.ServerListAck"
const SystemCheckSessionReq = "System.CheckSessionReq"
const SystemCheckSessionAck = "System.CheckSessionAck"
const SystemSessionUpdateReq = "System.SessionUpdateReq"
const SystemSessionUpdateAck = "System.SessionUpdateAck"

const EnterWorldJoinWorldReq = "EnterWorld.JoinWorldReq"
const EnterWorldJoinWorldAck = "EnterWorld.JoinWorldAck"
const EnterWorldUserInfoReq = "EnterWorld.UserInfoReq"
const EnterWorldUserInfoAck = "EnterWorld.UserInfoAck"
const EnterWorldUserLogoutReq = "EnterWorld.UserLogoutReq"
const EnterWorldGameServersAck = "EnterWorld.GameServersAck"
const EnterWorldUserLogoutAck = "EnterWorld.UserLogoutAck"

const EnterLoginDistributeWorldAck = "EnterLogin.DistributeWorldAck"
const EnterLoginRegisterAck = "EnterLogin.RegisterAck"
const EnterLoginLoginAck  = "EnterLogin.LoginAck"
const EnterLoginLoginReq = "EnterLogin.LoginReq"
const EnterLoginRegisterReq = "EnterLogin.RegisterReq"
const EnterLoginDistributeServerReq = "EnterLogin.DistributeServerReq"

const GateHandshake = "handshake"
const GateLoginServerReq = "gate.LoginServerReq"
const GateLoginServerAck = "gate.LoginServerAck"
const GateExitProxy = "gate.ExitProxy"

type BaseAck struct {
	Code   int
}
