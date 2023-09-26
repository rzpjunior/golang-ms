package service

import (
	"context"
	"net/url"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/repository"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type ISalesAssignmentObjectiveService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, codes []string, orderBy string) (res []*dto.SalesAssignmentObjectiveResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.SalesAssignmentObjectiveResponse, err error)
	Create(ctx context.Context, req dto.SalesAssignmentObjectiveRequestCreate) (res dto.SalesAssignmentObjectiveResponse, err error)
	Update(ctx context.Context, req dto.SalesAssignmentObjectiveRequestUpdate, id int64) (res dto.SalesAssignmentObjectiveResponse, err error)
	Archive(ctx context.Context, id int64) (res dto.SalesAssignmentObjectiveResponse, err error)
	UnArchive(ctx context.Context, id int64) (res dto.SalesAssignmentObjectiveResponse, err error)
}

type SalesAssignmentObjectiveService struct {
	opt                                opt.Options
	RepositorySalesAssignmentObjective repository.ISalesAssignmentObjectiveRepository
}

func NewSalesAssignmentObjectiveService() ISalesAssignmentObjectiveService {
	return &SalesAssignmentObjectiveService{
		opt:                                global.Setup.Common,
		RepositorySalesAssignmentObjective: repository.NewSalesAssignmentObjectiveRepository(),
	}
}

func (s *SalesAssignmentObjectiveService) Get(ctx context.Context, offset int, limit int, status int, search string, codes []string, orderBy string) (res []*dto.SalesAssignmentObjectiveResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentObjectiveService.Get")
	defer span.End()

	var salesAssignmentObjectives []*model.SalesAssignmentObjective
	salesAssignmentObjectives, total, err = s.RepositorySalesAssignmentObjective.Get(ctx, offset, limit, status, search, codes, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesAssignmentObjective := range salesAssignmentObjectives {
		var user *accountService.GetUserDetailResponse
		user, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
			Id: salesAssignmentObjective.CreatedBy,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("account", "user")
			return
		}

		res = append(res, &dto.SalesAssignmentObjectiveResponse{
			ID:            salesAssignmentObjective.ID,
			Code:          salesAssignmentObjective.Code,
			Name:          salesAssignmentObjective.Name,
			Objective:     salesAssignmentObjective.Objective,
			SurveyLink:    salesAssignmentObjective.SurveyLink,
			Status:        salesAssignmentObjective.Status,
			StatusConvert: statusx.ConvertStatusValue(salesAssignmentObjective.Status),
			CreatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.CreatedAt),
			CreatedBy: &dto.CreatedByResponse{
				ID:   user.Data.Id,
				Name: user.Data.Name,
			},
			UpdatedAt: timex.ToLocTime(ctx, salesAssignmentObjective.UpdatedAt),
		})
	}

	return
}

func (s *SalesAssignmentObjectiveService) GetByID(ctx context.Context, id int64) (res dto.SalesAssignmentObjectiveResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentObjectiveService.GetByID")
	defer span.End()

	var salesAssignmentObjective *model.SalesAssignmentObjective
	salesAssignmentObjective, err = s.RepositorySalesAssignmentObjective.GetByID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var user *accountService.GetUserDetailResponse
	user, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		Id: salesAssignmentObjective.CreatedBy,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user")
		return
	}

	res = dto.SalesAssignmentObjectiveResponse{
		ID:            salesAssignmentObjective.ID,
		Code:          salesAssignmentObjective.Code,
		Name:          salesAssignmentObjective.Name,
		Objective:     salesAssignmentObjective.Objective,
		SurveyLink:    salesAssignmentObjective.SurveyLink,
		Status:        salesAssignmentObjective.Status,
		StatusConvert: statusx.ConvertStatusValue(salesAssignmentObjective.Status),
		CreatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.CreatedAt),
		CreatedBy: &dto.CreatedByResponse{
			ID:   user.Data.Id,
			Name: user.Data.Name,
		},
		UpdatedAt: timex.ToLocTime(ctx, salesAssignmentObjective.UpdatedAt),
	}
	return
}

func (s *SalesAssignmentObjectiveService) Create(ctx context.Context, req dto.SalesAssignmentObjectiveRequestCreate) (res dto.SalesAssignmentObjectiveResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentObjectiveService.Create")
	defer span.End()

	// validate survey link
	_, err = url.ParseRequestURI(req.SurveyLink)
	if err != nil {
		err = edenlabs.ErrorInvalid("survey_link")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "SOB",
		Domain: "sales_assignment_objective",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "code_generator")
		return
	}
	code := codeGenerator.Data.Code

	salesAssignmentObjective := &model.SalesAssignmentObjective{
		Code:       code,
		Name:       req.Name,
		Objective:  req.Objective,
		SurveyLink: req.SurveyLink,
		Status:     statusx.ConvertStatusName(statusx.Active),
		CreatedAt:  time.Now(),
		CreatedBy:  ctx.Value(constants.KeyUserID).(int64),
	}

	err = s.RepositorySalesAssignmentObjective.Create(ctx, salesAssignmentObjective)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	res = dto.SalesAssignmentObjectiveResponse{
		ID:            salesAssignmentObjective.ID,
		Code:          salesAssignmentObjective.Code,
		Name:          salesAssignmentObjective.Name,
		Objective:     salesAssignmentObjective.Objective,
		SurveyLink:    salesAssignmentObjective.SurveyLink,
		Status:        salesAssignmentObjective.Status,
		StatusConvert: statusx.ConvertStatusValue(salesAssignmentObjective.Status),
		CreatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.UpdatedAt),
	}

	return
}

