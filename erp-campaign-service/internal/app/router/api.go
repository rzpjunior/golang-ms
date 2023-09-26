package router

import "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/handler"

var handlers = map[string]RouteHandlers{}

func init() {
	handlers["health_check"] = &handler.HealthCheckHandler{}
	handlers["banner"] = &handler.BannerHandler{}
	handlers["item_section"] = &handler.ItemSectionHandler{}
	handlers["push_notification"] = &handler.PushNotificationHandler{}
	handlers["membership"] = &handler.MembershipHandler{}
}
