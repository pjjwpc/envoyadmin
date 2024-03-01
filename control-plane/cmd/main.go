package main

import (
	"context"
	cpc "control-plane/internal/config"
	"control-plane/internal/db"
	ecpl "control-plane/internal/envoyserver"

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
