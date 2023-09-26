package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/service"
)

type CampaignGrpcHandler struct {
	Option                         global.HandlerOptions
	ServiceBanner                  service.IBannerService
	ServiceItemSection             service.IItemSectionService
	ServiceMembership              service.IMembershipService
	ServiceCustomerPointLog        service.ICustomerPointLogService
	ServiceTalon                   service.ITalonService
	ServiceCustomerPointSummary    service.ICustomerPointSummaryService
	ServicePushNotification        service.IPushNotificationService
	ServiceCustomerPointExpiration service.ICustomerPointExpirationService
}
