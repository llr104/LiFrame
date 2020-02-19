package slgdb

import "github.com/llr104/LiFrame/core/orm"

const (
	NationWei = iota
	NationShu
	NationWu
)

type Role struct {
	RoleId      int      `orm:"column(id);pk;auto" json:"roleId"`
	Name        string   `orm:"column(name);unique;size(16)" json:"name"`
	Nation      int8     `json:"nation"`
	Gold        uint32   `json:"gold"`
	Silver      uint32   `json:"silver"`
	Mine        uint32   `json:"mine"`
	Wood        uint32   `json:"wood"`
	UserId      uint32   `json:"userId"`
}

func (s *Role) TableName() string {
	return "tb_role"
}

func NewDefaultRole() Role{

	/*
	初始数据先写死，后面会改成配置，先做功能先
	*/
	r := Role{}
	r.Gold = 100000
	r.Silver = 100000
	r.Mine = 100000

	return r
}

func FindRoleByName(s *Role) error{
	return orm.NewOrm().Read(s, "name")
}

func FindRoleByUserId(s *Role) error{
	return orm.NewOrm().Read(s, "user_id")
}

func InsertRoleToDB(s *Role) (int64, error){
	return orm.NewOrm().Insert(s)
}
