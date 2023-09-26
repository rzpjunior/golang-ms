package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/service"
)

type LogisticGrpcHandler struct {
	Option                        global.HandlerOptions
	ServicesDeliveryRunSheet      service.IDeliveryRunSheetService
	ServicesDeliveryRunSheetItem  service.IDeliveryRunSheetItemService
	ServicesDeliveryRunReturn     service.IDeliveryRunReturnService
	ServicesDeliveryRunReturnItem service.IDeliveryRunReturnItemService
	ServicesAddressCoordinateLog  service.IAddressCoordinateLogService
	ServicesCourierLog            service.ICourierLogService
	ServicesMerchantDeliveryLog   service.IMerchantDeliveryLogService
	ServicesPostponeDeliveryLog   service.IPostponeDeliveryLogService
	ServicesControlTower          service.IControlTowerService
}
