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

func (h *AccountGrpcHandler) GetRoleDetail(ctx context.Context, req *accountService.GetRoleDetailRequest) (res *accountService.GetRoleDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetRoleDetail")
	defer span.End()

	sRole := service.ServiceRole()

	var role dto.RoleResponse
	role, err = sRole.GetByID(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &accountService.Role{
		Id:         role.ID,
		Code:       role.Code,
		Name:       role.Name,
		DivisionId: role.Division.ID,
		CreatedAt:  timestamppb.New(role.CreatedAt),
		UpdatedAt:  timestamppb.New(role.UpdatedAt),
		Status:     int32(role.Status),
		Note:       role.Note,
	}

	res = &accountService.GetRoleDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
