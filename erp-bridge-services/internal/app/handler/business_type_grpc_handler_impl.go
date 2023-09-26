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

// GetCustomerTypeList
func (h *BridgeGrpcHandler) GetCustomerTypeList(ctx context.Context, req *bridgeService.GetCustomerTypeListRequest) (res *bridgeService.GetCustomerTypeListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerTypeList")
	defer span.End()

	var CustomerTypes []dto.CustomerTypeResponse
	CustomerTypes, _, err = h.ServicesCustomerType.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.CustomerType
	for _, CustomerType := range CustomerTypes {
		data = append(data, &bridgeService.CustomerType{
			Id:           CustomerType.ID,
			Code:         CustomerType.Code,
			Description:  CustomerType.Description,
			GroupType:    CustomerType.GroupType,
			Abbreviation: CustomerType.Abbreviation,
			Status:       int32(CustomerType.Status),
			CreatedAt:    timestamppb.New(CustomerType.CreatedAt),
			UpdatedAt:    timestamppb.New(CustomerType.UpdatedAt),
		})
	}

	res = &bridgeService.GetCustomerTypeListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetCustomerTypeDetail(ctx context.Context, req *bridgeService.GetCustomerTypeDetailRequest) (res *bridgeService.GetCustomerTypeDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerTypeDetail")
	defer span.End()

	var CustomerType dto.CustomerTypeResponse
	CustomerType, err = h.ServicesCustomerType.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetCustomerTypeDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.CustomerType{
			Id:           CustomerType.ID,
			Code:         CustomerType.Code,
			Description:  CustomerType.Description,
			GroupType:    CustomerType.GroupType,
			Abbreviation: CustomerType.Abbreviation,
			Status:       int32(CustomerType.Status),
			CreatedAt:    timestamppb.New(CustomerType.CreatedAt),
			UpdatedAt:    timestamppb.New(CustomerType.UpdatedAt),
		},
	}
	return
}
