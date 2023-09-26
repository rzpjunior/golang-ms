package router

import (
	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/handler"
)

var handlers = map[string]RouteHandlers{}

func init() {
	handlers["health-check"] = &handler.HealthCheckHandler{}
	handlers["person"] = &handler.PersonHandler{}
}
