package dbobject

import "LiFrame/core/orm"

func Init() {
	orm.RegisterModel(new(User))
}
