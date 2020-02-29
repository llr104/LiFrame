package slgdb

import (
	"fmt"
	"github.com/llr104/LiFrame/core/orm"
)

type Farmland struct {
	Id      	int      `json:"Id"`
	Name        string   `json:"name" orm:"size(20)"`
	RoleId      uint32   `json:"roleId"`
	Level       int8     `json:"level"`
	Type        int8     `json:"type"`
	Yield       uint32   `json:"yield"`
}

func (s *Farmland) TableName() string {
	return "tb_farmland"
}

/*
新建角色农场类型建筑
*/
func NewRoleAllBFarmlands(roleId uint32) [] *Farmland{
	arr := make([] *Farmland, 16)
	for i:=0; i<16; i++ {
		d := Farmland{}
		d.Name = fmt.Sprintf("农场%d", i+1)
		d.Type = 0
		d.Level = 1
		d.RoleId = roleId
		d.Yield = 1000
		arr[i] = &d
	}
	return arr
}

func InsertFarmlandsToDB(arr []*Farmland) []*Farmland{
	orm.NewOrm().InsertMulti(1, arr)
	return arr
}

func ReadFarmlands(roleId uint32) []*Farmland{
	var farmlands []*Farmland
	qry := orm.NewOrm().QueryTable(&Farmland{}).Filter("role_id", roleId)
	qry.All(&farmlands)
	return farmlands
}

func UpdateFarmland(b *Farmland) {
	orm.NewOrm().Update(b)
}