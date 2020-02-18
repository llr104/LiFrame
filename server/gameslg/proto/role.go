package proto

import (
	"github.com/llr104/LiFrame/proto"
	"github.com/llr104/LiFrame/server/db/slgdb"
)

type QryRoleReq struct {

}

type QryRoleAck struct {
	proto.BaseAck
	Role  slgdb.Role   `json:"role"`
}

type NewRoleReq struct{
	Name 	string
	Nation  int8
}

type NewRoleAck struct{
	proto.BaseAck
	Role  slgdb.Role   `json:"role"`
}