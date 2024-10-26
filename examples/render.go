package main

import "runtime"

import (
	"github.com/phdavis1027/goecs/app"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	app := app.NewApp("Get Rect", 100)

	app.Init()

	app.Run()
}
