package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetCourierVendorList
func (h *BridgeGrpcHandler) GetCourierVendorList(ctx context.Context, req *bridgeService.GetCourierVendorListRequest) (res *bridgeService.GetCourierVendorListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCourierVendorList")
	defer span.End()

	var courierVendors []dto.CourierVendorResponse
	courierVendors, _, err = h.ServicesCourierVendor.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.SiteId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.CourierVendor
	for _, cVendor := range courierVendors {
		data = append(data, &bridgeService.CourierVendor{
			Id:     cVendor.ID,
			SiteId: cVendor.SiteID,
			Code:   cVendor.Code,
			Name:   cVendor.Name,
			Note:   cVendor.Note,
			Status: int32(cVendor.Status),
		})
	}

	res = &bridgeService.GetCourierVendorListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *BridgeGrpcHandler) GetCourierVendorDetail(ctx context.Context, req *bridgeService.GetCourierVendorDetailRequest) (res *bridgeService.GetCourierVendorDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCourierVendorDetail")
	defer span.End()

	var courierVendor dto.CourierVendorResponse
	courierVendor, err = h.ServicesCourierVendor.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetCourierVendorDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.CourierVendor{
			Id:     courierVendor.ID,
			SiteId: courierVendor.SiteID,
			Code:   courierVendor.Code,
			Name:   courierVendor.Name,
			Note:   courierVendor.Note,
			Status: int32(courierVendor.Status),
		},
	}

	return
}

func (h *BridgeGrpcHandler) GetCourierVendorGPList(ctx context.Context, req *bridgeService.GetCourierVendorGPListRequest) (res *bridgeService.GetCourierVendorGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCourierVendorGPList")
	defer span.End()

	res, err = h.ServicesCourierVendor.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetCourierVendorGPDetail(ctx context.Context, req *bridgeService.GetCourierVendorGPDetailRequest) (res *bridgeService.GetCourierVendorGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCourierVendorGPDetail")
	defer span.End()

	res, err = h.ServicesCourierVendor.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
