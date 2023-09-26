package provider

import (
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/client"
)

func ConnectAccountServiceGrpc(opt ProviderOptions) (accountClientService client.IAccountServiceGrpc, err error) {
	// setup client account service
	accountClientService, err = client.NewAccountServiceGrpc(client.AccountServiceGrpcOption{
		Host:                  opt.Common.Env.GetString("client.account_service_grpc.host"),
		Port:                  opt.Common.Env.GetInt("client.account_service_grpc.port"),
		Timeout:               opt.Common.Env.GetDuration("client.account_service_grpc.timeout"),
		MaxConcurrentRequests: opt.Common.Env.GetInt("client.account_service_grpc.max_concurrent_requests"),
		ErrorPercentThreshold: opt.Common.Env.GetInt("client.account_service_grpc.error_percent_threshold"),
		Tls:                   opt.Common.Env.GetBool("client.account_service_grpc.tls"),
		PemPath:               opt.Common.Env.GetString("client.account_service_grpc.pem_tls_path"),
		Secret:                opt.Common.Env.GetString("client.account_service_grpc.secret"),
		Realtime:              opt.Common.Env.GetBool("client.account_service_grpc.realtime"),
	})
	if err != nil {
		opt.Common.Logger.AddMessage(log.FatalLevel, fmt.Sprintf("Failed to connect client, error connect to account service | %v", err)).Print()
	}
	return
}
