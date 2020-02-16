package proto

const ProxyError  = "proxyError"
const AuthError  = "authError"
const GateHandshake = "handshake"
const GateLoginServerReq = "gate.LoginServerReq"
const GateLoginServerAck = "gate.LoginServerAck"
const GateExitProxy = "gate.ExitProxy"

const SystemShutDown = "System.ShutDown"
const SystemPing = "System.Ping"
const SystemPong = "System.Pong"
const SystemServerInfoReport = "System.ServerInfoReport"
const SystemServerListReq = "System.ServerListReq"
const SystemServerListAck = "System.ServerListAck"
const SystemCheckSessionReq = "System.CheckSessionReq"
const SystemCheckSessionAck = "System.CheckSessionAck"
const SystemSessionUpdateReq = "System.SessionUpdateReq"
const SystemSessionUpdateAck = "System.SessionUpdateAck"
const SystemUserOnOrOffReq = "System.UserOnOrOffReq"
const SystemUserOnOrOffAck = "System.UserOnOrOffAck"

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

const GameEnterGameReq  = "enterGameReq"
const GameEnterGameAck  = "enterGameAck"

type BaseAck struct {
	Code   int
}
