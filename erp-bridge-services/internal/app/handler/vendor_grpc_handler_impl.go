package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetVendorList(ctx context.Context, req *bridgeService.GetVendorListRequest) (res *bridgeService.GetVendorListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorList")
	defer span.End()

	var vendors []dto.VendorResponse
	vendors, _, err = h.ServicesVendor.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Vendor
	for _, vendor := range vendors {
		data = append(data, &bridgeService.Vendor{
			Id:                     vendor.ID,
			Code:                   vendor.Code,
			VendorOrganizationId:   vendor.VendorOrganizationID,
			VendorClassificationId: vendor.VendorClassificationID,
			SubDistrictId:          vendor.SubDistrictID,
			PicName:                vendor.PicName,
			Email:                  vendor.Email,
			PhoneNumber:            vendor.PhoneNumber,
			PaymentTermId:          vendor.PaymentTermID,
			Rejectable:             vendor.Rejectable,
			Returnable:             vendor.Returnable,
			Address:                vendor.Address,
			Note:                   vendor.Note,
			Status:                 vendor.Status,
			Latitude:               vendor.Latitude,
			Longitude:              vendor.Longitude,
			CreatedAt:              timestamppb.New(vendor.CreatedAt),
			CreatedBy:              vendor.CreatedBy,
		})
	}

	res = &bridgeService.GetVendorListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetVendorDetail(ctx context.Context, req *bridgeService.GetVendorDetailRequest) (res *bridgeService.GetVendorDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorDetail")
	defer span.End()

	var vendor dto.VendorResponse
	vendor, err = h.ServicesVendor.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetVendorDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Vendor{
			Id:                     vendor.ID,
			Code:                   vendor.Code,
			VendorOrganizationId:   vendor.VendorOrganizationID,
			VendorClassificationId: vendor.VendorClassificationID,
			SubDistrictId:          vendor.SubDistrictID,
			PicName:                vendor.PicName,
			Email:                  vendor.Email,
			PhoneNumber:            vendor.PhoneNumber,
			PaymentTermId:          vendor.PaymentTermID,
			Rejectable:             vendor.Rejectable,
			Returnable:             vendor.Returnable,
			Address:                vendor.Address,
			Note:                   vendor.Note,
			Status:                 vendor.Status,
			Latitude:               vendor.Latitude,
			Longitude:              vendor.Longitude,
			CreatedAt:              timestamppb.New(vendor.CreatedAt),
			CreatedBy:              vendor.CreatedBy,
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetVendorGPList(ctx context.Context, req *bridgeService.GetVendorGPListRequest) (res *bridgeService.GetVendorGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorList")
	defer span.End()

	res, err = h.ServicesVendor.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetVendorGPDetail(ctx context.Context, req *bridgeService.GetVendorGPDetailRequest) (res *bridgeService.GetVendorGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorDetail")
	defer span.End()

	// read detail vendor from GP.
	res, err = h.ServicesVendor.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) CreateVendor(ctx context.Context, req *bridgeService.CreateVendorRequest) (res *bridgeService.CreateVendorResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateVendor")
	defer span.End()

	_, err = h.ServicesVendor.CreateGP(ctx, &dto.CreateVendorGPRequest{
		InterID:             req.Interid,
		VendorID:            req.Vendorid,
		VendName:            req.Vendname,
		VendShnm:            req.Vendshnm,
		VndChknm:            req.Vndchknm,
		VendStts:            req.Vendstts,
		PrP_Vendor_Org_ID:   req.PrP_Vendor_Org_ID,
		PrP_Vendor_CLASF_ID: req.PrP_Vendor_CLASF_ID,
		VndCntct:            req.Vndcntct,
		AddresS1:            req.AddresS1,
		AddresS2:            req.AddresS2,
		Phnumbr1:            req.PhnumbR1,
		Phnumbr2:            req.PhnumbR2,
		Vaddcdpr:            req.Vaddcdpr,
		Vadcdpad:            req.Vadcdpad,
		Vadcdsfr:            req.Vadcdsfr,
		Vadcdtro:            req.Vadcdtro,
		Prp_payment_method:  req.PRP_Payment_Method,
		Prp_payment_term:    req.Pymtrmid,
	})

	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreateVendorResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}
