package gameslg

import (
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
	"github.com/llr104/LiFrame/utils"
	"sync"
	"time"
)

var playerMgr playerManager

func init() {
	playerMgr = playerManager{
		playerMaps:make(map[uint32] *playerData),
	}

	utils.Scheduler.NewTimerInterval(60*time.Second, utils.IntervalForever, step, []interface{}{})
}

func step(v ...interface{}) {
	playerMgr.step()
}

type playerManager struct {
	playerMaps 		map[uint32] *playerData
	mutex 			sync.RWMutex
}

func (s *playerManager) step(){
	s.mutex.Lock()
	defer s.mutex.Unlock()

	//给玩家算上单位时间内生产的产量，后续需要考虑仓库容量
	for _,v := range s.playerMaps{
		v.stepYield()
	}
}


func (s *playerManager) createPlayer(role*slgdb.Role) *playerData{
	s.mutex.Lock()
	defer s.mutex.Unlock()

	r,ok :=s.playerMaps[role.RoleId]
	if ok == false{
		r = newPlayerData(role)
	}
	s.playerMaps[role.RoleId] = r
	return r
}


func (s *playerManager) addPlayer(role*slgdb.Role, barracks[]*slgdb.Barrack, dwellingks []*slgdb.Dwelling,
	farmlands []*slgdb.Farmland, lumbers []*slgdb.Lumber, minefields []*slgdb.Mine) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	p := playerData{role:role,barracks:barracks,dwellingks:dwellingks,farmlands:farmlands,lumbers:lumbers,minefields:minefields}
	p.init()

	s.playerMaps[role.RoleId] = &p
}

func (s *playerManager) ReleasePlayer(userId uint32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	r := slgdb.Role{}
	r.UserId = userId
	slgdb.FindRoleByUserId(&r)

	player, ok := s.playerMaps[r.RoleId]
	if ok{
		player.saveToDB()
		delete(s.playerMaps, r.RoleId)
	}
}

func (s *playerManager) getRole(roleId uint32) *slgdb.Role {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	p, ok := s.playerMaps[roleId]
	if ok {
		return p.role
	}else{
		return nil
	}
}

func (s *playerManager) getBuilding(roleId uint32, buildingType int8) interface{}{
	s.mutex.Lock()
	defer s.mutex.Unlock()

	p, ok := s.playerMaps[roleId]
	if ok {
		return p.getBuilding(buildingType)
	}else {
		return nil
	}
}

func (s *playerManager) upBuilding(roleId uint32, buildId int, buildingType int8) (interface{}, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	p, ok := s.playerMaps[roleId]
	if ok {
		return p.upBuilding(buildId, buildingType)
	}else {
		return nil, false
	}
}

func (s *playerManager) getYield(roleId uint32,  buildingType int8) uint32 {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	p, ok := s.playerMaps[roleId]
	if ok {
		return p.getYield(buildingType)
	}else {
		return 0
	}
}

func (s *playerManager) getGenerals(roleId uint32) ([] *slgdb.General, bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	p, ok := s.playerMaps[roleId]
	if ok {
		return p.getGenerals(), true
	}else {
		return nil, false
	}
}


func (s *playerManager) getGeneral(roleId uint32, generalId uint32) (*slgdb.General, bool) {
	arr, ok := s.getGenerals(roleId)
	if ok {
		for _, v := range arr {
			if v.Id == generalId{
				return v, true
			}
		}
	}

	return nil, false
}


