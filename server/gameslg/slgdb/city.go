package slgdb

import "github.com/llr104/LiFrame/core/orm"

type City struct {
	Id       	int
	CId      	int16           `json:"cId" orm:"column(cId);description(城市配置Id)"`
	Name     	string			`json:"name" orm:"size(20)"`
	Nation   	int8			`json:"nation"`
	Capital  	bool			`json:"capital"`
	Adjacent	string          `json:"adjacent" orm:"column(adjacent);description(相邻城市)"`
}

func (s *City) TableName() string {
	return "tb_city"
}

func InsertCityToDB(s *City) (int64, error){
	return orm.NewOrm().Insert(s)
}