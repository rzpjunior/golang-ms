package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	dto "git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/dto"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

func NewServiceAuth() IAuthService {
	m := new(AuthService)
	m.opt = global.Setup.Common
	return m
}

type IAuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error)
	Session(c context.Context) (res dto.UserResponse, err error)
}

type AuthService struct {
	opt opt.Options
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AuthService.Login")
	defer span.End()

	var (
		loggedIn bool
	)
	// user role
	var userResponse *accountService.GetUserEmailAuthResponse
	userResponse, err = s.opt.Client.AccountServiceGrpc.GetUserEmailAuth(ctx, &accountService.GetUserEmailAuthRequest{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user")
		return
	}

	// check active user
	if userResponse.Data.Status != 1 {
		err = edenlabs.ErrorMustActive("email")
		return
	}

	// user role
	var userRoles *accountService.GetUserRoleByUserIdResponse
	userRoles, err = s.opt.Client.AccountServiceGrpc.GetUserRolesByUserId(ctx, &accountService.GetUserRoleByUserIdRequest{
		Id: userResponse.Data.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user_role")
		return
	}

	// config app
	var configAppEdnApp *configurationService.GetConfigAppListResponse
	configAppEdnApp, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configurationService.GetConfigAppListRequest{
		Application: 8,
		Attribute:   "admin_edn_role_id",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "config_app")
		return
	}

	for _, i := range configAppEdnApp.Data {
		salesAppRolesStr := strings.Split(i.Value, ",")
		for _, salesAppRole := range salesAppRolesStr {
			var salesAppRoleInt int
			salesAppRoleInt, err = strconv.Atoi(salesAppRole)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			for _, userRole := range userRoles.Data {
				if int(userRole.RoleId) == salesAppRoleInt {
					loggedIn = true
				}
			}
		}
	}

	if !loggedIn {
		err = edenlabs.ErrorInvalid("role")
		return
	}

	// validate password
	if err = utils.PasswordHash(userResponse.Data.Password, req.Password); err != nil {
		err = edenlabs.ErrorValidation("password", "The password is not match")
		return
	}

	jwtInit := jwt.NewJWT([]byte(s.opt.Config.Jwt.Key))
	uc := jwt.UserClaim{
		UserID:      userResponse.Data.Id,
		Permissions: []string{"edn_app"},
		ExpiresAt:   time.Now().Add(time.Hour * 18).Unix(),
		Timezone:    req.Timezone,
	}

	jwtGenerate, err := jwtInit.Create(uc)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// update token
	_, err = s.opt.Client.AccountServiceGrpc.UpdateUserEdnAppToken(ctx, &accountService.UpdateUserEdnAppTokenRequest{
		Id:               userResponse.Data.Id,
		ForceLogout:      2,
		EdnappLoginToken: jwtGenerate,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "UpdateUserSalesAppToken")
		return
	}

	res.Token = jwtGenerate
	res.User = &dto.UserResponse{
		ID:           userResponse.Data.Id,
		Name:         userResponse.Data.Name,
		Nickname:     userResponse.Data.Nickname,
		Email:        userResponse.Data.Email,
		ParentID:     userResponse.Data.ParentId,
		SiteID:       userResponse.Data.SiteIdGp,
		EmployeeCode: userResponse.Data.EmployeeCode,
		PhoneNumber:  userResponse.Data.PhoneNumber,
		Status:       int8(userResponse.Data.Status),
		CreatedAt:    userResponse.Data.CreatedAt.AsTime(),
		UpdatedAt:    userResponse.Data.UpdatedAt.AsTime(),
	}
	return
}

func (s *AuthService) Session(c context.Context) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(c, "AuthService.Session")
	defer span.End()

	var userResponse *accountService.GetUserDetailResponse

	userID := ctx.Value(constants.KeyUserID).(int64)

	// user detail
	userResponse, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		Id: userID,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user")

		return
	}

	// check active user
	if userResponse.Data.Status != 1 {
		err = edenlabs.ErrorMustActive("email")
		return
	}

	// get site
	// var site *bridgeService.GetSiteGPResponse
	// site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
	// 	Id: userResponse.Data.SiteIdGp,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "site")
	// 	return
	// }

	// get region
	// var region *bridgeService.GetAdmDivisionGPResponse
	// region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
	// 	Region: userResponse.Data.RegionIdGp,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "region")
	// 	return
	// }

	// fmt.Println(site)
	res = dto.UserResponse{
		ID:           userResponse.Data.Id,
		Name:         userResponse.Data.Name,
		Nickname:     userResponse.Data.Nickname,
		Email:        userResponse.Data.Email,
		ParentID:     userResponse.Data.ParentId,
		EmployeeCode: userResponse.Data.EmployeeCode,
		PhoneNumber:  userResponse.Data.PhoneNumber,
		Status:       int8(userResponse.Data.Status),
		CreatedAt:    userResponse.Data.CreatedAt.AsTime(),
		UpdatedAt:    userResponse.Data.UpdatedAt.AsTime(),
		SiteID:       userResponse.Data.SiteIdGp,
		Region: &dto.RegionResponse{
			// ID:          region.Data[0].Region,
			Code:        userResponse.Data.RegionIdGp,
			Description: userResponse.Data.RegionIdGp,
		},
	}

	return
}
