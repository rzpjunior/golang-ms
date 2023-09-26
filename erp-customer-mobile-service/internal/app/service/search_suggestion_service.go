package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
)

type ISearchSuggestionService interface {
	Get(ctx context.Context, req dto.SearchSuggestionRequest) (res []dto.SearchSuggestionResponse, err error)
}

type SearchSuggestionService struct {
	opt opt.Options
}

func NewSearchSuggestionService() ISearchSuggestionService {
	return &SearchSuggestionService{
		opt: global.Setup.Common,
	}
}

func (s *SearchSuggestionService) Get(ctx context.Context, req dto.SearchSuggestionRequest) (res []dto.SearchSuggestionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SearchSuggestionService.Get")
	defer span.End()

	//

	searchSuggestionList, err := s.opt.Client.CatalogServiceGrpc.GetItemList(ctx, &catalog_service.GetItemListRequest{
		Search: req.Data.Search,
		Status: 1,
		Limit:  5,
		Offset: 1,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range searchSuggestionList.Data {
		searchSuggestion := dto.SearchSuggestionResponse{
			ID:   v.Id,
			Name: v.Description,
			Code: v.Code,
		}
		res = append(res, searchSuggestion)
	}

	return
}
