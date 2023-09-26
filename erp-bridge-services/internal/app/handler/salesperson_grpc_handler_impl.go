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

func (h *BridgeGrpcHandler) GetSalespersonList(ctx context.Context, req *bridgeService.GetSalespersonListRequest) (res *bridgeService.GetSalespersonListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalespersonList")
	defer span.End()

	var salespersons []dto.SalespersonResponse
	salespersons, _, err = h.ServicesSalesperson.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Salesperson
	for _, salesperson := range salespersons {
		data = append(data, &bridgeService.Salesperson{
			Id:   salesperson.ID,
			Code: salesperson.Code,

			Firstname:  salesperson.FirstName,
			Middlename: salesperson.MiddleName,
			Lastname:   salesperson.LastName,
			Status:     int32(salesperson.Status),
			CreatedAt:  timestamppb.New(salesperson.CreatedAt),
			UpdatedAt:  timestamppb.New(salesperson.UpdatedAt),
		})
	}

	res = &bridgeService.GetSalespersonListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetSalespersonDetail(ctx context.Context, req *bridgeService.GetSalespersonDetailRequest) (res *bridgeService.GetSalespersonDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalespersonDetail")
	defer span.End()

	var salesperson dto.SalespersonResponse
	salesperson, err = h.ServicesSalesperson.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetSalespersonDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Salesperson{
			Id:   salesperson.ID,
			Code: salesperson.Code,

			Firstname:  salesperson.FirstName,
			Middlename: salesperson.MiddleName,
			Lastname:   salesperson.LastName,
			Status:     int32(salesperson.Status),
			CreatedAt:  timestamppb.New(salesperson.CreatedAt),
			UpdatedAt:  timestamppb.New(salesperson.UpdatedAt),
		},
	}
	return
}
