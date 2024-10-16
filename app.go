package main 

import (
  "log"
  "os"
) 

type App struct {
  logger  *log.Logger 
  name    string
}

func NewApp(name string) *App {
  return &App{
    logger: log.New(os.Stdout, name, log.LstdFlags),
    name: name,
  }
}
