// Copyright (c) Ollion
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"

	"github.com/cldcvr/terrarium/src/api/internal/config"
	"github.com/cldcvr/terrarium/src/api/service"
	"github.com/cldcvr/terrarium/src/api/transport"
	"github.com/cldcvr/terrarium/src/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/src/pkg/transporthelper"
	log "github.com/sirupsen/logrus"
)

func init() {
	config.LoggerConfig(log.StandardLogger())
	config.LoggerConfigDefault()
}

func main() {
	serviceInst, err := service.New()
	mustNotErr(err)

	transportInst := transport.NewTerrariumAPI(serviceInst)

	ctx := context.Background()
	server := transporthelper.NewServer(transporthelper.ServerOptions{
		HTTPPort: config.ServerHTTPPort(),
		GRPCPort: config.ServerGRPCPort(),
	})

	terrariumpb.RegisterTerrariumServiceServer(server.GRPCServer, transportInst)                // Registers transportInst as service implementation in server.GRPCServer
	err = terrariumpb.RegisterTerrariumServiceHandlerServer(ctx, server.HTTPMux, transportInst) // Adds Handlers to server.HttpMux for each endpoint in transportInst
	mustNotErr(err)

	// Now, that both server.GRPCServer & server.HttpMux are configured with transport layer, Run the server:
	err = server.Run(ctx)
	mustNotErr(err)
}

func mustNotErr(err error) {
	if err != nil {
		panic(err)
	}
}
