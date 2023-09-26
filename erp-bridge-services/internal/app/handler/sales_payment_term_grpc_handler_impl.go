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

func (h *BridgeGrpcHandler) GetSalesPaymentTermList(ctx context.Context, req *bridgeService.GetSalesPaymentTermListRequest) (res *bridgeService.GetSalesPaymentTermListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesPaymentTermList")
	defer span.End()

	var salesPaymentTerms []dto.SalesPaymentTermResponse
	salesPaymentTerms, _, err = h.ServicesSalesPaymentTerm.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.SalesPaymentTerm
	for _, salesPaymentTerm := range salesPaymentTerms {
		data = append(data, &bridgeService.SalesPaymentTerm{
			Id:          salesPaymentTerm.ID,
			Code:        salesPaymentTerm.Code,
			Description: salesPaymentTerm.Description,
			Status:      int32(salesPaymentTerm.Status),
			CreatedAt:   timestamppb.New(salesPaymentTerm.CreatedAt),
			UpdatedAt:   timestamppb.New(salesPaymentTerm.UpdatedAt),
		})
	}

	res = &bridgeService.GetSalesPaymentTermListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetSalesPaymentTermDetail(ctx context.Context, req *bridgeService.GetSalesPaymentTermDetailRequest) (res *bridgeService.GetSalesPaymentTermDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesPaymentTermDetail")
	defer span.End()

	var salesPaymentTerm dto.SalesPaymentTermResponse
	salesPaymentTerm, err = h.ServicesSalesPaymentTerm.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetSalesPaymentTermDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.SalesPaymentTerm{
			Id:          salesPaymentTerm.ID,
			Code:        salesPaymentTerm.Code,
			Description: salesPaymentTerm.Description,
			Status:      int32(salesPaymentTerm.Status),
			CreatedAt:   timestamppb.New(salesPaymentTerm.CreatedAt),
			UpdatedAt:   timestamppb.New(salesPaymentTerm.UpdatedAt),
		},
	}
	return
}
