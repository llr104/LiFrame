package gameslg

import (
	"github.com/llr104/LiFrame/server/db/slgdb"
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





func (s *playerManager) createPlayer(role* slgdb.Role) *playerData{
	s.mutex.Lock()
	defer s.mutex.Unlock()

	r,ok :=s.playerMaps[role.RoleId]
	if ok == false{
		r = newPlayerData(role)
	}
	s.playerMaps[role.RoleId] = r
	return r
}

func (s *playerManager) addPlayer(role* slgdb.Role, barracks[]*slgdb.Barrack, dwellingks []*slgdb.Dwelling,
	farmlands []*slgdb.Farmland, lumbers []*slgdb.Lumber, minefields []*slgdb.Mine) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	p := playerData{role:role,barracks:barracks,dwellingks:dwellingks,farmlands:farmlands,lumbers:lumbers,minefields:minefields}
	s.playerMaps[role.RoleId] = &p
}

func (s *playerManager) ReleasePlayer(userId uint32) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	r := slgdb.Role{}
	r.UserId = userId
	slgdb.FindRoleByUserId(&r)
	slgdb.UpdateRoleOffline(&r)

	delete(s.playerMaps, r.RoleId)
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