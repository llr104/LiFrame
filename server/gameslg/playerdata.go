package gameslg

import (
	"github.com/llr104/LiFrame/server/db/slgdb"
	"github.com/llr104/LiFrame/server/gameslg/slgproto"
	"math"
	"time"
)

func newPlayerData(role *slgdb.Role) *playerData {
	p := playerData{role:role,}
	p.init()

	return &p
}

type playerData struct {
	role *slgdb.Role

	barracks		[]*slgdb.Barrack
	dwellingks 		[]*slgdb.Dwelling
	farmlands		[]*slgdb.Farmland
	lumbers 		[]*slgdb.Lumber
	minefields 		[]*slgdb.Mine

	barrackYield		uint32
	dwellingkYield 		uint32
	farmlandYield		uint32
	lumberYield 		uint32
	minefieldYield 		uint32

}

func (s *playerData) init() {
	b := s.role.OffLineTime
	if b != 0{
		e := time.Now().Unix()
		diff := float64(e-b)/60.0
		//uint32(math.Ceil(float64(s.getYield(slgproto.Building_Barrack) / 60.0)))

		s.role.Mine += uint32(math.Ceil(float64(s.getYield(slgproto.Building_Minefield) / 60.0)*diff))
		s.role.Food += uint32(math.Ceil(float64(s.getYield(slgproto.Building_Lumberyard) / 60.0)*diff))
		s.role.Wood += uint32(math.Ceil(float64(s.getYield(slgproto.Building_Farmland) / 60.0)*diff))
		s.role.Silver += uint32(math.Ceil(float64(s.getYield(slgproto.Building_Dwelling) / 60.0)*diff))
	}

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
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = uint32(int(v.Level) * 1000)
				s.barrackYield = old + v.Yield

				slgdb.UpdateBarrack(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.Building_Dwelling {
		b := r.([]*slgdb.Dwelling)
		for _,v := range b{
			if v.Id == buildId && v.Level <= int8(100){
				v.Level++
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = uint32(int(v.Level) * 1000)
				s.dwellingkYield = old + v.Yield

				slgdb.UpdateDwelling(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.Building_Farmland {
		b := r.([]*slgdb.Farmland)
		for _,v := range b{
			if v.Id == buildId && v.Level <= int8(100){
				v.Level++
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = uint32(int(v.Level) * 1000)
				s.farmlandYield = old + v.Yield

				slgdb.UpdateFarmland(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.Building_Lumberyard {
		b := r.([]*slgdb.Lumber)
		for _,v := range b{
			if v.Id == buildId && v.Level <= int8(100){
				v.Level++
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = uint32(int(v.Level) * 1000)
				s.lumberYield = old + v.Yield

				slgdb.UpdateLumber(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.Building_Minefield {
		b := r.([]*slgdb.Mine)
		for _,v := range b{
			if v.Id == buildId && v.Level <= int8(100){
				v.Level++
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = uint32(int(v.Level) * 1000)
				s.minefieldYield = old + v.Yield

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
		if buildingType == slgproto.Building_Barrack {
			if s.barrackYield == 0{
				b := r.([]*slgdb.Barrack)
				for _,v := range b{
					s.barrackYield += v.Yield
				}
			}
			return s.barrackYield

		}else if buildingType == slgproto.Building_Dwelling {
			if s.dwellingkYield == 0{
				b := r.([]*slgdb.Dwelling)
				for _,v := range b{
					s.dwellingkYield += v.Yield
				}
			}
			return s.dwellingkYield

		}else if buildingType == slgproto.Building_Farmland {
			if s.farmlandYield == 0{
				b := r.([]*slgdb.Farmland)
				for _,v := range b{
					s.farmlandYield += v.Yield
				}
			}
			return s.farmlandYield

		}else if buildingType == slgproto.Building_Lumberyard {
			if s.lumberYield == 0{
				b := r.([]*slgdb.Lumber)
				for _,v := range b{
					s.lumberYield += v.Yield
				}
			}
			return s.lumberYield

		}else if buildingType == slgproto.Building_Minefield {
			if s.minefieldYield == 0{
				b := r.([]*slgdb.Mine)
				for _,v := range b{
					s.minefieldYield += v.Yield
				}
			}
			return s.minefieldYield
		}
		return 0
	}
}

func (s* playerData) stepYield() {
	//uint32(math.Ceil(float64(s.getYield(slgproto.Building_Barrack) / 60.0)))
	s.role.Mine += uint32(math.Ceil(float64(s.getYield(slgproto.Building_Minefield) / 60.0)))
	s.role.Food += uint32(math.Ceil(float64(s.getYield(slgproto.Building_Lumberyard) / 60.0)))
	s.role.Wood += uint32(math.Ceil(float64(s.getYield(slgproto.Building_Farmland) / 60.0)))
	s.role.Silver += uint32(math.Ceil(float64(s.getYield(slgproto.Building_Dwelling) / 60.0)))

}