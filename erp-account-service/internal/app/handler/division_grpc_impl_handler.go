package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/service"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *AccountGrpcHandler) GetDivisionDetail(ctx context.Context, req *accountService.GetDivisionDetailRequest) (res *accountService.GetDivisionDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetUserRoleByUserId")
	defer span.End()

	division := service.ServiceDivision()

	var divisionResponse dto.DivisionResponse
	divisionResponse, err = division.GetDetail(ctx, req.Id, req.Code)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &accountService.GetDivisionDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &accountService.Division{
			Id:        divisionResponse.ID,
			Code:      divisionResponse.Code,
			Name:      divisionResponse.Name,
			CreatedAt: timestamppb.New(divisionResponse.CreatedAt),
			UpdatedAt: timestamppb.New(divisionResponse.UpdatedAt),
			Status:    int32(divisionResponse.Status),
		},
	}
	return
}

func (h *AccountGrpcHandler) GetDivisionDefaultByCustomerType(ctx context.Context, req *accountService.GetDivisionDefaultByCustomerTypeRequest) (res *accountService.GetDivisionDefaultByCustomerTypeResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDivisionDefaultByCustomerType")
	defer span.End()

	division := service.ServiceDivision()

	var divisionResponse *dto.DivisionResponse
	divisionResponse, err = division.GetDivisonByCustomerType(ctx, req.CustomerTypeIdGp)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &accountService.GetDivisionDefaultByCustomerTypeResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &accountService.Division{
			Id:        divisionResponse.ID,
			Code:      divisionResponse.Code,
			Name:      divisionResponse.Name,
			CreatedAt: timestamppb.New(divisionResponse.CreatedAt),
			UpdatedAt: timestamppb.New(divisionResponse.UpdatedAt),
			Status:    int32(divisionResponse.Status),
		},
	}
	return
}
