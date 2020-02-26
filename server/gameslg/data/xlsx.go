package data

import "github.com/llr104/LiFrame/utils"

const XlsxBuilding = "conf/xlsx/building.xlsx"
const SheetDwelling  = "dwelling"

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
