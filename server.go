package main

import (
	"SecKill/data"
	"SecKill/engine"
	"fmt"
)

const port = 8000
func main() {
	router := engine.SeckillEngine()
	defer data.Close()
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		println("Error when running server. " + err.Error())
	}
}

