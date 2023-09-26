package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"

	"git.edenfarm.id/edenlabs/edenlabs/opt"

	"git.edenfarm.id/edenlabs/edenlabs/statusx"

	"time"

	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/repository"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

func ServiceUser() IUserService {
	m := new(UserService)
	m.opt = global.Setup.Common
	return m
}

type IUserService interface {
	Get(ctx context.Context, req *dto.GetUserRequest) (res []*dto.UserResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.UserResponse, err error)
	GetList(ctx context.Context, req *dto.GetUserRequest) (res []*dto.UserResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, email string, employee_code string) (res dto.UserResponse, err error)
	GetByEmail(ctx context.Context, email string) (res dto.UserResponse, err error)
	Create(ctx context.Context, req dto.UserRequestCreate) (res dto.UserResponse, err error)
	Update(ctx context.Context, req dto.UserRequestUpdate, id int64) (res dto.UserResponse, err error)
	UpdateProfile(ctx context.Context, req dto.ProfileRequestUpdate, id int64) (res dto.ProfileResponse, err error)
	Archive(ctx context.Context, id int64) (res dto.UserResponse, err error)
	UnArchive(ctx context.Context, id int64) (res dto.UserResponse, err error)
	ResetPassword(ctx context.Context, req dto.UserRequestResetPassword, id int64) (res dto.UserResponse, err error)
	UpdateSalesAppToken(ctx context.Context, req dto.UpdateSalesAppTokenRequest) (res dto.ProfileResponse, err error)
	GetBySalesAppLoginToken(ctx context.Context, token string) (res dto.UserResponse, err error)
	UpdateEdnAppToken(ctx context.Context, req dto.UpdateEdnAppTokenRequest) (res dto.ProfileResponse, err error)
	UpdatePurchaserAppToken(ctx context.Context, req dto.UpdatePurchaserAppTokenRequest) (res dto.ProfileResponse, err error)
	GetByEdnAppLoginToken(ctx context.Context, token string) (res dto.UserResponse, err error)
}

type UserService struct {
	opt opt.Options
}

