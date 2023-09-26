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

func (h *AccountGrpcHandler) GetUserRoleByUserId(ctx context.Context, req *accountService.GetUserRoleByUserIdRequest) (res *accountService.GetUserRoleByUserIdResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetUserRoleByUserId")
	defer span.End()

	sUserRole := service.ServiceUserRole()

	var userRolesRes dto.UserRoleByUserIdResponse
	userRolesRes, err = sUserRole.GetByUserID(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	datas := []*accountService.UserRole{}

	for _, userRole := range userRolesRes.Roles {
		datas = append(datas, &accountService.UserRole{
			Id:        userRole.ID,
			UserId:    userRole.UserID,
			RoleId:    userRole.RoleID,
			MainRole:  int32(userRole.MainRole),
			CreatedAt: timestamppb.New(userRole.CreatedAt),
			UpdatedAt: timestamppb.New(userRole.UpdatedAt),
			Status:    int32(userRole.Status),
		})
	}

	res = &accountService.GetUserRoleByUserIdResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    datas,
	}
	return
}
