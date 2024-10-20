package main

import (
	"log"
	"os"

	"github.com/phdavis1027/goecs/entity"
)

type App struct {
	logger *log.Logger
	ecs    *entity.ECS
	name   string
}

func NewApp(name string, ecsCap int) *App {
	return &App{
		logger: log.New(os.Stdout, name, log.LstdFlags),
		ecs:    entity.CreateEcsOfCapacity(ecsCap),
		name:   name,
	}
}

func (app *App) Run() {
	app.logger.Printf("Running app %s\n", app.name)

	app.logger.Println("Running attached Systems")

}
