package gameslg

import (
	"github.com/llr104/LiFrame/server/db/slgdb"
	"github.com/llr104/LiFrame/server/gameslg/slgproto"
)

var buildingMgr buildingManager

func init() {
	buildingMgr = buildingManager{}
}

type buildingManager struct {
	barrackMap 		map[uint32] *slgdb.Barrack
	dwellingkMap 	map[uint32] *slgdb.Dwelling
	farmlandMap 	map[uint32] *slgdb.Farmland
	lumberMap 		map[uint32] *slgdb.Lumber
	minefieldMap 	map[uint32] *slgdb.Mine
}

func (s *buildingManager) getBuilding(roleId uint32, buildingType int8) interface{} {
	if buildingType == slgproto.Building_Barrack {
		b, ok := s.barrackMap[roleId]
		if ok {
			return b
		}else{
			return slgdb.ReadBarracks(roleId)
		}
	}else if buildingType == slgproto.Building_Dwelling{
		b, ok := s.dwellingkMap[roleId]
		if ok {
			return b
		}else{
			return slgdb.ReadDwellings(roleId)
		}
	}else if buildingType == slgproto.Building_Farmland{
		b, ok := s.farmlandMap[roleId]
		if ok {
			return b
		}else{
			return slgdb.ReadFarmlands(roleId)
		}
	}else if buildingType == slgproto.Building_Lumberyard{
		b, ok := s.lumberMap[roleId]
		if ok {
			return b
		}else{
			return slgdb.ReadLumbers(roleId)
		}
	}else if buildingType == slgproto.Building_Minefield{
		b, ok := s.minefieldMap[roleId]
		if ok {
			return b
		}else{
			return slgdb.ReadMines(roleId)
		}
	}

	return nil
}