package slgdb

import (
	"github.com/llr104/LiFrame/core/orm"
	"math/rand"
	"time"
)

/*
将领
*/
type General struct {
	Id      	uint32    `json:"Id"`
	Name        string   `json:"name"`
	RoleId      uint32   `json:"roleId"`
	Level       int8     `json:"level"`
	Attack      int32    `json:"attack"`
	Defense     int32    `json:"defense"`
	SoldierNum  int16    `json:"soldierNum" orm:"description(当前士兵数量)"`
	SoldierMax  int16    `json:"soldierMax" orm:"description(最大士兵数量)"`
}

func RandomNewGeneral(roleId uint32) *General{
	g := General{}
	g.Level = 1
	g.RoleId = roleId
	g.SoldierMax = 1000
	g.SoldierNum = 1000

	/*
	属性先随机吧
	*/
	rand.Seed(time.Now().UnixNano())
	g.Attack = int32(rand.Intn(15) + 85)
	g.Defense = int32(rand.Intn(15) + 85)
	arr :=[...]string {"关羽","张飞","马超","赵云","黄忠","张辽","典韦","徐晃","许褚","夏侯渊","太史慈","甘宁","周泰","周瑜","孙策"}
	g.Name = arr[rand.Intn(len(arr))]

	InsertGeneral(&g)
	
	return &g
}

func ReadGenerals(roleId uint32) []*General{
	var generals []*General
	qry := orm.NewOrm().QueryTable(&General{}).Filter("role_id", roleId)
	qry.All(&generals)
	return generals
}

func InsertGeneral(g *General)  {
	orm.NewOrm().Insert(g)
}