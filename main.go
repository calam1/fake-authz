package main

import (
	"flag"
	"log"
	"net"

	auth "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
	auth2 "grainger.com/auth_proxy/v1/auth"
	"grainger.com/auth_proxy/v1/config"
	"grainger.com/auth_proxy/v1/health"
)

var (
	grpcport = flag.String("grpcport", ":50051", "grpcport")
	logger   *zap.Logger
)

func init() {
	log.Println("Initializing Config")
	config.ConfigSetup()
}

func initLogger() {
	cfg := zap.Config{
		Encoding:         "json",
		Level:            zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey: "message",

			LevelKey:    "level",
			EncodeLevel: zapcore.CapitalLevelEncoder,

			TimeKey:    "time",
			EncodeTime: zapcore.ISO8601TimeEncoder,

			CallerKey:    "caller",
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	logger, _ = cfg.Build()
	zap.ReplaceGlobals(logger)
}

func main() {
	initLogger()
	defer logger.Sync()

	lis, err := net.Listen("tcp", *grpcport)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{grpc.MaxConcurrentStreams(config.Configuration.ConcurrentStreams)}
	// opts = append(opts)

	s := grpc.NewServer(opts...)

	auth.RegisterAuthorizationServer(s, &auth2.AuthorizationServer{})
	healthpb.RegisterHealthServer(s, &health.HealthServer{})

	zap.L().Info("Starting gRPC ProductCore Auth Server", zap.String("grpcport", *grpcport))

	if config.Configuration.ShowReflection {
		reflection.Register(s)
	}

	s.Serve(lis)
}
