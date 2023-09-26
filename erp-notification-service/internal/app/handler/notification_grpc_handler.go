package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-notification-service/internal/app/service"
)

type NotificationGrpcHandler struct {
	Option                               global.HandlerOptions
	ServicesNotificationTransaction      service.INotificationTransactionService
	ServicesNotificationCampaign         service.INotificationCampaignService
	ServicesNotificationSite             service.INotificationSiteService
	ServicesNotificationPurchaser        service.INotificationPurchaserService
	ServicesNotificationCancelSalesOrder service.INotificationCancelSalesOrderService
	ServicesNotificationPaymentReminder  service.INotificationPaymentReminderService
}
