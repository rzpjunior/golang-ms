package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-audit-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/internal/app/repository"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type IAuditService interface {
	Get(ctx context.Context, offset int, limit int, auditType string, referenceID string) (res []dto.AuditResponseGet, total int64, err error)
	Create(ctx context.Context, req dto.AuditRequestCreate) (res dto.AuditResponseCreate, err error)
}

type AuditService struct {
	opt             opt.Options
	RepositoryAudit repository.IAuditRepository
}

func NewAuditService() IAuditService {
	return &AuditService{
		opt:             global.Setup.Common,
		RepositoryAudit: repository.NewAuditRepository(),
	}
}

func (s *AuditService) Get(ctx context.Context, offset int, limit int, auditType string, referenceID string) (res []dto.AuditResponseGet, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AuditService.Get")
	defer span.End()

	var auditLogs []*model.Log
	auditLogs, total, err = s.RepositoryAudit.Get(ctx, int64(offset), int64(limit), auditType, referenceID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, auditLog := range auditLogs {
		mongoId := auditLog.ID.Hex()
		var user *accountService.GetUserDetailResponse
		user, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
			Id: auditLog.UserID,
		})

		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("account", "user")
			return
		}
		res = append(res, dto.AuditResponseGet{
			ID:          mongoId,
			UserID:      auditLog.UserID,
			ReferenceID: auditLog.ReferenceID,
			Type:        auditLog.Type,
			Function:    auditLog.Function,
			CreatedAt:   auditLog.CreatedAt,
			Note:        auditLog.Note,
			User: &dto.UserResponse{
				ID:           user.Data.Id,
				Name:         user.Data.Name,
				Email:        user.Data.Email,
				MainRole:     user.Data.MainRole,
				Division:     user.Data.Division,
				EmployeeCode: user.Data.EmployeeCode,
				Nickname:     user.Data.Nickname,
				Status:       int8(user.Data.Status),
			},
		})
	}

	return
}

func (s *AuditService) Create(ctx context.Context, req dto.AuditRequestCreate) (res dto.AuditResponseCreate, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AuditService.Create")
	defer span.End()

	auditLog := &model.Log{
		UserID:      req.UserID,
		UserIdGp:    req.UserIdGp,
		ReferenceID: req.ReferenceID,
		Type:        req.Type,
		Function:    req.Function,
		CreatedAt:   time.Now(),
		Note:        req.Note,
	}

	span.AddEvent("creating new audit log")
	err = s.RepositoryAudit.Create(ctx, auditLog)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	span.AddEvent("audit log is created", trace.WithAttributes(attribute.String("log_id", auditLog.ID.String())))

	res = dto.AuditResponseCreate{
		UserID:      auditLog.UserID,
		UserIdGp:    auditLog.UserIdGp,
		ReferenceID: auditLog.ReferenceID,
		Type:        auditLog.Type,
		Function:    auditLog.Function,
		CreatedAt:   auditLog.CreatedAt,
		Note:        auditLog.Note,
	}

	return
}
