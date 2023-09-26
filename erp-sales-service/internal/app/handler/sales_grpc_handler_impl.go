package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/service"
)

type SalesGrpcHandler struct {
	Option                  global.HandlerOptions
	ServiceSalesOrder       service.ISalesOrderService
	ServiceSalesInvoice     service.ISalesInvoiceService
	ServicePaymentMethod    service.IPaymentMethodService
	ServicePaymentChannel   service.IPaymentChannelService
	ServicePaymentGroupComb service.IPaymentGroupCombService
}