func (s *UserService) Get(ctx context.Context, req *dto.GetUserRequest) (res []*dto.UserResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.Get")
	defer span.End()

	rUserRole := repository.RepositoryUserRole()
	rRole := repository.RepositoryRole()
	rUser := repository.RepositoryUser()
	rDivision := repository.RepositoryDivision()
	var users []*model.User

	users, _, err = rUser.Get(ctx, req)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, user := range users {
		// get user role
		var userRoles []*model.UserRole
		userRoles, err = rUserRole.GetByUserID(ctx, user.ID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		var mainRoleResponse *dto.RoleResponse
		var subRoleResponses []*dto.RoleResponse
		var skipped bool

		for _, userRole := range userRoles {
			// get role
			var role *model.Role
			role, err = rRole.GetByID(ctx, userRole.RoleID)
			if err != nil {
				err = edenlabs.ErrorNotFound("role")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			// get division
			var division *model.Division
			division, err = rDivision.GetDetail(ctx, role.DivisionID, "")
			if err != nil {
				err = edenlabs.ErrorNotFound("division")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			// check division id
			if req.DivisionID != 0 && division.ID != req.DivisionID {
				skipped = true
			}

			if userRole.MainRole == 1 {
				// check role id
				if req.RoleID != 0 && role.ID != req.RoleID {
					skipped = true
				}

				mainRoleResponse = &dto.RoleResponse{
					ID:   role.ID,
					Code: role.Code,
					Name: role.Name,
					Division: &dto.DivisionResponse{
						ID:            division.ID,
						Code:          division.Code,
						Name:          division.Name,
						Status:        division.Status,
						StatusConvert: statusx.ConvertStatusValue(division.Status),
						CreatedAt:     division.CreatedAt,
						UpdatedAt:     division.UpdatedAt,
					},
					Status:        role.Status,
					StatusConvert: statusx.ConvertStatusValue(role.Status),
					CreatedAt:     role.CreatedAt,
					UpdatedAt:     role.UpdatedAt,
				}

			} else {
				subRoleResponses = append(subRoleResponses, &dto.RoleResponse{
					ID:   role.ID,
					Code: role.Code,
					Name: role.Name,
					Division: &dto.DivisionResponse{
						ID:            division.ID,
						Code:          division.Code,
						Name:          division.Name,
						Status:        division.Status,
						StatusConvert: statusx.ConvertStatusValue(division.Status),
						CreatedAt:     division.CreatedAt,
						UpdatedAt:     division.UpdatedAt,
					},
					Status:        role.Status,
					StatusConvert: statusx.ConvertStatusValue(role.Status),
					CreatedAt:     role.CreatedAt,
					UpdatedAt:     role.UpdatedAt,
				})
			}
		}

		if !skipped {

			var site *bridgeService.GetSiteGPResponse
			site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
				Id: user.SiteIDGP,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "site")
				return
			}

			var admDivision *bridgeService.GetAdmDivisionGPResponse
			admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
				Id:   user.RegionIDGP,
				Type: "region",
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "adm_division")
				return
			}

			userResponse := &dto.UserResponse{
				ID:            user.ID,
				Name:          user.Name,
				Nickname:      user.Nickname,
				Email:         user.Email,
				EmployeeCode:  user.EmployeeCode,
				CreatedAt:     user.CreatedAt,
				UpdatedAt:     user.UpdatedAt,
				Status:        user.Status,
				StatusConvert: statusx.ConvertStatusValue(user.Status),
				MainRole:      mainRoleResponse,
				SubRoles:      subRoleResponses,
				Note:          user.Note,
			}
			if len(admDivision.Data) > 0 {
				userResponse.Region = &dto.RegionResponse{
					ID:          admDivision.Data[0].Code,
					Description: admDivision.Data[0].Region,
				}
			}

			if len(site.Data) > 0 {
				userResponse.Site = &dto.SiteResponse{
					ID:          site.Data[0].Locncode,
					Code:        site.Data[0].Locncode,
					Description: site.Data[0].Locndscr,
				}
			}

			if user.TerritoryIDGP != "" {
				var territory *bridgeService.GetSalesTerritoryGPResponse
				territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
					Id: user.TerritoryIDGP,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorRpcNotFound("bridge", "territory")
					return
				}
				if len(territory.Data) > 0 {
					userResponse.Territory = &dto.TerritoryResponse{
						ID:          territory.Data[0].Salsterr,
						Code:        territory.Data[0].Salsterr,
						Description: territory.Data[0].Slterdsc,
					}
				}
			}

			res = append(res, userResponse)
			total += 1
		}
	}

	return
}

