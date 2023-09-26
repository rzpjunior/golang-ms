package router

import "git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/handler"

var handlers = map[string]RouteHandlers{}

func init() {
	handlers["health_check"] = &handler.HealthCheckHandler{}
	handlers["item_category"] = &handler.ItemCategoryHandler{}
	handlers["item"] = &handler.ItemHandler{}
	handlers["item_section"] = &handler.ItemSectionHandler{}
	handlers["uom"] = &handler.UomHandler{}
	handlers["item_class"] = &handler.ItemClassHandler{}
}
