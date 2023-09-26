package repository

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/generator"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IBankRepository interface {
	Get(ctx context.Context, req *pb.GetBankListRequest) (Bankes []*model.Bank, count int64, err error)
	GetDetail(ctx context.Context, req *pb.GetBankDetailRequest) (Bank *model.Bank, err error)
}

type BankRepository struct {
	opt opt.Options
}

func NewBankRepository() IBankRepository {
	return &BankRepository{
		opt: global.Setup.Common,
	}
}

func (r *BankRepository) Get(ctx context.Context, req *pb.GetBankListRequest) (Bankes []*model.Bank, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "BankRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	db := r.opt.Database.Read

	qs := db.QueryTable(new(model.Bank))

	cond := orm.NewCondition()

	if req.Search != "" {
		condGroup := orm.NewCondition()
		condGroup = condGroup.And("description__icontains", req.Search).Or("code__icontains", req.Search)
		cond = cond.AndCond(condGroup)
	}

	if len(req.Status) != 0 {
		cond = cond.And("status__in", req.Status)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	count, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &Bankes)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *BankRepository) GetDetail(ctx context.Context, req *pb.GetBankDetailRequest) (Bank *model.Bank, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "BankRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil

	Bank = &model.Bank{}

	var cols []string

	if req.Id != 0 {
		cols = append(cols, "id")
		Bank.ID = req.Id
	}

	if req.Code != "" {
		cols = append(cols, "code")
		Bank.Code = req.Code
	}

	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, Bank, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *BankRepository) MockDatas(total int) (mockDatas []*model.Bank) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Bank{
				ID:              int64(i),
				Code:            fmt.Sprintf("DummyBank%d", i),
				Description:     fmt.Sprintf("Dummy Bank %d", i),
				Status:          1,
				CreatedAt:       generator.DummyTime(),
				UpdatedAt:       generator.DummyTime(),
				Value:           fmt.Sprintf("DummyValue%d", i),
				ImageUrl:        "https://sgp1.digitaloceanspaces.com/image-erp-dev-eden/item/vecteezy_modern-vector-graphic-troly-colorful-logo-good-for_11883295.jpg",
				PaymentGuideUrl: "https://www.littlethings.info/wp-content/uploads/2014/04/dummy-image-grey-e1398449111870.jpg",
				PublishIVA:      0,
				PublishFVA:      0,
			})
	}
	return
}
