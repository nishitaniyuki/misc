package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net"
	"os"

	"github.com/nishitaniyuki/misc/go/grpc_logging/pb"

	grpc_interceptors_logging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	otel_grpc "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	otel_sdk_trace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	grpc_reflection "google.golang.org/grpc/reflection"
	proto_json "google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type greeterServer struct {
	pb.UnimplementedGreeterServer
}

func (s *greeterServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	fields := grpc_interceptors_logging.ExtractFields(ctx)
	slog.InfoContext(ctx, "Hello", fields...)
	return &pb.HelloReply{Message: "Hello, " + req.GetName()}, nil
}

func init() {
	slog.SetDefault(slog.New(NewLogWithTraceHandler(slog.NewJSONHandler(os.Stderr, nil))))
}

func main() {
	traceProvider := otel_sdk_trace.NewTracerProvider()
	otel.SetTracerProvider(traceProvider)

	lis, err := net.Listen("tcp", ":5051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpc.UnaryServerInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
				pbReq, ok := req.(proto.Message)
				if !ok {
					return handler(ctx, req)
				}
				jsonData, err := proto_json.MarshalOptions{UseProtoNames: true}.Marshal(pbReq)
				if err != nil {
					return handler(ctx, req)
				}
				var fields map[string]interface{}
				if err := json.Unmarshal(jsonData, &fields); err != nil {
					return handler(ctx, req)
				}
				return handler(grpc_interceptors_logging.InjectFields(ctx, grpc_interceptors_logging.Fields{"grpc.args", fields}), req)
			}),
			grpc_interceptors_logging.UnaryServerInterceptor(
				grpc_interceptors_logging.LoggerFunc(func(ctx context.Context, lvl grpc_interceptors_logging.Level, msg string, fields ...any) {
					slog.Log(ctx, slog.Level(lvl), msg, fields...)
				}),
				grpc_interceptors_logging.WithLogOnEvents(grpc_interceptors_logging.FinishCall),
			),
		),
		grpc.StatsHandler(otel_grpc.NewServerHandler(otel_grpc.WithTracerProvider(traceProvider))),
	)
	pb.RegisterGreeterServer(s, &greeterServer{})
	grpc_reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
