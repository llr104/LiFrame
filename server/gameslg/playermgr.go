package gameslg

import "github.com/llr104/LiFrame/server/db/slgdb"

var playerMgr playerManager

func init() {
	playerMgr = playerManager{
		playerMaps:make(map[uint32] *playerData),
	}
}

type playerManager struct {
	playerMaps 		map[uint32] *playerData
}

func (s *playerManager) newPlayer(role* slgdb.Role) *playerData{
	r,ok :=s.playerMaps[role.RoleId]
	if ok{
		return r
	}else{
		p := newPlayerData(role)
		s.playerMaps[role.RoleId] = p
		return p
	}
}

func (s *playerManager) ReleasePlayer(userId uint32) {
	r := slgdb.Role{}
	r.UserId = userId
	slgdb.FindRoleByUserId(&r)
	delete(s.playerMaps, r.RoleId)
}

func (s *playerManager) getBuilding(roleId uint32, buildingType int8) interface{}{
	p, ok := s.playerMaps[roleId]
	if ok {
		return p.getBuilding(buildingType)
	}else {
		return nil
	}
}

func (s *playerManager) upBuilding(roleId uint32, buildId int, buildingType int8) (interface{}, bool) {
	p, ok := s.playerMaps[roleId]
	if ok {
		return p.upBuilding(buildId, buildingType)
	}else {
		return nil, false
	}
}

func (s *playerManager) getYield(roleId uint32,  buildingType int8) uint32 {
	p, ok := s.playerMaps[roleId]
	if ok {
		return p.getYield(buildingType)
	}else {
		return 0
	}
}