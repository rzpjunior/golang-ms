package repository

import (
	"context"
	"encoding/json"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-audit-service/internal/app/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IAuditRepository interface {
	Get(ctx context.Context, offset int64, limit int64, auditType string, referenceID string) (auditLogs []*model.Log, count int64, err error)
	Create(ctx context.Context, auditLog *model.Log) (err error)
}

type AuditRepository struct {
	opt opt.Options
}

func NewAuditRepository() IAuditRepository {
	return &AuditRepository{
		opt: global.Setup.Common,
	}
}

func (r *AuditRepository) Get(ctx context.Context, offset int64, limit int64, auditType string, referenceID string) (auditLogs []*model.Log, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AuditRepository.Get")
	defer span.End()

	db := r.opt.Mongox

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(limit).SetSkip(offset)

	var ret []byte
	ret, err = db.GetByFilter(ctx, "audit_log", &model.Log{
		Type:        auditType,
		ReferenceID: referenceID,
	}, opts)

	if err != nil {
		span.RecordError(err)
		return
	}

	err = json.Unmarshal(ret, &auditLogs)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *AuditRepository) Create(ctx context.Context, auditLog *model.Log) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "AuditRepository.Create")
	defer span.End()

	db := r.opt.Mongox
	db.CreateIndex(ctx, "audit_log", "user_id", false)
	db.CreateIndex(ctx, "audit_log", "user_id_gp", false)
	_, err = db.Insert(ctx, "audit_log", auditLog)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
