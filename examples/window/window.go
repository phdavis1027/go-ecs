package main

import "github.com/phdavis1027/goecs/app"

func main() {
	app := app.NewApp("My App", 1024*4)

  app.Run()
}
