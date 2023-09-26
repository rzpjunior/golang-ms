package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"

	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	dto "git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"

	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/repository"
)

func ServicePermission() IPermissionService {
	m := new(PermissionService)
	m.opt = global.Setup.Common
	return m
}

type IPermissionService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.PermissionResponse, total int64, err error)
	GetTree(ctx context.Context) (res []*dto.PermissionResponse, err error)
	GetPrivilege(ctx context.Context, UserID int64) (res []string, err error)
}

type PermissionService struct {
	opt opt.Options
}

func (s *PermissionService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.PermissionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PermissionService.Get")

	defer span.End()

	rPermission := repository.RepositoryPermission()
	var permissions []*model.Permission

	permissions, total, err = rPermission.Get(ctx, offset, limit, status, search, orderBy)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, permission := range permissions {
		res = append(res, &dto.PermissionResponse{
			ID:            permission.ID,
			Name:          permission.Name,
			Value:         permission.Value,
			CreatedAt:     permission.CreatedAt,
			UpdatedAt:     permission.UpdatedAt,
			Status:        permission.Status,
			StatusConvert: statusx.ConvertStatusValue(permission.Status),
		})
	}

	return
}

func (s *PermissionService) GetTree(ctx context.Context) (res []*dto.PermissionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PermissionService.GetTree")
	defer span.End()

	rPermission := repository.RepositoryPermission()
	var permissions []*model.Permission
	permissions, err = rPermission.GetTree(ctx)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, permission := range permissions {
		var childPermission []*dto.PermissionResponse
		for _, child := range permission.Child {
			var grandChildPermission []*dto.PermissionResponse
			for _, grandChild := range child.GrandChild {
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

		res = append(res, &dto.PermissionResponse{
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

	return
}

func (s *PermissionService) GetPrivilege(ctx context.Context, UserID int64) (res []string, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PermissionService.GetPrivilege")
	defer span.End()

	rPermission := repository.RepositoryPermission()
	res, err = rPermission.GetPrivilege(ctx, UserID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
