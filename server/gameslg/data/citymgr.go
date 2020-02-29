package data

import (
	"github.com/llr104/LiFrame/core/orm"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
	"github.com/llr104/LiFrame/server/gameslg/xlsx"
	"github.com/llr104/LiFrame/utils"
	"sync"
)

type cityManager struct {
	cityMap 		map[int16] *slgdb.City
	mutex 			sync.RWMutex
}

var CityMgr cityManager

func (s* cityManager) load() {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	count, err := orm.NewOrm().QueryTable(slgdb.City{}).Count()
	if err != nil{
		utils.Log.Error("WorldMap db error:%s", err.Error())
	}else{
		if count == 0 {
			/*
				还没有城市数据，创建城市数据
				先简单点一个10个城市，魏蜀吴开始每个国家各有三个城池，魏蜀吴中间有一个空白城连接，只有占领了空白城才能攻击其他国家
			*/
			t := utils.XlsxMgr.Get(xlsx.XlsxCity, xlsx.SheetWorldCity)
			n := t.GetCnt()
			for i:=0; i<n; i++ {
				cId, _ := t.GetInt("cId", i)
				name, _ := t.GetString("name", i)
				nation, _ := t.GetInt("nation", i)
				capital, _ := t.GetInt("capital", i)
				adjacent, _ := t.GetString("adjacent", i)
				c := slgdb.City{CId:int16(cId), Name:name, Nation:int8(nation), Capital:capital==1, Adjacent:adjacent}
				slgdb.InsertCityToDB(&c)
				s.cityMap[c.CId] = &c
			}

		}else {
			var citys []*slgdb.City
			orm.NewOrm().QueryTable(slgdb.City{}).All(&citys)
			for _,v := range citys{
				s.cityMap[v.CId] = v
			}
		}

		utils.Log.Info("读取城市数据成功:%v", s.cityMap)
	}
}

func (s* cityManager) Count() int{
	s.mutex.Lock()
	defer s.mutex.Unlock()

	return len(s.cityMap)
}

func (s* cityManager) CityMap() map[int16] *slgdb.City {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.cityMap
}
