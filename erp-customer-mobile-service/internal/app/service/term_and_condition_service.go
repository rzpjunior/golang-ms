package service

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/repository"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"github.com/labstack/echo/v4"
)

type ITermConditionService interface {
	Get(ctx context.Context, offset int, limit int) (res []dto.TermConditionResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64) (res dto.TermConditionResponse, err error)
	AcceptTNC(ctx context.Context, req *dto.RequestAcceptTNC) (res dto.TermConditionResponse, err error)
}

type TermConditionService struct {
	opt                     opt.Options
	RepositoryTermCondition repository.ITermConditionRepository
	RepositoryUserCustomer  repository.IUserCustomerRepository
}

func NewTermConditionService() ITermConditionService {
	return &TermConditionService{
		opt:                     global.Setup.Common,
		RepositoryTermCondition: repository.NewTermConditionRepository(),
		RepositoryUserCustomer:  repository.NewUserCustomerRepository(),
	}
}

func (s *TermConditionService) Get(ctx context.Context, offset int, limit int) (res []dto.TermConditionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TermConditionService.Get")
	defer span.End()

	var termConditions []*model.TermCondition
	termConditions, total, err = s.RepositoryTermCondition.Get(ctx, offset, limit)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, termCondition := range termConditions {
		res = append(res, dto.TermConditionResponse{
			ID:          strconv.Itoa(int(termCondition.ID)),
			Application: strconv.Itoa(int(termCondition.Application)),
			Version:     termCondition.Version,
			Title:       termCondition.Title,
			TitleValue:  termCondition.TitleValue,
			Content:     termCondition.Content,
		})
	}

	return
}

func (s *TermConditionService) GetDetail(ctx context.Context, id int64) (res dto.TermConditionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TermConditionService.GetDetail")
	defer span.End()

	var termCondition *model.TermCondition
	termCondition, err = s.RepositoryTermCondition.GetDetail(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.TermConditionResponse{
		ID:          strconv.Itoa(int(termCondition.ID)),
		Application: strconv.Itoa(int(termCondition.Application)),
		Version:     termCondition.Version,
		Title:       termCondition.Title,
		TitleValue:  termCondition.TitleValue,
		Content:     termCondition.Content,
	}

	return
}

func (s *TermConditionService) AcceptTNC(ctx context.Context, req *dto.RequestAcceptTNC) (res dto.TermConditionResponse, err error) {

	ctx, span := s.opt.Trace.Start(ctx, "TermConditionService.GetDetail")
	defer span.End()
	userCustomer := &model.UserCustomer{CustomerID: utils.ToInt64(req.Session.Customer.ID), Status: 1}
	userCustomer, err = s.RepositoryUserCustomer.GetDetail(ctx, userCustomer)
	if err != nil {
		err = echo.NewHTTPError(http.StatusRequestTimeout, "invalid or expired jwt token")
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return res, err
	}
	tncVersion, err := s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configuration_service.GetConfigAppDetailRequest{
		Attribute: "tnc_current_version",
	})
	if err != nil {
		err = echo.NewHTTPError(http.StatusRequestTimeout, "invalid")
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return res, err
	}

	userCustomer.TncAccAt = time.Now()
	userCustomer.TncAccVersion = tncVersion.Data.Value

	err = s.RepositoryUserCustomer.Update(ctx, userCustomer)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.TermConditionResponse{
		ID:          utils.ToString(tncVersion.Data.Id),
		Application: utils.ToString(tncVersion.Data.Application),
		Version:     tncVersion.Data.Value,
	}

	return
}
