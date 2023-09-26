package repository

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type ICustomerRepository interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, customerTypeId int64) (customer []*model.Customer, count int64, err error)
	GetDetail(ctx context.Context, id int64, code string, phoneNumber string) (address *model.Customer, err error)
}

type CustomerRepository struct {
	opt opt.Options
}

func NewCustomerRepository() ICustomerRepository {
	return &CustomerRepository{
		opt: global.Setup.Common,
	}
}

func (r *CustomerRepository) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, customerTypeId int64) (customer []*model.Customer, count int64, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerRepository.Get")
	defer span.End()

	dummies := r.MockDatas(10)
	return dummies, int64(len(dummies)), nil

	// db := r.opt.Database.Read

	// qs := db.QueryTable(new(model.Address))

	// cond := orm.NewCondition()

	// if search != "" {
	// 	condGroup := orm.NewCondition()
	// 	condGroup = condGroup.And("code", search)
	// 	cond = cond.AndCond(condGroup)
	// }

	// if status != 0 {
	// 	cond = cond.And("status", status)
	// }

	// if archetypeID != 0 {
	// 	cond = cond.And("archetype_id", archetypeID)
	// }

	// if admDivisionID != 0 {
	// 	cond = cond.And("admDivision_id", admDivisionID)
	// }

	// if siteID != 0 {
	// 	cond = cond.And("site_id", siteID)
	// }

	// if salespersonID != 0 {
	// 	cond = cond.And("salesperson_id", salespersonID)
	// }

	// if territoryID != 0 {
	// 	cond = cond.And("territory_id", territoryID)
	// }

	// if taxScheduleID != 0 {
	// 	cond = cond.And("tax_schedule_id", taxScheduleID)
	// }

	// qs = qs.SetCond(cond)

	// if orderBy != "" {
	// 	qs = qs.OrderBy(orderBy)
	// }

	// count, err = qs.Offset(offset).Limit(limit).AllWithCtx(ctx, &addresses)
	// if err != nil {
	// 	span.RecordError(err)
	// 	return
	// }

	return
}

func (r *CustomerRepository) GetDetail(ctx context.Context, id int64, code string, phoneNumber string) (customer *model.Customer, err error) {
	ctx, span := r.opt.Trace.Start(ctx, "CustomerRepository.GetDetail")
	defer span.End()

	return r.MockDatas(1)[0], nil

	// address = &model.Address{}

	// var cols []string

	// if id != 0 {
	// 	cols = append(cols, "id")
	// 	address.ID = id
	// }

	// if code != "" {
	// 	cols = append(cols, "code")
	// 	address.Code = code
	// }

	// if phoneNumber != "" {
	// 	cols = append(cols, "phone_number")
	// 	address.Phone1 = phoneNumber
	// }

	// db := r.opt.Database.Read
	// err = db.ReadWithCtx(ctx, address, cols...)
	// if err != nil {
	// 	span.RecordError(err)
	// 	return
	// }
}

func (r *CustomerRepository) MockDatas(total int) (mockDatas []*model.Customer) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.Customer{
				ID:                         int64(i),
				Code:                       fmt.Sprintf("CST%d", i),
				ReferralCode:               fmt.Sprintf("DUMMY REF CODE%d", i),
				Gender:                     1,
				Name:                       fmt.Sprintf("Dummy Customer %d", i),
				BirthDate:                  time.Now(),
				PicName:                    "Dummy PIC Name",
				PhoneNumber:                "85624464604",
				AltPhoneNumber:             "85624464604",
				Email:                      "dummy.faisal@edenfarm.id",
				Password:                   "$2a$10$jciWRgz5gMxTLnoy.7ZIJOBpn12WFr7mv8oAQ7YcJSZk9qqO2Uq/G",
				BillingAddress:             "Dummy BillingAddress",
				Note:                       "Dummy Note",
				ReferenceInfo:              "Dummy ReferenceInfo",
				TagCustomer:                "Dummy TagCustomer",
				Status:                     1,
				Suspended:                  2,
				UpgradeStatus:              0,
				CustomerGroup:              123,
				TagCustomerName:            "ABC Group",
				ReferrerCode:               "",
				CreatedAt:                  time.Now(),
				CreatedBy:                  1,
				LastUpdatedAt:              time.Now(),
				LastUpdatedBy:              1,
				TotalPoint:                 0,
				CustomerTypeCreditLimit:    1,
				EarnedPoint:                0,
				RedeemedPoint:              0,
				CustomCreditLimit:          0,
				CreditLimitAmount:          0,
				ProfileCode:                "",
				RemainingCreditLimitAmount: 0,
				AverageSales:               0,
				RemainingOutstanding:       0,
				OverdueDebt:                0,
				MerchantPhotosUrl:          "",
				KTPPhotosUrl:               "",
				KTPPhotosUrlArr:            nil,
				MerchantPhotosUrlArr:       nil,
				MembershipLevelID:          0,
				MembershipCheckpointID:     1,
				MembershipRewardID:         1,
				MembershipRewardAmount:     0,
				BirthDateString:            time.Now().Format("2006-01-02"),
				CustomerTypeID:             "BTY0009",
				SalesPaymentTermID:         1,
			})
	}
	mockDatas = append(mockDatas,
		&model.Customer{
			ID:                         int64(total + 1),
			Code:                       fmt.Sprintf("CST%d", total+1),
			ReferralCode:               fmt.Sprintf("DUMMY REF CODE%d", total+1),
			Gender:                     1,
			Name:                       fmt.Sprintf("Dummy Customer %d", total+1),
			BirthDate:                  time.Now(),
			PicName:                    "Dummy PIC Name",
			PhoneNumber:                "8000011113",
			AltPhoneNumber:             "8000011113",
			Email:                      "dummy.edenfarm@edenfarm.id",
			Password:                   "$2a$10$jciWRgz5gMxTLnoy.7ZIJOBpn12WFr7mv8oAQ7YcJSZk9qqO2Uq/G",
			BillingAddress:             "Dummy BillingAddress",
			Note:                       "Dummy Note",
			ReferenceInfo:              "Dummy ReferenceInfo",
			TagCustomer:                "Dummy TagCustomer",
			Status:                     1,
			Suspended:                  2,
			UpgradeStatus:              0,
			CustomerGroup:              123,
			TagCustomerName:            "ABC Group",
			ReferrerCode:               "",
			CreatedAt:                  time.Now(),
			CreatedBy:                  1,
			LastUpdatedAt:              time.Now(),
			LastUpdatedBy:              1,
			TotalPoint:                 0,
			CustomerTypeCreditLimit:    1,
			EarnedPoint:                0,
			RedeemedPoint:              0,
			CustomCreditLimit:          0,
			CreditLimitAmount:          0,
			ProfileCode:                "",
			RemainingCreditLimitAmount: 0,
			AverageSales:               0,
			RemainingOutstanding:       0,
			OverdueDebt:                0,
			MerchantPhotosUrl:          "",
			KTPPhotosUrl:               "",
			KTPPhotosUrlArr:            nil,
			MerchantPhotosUrlArr:       nil,
			MembershipLevelID:          0,
			MembershipCheckpointID:     1,
			MembershipRewardID:         1,
			MembershipRewardAmount:     0,
			BirthDateString:            time.Now().Format("2006-01-02"),
			CustomerTypeID:             "BTY0009",
			SalesPaymentTermID:         1,
		})
	return
}
