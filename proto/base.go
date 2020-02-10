package proto

const ProxyError  = "proxyError"
const AuthError  = "authError"
const GameEnterGameReq  = "enterGameReq"
const GameEnterGameAck  = "enterGameAck"

const MasterClientServerListAck = "MasterClient.ServerListAck"
const MasterClientPong  = "MasterClient.Pong"
const MasterClientShutDown  = "MasterClient.ShutDown"

const EnterLoginSessionUpdateReq = "EnterLogin.SessionUpdateReq"
const EnterLoginPing = "EnterLogin.Ping"
const EnterLoginCheckSessionReq = "EnterLogin.CheckSessionReq"
const EnterLoginDistributeWorldAck = "EnterLogin.DistributeWorldAck"
const EnterLoginSessionUpdateAck = "EnterLogin.SessionUpdateAck"
const EnterLoginRegisterAck = "EnterLogin.RegisterAck"
const EnterLoginLoginAck  = "EnterLogin.LoginAck"
const EnterLoginLoginReq = "EnterLogin.LoginReq"
const EnterLoginRegisterReq = "EnterLogin.RegisterReq"
const EnterLoginDistributeServerReq = "EnterLogin.DistributeServerReq"

const LoginClientPong = "LoginClient.Pong"
const EnterMasterPing = "EnterMaster.Ping"

const EnterMasterServerInfoReport = "EnterMaster.ServerInfoReport"
const EnterMasterServerListReq  = "EnterMaster.ServerListReq"
const EnterWorldCheckSessionAck = "EnterWorld.CheckSessionAck"
const EnterWorldJoinWorldReq  = "EnterWorld.JoinWorldReq"
const EnterWorldJoinWorldAck = "EnterWorld.JoinWorldAck"
const CommonWorldGameScenesAck = "CommonWorld.GameScenesAck"
const CommonWorldUserLogoutAck = "CommonWorld.UserLogoutAck"
const CommonWorldUserInfoReq = "CommonWorld.UserInfoReq"
const CommonWorldUserInfoAck = "CommonWorld.UserInfoAck"
const CommonWorldUserLogoutReq = "CommonWorld.UserLogoutReq"

const GateLoginServerReq = "gate.LoginServerReq"
const GateLoginServerAck = "gate.LoginServerAck"
const GateExitProxy = "gate.ExitProxy"
type BaseAck struct {
	Code   int
}
