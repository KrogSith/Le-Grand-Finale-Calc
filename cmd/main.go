package main

import "calculator/internal/application"

func main() {
	app := application.New()
	//app.Run()
	application.InitiateDatabase()
	app.RunServer()
}
