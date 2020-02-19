package slgdb


type Barrack struct {
	Id      	int      `json:"id"`
	Name        string   `json:"name"`
	RoleId      uint32   `json:"roleId"`
	Level       int8     `json:"level"`
	Type        int8     `json:"type"`
	Yield       uint32   `json:"yield"`
}

func (s *Barrack) TableName() string {
	return "tb_barrack"
}


