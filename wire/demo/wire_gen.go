// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package demo

// Injectors from wire.go:

func InitMission(name string) Mission {
	player := NewPlayer(name)
	monster := NewMonster()
	mission := NewMission(player, monster)
	return mission
}
