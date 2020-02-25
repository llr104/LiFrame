package slgproto

import (
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
)

type QryRoleReq struct {
	Type      	int8		`json:"type"`
}

type QryRoleAck struct {
	proto.BaseAck
	Role slgdb.Role `json:"role"`
	Type int8       `json:"type"`
}

type NewRoleReq struct{
	Name 	string
	Nation  int8
}

type NewRoleAck struct{
	proto.BaseAck
	Role slgdb.Role `json:"role"`
}