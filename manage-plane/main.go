package main

import (
	"manage-plane/router"
	"os"
)

func main() {
	r := router.InitRouter()
	r.Run(":8080")
	os.Exit(1)
}
