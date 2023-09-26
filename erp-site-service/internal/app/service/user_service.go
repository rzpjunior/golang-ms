package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
)

type IUserService interface {
	GetUser(ctx context.Context, req dto.GetUserRequest) (res []dto.UserResponse, total int64, err error)
}

type UserService struct {
	opt opt.Options
}

func NewServiceUser() IUserService {
	return &UserService{
		opt: global.Setup.Common,
	}
}

func (s *UserService) GetUser(ctx context.Context, req dto.GetUserRequest) (res []dto.UserResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.GetUser")
	defer span.End()

	var user *accountService.GetUserListResponse

	if user, err = s.opt.Client.AccountServiceGrpc.GetUserList(ctx, &accountService.GetUserListRequest{
		Limit:      int32(req.Limit),
		Offset:     int32(req.Offset),
		Status:     int32(req.Status),
		Search:     req.Search,
		OrderBy:    req.OrderBy,
		SiteId:     req.SiteId,
		DivisionId: req.DivisionId,
		RoleId:     req.RoleId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user")
		return
	}

	for _, user := range user.Data {
		res = append(res, dto.UserResponse{
			ID:           user.Id,
			Name:         user.Name,
			Nickname:     user.Nickname,
			Email:        user.Email,
			RegionID:     user.RegionId,
			ParentID:     user.ParentId,
			SiteID:       user.SiteId,
			TerritoryID:  user.TerritoryId,
			EmployeeCode: user.EmployeeCode,
			PhoneNumber:  user.PhoneNumber,
			Status:       int8(user.Status),
			CreatedAt:    user.CreatedAt.AsTime(),
			UpdatedAt:    user.UpdatedAt.AsTime(),
		})
	}

	total = int64(len(user.Data))

	return
}
