package game

import (
	"github.com/llr104/LiFrame/core/liFace"
)

var GUserMgr userMgr

func init() {
	GUserMgr = userMgr{make(map[uint32]userState)}
}

type userMgr struct {
	states 		map[uint32]userState
}

func (s *userMgr) UserChangeState(userId uint32, state int, sceneId int, conn liFace.IConnection){

	if state != GUserStateLeave {
		s.states[userId] = userState{UserId: userId, State:state, SceneId: sceneId, Conn:conn}
	}else{
		delete(s.states, userId)
	}
}

func (s *userMgr) UserIsIn(userId uint32) (bool, userState){
	state, ok := s.states[userId]
	return ok, state
}


