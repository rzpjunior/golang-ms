package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
)

type ISalesPersonService interface {
	Get(ctx context.Context, req dto.GetSalesPersonRequest) (res []dto.SalesPersonResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res dto.SalesPersonResponse, err error)
}

type SalesPersonService struct {
	opt opt.Options
}

func NewServiceSalesPerson() ISalesPersonService {
	return &SalesPersonService{
		opt: global.Setup.Common,
	}
}

func (s *SalesPersonService) Get(ctx context.Context, req dto.GetSalesPersonRequest) (res []dto.SalesPersonResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPersonService.Get")
	defer span.End()

	var (
		salesPerson *bridgeService.GetSalesPersonGPResponse
		status      int8
		statusGP    string
	)

	if req.Status != 0 {
		switch req.Status {
		case 1:
			statusGP = "0"
		case 7:
			statusGP = "1"
		default:
			statusGP = utils.ToString(req.Status)
		}
	}

	if salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPList(ctx, &bridgeService.GetSalesPersonGPListRequest{
		Limit:            int32(req.Limit),
		Offset:           int32(req.Offset),
		SalesTerritoryId: req.SalesTerritoryID,
		Status:           statusGP,
		Search:           req.Search,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales Person")
		return
	}

	for _, salesPerson := range salesPerson.Data {
		if salesPerson.Inactive == 0 {
			status = statusx.ConvertStatusName(statusx.Active)
		} else {
			status = statusx.ConvertStatusName(statusx.Archived)
		}

		res = append(res, dto.SalesPersonResponse{
			ID:               salesPerson.Slprsnid,
			Name:             salesPerson.Slprsnfn,
			MiddleName:       salesPerson.Sprsnsmn,
			LastName:         salesPerson.Sprsnsln,
			SalesTerritoryID: salesPerson.Salsterr,
			EmployeeID:       salesPerson.Employid,
			Status:           status,
			ConvertStatus:    statusx.ConvertStatusValue(status),
		})
	}

	total = int64(salesPerson.TotalRecords)

	return
}

func (s *SalesPersonService) GetDetail(ctx context.Context, id string) (res dto.SalesPersonResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPersonService.GetSalesPerson")
	defer span.End()

	var salesPerson *bridgeService.GetSalesPersonGPResponse

	if salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
		Id: id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales person")
		return
	}

	res = dto.SalesPersonResponse{
		ID:         salesPerson.Data[0].Slprsnid,
		Name:       salesPerson.Data[0].Slprsnfn,
		MiddleName: salesPerson.Data[0].Sprsnsmn,
		LastName:   salesPerson.Data[0].Sprsnsln,
	}

	return
}
