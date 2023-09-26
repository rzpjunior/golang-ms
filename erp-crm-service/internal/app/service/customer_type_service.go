package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ICustomerTypeService interface {
	Get(ctx context.Context, req *dto.CustomerTypeGetListRequest) (res []*dto.CustomerTypeResponse, total int64, err error)
	GetDetail(ctx context.Context, id string) (res *dto.CustomerTypeResponse, err error)
}

type CustomerTypeService struct {
	opt opt.Options
}

func NewCustomerTypeService() ICustomerTypeService {
	return &CustomerTypeService{
		opt: global.Setup.Common,
	}
}

func (s *CustomerTypeService) Get(ctx context.Context, req *dto.CustomerTypeGetListRequest) (res []*dto.CustomerTypeResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerType.Get")
	defer span.End()

	var (
		customerTypes *bridge_service.GetCustomerTypeGPResponse
		status        int8
		statusGP      string
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

	customerTypes, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPList(ctx, &bridge_service.GetCustomerTypeGPListRequest{
		Limit:       int32(req.Limit),
		Offset:      int32(req.Offset),
		Description: req.Search,
		Inactive:    statusGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range customerTypes.Data {

		if v.Inactive == 0 {
			status = statusx.ConvertStatusName(statusx.Active)
		} else {
			status = statusx.ConvertStatusName(statusx.Archived)
		}
		res = append(res, &dto.CustomerTypeResponse{
			ID:            v.GnL_Cust_Type_ID,
			Code:          v.GnL_Cust_Type_ID,
			Description:   v.GnL_CustType_Description,
			CustomerGroup: v.GnL_Cust_GroupDesc,
			Status:        status,
			ConvertStatus: statusx.ConvertStatusValue(status),
		})
	}

	total = int64(customerTypes.TotalRecords)

	return
}

func (s *CustomerTypeService) GetDetail(ctx context.Context, id string) (res *dto.CustomerTypeResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerType.GetDetail")
	defer span.End()

	var (
		customerType       *bridge_service.GetCustomerTypeGPResponse
		statusCustomerType int8
	)
	customerType, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if customerType.Data[0].Inactive == 0 {
		statusCustomerType = statusx.ConvertStatusName(statusx.Active)
	} else {
		statusCustomerType = statusx.ConvertStatusName(statusx.Archived)
	}

	res = &dto.CustomerTypeResponse{
		ID:            customerType.Data[0].GnL_Cust_Type_ID,
		Code:          customerType.Data[0].GnL_Cust_Type_ID,
		Description:   customerType.Data[0].GnL_CustType_Description,
		CustomerGroup: customerType.Data[0].GnL_Cust_GroupDesc,
		Status:        statusCustomerType,
		ConvertStatus: statusx.ConvertStatusValue(statusCustomerType),
	}

	return
}
