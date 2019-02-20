package main

func main() {
	app := App{}
	app.Initialize("localhost")
	app.Run(":8080")
}
