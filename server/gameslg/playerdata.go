package gameslg

import (
	"github.com/llr104/LiFrame/server/gameslg/data"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
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
	generalMap 		map[uint32] *slgdb.General

	barrackYield		uint32
	dwellingkYield 		uint32
	farmlandYield		uint32
	lumberYield 		uint32
	minefieldYield 		uint32

	barrackCapacity		uint32
	dwellingkCapacity	uint32
	farmlandCapacity	uint32
	lumberCapacity	 	uint32
	minefieldCapacity	uint32


}

func (s *playerData) init() {

	s.generalMap = make(map[uint32] *slgdb.General)
	b := s.role.OffLineTime
	if b != 0{
		e := time.Now().Unix()
		diff := float64(e-b)/60.0
		//uint32(math.Ceil(float64(s.getYield(slgproto.Building_Barrack) / 60.0)))

		s.role.Mine += uint32(math.Ceil(float64(s.getYield(slgproto.BuildingMinefield) / 60.0)*diff))
		s.role.Food += uint32(math.Ceil(float64(s.getYield(slgproto.BuildingLumberyard) / 60.0)*diff))
		s.role.Wood += uint32(math.Ceil(float64(s.getYield(slgproto.BuildingFarmland) / 60.0)*diff))
		s.role.Silver += uint32(math.Ceil(float64(s.getYield(slgproto.BuildingDwelling) / 60.0)*diff))

		s.checkCapacity()
	}

}

