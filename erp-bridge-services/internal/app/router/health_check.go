package router

import (
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/handler"
)

func init() {
	handlers["health_check"] = &handler.HealthCheckHandler{}
}
