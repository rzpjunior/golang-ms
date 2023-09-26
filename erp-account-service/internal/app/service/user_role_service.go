package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"

	dto "git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/repository"
)

func ServiceUserRole() IUserRoleService {
	m := new(UserRoleService)
	m.opt = global.Setup.Common
	return m
}

type IUserRoleService interface {
	GetByUserID(ctx context.Context, id int64) (res dto.UserRoleByUserIdResponse, err error)
}

type UserRoleService struct {
	opt opt.Options
}

func (s *UserRoleService) GetByUserID(ctx context.Context, id int64) (res dto.UserRoleByUserIdResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.GetByID")
	defer span.End()

	rUserRole := repository.RepositoryUserRole()

	// get user role
	var userRoles []*model.UserRole
	userRoles, err = rUserRole.GetByUserID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var userRolesResponse []*dto.UserRoleResponse
	for _, userRole := range userRoles {
		userRolesResponse = append(userRolesResponse, &dto.UserRoleResponse{
			ID:        userRole.ID,
			UserID:    userRole.UserID,
			RoleID:    userRole.RoleID,
			MainRole:  userRole.MainRole,
			CreatedAt: userRole.CreatedAt,
			UpdatedAt: userRole.UpdatedAt,
			Status:    userRole.Status,
		})
	}

	res = dto.UserRoleByUserIdResponse{
		Roles: userRolesResponse,
	}

	return
}
