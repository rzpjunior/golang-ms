package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/repository"
)

func ServiceMenu() IMenuService {
	m := new(MenuService)
	m.opt = global.Setup.Common
	return m
}

type IMenuService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.MenuResponse, total int64, err error)
}

type MenuService struct {
	opt opt.Options
}

func (s *MenuService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.MenuResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "MenuService.Get")

	defer span.End()
	rMenu := repository.RepositoryMenu()
	var menus []*model.Menu

	menus, total, err = rMenu.Get(ctx, offset, limit, status, search, orderBy)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, menu := range menus {
		var childResponse []*dto.MenuResponse
		for _, child := range menu.Child {
			childResponse = append(childResponse, &dto.MenuResponse{
				ID:       child.ID,
				ParentID: child.ParentID,
				Title:    child.Title,
				Url:      child.Url,
				Icon:     child.Icon,
			})
		}
		res = append(res, &dto.MenuResponse{
			ID:       menu.ID,
			ParentID: menu.ParentID,
			Title:    menu.Title,
			Icon:     menu.Icon,
			Child:    childResponse,
		})
	}

	return
}