func (s *UserService) GetByID(ctx context.Context, id int64) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.GetByID")
	defer span.End()

	rUserRole := repository.RepositoryUserRole()
	rRole := repository.RepositoryRole()
	rUser := repository.RepositoryUser()
	rDivision := repository.RepositoryDivision()
	var user *model.User
	var supervisor *model.User
	user, err = rUser.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if user.ParentID != 0 {
		supervisor, err = rUser.GetByID(ctx, user.ParentID)
	}
	// get user role
	var userRoles []*model.UserRole
	userRoles, err = rUserRole.GetByUserID(ctx, user.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var mainRoleResponse *dto.RoleResponse
	var subRoleResponses []*dto.RoleResponse

	for _, userRole := range userRoles {
		// get role
		var role *model.Role
		role, err = rRole.GetByID(ctx, userRole.RoleID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// get division
		var division *model.Division
		division, err = rDivision.GetDetail(ctx, role.DivisionID, "")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if userRole.MainRole == 1 {
			mainRoleResponse = &dto.RoleResponse{
				ID:   role.ID,
				Code: role.Code,
				Name: role.Name,
				Division: &dto.DivisionResponse{
					ID:            division.ID,
					Code:          division.Code,
					Name:          division.Name,
					Status:        division.Status,
					StatusConvert: statusx.ConvertStatusValue(division.Status),
					CreatedAt:     division.CreatedAt,
					UpdatedAt:     division.UpdatedAt,
				},
				Status:        role.Status,
				StatusConvert: statusx.ConvertStatusValue(role.Status),
				CreatedAt:     role.CreatedAt,
				UpdatedAt:     role.UpdatedAt,
			}

		} else {
			subRoleResponses = append(subRoleResponses, &dto.RoleResponse{
				ID:   role.ID,
				Code: role.Code,
				Name: role.Name,
				Division: &dto.DivisionResponse{
					ID:            division.ID,
					Code:          division.Code,
					Name:          division.Name,
					Status:        division.Status,
					StatusConvert: statusx.ConvertStatusValue(division.Status),
					CreatedAt:     division.CreatedAt,
					UpdatedAt:     division.UpdatedAt,
				},
				Status:        role.Status,
				StatusConvert: statusx.ConvertStatusValue(role.Status),
				CreatedAt:     role.CreatedAt,
				UpdatedAt:     role.UpdatedAt,
			})
		}
	}

	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: user.SiteIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	var admDivision *bridgeService.GetAdmDivisionGPResponse
	admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
		Id:   user.RegionIDGP,
		Type: "region",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm_division")
		return
	}

	res = dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Nickname: user.Nickname,
		Email:    user.Email,
		Password: "*****",
		ParentID: user.ParentID,
		Region: &dto.RegionResponse{
			ID:          admDivision.Data[0].Code,
			Description: admDivision.Data[0].Region,
		},
		Site: &dto.SiteResponse{
			ID:          site.Data[0].Locncode,
			Code:        site.Data[0].Locncode,
			Description: site.Data[0].Locndscr,
		},
		EmployeeCode:  user.EmployeeCode,
		PhoneNumber:   user.PhoneNumber,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Status:        user.Status,
		StatusConvert: statusx.ConvertStatusValue(user.Status),
		MainRole:      mainRoleResponse,
		SubRoles:      subRoleResponses,
		Note:          user.Note,
	}

	if user.TerritoryIDGP != "" {
		var territory *bridgeService.GetSalesTerritoryGPResponse
		territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: user.TerritoryIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "territory")
			return
		}

		res.Territory = &dto.TerritoryResponse{
			ID:          territory.Data[0].Salsterr,
			Code:        territory.Data[0].Salsterr,
			Description: territory.Data[0].Slterdsc,
		}
	}

	if supervisor != nil {
		res.Supervisor = &dto.Supervisor{
			ID:           supervisor.ID,
			Name:         supervisor.Name,
			EmployeeCode: supervisor.EmployeeCode,
		}
	}

	return
}

