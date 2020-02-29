package data

import "github.com/llr104/LiFrame/server/gameslg/slgdb"

func Init() {
	CityMgr = cityManager{
		cityMap:make(map[int16] *slgdb.City),
	}
	CityMgr.load()
}
