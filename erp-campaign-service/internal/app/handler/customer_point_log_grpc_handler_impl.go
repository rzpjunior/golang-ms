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

func (h *CampaignGrpcHandler) GetCustomerPointLogList(ctx context.Context, req *pb.GetCustomerPointLogListRequest) (res *pb.GetCustomerPointLogListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointLogList")
	defer span.End()

	param := &dto.CustomerPointLogRequestGet{
		Limit:           int64(req.Limit),
		Offset:          int64(req.Offset),
		CustomerID:      req.CustomerId,
		SalesOrderID:    req.SalesOrderId,
		TransactionType: int8(req.TransactionType),
		Status:          int8(req.Status),
		OrderBy:         req.OrderBy,
		CreatedDate:     req.CreatedDate,
	}

	var customerPointLogs []*dto.CustomerPointLogResponse
	customerPointLogs, _, err = h.ServiceCustomerPointLog.Get(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.CustomerPointLog
	for _, v := range customerPointLogs {
		data = append(data, &pb.CustomerPointLog{
			Id:               v.ID,
			CustomerId:       v.CustomerID,
			SalesOrderId:     v.SalesOrderID,
			EpCampaignId:     v.EPCampaignID,
			CurrentPointUsed: v.CurrentPointUsed,
			NextPointUsed:    v.NextPointUsed,
			PointValue:       v.PointValue,
			RecentPoint:      v.RecentPoint,
			Status:           int32(v.Status),
			TransactionType:  int32(v.TransactionType),
			CreatedDate:      v.CreatedDate,
			ExpiredDate:      v.ExpiredDate,
		})
	}

	res = &pb.GetCustomerPointLogListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetCustomerPointLogDetail(ctx context.Context, req *pb.GetCustomerPointLogDetailRequest) (res *pb.GetCustomerPointLogDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointLogDetail")
	defer span.End()

	param := &dto.CustomerPointLogRequestGetDetail{
		CustomerID:      req.CustomerId,
		SalesOrderID:    req.SalesOrderId,
		Status:          int8(req.Status),
		CreatedDate:     req.CreatedDate,
		TransactionType: int8(req.TransactionType),
	}

	var customerPointLog *dto.CustomerPointLogResponse
	customerPointLog, err = h.ServiceCustomerPointLog.GetDetail(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.CustomerPointLog

	data = &pb.CustomerPointLog{
		Id:               customerPointLog.ID,
		CustomerId:       customerPointLog.CustomerID,
		SalesOrderId:     customerPointLog.SalesOrderID,
		EpCampaignId:     customerPointLog.EPCampaignID,
		CurrentPointUsed: customerPointLog.CurrentPointUsed,
		NextPointUsed:    customerPointLog.NextPointUsed,
		PointValue:       customerPointLog.PointValue,
		RecentPoint:      customerPointLog.RecentPoint,
		Status:           int32(customerPointLog.Status),
		TransactionType:  int32(customerPointLog.TransactionType),
		CreatedDate:      customerPointLog.CreatedDate,
		ExpiredDate:      customerPointLog.ExpiredDate,
	}

	res = &pb.GetCustomerPointLogDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) GetCustomerPointLogDetailHistoryMobile(ctx context.Context, req *pb.GetCustomerPointLogDetailRequest) (res *pb.GetCustomerPointLogDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointLogDetail")
	defer span.End()

	param := &dto.CustomerPointLogRequestGetDetail{
		CustomerID:      req.CustomerId,
		SalesOrderID:    req.SalesOrderId,
		Status:          int8(req.Status),
		CreatedDate:     req.CreatedDate,
		TransactionType: int8(req.TransactionType),
	}

	var customerPointLog *dto.PointHistoryList
	customerPointLog, err = h.ServiceCustomerPointLog.GetDetailHistoryMobile(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.CustomerPointLog

	data = &pb.CustomerPointLog{
		PointValue:  customerPointLog.PointValue,
		Status:      int32(customerPointLog.Status),
		CreatedDate: customerPointLog.CreatedDate,
		StatusType:  customerPointLog.StatusType,
	}

	res = &pb.GetCustomerPointLogDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CampaignGrpcHandler) CreateCustomerPointLog(ctx context.Context, req *pb.CreateCustomerPointLogRequest) (res *pb.CreateCustomerPointLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerPointLogDetail")
	defer span.End()

	param := &dto.CustomerPointLogRequestCreate{
		CustomerID:       req.CustomerId,
		SalesOrderID:     req.SalesOrderId,
		EPCampaignID:     req.EpCampaignId,
		PointValue:       req.PointValue,
		RecentPoint:      req.RecentPoint,
		CurrentPointUsed: req.CurrentPointUsed,
		NextPointUsed:    req.NextPointUsed,
		Status:           int8(req.Status),
		TransactionType:  int8(req.TransactionType),
		Note:             req.Note,
		CreatedDate:      req.CreatedDate,
		ExpiredDate:      req.ExpiredDate,
	}
	var customerPointLogID int64
	customerPointLogID, err = h.ServiceCustomerPointLog.Create(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.CreateCustomerPointLogResponse{
		Code:               int32(codes.OK),
		Message:            codes.OK.String(),
		CustomerPointLogId: customerPointLogID,
	}

	return
}

func (h *CampaignGrpcHandler) GetReferralHistory(ctx context.Context, req *pb.GetReferralHistoryRequest) (res *pb.GetReferralHistoryResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetReferralHistory")
	defer span.End()
	var referralHistory *dto.ReferralHistoryReturn
	referralHistory, err = h.ServiceCustomerPointLog.GetReferralDataCustomer(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var dataReferral []*pb.ReferralList
	var dataReferralPoint []*pb.ReferralPointList
	for _, v := range referralHistory.Detail.ReferralList {
		dataReferral = append(dataReferral, &pb.ReferralList{
			Name:      v.Name,
			CreatedAt: timestamppb.New(v.CreatedAt),
		})
	}
	for _, v := range referralHistory.Detail.ReferralPointList {
		dataReferralPoint = append(dataReferralPoint, &pb.ReferralPointList{
			Name:       v.Name,
			CreatedAt:  timestamppb.New(v.CreatedAt),
			PointValue: v.PointValue,
		})
	}
	res = &pb.GetReferralHistoryResponse{
		Code:              int32(codes.OK),
		Message:           codes.OK.String(),
		TotalPoint:        referralHistory.Summary.TotalPoint,
		TotalReferral:     referralHistory.Summary.TotalReferral,
		DataReferral:      dataReferral,
		DataReferralPoint: dataReferralPoint,
	}

	return
}

func (h *CampaignGrpcHandler) CancelCustomerPointLog(ctx context.Context, req *pb.CancelCustomerPointLogRequest) (res *pb.CancelCustomerPointLogResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetReferralHistory")
	defer span.End()

	param := &dto.CustomerPointLogRequestCancel{
		CustomerID:   req.CustomerId,
		SalesOrderID: req.SalesOrderId,
	}

	_, err = h.ServiceCustomerPointLog.Cancel(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.CancelCustomerPointLogResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}

	return
}
