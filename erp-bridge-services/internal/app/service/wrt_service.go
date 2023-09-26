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

type IWrtService interface {
	Get(ctx context.Context, offset, limit int, regionID int64, search string) (res []dto.WrtResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.WrtResponse, err error)
	GetGP(ctx context.Context, req *pb.GetWrtGPListRequest) (res *pb.GetWrtGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetWrtGPDetailRequest) (res *pb.GetWrtGPResponse, err error)
}

type WrtService struct {
	opt              opt.Options
	RepositoryWrt    repository.IWrtRepository
	RepositoryRegion repository.IRegionRepository
}

func NewWrtService() IWrtService {
	return &WrtService{
		opt:              global.Setup.Common,
		RepositoryWrt:    repository.NewWrtRepository(),
		RepositoryRegion: repository.NewRegionRepository(),
	}
}

func (s *WrtService) Get(ctx context.Context, offset, limit int, regionID int64, search string) (res []dto.WrtResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.Get")
	defer span.End()

	var (
		wrtList      []*model.Wrt
		region       *model.Region
		regionDetail *dto.RegionResponse
	)

	wrtList, total, err = s.RepositoryWrt.Get(ctx, offset, limit, regionID, search)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, wrt := range wrtList {

		region, err = s.RepositoryRegion.GetDetail(ctx, wrt.RegionID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("region_id")
			return
		}

		regionDetail = &dto.RegionResponse{
			ID:          region.ID,
			Code:        region.Code,
			Description: region.Description,
		}

		res = append(res, dto.WrtResponse{
			ID:        wrt.ID,
			RegionID:  wrt.RegionID,
			Code:      wrt.Code,
			StartTime: wrt.StartTime,
			EndTime:   wrt.EndTime,
			Region:    regionDetail,
		})
	}

	return
}

func (s *WrtService) GetDetail(ctx context.Context, id int64, code string) (res dto.WrtResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.GetDetail")
	defer span.End()

	var (
		Wrt          *model.Wrt
		region       *model.Region
		regionDetail *dto.RegionResponse
	)
	Wrt, err = s.RepositoryWrt.GetDetail(ctx, id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	region, err = s.RepositoryRegion.GetDetail(ctx, Wrt.RegionID, "")
	if err != nil {
		err = edenlabs.ErrorInvalid("region_id")
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	regionDetail = &dto.RegionResponse{
		ID:          region.ID,
		Code:        region.Code,
		Description: region.Description,
		Status:      region.Status,
	}

	res = dto.WrtResponse{
		ID:        Wrt.ID,
		RegionID:  Wrt.RegionID,
		Code:      Wrt.Code,
		StartTime: Wrt.StartTime,
		EndTime:   Wrt.EndTime,
		Region:    regionDetail,
	}

	return
}

// FOR GRPC PART
func (s *WrtService) GetGP(ctx context.Context, req *pb.GetWrtGPListRequest) (res *pb.GetWrtGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	//TODO: IF THERES ANY PARAMS ADDED PLEASE FOLLOW THIS CODE
	if req.Status != "" {
		params["inactive"] = req.Status
	}

	if req.Search != "" {
		req.Search = url.PathEscape(req.Search)
		params["gnl_wrt_id"] = req.Search
	}

	if req.GnlRegion != "" {
		req.GnlRegion = url.PathEscape(req.GnlRegion)
		params["gnl_region"] = req.GnlRegion
	}

	if req.OrderBy != "" {
		req.OrderBy = url.PathEscape(req.OrderBy)
		params["orderby"] = req.OrderBy
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "wrt/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *WrtService) GetDetailGP(ctx context.Context, req *pb.GetWrtGPDetailRequest) (res *pb.GetWrtGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "WrtService.GetDetailGP")
	defer span.End()

	req.Id = url.PathEscape(req.Id)
	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "wrt/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 {
		err = edenlabs.ErrorNotFound("wrt")
	}

	return
}
