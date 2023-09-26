package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-inventory-service/internal/app/service"
)

type CatalogGrpcHandler struct {
	Option               global.HandlerOptions
	ServicesItemImage    service.IItemImageService
	ServicesItemCategory service.IItemCategoryService
	ServicesItem         service.IItemService
}
