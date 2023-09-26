package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
)

type ISalesPriceLevelService interface {
	Get(ctx context.Context, req *dto.GetSalesPriceLevelRequest) (res []*dto.SalesPriceLevelResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res *dto.SalesPriceLevelResponse, err error)
}

type SalesPriceLevelService struct {
	opt opt.Options
}

func NewServiceSalesPriceLevel() ISalesPriceLevelService {
	return &SalesPriceLevelService{
		opt: global.Setup.Common,
	}
}

func (s *SalesPriceLevelService) Get(ctx context.Context, req *dto.GetSalesPriceLevelRequest) (res []*dto.SalesPriceLevelResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPriceLevelService.Get")
	defer span.End()

	var (
		salesPriceLevels *bridgeService.GetSalesPriceLevelResponse
	)

	if salesPriceLevels, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelList(ctx, &bridgeService.GetSalesPriceLevelListRequest{
		Limit:         req.Limit,
		Offset:        req.Offset,
		GnlCustTypeId: req.CustomerTypeID,
		GnlRegion:     req.RegionID,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales price level")
		return
	}

	for _, salesPriceLevel := range salesPriceLevels.Data {
		res = append(res, &dto.SalesPriceLevelResponse{
			ID:             salesPriceLevel.Prclevel,
			Description:    salesPriceLevel.Prclevel,
			CustomerTypeID: salesPriceLevel.GnlCustTypeId,
			RegionID:       salesPriceLevel.GnlRegion,
		})
	}

	total = int64(salesPriceLevels.TotalRecords)

	return
}

func (s *SalesPriceLevelService) GetDetail(ctx context.Context, id string) (res *dto.SalesPriceLevelResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPriceLevelService.GetSalesPriceLevel")
	defer span.End()

	var salesPriceLevel *bridgeService.GetSalesPriceLevelResponse

	if salesPriceLevel, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelDetail(ctx, &bridgeService.GetSalesPriceLevelDetailRequest{
		Id: id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales price level")
		return
	}

	res = &dto.SalesPriceLevelResponse{
		ID:             salesPriceLevel.Data[0].Prclevel,
		Description:    salesPriceLevel.Data[0].Prclevel,
		CustomerTypeID: salesPriceLevel.Data[0].GnlCustTypeId,
		RegionID:       salesPriceLevel.Data[0].GnlRegion,
	}

	return
}
