package gameslg

import (
	"github.com/llr104/LiFrame/core/liNet"
	"github.com/llr104/LiFrame/core/orm"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
	"github.com/llr104/LiFrame/utils"
)

/*
世界地图
*/
var WorldMap worldMap

func init() {
	WorldMap = worldMap{cityMap:make(map[int] *slgdb.City)}
}

func (s *worldMap) Load() {
	count, err := orm.NewOrm().QueryTable(slgdb.City{}).Count()
	if err != nil{
		utils.Log.Error("WorldMap db error:%s", err.Error())
	}else{
		if count == 0 {
			/*
				还没有城市数据，创建城市数据
				先简单点一个9个城市1个关卡，魏蜀吴开始每个国家各有三个城池，魏蜀吴中间有一个关卡连接，只有占领了关卡才能攻击其他国家
			*/
			{
				c1 := slgdb.City{Name:"魏都城",Nation:slgdb.NationWei}
				slgdb.InsertCityToDB(&c1)
				s.cityMap[c1.Id] = &c1

				c2 := slgdb.City{Name:"魏副城1",Nation:slgdb.NationWei}
				slgdb.InsertCityToDB(&c2)
				s.cityMap[c2.Id] = &c2

				c3 := slgdb.City{Name:"魏副城2",Nation:slgdb.NationWei}
				slgdb.InsertCityToDB(&c3)
				s.cityMap[c3.Id] = &c3
			}

			{
				c1 := slgdb.City{Name:"蜀都城",Nation:slgdb.NationShu}
				slgdb.InsertCityToDB(&c1)
				s.cityMap[c1.Id] = &c1

				c2 := slgdb.City{Name:"蜀副城1",Nation:slgdb.NationShu}
				slgdb.InsertCityToDB(&c2)
				s.cityMap[c2.Id] = &c2

				c3 := slgdb.City{Name:"蜀副城2",Nation:slgdb.NationShu}
				slgdb.InsertCityToDB(&c3)
				s.cityMap[c3.Id] = &c3
			}

			{
				c1 := slgdb.City{Name:"吴都城",Nation:slgdb.NationWu}
				slgdb.InsertCityToDB(&c1)
				s.cityMap[c1.Id] = &c1

				c2 := slgdb.City{Name:"吴副城1",Nation:slgdb.NationWu}
				slgdb.InsertCityToDB(&c2)
				s.cityMap[c2.Id] = &c2

				c3 := slgdb.City{Name:"吴副城2",Nation:slgdb.NationWu}
				slgdb.InsertCityToDB(&c3)
				s.cityMap[c3.Id] = &c3
			}

			{
				c := slgdb.City{Name:"中立城",Nation:slgdb.NationOther}
				slgdb.InsertCityToDB(&c)
			}

		}else {
			var citys []*slgdb.City
			orm.NewOrm().QueryTable(slgdb.City{}).All(&citys)
			for _,v := range citys{
				s.cityMap[v.Id] = v
			}
		}

		utils.Log.Info("读取城市数据成功:%v", s.cityMap)
	}
}

type worldMap struct {
	liNet.BaseRouter
	cityMap 		map[int] *slgdb.City
}

func (s *worldMap) NameSpace() string {
	return "worldMap"
}
