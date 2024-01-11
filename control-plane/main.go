package main

import (
	"context"
	cpc "control-plane/config"
	"control-plane/db"
	ecpl "control-plane/envoyserver"

	"github.com/envoyproxy/go-control-plane/pkg/server/v3"
)

func main() {
	db.InitDb()
	cache := ecpl.InitCache()
	// Run the xDS server
	ctx := context.Background()
	cb := &ecpl.Callbacks{Debug: false}

	srv := server.NewServer(ctx, cache, cb)

	ecpl.RunServer(ctx, srv, cpc.Config.Port)
}
