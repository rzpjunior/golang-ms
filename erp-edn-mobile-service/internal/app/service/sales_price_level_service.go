package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceSalesPriceLevel() ISalesPriceLevelService {
	m := new(SalesPriceLevelService)
	m.opt = global.Setup.Common
	return m
}

type ISalesPriceLevelService interface {
	GetGP(ctx context.Context, req dto.GetSalesPriceLevelListRequest) (res []*dto.SalesPriceLevel, total int64, err error)
	GetDetaiGPlById(ctx context.Context, id string) (res *dto.SalesPriceLevel, err error)
}

type SalesPriceLevelService struct {
	opt opt.Options
}

func (s *SalesPriceLevelService) GetGP(ctx context.Context, req dto.GetSalesPriceLevelListRequest) (res []*dto.SalesPriceLevel, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPriceLevelService.GetGP")
	defer span.End()

	// get price level from bridge
	var priceLevels *bridgeService.GetSalesPriceLevelResponse
	priceLevels, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelList(ctx, &bridgeService.GetSalesPriceLevelListRequest{
		Limit:         int64(req.Limit),
		Offset:        int64(req.Offset),
		GnlCustTypeId: req.CustTypeID,
		GnlRegion:     req.RegionID,
		Prclevel:      req.PriceLevel,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "price level")
		return
	}

	datas := []*dto.SalesPriceLevel{}
	for _, pl := range priceLevels.Data {
		datas = append(datas, &dto.SalesPriceLevel{
			RegionID:   pl.GnlRegion,
			CustTypeID: pl.GnlCustTypeId,
			PriceLevel: pl.Prclevel,
		})
	}

	total = int64(priceLevels.TotalRecords)
	res = datas

	return
}

func (s *SalesPriceLevelService) GetDetaiGPlById(ctx context.Context, id string) (res *dto.SalesPriceLevel, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PriceLevelService.GetDetaiGPlById")
	defer span.End()

	// get price level from bridge
	var pl *bridgeService.GetSalesPriceLevelResponse
	pl, err = s.opt.Client.BridgeServiceGrpc.GetSalesPriceLevelDetail(ctx, &bridgeService.GetSalesPriceLevelDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "price level")
		return
	}

	if len(pl.Data) > 0 {
		res = &dto.SalesPriceLevel{
			RegionID:   pl.Data[0].GnlRegion,
			CustTypeID: pl.Data[0].GnlCustTypeId,
			PriceLevel: pl.Data[0].Prclevel,
		}
	}

	return
}
