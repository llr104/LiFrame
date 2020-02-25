package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/llr104/LiFrame/core/orm"
	"github.com/llr104/LiFrame/server/db/dbobject"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
	"github.com/llr104/LiFrame/utils"
)

func main() {

	dbInfo := "root:123456abc@tcp(127.0.0.1:3306)/li_db?charset=utf8"

	maxIdle := 30
	maxConn := 30

	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default","mysql",dbInfo,maxIdle,maxConn)
	dbobject.Init()
	slgdb.Init()

	orm.RunCommand()

	_, err:= orm.NewOrm().QueryTable(dbobject.User{}).Count()
	if err != nil{
		utils.Log.Error("db err:%s",err.Error())
	}else{
		utils.Log.Info("db is connect")
	}

}
