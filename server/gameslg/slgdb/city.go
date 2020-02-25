package slgdb

import "github.com/llr104/LiFrame/core/orm"

type City struct {
	Id       int
	Name     string				`json:"name"`
	Nation   int8				`json:"nation"`
}

func (s *City) TableName() string {
	return "tb_city"
}

func InsertCityToDB(s *City) (int64, error){
	return orm.NewOrm().Insert(s)
}