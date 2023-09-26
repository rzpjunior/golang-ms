package repository

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
)

type ICodeGeneratorRepository interface {
	GetLastCode(ctx context.Context, domain string, format string) (codeGenerator *model.CodeGenerator, err error)
	Create(ctx context.Context, codeGenerator *model.CodeGenerator) (err error)
	CreateReferral(ctx context.Context, codeGenerator *model.CodeGeneratorReferral) (err error)
}

type CodeGeneratorRepository struct {
	opt opt.Options
}

func NewCodeGeneratorRepository() ICodeGeneratorRepository {
	return &CodeGeneratorRepository{
		opt: global.Setup.Common,
	}
}

func (r *CodeGeneratorRepository) GetLastCode(ctx context.Context, domain string, format string) (codeGenerator *model.CodeGenerator, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CodeGeneratorRepository.GetByID")
	defer span.End()

	db := r.opt.Database.Read

	cg := model.CodeGenerator{}

	err = db.QueryTable(new(model.CodeGenerator)).Filter("code__icontains", format).Filter("domain", domain).OrderBy("-id").Limit(1).OneWithCtx(ctx, &cg)
	if err != nil {
		span.RecordError(err)
		return
	}

	codeGenerator = &cg

	return
}

func (r *CodeGeneratorRepository) Create(ctx context.Context, codeGenerator *model.CodeGenerator) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CodeGeneratorRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, codeGenerator)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *CodeGeneratorRepository) CreateReferral(ctx context.Context, codeGenerator *model.CodeGeneratorReferral) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CodeGeneratorRepository.CreateReferral")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.InsertWithCtx(ctx, codeGenerator)
	if err != nil {
		span.RecordError(err)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}
