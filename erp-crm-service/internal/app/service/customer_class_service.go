package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ICustomerClassService interface {
	Get(ctx context.Context, req *dto.CustomerClassGetListRequest) (res []*dto.CustomerClassResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res *dto.CustomerClassResponse, err error)
}

type CustomerClassService struct {
	opt opt.Options
}

func NewCustomerClassService() ICustomerClassService {
	return &CustomerClassService{
		opt: global.Setup.Common,
	}
}

func (s *CustomerClassService) Get(ctx context.Context, req *dto.CustomerClassGetListRequest) (res []*dto.CustomerClassResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerClass.Get")
	defer span.End()

	var (
		customerClasss *bridge_service.GetCustomerClassResponse
	)

	customerClasss, err = s.opt.Client.BridgeServiceGrpc.GetCustomerClassList(ctx, &bridge_service.GetCustomerClassListRequest{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Clasdscr: req.Search,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range customerClasss.Data {

		res = append(res, &dto.CustomerClassResponse{
			ID:          v.Classid,
			Code:        v.Classid,
			Description: v.Clasdscr,
		})
	}

	total = int64(customerClasss.TotalRecords)

	return
}

func (s *CustomerClassService) GetDetail(ctx context.Context, id string) (res *dto.CustomerClassResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerClass.GetDetail")
	defer span.End()

	var (
		customerClass *bridge_service.GetCustomerClassResponse
	)
	customerClass, err = s.opt.Client.BridgeServiceGrpc.GetCustomerClassDetail(ctx, &bridge_service.GetCustomerClassDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.CustomerClassResponse{
		ID:          customerClass.Data[0].Classid,
		Code:        customerClass.Data[0].Classid,
		Description: customerClass.Data[0].Clasdscr,
	}

	return
}