func (s *UserService) GetList(ctx context.Context, req *dto.GetUserRequest) (res []*dto.UserResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.Get")
	defer span.End()

	rUserRole := repository.RepositoryUserRole()
	rRole := repository.RepositoryRole()
	rUser := repository.RepositoryUser()
	rDivision := repository.RepositoryDivision()
	var users []*model.User

	users, _, err = rUser.Get(ctx, req)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, user := range users {
		// get user role
		var userRoles []*model.UserRole
		userRoles, err = rUserRole.GetByUserID(ctx, user.ID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		var mainRoleResponse *dto.RoleResponse
		var subRoleResponses []*dto.RoleResponse
		var skipped bool

		for _, userRole := range userRoles {
			// get role
			var role *model.Role
			role, err = rRole.GetByID(ctx, userRole.RoleID)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			// get division
			var division *model.Division
			division, err = rDivision.GetDetail(ctx, role.DivisionID, "")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			// check division id
			if req.DivisionID != 0 && division.ID != req.DivisionID {
				skipped = true
			}

			if userRole.MainRole == 1 {
				// check role id
				if req.RoleID != 0 && role.ID != req.RoleID {
					skipped = true
				}

				mainRoleResponse = &dto.RoleResponse{
					ID:   role.ID,
					Code: role.Code,
					Name: role.Name,
					Division: &dto.DivisionResponse{
						ID:            division.ID,
						Code:          division.Code,
						Name:          division.Name,
						Status:        division.Status,
						StatusConvert: statusx.ConvertStatusValue(division.Status),
						CreatedAt:     division.CreatedAt,
						UpdatedAt:     division.UpdatedAt,
					},
					Status:        role.Status,
					StatusConvert: statusx.ConvertStatusValue(role.Status),
					CreatedAt:     role.CreatedAt,
					UpdatedAt:     role.UpdatedAt,
				}

			} else {
				subRoleResponses = append(subRoleResponses, &dto.RoleResponse{
					ID:   role.ID,
					Code: role.Code,
					Name: role.Name,
					Division: &dto.DivisionResponse{
						ID:            division.ID,
						Code:          division.Code,
						Name:          division.Name,
						Status:        division.Status,
						StatusConvert: statusx.ConvertStatusValue(division.Status),
						CreatedAt:     division.CreatedAt,
						UpdatedAt:     division.UpdatedAt,
					},
					Status:        role.Status,
					StatusConvert: statusx.ConvertStatusValue(role.Status),
					CreatedAt:     role.CreatedAt,
					UpdatedAt:     role.UpdatedAt,
				})
			}
		}

		if !skipped {
			res = append(res, &dto.UserResponse{
				ID:            user.ID,
				Name:          user.Name,
				Nickname:      user.Nickname,
				Email:         user.Email,
				EmployeeCode:  user.EmployeeCode,
				CreatedAt:     user.CreatedAt,
				UpdatedAt:     user.UpdatedAt,
				Status:        user.Status,
				StatusConvert: statusx.ConvertStatusValue(user.Status),
				MainRole:      mainRoleResponse,
				SubRoles:      subRoleResponses,
				Note:          user.Note,
				Region: &dto.RegionResponse{
					ID: user.RegionIDGP,
				},
				Site: &dto.SiteResponse{
					ID: user.SiteIDGP,
				},
				Territory: &dto.TerritoryResponse{
					ID: user.TerritoryIDGP,
				},
			})
			total += 1
		}
	}

	return
}
func (s *UserService) GetDetail(ctx context.Context, id int64, email string, employee_code string) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.GetDetail")
	defer span.End()

	rUser := repository.RepositoryUser()
	rUserRole := repository.RepositoryUserRole()
	rRole := repository.RepositoryRole()
	rDivision := repository.RepositoryDivision()
	var (
		user *model.User
	)

	if id != 0 {
		user, err = rUser.GetByID(ctx, id)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else if email != "" {
		user, err = rUser.GetByEmail(ctx, email)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else if employee_code != "" {
		user, err = rUser.GetByEmployeeCode(ctx, employee_code)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// get user role
	var userRoles []*model.UserRole
	userRoles, err = rUserRole.GetByUserID(ctx, user.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var mainRoleResponse *dto.RoleResponse
	var subRoleResponses []*dto.RoleResponse

	for _, userRole := range userRoles {
		// get role
		var role *model.Role
		role, err = rRole.GetByID(ctx, userRole.RoleID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// get division
		var division *model.Division
		division, err = rDivision.GetDetail(ctx, role.DivisionID, "")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if userRole.MainRole == 1 {
			mainRoleResponse = &dto.RoleResponse{
				ID:   role.ID,
				Code: role.Code,
				Name: role.Name,
				Division: &dto.DivisionResponse{
					ID:            division.ID,
					Code:          division.Code,
					Name:          division.Name,
					Status:        division.Status,
					StatusConvert: statusx.ConvertStatusValue(division.Status),
					CreatedAt:     division.CreatedAt,
					UpdatedAt:     division.UpdatedAt,
				},
				Status:        role.Status,
				StatusConvert: statusx.ConvertStatusValue(role.Status),
				CreatedAt:     role.CreatedAt,
				UpdatedAt:     role.UpdatedAt,
			}
		} else {
			subRoleResponses = append(subRoleResponses, &dto.RoleResponse{
				ID:   role.ID,
				Code: role.Code,
				Name: role.Name,
				Division: &dto.DivisionResponse{
					ID:            division.ID,
					Code:          division.Code,
					Name:          division.Name,
					Status:        division.Status,
					StatusConvert: statusx.ConvertStatusValue(division.Status),
					CreatedAt:     division.CreatedAt,
					UpdatedAt:     division.UpdatedAt,
				},
				Status:        role.Status,
				StatusConvert: statusx.ConvertStatusValue(role.Status),
				CreatedAt:     role.CreatedAt,
				UpdatedAt:     role.UpdatedAt,
			})
		}
	}

	res = dto.UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		Nickname:     user.Nickname,
		Email:        user.Email,
		ParentID:     user.ParentID,
		EmployeeCode: user.EmployeeCode,
		PhoneNumber:  user.PhoneNumber,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Status:       user.Status,
		Note:         user.Note,
		MainRole:     mainRoleResponse,
		SubRoles:     subRoleResponses,
		Region: &dto.RegionResponse{
			ID: user.RegionIDGP,
		},
		Site: &dto.SiteResponse{
			ID: user.SiteIDGP,
		},
		Territory: &dto.TerritoryResponse{
			ID: user.TerritoryIDGP,
		},
		PurchaserAppNotifToken: user.PurchaserAppNotifToken,
	}

	return
}

func (s *UserService) GetByEmail(ctx context.Context, email string) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.GetByEmail")
	defer span.End()

	rUser := repository.RepositoryUser()
	var user *model.User
	user, err = rUser.GetByEmail(ctx, email)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: user.SiteIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	var region *bridgeService.GetAdmDivisionGPResponse
	region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
		Id:   user.RegionIDGP,
		Type: "region",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	var territory *bridgeService.GetSalesTerritoryGPResponse
	territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
		Id: user.TerritoryIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "territory")
		return
	}

	res = dto.UserResponse{
		ID:       user.ID,
		Name:     user.Name,
		Nickname: user.Nickname,
		Email:    user.Email,
		Password: "*****",
		ParentID: user.ParentID,
		Region: &dto.RegionResponse{
			ID:          region.Data[0].Code,
			Code:        region.Data[0].Code,
			Description: region.Data[0].Region,
		},
		EmployeeCode:  user.EmployeeCode,
		PhoneNumber:   user.PhoneNumber,
		Status:        user.Status,
		StatusConvert: statusx.ConvertStatusValue(user.Status),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Note:          user.Note,
	}
	if len(site.Data) > 0 {
		res.Site = &dto.SiteResponse{
			ID:          site.Data[0].Locncode,
			Code:        site.Data[0].Locncode,
			Description: site.Data[0].Locndscr,
		}
	}
	if len(territory.Data) > 0 {
		res.Territory = &dto.TerritoryResponse{
			ID:          territory.Data[0].Salsterr,
			Code:        territory.Data[0].Salsterr,
			Description: territory.Data[0].Slterdsc,
		}
	}
	return
}

func (s *UserService) Create(ctx context.Context, req dto.UserRequestCreate) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.Create")
	defer span.End()

	rUserRole := repository.RepositoryUserRole()
	rRole := repository.RepositoryRole()
	rUser := repository.RepositoryUser()

	// validate parent_id
	if req.ParentID != 0 {
		_, err = rUser.GetByID(ctx, req.ParentID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// validate employee code
	var existsUser *model.User
	existsUser, _ = rUser.GetByEmployeeCode(ctx, req.EmployeeCode)
	if existsUser.ID != 0 {
		err = edenlabs.ErrorValidation("employee_code", "The employee code is already exists")
		return
	}

	// validate main role
	var mainRole *model.Role
	mainRole, _ = rRole.GetByID(ctx, req.MainRole)
	if mainRole == nil {
		err = edenlabs.ErrorValidation("main_role", "The main roles is invalid")
		return
	}

	// validate sub roles
	for _, subRoleID := range req.SubRoles {
		var subRole *model.Role
		subRole, _ = rRole.GetByID(ctx, subRoleID)
		if subRole == nil {
			err = edenlabs.ErrorValidation("sub_role", "The sub roles is invalid")
			return
		}
		if subRole.ID == mainRole.ID {
			err = edenlabs.ErrorValidation("sub_role", "The sub roles is same as main role")
			return
		}
		if subRole.DivisionID != mainRole.DivisionID {
			err = edenlabs.ErrorValidation("sub_role", "The sub roles division is different as main role division")
			return
		}
	}

	passwordHash, err := utils.PasswordHasher(req.Password)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	user := &model.User{
		Name:          req.Name,
		Nickname:      req.Nickname,
		Email:         req.Email,
		Password:      passwordHash,
		RegionIDGP:    req.RegionID,
		ParentID:      req.ParentID,
		SiteIDGP:      req.SiteID,
		TerritoryIDGP: req.TerritoryID,
		EmployeeCode:  req.EmployeeCode,
		PhoneNumber:   req.PhoneNumber,
		Note:          req.Note,
		CreatedAt:     time.Now(),
		Status:        1,
	}

	span.AddEvent("creating user")
	err = rUser.Create(ctx, user)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// create main role
	span.AddEvent("creating main role")
	err = rUserRole.Create(ctx, &model.UserRole{
		UserID:    user.ID,
		RoleID:    mainRole.ID,
		CreatedAt: time.Now(),
		MainRole:  1,
		Status:    1,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// create sub roles
	span.AddEvent("creating sub roles")
	for _, subRoleID := range req.SubRoles {
		err = rUserRole.Create(ctx, &model.UserRole{
			UserID:    user.ID,
			RoleID:    subRoleID,
			CreatedAt: time.Now(),
			MainRole:  0,
			Status:    1,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	res = dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Status:    user.Status,
		Note:      user.Note,
	}

	return
}

func (s *UserService) Update(ctx context.Context, req dto.UserRequestUpdate, id int64) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.Update")
	defer span.End()

	rUserRole := repository.RepositoryUserRole()
	rRole := repository.RepositoryRole()
	rUser := repository.RepositoryUser()
	user := &model.User{
		ID:            id,
		Name:          req.Name,
		Nickname:      req.Nickname,
		RegionIDGP:    req.RegionID,
		ParentID:      req.ParentID,
		SiteIDGP:      req.SiteID,
		TerritoryIDGP: req.TerritoryID,
		PhoneNumber:   req.PhoneNumber,
		Note:          req.Note,
		UpdatedAt:     time.Now(),
	}

	// validate data is exist
	var userOld *model.User
	userOld, err = rUser.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if userOld.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	// validate main role
	var mainRole *model.Role
	mainRole, _ = rRole.GetByID(ctx, req.MainRole)
	if mainRole == nil {
		err = edenlabs.ErrorValidation("main_role", "The main roles is invalid")
		return
	}

	// validate sub roles
	for _, subRoleID := range req.SubRoles {
		var subRole *model.Role
		subRole, _ = rRole.GetByID(ctx, subRoleID)
		if subRole == nil {
			err = edenlabs.ErrorValidation("sub_role", "The sub roles is invalid")
			return
		}
		if subRole.ID == mainRole.ID {
			err = edenlabs.ErrorValidation("sub_role", "The sub roles is same as main role")
			return
		}
		if subRole.DivisionID != mainRole.DivisionID {
			err = edenlabs.ErrorValidation("sub_role", "The sub roles division is different as main role division")
			return
		}
	}

	// validate division same as old user
	var userRoleOlds []*model.UserRole
	userRoleOlds, err = rUserRole.GetActiveByUserID(ctx, userOld.ID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, userRolesOld := range userRoleOlds {
		var roleOld *model.Role
		roleOld, err = rRole.GetByID(ctx, userRolesOld.RoleID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if userRolesOld.MainRole == 1 && mainRole.DivisionID != roleOld.DivisionID {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("role", "The role is not same as old division")
			return
		}

		for _, subRoleID := range req.SubRoles {
			var subRoleNew *model.Role
			subRoleNew, err = rRole.GetByID(ctx, subRoleID)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			if subRoleNew.DivisionID != roleOld.DivisionID {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("sub_role", "The sub roles division is different as old main role division")
				return
			}
		}
	}

	span.AddEvent("updating user")
	err = rUser.Update(ctx, user, "Name", "Nickname", "RegionIDGP", "ParentID", "SiteIDGP", "TerritoryIDGP", "PhoneNumber", "Note", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// delete old main role & sub roles
	for _, userRoleOld := range userRoleOlds {
		err = rUserRole.Delete(ctx, userRoleOld.ID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// create new sub roles
	for _, roleID := range req.SubRoles {
		err = rUserRole.Create(ctx, &model.UserRole{
			UserID:    user.ID,
			RoleID:    roleID,
			MainRole:  0,
			Status:    1,
			CreatedAt: time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// create new main role
	err = rUserRole.Create(ctx, &model.UserRole{
		UserID:    user.ID,
		RoleID:    mainRole.ID,
		MainRole:  1,
		Status:    1,
		CreatedAt: time.Now(),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Nickname:  user.Nickname,
		Email:     user.Email,
		Password:  "*****",
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Status:    user.Status,
		Note:      user.Note,
	}

	return
}

func (s *UserService) UpdateProfile(ctx context.Context, req dto.ProfileRequestUpdate, id int64) (res dto.ProfileResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.Update")
	defer span.End()

	rUser := repository.RepositoryUser()
	user := &model.User{
		ID:          id,
		Nickname:    req.Nickname,
		PhoneNumber: req.PhoneNumber,
		UpdatedAt:   time.Now(),
	}

	// validate data is exist
	var userOld *model.User
	userOld, err = rUser.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if userOld.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	err = rUser.Update(ctx, user, "Nickname", "PhoneNumber", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ProfileResponse{
		ID:           user.ID,
		Name:         user.Name,
		Nickname:     user.Nickname,
		Email:        user.Email,
		Password:     "*****",
		RegionID:     user.RegionID,
		ParentID:     user.ParentID,
		SiteID:       user.SiteID,
		TerritoryID:  user.TerritoryID,
		EmployeeCode: user.EmployeeCode,
		PhoneNumber:  user.PhoneNumber,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Status:       user.Status,
	}

	return
}

func (s *UserService) Archive(ctx context.Context, id int64) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.Archive")
	defer span.End()

	rUserRole := repository.RepositoryUserRole()
	rUser := repository.RepositoryUser()
	// validate data is exist
	var userOld *model.User
	userOld, err = rUser.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if userOld.Status == 2 {
		err = edenlabs.ErrorValidation("status", "The status has been archived")
		return
	}

	if userOld.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	user := &model.User{
		ID:        id,
		Status:    2,
		UpdatedAt: time.Now(),
	}

	err = rUser.Update(ctx, user, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = rUserRole.ArchiveByUserID(ctx, user.ID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Status:    user.Status,
		Note:      user.Note,
	}

	return
}

func (s *UserService) UnArchive(ctx context.Context, id int64) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.UnArchive")
	defer span.End()

	rUser := repository.RepositoryUser()
	// validate data is exist
	var userOld *model.User
	userOld, err = rUser.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if userOld.Status != 2 {
		err = edenlabs.ErrorValidation("status", "The status is not archived")
		return
	}

	user := &model.User{
		ID:        id,
		Status:    1,
		UpdatedAt: time.Now(),
	}

	err = rUser.Update(ctx, user, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Status:    user.Status,
		Note:      user.Note,
	}

	return
}

func (s *UserService) ResetPassword(ctx context.Context, req dto.UserRequestResetPassword, id int64) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.ResetPassword")
	defer span.End()

	rUser := repository.RepositoryUser()
	// validate data is exist
	var userOld *model.User
	userOld, err = rUser.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if userOld.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	passwordHash, err := utils.PasswordHasher(req.Password)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	user := &model.User{
		ID:        id,
		Password:  passwordHash,
		UpdatedAt: time.Now(),
	}

	err = rUser.Update(ctx, user, "Password", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Status:    user.Status,
	}

	return
}

func (s *UserService) UpdateSalesAppToken(ctx context.Context, req dto.UpdateSalesAppTokenRequest) (res dto.ProfileResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.Update")
	defer span.End()

	// validate data is exist
	rUser := repository.RepositoryUser()
	var u *model.User
	u, err = rUser.GetByID(ctx, req.Id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if u.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	user := &model.User{
		ID:                 req.Id,
		ForceLogout:        req.ForceLogout,
		SalesAppLoginToken: req.SalesAppLoginToken,
		SalesAppNotifToken: req.SalesAppNotifToken,
		UpdatedAt:          time.Now(),
	}

	err = rUser.Update(ctx, user, "ForceLogout", "SalesAppLoginToken", "SalesAppNotifToken", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ProfileResponse{
		ID:           user.ID,
		Name:         user.Name,
		Nickname:     user.Nickname,
		Email:        user.Email,
		Password:     "*****",
		RegionID:     user.RegionID,
		ParentID:     user.ParentID,
		SiteID:       user.SiteID,
		TerritoryID:  user.TerritoryID,
		EmployeeCode: user.EmployeeCode,
		PhoneNumber:  user.PhoneNumber,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Status:       user.Status,
	}

	return
}

func (s *UserService) GetBySalesAppLoginToken(ctx context.Context, token string) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.GetBySalesAppLoginToken")
	defer span.End()

	rUser := repository.RepositoryUser()
	var user *model.User
	user, err = rUser.GetBySalesAppLoginToken(ctx, token)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Nickname:      user.Nickname,
		Email:         user.Email,
		Password:      "*****",
		ParentID:      user.ParentID,
		EmployeeCode:  user.EmployeeCode,
		PhoneNumber:   user.PhoneNumber,
		Status:        user.Status,
		StatusConvert: statusx.ConvertStatusValue(user.Status),
		ForceLogout:   int8(user.ForceLogout),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Note:          user.Note,
	}

	return
}

func (s *UserService) UpdateEdnAppToken(ctx context.Context, req dto.UpdateEdnAppTokenRequest) (res dto.ProfileResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.Update")
	defer span.End()

	// validate data is exist
	rUser := repository.RepositoryUser()
	var u *model.User
	u, err = rUser.GetByID(ctx, req.Id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if u.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	user := &model.User{
		ID:               req.Id,
		ForceLogout:      req.ForceLogout,
		EdnAppLoginToken: req.EdnAppLoginToken,
		UpdatedAt:        time.Now(),
	}

	err = rUser.Update(ctx, user, "ForceLogout", "EdnAppLoginToken", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ProfileResponse{
		ID:           user.ID,
		Name:         user.Name,
		Nickname:     user.Nickname,
		Email:        user.Email,
		Password:     "*****",
		RegionID:     user.RegionID,
		ParentID:     user.ParentID,
		SiteID:       user.SiteID,
		TerritoryID:  user.TerritoryID,
		EmployeeCode: user.EmployeeCode,
		PhoneNumber:  user.PhoneNumber,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Status:       user.Status,
	}

	return
}

func (s *UserService) UpdatePurchaserAppToken(ctx context.Context, req dto.UpdatePurchaserAppTokenRequest) (res dto.ProfileResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.Update")
	defer span.End()

	// validate data is exist
	rUser := repository.RepositoryUser()
	var u *model.User
	u, err = rUser.GetByID(ctx, req.Id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if u.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	user := &model.User{
		ID:                     req.Id,
		ForceLogout:            req.ForceLogout,
		PurchaserAppNotifToken: req.PurchaserAppNotifToken,
		UpdatedAt:              time.Now(),
	}

	err = rUser.Update(ctx, user, "ForceLogout", "purchaser_notif_token", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.ProfileResponse{
		ID:           user.ID,
		Name:         user.Name,
		Nickname:     user.Nickname,
		Email:        user.Email,
		Password:     "*****",
		RegionID:     user.RegionID,
		ParentID:     user.ParentID,
		SiteID:       user.SiteID,
		TerritoryID:  user.TerritoryID,
		EmployeeCode: user.EmployeeCode,
		PhoneNumber:  user.PhoneNumber,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Status:       user.Status,
	}

	return
}

func (s *UserService) GetByEdnAppLoginToken(ctx context.Context, token string) (res dto.UserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UserService.GetByEdnAppLoginToken")
	defer span.End()

	rUser := repository.RepositoryUser()
	var user *model.User
	user, err = rUser.GetByEdnAppLoginToken(ctx, token)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.UserResponse{
		ID:            user.ID,
		Name:          user.Name,
		Nickname:      user.Nickname,
		Email:         user.Email,
		Password:      "*****",
		ParentID:      user.ParentID,
		EmployeeCode:  user.EmployeeCode,
		PhoneNumber:   user.PhoneNumber,
		Status:        user.Status,
		StatusConvert: statusx.ConvertStatusValue(user.Status),
		ForceLogout:   int8(user.ForceLogout),
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		Note:          user.Note,
	}

	return
}
