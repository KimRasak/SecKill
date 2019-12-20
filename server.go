package main

import (
	"SecKill/data"
	"SecKill/engine"
	"fmt"
	"net/http"
	_ "net/http/pprof"
)

const port = 20080
func main() {
	router := engine.SeckillEngine()
	defer data.Close()

	go func() {
		fmt.Println("pprof start...")
		fmt.Println(http.ListenAndServe(":9876", nil))
	}()

	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		println("Error when running server. " + err.Error())
	}
}

