package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
)

type IVendorOrganizationService interface {
	Get(ctx context.Context, req *dto.VendorOrganizationListRequest) (res []*dto.VendorOrganizationResponse, total int64, err error)
	GetByID(ctx context.Context, id string) (res *dto.VendorOrganizationResponse, err error)
}

type VendorOrganizationService struct {
	opt opt.Options
}

func NewVendorOrganizationService() IVendorOrganizationService {
	return &VendorOrganizationService{
		opt: global.Setup.Common,
	}
}

func (s *VendorOrganizationService) Get(ctx context.Context, req *dto.VendorOrganizationListRequest) (res []*dto.VendorOrganizationResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorOrganizationService.Get")
	defer span.End()

	var vendorOrganizations *bridgeService.GetVendorOrganizationGPResponse
	vendorOrganizations, err = s.opt.Client.BridgeServiceGrpc.GetVendorOrganizationGPList(ctx, &bridgeService.GetVendorOrganizationGPListRequest{
		Limit:            req.Limit,
		Offset:           req.Offset,
		PrpVendorOrgDesc: req.Search,
		Status:           "0",
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor_organization")
		return
	}

	for _, vendorOrganization := range vendorOrganizations.Data {

		res = append(res, &dto.VendorOrganizationResponse{
			ID:          vendorOrganization.PrpVendorOrgId,
			Code:        vendorOrganization.PrpVendorOrgId,
			Description: vendorOrganization.PrpVendorOrgDesc,
		})
	}
	total = int64(len(res))

	return
}

func (s *VendorOrganizationService) GetByID(ctx context.Context, id string) (res *dto.VendorOrganizationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorOrganizationService.GetByID")
	defer span.End()

	var vendorOrganization *bridgeService.GetVendorOrganizationGPResponse
	vendorOrganization, err = s.opt.Client.BridgeServiceGrpc.GetVendorOrganizationGPDetail(ctx, &bridgeService.GetVendorOrganizationGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor_organization")
		return
	}

	res = &dto.VendorOrganizationResponse{
		ID:          vendorOrganization.Data[0].PrpVendorOrgId,
		Code:        vendorOrganization.Data[0].PrpVendorOrgId,
		Description: vendorOrganization.Data[0].PrpVendorOrgDesc,
	}

	return
}
