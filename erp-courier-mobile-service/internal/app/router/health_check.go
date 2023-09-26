package router

import (
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/handler"
)

func init() {
	handlers["health_check"] = &handler.HealthCheckHandler{}
}
