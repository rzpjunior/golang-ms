package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/service"
)

type SiteGrpcHandler struct {
	Option               global.HandlerOptions
	ServicesKoli         service.IKoliService
	ServicesDeliveryKoli service.IDeliveryKoliService
	ServicesPickingOrder service.IPickingOrderService
}
