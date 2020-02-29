package slgdb

import (
	"fmt"
	"github.com/llr104/LiFrame/core/orm"
)

type Lumber struct {
	Id      	int      `json:"Id"`
	Name        string   `json:"name" orm:"size(20)"`
	RoleId      uint32   `json:"roleId"`
	Level       int8     `json:"level"`
	Type        int8     `json:"type"`
	Yield       uint32   `json:"yield"`
}

func (s *Lumber) TableName() string {
	return "tb_lumber"
}

/*
新建角色木材类型建筑
*/
func NewRoleAllBLumbers(roleId uint32) [] *Lumber{
	arr := make([] *Lumber, 16)
	for i:=0; i<16; i++ {
		d := Lumber{}
		d.Name = fmt.Sprintf("木材%d", i+1)
		d.Type = 0
		d.Level = 1
		d.RoleId = roleId
		d.Yield = 1000
		arr[i] = &d
	}
	return arr
}

func InsertLumbersToDB(arr []*Lumber) []*Lumber{
	orm.NewOrm().InsertMulti(1, arr)
	return arr
}


func ReadLumbers(roleId uint32) []*Lumber{
	var lumbers []*Lumber
	qry := orm.NewOrm().QueryTable(&Lumber{}).Filter("role_id", roleId)
	qry.All(&lumbers)
	return lumbers
}

func UpdateLumber(b *Lumber) {
	orm.NewOrm().Update(b)
}