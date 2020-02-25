package data

import (
	"fmt"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
	"math/rand"
	"time"
)

type NpcScene struct {
	Id       uint16
	Name     string           `json:"name"`
	Generals []*slgdb.General `json:"generals"`
}

func RandomNPCScene(Id uint16) *NpcScene {
	s := NpcScene{}
	s.Id = Id
	s.Name = fmt.Sprintf("npc场景 %d", s.Id)

	arr1 :=[...]string {"吕布","高顺","貂蝉","张角","张宝","文丑","颜良","张郃","张梁","张绣"}
	arr2 :=[...]string {"袁术","孙坚","廖化","周仓","吕蒙","陆逊","丁原","祖茂","蔡瑁","张允"}
	arr3 :=[...]string {"袁绍","黄盖","陈到","赵云","刘封","姜维","郝昭","曹彰","曹洪","夏侯霸"}
	s.Generals = make([]*slgdb.General, 3)

	rand.Seed(time.Now().UnixNano())
	for i:=0; i<3; i++ {
		g := slgdb.General{}
		g.Level = 1
		g.RoleId = 0
		g.SoldierMax = 1000
		g.SoldierNum = 1000

		/*
			属性先随机吧
		*/
		if i == 0{
			g.Attack = int32(rand.Intn(25) + 75)
			g.Defense = int32(rand.Intn(25) + 75)
			g.Name = arr1[rand.Intn(len(arr1))]
		}else if i == 1 {
			g.Attack = int32(rand.Intn(25) + 75)
			g.Defense = int32(rand.Intn(25) + 75)
			g.Name = arr2[rand.Intn(len(arr2))]
		}else{
			g.Attack = int32(rand.Intn(25) + 75)
			g.Defense = int32(rand.Intn(25) + 75)
			g.Name = arr3[rand.Intn(len(arr3))]
		}
		s.Generals[i] = &g
	}

	return &s
}
