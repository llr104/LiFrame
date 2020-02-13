package game

import (
	"fmt"
	"github.com/llr104/LiFrame/utils"
	"math/rand"
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

func (s*scene1) calMonsterPosition()(bool, int, int){
	if len(s.monsters) == 0{
		x := rand.Intn(1000)+100
		y := rand.Intn(360)+200
		return true, x, y
	}


	for i:=0; i<=100; i++ {
		x := rand.Intn(1160)+40
		y := rand.Intn(400)+180

		ok := true
		for _, v := range s.monsters  {
			d := (x-v.X)*(x-v.X) + (y-v.Y)*(y-v.Y)
			if d <4900{
				ok = false
				break
			}
		}

		if ok{
			return true, x, y
		}
	}
	return false, 0, 0
}

func (s*scene1) step(){
	s.lock.Lock()
	defer s.lock.Unlock()

	if len(s.monsters) <= 15{

		m := newRandomMonster()
		ok, x, y := s.calMonsterPosition()
		if ok == false{
			return
		}

		s.curId++
		m.X = x
		m.Y = y
		m.Name = fmt.Sprintf("陪练 %d", s.curId)

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





