package game

type iScene interface {
	EnterScene(userId uint32) bool
	ExitScene(userId uint32) bool
	GameMessage(userId uint32, msgName string, data []byte)
	UserOffLine(userId uint32) bool
	SendMessageToUser(userId uint32, msgName string, msg interface{})
    SendMessageToAll(msgName string, msg interface{})
	Name()string
	Id()int
	SetName(string)
	SetId(int)
}