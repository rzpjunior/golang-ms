package router

import "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/handler"

var handlers = map[string]RouteHandlers{}

func init() {
	handlers["health_check"] = &handler.HealthCheckHandler{}
	handlers["app"] = &handler.ApplicationConfigHandler{}
	handlers["region_policy"] = &handler.RegionPolicyHandler{}
	handlers["day_off"] = &handler.DayOffHandler{}
	handlers["glossary"] = &handler.GlossaryHandler{}
	handlers["wrt"] = &handler.WrtHandler{}
	handlers["region"] = &handler.RegionHandler{}
	handlers["adm_division"] = &handler.AdmDivisionHandler{}
}
