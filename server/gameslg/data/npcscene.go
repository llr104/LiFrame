package data

import (
	"fmt"
	"github.com/llr104/LiFrame/server/gameslg/slgdb"
	"math/rand"
	"time"
)

type NpcScene struct {
	Id       uint16
	Name     string           `json:"name"`
	Generals []*slgdb.General `json:"generals"`
}

func RandomNPCScene(Id uint16) *NpcScene {
	s := NpcScene{}
	s.Id = Id
	s.Name = fmt.Sprintf("npc场景 %d", s.Id)

	s.Generals = make([]*slgdb.General, 3)
	rand.Seed(time.Now().UnixNano())

	for i:=0; i<3; i++ {
		g := slgdb.RandomNPCNewGeneral()
		s.Generals[i] = g
	}

	return &s
}
