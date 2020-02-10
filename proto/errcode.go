package proto

const (
	Code_Success = iota   	//0
	Code_Illegal
	Code_User_Not_Exist
	Code_User_Error 		//账号或密码错误
	Code_User_Exist
	Code_User_Forbid
	Code_Not_Server
	Code_Session_Error
	Code_EnterGameError
	Code_EnterSceneError
	Code_ExitSceneError
)