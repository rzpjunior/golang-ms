package provider

import (
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/client"
)

func ConnectCatalogServiceGrpc(opt ProviderOptions) (catalogClientService client.ICatalogServiceGrpc, err error) {
	// setup client catalog service
	catalogClientService, err = client.NewCatalogServiceGrpc(client.CatalogServiceGrpcOption{
		Host:                  opt.Common.Env.GetString("client.catalog_service_grpc.host"),
		Port:                  opt.Common.Env.GetInt("client.catalog_service_grpc.port"),
		Timeout:               opt.Common.Env.GetDuration("client.catalog_service_grpc.timeout"),
		MaxConcurrentRequests: opt.Common.Env.GetInt("client.catalog_service_grpc.max_concurrent_requests"),
		ErrorPercentThreshold: opt.Common.Env.GetInt("client.catalog_service_grpc.error_percent_threshold"),
		Tls:                   opt.Common.Env.GetBool("client.catalog_service_grpc.tls"),
		PemPath:               opt.Common.Env.GetString("client.catalog_service_grpc.pem_tls_path"),
		Secret:                opt.Common.Env.GetString("client.catalog_service_grpc.secret"),
		Realtime:              opt.Common.Env.GetBool("client.catalog_service_grpc.realtime"),
	})
	if err != nil {
		opt.Common.Logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to connect client, error connect to catalog service | %v", err)).Print()
	}
	return
}
