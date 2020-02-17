package dbobject

import (
	"github.com/llr104/LiFrame/core/orm"
)

func Init() {
	orm.RegisterModel(new(User))
}

