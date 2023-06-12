package config

import (
	"github.com/cldcvr/terrarium/api/pkg/env"
)

// LocalHTTPPort port to run http server on
func LocalHTTPPort() int {
	return env.GetEnvInt("LOCAL_HTTP_PORT", 3000)
}

// LocalGRPCPort port to run GRPC server on
func LocalGRPCPort() int {
	return env.GetEnvInt("LOCAL_GRPC_PORT", 10000)
}
