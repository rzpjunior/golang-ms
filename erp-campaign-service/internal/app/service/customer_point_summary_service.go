package service

import (
	"context"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/repository"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ICustomerPointSummaryService interface {
	GetDetail(ctx context.Context, req *dto.CustomerPointSummaryRequestGetDetail) (res *dto.CustomerPointSummaryResponse, err error)
	Create(ctx context.Context, req *dto.CustomerPointSummaryRequestCreate) (err error)
	Update(ctx context.Context, req *dto.CustomerPointSummaryRequestUpdate) (err error)
}

type CustomerPointSummaryService struct {
	opt                            opt.Options
	RepositoryCustomerPointSummary repository.ICustomerPointSummaryRepository
}

func NewCustomerPointSummaryService() ICustomerPointSummaryService {
	return &CustomerPointSummaryService{
		opt:                            global.Setup.Common,
		RepositoryCustomerPointSummary: repository.NewCustomerPointSummaryRepository(),
	}
}

func (s *CustomerPointSummaryService) GetDetail(ctx context.Context, req *dto.CustomerPointSummaryRequestGetDetail) (res *dto.CustomerPointSummaryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointSummaryService.GetByID")
	defer span.End()

	var customerPointSummary *model.CustomerPointSummary
	customerPointSummary, err = s.RepositoryCustomerPointSummary.GetDetail(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.CustomerPointSummaryResponse{
		ID:            customerPointSummary.ID,
		CustomerID:    customerPointSummary.CustomerID,
		EarnedPoint:   customerPointSummary.EarnedPoint,
		RedeemedPoint: customerPointSummary.RedeemedPoint,
		SummaryDate:   customerPointSummary.SummaryDate,
	}

	return
}

func (s *CustomerPointSummaryService) Create(ctx context.Context, req *dto.CustomerPointSummaryRequestCreate) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointSummaryService.Create")
	defer span.End()

	customerPointSummary := &model.CustomerPointSummary{
		CustomerID:    req.CustomerID,
		EarnedPoint:   req.EarnedPoint,
		RedeemedPoint: req.RedeemedPoint,
		SummaryDate:   req.SummaryDate,
	}

	span.AddEvent("creating new customer point summary")
	err = s.RepositoryCustomerPointSummary.Create(ctx, customerPointSummary)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("customer point summary is created", trace.WithAttributes(attribute.Int64("customer_point_summary_id", customerPointSummary.ID)))

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      22,
			ReferenceId: strconv.Itoa(int(customerPointSummary.ID)),
			Type:        "customer_point_summary",
			Function:    "create",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *CustomerPointSummaryService) Update(ctx context.Context, req *dto.CustomerPointSummaryRequestUpdate) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointSummaryService.Update")
	defer span.End()

	customerPointSummary := &model.CustomerPointSummary{
		ID:            req.ID,
		CustomerID:    req.CustomerID,
		EarnedPoint:   req.EarnedPoint,
		RedeemedPoint: req.RedeemedPoint,
	}

	span.AddEvent("updating new customer point summary")
	err = s.RepositoryCustomerPointSummary.Update(ctx, customerPointSummary, req.FieldUpdate...)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("customer point summary is updated", trace.WithAttributes(attribute.Int64("customer_point_summary_id", customerPointSummary.ID)))

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      22,
			ReferenceId: strconv.Itoa(int(customerPointSummary.ID)),
			Type:        "customer_point_summary",
			Function:    "update",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
