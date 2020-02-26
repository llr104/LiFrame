package slgdb

import (
	"github.com/llr104/LiFrame/core/orm"
)

func Init() {
	orm.RegisterModel(new(Role))
	orm.RegisterModel(new(Barrack))
	orm.RegisterModel(new(Dwelling))
	orm.RegisterModel(new(Farmland))
	orm.RegisterModel(new(Lumber))
	orm.RegisterModel(new(Mine))
	orm.RegisterModel(new(General))
	orm.RegisterModel(new(City))
}
