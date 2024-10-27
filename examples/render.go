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

	if _, err := app.Init(); err != nil {
		panic(err)
	}

	app.Run()
}
