package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ICustomerTypeService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.CustomerTypeResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.CustomerTypeResponse, err error)
	GetGP(ctx context.Context, req *pb.GetCustomerTypeGPListRequest) (res *pb.GetCustomerTypeGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetCustomerTypeGPDetailRequest) (res *pb.GetCustomerTypeGPResponse, err error)
}

type CustomerTypeService struct {
	opt                    opt.Options
	RepositoryCustomerType repository.ICustomerTypeRepository
}

func NewCustomerTypeService() ICustomerTypeService {
	return &CustomerTypeService{
		opt:                    global.Setup.Common,
		RepositoryCustomerType: repository.NewCustomerTypeRepository(),
	}
}

func (s *CustomerTypeService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.CustomerTypeResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerTypeService.Get")
	defer span.End()

	var CustomerTypes []*model.CustomerType
	CustomerTypes, total, err = s.RepositoryCustomerType.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, CustomerType := range CustomerTypes {
		res = append(res, dto.CustomerTypeResponse{
			ID:            CustomerType.ID,
			Code:          CustomerType.Code,
			Description:   CustomerType.Description,
			GroupType:     CustomerType.GroupType,
			Abbreviation:  CustomerType.Abbreviation,
			Status:        CustomerType.Status,
			StatusConvert: statusx.ConvertStatusValue(CustomerType.Status),
			CreatedAt:     timex.ToLocTime(ctx, CustomerType.CreatedAt),
			UpdatedAt:     timex.ToLocTime(ctx, CustomerType.UpdatedAt),
		})
	}

	return
}

func (s *CustomerTypeService) GetDetail(ctx context.Context, id int64, code string) (res dto.CustomerTypeResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerTypeService.GetDetail")
	defer span.End()

	var CustomerType *model.CustomerType
	CustomerType, err = s.RepositoryCustomerType.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.CustomerTypeResponse{
		ID:            CustomerType.ID,
		Code:          CustomerType.Code,
		Description:   CustomerType.Description,
		GroupType:     CustomerType.GroupType,
		Abbreviation:  CustomerType.Abbreviation,
		Status:        CustomerType.Status,
		StatusConvert: statusx.ConvertStatusValue(CustomerType.Status),
		CreatedAt:     timex.ToLocTime(ctx, CustomerType.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, CustomerType.UpdatedAt),
	}

	return
}

func (s *CustomerTypeService) GetGP(ctx context.Context, req *pb.GetCustomerTypeGPListRequest) (res *pb.GetCustomerTypeGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerTypeService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.CustomerGroupId != "" {
		req.CustomerGroupId = url.PathEscape(req.CustomerGroupId)
		params["GNL_Cust_Group"] = req.CustomerGroupId
	}

	if req.Description != "" {
		req.Description = url.PathEscape(req.Description)
		params["GNL_CustType_Description"] = req.Description
	}

	if req.Inactive != "" {
		params["INACTIVE"] = req.Inactive
	}

	if req.Id != "" {
		req.Id = url.PathEscape(req.Id)
		params["gnl_cust_type_id"] = req.Id
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "CustomerType/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *CustomerTypeService) GetDetailGP(ctx context.Context, req *pb.GetCustomerTypeGPDetailRequest) (res *pb.GetCustomerTypeGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerTypeService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "CustomerType/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 {
		err = edenlabs.ErrorNotFound("customer_type")
	}

	return
}
