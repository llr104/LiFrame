package xlsx

import (
	"github.com/llr104/LiFrame/utils"
)


const XlsxBuilding = "building.xlsx"
const SheetDwelling  = "dwelling"
const SheetBarrack  = "barrack"
const SheetFarmland  = "farmland"
const SheetLumber  = "lumber"
const SheetMine  = "mine"

const XlsxGeneral = "general.xlsx"
const SheetBase = "base"
const XlsxCity = "city.xlsx"
const SheetWorldCity  = "worldcity"

func Init(xlsxDir string) {
	utils.XlsxMgr.SetRootDir(xlsxDir)
	utils.XlsxMgr.Load(XlsxBuilding)
	utils.XlsxMgr.Load(XlsxGeneral)
	utils.XlsxMgr.Load(XlsxCity)
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


func DwellingMaxLevel() int8 {
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetDwelling)
	if t != nil {
		c := t.GetCnt()
		return int8(c)
	}else{
		return 1
	}
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


func MineMaxLevel() int8 {
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetMine)
	if t != nil {
		c := t.GetCnt()
		return int8(c)
	}else{
		return 1
	}
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

func FarmlandMaxLevel() int8 {
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetFarmland)
	if t != nil {
		c := t.GetCnt()
		return int8(c)
	}else{
		return 1
	}
}


func BarrackCapacity(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetBarrack)
	if t != nil {
		if int(level) < t.GetCnt(){
			capacity, _ := t.GetInt("capacity", int(level))
			return uint32(capacity)
		}
	}else{
		return 0
	}
	return 0
}

func DwellingCapacity(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetDwelling)
	if t != nil {
		if int(level) < t.GetCnt(){
			capacity, _ := t.GetInt("capacity", int(level))
			return uint32(capacity)
		}
	}else{
		return 0
	}
	return 0
}

func FarmlandCapacity(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetFarmland)
	if t != nil {
		if int(level) < t.GetCnt(){
			capacity, _ := t.GetInt("capacity", int(level))
			return uint32(capacity)
		}
	}else{
		return 0
	}
	return 0
}



func LumberCapacity(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetLumber)
	if t != nil {
		if int(level) < t.GetCnt(){
			capacity, _ := t.GetInt("capacity", int(level))
			return uint32(capacity)
		}
	}else{
		return 0
	}
	return 0
}

func MineCapacity(level int8) uint32{
	t := utils.XlsxMgr.Get(XlsxBuilding, SheetMine)
	if t != nil {
		if int(level) < t.GetCnt(){
			capacity, _ := t.GetInt("capacity", int(level))
			return uint32(capacity)
		}
	}else{
		return 0
	}
	return 0
}


