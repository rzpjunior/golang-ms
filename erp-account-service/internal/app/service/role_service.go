package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"

	"time"

	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	dto "git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func ServiceRole() IRoleService {
	m := new(RoleService)
	m.opt = global.Setup.Common
	return m
}

type IRoleService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, divisionID int64) (res []*dto.RoleResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.RoleResponse, err error)
	Create(ctx context.Context, req dto.RoleRequestCreate) (res dto.RoleResponse, err error)
	Update(ctx context.Context, req dto.RoleRequestUpdate, id int64) (res dto.RoleResponse, err error)
	Archive(ctx context.Context, id int64) (res dto.RoleResponse, err error)
	UnArchive(ctx context.Context, id int64) (res dto.RoleResponse, err error)
}

type RoleService struct {
	opt opt.Options
}

func (s *RoleService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, divisionID int64) (res []*dto.RoleResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RoleService.Get")

	defer span.End()

	rRole := repository.RepositoryRole()
	rDivision := repository.RepositoryDivision()
	var roles []*model.Role

	roles, total, err = rRole.Get(ctx, offset, limit, status, search, orderBy, divisionID)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, role := range roles {

		// get division
		var division *model.Division
		division, err = rDivision.GetDetail(ctx, role.DivisionID, "")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		res = append(res, &dto.RoleResponse{
			ID:   role.ID,
			Code: role.Code,
			Name: role.Name,
			Division: &dto.DivisionResponse{
				ID:            division.ID,
				Code:          division.Code,
				Name:          division.Name,
				Status:        division.Status,
				StatusConvert: statusx.ConvertStatusValue(division.Status),
			},
			CreatedAt:     role.CreatedAt,
			UpdatedAt:     role.UpdatedAt,
			Status:        role.Status,
			StatusConvert: statusx.ConvertStatusValue(role.Status),
			Note:          role.Note,
		})
	}

	return
}

