package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IDeliveryFeeService interface {
	Get(ctx context.Context, req *pb.GetDeliveryFeeListRequest) (res []dto.DeliveryFeeResponse, total int64, err error)
	GetDetail(ctx context.Context, req *pb.GetDeliveryFeeDetailRequest) (res dto.DeliveryFeeResponse, err error)
	GetGP(ctx context.Context, req *pb.GetDeliveryFeeGPListRequest) (res *pb.GetDeliveryFeeGPListResponse, err error)
}

type DeliveryFeeService struct {
	opt                   opt.Options
	RepositoryDeliveryFee repository.IDeliveryFeeRepository
}

func NewDeliveryFeeService() IDeliveryFeeService {
	return &DeliveryFeeService{
		opt:                   global.Setup.Common,
		RepositoryDeliveryFee: repository.NewDeliveryFeeRepository(),
	}
}

func (s *DeliveryFeeService) Get(ctx context.Context, req *pb.GetDeliveryFeeListRequest) (res []dto.DeliveryFeeResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryFeeService.Get")
	defer span.End()

	var DeliveryFeees []*model.DeliveryFee
	DeliveryFeees, total, err = s.RepositoryDeliveryFee.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, DeliveryFee := range DeliveryFeees {
		res = append(res, dto.DeliveryFeeResponse{
			ID:            DeliveryFee.ID,
			Code:          DeliveryFee.Code,
			Name:          DeliveryFee.Name,
			Note:          DeliveryFee.Note,
			Status:        DeliveryFee.Status,
			MinimumOrder:  DeliveryFee.MinimumOrder,
			DeliveryFee:   DeliveryFee.DeliveryFee,
			RegionId:      DeliveryFee.RegionId,
			CutomerTypeId: DeliveryFee.CutomerTypeId,
		})
	}

	return
}

func (s *DeliveryFeeService) GetDetail(ctx context.Context, req *pb.GetDeliveryFeeDetailRequest) (res dto.DeliveryFeeResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryFeeService.GetDetail")
	defer span.End()

	var DeliveryFee *model.DeliveryFee
	DeliveryFee, err = s.RepositoryDeliveryFee.GetDetail(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.DeliveryFeeResponse{
		ID:            DeliveryFee.ID,
		Code:          DeliveryFee.Code,
		Name:          DeliveryFee.Name,
		Note:          DeliveryFee.Note,
		Status:        DeliveryFee.Status,
		MinimumOrder:  DeliveryFee.MinimumOrder,
		DeliveryFee:   DeliveryFee.DeliveryFee,
		RegionId:      DeliveryFee.RegionId,
		CutomerTypeId: DeliveryFee.CutomerTypeId,
	}

	return
}

func (s *DeliveryFeeService) GetGP(ctx context.Context, req *pb.GetDeliveryFeeGPListRequest) (res *pb.GetDeliveryFeeGPListResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryFeeService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.GnlRegion != "" {
		req.GnlRegion = url.PathEscape(req.GnlRegion)
		params["gnl_region"] = req.GnlRegion
	}

	if req.GnlCustTypeId != "" {
		params["gnl_cust_type_id"] = req.GnlCustTypeId
	}

	// if req.Minorqty != 0 {
	// 	params["minorqty"] = req.Minorqty
	// }
	// if req.GnlDeliveryFee != 0.0 {
	// 	params["gnl_delivery_fee"] = req.GnlDeliveryFee
	// }

	err = global.HttpRestApiToMicrosoftGP("GET", "DeliveryFee/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
