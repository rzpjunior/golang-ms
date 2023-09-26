package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CampaignGrpcHandler) GetCustomerPointSummaryDetail(ctx context.Context, req *pb.GetCustomerPointSummaryRequestDetail) (res *pb.GetCustomerPointSummaryDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointSummaryDetail")
	defer span.End()

	param := &dto.CustomerPointSummaryRequestGetDetail{
		ID:          req.Id,
		CustomerID:  req.CustomerId,
		SummaryDate: req.SummaryDate,
	}

	var customerPointLog *dto.CustomerPointSummaryResponse
	customerPointLog, err = h.ServiceCustomerPointSummary.GetDetail(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.CustomerPointSummary

	data = &pb.CustomerPointSummary{
		Id:            customerPointLog.ID,
		CustomerId:    customerPointLog.CustomerID,
		EarnedPoint:   customerPointLog.EarnedPoint,
		RedeemedPoint: customerPointLog.RedeemedPoint,
		SummaryDate:   customerPointLog.SummaryDate,
	}

	res = &pb.GetCustomerPointSummaryDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) CreateCustomerPointSummary(ctx context.Context, req *pb.CreateCustomerPointSummaryRequest) (res *pb.CreateCustomerPointSummaryResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointSummaryDetail")
	defer span.End()

	param := &dto.CustomerPointSummaryRequestCreate{
		CustomerID:    req.CustomerId,
		EarnedPoint:   req.EarnedPoint,
		RedeemedPoint: req.RedeemedPoint,
		SummaryDate:   req.SummaryDate,
	}

	err = h.ServiceCustomerPointSummary.Create(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.CreateCustomerPointSummaryResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}

func (h *CampaignGrpcHandler) UpdateCustomerPointSummary(ctx context.Context, req *pb.UpdateCustomerPointSummaryRequest) (res *pb.UpdateCustomerPointSummaryResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointSummaryDetail")
	defer span.End()

	param := &dto.CustomerPointSummaryRequestUpdate{
		ID:            req.Id,
		CustomerID:    req.CustomerId,
		EarnedPoint:   req.EarnedPoint,
		RedeemedPoint: req.RedeemedPoint,
		FieldUpdate:   req.FieldUpdate,
	}

	err = h.ServiceCustomerPointSummary.Update(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.UpdateCustomerPointSummaryResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}
