package slgdb

import (
	"github.com/llr104/LiFrame/core/orm"
	"github.com/llr104/LiFrame/server/gameslg/xlsx"
	"github.com/llr104/LiFrame/utils"
	"math/rand"
	"time"
)

/*
将领
*/
type General struct {
	Id      		uint32   `json:"Id"`
	GID				uint32   `json:"gId" orm:"column(gId);description(武将配置Id)"`
	Name        	string   `json:"name" orm:"size(20)"`
	RoleId      	uint32   `json:"roleId"`
	Attack      	int32    `json:"attack"`
	Defense     	int32    `json:"defense"`
	AttackRate		int32    `json:"attack_rate"`
	DefenseRate		int32    `json:"defense_rate"`

	SoldierNum  	int16    `json:"soldierNum" orm:"description(当前士兵数量)"`
	SoldierMax  	int16    `json:"soldierMax" orm:"description(最大士兵数量)"`
	Level       	int8     `json:"level"`
	Exp             int32	 `json:"exp"`
	CityId      	int16    `json:"cityId" orm:"description(目前驻守的城池Id)"`
}

func (s *General) TableName() string {
	return "tb_general"
}

func RandomNPCNewGeneral() *General{
	/*
		随机一个武将
	*/
	rand.Seed(time.Now().UnixNano())
	t := utils.XlsxMgr.Get(xlsx.XlsxGeneral, xlsx.SheetBase)
	i := rand.Intn(t.GetCnt())

	g := General{}
	g.Level = 1
	g.RoleId = 0
	g.Exp = 0

	gId, _ := t.GetInt("gId" , i)
	name, _ := t.GetString("name", i)
	attack , _ := t.GetInt("attack", i)
	defense, _ := t.GetInt("defense", i)
	soldier, _ := t.GetInt("soldier", i)
	attackRate, _ := t.GetInt("attack_rate", i)
	defenseRate, _ := t.GetInt("defense_rate", i)

	g.GID = uint32(gId)
	g.Name = name
	g.Attack = int32(attack)
	g.Defense = int32(defense)
	g.DefenseRate = int32(defenseRate)
	g.AttackRate = int32(attackRate)
	g.SoldierNum = int16(soldier)
	g.SoldierMax = int16(soldier)
	return &g
}

func RandomNewGeneral(roleId uint32) *General{

	g := RandomNPCNewGeneral()
	g.RoleId = roleId
	InsertGeneral(g)
	
	return g
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

func UpdateGeneral(b *General) {
	orm.NewOrm().Update(b)
}