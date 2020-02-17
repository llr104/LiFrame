package game

import (
	"encoding/json"
	"github.com/llr104/LiFrame/core/liFace"
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/gameutils"
	"github.com/llr104/LiFrame/utils"
)

var game mainLogic

func init() {
	game = mainLogic{make(map[int]iScene), false}
	gameutils.STS.SetGame(&game)

	s1 := NewScene1()
	s1.SetId(0)
	s1.SetName("场景1")
	game.scenes[0] = s1

	s2 := NewScene1()
	s2.SetId(1)
	s2.SetName("场景2")
	game.scenes[1] = s2
}

type mainLogic struct {
	scenes     map[int]iScene
	isShutDown bool
}

func (s *mainLogic) enterGame(req proto.EnterGameReq) bool{
	if s.isShutDown {
		return false
	}

	return true
}

func (s *mainLogic) enterScene(userId uint32, sceneId int, conn liFace.IConnection) bool{

	ea := enterSceneAck{}
	ea.SceneId = sceneId

	scene, ok := s.scenes[sceneId]
	if ok == false {
		ea.Code = proto.Code_EnterSceneError
		data, _ := json.Marshal(ea)
		conn.SendMsg(protoEnterSceneAck, data)
		return false
	}

	isIn, state := gameutils.GUserMgr.UserIsIn(userId)
	if isIn{
		if state.SceneId != sceneId{
			t := s.scenes[state.SceneId]
			t.ExitScene(userId)
		}
	}
	ok = scene.EnterScene(userId)
	if ok {
		ea.Code = proto.Code_Success
		ea.SceneName = scene.Name()
		gameutils.GUserMgr.UserChangeState(userId, gameutils.GUserStateOnline, sceneId, conn)
	}else{
		ea.Code = proto.Code_EnterSceneError
	}

	data, _ := json.Marshal(ea)
	conn.SendMsg(protoEnterSceneAck, data)

	return ok
}

func (s *mainLogic) exitScene(userId uint32, sceneId int, conn liFace.IConnection){
	scene, ok := s.scenes[sceneId]
	ea := exitSceneAck{}
	ea.SceneId = sceneId

	if ok == false {
		ea.Code = proto.Code_ExitSceneError
		data, _ := json.Marshal(ea)
		conn.SendMsg(protoExitSceneAck, data)
		return
	}

	isIn, state := gameutils.GUserMgr.UserIsIn(userId)
	if isIn{
		if state.SceneId != sceneId{
			ea.Code = proto.Code_ExitSceneError
		}else{
			scene.ExitScene(userId)
			gameutils.GUserMgr.UserChangeState(userId, gameutils.GUserStateLeave, -1, conn)
			ea.Code = proto.Code_Success
		}
	}else{
		ea.Code = proto.Code_ExitSceneError
	}

	data, _ := json.Marshal(ea)
	conn.SendMsg(protoExitSceneAck, data)
}

func (s *mainLogic) gameMessage(userId uint32, msgName string, data []byte, conn liFace.IConnection){
	if s.isShutDown {
		return
	}

	if msgName == protoSceneListReq{
		a := sceneListAck{}
		a.SceneId = make([]int, len(s.scenes))
		a.SceneName = make([]string, len(s.scenes))
		for k, v := range s.scenes{
			a.SceneName[k] = v.Name()
			a.SceneId[k] = v.Id()
		}
		data, _ := json.Marshal(a)
		conn.SendMsg(protoSceneListAck, data)

	} else if msgName == protoEnterSceneReq{
		e := enterSceneReq{}
		json.Unmarshal(data, &e)
		s.enterScene(userId, e.SceneId, conn)
	}else if msgName == protoExitSceneReq{
		e := exitSceneReq{}
		json.Unmarshal(data, &e)
		s.exitScene(userId, e.SceneId, conn)
	}else{
		if ok, state := gameutils.GUserMgr.UserIsIn(userId); ok {
			if t, isOk := s.scenes[state.SceneId]; isOk{
				t.GameMessage(userId, msgName, data)
			}
		}
	}
}



/*
返回true:用户离开了游戏，返回false:用户断线，保留用户的游戏状态
*/
func (s *mainLogic) UserOffLine(userId uint32) bool{
	if ok, state := gameutils.GUserMgr.UserIsIn(userId); ok {
		if t, isOk := s.scenes[state.SceneId]; isOk{
			return t.UserOffLine(userId)
		}
	}
	return true
}

func (s *mainLogic) UserOnLine(userId uint32){
	utils.Log.Info("UserOnLine: %d", userId)
}

func (s *mainLogic) UserLogout(userId uint32) bool{
	return s.UserOffLine(userId)
}


func (s *mainLogic) ShutDown(){
	utils.Log.Info("ShutDown")
}
