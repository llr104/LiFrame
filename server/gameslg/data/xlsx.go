package data

import "github.com/llr104/LiFrame/utils"

const XlsxBuilding = "conf/xlsx/building.xlsx"
const SheetDwelling  = "dwelling"
const SheetBarrack  = "barrack"
const SheetFarmland  = "farmland"
const SheetLumber  = "lumber"
const SheetMine  = "mine"

func init() {
	utils.XlsxMgr.Load(XlsxBuilding)
}

func DwellingMaxLevel() int8 {
	 t := utils.XlsxMgr.Get(XlsxBuilding, SheetDwelling)
	 if t != nil {
		c := t.GetCnt()
		return int8(c)
	 }else{
		 return 1
	 }
}

func DwellingYield(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetDwelling)
	if t != nil {
		if int(level) < t.GetCnt(){
			yield, _ := t.GetInt("yield", int(level))
			return uint32(yield)
		}
	}else{
		return 0
	}
	return 0
}


func BarrackMaxLevel() int8 {
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetBarrack)
	if t != nil {
		c := t.GetCnt()
		return int8(c)
	}else{
		return 1
	}
}

func BarrackYield(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetBarrack)
	if t != nil {
		if int(level) < t.GetCnt(){
			yield, _ := t.GetInt("yield", int(level))
			return uint32(yield)
		}
	}else{
		return 0
	}
	return 0
}


func FarmlandMaxLevel() int8 {
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetFarmland)
	if t != nil {
		c := t.GetCnt()
		return int8(c)
	}else{
		return 1
	}
}

func FarmlandYield(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetFarmland)
	if t != nil {
		if int(level) < t.GetCnt(){
			yield, _ := t.GetInt("yield", int(level))
			return uint32(yield)
		}
	}else{
		return 0
	}
	return 0
}



func LumberMaxLevel() int8 {
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetLumber)
	if t != nil {
		c := t.GetCnt()
		return int8(c)
	}else{
		return 1
	}
}

func LumberYield(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetLumber)
	if t != nil {
		if int(level) < t.GetCnt(){
			yield, _ := t.GetInt("yield", int(level))
			return uint32(yield)
		}
	}else{
		return 0
	}
	return 0
}


func MineMaxLevel() int8 {
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetMine)
	if t != nil {
		c := t.GetCnt()
		return int8(c)
	}else{
		return 1
	}
}

func MineYield(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetMine)
	if t != nil {
		if int(level) < t.GetCnt(){
			yield, _ := t.GetInt("yield", int(level))
			return uint32(yield)
		}
	}else{
		return 0
	}
	return 0
}