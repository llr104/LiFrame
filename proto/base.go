package proto

const ProxyError  = "proxyError"
const AuthError  = "authError"
const GateHandshake = "handshake"
const GateLoginServerReq = "gate.LoginServerReq"
const GateExitProxy = "gate.ExitProxy"

const SystemShutDown = "System.ShutDown"
const SystemPing = "System.Ping"
const SystemPong = "System.Pong"
const SystemServerInfoReport = "System.ServerInfoReport"
const SystemServerListReq = "System.ServerListReq"
const SystemCheckSessionReq = "System.CheckSessionReq"
const SystemSessionUpdateReq = "System.SessionUpdateReq"
const SystemUserOnOrOffReq = "System.UserOnOrOffReq"

const EnterWorldJoinWorldReq = "EnterWorld.JoinWorldReq"
const EnterWorldSession = "EnterWorld.Session"

const EnterLoginLoginReq = "EnterLogin.LoginReq"
const EnterLoginRegisterReq = "EnterLogin.RegisterReq"


const GameEnterGameReq  = "enterGameReq"
const GameEnterGameAck  = "enterGameAck"

type BaseAck struct {
	Code   int
}
