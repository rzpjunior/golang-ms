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

func ServiceProfile() IProfileService {
	m := new(ProfileService)
	m.opt = global.Setup.Common
	return m
}

type IProfileService interface {
	GetMenu(ctx context.Context, UserID int64) (res []*dto.MenuResponse, err error)
}

type ProfileService struct {
	opt opt.Options
}

func (s *ProfileService) GetMenu(ctx context.Context, userID int64) (res []*dto.MenuResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ProfileService.GetMenu")
	defer span.End()
	rMenu := repository.RepositoryMenu()

	var menus []*model.Menu
	menus, err = rMenu.GetByUserID(ctx, userID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var menuResponse *dto.MenuResponse
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
		menuResponse = &dto.MenuResponse{
			ID:       menu.ID,
			ParentID: menu.ParentID,
			Url:      menu.Url,
			Title:    menu.Title,
			Icon:     menu.Icon,
			Child:    childResponse,
		}
		res = append(res, menuResponse)
	}

	return
}
