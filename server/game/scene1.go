package game

import (
	"LiFrame/utils"
	"time"
)

type scene1 struct {
	sceneBase

}

func NewScene1() iScene {
	s := scene1{}
	s.init()
	s.maxUser = 1000
	s.start()

	return &s
}

func (s*scene1) start(){
	utils.Scheduler.NewTimerInterval(time.Duration(2)*time.Second, utils.IntervalForever, update, []interface{}{s})
}

func update(args ...interface{}){
	s := args[0].(*scene1)
	s.step()
}

func (s*scene1) step(){
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.monsters) <= 20{
		s.curId++
		m := newRandomMonster()
		m.Id = s.curId
		s.monsters[s.curId] = m

		p := monsterPush{}
		p.Monsters = make( map[uint32] *monster)
		p.Monsters[s.curId] = m
		s.SendMessageToAll(protoMonsterPush, p)

		if s.curId > 1000000000{
			s.curId = 0
		}

		utils.Log.Info("cur monster count:%d", len(s.monsters))
	}
}





