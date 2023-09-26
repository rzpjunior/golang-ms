package router

import (
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/handler"
)

var handlers = map[string]RouteHandlers{}

func init() {
	handlers["auth"] = &handler.AuthHandler{}
	handlers["archetype"] = &handler.ArchetypeHandler{}
	handlers["customer"] = &handler.CustomerHandler{}
	handlers["address"] = &handler.AddressHandler{}
	handlers["site"] = &handler.SiteHandler{}
	handlers["vendor"] = &handler.VendorHandler{}
	handlers["region"] = &handler.RegionHandler{}
	handlers["adm_division"] = &handler.AdmDivisionHandler{}
	handlers["item"] = &handler.ItemHandler{}
	handlers["purchase_order"] = &handler.PurchaseOrderHandler{}
	handlers["item/transfer"] = &handler.ItemTransferHandler{}
	handlers["receiving"] = &handler.ReceivingHandler{}
	handlers["sales/invoice"] = &handler.SalesInvoiceHandler{}
	handlers["sales/order"] = &handler.SalesOrderHandler{}
	handlers["sales/payment"] = &handler.SalesPaymentHandler{}
	handlers["sales/person"] = &handler.SalespersonHandler{}
	handlers["sales/price_level"] = &handler.SalesPriceLevelHandler{}
	handlers["wrt"] = &handler.WrtHandler{}
	handlers["payment/checkbook"] = &handler.CheckbookHandler{}
	handlers["payment/method"] = &handler.PaymentMethodHandler{}
	handlers["payment/term"] = &handler.PaymentTermHandler{}
	handlers["health_check"] = &handler.HealthCheckHandler{}
	handlers["upload"] = &handler.UploadHandler{}
	handlers["customer/class"] = &handler.CustomerClassHandler{}
}