func (s *playerData) getBuilding(buildingType int8) interface{} {
	if buildingType == slgproto.BuildingBarrack {
		if s.barracks == nil{
			r := slgdb.ReadBarracks(s.role.RoleId)
			s.barracks = r
		}
		return s.barracks
	}else if buildingType == slgproto.BuildingDwelling {
		if s.dwellingks == nil{
			r := slgdb.ReadDwellings(s.role.RoleId)
			s.dwellingks = r
		}
		return s.dwellingks
	}else if buildingType == slgproto.BuildingFarmland {
		if s.farmlands == nil{
			r := slgdb.ReadFarmlands(s.role.RoleId)
			s.farmlands = r
		}
		return s.farmlands

	}else if buildingType == slgproto.BuildingLumberyard {
		if s.lumbers == nil{
			r := slgdb.ReadLumbers(s.role.RoleId)
			s.lumbers = r
		}
		return s.lumbers

	}else if buildingType == slgproto.BuildingMinefield {
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
	if buildingType == slgproto.BuildingBarrack {
		b := r.([]*slgdb.Barrack)

		for _,v := range b{
			if v.Id == buildId && v.Level < data.BarrackMaxLevel(){
				v.Level++
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = data.BarrackYield(v.Level)
				s.barrackYield = old + v.Yield

				slgdb.UpdateBarrack(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.BuildingDwelling {
		b := r.([]*slgdb.Dwelling)
		for _,v := range b{
			if v.Id == buildId && v.Level < data.DwellingMaxLevel(){
				v.Level++
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = data.DwellingYield(v.Level)
				s.dwellingkYield = old + v.Yield

				slgdb.UpdateDwelling(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.BuildingFarmland {
		b := r.([]*slgdb.Farmland)
		for _,v := range b{
			if v.Id == buildId && v.Level < data.FarmlandMaxLevel(){
				v.Level++
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = data.FarmlandYield(v.Level)
				s.farmlandYield = old + v.Yield

				slgdb.UpdateFarmland(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.BuildingLumberyard {
		b := r.([]*slgdb.Lumber)
		for _,v := range b{
			if v.Id == buildId && v.Level < data.LumberMaxLevel(){
				v.Level++
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = data.LumberYield(v.Level)
				s.lumberYield = old + v.Yield

				slgdb.UpdateLumber(v)
				return v, true
			}
		}
	}else if buildingType == slgproto.BuildingMinefield {
		b := r.([]*slgdb.Mine)
		for _,v := range b{
			if v.Id == buildId && v.Level < data.MineMaxLevel(){
				v.Level++
				old := s.getYield(buildingType)
				old -= v.Yield
				v.Yield = data.MineYield(v.Level)
				s.minefieldYield = old + v.Yield

				slgdb.UpdateMine(v)
				return v, true
			}
		}
	}

	return nil, false
}

func (s* playerData) getCapacity(buildingType int8) uint32{
	r := s.getBuilding(buildingType)
	if r == nil{
		return  0
	}else {
		if buildingType == slgproto.BuildingBarrack {
			if s.barrackCapacity == 0{
				b := r.([]*slgdb.Barrack)
				for _,v := range b{
					s.barrackCapacity += data.BarrackCapacity(v.Level)
				}
			}
			return s.barrackCapacity

		}else if buildingType == slgproto.BuildingDwelling {
			if s.dwellingkCapacity == 0{
				b := r.([]*slgdb.Dwelling)
				for _,v := range b{
					s.dwellingkCapacity += data.DwellingCapacity(v.Level)
				}
			}
			return s.dwellingkCapacity

		}else if buildingType == slgproto.BuildingFarmland {
			if s.farmlandCapacity == 0{
				b := r.([]*slgdb.Farmland)
				for _,v := range b{
					s.farmlandCapacity += data.FarmlandCapacity(v.Level)
				}
			}
			return s.farmlandCapacity

		}else if buildingType == slgproto.BuildingLumberyard {
			if s.lumberCapacity == 0{
				b := r.([]*slgdb.Lumber)
				for _,v := range b{
					s.lumberCapacity += data.LumberCapacity(v.Level)
				}
			}
			return s.lumberYield

		}else if buildingType == slgproto.BuildingMinefield {
			if s.minefieldCapacity == 0{
				b := r.([]*slgdb.Mine)
				for _,v := range b{
					s.minefieldCapacity +=  data.MineCapacity(v.Level)
				}
			}
			return s.minefieldCapacity
		}
		return 0
	}
}

func (s *playerData) getYield(buildingType int8) uint32 {
	r := s.getBuilding(buildingType)
	if r == nil{
		return  0
	}else {
		if buildingType == slgproto.BuildingBarrack {
			if s.barrackYield == 0{
				b := r.([]*slgdb.Barrack)
				for _,v := range b{
					s.barrackYield += v.Yield
				}
			}
			return s.barrackYield

		}else if buildingType == slgproto.BuildingDwelling {
			if s.dwellingkYield == 0{
				b := r.([]*slgdb.Dwelling)
				for _,v := range b{
					s.dwellingkYield += v.Yield
				}
			}
			return s.dwellingkYield

		}else if buildingType == slgproto.BuildingFarmland {
			if s.farmlandYield == 0{
				b := r.([]*slgdb.Farmland)
				for _,v := range b{
					s.farmlandYield += v.Yield
				}
			}
			return s.farmlandYield

		}else if buildingType == slgproto.BuildingLumberyard {
			if s.lumberYield == 0{
				b := r.([]*slgdb.Lumber)
				for _,v := range b{
					s.lumberYield += v.Yield
				}
			}
			return s.lumberYield

		}else if buildingType == slgproto.BuildingMinefield {
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

func (s* playerData) getGenerals() [] *slgdb.General {
	if len(s.generalMap) == 0{
		generals := slgdb.ReadGenerals(s.role.RoleId)
		if len(generals) == 0{
			//如果没有武将，默认给几个吧
			for i:=0; i<3; i++ {
				g := slgdb.RandomNewGeneral(s.role.RoleId)
				s.generalMap[g.Id] = g
			}
		}else{
			for _, g := range generals{
				s.generalMap[g.Id] = g
			}
		}
	}

	//转成数组给客户端吧
	n := len(s.generalMap)
	arr := make([]*slgdb.General, n)
	i := 0
	for _, g := range s.generalMap{
		arr[i] = g
		i++
	}

	return arr
}

func  (s* playerData) saveToDB(){
	/*
	暂时傻瓜式写库，后面会优化，减少不必要的写入
	*/
	slgdb.UpdateRoleOffline(s.role)

	for _, v := range s.minefields  {
		slgdb.UpdateMine(v)
	}

	for _, v := range s.lumbers  {
		slgdb.UpdateLumber(v)
	}

	for _, v := range s.farmlands  {
		slgdb.UpdateFarmland(v)
	}

	for _, v := range s.dwellingks  {
		slgdb.UpdateDwelling(v)
	}

	for _, v := range s.barracks  {
		slgdb.UpdateBarrack(v)
	}

	for _, v := range s.minefields  {
		slgdb.UpdateMine(v)
	}

	for _, v := range s.generalMap{
		slgdb.UpdateGeneral(v)
	}
}

func (s* playerData) stepYield() {
	//uint32(math.Ceil(float64(s.getYield(slgproto.Building_Barrack) / 60.0)))

	s.role.Mine += uint32(math.Ceil(float64(s.getYield(slgproto.BuildingMinefield) / 60.0)))
	s.role.Food += uint32(math.Ceil(float64(s.getYield(slgproto.BuildingLumberyard) / 60.0)))
	s.role.Wood += uint32(math.Ceil(float64(s.getYield(slgproto.BuildingFarmland) / 60.0)))
	s.role.Silver += uint32(math.Ceil(float64(s.getYield(slgproto.BuildingDwelling) / 60.0)))
	s.checkCapacity()
}

func (s* playerData) checkCapacity(){
	/*
		需要考虑仓库容量，不能爆仓
	*/
	maxM := s.getCapacity(slgproto.BuildingMinefield)
	maxL := s.getCapacity(slgproto.BuildingLumberyard)
	maxF := s.getCapacity(slgproto.BuildingFarmland)
	maxD := s.getCapacity(slgproto.BuildingDwelling)

	s.role.Mine = uint32(math.Min(float64(maxM), float64(s.role.Mine)))
	s.role.Food = uint32(math.Min(float64(maxL), float64(s.role.Food)))
	s.role.Wood = uint32(math.Min(float64(maxF), float64(s.role.Wood)))
	s.role.Silver = uint32(math.Min(float64(maxD), float64(s.role.Silver)))
}
