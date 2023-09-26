package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
)

type CustomerMobileGrpcHandler struct {
	Option              global.HandlerOptions
	ServiceUserCustomer service.IUserCustomerService
}
