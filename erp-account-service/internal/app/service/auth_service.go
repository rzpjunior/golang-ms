package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/repository"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	dto "git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
)

func ServiceAuth() IAuthService {
	m := new(AuthService)
	m.opt = global.Setup.Common
	return m
}

type IAuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error)
	GetUserEmail(ctx context.Context, req dto.UserPasswordRequest) (res dto.UserPasswordResponse, err error)
}

type AuthService struct {
	opt opt.Options
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error) {
	var user *model.User
	rAuth := repository.RepositoryAuth()
	rUserRole := repository.RepositoryUserRole()
	rRolePermission := repository.RepositoryRolePermssion()
	rPermission := repository.RepositoryPermission()
	user, err = rAuth.GetByEmail(ctx, req.Email, req.Password)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("email", "Please recheck Email or Password")
		return
	}

	// validate password
	if err = utils.PasswordHash(user.Password, req.Password); err != nil {
		err = edenlabs.ErrorValidation("password", "The password is not match")
		return
	}

	// user role
	var userRoles []*model.UserRole
	userRoles, err = rUserRole.GetActiveByUserID(ctx, user.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var permissions []string
	for _, userRole := range userRoles {
		// role permission
		var rolePermissions []*model.RolePermission
		rolePermissions, err = rRolePermission.GetByRoleID(ctx, userRole.RoleID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		for _, rolePermission := range rolePermissions {
			// permission
			var permission *model.Permission
			permission, err = rPermission.GetByID(ctx, rolePermission.PermissionID)
			if err != nil {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			permissions = append(permissions, permission.Value)
		}
	}

	jwtInit := jwt.NewJWT([]byte(s.opt.Config.Jwt.Key))
	uc := jwt.UserClaim{
		UserID:      user.ID,
		Permissions: permissions,
		ExpiresAt:   time.Now().Add(time.Hour * 4).Unix(),
		Timezone:    req.Timezone,
	}

	jwtGenerate, err := jwtInit.Create(uc)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res.Token = jwtGenerate
	return
}

func (s *AuthService) GetUserEmail(ctx context.Context, req dto.UserPasswordRequest) (res dto.UserPasswordResponse, err error) {
	var user *model.User
	rAuth := repository.RepositoryAuth()
	user, err = rAuth.GetByEmail(ctx, req.Email, req.Password)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("email", "Please recheck Email or Password")
		return
	}

	// validate password
	if err = utils.PasswordHash(user.Password, req.Password); err != nil {
		err = edenlabs.ErrorValidation("password", "The password is not match")
		return
	}

	userRes := dto.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Nickname:      user.Nickname,
		Email:         user.Email,
		Password:      user.Password,
		ParentID:      user.ParentID,
		EmployeeCode:  user.EmployeeCode,
		PhoneNumber:   user.PhoneNumber,
		Status:        user.Status,
		StatusConvert: statusx.ConvertStatusValue(user.Status),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Site: &dto.SiteResponse{
			Code: user.SiteIDGP,
		},
		Region: &dto.RegionResponse{
			Code: user.RegionIDGP,
		},
	}
	res = dto.UserPasswordResponse{
		UserResponse: userRes,
	}

	return
}
