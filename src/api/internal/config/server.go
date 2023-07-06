package config

import "github.com/cldcvr/terrarium/src/pkg/confighelper"

// ServerHTTPPort port to run http server on
func ServerHTTPPort() int {
	return confighelper.MustGetInt("server.http_port")
}

// ServerGRPCPort port to run GRPC server on
func ServerGRPCPort() int {
	return confighelper.MustGetInt("server.grpc_port")
}
