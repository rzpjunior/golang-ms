package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetVendorOrganizationList(ctx context.Context, req *bridgeService.GetVendorOrganizationListRequest) (res *bridgeService.GetVendorOrganizationListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorOrganizationList")
	defer span.End()

	var vendorOrganizations []dto.VendorOrganizationResponse
	vendorOrganizations, _, err = h.ServicesVendorOrganization.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.VendorOrganization
	for _, vendorOrganization := range vendorOrganizations {
		data = append(data, &bridgeService.VendorOrganization{
			Id:                     vendorOrganization.ID,
			Code:                   vendorOrganization.Code,
			VendorClassificationId: vendorOrganization.VendorClassificationID,
			SubDistrictId:          vendorOrganization.SubDistrictID,
			PaymentTermId:          vendorOrganization.PaymentTermID,
			Name:                   vendorOrganization.Name,
			Address:                vendorOrganization.Address,
			Note:                   vendorOrganization.Note,
			Status:                 vendorOrganization.Status,
		})
	}

	res = &bridgeService.GetVendorOrganizationListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetVendorOrganizationDetail(ctx context.Context, req *bridgeService.GetVendorOrganizationDetailRequest) (res *bridgeService.GetVendorOrganizationDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorOrganizationDetail")
	defer span.End()

	var vendorOrganization dto.VendorOrganizationResponse
	vendorOrganization, err = h.ServicesVendorOrganization.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetVendorOrganizationDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.VendorOrganization{
			Id:                     vendorOrganization.ID,
			Code:                   vendorOrganization.Code,
			VendorClassificationId: vendorOrganization.VendorClassificationID,
			SubDistrictId:          vendorOrganization.SubDistrictID,
			PaymentTermId:          vendorOrganization.PaymentTermID,
			Name:                   vendorOrganization.Name,
			Address:                vendorOrganization.Address,
			Note:                   vendorOrganization.Note,
			Status:                 vendorOrganization.Status,
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetVendorOrganizationGPList(ctx context.Context, req *bridgeService.GetVendorOrganizationGPListRequest) (res *bridgeService.GetVendorOrganizationGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorOrganizationGPList")
	defer span.End()

	res, err = h.ServicesVendorOrganization.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetVendorOrganizationGPDetail(ctx context.Context, req *bridgeService.GetVendorOrganizationGPDetailRequest) (res *bridgeService.GetVendorOrganizationGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorOrganizationGPDetail")
	defer span.End()

	res, err = h.ServicesVendorOrganization.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
