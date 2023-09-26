package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetVendorClassificationList(ctx context.Context, req *bridgeService.GetVendorClassificationListRequest) (res *bridgeService.GetVendorClassificationListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorClassificationList")
	defer span.End()

	var vendorClassifications []dto.VendorClassificationResponse
	vendorClassifications, _, err = h.ServicesVendorClassification.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.VendorClassification
	for _, vendorClassification := range vendorClassifications {
		data = append(data, &bridgeService.VendorClassification{
			Id:            vendorClassification.ID,
			CommodityCode: vendorClassification.CommodityCode,
			CommodityName: vendorClassification.CommodityName,
			BadgeCode:     vendorClassification.BadgeCode,
			BadgeName:     vendorClassification.BadgeName,
			TypeCode:      vendorClassification.TypeCode,
			TypeName:      vendorClassification.TypeName,
		})
	}

	res = &bridgeService.GetVendorClassificationListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetVendorClassificationDetail(ctx context.Context, req *bridgeService.GetVendorClassificationDetailRequest) (res *bridgeService.GetVendorClassificationDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorClassificationDetail")
	defer span.End()

	var vendorClassification dto.VendorClassificationResponse
	vendorClassification, err = h.ServicesVendorClassification.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetVendorClassificationDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.VendorClassification{
			Id:            vendorClassification.ID,
			CommodityCode: vendorClassification.CommodityCode,
			CommodityName: vendorClassification.CommodityName,
			BadgeCode:     vendorClassification.BadgeCode,
			BadgeName:     vendorClassification.BadgeName,
			TypeCode:      vendorClassification.TypeCode,
			TypeName:      vendorClassification.TypeName,
		},
	}
	return
}

func (h *BridgeGrpcHandler) GetVendorClassificationGPList(ctx context.Context, req *bridgeService.GetVendorClassificationGPListRequest) (res *bridgeService.GetVendorClassificationGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorClassificationGPList")
	defer span.End()

	res, err = h.ServicesVendorClassification.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetVendorClassificationGPDetail(ctx context.Context, req *bridgeService.GetVendorClassificationGPDetailRequest) (res *bridgeService.GetVendorClassificationGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetVendorClassificationGPDetail")
	defer span.End()

	res, err = h.ServicesVendorClassification.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
