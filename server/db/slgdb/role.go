package slgdb

import "github.com/llr104/LiFrame/core/orm"

const (
	NationWei = iota
	NationShu
	NationWu
)

type Role struct {
	RoleId      int      `orm:"column(id);pk;auto"`
	Name        string   `orm:"column(name);unique;size(16)"`
	Nation      int8
	Gold        uint32
	Silver      uint32
	Iron        uint32
	Stone       uint32
	Wood        uint32
	Food        uint32
	UserId      uint32

}

func (s *Role) TableName() string {
	return "tb_role"
}

func NewDefaultRole() Role{
	r := Role{}
	r.Gold = 100000
	r.Food = 100000
	r.Iron = 100000
	r.Silver = 100000
	r.Stone = 100000
	r.Wood = 100000
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
