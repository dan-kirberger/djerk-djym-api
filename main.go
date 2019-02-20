package main

func main() {
	app := App{}
	app.Initialize("mongodb://localhost:27017")
	app.Run(":8080")
}
