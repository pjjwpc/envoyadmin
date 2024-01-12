package main

import (
	_ "manage-plane/config"
	"manage-plane/db"
	"manage-plane/router"
	"os"
)

func main() {
	db.InitDb()
	r := router.InitRouter()
	r.Run(":8080")
	os.Exit(1)
}
