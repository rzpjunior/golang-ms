package service

import (
	"context"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/repository"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type ISalesFailedVisitService interface {
	Get(ctx context.Context, offset int, limit int, salesFailedVisitItem int64, failedStatus int32, orderBy string) (res []*dto.SalesFailedVisitResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.SalesFailedVisitResponse, err error)
	Create(ctx context.Context, req *dto.SalesFailedVisitRequest) (res dto.SalesFailedVisitResponse, err error)
	Update(ctx context.Context, id int64, req *dto.SalesFailedVisitRequest) (res dto.SalesFailedVisitResponse, err error)
	SubmitTaskFailed(ctx context.Context, req dto.SalesFailedVisitRequest) (res *dto.SalesFailedVisitResponse, err error)
}

type SalesFailedVisitService struct {
	opt                           opt.Options
	RepositorySalesFailedVisit    repository.ISalesFailedVisitRepository
	RepositorySalesAssignmentItem repository.ISalesAssignmentItemRepository
}

func NewSalesFailedVisitService() ISalesFailedVisitService {
	return &SalesFailedVisitService{
		opt:                           global.Setup.Common,
		RepositorySalesFailedVisit:    repository.NewSalesFailedVisitRepository(),
		RepositorySalesAssignmentItem: repository.NewSalesAssignmentItemRepository(),
	}
}

func (s *SalesFailedVisitService) Get(ctx context.Context, offset int, limit int, salesFailedVisitItem int64, failedStatus int32, orderBy string) (res []*dto.SalesFailedVisitResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesFailedVisitService.Get")
	defer span.End()

	var salesFailedVisits []*model.SalesFailedVisit
	salesFailedVisits, total, err = s.RepositorySalesFailedVisit.Get(ctx, offset, limit, salesFailedVisitItem, failedStatus, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesFailedVisit := range salesFailedVisits {
		res = append(res, &dto.SalesFailedVisitResponse{
			ID:                    salesFailedVisit.ID,
			SalesAssignmentItemID: salesFailedVisit.SalesAssignmentItemID,
			FailedStatus:          salesFailedVisit.FailedStatus,
			DescriptionFailed:     salesFailedVisit.DescriptionFailed,
			FailedImage:           strings.Split(salesFailedVisit.FailedImage, ","),
		})
	}

	return
}

func (s *SalesFailedVisitService) GetByID(ctx context.Context, id int64) (res dto.SalesFailedVisitResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesFailedVisitService.GetByID")
	defer span.End()

	var salesFailedVisit *model.SalesFailedVisit
	salesFailedVisit, err = s.RepositorySalesFailedVisit.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesFailedVisitResponse{
		ID:                    salesFailedVisit.ID,
		SalesAssignmentItemID: salesFailedVisit.SalesAssignmentItemID,
		FailedStatus:          salesFailedVisit.FailedStatus,
		DescriptionFailed:     salesFailedVisit.DescriptionFailed,
		FailedImage:           strings.Split(salesFailedVisit.FailedImage, ","),
	}

	return
}

func (s *SalesFailedVisitService) Create(ctx context.Context, req *dto.SalesFailedVisitRequest) (res dto.SalesFailedVisitResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesFailedVisitService.Create")
	defer span.End()

	var salesFailedVisit *model.SalesFailedVisit
	id, err := s.RepositorySalesFailedVisit.Create(ctx, &model.SalesFailedVisit{
		SalesAssignmentItemID: req.SalesAssignmentItemID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	salesFailedVisit, err = s.RepositorySalesFailedVisit.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesFailedVisitResponse{
		ID:                    salesFailedVisit.ID,
		SalesAssignmentItemID: salesFailedVisit.SalesAssignmentItemID,
		FailedStatus:          salesFailedVisit.FailedStatus,
		DescriptionFailed:     salesFailedVisit.DescriptionFailed,
		FailedImage:           strings.Split(salesFailedVisit.FailedImage, ","),
	}

	return
}

func (s *SalesFailedVisitService) Update(ctx context.Context, id int64, req *dto.SalesFailedVisitRequest) (res dto.SalesFailedVisitResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesFailedVisitService.Create")
	defer span.End()

	var salesFailedVisit *model.SalesFailedVisit
	err = s.RepositorySalesFailedVisit.Update(ctx, &model.SalesFailedVisit{
		ID:                id,
		FailedStatus:      req.FailedStatus,
		DescriptionFailed: req.DescriptionFailed,
		FailedImage:       req.FailedImage,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	salesFailedVisit, err = s.RepositorySalesFailedVisit.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesFailedVisitResponse{
		ID:                    salesFailedVisit.ID,
		SalesAssignmentItemID: salesFailedVisit.SalesAssignmentItemID,
		FailedStatus:          salesFailedVisit.FailedStatus,
		DescriptionFailed:     salesFailedVisit.DescriptionFailed,
		FailedImage:           strings.Split(salesFailedVisit.FailedImage, ","),
	}

	return
}

func (s *SalesFailedVisitService) SubmitTaskFailed(ctx context.Context, req dto.SalesFailedVisitRequest) (res *dto.SalesFailedVisitResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesFailedVisitService.SubmitTaskFailed")
	defer span.End()

	// validate sales assignment object is exist
	var salesAssignmentItem *model.SalesAssignmentItem
	salesAssignmentItem, err = s.RepositorySalesAssignmentItem.GetByID(ctx, req.SalesAssignmentItemID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if salesAssignmentItem.Status != statusx.ConvertStatusName(statusx.Active) {
		err = edenlabs.ErrorMustActive("status")
		return
	}

	if len(strings.Split(req.FailedImage, ",")) > 7 {
		edenlabs.ErrorMustEqualOrLess("task_image_urls", "7")
		return
	}

	// get glossary from configuration for validation task answer
	_, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryList(ctx, &configurationService.GetGlossaryListRequest{
		Table:     "sales_failed_visit",
		Attribute: "failed_status",
		ValueInt:  int32(req.FailedStatus),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}

	salesAssignmentItem = &model.SalesAssignmentItem{
		ID:         req.SalesAssignmentItemID,
		SubmitDate: time.Now(),
		Status:     14,
	}

	err = s.RepositorySalesAssignmentItem.Update(ctx, salesAssignmentItem, "SubmitDate", "Status")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	slsFailed := &model.SalesFailedVisit{
		SalesAssignmentItemID: req.SalesAssignmentItemID,
		FailedStatus:          req.FailedStatus,
		DescriptionFailed:     req.DescriptionFailed,
		FailedImage:           req.FailedImage,
	}

	_, err = s.RepositorySalesFailedVisit.Create(ctx, slsFailed)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.SalesFailedVisitResponse{
		ID:                    slsFailed.ID,
		SalesAssignmentItemID: slsFailed.SalesAssignmentItemID,
		FailedStatus:          slsFailed.FailedStatus,
		DescriptionFailed:     slsFailed.DescriptionFailed,
		FailedImage:           strings.Split(slsFailed.FailedImage, ","),
	}

	return
}
