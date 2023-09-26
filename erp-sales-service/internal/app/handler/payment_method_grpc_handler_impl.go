package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *SalesGrpcHandler) GetPaymentMethodList(ctx context.Context, req *pb.GetPaymentMethodListRequest) (res *pb.GetPaymentMethodListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetPaymentMethodList")
	defer span.End()
	var pm []*model.PaymentMethod
	pm, err = h.ServicePaymentMethod.GetListGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.PaymentMethod
	for _, paymentMethod := range pm {
		data = append(data, &pb.PaymentMethod{
			Id:          paymentMethod.ID,
			Code:        paymentMethod.Code,
			Name:        paymentMethod.Name,
			Note:        paymentMethod.Note,
			Status:      int32(paymentMethod.Status),
			Publish:     int32(paymentMethod.Publish),
			Maintenance: int32(paymentMethod.Maintenance),
		})
	}

	res = &pb.GetPaymentMethodListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
