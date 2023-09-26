package router

import "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/handler"

func init() {
	handlers["packing_order"] = &handler.PackingOrderHandler{}
	handlers["packing_order/pack"] = &handler.PackingOrderPackHandler{}
}