func (s *RoleService) GetByID(ctx context.Context, id int64) (res dto.RoleResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RoleService.GetByID")
	defer span.End()

	rRole := repository.RepositoryRole()
	rPermission := repository.RepositoryPermission()
	rDivision := repository.RepositoryDivision()
	var role *model.Role
	role, err = rRole.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var permissions []*model.Permission
	permissions, err = rPermission.GetTreeByRoleID(ctx, role.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var permissionIDs []int64
	var permissionsResponse []*dto.PermissionResponse
	for _, permission := range permissions {
		permissionIDs = append(permissionIDs, permission.ID)
		var childPermission []*dto.PermissionResponse
		for _, child := range permission.Child {
			permissionIDs = append(permissionIDs, child.ID)
			var grandChildPermission []*dto.PermissionResponse
			for _, grandChild := range child.GrandChild {
				permissionIDs = append(permissionIDs, grandChild.ID)
				grandChildPermission = append(grandChildPermission, &dto.PermissionResponse{
					ID:            grandChild.ID,
					ParentID:      grandChild.ParentID,
					Name:          grandChild.Name,
					Value:         grandChild.Value,
					CreatedAt:     grandChild.CreatedAt,
					UpdatedAt:     grandChild.UpdatedAt,
					Status:        grandChild.Status,
					StatusConvert: statusx.ConvertStatusValue(grandChild.Status),
				})
			}

			childPermission = append(childPermission, &dto.PermissionResponse{
				ID:            child.ID,
				ParentID:      child.ParentID,
				Name:          child.Name,
				Value:         child.Value,
				CreatedAt:     child.CreatedAt,
				UpdatedAt:     child.UpdatedAt,
				Status:        child.Status,
				StatusConvert: statusx.ConvertStatusValue(child.Status),
				GrandChild:    grandChildPermission,
			})
		}

		permissionsResponse = append(permissionsResponse, &dto.PermissionResponse{
			ID:            permission.ID,
			ParentID:      permission.ParentID,
			Name:          permission.Name,
			Value:         permission.Value,
			CreatedAt:     permission.CreatedAt,
			UpdatedAt:     permission.UpdatedAt,
			Status:        permission.Status,
			StatusConvert: statusx.ConvertStatusValue(permission.Status),
			Child:         childPermission,
		})
	}

	// get division
	var division *model.Division
	division, err = rDivision.GetDetail(ctx, role.DivisionID, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.RoleResponse{
		ID:            role.ID,
		Code:          role.Code,
		Name:          role.Name,
		Permissions:   permissionsResponse,
		PermissionIDs: permissionIDs,
		Division: &dto.DivisionResponse{
			ID:            division.ID,
			Code:          division.Code,
			Name:          division.Name,
			Status:        division.Status,
			StatusConvert: statusx.ConvertStatusValue(division.Status),
			CreatedAt:     division.CreatedAt,
			UpdatedAt:     division.UpdatedAt,
		},
		CreatedAt:     role.CreatedAt,
		UpdatedAt:     role.UpdatedAt,
		Status:        role.Status,
		StatusConvert: statusx.ConvertStatusValue(role.Status),
		Note:          role.Note,
	}

	return
}

func (s *RoleService) Create(ctx context.Context, req dto.RoleRequestCreate) (res dto.RoleResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RoleService.Create")
	defer span.End()

	rPermission := repository.RepositoryPermission()
	rRole := repository.RepositoryRole()
	rRolePermission := repository.RepositoryRolePermssion()
	rDivision := repository.RepositoryDivision()

	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "ROL",
		Domain: "role",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "generate_code")
		return
	}

	role := &model.Role{
		Code:       codeGenerator.Data.Code,
		DivisionID: req.DivisionID,
		Name:       req.Name,
		CreatedAt:  time.Now(),
		Note:       req.Note,
		Status:     1,
	}

	// validate role name
	var roleExisting *model.Role
	span.AddEvent("validate role name")
	roleExisting, _ = rRole.GetByName(ctx, req.Name)
	if roleExisting.ID != 0 {
		err = edenlabs.ErrorValidation("role", "The role already exists")
		return
	}

	// validate division id
	span.AddEvent("validate division id")
	_, err = rDivision.GetDetail(ctx, req.DivisionID, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("division is validated", trace.WithAttributes(attribute.Int64("division_id", req.DivisionID)))

	// validate permission id
	span.AddEvent("validate permission_id")
	var permissions []*dto.PermissionResponse
	for _, permission_id := range req.Permissions {
		var permission *model.Permission
		permission, err = rPermission.GetByID(ctx, permission_id)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		permissions = append(permissions, &dto.PermissionResponse{
			ID:        permission.ID,
			Name:      permission.Name,
			CreatedAt: permission.CreatedAt,
			UpdatedAt: permission.UpdatedAt,
			Status:    permission.Status,
		})
	}
	span.AddEvent("permission_id is validated", trace.WithAttributes(attribute.Int64Slice("permissions", req.Permissions)))

	// create role
	span.AddEvent("creating new role")
	err = rRole.Create(ctx, role)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("role is created", trace.WithAttributes(attribute.Int64("role_id", role.ID)))

	// create role permission
	span.AddEvent("creating role permission id")
	for _, permission := range req.Permissions {
		rolePermission := model.RolePermission{
			RoleID:       role.ID,
			PermissionID: permission,
			CreatedAt:    time.Now(),
			Status:       1,
		}
		err = rRolePermission.Create(ctx, &rolePermission)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}
	span.AddEvent("role permission is created", trace.WithAttributes(attribute.Int64("role_id", role.ID)))

	res = dto.RoleResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
		Status:      role.Status,
	}

	return
}

