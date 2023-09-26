package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *SalesGrpcHandler) GetPaymentChannelList(ctx context.Context, req *pb.GetPaymentChannelListRequest) (res *pb.GetPaymentChannelListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "SalesGrpcHandler.GetPaymentChannelList")
	defer span.End()
	var pm []*model.PaymentChannel
	pm, err = h.ServicePaymentChannel.GetListGRPC(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.PaymentChannel
	for _, paymentChannel := range pm {
		data = append(data, &pb.PaymentChannel{
			Id:              paymentChannel.ID,
			Code:            paymentChannel.Code,
			Name:            paymentChannel.Name,
			Note:            paymentChannel.Note,
			Status:          int32(paymentChannel.Status),
			Value:           paymentChannel.Value,
			ImageUrl:        paymentChannel.ImageUrl,
			PublishIva:      int32(paymentChannel.PublishIva),
			PublishFva:      int32(paymentChannel.PublishFva),
			PaymentMethodId: paymentChannel.PaymentMethodID,
			PaymentGuideUrl: paymentChannel.PaymentGuideURL,
		})
	}

	res = &pb.GetPaymentChannelListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
