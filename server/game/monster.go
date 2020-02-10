package game

import (
	"fmt"
	"math/rand"
	"time"
)

type monster struct {
	Level int       `json:"level"`
	Hp    int       `json:"hp"`
	X     int       `json:"x"`
	Y     int       `json:"y"`
	Name  string    `json:"name"`
	Id    uint32    `json:"Id"`
}

type player struct {
	UserId  uint32    `json:"userId"`
	Name    string    `json:"name"`
	X       int       `json:"x"`
	Y       int       `json:"y"`
}

func newRandomMonster() *monster{
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(99)
	m := monster{}
	m.Level = n+1
	m.Hp = 10*m.Level

	x := rand.Intn(1280)
	y := rand.Intn(720)
	m.X = x
	m.Y = y

	na := rand.Intn(10)+1
	m.Name = fmt.Sprintf("monster %d", na)

	return &m
}

