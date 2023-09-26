package repository

import (
	"context"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-promotion-service/internal/app/model"
)

type IVoucherRepository interface {
	Get(ctx context.Context, req *dto.VoucherRequestGet) (vouchers []*model.Voucher, count int64, err error)
	GetDetail(ctx context.Context, id int64) (voucher *model.Voucher, err error)
	Create(ctx context.Context, voucher *model.Voucher) (err error)
	IsRedeemCodeExist(ctx context.Context, redeemCode string) (isExist bool)
	Update(ctx context.Context, voucher *model.Voucher, columns ...string) (err error)
	GetMobileVoucherList(ctx context.Context, req *dto.VoucherRequestGetMobileVoucherList) (vouchers []*model.Voucher, count int64, err error)
	GetMobileVoucherDetail(ctx context.Context, req *dto.VoucherRequestGetMobileVoucherDetail) (voucher *model.Voucher, err error)
}

type VoucherRepository struct {
	opt opt.Options
}

func NewVoucherRepository() IVoucherRepository {
	return &VoucherRepository{
		opt: global.Setup.Common,
	}
}

func (r *VoucherRepository) Get(ctx context.Context, req *dto.VoucherRequestGet) (vouchers []*model.Voucher, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherRepository.Get")
	defer span.End()

	db := r.opt.Database.Read

	cond := orm.NewCondition()

	qs := db.QueryTable(new(model.Voucher))

	if req.Search != "" {
		cond1 := orm.NewCondition()
		cond1 = cond1.And("name__icontains", req.Search).Or("redeem_code__icontains", req.Search)
		cond = cond.AndCond(cond1)
	}

	if req.RegionID != "" {
		cond = cond.And("region_id_gp", req.RegionID).Or("region_id_gp", "")
	}

	if req.ArchetypeID != "" {
		cond1 := orm.NewCondition()
		cond2 := orm.NewCondition()
		cond3 := orm.NewCondition()

		cond1 = cond1.And("customer_type_id_gp", "").And("archetype_id_gp", "")
		cond2 = cond2.And("customer_type_id_gp", req.CustomerTypeID).And("archetype_id_gp", "")
		cond3 = cond3.Or("archetype_id_gp", req.ArchetypeID)
		cond1 = cond1.OrCond(cond2)
		cond1 = cond1.OrCond(cond3)
		cond = cond.AndCond(cond1)
	}

	if req.CustomerID != 0 {
		cond = cond.And("customer_id", req.CustomerID)
	}

	if req.MembershipLevelID != 0 {
		cond = cond.And("membership_level_id", req.MembershipLevelID)
	}

	if req.MembershipCheckpointID != 0 {
		cond = cond.And("membership_checkpoint_id", req.MembershipCheckpointID)
	}

	if req.Type != 0 {
		cond = cond.And("type", req.Type)
	}

	if req.Status != 0 {
		cond = cond.And("status", req.Status)
	}

	qs = qs.SetCond(cond)

	if req.OrderBy != "" {
		qs = qs.OrderBy(req.OrderBy)
	}

	_, err = qs.Offset(req.Offset).Limit(req.Limit).AllWithCtx(ctx, &vouchers)
	if err != nil {
		span.RecordError(err)
		return
	}
	count, err = qs.Count()
	return
}

func (r *VoucherRepository) GetDetail(ctx context.Context, id int64) (voucher *model.Voucher, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherRepository.GetByID")
	defer span.End()

	voucher = &model.Voucher{
		ID: id,
	}
	db := r.opt.Database.Read
	err = db.ReadWithCtx(ctx, voucher, "id")
	if err != nil {
		span.RecordError(err)
		return
	}

	return
}

func (r *VoucherRepository) Create(ctx context.Context, voucher *model.Voucher) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherRepository.Create")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}
	voucher.ID, err = tx.InsertWithCtx(ctx, voucher)
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

// IsRedeemCodeExist : to check if there already exist redeem code at based on parameters
func (r *VoucherRepository) IsRedeemCodeExist(ctx context.Context, redeemCode string) (isExist bool) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherRepository.IsRedeemCodeExist")
	defer span.End()

	db := r.opt.Database.Read

	m := new(model.Voucher)

	isExist = db.QueryTable(m).Exclude("status", statusx.Archived).Filter("redeem_code", redeemCode).Exist()

	return isExist
}

