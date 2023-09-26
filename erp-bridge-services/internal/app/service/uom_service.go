package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IUomService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.UomResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.UomResponse, err error)
	GetGP(ctx context.Context, req *pb.GetUomGPListRequest) (res *pb.GetUomGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetUomGPDetailRequest) (res *pb.GetUomGPResponse, err error)
}

type UomService struct {
	opt           opt.Options
	RepositoryUom repository.IUomRepository
}

func NewUomService() IUomService {
	return &UomService{
		opt:           global.Setup.Common,
		RepositoryUom: repository.NewUomRepository(),
	}
}

func (s *UomService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.UomResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UomService.Get")
	defer span.End()

	var uoms []*model.Uom
	uoms, total, err = s.RepositoryUom.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, uom := range uoms {
		res = append(res, dto.UomResponse{
			ID:             uom.ID,
			Code:           uom.Code,
			Description:    uom.Description,
			Status:         uom.Status,
			DecimalEnabled: uom.DecimalEnabled,
			CreatedAt:      uom.CreatedAt,
			UpdatedAt:      uom.UpdatedAt,
		})
	}

	return
}

func (s *UomService) GetDetail(ctx context.Context, id int64, code string) (res dto.UomResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UomService.GetDetail")
	defer span.End()

	var uom *model.Uom
	uom, err = s.RepositoryUom.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.UomResponse{
		ID:             uom.ID,
		Code:           uom.Code,
		Description:    uom.Description,
		Status:         uom.Status,
		DecimalEnabled: uom.DecimalEnabled,
		CreatedAt:      uom.CreatedAt,
		UpdatedAt:      uom.UpdatedAt,
	}

	return
}

func (s *UomService) GetGP(ctx context.Context, req *pb.GetUomGPListRequest) (res *pb.GetUomGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UomService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Search != "" {
		req.Search = url.PathEscape(req.Search)
		params["umschdsc"] = req.Search
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "uom/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *UomService) GetDetailGP(ctx context.Context, req *pb.GetUomGPDetailRequest) (res *pb.GetUomGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UomService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "uom/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 {
		err = edenlabs.ErrorNotFound("uom")
	}

	return
}
