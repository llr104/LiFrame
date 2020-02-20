package gameslg

import (
	"github.com/llr104/LiFrame/server/db/slgdb"
	"github.com/llr104/LiFrame/server/gameslg/slgproto"
)

func newPlayerData(role *slgdb.Role) *playerData {
	p := playerData{role:role,}
	return &p
}

type playerData struct {
	role *slgdb.Role

	barracks		[]*slgdb.Barrack
	dwellingks 		[]*slgdb.Dwelling
	farmlands		[]*slgdb.Farmland
	lumbers 		[]*slgdb.Lumber
	minefields 		[]*slgdb.Mine
}


func (s *playerData) getBuilding(buildingType int8) interface{} {
	if buildingType == slgproto.Building_Barrack {
		if s.barracks == nil{
			r := slgdb.ReadBarracks(s.role.RoleId)
			s.barracks = r
		}
		return s.barracks
	}else if buildingType == slgproto.Building_Dwelling{
		if s.dwellingks == nil{
			r := slgdb.ReadDwellings(s.role.RoleId)
			s.dwellingks = r
		}
		return s.dwellingks
	}else if buildingType == slgproto.Building_Farmland{
		if s.farmlands == nil{
			r := slgdb.ReadFarmlands(s.role.RoleId)
			s.farmlands = r
		}
		return s.farmlands

	}else if buildingType == slgproto.Building_Lumberyard{
		if s.lumbers == nil{
			r := slgdb.ReadLumbers(s.role.RoleId)
			s.lumbers = r
		}
		return s.lumbers

	}else if buildingType == slgproto.Building_Minefield{
		if s.minefields == nil{
			r := slgdb.ReadMines(s.role.RoleId)
			s.minefields = r
		}
		return s.minefields
	}

	return nil
}

func (s *playerData) upBuilding(buildId int, buildingType int8) (interface{}, bool) {
	r := s.getBuilding(buildingType)
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

func (s *playerData) getYield(buildingType int8) uint32 {
	r := s.getBuilding(buildingType)
	if r == nil{
		return  0
	}else {
		var t uint32 = 0
		if buildingType == slgproto.Building_Barrack {
			b := r.([]*slgdb.Barrack)
			for _,v := range b{
				t += v.Yield
			}
		}else if buildingType == slgproto.Building_Dwelling {
			b := r.([]*slgdb.Dwelling)
			for _,v := range b{
				t += v.Yield
			}
		}else if buildingType == slgproto.Building_Farmland {
			b := r.([]*slgdb.Farmland)
			for _,v := range b{
				t += v.Yield
			}
		}else if buildingType == slgproto.Building_Lumberyard {
			b := r.([]*slgdb.Lumber)
			for _,v := range b{
				t += v.Yield
			}
		}else if buildingType == slgproto.Building_Minefield {
			b := r.([]*slgdb.Mine)
			for _,v := range b{
				t += v.Yield
			}
		}
		return t
	}
}