package slgdb

import (
	"fmt"
	"github.com/llr104/LiFrame/core/orm"
)

type Mine struct {
	Id      	int      `json:"Id"`
	Name        string   `json:"name" orm:"size(20)"`
	RoleId      uint32   `json:"roleId"`
	Level       int8     `json:"level"`
	Type        int8     `json:"type"`
	Yield       uint32   `json:"yield"`
}

func (s *Mine) TableName() string {
	return "tb_mine"
}

/*
新建角色矿场类型建筑
*/
func NewRoleAllMines(roleId uint32) [] *Mine{
	arr := make([] *Mine, 16)
	for i:=0; i<16; i++ {
		d := Mine{}
		d.Name = fmt.Sprintf("矿场%d", i+1)
		d.Type = 0
		d.Level = 1
		d.RoleId = roleId
		d.Yield = 1000
		arr[i] = &d
	}
	return arr
}

func InsertMinesToDB(arr []*Mine) []*Mine{
	orm.NewOrm().InsertMulti(1, arr)
	return arr
}


func ReadMines(roleId uint32) []*Mine{
	var mines []*Mine
	qry := orm.NewOrm().QueryTable(&Mine{}).Filter("role_id", roleId)
	qry.All(&mines)
	return mines
}

func UpdateMine(b *Mine) {
	orm.NewOrm().Update(b)
}