package main

import (
	_ "manage-plane/internal/config"
	"manage-plane/internal/db"
	"manage-plane/internal/router"
	"os"
)

func main() {
	db.InitDb()
	r := router.InitRouter()
	r.Run(":8080")
	os.Exit(1)
}
