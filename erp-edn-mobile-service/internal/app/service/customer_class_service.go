package service

import (
	"context"

	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceCustomerClass() ICustomerClassService {
	m := new(CustomerClassService)
	m.opt = global.Setup.Common
	return m
}

type ICustomerClassService interface {
	GetCustomerClass(ctx context.Context, req dto.CustomerClassRequest) (res []*dto.CustomerClassResponse, total int64, err error)
}

type CustomerClassService struct {
	opt opt.Options
}

func (s *CustomerClassService) GetCustomerClass(ctx context.Context, req dto.CustomerClassRequest) (res []*dto.CustomerClassResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerClassService.GetCustomerClass")
	defer span.End()

	var custClassResponse *bridgeService.GetCustomerClassResponse
	custClassResponse, err = s.opt.Client.BridgeServiceGrpc.GetCustomerClassList(ctx, &bridgeService.GetCustomerClassListRequest{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Clasdscr: req.Search,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer_class")
		return
	}
	
	datas := []*dto.CustomerClassResponse{}
	for _, customerCls := range custClassResponse.Data {
		datas = append(datas, &dto.CustomerClassResponse{
			ID:                  customerCls.Classid,
			CreditLimitTypeDesc: customerCls.CrlmttypDesc,
			CreditLimitType:     customerCls.Crlmttyp,
			Description:         customerCls.Clasdscr,
			CreditLimitAmount:   customerCls.Crlmtamt,
		})
	}

	total = int64(len(datas))
	res = datas

	return
}
