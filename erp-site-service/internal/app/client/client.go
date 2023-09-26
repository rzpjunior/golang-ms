package client

import (
	"git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/client"
)

// Clients all client object injected here
type Clients struct {
	AccountServiceGrpc       client.IAccountServiceGrpc
	BridgeServiceGrpc        client.IBridgeServiceGrpc
	ConfigurationServiceGrpc client.IConfigurationServiceGrpc
	CatalogServiceGrpc       client.ICatalogServiceGrpc
}
