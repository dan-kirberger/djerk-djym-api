package main

import "log"

func main() {
	log.Println("Lets start")
	app := App{}
	app.Initialize()
	app.Run(":8080")
}
