package main

import (
	"runtime"

	"github.com/phdavis1027/goecs/app"
	"github.com/pkg/profile"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	defer profile.Start(profile.GoroutineProfile).Stop()

	app.NewApp("Get Rect", 100).Run()
}