func (r *VoucherRepository) Update(ctx context.Context, voucher *model.Voucher, columns ...string) (err error) {
	ctx, span := r.opt.Trace.Start(ctx, "ItemSectionRepository.Update")
	defer span.End()

	db := r.opt.Database.Write
	tx, err := db.BeginWithCtx(ctx)
	if err != nil {
		return
	}

	_, err = tx.UpdateWithCtx(ctx, voucher, columns...)

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

func (r *VoucherRepository) GetMobileVoucherList(ctx context.Context, req *dto.VoucherRequestGetMobileVoucherList) (vouchers []*model.Voucher, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherRepository.GetMobileVoucherDetail")
	defer span.End()

	db := r.opt.Database.Read

	type VoucherCount struct {
		VoucherId     int64 `orm:"column(voucher_id);auto" json:"-"`
		TotalUserUsed int64 `orm:"column(total_user_used);null" json:"-"`
	}
	var (
		voucherCount                           []*VoucherCount
		whereLevel, whereCheckpoint, whereType string
	)

	if req.Limit == 0 {
		req.Limit = 20
	}

	// set level filter
	if req.MembershipLevel == req.CustomerLevel {
		whereLevel = "v.membership_level_id <= ?"
		whereCheckpoint = " AND v.membership_checkpoint_id <= ?"
	} else {
		whereLevel = "v.membership_level_id = ?"
		whereCheckpoint = " AND v.membership_checkpoint_id = ?"
	}

	// set membership filter
	whereMembership := "AND (v.membership_level_id = 0 OR (v.membership_level_id != 0 AND " + whereLevel + ")) AND (v.membership_checkpoint_id = 0 OR (v.membership_checkpoint_id != 0" + whereCheckpoint + ")) "
	orderBy := "ORDER BY v.end_time DESC "
	if req.IsMembershipOnly {
		whereMembership = "AND v.membership_level_id != 0 AND " + whereLevel + " AND v.membership_checkpoint_id != 0" + whereCheckpoint + " "
		orderBy = "ORDER BY v.membership_level_id, v.membership_checkpoint_id, v.type "
	}

	values := []interface{}{req.RegionID, req.CustomerTypeID, req.ArchetypeID, req.CustomerID, req.MembershipLevel}
	// remove membership checkpoint filter if it's value is -1
	// and add it's value into value slice if not
	if req.MembershipCheckpoint == -1 {
		whereMembership = strings.ReplaceAll(whereMembership, " AND v.membership_checkpoint_id <= ?", "")
		whereMembership = strings.ReplaceAll(whereMembership, " AND v.membership_checkpoint_id = ?", "")
	} else {
		values = append(values, req.MembershipCheckpoint)
	}

	// Filter to get voucher delivery and non delivery
	if req.Category == 1 {
		whereType = " AND v.type = 2 "
	} else if req.Category == 2 {
		whereType = " AND v.type != 2 "
	}

	values = append(values, req.Offset, req.Limit)

	// get list of voucher data
	q := "SELECT v.id, v.image_url, v.redeem_code, v.name, v.min_order, v.end_time, v.start_time, v.user_quota, v.voucher_item, v.membership_level_id, v.membership_checkpoint_id, v.type, v.status, v.disc_amount " +
		"FROM voucher v " +
		"WHERE v.status = 1 AND v.rem_overall_quota > 0 AND CURRENT_TIMESTAMP() BETWEEN v.start_time AND v.end_time AND (v.region_id_gp = '' OR v.region_id_gp = ?) " +
		"AND ((v.customer_type_id_gp = '' AND v.archetype_id_gp = '') OR (v.customer_type_id_gp = ? AND v.archetype_id_gp = '') OR v.archetype_id_gp = ?) AND (v.customer_id = 0 OR v.customer_id IS NULL OR v.customer_id = ?) AND v.type !=3 " +
		whereMembership + whereType +
		orderBy +
		"LIMIT ?, ?"
	if count, err = db.Raw(q, values).QueryRows(&vouchers); err != nil || count == 0 {
		return nil, 0, err
	}

	voucherIdArr := []int64{}
	qMark := ""
	for _, v := range vouchers {
		qMark += "?,"
		voucherIdArr = append(voucherIdArr, v.ID)
	}
	qMark = qMark[:len(qMark)-1]

	// get usage of voucher based on voucher id got from previous query
	q = "SELECT vl.voucher_id, COUNT(vl.voucher_id) AS total_user_used " +
		"FROM voucher_log vl " +
		"WHERE vl.status = 1 AND vl.customer_id = ? AND vl.voucher_id IN (" + qMark + ") " +
		"GROUP BY vl.voucher_id"
	db.Raw(q, req.CustomerID, voucherIdArr).QueryRows(&voucherCount)

	// set up a map to hold how many voucher did the user had used
	voucherMap := make(map[int64]int64)
	for _, v := range voucherCount {
		voucherMap[v.VoucherId] = v.TotalUserUsed
	}

	// start remove voucher that has been fully used
	i := 0
	for i < len(vouchers) {
		v := vouchers[i]
		if _, isExist := voucherMap[v.ID]; isExist {
			v.TotalUserUsed = voucherMap[v.ID]
		}

		v.RemUserQuota = v.UserQuota - v.TotalUserUsed

		// if remaining quota has been reached, remove voucher data from array
		if v.RemUserQuota <= 0 {
			vouchers = append(vouchers[:i], vouchers[i+1:]...)
			continue
		}

		i++
		count--
	}

	return
}

func (r *VoucherRepository) GetMobileVoucherDetail(ctx context.Context, req *dto.VoucherRequestGetMobileVoucherDetail) (voucher *model.Voucher, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "VoucherRepository.GetMobileVoucherDetail")
	defer span.End()

	db := r.opt.Database.Read

	voucher = &model.Voucher{}
	var cols []string

	if req.RedeemCode != "" {
		voucher.RedeemCode = req.RedeemCode
		cols = append(cols, "redeem_code")
	}
	if req.Status != 0 {
		voucher.Status = req.Status
		cols = append(cols, "status")
	}
	if req.Code != "" {
		voucher.Code = req.Code
		cols = append(cols, "code")
	}
	err = db.ReadWithCtx(ctx, voucher, cols...)
	if err != nil {
		span.RecordError(err)
		return
	}

	qs := db.QueryTable(new(model.VoucherLog))

	cond := orm.NewCondition()

	cond = cond.And("voucher_id", voucher.ID).And("status", statusx.ConvertStatusName("Active")).And("customer_id", req.CustomerID)

	voucher.TotalUserUsed, err = qs.SetCond(cond).Count()
	if err != nil {
		span.RecordError(err)
		return
	}

	voucher.RemUserQuota = voucher.UserQuota - voucher.TotalUserUsed

	return voucher, nil
}
