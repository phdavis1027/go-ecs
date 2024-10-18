package main

import (
	"log"
	"os"

	"github.com/phdavis1027/goecs/entity"
) 

type App struct {
  logger  *log.Logger 
  ecs     *entity.ECS
  name    string
}

func NewApp(name string, ecsCap int) *App {
  return &App{
    logger: log.New(os.Stdout, name, log.LstdFlags),
    ecs: entity.CreateEcsOfCapacity(ecsCap),
    name: name,
  }
}

func (app *App) Run() {
  log.Println("Running app")



  for _, system := range app.ecs.Systems {
    // SAFETY: OnTock is only ever read-only w/r/t the ECS 
    system.OnTick(app.ecs) 
  }
}
