package router

import (
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/handler"
)

func init() {
	handlers["app"] = &handler.CourierAppHandler{}
	handlers["app/delivery"] = &handler.DeliveryRunReturnHandler{}
}
