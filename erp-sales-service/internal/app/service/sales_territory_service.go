package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
)

type ISalesTerritoryService interface {
	Get(ctx context.Context, req dto.GetSalesTerritoryRequest) (res []dto.SalesTerritoryResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res dto.SalesTerritoryResponse, err error)
}

type SalesTerritoryService struct {
	opt opt.Options
}

func NewServiceSalesTerritory() ISalesTerritoryService {
	return &SalesTerritoryService{
		opt: global.Setup.Common,
	}
}

func (s *SalesTerritoryService) Get(ctx context.Context, req dto.GetSalesTerritoryRequest) (res []dto.SalesTerritoryResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesTerritoryService.Get")
	defer span.End()

	var territory *bridgeService.GetSalesTerritoryGPResponse

	if territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPList(ctx, &bridgeService.GetSalesTerritoryGPListRequest{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
		Search: req.Search,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales Person")
		return
	}

	for _, territory := range territory.Data {

		layout := "2006-01-02T15:04:05"
		createdAt, _ := time.Parse(layout, territory.Creatddt)
		updatedAt, _ := time.Parse(layout, territory.Modifdt)

		res = append(res, dto.SalesTerritoryResponse{
			ID:            territory.Salsterr,
			Description:   territory.Slterdsc,
			SalespersonID: territory.Slprsnid,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		})
	}

	total = int64(territory.TotalRecords)

	return
}

func (s *SalesTerritoryService) GetDetail(ctx context.Context, id string) (res dto.SalesTerritoryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesTerritoryService.GetSalesTerritory")
	defer span.End()

	var territory *bridgeService.GetSalesTerritoryGPResponse

	if territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
		Id: id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales territory")
		return
	}

	layout := "2006-01-02T15:04:05"
	createdAt, _ := time.Parse(layout, territory.Data[0].Creatddt)
	updatedAt, _ := time.Parse(layout, territory.Data[0].Modifdt)

	res = dto.SalesTerritoryResponse{
		ID:            territory.Data[0].Salsterr,
		Description:   territory.Data[0].Slterdsc,
		SalespersonID: territory.Data[0].Slprsnid,
		CreatedAt:     createdAt,
		UpdatedAt:     updatedAt,
	}

	return
}
