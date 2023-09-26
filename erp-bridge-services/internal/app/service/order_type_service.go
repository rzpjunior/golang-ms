package service

import (
	"context"
	"strconv"

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

type IOrderTypeService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.OrderTypeResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.OrderTypeResponse, err error)
	GetGP(ctx context.Context, req *pb.GetOrderTypeGPListRequest) (res *pb.GetOrderTypeGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetOrderTypeGPDetailRequest) (res *pb.GetOrderTypeGPResponse, err error)
}

type OrderTypeService struct {
	opt                 opt.Options
	RepositoryOrderType repository.IOrderTypeRepository
}

func NewOrderTypeService() IOrderTypeService {
	return &OrderTypeService{
		opt:                 global.Setup.Common,
		RepositoryOrderType: repository.NewOrderTypeRepository(),
	}
}

func (s *OrderTypeService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.OrderTypeResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "OrderTypeService.Get")
	defer span.End()

	var orderTypes []*model.OrderType
	orderTypes, total, err = s.RepositoryOrderType.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, orderType := range orderTypes {
		res = append(res, dto.OrderTypeResponse{
			ID:            orderType.ID,
			Code:          orderType.Code,
			Description:   orderType.Description,
			Status:        orderType.Status,
			StatusConvert: statusx.ConvertStatusValue(orderType.Status),
			CreatedAt:     timex.ToLocTime(ctx, orderType.CreatedAt),
			UpdatedAt:     timex.ToLocTime(ctx, orderType.UpdatedAt),
		})
	}

	return
}

func (s *OrderTypeService) GetDetail(ctx context.Context, id int64, code string) (res dto.OrderTypeResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "OrderTypeService.GetDetail")
	defer span.End()

	var orderType *model.OrderType
	orderType, err = s.RepositoryOrderType.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.OrderTypeResponse{
		ID:            orderType.ID,
		Code:          orderType.Code,
		Description:   orderType.Description,
		Status:        orderType.Status,
		StatusConvert: statusx.ConvertStatusValue(orderType.Status),
		CreatedAt:     timex.ToLocTime(ctx, orderType.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, orderType.UpdatedAt),
	}

	return
}

func (s *OrderTypeService) GetGP(ctx context.Context, req *pb.GetOrderTypeGPListRequest) (res *pb.GetOrderTypeGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "OrderTypeService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "OrderType/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *OrderTypeService) GetDetailGP(ctx context.Context, req *pb.GetOrderTypeGPDetailRequest) (res *pb.GetOrderTypeGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "OrderTypeService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "OrderType/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
