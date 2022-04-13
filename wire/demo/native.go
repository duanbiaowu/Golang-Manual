package demo

import "fmt"

type Monster struct {
	Name string
}

func NewMonster() Monster {
	return Monster{"Kitty"}
}

type Player struct {
	Name string
}

func NewPlayer(name string) Player {
	return Player{name}
}

type Mission struct {
	Player  Player
	Monster Monster
}

func NewMission(p Player, m Monster) Mission {
	return Mission{p, m}
}

func (m Mission) Start() {
	fmt.Printf("%s defeats %s, world peace!\n", m.Player.Name, m.Monster.Name)
}

func InitMissionNative(name string) Mission {
	player := NewPlayer(name)
	monster := NewMonster()
	return NewMission(player, monster)
}
