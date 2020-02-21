package game

import (
	"encoding/json"
	"github.com/llr104/LiFrame/server/db/dbobject"
	"math/rand"
	"sync"
	"time"
)

type sceneBase struct {
	maxUser  	int
	users    	map[uint32] *dbobject.User
	players  	map[uint32] *player
	monsters 	map[uint32] *monster
	curId    	uint32
	lock     	sync.Mutex
	sceneName 	string
	sceneId     int
}

func  (s *sceneBase) init() {
	s.users = make(map[uint32] *dbobject.User)
	s.players = make(map[uint32] *player)
	s.monsters = make(map[uint32] *monster)
}

func (s *sceneBase) Name() string{
	return s.sceneName
}

func (s *sceneBase) Id() int{
	return s.sceneId
}

func (s *sceneBase) SetName(n string){
	s.sceneName = n
}

func (s *sceneBase) SetId(id int){
	s.sceneId = id
}

func (s *sceneBase) FindUser(userId uint32) (*dbobject.User, bool) {
	u, ok := s.users[userId]
	return u, ok
}

func (s *sceneBase) DelUser(userId uint32) {
	delete(s.users, userId)
	delete(s.players, userId)
}


func (s*scene1) EnterScene(userId uint32) bool{
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.users) < s.maxUser {
		u := dbobject.User{Id: userId}
		dbobject.FindUserById(&u)
		s.users[userId] = &u

		rand.Seed(time.Now().UnixNano())
		x := rand.Intn(1280)
		y := rand.Intn(720)
		s.players[userId] = &player{UserId:userId, Name:u.Name,X:x,Y:y}

		up := userPush{}
		up.Players = s.players
		s.SendMessageToAll(protoUserPush, up)

		return true
	}

	return false
}

func (s *scene1) ExitScene(userId uint32) bool{

	s.UserOffLine(userId)
	return true
}


func (s*scene1) GameMessage(userId uint32, msgName string, data []byte){
	s.lock.Lock()
	defer s.lock.Unlock()

	if msgName == protoSceneReq {
		up := sceneData{}
		up.Players = s.players
		up.Monsters = s.monsters
		s.SendMessageToUser(userId, protoSceneAck, up)
	}else if msgName == protoMoveReq {
		m := move{}
		json.Unmarshal(data, &m)
		s.SendMessageToAll(protoMovePush, m)
		p, ok := s.players[m.UserId]
		if ok {
			p.X = m.TX
			p.Y = m.TY
		}

	}else if msgName == protoAttackReq {
		a := attackReq{}

		json.Unmarshal(data, &a)
		m, ok := s.monsters[a.MonsterId]
		if ok {
			m.Hp -= a.Hurt
			if m.Hp <0 {
				m.Hp = 0
			}

			ap := attackPush{}
			ap.Hurt = a.Hurt
			ap.UserId = a.UserId
			ap.MonsterHp = m.Hp
			ap.MonsterId = a.MonsterId

			s.SendMessageToAll(protoAttackPush, ap)

			if m.Hp == 0{
				delete(s.monsters, a.MonsterId)
			}
		}

	}
}

func (s*scene1) UserOffLine(userId uint32) bool{
	s.lock.Lock()
	defer s.lock.Unlock()
	s.DelUser(userId)

	up := userPush{}
	up.Players = s.players
	s.SendMessageToAll(protoUserPush, up)

	return true
}



/*
发送消息给指定用户
*/
func (s *sceneBase) SendMessageToUser(userId uint32, msgName string, msg interface{}) {
	ok, state := GUserMgr.UserIsIn(userId)
	if ok && state.State == GUserStateOnline {
		data, err := json.Marshal(msg)
		if err == nil{
			state.Conn.SendMsg(msgName, data)
		}
	}
}

/*
发送消息给table的所有用户
*/
func (s *sceneBase) SendMessageToAll(msgName string, msg interface{}) {
	for id,_ := range s.users {
		s.SendMessageToUser(id, msgName, msg)
	}
}