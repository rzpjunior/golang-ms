package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ICourierVendorService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID int64) (res []dto.CourierVendorResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.CourierVendorResponse, err error)
	GetGP(ctx context.Context, req *pb.GetCourierVendorGPListRequest) (res *pb.GetCourierVendorGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetCourierVendorGPDetailRequest) (res *pb.GetCourierVendorGPResponse, err error)
}

type CourierVendorService struct {
	opt                     opt.Options
	RepositoryCourierVendor repository.ICourierVendorRepository
}

func NewCourierVendorService() ICourierVendorService {
	return &CourierVendorService{
		opt:                     global.Setup.Common,
		RepositoryCourierVendor: repository.NewCourierVendorRepository(),
	}
}

func (s *CourierVendorService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID int64) (res []dto.CourierVendorResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierVendorService.Get")
	defer span.End()

	var courierVendors []*model.CourierVendor
	courierVendors, total, err = s.RepositoryCourierVendor.Get(ctx, offset, limit, status, search, orderBy, siteID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, cVendor := range courierVendors {
		res = append(res, dto.CourierVendorResponse{
			ID:            cVendor.ID,
			Code:          cVendor.Code,
			Name:          cVendor.Name,
			Note:          cVendor.Note,
			Status:        cVendor.Status,
			StatusConvert: statusx.ConvertStatusValue(cVendor.Status),
			SiteID:        cVendor.SiteID,
		})
	}

	return
}

func (s *CourierVendorService) GetDetail(ctx context.Context, id int64, code string) (res dto.CourierVendorResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierVendorService.GetDetail")
	defer span.End()

	var courierVendor *model.CourierVendor
	courierVendor, err = s.RepositoryCourierVendor.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.CourierVendorResponse{
		ID:            courierVendor.ID,
		Code:          courierVendor.Code,
		Name:          courierVendor.Name,
		Note:          courierVendor.Note,
		Status:        courierVendor.Status,
		StatusConvert: statusx.ConvertStatusValue(courierVendor.Status),
		SiteID:        courierVendor.SiteID,
	}

	return
}

func (s *CourierVendorService) GetGP(ctx context.Context, req *pb.GetCourierVendorGPListRequest) (res *pb.GetCourierVendorGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierVendorService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Status == 0 || req.Status == 1 {
		params["status"] = strconv.Itoa(int(req.Status))
	}

	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}

	if req.GnlCourierVendorId != "" {
		params["gnl_courier_vendor_id"] = req.GnlCourierVendorId
	}

	if req.GnlCourierVendorName != "" {
		params["gnl_courier_vendor_name"] = req.GnlCourierVendorName
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "couriervendor/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *CourierVendorService) GetDetailGP(ctx context.Context, req *pb.GetCourierVendorGPDetailRequest) (res *pb.GetCourierVendorGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierVendorService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "couriervendor/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
