package router

import (
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/handler"
)

func init() {
	handlers["auth"] = &handler.AuthHandler{}
	handlers["config"] = &handler.ConfigHandler{}
	handlers["adm_division"] = &handler.AdmDivisionHandler{}
	handlers["item"] = &handler.ItemHandler{}
	handlers["vendor"] = &handler.VendorHandler{}
	handlers["vendor_organization"] = &handler.VendorOrganizationHandler{}
	handlers["purchase/plan"] = &handler.PurchasePlanHandler{}
	handlers["purchase/order"] = &handler.PurchaseOrderHandler{}
	handlers["purchase/order/consolidated_shipment"] = &handler.ConsolidatedShipmentHandler{}
	handlers["payment_term"] = &handler.PaymentTermHandler{}
	handlers["payment_method"] = &handler.PaymentMethodHandler{}
	handlers["upload"] = &handler.UploadHandler{}
	handlers["field_purchaser"] = &handler.FieldPurchaserHandler{}
}
