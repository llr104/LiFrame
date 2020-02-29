package db
import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/llr104/LiFrame/core/orm"
	"github.com/llr104/LiFrame/server/db/dbobject"
	"github.com/llr104/LiFrame/utils"
)

func Init()  {

	dbConfig := utils.GlobalObject.AppConfig.DataBase
	if  dbConfig.Name == "" || dbConfig.User == "" || dbConfig.Password == "" || dbConfig.IP == "" || dbConfig.Port == 0{
		utils.Log.Info("no database")
	}else{
		dbPort := fmt.Sprintf("%d", dbConfig.Port)
		dbInfo := dbConfig.User+":"+dbConfig.Password+"@tcp("+dbConfig.IP+":"+dbPort+")/"+dbConfig.Name+"?charset=utf8"

		maxIdle := 30
		maxConn := 30

		orm.RegisterDriver("mysql", orm.DRMySQL)
		orm.RegisterDataBase("default","mysql",dbInfo,maxIdle,maxConn)

		dbobject.Init()
		_, err:= orm.NewOrm().QueryTable(dbobject.User{}).Count()
		if err != nil{
			utils.Log.Error("db err:%s",err.Error())
		}else{
			utils.Log.Info("db is connect")
		}

		user := dbobject.User{}
		user.Name = "test01"
		user.Password = "123456"

		if err := dbobject.FindUserByNP(&user); err!=nil {
			utils.Log.Info("initDataBase FindUserByNP error:",err.Error())
		}
	}
}