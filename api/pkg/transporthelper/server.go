package transporthelper

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"runtime/debug"

	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rotisserie/eris"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type CtxKey string
type ctxPrivateKey string

const (
	MdKeyGRPCMethod    = "grpc_method"
	MdKeyHTTPMethod    = "http_method"
	MdKeyHTTPHost      = "http_host"
	MdKeyHTTPPath      = "http_path"
	MdKeyHTTPRemote    = "http_remote"
	MdKeyHTTPRequestID = "request_id"

	HTTPHdrRequestID = "X-Request-Id"
)

// Server represents the server object.
type Server struct {
	HTTPMux    *runtime.ServeMux
	HTTPMw     []func(http.Handler) http.Handler
	GRPCServer *grpc.Server
	Options    ServerOptions
}

// ServerOptions holds the server configuration options.
type ServerOptions struct {
	HTTPPort int
	GRPCPort int
}

// NewServer creates a new instance of the Server.
func NewServer(options ServerOptions) *Server {
	return &Server{
		HTTPMux:    createRuntimeServeMux(),
		HTTPMw:     []func(http.Handler) http.Handler{appendReqIDToHTTPContext},
		GRPCServer: createGRPCServer(),
		Options:    options,
	}
}

func createRuntimeServeMux() *runtime.ServeMux {
	return runtime.NewServeMux(
		runtime.WithForwardResponseOption(appendReqIDToHTTPHeader),
		runtime.WithMetadata(extractHTTPRequestMetadata),
		runtime.WithIncomingHeaderMatcher(runtime.DefaultHeaderMatcher),
		runtime.WithErrorHandler(runtime.DefaultHTTPErrorHandler),
		runtime.WithRoutingErrorHandler(runtime.DefaultRoutingErrorHandler),
		runtime.WithMarshalerOption("application/json", createJSONBodyMarshaler()),
	)
}

func extractHTTPRequestMetadata(ctx context.Context, r *http.Request) metadata.MD {
	return metadata.New(map[string]string{
		MdKeyHTTPMethod: r.Method,
		MdKeyHTTPHost:   r.URL.Host,
		MdKeyHTTPPath:   r.URL.Path,
		MdKeyHTTPRemote: r.RemoteAddr,
	})
}

func createJSONBodyMarshaler() *runtime.HTTPBodyMarshaler {
	return &runtime.HTTPBodyMarshaler{
		Marshaler: &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				EmitUnpopulated: true,
				UseEnumNumbers:  false,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				AllowPartial:   false,
				DiscardUnknown: false,
			},
		},
	}
}

func createGRPCServer() *grpc.Server {
	return grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			appendReqIDToGRPCContext,
			appendDetailsToGRPCError,
		),
	)
}

// Run starts the servers and blocks until all are running.
func (s *Server) Run(ctx context.Context) error {
	grpcChan, err := s.runGRPCAsync(ctx)
	if err != nil {
		return err
	}

	httpChan, err := s.runHTTPAsync(ctx)
	if err != nil {
		return err
	}

	select {
	case err = <-grpcChan:
		log.Default().Printf("gRPC server has stopped: %+v", err)
	case err = <-httpChan:
		log.Default().Printf("HTTP server has stopped: %+v", err)
	}

	return err
}
func (s *Server) runGRPCAsync(ctx context.Context) (chan error, error) {
	grpcHost := fmt.Sprintf(":%d", s.Options.GRPCPort)
	log.Default().Printf("starting gRPC server on %s...", grpcHost)

	return runAsync(ctx, grpcHost, s.GRPCServer)
}

func (s *Server) runHTTPAsync(ctx context.Context) (chan error, error) {
	var httpHandler http.Handler = s.HTTPMux
	for _, mw := range s.HTTPMw {
		httpHandler = mw(httpHandler)
	}

	httpHost := fmt.Sprintf(":%d", s.Options.HTTPPort)
	log.Default().Printf("starting HTTP server on %s...", httpHost)

	return runAsync(ctx, httpHost, &http.Server{Addr: httpHost, Handler: httpHandler})
}

func runAsync(ctx context.Context, host string, srv interface{ Serve(net.Listener) error }) (chan error, error) {
	listener, err := net.Listen("tcp", host)
	if err != nil {
		return nil, err
	}

	result := make(chan error)
	go func() {
		defer close(result)
		defer func() {
			if p := recover(); p != nil {
				result <- fmt.Errorf("server task has failed with unhandled panic: %v\n%s", p, debug.Stack())
			}
		}()
		result <- srv.Serve(listener)
	}()

	return result, nil
}

func appendReqIDToGRPCContext(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	ctx = appendRequestID(ctx)
	ctx = context.WithValue(ctx, ctxPrivateKey(MdKeyGRPCMethod), info.FullMethod)
	md := metadata.Pairs(HTTPHdrRequestID, GetReqIdFromCtx(ctx))
	grpc.SetHeader(ctx, md)
	return handler(ctx, req)
}

func appendDetailsToGRPCError(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	resp, err := handler(ctx, req)
	if err != nil {
		e := status.Convert(err)
		errMap, _ := structpb.NewStruct(eris.ToJSON(err, true))
		e.WithDetails(errMap)
		return resp, e.Err()
	}
	return resp, err
}

func appendReqIDToHTTPContext(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := appendRequestID(r.Context())
		req := r.WithContext(ctx)
		h.ServeHTTP(w, req)
	})
}

func appendRequestID(ctx context.Context) context.Context {
	key := ctxPrivateKey(MdKeyHTTPRequestID)

	if ctx.Value(key) != nil {
		return ctx
	}

	return context.WithValue(ctx, key, uuid.NewString())
}

func appendReqIDToHTTPHeader(ctx context.Context, w http.ResponseWriter, _ proto.Message) error {
	w.Header().Set(HTTPHdrRequestID, GetReqIdFromCtx(ctx))
	return nil
}

func GetReqIdFromCtx(ctx context.Context) string {
	val := ctx.Value(ctxPrivateKey(MdKeyHTTPRequestID))
	if val != nil {
		return val.(string)
	}

	return ""
}