func (s *SalesAssignmentObjectiveService) Update(ctx context.Context, req dto.SalesAssignmentObjectiveRequestUpdate, id int64) (res dto.SalesAssignmentObjectiveResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentObjectiveService.Update")
	defer span.End()

	// validate sales assignment object is exist
	var salesAssignmentObjectiveOld *model.SalesAssignmentObjective
	salesAssignmentObjectiveOld, err = s.RepositorySalesAssignmentObjective.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if salesAssignmentObjectiveOld.Status != statusx.ConvertStatusName(statusx.Active) {
		err = edenlabs.ErrorMustActive("status")
		return
	}

	// validate survey link
	_, err = url.ParseRequestURI(req.SurveyLink)
	if err != nil {
		err = edenlabs.ErrorInvalid("survey_link")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	salesAssignmentObjective := &model.SalesAssignmentObjective{
		ID:         id,
		Objective:  req.Objective,
		SurveyLink: req.SurveyLink,
		Status:     statusx.ConvertStatusName(statusx.Active),
		UpdatedAt:  time.Now(),
	}

	err = s.RepositorySalesAssignmentObjective.Update(ctx, salesAssignmentObjective, "Objective", "SurveyLink", "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesAssignmentObjectiveResponse{
		ID:            salesAssignmentObjective.ID,
		Code:          salesAssignmentObjective.Code,
		Name:          salesAssignmentObjective.Name,
		Objective:     salesAssignmentObjective.Objective,
		SurveyLink:    salesAssignmentObjective.SurveyLink,
		Status:        salesAssignmentObjective.Status,
		StatusConvert: statusx.ConvertStatusValue(salesAssignmentObjective.Status),
		CreatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.UpdatedAt),
	}

	return
}

func (s *SalesAssignmentObjectiveService) Archive(ctx context.Context, id int64) (res dto.SalesAssignmentObjectiveResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentObjectiveService.Archive")
	defer span.End()

	// validate sales assignment object is exist
	var salesAssignmentObjectiveOld *model.SalesAssignmentObjective
	salesAssignmentObjectiveOld, err = s.RepositorySalesAssignmentObjective.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if salesAssignmentObjectiveOld.Status != statusx.ConvertStatusName(statusx.Active) {
		err = edenlabs.ErrorMustActive("status")
		return
	}

	salesAssignmentObjective := &model.SalesAssignmentObjective{
		ID:        id,
		Status:    statusx.ConvertStatusName(statusx.Archived),
		UpdatedAt: time.Now(),
	}

	err = s.RepositorySalesAssignmentObjective.Update(ctx, salesAssignmentObjective, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesAssignmentObjectiveResponse{
		ID:            salesAssignmentObjective.ID,
		Code:          salesAssignmentObjective.Code,
		Name:          salesAssignmentObjective.Name,
		Objective:     salesAssignmentObjective.Objective,
		SurveyLink:    salesAssignmentObjective.SurveyLink,
		Status:        salesAssignmentObjective.Status,
		StatusConvert: statusx.ConvertStatusValue(salesAssignmentObjective.Status),
		CreatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.UpdatedAt),
	}

	return
}

func (s *SalesAssignmentObjectiveService) UnArchive(ctx context.Context, id int64) (res dto.SalesAssignmentObjectiveResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentObjectiveService.UnArchive")
	defer span.End()

	// validate sales assignment object is exist
	var salesAssignmentObjectiveOld *model.SalesAssignmentObjective
	salesAssignmentObjectiveOld, err = s.RepositorySalesAssignmentObjective.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if salesAssignmentObjectiveOld.Status != statusx.ConvertStatusName(statusx.Archived) {
		err = edenlabs.ErrorMustArchived("status")
		return
	}

	salesAssignmentObjective := &model.SalesAssignmentObjective{
		ID:        id,
		Status:    statusx.ConvertStatusName(statusx.Active),
		UpdatedAt: time.Now(),
	}

	err = s.RepositorySalesAssignmentObjective.Update(ctx, salesAssignmentObjective, "Status", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesAssignmentObjectiveResponse{
		ID:            salesAssignmentObjective.ID,
		Code:          salesAssignmentObjective.Code,
		Name:          salesAssignmentObjective.Name,
		Objective:     salesAssignmentObjective.Objective,
		SurveyLink:    salesAssignmentObjective.SurveyLink,
		Status:        salesAssignmentObjective.Status,
		StatusConvert: statusx.ConvertStatusValue(salesAssignmentObjective.Status),
		CreatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, salesAssignmentObjective.UpdatedAt),
	}

	return
}
