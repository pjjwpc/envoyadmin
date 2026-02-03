package main

import (
	_ "manage-plane/internal/config"
	"manage-plane/internal/db"
	"manage-plane/internal/router"
	"manage-plane/internal/service"
	"os"
)

func main() {
	db.InitDb()
	go service.StartNotifyServer()
	r := router.InitRouter()
	r.Run(":8080")
	os.Exit(1)
}
