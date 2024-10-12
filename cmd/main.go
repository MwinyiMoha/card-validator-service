package main

import (
	"card-validator-service/internal/api"
	"card-validator-service/internal/core/application"
	protos "card-validator-service/internal/gen"
	"card-validator-service/internal/validation"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpclogging "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/mwinyimoha/card-validator-utils/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const serverPort = 8080

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	logger, err := logging.NewLoggerConfig().BuildLogger()
	if err != nil {
		log.Fatal("could not initialize logging", err)
	}

	defer logger.Sync()

	val := validation.New()
	svc := application.NewService(val)
	srv := api.NewServer(svc)

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpclogging.UnaryServerInterceptor(api.InterceptorLogger(logger)),
			recovery.UnaryServerInterceptor(),
		),
	)
	protos.RegisterCardValidatorServiceServer(s, srv)
	reflection.Register(s)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", serverPort))
	if err != nil {
		logger.Fatal("could not create listener", zap.String("original_error", err.Error()))
	}

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func(c chan os.Signal) {
		logger.Info(
			"starting gRPC server",
			zap.Int("port", serverPort),
		)

		if err := s.Serve(lis); err != nil {
			logger.Error("could not start server", zap.String("original_error", err.Error()))
			c <- syscall.SIGTERM
		}
	}(ch)

	received := <-ch

	func() {
		logger.Info("initiating graceful shutdown", zap.String("OS signal", received.String()))
		s.GracefulStop()
	}()
}
