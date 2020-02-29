package dbobject

import "github.com/llr104/LiFrame/core/orm"

const (
	UserStateNormal = iota
	UserStateForbid
)

type User struct {
	Id      	  uint32
	Name     	  string	`orm:"size(64)"`
	Password 	  string	`orm:"size(128)"`
	LoginTimes    int
	LastLoginIp   string	`orm:"size(64)"`
	LastLoginTime int64
	LogoutTime	  int64
	IsOnline	  bool
	State      	  int8
	Gold		  int64
}

func (u *User) TableName() string {
	return "tb_user"
}

func FindUserById(u* User) error{
	return orm.NewOrm().Read(u, "id")
}

func FindUserByNP(u *User) error{
	return orm.NewOrm().Read(u, "name", "Password")
}

func FindUserByName(u *User) error{
	return orm.NewOrm().Read(u, "name")
}

func UpdateUserToDB(u *User) (int64, error){
	return orm.NewOrm().Update(u)
}

func UpdateUserGold(u *User) (int64, error) {
	return orm.NewOrm().Update(u, "gold")
}

func InsertUserToDB(u *User) (int64, error){
	return orm.NewOrm().Insert(u)
}
