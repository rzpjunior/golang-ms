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

type IVendorOrganizationService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.VendorOrganizationResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.VendorOrganizationResponse, err error)
	GetGP(ctx context.Context, req *pb.GetVendorOrganizationGPListRequest) (res *pb.GetVendorOrganizationGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetVendorOrganizationGPDetailRequest) (res *pb.GetVendorOrganizationGPResponse, err error)
}

type VendorOrganizationService struct {
	opt                          opt.Options
	RepositoryVendorOrganization repository.IVendorOrganizationRepository
}

func NewVendorOrganizationService() IVendorOrganizationService {
	return &VendorOrganizationService{
		opt:                          global.Setup.Common,
		RepositoryVendorOrganization: repository.NewVendorOrganizationRepository(),
	}
}

func (s *VendorOrganizationService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.VendorOrganizationResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorOrganizationService.Get")
	defer span.End()

	var vendorOrganizations []*model.VendorOrganization
	vendorOrganizations, total, err = s.RepositoryVendorOrganization.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, vendorOrganization := range vendorOrganizations {
		res = append(res, dto.VendorOrganizationResponse{
			ID:                     vendorOrganization.ID,
			Code:                   vendorOrganization.Code,
			VendorClassificationID: vendorOrganization.VendorClassificationID,
			SubDistrictID:          vendorOrganization.SubDistrictID,
			PaymentTermID:          vendorOrganization.PaymentTermID,
			Name:                   vendorOrganization.Name,
			Address:                vendorOrganization.Address,
			Note:                   vendorOrganization.Note,
			Status:                 vendorOrganization.Status,
		})
	}

	return
}

func (s *VendorOrganizationService) GetDetail(ctx context.Context, id int64, code string) (res dto.VendorOrganizationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorOrganizationService.GetDetail")
	defer span.End()

	var vendorOrganization *model.VendorOrganization
	vendorOrganization, err = s.RepositoryVendorOrganization.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.VendorOrganizationResponse{
		ID:                     vendorOrganization.ID,
		Code:                   vendorOrganization.Code,
		VendorClassificationID: vendorOrganization.VendorClassificationID,
		SubDistrictID:          vendorOrganization.SubDistrictID,
		PaymentTermID:          vendorOrganization.PaymentTermID,
		Name:                   vendorOrganization.Name,
		Address:                vendorOrganization.Address,
		Note:                   vendorOrganization.Note,
		Status:                 vendorOrganization.Status,
	}

	return
}

func (s *VendorOrganizationService) GetGP(ctx context.Context, req *pb.GetVendorOrganizationGPListRequest) (res *pb.GetVendorOrganizationGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorOrganizationService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.PrpVendorOrgDesc != "" {
		params["prp_vendor_org_desc"] = req.PrpVendorOrgDesc
	}

	if req.Status != "" {
		params["Inactive"] = req.Status
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "vendororganization/getall", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *VendorOrganizationService) GetDetailGP(ctx context.Context, req *pb.GetVendorOrganizationGPDetailRequest) (res *pb.GetVendorOrganizationGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorOrganizationService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      url.PathEscape(req.Id),
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "vendororganization/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
