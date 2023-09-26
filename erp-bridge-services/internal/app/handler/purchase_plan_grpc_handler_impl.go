package handler

import (
	context "context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetPurchasePlanList(ctx context.Context, req *bridgeService.GetPurchasePlanListRequest) (res *bridgeService.GetPurchasePlanListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchasePlanList")
	defer span.End()

	var purchasePlans []dto.PurchasePlanResponse
	purchasePlans, _, err = h.ServicesPurchasePlan.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.PurchasePlan
	for _, purchasePlan := range purchasePlans {
		data = append(data, &bridgeService.PurchasePlan{
			Id:                   purchasePlan.ID,
			Code:                 purchasePlan.Code,
			VendorOrganizationId: purchasePlan.VendorOrganizationID,
			SiteId:               purchasePlan.SiteID,
			RecognitionDate:      timestamppb.New(purchasePlan.RecognitionDate),
			EtaDate:              timestamppb.New(purchasePlan.EtaDate),
			EtaTime:              purchasePlan.EtaTime,
			TotalPrice:           purchasePlan.TotalPrice,
			TotalWeight:          purchasePlan.TotalWeight,
			TotalPurchasePlanQty: purchasePlan.TotalPurchasePlanQty,
			TotalPurchaseQty:     purchasePlan.TotalPurchaseQty,
			Note:                 purchasePlan.Note,
			Status:               purchasePlan.Status,
			AssignedTo:           purchasePlan.AssignedTo,
			AssignedBy:           purchasePlan.AssignedBy,
			AssignedAt:           timestamppb.New(purchasePlan.AssignedAt),
			CreatedAt:            timestamppb.New(purchasePlan.CreatedAt),
			CreatedBy:            purchasePlan.CreatedBy,
		})
	}

	res = &bridgeService.GetPurchasePlanListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetPurchasePlanDetail(ctx context.Context, req *bridgeService.GetPurchasePlanDetailRequest) (res *bridgeService.GetPurchasePlanDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchasePlanDetail")
	defer span.End()

	var purchasePlan dto.PurchasePlanResponse
	purchasePlan, err = h.ServicesPurchasePlan.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetPurchasePlanDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.PurchasePlan{
			Id:                   purchasePlan.ID,
			Code:                 purchasePlan.Code,
			VendorOrganizationId: purchasePlan.VendorOrganizationID,
			SiteId:               purchasePlan.SiteID,
			RecognitionDate:      timestamppb.New(purchasePlan.RecognitionDate),
			EtaDate:              timestamppb.New(purchasePlan.EtaDate),
			EtaTime:              purchasePlan.EtaTime,
			TotalPrice:           purchasePlan.TotalPrice,
			TotalWeight:          purchasePlan.TotalWeight,
			TotalPurchasePlanQty: purchasePlan.TotalPurchasePlanQty,
			TotalPurchaseQty:     purchasePlan.TotalPurchaseQty,
			Note:                 purchasePlan.Note,
			Status:               purchasePlan.Status,
			AssignedTo:           purchasePlan.AssignedTo,
			AssignedBy:           purchasePlan.AssignedBy,
			AssignedAt:           timestamppb.New(purchasePlan.AssignedAt),
			CreatedAt:            timestamppb.New(purchasePlan.CreatedAt),
			CreatedBy:            purchasePlan.CreatedBy,
		},
	}
	return
}

func (h *BridgeGrpcHandler) AssignPurchasePlanGP(ctx context.Context, req *bridgeService.AssignPurchasePlanGPRequest) (res *bridgeService.AssignPurchasePlanGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.AssignPurchasePlan")
	defer span.End()
	// var purchasePlan *dto.AssignPurchasePlanGPResponse

	_, err = h.ServicesPurchasePlan.AssignPurchasePlanGP(ctx, &dto.AssignPurchasePlanGPRequest{
		PrpPurchaseplanNo:   req.PrpPurchaseplanNo,
		PrpPurchaseplanUser: req.PrpPurchaseplanUser,
	})

	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	h.Option.Common.Logger.AddMessage(log.InfoLevel, "AssignPurchasePlan").Print()

	res = &bridgeService.AssignPurchasePlanGPResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}

func (h *BridgeGrpcHandler) CancelAssignPurchasePlan(ctx context.Context, req *bridgeService.CancelAssignPurchasePlanRequest) (res *bridgeService.CancelAssignPurchasePlanResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CancelAssignPurchasePlan")
	defer span.End()

	h.Option.Common.Logger.AddMessage(log.InfoLevel, "CancelAssignPurchasePlan").Print()

	res = &bridgeService.CancelAssignPurchasePlanResponse{}

	return
}

func (h *BridgeGrpcHandler) GetPurchasePlanGPList(ctx context.Context, req *bridgeService.GetPurchasePlanGPListRequest) (res *bridgeService.GetPurchasePlanGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchasePlanGPList")
	defer span.End()
	fmt.Println(ctx)

	res, err = h.ServicesPurchasePlan.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}

func (h *BridgeGrpcHandler) GetPurchasePlanGPDetail(ctx context.Context, req *bridgeService.GetPurchasePlanGPDetailRequest) (res *bridgeService.GetPurchasePlanGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetPurchasePlanGPDetail")
	defer span.End()

	res, err = h.ServicesPurchasePlan.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return
}
