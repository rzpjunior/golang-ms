package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *CampaignGrpcHandler) GetCustomerPointExpirationDetail(ctx context.Context, req *pb.GetCustomerPointExpirationDetailRequest) (res *pb.GetCustomerPointExpirationDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointExpirationDetail")
	defer span.End()

	var customerPointExpiration *dto.CustomerPointExpirationResponse
	customerPointExpiration, err = h.ServiceCustomerPointExpiration.GetDetail(ctx, req.Id, req.CustomerId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.CustomerPointExpiration

	data = &pb.CustomerPointExpiration{
		Id:                 customerPointExpiration.ID,
		CustomerId:         customerPointExpiration.CustomerID,
		CurrentPeriodPoint: customerPointExpiration.CurrentPeriodPoint,
		NextPeriodPoint:    customerPointExpiration.NextPeriodPoint,
		CurrentPeriodDate:  timestamppb.New(customerPointExpiration.CurrentPeriodDate),
		NextPeriodDate:     timestamppb.New(customerPointExpiration.NextPeriodDate),
		LastUpdatedAt:      timestamppb.New(customerPointExpiration.LastUpdatedAt),
	}

	res = &pb.GetCustomerPointExpirationDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

// func (h *CampaignGrpcHandler) CreateCustomerPointExpiration(ctx context.Context, req *pb.CreateCustomerPointExpirationRequest) (res *pb.CreateCustomerPointExpirationResponse, err error) {
// 	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointExpirationDetail")
// 	defer span.End()

// 	param := &dto.CustomerPointExpirationRequestCreate{
// 		CustomerID:     req.CustomerId,
// 		EarnedPoint:    req.EarnedPoint,
// 		RedeemedPoint:  req.RedeemedPoint,
// 		ExpirationDate: req.ExpirationDate,
// 	}

// 	err = h.ServiceCustomerPointExpiration.Create(ctx, param)
// 	if err != nil {
// 		err = status.New(codes.NotFound, err.Error()).Err()
// 		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
// 		return
// 	}

// 	res = &pb.CreateCustomerPointExpirationResponse{
// 		Code:    int32(codes.OK),
// 		Message: codes.OK.String(),
// 	}

// 	return
// }

// func (h *CampaignGrpcHandler) UpdateCustomerPointExpiration(ctx context.Context, req *pb.UpdateCustomerPointExpirationRequest) (res *pb.UpdateCustomerPointExpirationResponse, err error) {
// 	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointExpirationDetail")
// 	defer span.End()

// 	param := &dto.CustomerPointExpirationRequestUpdate{
// 		ID:            req.Id,
// 		CustomerID:    req.CustomerId,
// 		EarnedPoint:   req.EarnedPoint,
// 		RedeemedPoint: req.RedeemedPoint,
// 		FieldUpdate:   req.FieldUpdate,
// 	}

// 	err = h.ServiceCustomerPointExpiration.Update(ctx, param)
// 	if err != nil {
// 		err = status.New(codes.NotFound, err.Error()).Err()
// 		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
// 		return
// 	}

// 	res = &pb.UpdateCustomerPointExpirationResponse{
// 		Code:    int32(codes.OK),
// 		Message: codes.OK.String(),
// 	}

// 	return
// }
