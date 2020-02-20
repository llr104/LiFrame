package gameslg

import (
	"github.com/llr104/LiFrame/server/db/slgdb"
	"github.com/llr104/LiFrame/server/gameslg/slgproto"
)

var buildingMgr buildingManager

func init() {
	buildingMgr = buildingManager{
		barrackMap:make(map[uint32] []*slgdb.Barrack),
		dwellingkMap:make(map[uint32] []*slgdb.Dwelling),
		farmlandMap:make(map[uint32] []*slgdb.Farmland),
		lumberMap:make(map[uint32] []*slgdb.Lumber),
		minefieldMap:make(map[uint32] []*slgdb.Mine),
	}
}

type buildingManager struct {
	barrackMap 		map[uint32] []*slgdb.Barrack
	dwellingkMap 	map[uint32] []*slgdb.Dwelling
	farmlandMap 	map[uint32] []*slgdb.Farmland
	lumberMap 		map[uint32] []*slgdb.Lumber
	minefieldMap 	map[uint32] []*slgdb.Mine
}

func (s *buildingManager) getBuilding(roleId uint32, buildingType int8) interface{} {
	if buildingType == slgproto.Building_Barrack {
		b, ok := s.barrackMap[roleId]
		if ok {
			return b
		}else{
			r := slgdb.ReadBarracks(roleId)
			s.barrackMap[roleId] = r
			return r
		}
	}else if buildingType == slgproto.Building_Dwelling{
		b, ok := s.dwellingkMap[roleId]
		if ok {
			return b
		}else{
			r := slgdb.ReadDwellings(roleId)
			s.dwellingkMap[roleId] = r
			return r
		}
	}else if buildingType == slgproto.Building_Farmland{
		b, ok := s.farmlandMap[roleId]
		if ok {
			return b
		}else{
			r := slgdb.ReadFarmlands(roleId)
			s.farmlandMap[roleId] = r
			return r
		}
	}else if buildingType == slgproto.Building_Lumberyard{
		b, ok := s.lumberMap[roleId]
		if ok {
			return b
		}else{
			r := slgdb.ReadLumbers(roleId)
			s.lumberMap[roleId] = r
			return r
		}
	}else if buildingType == slgproto.Building_Minefield{
		b, ok := s.minefieldMap[roleId]
		if ok {
			return b
		}else{
			r := slgdb.ReadMines(roleId)
			s.minefieldMap[roleId] = r
			return r
		}
	}

	return nil
}

func (s *buildingManager) upBuilding(roleId uint32, buildId int, buildingType int8) (interface{}, bool) {
	 r := s.getBuilding(roleId, buildingType)
	 if r == nil{
	 	return nil, false
	 }

	 /*
	 先简单升级，不考虑资源消耗
	 */
	if buildingType == slgproto.Building_Barrack {
		b := r.([]*slgdb.Barrack)
		for _,v := range b{
			if v.Id == buildId && v.Level <= int8(100){
				v.Level++
				v.Yield = uint32(int(v.Level) * 1000)
				slgdb.UpdateBarrack(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.Building_Dwelling {
		b := r.([]*slgdb.Dwelling)
		for _,v := range b{
			if v.Id == buildId && v.Level <= int8(100){
				v.Level++
				v.Yield = uint32(int(v.Level) * 1000)
				slgdb.UpdateDwelling(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.Building_Farmland {
		b := r.([]*slgdb.Farmland)
		for _,v := range b{
			if v.Id == buildId && v.Level <= int8(100){
				v.Level++
				v.Yield = uint32(int(v.Level) * 1000)
				slgdb.UpdateFarmland(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.Building_Lumberyard {
		b := r.([]*slgdb.Lumber)
		for _,v := range b{
			if v.Id == buildId && v.Level <= int8(100){
				v.Level++
				v.Yield = uint32(int(v.Level) * 1000)
				slgdb.UpdateLumber(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.Building_Minefield {
		b := r.([]*slgdb.Mine)
		for _,v := range b{
			if v.Id == buildId && v.Level <= int8(100){
				v.Level++
				v.Yield = uint32(int(v.Level) * 1000)
				slgdb.UpdateMine(v)
				return v, true
			}
		}
	}

	return nil, false
}
