package main

import "go/project/app"

func main() {
	app := app.NewApp()

	app.StartServer()
}
