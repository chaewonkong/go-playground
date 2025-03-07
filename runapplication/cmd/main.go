package main

import (
	"runapplication/app"
	"runapplication/app/server"
)

func main() {
	app.RunApplication("gateway", &server.Factory{})
}
