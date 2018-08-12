package main

import (
	"fmt"
)

func main() {
	fmt.Println("Simple Inventory with Go")

	app := App{}
	app.Initalize()
	app.Run(":8082")
}
