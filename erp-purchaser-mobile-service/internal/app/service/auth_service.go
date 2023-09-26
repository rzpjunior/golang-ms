package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
)

type IAuthService interface {
	Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error)
	Session(c context.Context) (res dto.UserResponse, err error)
}

type AuthService struct {
	opt opt.Options
}

func NewAuthService() IAuthService {
	return &AuthService{
		opt: global.Setup.Common,
	}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AuthService.Login")
	defer span.End()

	// user
	var userAuth *accountService.GetUserEmailAuthResponse
	userAuth, err = s.opt.Client.AccountServiceGrpc.GetUserEmailAuth(ctx, &accountService.GetUserEmailAuthRequest{
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
	if userAuth.Data.Status != 1 {
		err = edenlabs.ErrorMustActive("email")
		return
	}

	// user role
	var userRoles *accountService.GetUserRoleByUserIdResponse
	userRoles, err = s.opt.Client.AccountServiceGrpc.GetUserRolesByUserId(ctx, &accountService.GetUserRoleByUserIdRequest{
		Id: userAuth.Data.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user_role")
		return
	}

	// config app
	var configAppSalesApp *configurationService.GetConfigAppListResponse
	configAppSalesApp, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configurationService.GetConfigAppListRequest{
		Application: 7,
		Attribute:   "fieldpurchaser_role_id",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "config_app")
		return
	}

	var loggedIn bool

	for _, i := range configAppSalesApp.Data {
		purchaserAppRolesStr := strings.Split(i.Value, ",")
		for _, purchaserAppRole := range purchaserAppRolesStr {
			var purchaserAppRoleInt int
			purchaserAppRoleInt, err = strconv.Atoi(purchaserAppRole)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			for _, userRole := range userRoles.Data {
				if int(userRole.RoleId) == purchaserAppRoleInt {
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
	if err = utils.PasswordHash(userAuth.Data.Password, req.Password); err != nil {
		err = edenlabs.ErrorValidation("password", "The password is not match")
		return
	}

	var user *accountService.GetUserDetailResponse
	user, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		Id: userAuth.Data.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user")
		return
	}

	var permissions []string
	switch user.Data.MainRole {
	case "Sourcing Manager":
		permissions = append(permissions, "purchaser_app_manager")
		permissions = append(permissions, "purchaser_app")
	case "Purchasing Manager":
		permissions = append(permissions, "purchaser_app_manager")
		permissions = append(permissions, "purchaser_app")
	default:
		permissions = append(permissions, "purchaser_app")

	}

	jwtInit := jwt.NewJWT([]byte(s.opt.Config.Jwt.Key))
	uc := jwt.UserClaim{
		UserID:      userAuth.Data.Id,
		Permissions: permissions,
		ExpiresAt:   time.Now().Add(time.Hour * 18).Unix(),
		Timezone:    req.Timezone,
	}

	jwtGenerate, err := jwtInit.Create(uc)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// get site
	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: userAuth.Data.SiteIdGp,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	// get region
	// var region *bridgeService.GetRegionDetailResponse
	// region, err = s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridgeService.GetRegionDetailRequest{
	// 	Id: userAuth.Data.RegionId,
	// })

	var region *bridgeService.GetAdmDivisionGPResponse
	region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
		Region: userAuth.Data.RegionIdGp,
		Limit:  1,
		Offset: 0,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "region")
		return
	}

	res.Token = jwtGenerate
	res.User = &dto.UserResponse{
		ID:           user.Data.Id,
		Name:         user.Data.Name,
		Nickname:     user.Data.Nickname,
		Email:        user.Data.Email,
		ParentID:     user.Data.ParentId,
		EmployeeCode: user.Data.EmployeeCode,
		PhoneNumber:  user.Data.PhoneNumber,
		Division:     user.Data.Division,
		MainRole:     user.Data.MainRole,
		Status:       int8(user.Data.Status),
		CreatedAt:    user.Data.CreatedAt.AsTime(),
		UpdatedAt:    user.Data.UpdatedAt.AsTime(),
		Site: &dto.SiteResponse{
			ID:          site.Data[0].Locncode,
			Code:        site.Data[0].Locncode,
			Description: site.Data[0].Locndscr,
		},
		Region: &dto.RegionResponse{
			// ID:          region.Data[0].Id,
			Code:        region.Data[0].Region,
			Description: region.Data[0].Region,
		},
	}

	// var userUpdate *accountService.GetUserDetailResponse
	_, err = s.opt.Client.AccountServiceGrpc.UpdateUserPurchaserAppToken(ctx, &accountService.UpdateUserPurchaserAppTokenRequest{
		Id:                     userAuth.Data.Id,
		PurchaserappNotifToken: req.FCM,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user")
		return
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
	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: userResponse.Data.SiteIdGp,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	// get region
	var region *bridgeService.GetRegionDetailResponse
	region, err = s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridgeService.GetRegionDetailRequest{
		Id: userResponse.Data.RegionId,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "region")
		return
	}

	res = dto.UserResponse{
		ID:           userResponse.Data.Id,
		Name:         userResponse.Data.Name,
		Nickname:     userResponse.Data.Nickname,
		Email:        userResponse.Data.Email,
		ParentID:     userResponse.Data.ParentId,
		EmployeeCode: userResponse.Data.EmployeeCode,
		PhoneNumber:  userResponse.Data.PhoneNumber,
		Division:     userResponse.Data.Division,
		MainRole:     userResponse.Data.MainRole,
		Status:       int8(userResponse.Data.Status),
		CreatedAt:    userResponse.Data.CreatedAt.AsTime(),
		UpdatedAt:    userResponse.Data.UpdatedAt.AsTime(),
		Site: &dto.SiteResponse{
			ID:          site.Data[0].Locncode,
			Code:        site.Data[0].Locncode,
			Description: site.Data[0].Locndscr,
		},
		Region: &dto.RegionResponse{
			ID:          region.Data.Id,
			Code:        region.Data.Code,
			Description: region.Data.Description,
		},
	}

	return
}
