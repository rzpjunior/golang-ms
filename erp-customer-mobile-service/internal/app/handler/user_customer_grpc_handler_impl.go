package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/customer_mobile_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CustomerMobileGrpcHandler) GetUserCustomerDetail(ctx context.Context, req *pb.GetUserCustomerDetailRequest) (res *pb.GetUserCustomerDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetItemSectionDetail")
	defer span.End()

	param := &dto.GetDetailUserCustomerRequest{
		CustomerID: req.CustomerId,
	}
	var userCustomer *dto.UserCustomerResponse
	userCustomer, err = h.ServiceUserCustomer.GetDetail(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.GetUserCustomerDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.UserCustomer{
			Id:            userCustomer.ID,
			CustomerId:    userCustomer.CustomerID,
			FirebaseToken: userCustomer.FirebaseToken,
		},
	}
	return
}
