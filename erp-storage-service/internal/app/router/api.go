package router

import "git.edenfarm.id/project-version3/erp-services/erp-storage-service/internal/app/handler"

var handlers = map[string]RouteHandlers{}

func init() {
	handlers["health-check"] = &handler.HealthCheckHandler{}
	handlers["upload"] = &handler.UploadHandler{}
}
