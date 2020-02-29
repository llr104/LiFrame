package slgdb

import (
	"fmt"
	"github.com/llr104/LiFrame/core/orm"
)

type Dwelling struct {
	Id      	int      `json:"Id"`
	Name        string   `json:"name" orm:"size(20)"`
	RoleId      uint32   `json:"roleId"`
	Level       int8     `json:"level"`
	Type        int8     `json:"type"`
	Yield       uint32   `json:"yield"`
}

func (s *Dwelling) TableName() string {
	return "tb_dwelling"
}



/*
新建角色民居类型建筑
*/
func NewRoleAllDwellings(roleId uint32) [] *Dwelling{
	arr := make([] *Dwelling, 16)
	for i:=0; i<16; i++ {
		d := Dwelling{}
		d.Name = fmt.Sprintf("民居%d", i+1)
		d.Type = 0
		d.Level = 1
		d.RoleId = roleId
		d.Yield = 1000
		arr[i] = &d
	}
	return arr
}

func InsertDwellingsToDB(arr []*Dwelling) []*Dwelling{
	orm.NewOrm().InsertMulti(1, arr)
	return arr
}

func ReadDwellings(roleId uint32) []*Dwelling{
	var dwellings []*Dwelling
	qry := orm.NewOrm().QueryTable(&Dwelling{}).Filter("role_id", roleId)
	qry.All(&dwellings)
	return dwellings
}

func UpdateDwelling(b *Dwelling) {
	orm.NewOrm().Update(b)
}