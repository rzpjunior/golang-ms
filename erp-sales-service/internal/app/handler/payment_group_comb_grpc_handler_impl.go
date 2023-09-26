package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *SalesGrpcHandler) GetPaymentGroupCombList(ctx context.Context, req *pb.GetPaymentGroupCombListRequest) (res *pb.GetPaymentGroupCombListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetPaymentGroupCombList")
	defer span.End()
	var pm []*model.PaymentGroupComb
	pm, err = h.ServicePaymentGroupComb.GetListGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.PaymentGroupComb
	for _, paymentGroupComb := range pm {
		data = append(data, &pb.PaymentGroupComb{
			Id:              paymentGroupComb.ID,
			PaymentGroupSls: paymentGroupComb.PaymentGroupSls,
			TermPaymentSls:  paymentGroupComb.TermPaymentSls,
		})
	}

	res = &pb.GetPaymentGroupCombListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
