package game

import (
	"github.com/llr104/LiFrame/proto"
)

const protoHeartBeatReq = "heartBeatReq"
const protoHeartBeatAck = "heartBeatAck"
const protoLogoutReq = "logoutReq"
const protoLogoutAck = "logoutAck"

const protoSceneListReq = "sceneListReq"
const protoSceneListAck = "sceneListAck"

const protoEnterSceneReq = "enterSceneReq"
const protoEnterSceneAck = "enterSceneAck"
const protoExitSceneReq = "exitSceneReq"
const protoExitSceneAck = "exitSceneAck"
const protoSceneReq = "sceneReq"
const protoSceneAck = "sceneAck"
const protoMoveReq = "moveReq"
const protoMovePush = "movePush"
const protoAttackReq = "attackReq"
const protoAttackPush = "attackPush"
const protoMonsterPush = "monsterPush"
const protoUserPush = "userPush"


type sceneData struct {
	Players 		map[uint32] *player 	`json:"players"`
	Monsters        map[uint32] *monster 	`json:"monsters"`
}

type heartBeat struct {
	ClientTimeStamp int64		`json:"clientTimeStamp"`
	ServerTimeStamp int64		`json:"serverTimeStamp"`
}

type sceneListAck struct {
	SceneId     []int      	`json:"sceneId"`
	SceneName   []string     `json:"sceneName"`
}

type enterSceneReq struct {
	SceneId     int      	`json:"sceneId"`
}

type enterSceneAck struct {
	proto.BaseAck
	SceneId     int      `json:"sceneId"`
	SceneName   string    `json:"sceneName"`
}

type exitSceneReq struct {
	SceneId     int      `json:"sceneId"`
}

type exitSceneAck struct {
	proto.BaseAck
	SceneId     int      `json:"sceneId"`
}

type monsterPush struct {
	Monsters        map[uint32] *monster 	`json:"monsters"`
}

type userPush struct {
	Players 		map[uint32] *player 	`json:"players"`
}

type move struct {
	SX		int      `json:"sx"`
	SY      int		 `json:"sy"`
	TX      int      `json:"tx"`
	TY      int		 `json:"ty"`
	UserId  uint32   `json:"userId"`
}

type attackReq struct {
	UserId  	uint32   `json:"userId"`
	MonsterId  	uint32   `json:"monsterId"`
	Hurt    	int		 `json:"hurt"`
}

type attackPush struct {
	UserId  	uint32   `json:"userId"`
	MonsterId  	uint32   `json:"monsterId"`
	Hurt    	int		 `json:"hurt"`
	MonsterHp   int      `json:"monsterHp"`
}

