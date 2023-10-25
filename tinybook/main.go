package main

func main() {
	engine := InitWebServer()
	engine.Run(":8081")
}
