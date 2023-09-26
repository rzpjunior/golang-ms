package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func NewServiceVendor() IVendorService {
	m := new(VendorService)
	m.opt = global.Setup.Common
	return m
}

type IVendorService interface {
	GetVendors(ctx context.Context, req dto.VendorListRequest) (res []*dto.VendorResponse, err error)
	GetVendorDetailById(ctx context.Context, req dto.VendorDetailRequest) (res *dto.VendorResponse, err error)
	GetListGp(ctx context.Context, req dto.GetVendorGPListRequest) (res []*dto.VendorGP, total int64, err error)
	GetDetailGp(ctx context.Context, id string) (res *dto.VendorGP, err error)
}

type VendorService struct {
	opt opt.Options
}

func (s *VendorService) GetVendors(ctx context.Context, req dto.VendorListRequest) (res []*dto.VendorResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.GetVendors")
	defer span.End()

	// get vendor from bridge
	var vendorRes *bridgeService.GetVendorListResponse
	vendorRes, err = s.opt.Client.BridgeServiceGrpc.GetVendorList(ctx, &bridgeService.GetVendorListRequest{
		Limit:   req.Limit,
		Offset:  req.Offset,
		Status:  req.Status,
		Search:  req.Search,
		OrderBy: req.OrderBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	datas := []*dto.VendorResponse{}
	for _, vendor := range vendorRes.Data {
		datas = append(datas, &dto.VendorResponse{
			// ID:                     vendor.Id,
			Code: vendor.Code,
			// VendorOrganizationID:   vendor.VendorOrganizationId,
			VendorClassificationID: vendor.VendorClassificationId,
			SubDistrictID:          vendor.SubDistrictId,
			PaymentTermID:          vendor.PaymentTermId,
			PicName:                vendor.PicName,
			Email:                  vendor.Email,
			PhoneNumber:            vendor.PhoneNumber,
			Rejectable:             vendor.Rejectable,
			Returnable:             vendor.Rejectable,
			Address:                vendor.Address,
			Note:                   vendor.Note,
			Status:                 vendor.Status,
			Latitude:               vendor.Latitude,
			Longitude:              vendor.Longitude,
			CreatedAt:              vendor.CreatedAt.AsTime(),
		})
	}
	res = datas

	return
}

func (s *VendorService) GetVendorDetailById(ctx context.Context, req dto.VendorDetailRequest) (res *dto.VendorResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.GetVendorDetailById")
	defer span.End()

	// get Vendor from bridge
	var vendorRes *bridgeService.GetVendorDetailResponse
	vendorRes, err = s.opt.Client.BridgeServiceGrpc.GetVendorDetail(ctx, &bridgeService.GetVendorDetailRequest{
		Id: int64(req.Id),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
		return
	}

	res = &dto.VendorResponse{
		// ID:                     vendorRes.Data.Id,
		Code: vendorRes.Data.Code,
		// VendorOrganizationID:   vendorRes.Data.VendorOrganizationId,
		VendorClassificationID: vendorRes.Data.VendorClassificationId,
		SubDistrictID:          vendorRes.Data.SubDistrictId,
		PaymentTermID:          vendorRes.Data.PaymentTermId,
		PicName:                vendorRes.Data.PicName,
		Email:                  vendorRes.Data.Email,
		PhoneNumber:            vendorRes.Data.PhoneNumber,
		Rejectable:             vendorRes.Data.Rejectable,
		Returnable:             vendorRes.Data.Rejectable,
		Address:                vendorRes.Data.Address,
		Note:                   vendorRes.Data.Note,
		Status:                 vendorRes.Data.Status,
		Latitude:               vendorRes.Data.Latitude,
		Longitude:              vendorRes.Data.Longitude,
		CreatedAt:              vendorRes.Data.CreatedAt.AsTime(),
	}

	return
}

func (s *VendorService) GetListGp(ctx context.Context, req dto.GetVendorGPListRequest) (res []*dto.VendorGP, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.GetVendorGp")
	defer span.End()

	var vendor *bridgeService.GetVendorGPResponse

	if vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPList(ctx, &bridgeService.GetVendorGPListRequest{
		Limit:  int32(req.Limit),
		Offset: int32(req.Offset),
		Search: req.Search,
	}); err != nil || !vendor.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
		return
	}

	total = int64(len(vendor.Data))
	for _, v := range vendor.Data {

		res = append(res, &dto.VendorGP{
			VendorId: v.VENDORID,
			VendName: v.VENDNAME,
			Address:  v.ADDRESS1 + v.ADDRESS2 + v.ADDRESS3,
			Inactive: int32(v.VENDSTTS),
		})

	}

	return
}

func (s *VendorService) GetDetailGp(ctx context.Context, id string) (res *dto.VendorGP, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.GetDetailVendorGp")
	defer span.End()

	var vendor *bridgeService.GetVendorGPResponse

	if vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
		Id: id,
	}); err != nil || !vendor.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
		return
	}

	if len(vendor.Data) > 0 {
		res = &dto.VendorGP{
			VendorId: vendor.Data[0].VENDORID,
			VendName: vendor.Data[0].VENDNAME,
			Address:  vendor.Data[0].ADDRESS1 + vendor.Data[0].ADDRESS2 + vendor.Data[0].ADDRESS3,
			Inactive: int32(vendor.Data[0].VENDSTTS),
		}
	}

	return
}
