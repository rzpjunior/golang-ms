package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IDivisionRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (divisions []*model.Division, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (division *model.Division, err error)
}

type DivisionRepository struct {
	opt opt.Options
}

func NewDivisionRepository() IDivisionRepository {
	return &DivisionRepository{
		opt: global.Setup.Common,
	}
}

func (r *DivisionRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (divisions []*model.Division, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ArchetypeRepository.Get")
	defer span.End()

	// RETURN DUMMIES
	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil
}

func (r *DivisionRepository) GetDetail(ctx context.Context, id int64, code string) (division *model.Division, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "DivisionRepository.GetDetail")
	defer span.End()

	// RETURN DUMMIES
	return r.MockDatas(1)[0], nil
}

func (r *DivisionRepository) MockDatas(total int) (mockDatas []*model.Division) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Division{
				ID:     int64(i),
				Code:   fmt.Sprintf("DIV000%d", i),
				Name:   fmt.Sprintf("Dummy Divsion 000%d", i),
				Note:   fmt.Sprintf("Dummy Divsion Note 000%d", i),
				Status: 1,
			})
	}

	return mockDatas
}
