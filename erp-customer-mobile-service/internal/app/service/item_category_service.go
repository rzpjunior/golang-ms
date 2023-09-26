package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
)

type IItemCategoryService interface {
	Get(ctx context.Context, req dto.ItemCategoryMobileRequest) (res []dto.ItemCategoryResponse, err error)
}

type ItemCategoryService struct {
	opt opt.Options
}

func NewItemCategoryService() IItemCategoryService {
	return &ItemCategoryService{
		opt: global.Setup.Common,
	}
}

func (s *ItemCategoryService) Get(ctx context.Context, req dto.ItemCategoryMobileRequest) (res []dto.ItemCategoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ItemCategoryService.Get")
	defer span.End()

	var (
		itemCategoryList *catalog_service.GetItemCategoryListResponse
		addressDetail    *bridge_service.GetAddressGPResponse
		admDivision      *bridge_service.GetAdmDivisionGPResponse
		regionID         string
	)

	// Set default area to jkt
	regionID = "Greater Jakarta"

	if req.Data.AddressID != "" {
		addressDetail, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
			Id: req.Data.AddressID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
			AdmDivisionCode: addressDetail.Data[0].AdministrativeDiv.GnlAdministrativeCode,
			Limit:           10,
			Offset:          1,
		})
		if err != nil || len(admDivision.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		regionID = admDivision.Data[0].Region
	}

	itemCategoryList, err = s.opt.Client.CatalogServiceGrpc.GetItemCategoryList(ctx, &catalog_service.GetItemCategoryListRequest{
		RegionId: regionID,
		Status:   int32(statusx.ConvertStatusName("Active")),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range itemCategoryList.Data {
		itemCategory := dto.ItemCategoryResponse{
			ID:       strconv.Itoa(int(v.Id)),
			Region:   v.RegionId,
			Name:     v.Name,
			ImageUrl: v.ImageUrl,
			Status:   strconv.Itoa(int(v.Status)),
		}
		res = append(res, itemCategory)
	}

	return
}