func (s *RoleService) Update(ctx context.Context, req dto.RoleRequestUpdate, id int64) (res dto.RoleResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RoleService.Update")
	defer span.End()

	rPermission := repository.RepositoryPermission()
	rRole := repository.RepositoryRole()
	rRolePermission := repository.RepositoryRolePermssion()
	rDivision := repository.RepositoryDivision()
	// validate role id
	var roleOld *model.Role
	roleOld, err = rRole.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if roleOld.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	// validate division id
	span.AddEvent("validate division id")
	_, err = rDivision.GetDetail(ctx, req.DivisionID, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("division is validated", trace.WithAttributes(attribute.Int64("division_id", req.DivisionID)))

	// validate permission id
	span.AddEvent("validate permission_id")
	var permissions []*dto.PermissionResponse
	for _, permission_id := range req.Permissions {
		var permission *model.Permission
		permission, err = rPermission.GetByID(ctx, permission_id)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		permissions = append(permissions, &dto.PermissionResponse{
			ID:        permission.ID,
			Name:      permission.Name,
			CreatedAt: permission.CreatedAt,
			UpdatedAt: permission.UpdatedAt,
			Status:    permission.Status,
		})
	}
	span.AddEvent("permission_id is validated", trace.WithAttributes(attribute.Int64Slice("permissions", req.Permissions)))

	// update role
	role := &model.Role{
		ID:         id,
		Note:       req.Note,
		DivisionID: req.DivisionID,
		Name:       req.Name,
		UpdatedAt:  time.Now(),
	}
	err = rRole.Update(ctx, role, "Note", "DivisionID", "Name", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// get old role permission
	var rolePermissions []*model.RolePermission
	rolePermissions, err = rRolePermission.GetByRoleID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// delete old permission
	for _, rolePermission := range rolePermissions {
		err = rRolePermission.Delete(ctx, rolePermission.ID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// insert new role permission
	for _, reqPermissionID := range req.Permissions {
		err = rRolePermission.Create(ctx, &model.RolePermission{
			RoleID:       id,
			PermissionID: reqPermissionID,
			Status:       1,
			CreatedAt:    time.Now(),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	res = dto.RoleResponse{
		ID:          role.ID,
		Code:        role.Code,
		Name:        role.Name,
		Permissions: permissions,
		CreatedAt:   role.CreatedAt,
		UpdatedAt:   role.UpdatedAt,
		Status:      role.Status,
	}

	return
}

func (s *RoleService) Archive(ctx context.Context, id int64) (res dto.RoleResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RoleService.Archive")
	defer span.End()

	rRole := repository.RepositoryRole()
	rUserRole := repository.RepositoryUserRole()
	// validate data is exist
	var roleOld *model.Role
	roleOld, err = rRole.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var userRolesActive []*model.UserRole
	userRolesActive, _ = rUserRole.GetByRoleID(ctx, id)
	if len(userRolesActive) != 0 {
		err = edenlabs.ErrorValidation("status", "The role still have active users, please check users active")
		return
	}

	if roleOld.Status == 2 {
		err = edenlabs.ErrorValidation("status", "The status has been archived")
		return
	}

	if roleOld.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	role := &model.Role{
		ID:        id,
		Status:    2,
		UpdatedAt: time.Now(),
	}

	err = rRole.Update(ctx, role, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.RoleResponse{
		ID:        role.ID,
		Code:      role.Code,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		Status:    role.Status,
	}

	return
}

func (s *RoleService) UnArchive(ctx context.Context, id int64) (res dto.RoleResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "RoleService.UnArchive")
	defer span.End()

	rRole := repository.RepositoryRole()
	// validate data is exist
	var roleOld *model.Role
	roleOld, err = rRole.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if roleOld.Status != 2 {
		err = edenlabs.ErrorValidation("status", "The status is not archived")
		return
	}

	role := &model.Role{
		ID:        id,
		Status:    1,
		UpdatedAt: time.Now(),
	}

	err = rRole.Update(ctx, role, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.RoleResponse{
		ID:        role.ID,
		Code:      role.Code,
		Name:      role.Name,
		CreatedAt: role.CreatedAt,
		UpdatedAt: role.UpdatedAt,
		Status:    role.Status,
	}

	return
}
