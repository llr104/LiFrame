package slgdb

import (
	"fmt"
	"github.com/llr104/LiFrame/core/orm"
)

type Barrack struct {
	Id      	int      `json:"id"`
	Name        string   `json:"name"`
	RoleId      uint32   `json:"roleId"`
	Level       int8     `json:"level"`
	Type        int8     `json:"type"`
	Yield       uint32   `json:"yield"`
}

func (s *Barrack) TableName() string {
	return "tb_barrack"
}

/*
新建角色兵营类型建筑
*/
func NewRoleAllBarracks(roleId uint32) [] Barrack{
	arr := make([] Barrack, 16)
	for i:=0; i<16; i++ {
		d := Barrack{}
		d.Name = fmt.Sprintf("兵营%d", i)
		d.Type = 0
		d.Level = 1
		d.RoleId = roleId
		d.Yield = 1000
		arr[i] = d
	}
	return arr
}

func InsertBarracksToDB(arr []Barrack) []Barrack{
	orm.NewOrm().InsertMulti(len(arr), arr)
	return arr
}


