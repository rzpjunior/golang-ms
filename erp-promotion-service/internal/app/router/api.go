package router

import "git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/handler"

var handlers = map[string]RouteHandlers{}

func init() {
	handlers["health_check"] = &handler.HealthCheckHandler{}
	handlers["voucher"] = &handler.VoucherHandler{}
}
