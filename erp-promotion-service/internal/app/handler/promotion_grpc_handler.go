package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/service"
)

type PromotionGrpcHandler struct {
	Option                  global.HandlerOptions
	ServicesVoucher         service.IVoucherService
	ServicesVoucherItem     service.IVoucherItemService
	ServicesVoucherLog      service.IVoucherLogService
	ServicesPriceTieringLog service.IPriceTieringLogService
}
