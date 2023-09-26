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

type ICustomerTypeService interface {
	Get(ctx context.Context, req dto.GetCustomerTypeRequest) (res []dto.CustomerTypeResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res dto.CustomerTypeResponse, err error)
}

type CustomerTypeService struct {
	opt opt.Options
}

func NewServiceCustomerType() ICustomerTypeService {
	return &CustomerTypeService{
		opt: global.Setup.Common,
	}
}

func (s *CustomerTypeService) Get(ctx context.Context, req dto.GetCustomerTypeRequest) (res []dto.CustomerTypeResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerTypeService.Get")
	defer span.End()

	var customerType *bridgeService.GetCustomerTypeGPResponse

	if customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPList(ctx, &bridgeService.GetCustomerTypeGPListRequest{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "Customer type")
		return
	}

	for _, customerType := range customerType.Data {
		res = append(res, dto.CustomerTypeResponse{
			ID:   customerType.GnL_Cust_Type_ID,
			Name: customerType.GnL_CustType_Description,
		})
	}

	total = int64(len(customerType.Data))

	return
}

func (s *CustomerTypeService) GetDetail(ctx context.Context, id string) (res dto.CustomerTypeResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerTypeService.GetCustomerType")
	defer span.End()

	var customerType *bridgeService.GetCustomerTypeGPResponse

	if customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridgeService.GetCustomerTypeGPDetailRequest{
		Id: id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "Customer type")
		return
	}

	res = dto.CustomerTypeResponse{
		ID:   customerType.Data[0].GnL_Cust_Type_ID,
		Name: customerType.Data[0].GnL_CustType_Description,
	}

	return
}
