package main

import (
	"context"

	"github.com/cldcvr/terrarium/api/internal/config"
	"github.com/cldcvr/terrarium/api/pkg/pb/terrariumpb"
	"github.com/cldcvr/terrarium/api/pkg/transporthelper"
	"github.com/cldcvr/terrarium/api/service"
	"github.com/cldcvr/terrarium/api/transport"
	log "github.com/sirupsen/logrus"
)

func init() {
	config.SetupLogger(log.StandardLogger())
}

func main() {
	serviceInst := service.New()

	transportInst := transport.NewTerrariumAPI(serviceInst)

	ctx := context.Background()
	server := transporthelper.NewServer(transporthelper.ServerOptions{
		HTTPPort: config.LocalHTTPPort(),
		GRPCPort: config.LocalGRPCPort(),
	})

	terrariumpb.RegisterTerrariumServiceServer(server.GRPCServer, transportInst)                 // Registers transportInst as service implementation in server.GRPCServer
	err := terrariumpb.RegisterTerrariumServiceHandlerServer(ctx, server.HTTPMux, transportInst) // Adds Handlers to server.HttpMux for each endpoint in transportInst
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
