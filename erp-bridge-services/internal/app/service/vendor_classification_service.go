package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IVendorClassificationService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.VendorClassificationResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.VendorClassificationResponse, err error)
	GetGP(ctx context.Context, req *pb.GetVendorClassificationGPListRequest) (res *pb.GetVendorClassificationGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetVendorClassificationGPDetailRequest) (res *pb.GetVendorClassificationGPResponse, err error)
}

type VendorClassificationService struct {
	opt                            opt.Options
	RepositoryVendorClassification repository.IVendorClassificationRepository
}

func NewVendorClassificationService() IVendorClassificationService {
	return &VendorClassificationService{
		opt:                            global.Setup.Common,
		RepositoryVendorClassification: repository.NewVendorClassificationRepository(),
	}
}

func (s *VendorClassificationService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.VendorClassificationResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorClassificationService.Get")
	defer span.End()

	var vendorClassifications []*model.VendorClassification
	vendorClassifications, total, err = s.RepositoryVendorClassification.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, vendorClassification := range vendorClassifications {
		res = append(res, dto.VendorClassificationResponse{
			ID:            vendorClassification.ID,
			CommodityCode: vendorClassification.CommodityCode,
			CommodityName: vendorClassification.CommodityName,
			BadgeCode:     vendorClassification.BadgeCode,
			BadgeName:     vendorClassification.BadgeName,
			TypeCode:      vendorClassification.TypeCode,
			TypeName:      vendorClassification.TypeName,
		})
	}

	return
}

func (s *VendorClassificationService) GetDetail(ctx context.Context, id int64, code string) (res dto.VendorClassificationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorClassificationService.GetDetail")
	defer span.End()

	var vendorClassification *model.VendorClassification
	vendorClassification, err = s.RepositoryVendorClassification.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.VendorClassificationResponse{
		ID:            vendorClassification.ID,
		CommodityCode: vendorClassification.CommodityCode,
		CommodityName: vendorClassification.CommodityName,
		BadgeCode:     vendorClassification.BadgeCode,
		BadgeName:     vendorClassification.BadgeName,
		TypeCode:      vendorClassification.TypeCode,
		TypeName:      vendorClassification.TypeName,
	}

	return
}

func (s *VendorClassificationService) GetGP(ctx context.Context, req *pb.GetVendorClassificationGPListRequest) (res *pb.GetVendorClassificationGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorClassificationService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "VendorClassification/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *VendorClassificationService) GetDetailGP(ctx context.Context, req *pb.GetVendorClassificationGPDetailRequest) (res *pb.GetVendorClassificationGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorClassificationService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "VendorClassification/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
