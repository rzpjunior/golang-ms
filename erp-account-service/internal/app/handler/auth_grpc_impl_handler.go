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

func (h *AccountGrpcHandler) GetUserEmailAuth(ctx context.Context, req *accountService.GetUserEmailAuthRequest) (res *accountService.GetUserEmailAuthResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetUserRoleByUserId")
	defer span.End()

	sAuth := service.ServiceAuth()

	var userRolesRes dto.UserPasswordResponse
	userRolesRes, err = sAuth.GetUserEmail(ctx, dto.UserPasswordRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &accountService.GetUserEmailAuthResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &accountService.User{
			Id:           userRolesRes.ID,
			Name:         userRolesRes.Name,
			Nickname:     userRolesRes.Nickname,
			Email:        userRolesRes.Email,
			Password:     userRolesRes.Password,
			SiteIdGp:     userRolesRes.Site.Code,
			EmployeeCode: userRolesRes.EmployeeCode,
			PhoneNumber:  userRolesRes.PhoneNumber,
			Status:       int32(userRolesRes.Status),
			CreatedAt:    timestamppb.New(userRolesRes.CreatedAt),
			UpdatedAt:    timestamppb.New(userRolesRes.UpdatedAt),
			RegionIdGp:   userRolesRes.Region.Code,
		},
	}
	return
}
