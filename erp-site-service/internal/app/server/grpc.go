package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/handler"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/service"

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/site_service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer() (err error) {
	url := fmt.Sprintf("%s:%d", global.Setup.Common.Config.Grpc.Host, global.Setup.Common.Config.Grpc.Port)
	listen, err := net.Listen("tcp", url)
	if err != nil {
		logrus.Errorf("[GRPC] Fail to start listen and server: %v", err)
		return
	}

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(func(method string, err error) bool {
			if global.Setup.Common.Config.App.Debug && err == nil {
				return false
			}
			return true
		}),
	}

	logger := logrus.NewEntry(global.Setup.Common.Logger.Logger())

	// init grpc server
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logger, opts...),
			grpc_recovery.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(logger, opts...),
			grpc_recovery.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(),
		)),
	)

	// setup handler
	siteGrpcHandler := &handler.SiteGrpcHandler{
		Option:               global.Setup,
		ServicesKoli:         service.NewKoliService(),
		ServicesDeliveryKoli: service.NewDeliveryKoliService(),
		ServicesPickingOrder: service.NewPickingOrderService(),
	}

	pb.RegisterSiteServiceServer(srv, siteGrpcHandler)
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())

	reflection.Register(srv)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logrus.Infoln("[GRPC] Server is shutting down")
			srv.GracefulStop()
			logrus.Info("[GRPC] Bye")
			return
		}
	}()

	logrus.Infof("[GRPC] GRPC serve at %s\n", url)
	srv.Serve(listen)

	return
}
