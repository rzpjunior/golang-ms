package cronjob

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/cuxs"
	"git.edenfarm.id/cuxs/dbredis"
	"git.edenfarm.id/cuxs/mongodb"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/log"
	"git.edenfarm.id/project-version2/api/service/talon"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
	. "github.com/ahmetb/go-linq/v3"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

var GetAreaPolicy []areaPolicy

func (h *Handler) URLMapping(r *echo.Group) {
	r.GET("", h.remindPayment)
	r.GET("/price/scheduler", h.priceScheduler)
	r.GET("/loyalty", h.LoyaltyScheduler)
	r.GET("/merchant/suspension", h.suspendMerchantScheduler)
	r.GET("/notification_campaign", h.notificationCampaign)
	r.GET("/sales/assignment/auto_failed", h.autoFailedSalesAssignmentScheduler)
	r.GET("/summary-notification-campaign", h.summaryNotificationCampaign)
	r.GET("/voucher/anniversary-join", h.createAnniversaryJoinVoucher)
	r.GET("/expired-eden-point", h.ExpiredEdenPoint)
}

type mplQuery struct {
	SalesOrderId      int64   `orm:"column(sales_order_id)"`
	DeliveryReturnId  int64   `orm:"column(delivery_return_id);null"`
	MerchantId        int64   `orm:"column(merchant_id)"`
	TotalPoint        float64 `orm:"column(total_point)"`
	RecentPoint       float64 `orm:"column(recent_point)"`
	TalonCampaignID   int64   `orm:"-"`
	TalonCampaignName string  `orm:"-"`
	TalonMultiplier   int8    `orm:"-"`
	TransactionType   int8    `orm:"-"`
	RefereeID         int64   `orm:"-"`
	ReferrerID        int64   `orm:"-"`
}

type soQuery struct {
	SalesOrderId          int64   `orm:"column(sales_order_id)"`
	DeliveryReturnId      int64   `orm:"column(delivery_return_id);null"`
	AmountPayment         float64 `orm:"column(amount_payment)"`
	MerchantId            int64   `orm:"column(merchant_id)"`
	FirebaseToken         string  `orm:"column(firebase_token)"`
	TotalPoint            float64 `orm:"column(total_point)"`
	IntegrationCode       string  `orm:"column(integration_code)"`
	ProfileCode           string  `orm:"column(profile_code)"`
	Archetype             string  `orm:"column(archetype_name)"`
	PriceSet              string  `orm:"column(price_set_name)"`
	RedeemAmount          float64 `orm:"column(point_redeem_amount)"`
	ReferrerCode          string  `orm:"column(referrer_code)"`
	SalesInvoiceID        int64   `orm:"column(sales_invoice_id)"`
	VoucherAmount         float64 `orm:"column(vou_disc_amount)"`
	VoucherType           int8    `orm:"column(voucher_type)"`
	VoucherID             int64   `orm:"column(voucher_id)"`
	OrderTypeID           int64   `orm:"column(order_type_sls_id)"`
	Code                  string  `orm:"-"`
	MembershipLevel       string  `orm:"-"`
	MembershipCheckpoint  string  `orm:"-"`
	MembershipRewardLevel string  `orm:"="`
}

func (h *Handler) LoyaltyScheduler(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	o := orm.NewOrm()
	o.Using("read_only")
	ormWrite := orm.NewOrm()

	// dateDiff : to hold date differences parameter
	dateDiff := "2"
	if ctx.QueryParam("date_diff") != "" {
		dateDiff = ctx.QueryParam("date_diff")
	}

	// customDate : to set custom current date, be used in testing
	customDate := ""
	if ctx.QueryParam("custom_date") != "" {
		customDate = ctx.QueryParam("custom_date")
	}

	currentDate := time.Now().Format("2006-01-02")

	type merchantPoint struct {
		MerchantId  int64
		RecentPoint float64
	}
	type calcReturn struct {
		Amount        float64 `orm:"column(amount)"`
		SalesOrderId  int64   `orm:"column(sales_order_id)"`
		DeliveryFee   float64 `orm:"column(delivery_fee)"`
		VouDiscAmount float64 `orm:"column(vou_disc_amount)"`
		VoucherType   int64   `orm:"column(voucher_type);null"`
	}
	type merchantData struct {
		TotalPoint   float64
		LevelID      int64
		CheckpointID int64
		RewardID     int64
		RewardAmount float64
	}
	var (
		tempMerchant                                                                                     []merchantPoint
		salesOrderQuery                                                                                  []soQuery
		salesOrder                                                                                       []soQuery
		tempSalesOrder                                                                                   []soQuery
		tempMplQuery                                                                                     []mplQuery
		salesOrderDR                                                                                     []soQuery
		salesInvoice                                                                                     *model.SalesInvoice
		merchant                                                                                         *model.Merchant
		orderType                                                                                        *model.OrderType
		csr                                                                                              *model.CustomerSessionReturn
		paymentMethod                                                                                    string
		merchantDatas                                                                                    = make(map[int64]*merchantData)
		redeemedData                                                                                     []map[string]interface{}
		lrecentPoint                                                                                     = make(map[int64]float64)
		vouchers                                                                                         = make(map[int64]map[string]*model.Voucher)
		expiredMonth, quartalPhase, currentDay, currentYear, currentMonthExpiration, nextMonthExpiration int
		currentMonth                                                                                     time.Month
		isExpiredNextPeriod                                                                              bool
	)

	_, e = o.Raw("select "+
		"so.id as sales_order_id "+
		", so.code "+
		", dr.id as delivery_return_id "+
		", SUM(sp.amount) as amount_payment "+
		", um.firebase_token firebase_token "+
		", m.id as merchant_id  "+
		", m.total_point "+
		", so.eden_point_campaign_id "+
		", so.integration_code "+
		", if (m.profile_code = '', m.code, m.profile_code) profile_code "+
		", a.name archetype_name "+
		", ps.name price_set_name "+
		", so.point_redeem_amount "+
		", m.referrer_code "+
		", v.type voucher_type "+
		", so.vou_disc_amount "+
		", si.id sales_invoice_id "+
		", v.id voucher_id "+
		", so.order_type_sls_id "+
		"from sales_order so "+
		"JOIN delivery_order do ON do.sales_order_id = so.id and do.status = 2 "+
		"LEFT JOIN delivery_return dr ON dr.delivery_order_id = do.id and dr.status = 2 "+
		"JOIN sales_invoice si ON si.sales_order_id = so.id  "+
		"JOIN sales_payment sp ON sp.sales_invoice_id = si.id "+
		"JOIN branch b ON b.id = so.branch_id "+
		"JOIN merchant m ON m.id = b.merchant_id "+
		"JOIN archetype a ON b.archetype_id = a.id "+
		"JOIN price_set ps on b.price_set_id = ps.id "+
		"INNER JOIN user_merchant um ON um.id = m.user_merchant_id "+
		"LEFT JOIN voucher v ON so.voucher_id = v.id "+
		"JOIN config_app ca ON ca.`attribute` = 'edenpoint_business_type' and find_in_set(m.business_type_id, ca.value) > 0 "+
		"JOIN config_app ca2 ON ca2.`attribute` = 'edenpoint_order_channel' and find_in_set(so.order_channel, ca2.value) > 0 "+
		"WHERE so.status = 2 "+
		"and sp.status = 2 "+
		"and sp.amount != 0 "+
		"and DATE_ADD(CURRENT_DATE, INTERVAL -"+dateDiff+" DAY) = DATE_FORMAT(so.finished_at, ?) "+
		"group by sales_order_id order by sales_order_id desc;", "%Y-%m-%d").QueryRows(&salesOrderQuery)

	if e != nil {
		return e
	}

	//getsales order first
	for _, item := range salesOrderQuery {
		//cek ada di mpl ato engga so tersebut
		//check if exist in merchant point log
		filter := map[string]interface{}{
			"merchant_id":    item.MerchantId,
			"sales_order_id": item.SalesOrderId,
			"status":         int8(1),
		}
		exclude := map[string]interface{}{}
		if _, isExistMerchantPointLog, err := repository.CheckMerchantPointLogData(filter, exclude); err == nil && isExistMerchantPointLog >= 1 {
			continue
		}

		//cek ada delivery return ato engga
		//kalo ada continue
		if item.DeliveryReturnId != 0 {
			salesOrderDR = append(salesOrderDR, item)
			cntMerchantId := From(tempMerchant).Where(
				func(f interface{}) bool { return f.(merchantPoint).MerchantId == item.MerchantId },
			).Count()
			if cntMerchantId == 0 {
				tempMerchantLocal := merchantPoint{MerchantId: item.MerchantId, RecentPoint: item.TotalPoint}
				tempMerchant = append(tempMerchant, tempMerchantLocal)
			}

		} else { //kalo ga ada delivery return
			salesOrder = append(salesOrder, item)
			cntMerchantId := From(tempMerchant).Where(
				func(f interface{}) bool { return f.(merchantPoint).MerchantId == item.MerchantId },
			).Count()
			if cntMerchantId == 0 {
				tempMerchantLocal := merchantPoint{MerchantId: item.MerchantId, RecentPoint: item.TotalPoint}
				tempMerchant = append(tempMerchant, tempMerchantLocal)
			}
		}

		if item.RedeemAmount > 0 {
			redeemData := map[string]interface{}{
				"merchant_profile": item.ProfileCode,
				"redeemed_amount":  item.RedeemAmount,
			}

			redeemedData = append(redeemedData, redeemData)
		}
	}

	//loop the sales order that doesnt have delivery return
	//and add it to temporary variable for inserting to database
	for _, item := range salesOrder {
		var (
			itemList     []*model.SessionItemData
			referrerData []string
		)

		if val, isExist := lrecentPoint[item.MerchantId]; !isExist || (isExist && val == 0) {
			currentPointReferrer := 0.00
			e = o.Raw("select total_point from merchant where id = ?", item.MerchantId).QueryRow(&currentPointReferrer)
			lrecentPoint[item.MerchantId] = currentPointReferrer
		}

		//get calculate earning point
		var CalculateReturn []calcReturn

		_, e := o.Raw("SELECT "+
			"SUM(sp.amount) amount "+
			", si.sales_order_id "+
			", so.delivery_fee "+
			", so.vou_disc_amount "+
			", v.`type` as voucher_type "+
			"from sales_payment sp "+
			"JOIN sales_invoice si ON si.id = sp.sales_invoice_id "+
			"JOIN sales_order so ON so.id = si.sales_order_id "+
			"LEFT JOIN voucher v ON v.redeem_code = so.vou_redeem_code "+
			"WHERE so.id = ? "+ // -- Impacted Sales Order
			"and sp.status = 2 "+
			"group by sales_order_id ; ", item.SalesOrderId).QueryRows(&CalculateReturn)

		if e != nil {
			return e
		}

		// start update customer session in talon
		salesInvoice, _ = repository.GetSalesInvoice("id", item.SalesInvoiceID)
		for _, v := range salesInvoice.SalesInvoiceItems {
			var (
				itemData                    *model.SessionItemData
				parentName, grandparentName string
			)

			v.Product.Read("ID")
			v.Product.Category.Read("ID")
			v.SalesOrderItem.Read("ID")
			if v.Product.Category.ParentID != 0 {
				v.Product.Category.Parent = &model.Category{ID: v.Product.Category.ParentID}
				v.Product.Category.Parent.Read("ID")
				parentName = v.Product.Category.Parent.Name
			}
			if v.Product.Category.GrandParentID != 0 {
				v.Product.Category.GrandParent = &model.Category{ID: v.Product.Category.GrandParentID}
				v.Product.Category.GrandParent.Read("ID")
				grandparentName = v.Product.Category.GrandParent.Name
			}

			itemData = &model.SessionItemData{
				ProductName:  v.Product.Name,
				ProductCode:  v.Product.Code,
				CategoryName: v.Product.Category.Name,
				UnitPrice:    v.UnitPrice - v.SalesOrderItem.UnitPriceDiscount,
				OrderQty:     1,
				UnitWeight:   v.InvoiceQty,
				Attributes: map[string]string{
					"parent_category":       parentName,
					"grand_parent_category": grandparentName,
				},
			}
			itemList = append(itemList, itemData)
		}

		orderType = &model.OrderType{ID: item.OrderTypeID}
		if e = orderType.Read("ID"); e != nil {
			continue
		}

		merchant = &model.Merchant{ID: item.MerchantId}
		if e = merchant.Read("ID"); e != nil {
			continue
		}
		if e = merchant.FinanceArea.Read("ID"); e != nil {
			continue
		}
		if e = merchant.BusinessType.Read("ID"); e != nil {
			continue
		}
		// start set if merchant has referral
		if merchant.Referrer != nil && merchant.Referrer.ID != 0 {
			if e = merchant.Referrer.Read("ID"); e != nil {
				continue
			}
			if e = merchant.Referrer.FinanceArea.Read("ID"); e != nil {
				continue
			}
			if e = merchant.Referrer.BusinessType.Read("ID"); e != nil {
				continue
			}
			paymentMethod = ""
			if merchant.Referrer.PaymentMethod != nil {
				if e = merchant.Referrer.PaymentMethod.Read("ID"); e != nil {
					continue
				}
				paymentMethod = merchant.Referrer.PaymentMethod.Name
			}

			if merchant.Referrer.ProfileCode == "" {
				merchant.Referrer.ProfileCode = merchant.Referrer.Code
			}
			e = talon.UpdateCustomerProfileTalon(merchant.Referrer.ProfileCode, merchant.Referrer.TagCustomer, merchant.Referrer.FinanceArea.Name, merchant.Referrer.BusinessType.Name, paymentMethod, merchant.Referrer.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
			_, e = ormWrite.Raw("update merchant set profile_code = ? where id = ?", merchant.Referrer.ProfileCode, merchant.Referrer.ID).Exec()
			referrerData = []string{merchant.Referrer.ProfileCode, merchant.Referrer.ReferralCode}
		}
		// end set if merchant has referral
		paymentMethod = ""
		if merchant.PaymentMethod != nil {
			if e = merchant.PaymentMethod.Read("ID"); e != nil {
				continue
			}
			paymentMethod = merchant.PaymentMethod.Name
		}

		if merchant.ProfileCode == "" {
			merchant.ProfileCode = merchant.Code
		}
		e = talon.UpdateCustomerProfileTalon(merchant.ProfileCode, merchant.TagCustomer, merchant.FinanceArea.Name, merchant.BusinessType.Name, paymentMethod, merchant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), referrerData...)
		_, e = ormWrite.Raw("update merchant set profile_code = ? where id = ?", merchant.ProfileCode, merchant.ID).Exec()

		voucherAmount := 0.00
		if item.VoucherAmount > 0 && (item.VoucherType == 1 || item.VoucherType == 2) {
			voucherAmount = item.VoucherAmount
		}

		if item.IntegrationCode == "" {
			item.IntegrationCode = strings.ReplaceAll(time.Now().Format("20060102150405.99"), ".", "") + item.ProfileCode
			_, e = ormWrite.Raw("update sales_order set integration_code = ? where id = ?", item.IntegrationCode, item.SalesOrderId).Exec()
		}
		csr, e = talon.UpdateCustomerSessionTalon("closed", "false", item.IntegrationCode, item.ProfileCode, item.Archetype, item.PriceSet, item.ReferrerCode, itemList, item.RedeemAmount > 0, voucherAmount, orderType.Name)
		// end update customer session in talon

		// start loop through triggered efffects
		for _, v := range csr.Effects {
			if v.EffectType == "addLoyaltyPoints" && v.Props.SubLedgerID == "" {
				var (
					tempMpl     mplQuery
					tempSOQuery soQuery
				)

				tempMpl, tempSOQuery, lrecentPoint, e = setUpLoyaltyPoints(csr.CustomerSession.ProfileCode, v, item, lrecentPoint)

				tempMplQuery = append(tempMplQuery, tempMpl)
				tempSalesOrder = append(tempSalesOrder, tempSOQuery)
			}
		}
		// end loop through triggered efffects

		// start update membership level & checkpoint of merchant
		var (
			membershipName, membershipCheckpointStr, q                string
			membershipLevel                                           = new(model.MembershipLevel)
			membershipCheckpoint                                      = new(model.MembershipCheckpoint)
			membershipRewards                                         []*model.MembershipReward
			membershipRewardID                                        int64
			membershipRewardAmount, totalFreshAmount, rewardMinAmount float64
			rewardMinLevel, membershipRewardLevel                     int8
		)

		attributes := reflect.ValueOf(csr.CustomerProfile.Attributes)
		for _, v := range attributes.MapKeys() {
			if v.String() == "membership_level" {
				membershipName = attributes.MapIndex(v).Interface().(string)
				continue
			}

			if v.String() == "membership_checkpoint" {
				membershipCheckpointStr = attributes.MapIndex(v).Interface().(string)
				continue
			}

			if v.String() == "fresh_product_revenue" {
				totalFreshAmount = attributes.MapIndex(v).Interface().(float64)
				continue
			}
		}

		// get membership level data
		if membershipName != "" {
			q = "select * from membership_level ml where ml.status = 1 and ml.name = ?"
			if e = o.Raw(q, membershipName).QueryRow(&membershipLevel); e != nil {
				continue
			}
		}

		// get membership checkpoint data
		if membershipCheckpointStr != "" {
			q = "select * from membership_checkpoint mc where mc.status = 1 and mc.checkpoint = ?"
			if e = o.Raw(q, membershipCheckpointStr).QueryRow(&membershipCheckpoint); e != nil {
				continue
			}
		}

		// check whether merchant deserve to get reward
		if membershipLevel.ID != 0 {
			q = "select ca.value from config_app ca where ca.attribute = 'minimum_level_membership_reward'"
			e = o.Raw(q).QueryRow(&rewardMinLevel)
		}

		if membershipLevel.Level >= rewardMinLevel {
			membershipRewardID = merchant.MembershipRewardID
			merchant.MembershipReward = &model.MembershipReward{ID: membershipRewardID}
			e = merchant.MembershipReward.Read("ID")

			// get membership reward data
			q = "select mr.* from membership_reward mr where mr.status = 1 and mr.reward_level >= ? order by mr.reward_level asc"
			_, e = o.Raw(q, merchant.MembershipReward.RewardLevel).QueryRows(&membershipRewards)

			// get minimum amount to get reward
			q = "select max(target_amount) " +
				"from membership_level ml " +
				"join membership_checkpoint mc on ml.id = mc.membership_level_id " +
				"join config_app ca on ca.`attribute` = 'minimum_level_membership_reward' and ml.`level` < ca.value"
			if e = o.Raw(q).QueryRow(&rewardMinAmount); e != nil {
				continue
			}

			membershipRewardAmount = totalFreshAmount - rewardMinAmount

			var voucherSnk string
			q = "select value from config_app where attribute = ?"
			o.Raw(q, "membership_reward_voucher_snk").QueryRow(&voucherSnk)

			voucherMap := make(map[string]*model.Voucher)
			for _, v := range membershipRewards {
				membershipRewardID = v.ID
				if membershipRewardAmount < v.MaxAmount {
					break
				}
				membershipRewardLevel = v.RewardLevel

				var (
					voucherImage  string
					voucherQuota  int64
					voucherAmount float64
				)

				// create voucher object for each reward level
				if v.RewardLevel <= 3 {
					loc, _ := time.LoadLocation("Asia/Jakarta")
					currTime := time.Now().In(loc)
					currYear := currTime.Year()
					currDateStr := currTime.Format("2006-01-02 00:00:00")
					currTime, _ = time.Parse("2006-01-02 15:04:05", currDateStr)
					currTime = currTime.In(loc)

					o.Raw(q, "membership_reward_voucher_quota_"+strconv.Itoa(int(v.RewardLevel))).QueryRow(&voucherQuota)
					o.Raw(q, "membership_reward_voucher_amount_"+strconv.Itoa(int(v.RewardLevel))).QueryRow(&voucherAmount)
					o.Raw(q, "membership_reward_voucher_image_"+strconv.Itoa(int(v.RewardLevel))).QueryRow(&voucherImage)

					voucher := &model.Voucher{
						Area:            &model.Area{ID: 1},
						Archetype:       &model.Archetype{ID: 22},
						RedeemCode:      "KJTNDSKN" + currTime.Format("2006"),
						Type:            1,
						Name:            "Voucher Belanja " + strconv.Itoa(int(voucherAmount)),
						StartTimestamp:  currTime,
						EndTimestamp:    time.Date(currYear, 12, 31, 23, 59, 59, 0, loc),
						OverallQuota:    voucherQuota,
						RemOverallQuota: voucherQuota,
						UserQuota:       voucherQuota,
						DiscAmount:      voucherAmount,
						Status:          1,
						ChannelVoucher:  "2,3",
						MerchantID:      merchant.ID,
					}

					voucher.VoucherContent = &model.VoucherContent{
						Voucher:        voucher,
						ImageUrl:       voucherImage,
						TermConditions: voucherSnk,
					}

					voucherMap[voucher.Name] = voucher
				}
			}
			vouchers[merchant.ID] = voucherMap

			if membershipRewardLevel > 0 {
				tempSOQuery := soQuery{
					SalesOrderId:          item.SalesOrderId,
					MerchantId:            item.MerchantId,
					FirebaseToken:         item.FirebaseToken,
					Code:                  "NOT0031",
					MembershipRewardLevel: strconv.Itoa(int(membershipRewardLevel)),
				}
				tempSalesOrder = append(tempSalesOrder, tempSOQuery)
			}
		}

		// set merchant membership related data
		totalPoint := float64(0)
		if _, isExist := merchantDatas[merchant.ID]; isExist {
			totalPoint += merchantDatas[merchant.ID].TotalPoint
		}
		merchantDatas[merchant.ID] = &merchantData{
			LevelID:      membershipLevel.ID,
			CheckpointID: membershipCheckpoint.ID,
			RewardID:     membershipRewardID,
			RewardAmount: membershipRewardAmount,
			TotalPoint:   totalPoint,
		}

		merchant.MembershipLevel = &model.MembershipLevel{ID: merchant.MembershipLevelID}
		merchant.MembershipLevel.Read("ID")
		if merchant.MembershipLevel.Level < membershipLevel.Level && membershipLevel.Level > 0 {
			code := "NOT0032"
			if membershipLevel.Level == 3 {
				code = "NOT0033"
			}

			tempSOQuery := soQuery{
				SalesOrderId:         item.SalesOrderId,
				MerchantId:           item.MerchantId,
				FirebaseToken:        item.FirebaseToken,
				Code:                 code,
				MembershipLevel:      membershipLevel.Name,
				MembershipCheckpoint: strconv.Itoa(int(membershipCheckpoint.Checkpoint)),
			}
			tempSalesOrder = append(tempSalesOrder, tempSOQuery)
		} else {
			merchant.MembershipCheckpoint = &model.MembershipCheckpoint{ID: merchant.MembershipCheckpointID}
			merchant.MembershipCheckpoint.Read("ID")
			if merchant.MembershipCheckpoint.Checkpoint < membershipCheckpoint.Checkpoint && membershipCheckpoint.Checkpoint > 0 {
				tempSOQuery := soQuery{
					SalesOrderId:         item.SalesOrderId,
					MerchantId:           item.MerchantId,
					FirebaseToken:        item.FirebaseToken,
					Code:                 "NOT0030",
					MembershipLevel:      membershipLevel.Name,
					MembershipCheckpoint: strconv.Itoa(int(membershipCheckpoint.Checkpoint)),
				}
				tempSalesOrder = append(tempSalesOrder, tempSOQuery)
			}
		}
		// end update membership level & checkpoint of merchant

		// add extra edenpoint from voucher
		if item.VoucherType == 4 {
			voucher := &model.Voucher{ID: item.VoucherID}
			if e = voucher.Read("ID"); e != nil {
				continue
			}
			lrecentPoint[item.MerchantId] += voucher.DiscAmount
			tempMpl := mplQuery{
				SalesOrderId:     item.SalesOrderId,
				DeliveryReturnId: item.DeliveryReturnId,
				MerchantId:       item.MerchantId,
				TotalPoint:       voucher.DiscAmount,
				RecentPoint:      lrecentPoint[item.MerchantId],
				TransactionType:  8,
			}
			tempMplQuery = append(tempMplQuery, tempMpl)
		}
	}

	for _, item := range salesOrderDR {
		var (
			itemList     []*model.SessionItemData
			referrerData []string
		)

		if val, isExist := lrecentPoint[item.MerchantId]; !isExist || (isExist && val == 0) {
			currentPointReferrer := 0.00
			e = o.Raw("select total_point from merchant where id = ?", item.MerchantId).QueryRow(&currentPointReferrer)
			lrecentPoint[item.MerchantId] = currentPointReferrer
		}

		//recalculate amount payment total
		var DeliveryReturnCalculate []struct {
			ProductID      int64   `orm:"column(product_id)"`
			SalesOrderID   int64   `orm:"column(sales_order_id)"`
			ReceiveQty     float64 `orm:"column(receive_qty)"`
			DeliverQty     float64 `orm:"column(deliver_qty)"`
			ReturnGoodQty  float64 `orm:"column(return_good_qty)"`
			ReturnWasteQty float64 `orm:"column(return_waste_qty)"`
			InvoiceQty     float64 `orm:"column(invoice_qty)"`
			UnitPrice      float64 `orm:"column(unit_price)"`
		}
		//get delivery return item
		_, e = o.Raw("SELECT "+
			"dri.product_id "+
			", so.id as sales_order_id "+
			", doi.receive_qty "+
			", doi.deliver_qty "+
			", dri.return_good_qty "+ // return qty (return_good_qty + return_waste_qty)
			", dri.return_waste_qty "+ // return qty (return_good_qty + return_waste_qty)
			", sii.invoice_qty "+
			", sii.unit_price "+
			"from delivery_return_item dri "+
			"JOIN delivery_return dr ON dr.id = dri.delivery_return_id "+
			"JOIN delivery_order do ON do.id = dr.delivery_order_id "+
			"JOIN delivery_order_item doi ON doi.delivery_order_id = do.id and doi.product_id = dri.product_id "+
			"JOIN sales_order so ON so.id = do.sales_order_id "+
			"JOIN sales_invoice si ON si.sales_order_id = so.id "+
			"JOIN sales_invoice_item sii ON sii.sales_invoice_id = si.id and sii.product_id = dri.product_id "+
			"where dri.delivery_return_id = ?;", item.DeliveryReturnId).QueryRows(&DeliveryReturnCalculate)

		if e != nil {
			return e
		}
		totalAmountReturn := 0.0
		for _, itemDR := range DeliveryReturnCalculate {
			if itemDR.InvoiceQty == 0 {
				continue
			} else {
				returnQty := itemDR.ReturnGoodQty + itemDR.ReturnWasteQty
				if itemDR.InvoiceQty == itemDR.DeliverQty-(returnQty) {
					continue
				} else {
					totalAmountReturn = totalAmountReturn + (returnQty * itemDR.UnitPrice)
				}
			}
		}
		//get calculate earning point
		var CalculateReturn []calcReturn

		_, e := o.Raw("SELECT "+
			"SUM(sp.amount) amount "+
			", si.sales_order_id "+
			", so.delivery_fee "+
			", so.vou_disc_amount "+
			", v.`type` as voucher_type "+
			"from sales_payment sp "+
			"JOIN sales_invoice si ON si.id = sp.sales_invoice_id "+
			"JOIN sales_order so ON so.id = si.sales_order_id "+
			"LEFT JOIN voucher v ON v.redeem_code = so.vou_redeem_code "+
			"WHERE so.id = ? "+ // -- Impacted Sales Order
			"and sp.status = 2 "+
			"group by sales_order_id ; ", item.SalesOrderId).QueryRows(&CalculateReturn)

		if e != nil {
			return e
		}

		// start update customer session in talon
		salesInvoice, _ = repository.GetSalesInvoice("id", item.SalesInvoiceID)
		for _, v := range salesInvoice.SalesInvoiceItems {
			var (
				itemData                    *model.SessionItemData
				parentName, grandparentName string
			)

			v.Product.Read("ID")
			v.Product.Category.Read("ID")
			v.SalesOrderItem.Read("ID")
			if v.Product.Category.ParentID != 0 {
				v.Product.Category.Parent = &model.Category{ID: v.Product.Category.ParentID}
				v.Product.Category.Parent.Read("ID")
				parentName = v.Product.Category.Parent.Name
			}
			if v.Product.Category.GrandParentID != 0 {
				v.Product.Category.GrandParent = &model.Category{ID: v.Product.Category.GrandParentID}
				v.Product.Category.GrandParent.Read("ID")
				grandparentName = v.Product.Category.GrandParent.Name
			}

			type QtyData struct {
				InvoiceQty     float64 `orm:"column(invoice_qty)"`
				DeliveryQty    float64 `orm:"column(deliver_qty)"`
				ReturnWasteQty float64 `orm:"column(return_good_qty)"`
				ReturnGoodQty  float64 `orm:"column(return_waste_qty)"`
			}
			var qtyData *QtyData

			qty := v.InvoiceQty

			if e = o.Raw("select dri.return_good_qty, dri.return_waste_qty, doi.deliver_qty, sii.invoice_qty "+
				"from delivery_return_item dri "+
				"join delivery_return dr on dri.delivery_return_id = dr.id "+
				"join delivery_order_item doi on dri.delivery_order_item_id = doi.id "+
				"join sales_invoice_item sii on doi.sales_order_item_id = sii.sales_order_item_id "+
				"where doi.sales_order_item_id = ? and dri.product_id = ? and dr.status = 2", v.SalesOrderItem.ID, v.Product.ID).QueryRow(&qtyData); e == nil {

				if qtyData.InvoiceQty == qtyData.DeliveryQty-(qtyData.ReturnWasteQty+qtyData.ReturnGoodQty) {
					qty = qtyData.InvoiceQty
				} else {
					qty = qtyData.InvoiceQty - (qtyData.ReturnGoodQty + qtyData.ReturnWasteQty)
				}
			}
			if qty > 0 {
				itemData = &model.SessionItemData{
					ProductName:  v.Product.Name,
					ProductCode:  v.Product.Code,
					CategoryName: v.Product.Category.Name,
					UnitPrice:    v.UnitPrice - v.SalesOrderItem.UnitPriceDiscount,
					OrderQty:     1,
					UnitWeight:   qty,
					Attributes: map[string]string{
						"parent_category":       parentName,
						"grand_parent_category": grandparentName,
					},
				}
				itemList = append(itemList, itemData)
			}
		}

		orderType = &model.OrderType{ID: item.OrderTypeID}
		if e = orderType.Read("ID"); e != nil {
			continue
		}

		merchant = &model.Merchant{ID: item.MerchantId}
		if e = merchant.Read("ID"); e != nil {
			continue
		}
		if e = merchant.FinanceArea.Read("ID"); e != nil {
			continue
		}
		if e = merchant.BusinessType.Read("ID"); e != nil {
			continue
		}
		// start set if merchant has referral
		if merchant.Referrer != nil && merchant.Referrer.ID != 0 {
			if e = merchant.Referrer.Read("ID"); e != nil {
				continue
			}
			if e = merchant.Referrer.FinanceArea.Read("ID"); e != nil {
				continue
			}
			if e = merchant.Referrer.BusinessType.Read("ID"); e != nil {
				continue
			}
			paymentMethod = ""
			if merchant.Referrer.PaymentMethod != nil {
				if e = merchant.Referrer.PaymentMethod.Read("ID"); e != nil {
					continue
				}
				paymentMethod = merchant.Referrer.PaymentMethod.Name
			}

			if merchant.Referrer.ProfileCode == "" {
				merchant.Referrer.ProfileCode = merchant.Referrer.Code
			}
			e = talon.UpdateCustomerProfileTalon(merchant.Referrer.ProfileCode, merchant.Referrer.TagCustomer, merchant.Referrer.FinanceArea.Name, merchant.Referrer.BusinessType.Name, paymentMethod, merchant.Referrer.CreatedAt.Format("2006-01-02T15:04:05Z07:00"))
			_, e = ormWrite.Raw("update merchant set profile_code = ? where id = ?", merchant.Referrer.ProfileCode, merchant.Referrer.ID).Exec()
			referrerData = []string{merchant.Referrer.ProfileCode, merchant.Referrer.ReferralCode}
		}
		// end set if merchant has referral
		paymentMethod = ""
		if merchant.PaymentMethod != nil {
			if e = merchant.PaymentMethod.Read("ID"); e != nil {
				continue
			}
			paymentMethod = merchant.PaymentMethod.Name
		}

		if merchant.ProfileCode == "" {
			merchant.ProfileCode = merchant.Code
		}
		e = talon.UpdateCustomerProfileTalon(merchant.ProfileCode, merchant.TagCustomer, merchant.FinanceArea.Name, merchant.BusinessType.Name, paymentMethod, merchant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), referrerData...)
		_, e = ormWrite.Raw("update merchant set profile_code = ? where id = ?", merchant.ProfileCode, merchant.ID).Exec()

		voucherAmount := 0.00
		if item.VoucherAmount > 0 && (item.VoucherType == 1 || item.VoucherType == 2) {
			voucherAmount = item.VoucherAmount
		}

		if item.IntegrationCode == "" {
			item.IntegrationCode = strings.ReplaceAll(time.Now().Format("20060102150405.99"), ".", "") + item.ProfileCode
			_, e = ormWrite.Raw("update sales_order set integration_code = ? where id = ?", item.IntegrationCode, item.SalesOrderId).Exec()
		}
		csr, e = talon.UpdateCustomerSessionTalon("closed", "false", item.IntegrationCode, item.ProfileCode, item.Archetype, item.PriceSet, item.ReferrerCode, itemList, item.RedeemAmount > 0, voucherAmount, orderType.Name)
		// end update customer session in talon

		// start loop through triggered efffects
		for _, v := range csr.Effects {
			if v.EffectType == "addLoyaltyPoints" && v.Props.SubLedgerID == "" {
				var (
					tempMpl     mplQuery
					tempSOQuery soQuery
				)

				tempMpl, tempSOQuery, lrecentPoint, e = setUpLoyaltyPoints(csr.CustomerSession.ProfileCode, v, item, lrecentPoint)

				tempMplQuery = append(tempMplQuery, tempMpl)
				tempSalesOrder = append(tempSalesOrder, tempSOQuery)
			}
		}
		// end loop through triggered efffects

		// start update membership level & checkpoint of merchant
		var (
			membershipName, membershipCheckpointStr, q                string
			membershipLevel                                           = new(model.MembershipLevel)
			membershipCheckpoint                                      = new(model.MembershipCheckpoint)
			membershipRewards                                         []*model.MembershipReward
			membershipRewardID                                        int64
			membershipRewardAmount, totalFreshAmount, rewardMinAmount float64
			rewardMinLevel, membershipRewardLevel                     int8
		)

		attributes := reflect.ValueOf(csr.CustomerProfile.Attributes)
		for _, v := range attributes.MapKeys() {
			if v.String() == "membership_level" {
				membershipName = attributes.MapIndex(v).Interface().(string)
				continue
			}

			if v.String() == "membership_checkpoint" {
				membershipCheckpointStr = attributes.MapIndex(v).Interface().(string)
				continue
			}

			if v.String() == "fresh_product_revenue" {
				totalFreshAmount = attributes.MapIndex(v).Interface().(float64)
				continue
			}
		}

		// get membership level data
		if membershipName != "" {
			q = "select * from membership_level ml where ml.status = 1 and ml.name = ?"
			if e = o.Raw(q, membershipName).QueryRow(&membershipLevel); e != nil {
				continue
			}
		}

		// get membership checkpoint data
		if membershipCheckpointStr != "" {
			q = "select * from membership_checkpoint mc where mc.status = 1 and mc.checkpoint = ?"
			if e = o.Raw(q, membershipCheckpointStr).QueryRow(&membershipCheckpoint); e != nil {
				continue
			}
		}

		// check whether merchant deserve to get reward
		if membershipLevel.ID != 0 {
			q = "select ca.value from config_app ca where ca.attribute = 'minimum_level_membership_reward'"
			e = o.Raw(q).QueryRow(&rewardMinLevel)
		}

		if membershipLevel.Level >= rewardMinLevel {
			membershipRewardID = merchant.MembershipRewardID
			merchant.MembershipReward = &model.MembershipReward{ID: membershipRewardID}
			e = merchant.MembershipReward.Read("ID")

			// get membership reward data
			q = "select mr.* from membership_reward mr where mr.status = 1 and mr.reward_level >= ? order by mr.reward_level asc"
			_, e = o.Raw(q, merchant.MembershipReward.RewardLevel).QueryRows(&membershipRewards)

			// get minimum amount to get reward
			q = "select max(target_amount) " +
				"from membership_level ml " +
				"join membership_checkpoint mc on ml.id = mc.membership_level_id " +
				"join config_app ca on ca.`attribute` = 'minimum_level_membership_reward' and ml.`level` < ca.value"
			if e = o.Raw(q).QueryRow(&rewardMinAmount); e != nil {
				continue
			}

			membershipRewardAmount = totalFreshAmount - rewardMinAmount

			var voucherSnk string
			q = "select value from config_app where attribute = ?"
			o.Raw(q, "membership_reward_voucher_snk").QueryRow(&voucherSnk)

			voucherMap := make(map[string]*model.Voucher)
			for _, v := range membershipRewards {
				membershipRewardID = v.ID
				if membershipRewardAmount < v.MaxAmount {
					break
				}
				membershipRewardLevel = v.RewardLevel

				var (
					voucherImage  string
					voucherQuota  int64
					voucherAmount float64
				)

				// create voucher object for each reward level
				if v.RewardLevel <= 3 {
					loc, _ := time.LoadLocation("Asia/Jakarta")
					currTime := time.Now().In(loc)
					currYear := currTime.Year()
					currDateStr := currTime.Format("2006-01-02 00:00:00")
					currTime, _ = time.Parse("2006-01-02 15:04:05", currDateStr)
					currTime = currTime.In(loc)

					o.Raw(q, "membership_reward_voucher_quota_"+strconv.Itoa(int(v.RewardLevel))).QueryRow(&voucherQuota)
					o.Raw(q, "membership_reward_voucher_amount_"+strconv.Itoa(int(v.RewardLevel))).QueryRow(&voucherAmount)
					o.Raw(q, "membership_reward_voucher_image_"+strconv.Itoa(int(v.RewardLevel))).QueryRow(&voucherImage)

					voucher := &model.Voucher{
						Area:            &model.Area{ID: 1},
						Archetype:       &model.Archetype{ID: 22},
						RedeemCode:      "KJTNDSKN" + currTime.Format("2006"),
						Type:            1,
						Name:            "Voucher Belanja " + strconv.Itoa(int(voucherAmount)),
						StartTimestamp:  currTime,
						EndTimestamp:    time.Date(currYear, 12, 31, 23, 59, 59, 0, loc),
						OverallQuota:    voucherQuota,
						RemOverallQuota: voucherQuota,
						UserQuota:       voucherQuota,
						DiscAmount:      voucherAmount,
						Status:          1,
						ChannelVoucher:  "2,3",
						MerchantID:      merchant.ID,
					}

					voucher.VoucherContent = &model.VoucherContent{
						Voucher:        voucher,
						ImageUrl:       voucherImage,
						TermConditions: voucherSnk,
					}

					voucherMap[voucher.Name] = voucher
				}
			}
			vouchers[merchant.ID] = voucherMap

			if membershipRewardLevel > 0 {
				tempSOQuery := soQuery{
					SalesOrderId:          item.SalesOrderId,
					MerchantId:            item.MerchantId,
					FirebaseToken:         item.FirebaseToken,
					Code:                  "NOT0031",
					MembershipRewardLevel: strconv.Itoa(int(membershipRewardLevel)),
				}
				tempSalesOrder = append(tempSalesOrder, tempSOQuery)
			}
		}

		// set merchant membership related data
		totalPoint := float64(0)
		if _, isExist := merchantDatas[merchant.ID]; isExist {
			totalPoint += merchantDatas[merchant.ID].TotalPoint
		}
		merchantDatas[merchant.ID] = &merchantData{
			LevelID:      membershipLevel.ID,
			CheckpointID: membershipCheckpoint.ID,
			RewardID:     membershipRewardID,
			RewardAmount: membershipRewardAmount,
			TotalPoint:   totalPoint,
		}

		merchant.MembershipLevel = &model.MembershipLevel{ID: merchant.MembershipLevelID}
		merchant.MembershipLevel.Read("ID")
		if merchant.MembershipLevel.Level < membershipLevel.Level && membershipLevel.Level > 0 {
			code := "NOT0032"
			if membershipLevel.Level == 3 {
				code = "NOT0033"
			}

			tempSOQuery := soQuery{
				SalesOrderId:         item.SalesOrderId,
				MerchantId:           item.MerchantId,
				FirebaseToken:        item.FirebaseToken,
				Code:                 code,
				MembershipLevel:      membershipLevel.Name,
				MembershipCheckpoint: strconv.Itoa(int(membershipCheckpoint.Checkpoint)),
			}
			tempSalesOrder = append(tempSalesOrder, tempSOQuery)
		} else {
			merchant.MembershipCheckpoint = &model.MembershipCheckpoint{ID: merchant.MembershipCheckpointID}
			merchant.MembershipCheckpoint.Read("ID")
			if merchant.MembershipCheckpoint.Checkpoint < membershipCheckpoint.Checkpoint && membershipCheckpoint.Checkpoint > 0 {
				tempSOQuery := soQuery{
					SalesOrderId:         item.SalesOrderId,
					MerchantId:           item.MerchantId,
					FirebaseToken:        item.FirebaseToken,
					Code:                 "NOT0030",
					MembershipLevel:      membershipLevel.Name,
					MembershipCheckpoint: strconv.Itoa(int(membershipCheckpoint.Checkpoint)),
				}
				tempSalesOrder = append(tempSalesOrder, tempSOQuery)
			}
		}
		// end update membership level & checkpoint of merchant

		// add extra edenpoint from voucher
		if item.VoucherType == 4 {
			voucher := &model.Voucher{ID: item.VoucherID}
			if e = voucher.Read("ID"); e != nil {
				continue
			}
			lrecentPoint[item.MerchantId] += voucher.DiscAmount
			tempMpl := mplQuery{
				SalesOrderId:     item.SalesOrderId,
				DeliveryReturnId: item.DeliveryReturnId,
				MerchantId:       item.MerchantId,
				TotalPoint:       voucher.DiscAmount,
				RecentPoint:      lrecentPoint[item.MerchantId],
				TransactionType:  8,
			}
			tempMplQuery = append(tempMplQuery, tempMpl)
		}
	}

	o.Using("default")
	if e = o.Begin(); e != nil {
		return e
	}

	// Get term and tolerance expiration eden point
	termExpiredConfig, _ := repository.GetConfigApp("attribute", "eden_point_expiration_term")
	toleranceExpiredConfig, _ := repository.GetConfigApp("attribute", "eden_point_expiration_tolerance")

	termExpiredEdenpoint, _ := strconv.Atoi(termExpiredConfig.Value)
	toleranceExpiredEdenpoint, _ := strconv.Atoi(toleranceExpiredConfig.Value)

	currentTime := time.Now()
	// Set current time based on the custom date
	if customDate != "" {
		currentTime, _ = time.Parse("2006-01-02", customDate)
		currentDate = customDate
	}

	// Get Current day, month and year
	currentDay = currentTime.Day()
	currentMonth = currentTime.Month()
	currentYear = currentTime.Year()

	// Check current quartal phase
	quartalPhase = (int(currentMonth) / termExpiredEdenpoint)

	// Check if the current date is the same with expired month per phase
	if int(currentMonth)%termExpiredEdenpoint == 0 {
		// Get last day in current month
		lastDateCurrentMonth := time.Date(currentYear, currentMonth+1, 0, 0, 0, 0, 0, currentTime.Location())
		lastDay := lastDateCurrentMonth.Day()
		// Set up expiration month
		currentMonthExpiration = (quartalPhase * termExpiredEdenpoint)
		nextMonthExpiration = ((quartalPhase + 1) * termExpiredEdenpoint)

		// If the current days more than equal tolerance days on current month, expired month set to next quartal
		if currentDay >= (lastDay - toleranceExpiredEdenpoint) {
			// Expired in the next phase
			quartalPhase += 1
			isExpiredNextPeriod = true
			nextMonthExpiration = (quartalPhase * termExpiredEdenpoint)
		}
	} else {
		// Expired in the next phase
		quartalPhase += 1
		currentMonthExpiration = (quartalPhase * termExpiredEdenpoint)
		nextMonthExpiration = ((quartalPhase + 1) * termExpiredEdenpoint)
	}

	// Get current expiration date and next expiration date
	currentPeriodDate := time.Date(currentYear, time.Month(currentMonthExpiration)+1, 0, 0, 0, 0, 0, currentTime.Location())
	nextPeriodDate := time.Date(currentYear, time.Month(nextMonthExpiration)+1, 0, 0, 0, 0, 0, currentTime.Location())

	// Get Actual Expired month
	expiredMonth = (quartalPhase * termExpiredEdenpoint) + 1

	// Get actual expired date eden point
	expiredDate := time.Date(currentYear, time.Month(expiredMonth), 0, 0, 0, 0, 0, currentTime.Location()).Format("2006-01-02")

	pointSummary := make(map[int64]float64)
	//loop temporary merchant point log local to insert to database
	for _, item := range tempMplQuery {
		var (
			columnStr, qMarkStr string
			colValue            []interface{}
		)

		columnStr = "merchant_id, sales_order_id, point_value, recent_point, status, created_date, expired_date, campaign_id, campaign_name, campaign_multiplier, transaction_type, referrer_id, referee_id"
		qMarkStr = "?, ?, ?, ?, 1, ?, ?, ?, ?, ?, ?, ?, ?"
		colValue = []interface{}{item.MerchantId, item.SalesOrderId, item.TotalPoint, item.RecentPoint, currentDate, expiredDate, item.TalonCampaignID, item.TalonCampaignName, item.TalonMultiplier, item.TransactionType, item.ReferrerID, item.RefereeID}

		_, e = o.Raw("INSERT INTO merchant_point_log "+
			"("+columnStr+") "+
			" VALUES("+qMarkStr+"); ", colValue).Exec()

		if e != nil {
			o.Rollback()
			return e
		}

		if _, isExist := merchantDatas[item.MerchantId]; isExist {
			merchantDatas[item.MerchantId].TotalPoint = item.RecentPoint
		} else {
			merchantDatas[item.MerchantId] = &merchantData{TotalPoint: item.RecentPoint}
		}

		pointSummary[item.MerchantId] += item.TotalPoint
	}

	//looping merchant to update merchant total point
	for merchantID, v := range merchantDatas {
		var (
			setString string
			setValue  []interface{}
		)
		if v.CheckpointID != 0 {
			setString += "membership_checkpoint_id = ?, "
			setValue = append(setValue, v.CheckpointID)
		}
		if v.LevelID != 0 {
			setString += "membership_level_id = ?, "
			setValue = append(setValue, v.LevelID)
		}
		if v.RewardAmount != 0 {
			setString += "membership_reward_amount = ?, "
			setValue = append(setValue, v.RewardAmount)
		}
		if v.RewardID != 0 {
			setString += "membership_reward_id = ?, "
			setValue = append(setValue, v.RewardID)
		}
		if v.TotalPoint != 0 {
			setString += "total_point = ?, "
			setValue = append(setValue, v.TotalPoint)
		}
		setString = strings.TrimSuffix(setString, ", ")

		_, e = o.Raw("UPDATE merchant SET "+setString+" WHERE id = ?", setValue, merchantID).Exec()

		key := "total_current_eden_point"
		wib, _ := time.LoadLocation("Asia/Jakarta")
		currentTime, _ := time.ParseInLocation("2006-01-02 15:04:05", time.Now().Format("2006-01-02 15:04:05"), wib)
		totalCurrentPoint := 0.00
		if dbredis.Redis.CheckExistKey(key) {
			var merchantPoint float64
			e = o.Raw("select total_point from merchant where id = ?", merchantID).QueryRow(&merchantPoint)
			dbredis.Redis.GetCache(key, &totalCurrentPoint)
			totalCurrentPoint = totalCurrentPoint - merchantPoint + v.TotalPoint
		} else {
			totalCurrentPoint = v.TotalPoint
		}
		dbredis.Redis.SetCache(key, totalCurrentPoint, 0)
		dbredis.Redis.SetCache(key+"_updated_date", currentTime, 0)
	}

	// insert or update into merchant point summary
	for i, v := range pointSummary {
		var isExist, isExistMerchantPointExpiration bool
		if e = o.Raw("select exists(select id from merchant_point_summary where summary_date = ? and merchant_id = ?)", currentDate, i).QueryRow(&isExist); e != nil || (e == nil && !isExist) {
			_, e = o.Raw("insert into merchant_point_summary (merchant_id, summary_date, earned_point, redeemed_point) values (?, ?, ?, 0)", i, currentDate, v).Exec()
		} else {
			_, e = o.Raw("update merchant_point_summary set earned_point = earned_point + ? where summary_date = ? and merchant_id = ?", v, currentDate, i).Exec()
		}

		// Check if merchant exists in merchant_point_expiration
		if e = o.Raw("select exists(select merchant_id from merchant_point_expiration where merchant_id = ?)", i).QueryRow(&isExistMerchantPointExpiration); e != nil {
			continue
		}
		merchantPointExpiration := &model.MerchantPointExpiration{ID: i}
		// if merchant doesn't exist in merchant_point_expiration, insert new data
		if !isExistMerchantPointExpiration {
			// Set point based on current date, point goes to current or next period
			if isExpiredNextPeriod {
				merchantPointExpiration.NextPeriodPoint = v
			} else {
				merchantPointExpiration.CurrentPeriodPoint = v
			}
			merchantPointExpiration.CurrentPeriodDate = currentPeriodDate
			merchantPointExpiration.NextPeriodDate = nextPeriodDate
			merchantPointExpiration.LastUpdatedAt = time.Now()
			if _, e = o.Insert(merchantPointExpiration); e != nil {
				continue
			}
			// if merchant exist in merchant_point_expiration, update data
		} else {
			if e = merchantPointExpiration.Read("ID"); e != nil {
				continue
			}
			// Set point based on current date, point goes to current or next period
			if isExpiredNextPeriod {
				merchantPointExpiration.NextPeriodPoint += v
			} else {
				merchantPointExpiration.CurrentPeriodPoint += v
			}

			merchantPointExpiration.LastUpdatedAt = time.Now()
			if e = merchantPointExpiration.Save(); e != nil {
				continue
			}
		}
	}
	// insert new voucher if merchant is in membership reward level
	for _, v := range vouchers {
		for _, val := range v {
			var code string
			if code, e = util.CheckTable("voucher"); e == nil {
				code, e = util.GenerateCode(code, "voucher")
			}

			val.Code = code
			val.RedeemCode += code[3:]
			if _, e = o.Insert(val); e != nil {
				continue
			}

			if _, e = o.Insert(val.VoucherContent); e != nil {
				continue
			}

			if e = log.AuditLogByUser(&model.Staff{ID: 222}, val.ID, "voucher", "auto_create", "membership voucher"); e != nil {
				continue
			}
		}
	}

	for _, v := range redeemedData {
		e = talon.ChangeTalonPoints("deduct_points", "redeem points", v["merchant_profile"].(string), v["redeemed_amount"].(float64))
	}

	if e = o.Commit(); e != nil {
		return e
	}
	if e != nil {
		return e
	}

	o.Using("read_only")
	//send notif
	for _, item := range tempSalesOrder {
		mn := &util.MessageNotification{}
		if e = o.Raw("SELECT message, title, type FROM notification WHERE code= ?", item.Code).QueryRow(&mn); e != nil {
			return e
		}

		modelNotif := &util.ModelNotification{}

		modelNotif.SendTo = item.FirebaseToken
		modelNotif.Title = mn.Title
		modelNotif.Message = mn.Message
		modelNotif.Type = mn.Type
		modelNotif.RefID = item.SalesOrderId
		modelNotif.MerchantID = item.MerchantId
		modelNotif.ServerKey = util.ServerKeyFireBase

		if e = util.PostModelNotification(modelNotif); e != nil {
			return e
		}

	}
	return
}

func (h *Handler) priceScheduler(c echo.Context) (e error) {
	var err error

	o := orm.NewOrm()
	orSelect := orm.NewOrm()

	orSelect.Using("read_only")
	o.Begin()

	var priceSetSch []model.PriceSchedule

	orSelect.Raw("SELECT * FROM price_schedule where status = 1").QueryRows(&priceSetSch)

	layoutDate := "2006-01-02"
	layoutTime := "15:04"
	today := time.Now()
	currentDate, _ := time.Parse(layoutDate, fmt.Sprintf("%d-%02d-%02d", today.Year(), today.Month(), today.Day()))
	currentTime, _ := time.Parse(layoutTime, fmt.Sprintf("%02d:%02d", today.Hour(), today.Minute()))

	for _, ap := range priceSetSch {

		formattedDate, _ := time.Parse(layoutDate, ap.ScheduleDate)
		formattedTime, _ := time.Parse(layoutTime, ap.ScheduleTime)

		if currentDate == formattedDate {
			if currentTime.Equal(formattedTime) || currentTime.After(formattedTime) {

				pss := &model.PriceSchedule{
					ID:     ap.ID,
					Status: 2,
				}
				if _, e = o.Update(pss, "Status"); e == nil {
					orSelect.LoadRelated(pss, "PriceScheduleDumps", 1, 10000)

					for _, row := range pss.PriceScheduleDumps {

						p := &model.Price{Product: row.Product, PriceSet: row.PriceSet}
						err = p.Read("Product", "PriceSet")

						p = &model.Price{
							ID:        p.ID,
							UnitPrice: math.Round(row.UnitPrice),
						}

						if _, err := o.Update(p, "unit_price"); err != nil {
							o.Rollback()
						}

						pl := &model.PriceLog{
							PriceID:   p.ID,
							UnitPrice: p.UnitPrice,
							CreatedAt: time.Now(),
							CreatedBy: &model.Staff{ID: 222},
						}
						if _, err = o.Insert(pl); err != nil {
							o.Rollback()
						}
					}

					if err == nil {
						o.Raw("DELETE FROM price_schedule_dump where price_schedule_id = ?", ap.ID).Exec()
					}

					err := log.AuditLogByUser(&model.Staff{ID: 222}, pss.ID, "product price schedule", "update", "Update Product Price Using Scheduler")
					if err != nil {
						o.Rollback()
						return
					}
				}
			}
		}
	}

	if err == nil {
		o.Commit()
	}
	return
}

// read : function to get requested data based on parameters
func (h *Handler) remindPayment(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)

	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	orSelect.Raw("SELECT order_time_limit,area_id  FROM area_policy").QueryRows(&GetAreaPolicy)

	for _, ap := range GetAreaPolicy {
		areaPolicyObject := ap // create a new "areaPolicyObject" variable on each iteration (goroutines on loop iterator variables)
		// this function for push notif before
		layout := "15:04:05"
		t, _ := time.Parse(layout, ap.OrderTimeLimit+":00")
		t2, _ := time.Parse("15:04:05", t.Add(time.Hour*-1).Format("15:04:05"))
		t3, _ := time.Parse("15:04:05", time.Now().Format("15:04:05"))
		t4, _ := time.Parse("15:04:05", t.Add(time.Minute*10).Format("15:04:05"))

		if t3.After(t2) && t3.Before(t) {
			pushNotification(areaPolicyObject, "reminderPayment")
		}
		if t3.After(t) && t3.Before(t4) {
			pushNotification(areaPolicyObject, "cancelled")
		}

	}

	return ctx.Serve(e)
}

type areaPolicy struct {
	OrderTimeLimit string `orm:"column(order_time_limit);null" json:"order_time_limit"`
	Area           string `orm:"column(area_id)" json:"area_id"`
}

type so struct {
	ID            int64  `orm:"column(id);" json:"id"`
	Code          string `orm:"column(code);size(45);null" json:"code"`
	FirebaseToken string `orm:"column(firebase_token);null" json:"firebase_token"`
	MerchantID    int64  `orm:"column(merchant_id);null" json:"merchant_id"`
}

func pushNotification(areaPolicy areaPolicy, status string) {
	var bulan string
	var salesOrder []so
	var voucID string
	var pointID string
	var recentPoint float64
	var pointValue float64

	orSelect := orm.NewOrm()
	orSelect.Using("read_only")
	mn := &util.MessageNotification{}
	modelNotif := &util.ModelNotification{}
	currentTime := time.Now()
	if status == "cancelled" {
		orSelect.Raw("SELECT so.id id, so.code code, um.firebase_token firebase_token, m.id merchant_id "+
			"FROM sales_order so "+
			"INNER JOIN branch b ON b.id = so.branch_id "+
			"INNER JOIN merchant m ON m.id = b.merchant_id "+
			"INNER JOIN user_merchant um ON um.id = m.user_merchant_id "+
			"WHERE so.area_id = ? AND so.status= 1 AND DATE(so.delivery_date) = ? AND so.payment_group_sls_id = 1 AND so.has_ext_invoice = 2 AND so.payment_reminder= 1;", areaPolicy.Area, currentTime.Add(time.Hour*24).Format("2006-01-02")).QueryRows(&salesOrder)
	} else if status == "reminderPayment" {
		orSelect.Raw("SELECT so.id id, so.code code, um.firebase_token firebase_token, m.id merchant_id "+
			"FROM sales_order so "+
			"INNER JOIN branch b ON b.id = so.branch_id "+
			"INNER JOIN merchant m ON m.id = b.merchant_id "+
			"INNER JOIN user_merchant um ON um.id = m.user_merchant_id "+
			"WHERE so.area_id = ? AND so.status= 1 AND DATE(so.delivery_date) = ? AND so.payment_group_sls_id = 1  AND so.payment_reminder= 2;", areaPolicy.Area, currentTime.Add(time.Hour*24).Format("2006-01-02")).QueryRows(&salesOrder)
	}

	if len(salesOrder) > 0 {
		for _, s := range salesOrder {
			if status == "cancelled" {
				orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0005'").QueryRow(&mn)
				mn.Message = util.ReplaceNotificationSalesOrder(mn.Message, "#sales_order_code#", s.Code)

				// returns voucher remaining if sales order have voucher
				orSelect.Raw("SELECT voucher_id from sales_order WHERE id = ?;", s.ID).QueryRow(&voucID)
				if voucID != "" {
					orm.NewOrm().Raw("UPDATE voucher_log SET status = 3 WHERE sales_order_id = ? AND voucher_id = ?;", s.ID, voucID).Exec()
					orm.NewOrm().Raw("UPDATE voucher SET rem_overall_quota = rem_overall_quota + 1 WHERE id = ?;", voucID).Exec()
				}
				//returns point remaining if sales order using point
				orSelect.Raw("SELECT point_redeem_id FROM sales_order WHERE id = ?;", s.ID).QueryRow(&pointID)
				if pointID != "" {
					orSelect.Raw("SELECT recent_point FROM merchant_point_log where merchant_id = ? order by id desc limit 1 ", s.MerchantID).QueryRow(&recentPoint)
					orSelect.Raw("SELECT point_value FROM merchant_point_log where merchant_id = ? and id = ? order by id desc limit 1 ", s.MerchantID, pointID).QueryRow(&pointValue)
					totalPoint := pointValue + recentPoint
					orm.NewOrm().Raw(
						"UPDATE merchant_point_log SET status = 4, note = 'Cancellation due to cancel sales order'"+
							"WHERE id = ? and sales_order_id = ? and status = 2 ", pointID, s.ID).Exec()
					orm.NewOrm().Raw("INSERT INTO merchant_point_log"+
						"(merchant_id,sales_order_id,point_value,recent_point,status,created_date,note)"+
						"SELECT merchant_id,sales_order_id,point_value,? ,1, ? ,"+
						"'Point Issued From Cancellation Redeem' FROM merchant_point_log mpl where mpl.id = ? ", time.Now(), totalPoint, pointID).Exec()
					orm.NewOrm().Raw("UPDATE merchant SET total_point = ? WHERE id= ?", totalPoint, s.MerchantID).Exec()
				}

				orm.NewOrm().Raw("UPDATE sales_order SET status = 3, cancel_type = 2 WHERE id = ?;", s.ID).Exec()

				modelNotif.SendTo = s.FirebaseToken
				modelNotif.Title = mn.Title
				modelNotif.Message = mn.Message
				modelNotif.Type = "1"
				modelNotif.RefID = s.ID
				modelNotif.MerchantID = s.MerchantID
				modelNotif.ServerKey = util.ServerKeyFireBase
				util.PostModelNotification(modelNotif)

			} else if status == "reminderPayment" {
				orm.NewOrm().Raw("UPDATE sales_order SET  payment_reminder = 1  WHERE id = ?;", s.ID).Exec()
				orSelect.Raw("SELECT message, title FROM notification WHERE code= 'NOT0007'").QueryRow(&mn)
				year, month, day := time.Now().Date()
				if month.String() == "January" {
					bulan = "Jan"
				} else if month.String() == "February" {
					bulan = "Feb"
				} else if month.String() == "March" {
					bulan = "Mar"
				} else if month.String() == "April" {
					bulan = "Apr"
				} else if month.String() == "May" {
					bulan = "Mei"
				} else if month.String() == "June" {
					bulan = "Jun"
				} else if month.String() == "July" {
					bulan = "Jul"
				} else if month.String() == "August" {
					bulan = "Ags"
				} else if month.String() == "September" {
					bulan = "Sep"
				} else if month.String() == "October" {
					bulan = "Okt"
				} else if month.String() == "November" {
					bulan = "Nov"
				} else if month.String() == "December" {
					bulan = "Des"
				}
				mn.Message = util.ReplaceNotificationSalesOrder(mn.Message, "#current_date#", strconv.Itoa(day)+" "+bulan+" "+strconv.Itoa(year))
				mn.Message = util.ReplaceNotificationSalesOrder(mn.Message, "#time_limit#", areaPolicy.OrderTimeLimit)
				mn.Message = util.ReplaceNotificationSalesOrder(mn.Message, "#sales_order_code#", s.Code)

				modelNotif.SendTo = s.FirebaseToken
				modelNotif.Title = mn.Title
				modelNotif.Message = mn.Message
				modelNotif.Type = "1"
				modelNotif.RefID = s.ID
				modelNotif.MerchantID = s.MerchantID
				modelNotif.ServerKey = util.ServerKeyFireBase
				util.PostModelNotification(modelNotif)
			}
		}
	}
	return
}

func (h *Handler) suspendMerchantScheduler(c echo.Context) (e error) {
	var err error

	o := orm.NewOrm()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	o.Begin()

	var merchants []model.Merchant

	orSelect.Raw(`SELECT m.id id FROM merchant m
					LEFT JOIN branch b ON b.merchant_id = m.id
					LEFT JOIN sales_order so ON so.branch_id = b.id
					LEFT JOIN sales_invoice si ON si.sales_order_id = so.id
					WHERE (si.status = 1 OR si.status = 6)
					AND DATE_ADD(si.due_date, INTERVAL 1 DAY) < NOW()
					AND m.suspended != 1 AND so.term_payment_sls_id NOT IN (7,10,11,14) AND so.order_type_sls_id !=13 
				`).QueryRows(&merchants)

	for _, m := range merchants {
		merchantToUpdate := &model.Merchant{
			ID:        m.ID,
			Suspended: 1,
		}

		if _, err = o.Update(merchantToUpdate, "Suspended"); err != nil {
			o.Rollback()
		}
	}

	o.Commit()

	return
}

func (h *Handler) notificationCampaign(c echo.Context) (e error) {
	o := orm.NewOrm()
	// set charset to utf8mb4 for supporting emoji connection
	o.Raw("SET NAMES 'utf8mb4'").Exec()

	var notificationCampaigns []model.NotificationCampaign
	o.Raw("SELECT `id`, `code`, `campaign_name`, `area`, `archetype`, `redirect_to`, `redirect_value`, `title`, `message`, `push_now`, `scheduled_at`, `status` FROM `notification_campaign` WHERE `scheduled_at` <= CURRENT_TIMESTAMP and `status` = '1' and `push_now` = '2'").QueryRows(&notificationCampaigns)
	for _, m := range notificationCampaigns {
		idAreas := strings.Split(m.Area, ",")
		idArchetype := strings.Split(m.Archetype, ",")

		var areasID []int64
		var archetypesID []int64

		// get areas array of id to models
		for _, a := range idAreas {
			id, _ := strconv.Atoi(a)
			areasID = append(areasID, int64(id))
		}

		// get archetypes array of id to models
		for _, a := range idArchetype {
			id, _ := strconv.Atoi(a)
			archetypesID = append(archetypesID, int64(id))
		}

		// checking redirect_to to glossary
		var redirectToGlossary *model.Glossary
		redirectToGlossary, e = repository.GetGlossaryMultipleValue("table", "notification_campaign", "attribute", "redirect_to", "value_int", m.RedirectTo)
		if e != nil {
			return
		}

		// set redirect_to_name
		redirectToName := redirectToGlossary.ValueName

		// set value name by redirect_to
		switch redirectToGlossary.ValueName {
		case "Product":
			id, _ := strconv.Atoi(m.RedirectValue)
			m.RedirectValue = common.Encrypt(id)
		case "Product Tag":
			id, _ := strconv.Atoi(m.RedirectValue)
			m.RedirectValue = common.Encrypt(id)
		case "URL":
		default:
		}

		messageNotif := &util.MessageNotificationCampaign{
			ID:             common.Encrypt(m.ID),
			Code:           m.Code,
			CampaignName:   m.CampaignName,
			Area:           areasID,
			Archetype:      archetypesID,
			RedirectTo:     m.RedirectTo,
			RedirectToName: redirectToName,
			RedirectValue:  m.RedirectValue,
			Title:          m.Title,
			Message:        m.Message,
			ServerKey:      util.CampaignServerKeyFireBase,
		}

		// push notification
		if e = util.PostModelNotificationCampaign(messageNotif); e != nil {
			return
		}
	}

	return
}

func (h *Handler) autoFailedSalesAssignmentScheduler(c echo.Context) (e error) {
	var err error

	o := orm.NewOrm()
	orSelect := orm.NewOrm()
	orSelect.Using("read_only")

	o.Begin()

	loc, _ := time.LoadLocation("Asia/Jakarta")
	currentTime := time.Now().In(loc)
	var task []model.SalesAssignmentItem

	orSelect.Raw(`SELECT * FROM sales_assignment_item
					WHERE end_date = ?
					AND status = 1`,
		currentTime.Format("2006-01-02")).QueryRows(&task)

	if task != nil {
		reason := model.Glossary{}
		orSelect.Raw("SELECT * FROM glossary WHERE `table` = 'sales_failed_visit' AND `attribute` = 'failed_status' AND `value_int` = 6").QueryRow(&reason)
		for _, m := range task {
			saiToUpdate := &model.SalesAssignmentItem{
				ID:         m.ID,
				SubmitDate: currentTime,
				Status:     14,
			}

			if _, err = o.Update(saiToUpdate, "SubmitDate", "Status"); err != nil {
				o.Rollback()
			}

			failed := model.SalesFailedVisit{
				SalesAssignmentItem: &m,
				FailedStatus:        reason.ValueInt,
				DescriptionFailed:   reason.ValueName,
			}
			if _, err = o.Insert(&failed); err != nil {
				o.Rollback()
			}
			o.Commit()

			o.Begin()
			var sai []model.SalesAssignmentItem
			orSelect.Raw("SELECT * FROM sales_assignment_item WHERE sales_assignment_id = ? AND id != ? AND status = 1", m.SalesAssignment.ID, m.ID).QueryRows(&sai)
			if sai == nil {
				saToUpdate := &model.SalesAssignment{
					ID:     m.SalesAssignment.ID,
					Status: 2,
				}
				if _, err = o.Update(saToUpdate, "Status"); err != nil {
					o.Rollback()
				}
			}
			o.Commit()
		}
	}
	return
}

func (h *Handler) summaryNotificationCampaign(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	o := orm.NewOrm()
	o.Using("read_only")
	var qryStr, period string

	if ctx.QueryParam("period") != "" {
		period = ctx.QueryParam("period")
	}

	var notificationCampaigns []*model.NotificationCampaign
	switch period {
	case "hourly":
		qryStr = "SELECT `id`, `code`, `campaign_name`, `status` FROM `notification_campaign` WHERE `scheduled_at` >= date_sub(NOW(), interval 1 week)"
	case "daily":
		qryStr = "SELECT `id`, `code`, `campaign_name`, `status` FROM `notification_campaign` WHERE `scheduled_at` >= date_sub(NOW(), interval 3 month)"
	default:
		return ctx.Serve(e)
	}

	o.Raw(qryStr).QueryRows(&notificationCampaigns)

	md := mongodb.NewMongo()

	for _, m := range notificationCampaigns {
		idStr := common.Encrypt(m.ID)
		// get total success sent for item from mongoDB
		totalSuccessSent, e := md.GetCountDataWithFilter("Notification_Campaign_Item", &model.NotificationCampaignItem{
			NotificationCampaignID: idStr,
			Sent:                   1,
		})
		if e != nil {
			md.DisconnectMongoClient()
			return e
		}

		// get total failed sent for item from mongoDB
		totalFailedSent, e := md.GetCountDataWithFilter("Notification_Campaign_Item", &model.NotificationCampaignItem{
			NotificationCampaignID: idStr,
			Sent:                   2,
		})
		if e != nil {
			md.DisconnectMongoClient()
			return e
		}

		// get total failed sent for item from mongoDB
		totalOpened, e := md.GetCountDataWithFilter("Notification_Campaign_Item", &model.NotificationCampaignItem{
			NotificationCampaignID: idStr,
			Opened:                 1,
		})
		if e != nil {
			md.DisconnectMongoClient()
			return e
		}

		m.SuccessSent = totalSuccessSent
		m.FailedSent = totalFailedSent
		m.Open = totalOpened

		e = m.Save("SuccessSent", "FailedSent", "Open")
		if e != nil {
			return e
		}
	}

	md.DisconnectMongoClient()
	return
}

// setUpLoyaltyPoints : function to set up eden points for merchant from talon.one's session return
func setUpLoyaltyPoints(profileCode string, v model.Effects, item soQuery, lrecentPoint map[int64]float64) (varMplQuery mplQuery, varSoQuery soQuery, lrecentPointReturn map[int64]float64, err error) {
	var campaignMap = make(map[int]*model.CampaignDetail)

	o := orm.NewOrm()
	o.Using("read_only")

	campaign, isExist := campaignMap[v.CampaignID]
	if !isExist {
		campaign, _ = talon.GetCampaignDetail(strconv.Itoa(v.CampaignID))
		campaignMap[v.CampaignID] = campaign
	}

	getPoint := v.Props.Value.(float64)
	campaignTags := "," + strings.Join(campaign.Tags, ",") + ","
	varSoQuery = soQuery{
		SalesOrderId:     item.SalesOrderId,
		DeliveryReturnId: item.DeliveryReturnId,
		AmountPayment:    item.AmountPayment,
		MerchantId:       item.MerchantId,
		FirebaseToken:    item.FirebaseToken,
		Code:             "NOT0012",
	}
	varMplQuery = mplQuery{
		SalesOrderId:      item.SalesOrderId,
		DeliveryReturnId:  item.DeliveryReturnId,
		MerchantId:        item.MerchantId,
		TotalPoint:        getPoint,
		TalonCampaignID:   int64(v.CampaignID),
		TalonCampaignName: campaign.Name,
		TalonMultiplier:   int8(campaign.Attributes.Multiplier),
		TransactionType:   2,
	}

	// start set transaction type and merchant id that get point based on tags
	if strings.Contains(campaignTags, ",default,") {
		varMplQuery.TransactionType = 1
	} else if strings.Contains(campaignTags, ",referral,") {
		varMplQuery.TransactionType = 4
		merchant, _ := repository.GetMerchant("profile_code", v.Props.RecipientIntegrationID)
		varMplQuery.MerchantId = merchant.ID
		if merchant.Referrer != nil {
			varMplQuery.ReferrerID = merchant.Referrer.ID
		}
		if profileCode != v.Props.RecipientIntegrationID {
			varMplQuery.TransactionType = 3
			varMplQuery.RefereeID = item.MerchantId
			varMplQuery.ReferrerID = 0

			if val, isExist := lrecentPoint[varMplQuery.MerchantId]; !isExist || (isExist && val == 0) {
				currentPointReferrer := 0.00
				err = o.Raw("select total_point from merchant where id = ?", varMplQuery.MerchantId).QueryRow(&currentPointReferrer)
				lrecentPoint[varMplQuery.MerchantId] = currentPointReferrer
			}
		}
	}
	sumPoint := lrecentPoint[varMplQuery.MerchantId] + getPoint
	lrecentPoint[varMplQuery.MerchantId] = sumPoint
	varSoQuery.TotalPoint = getPoint
	varMplQuery.RecentPoint = sumPoint
	lrecentPointReturn = lrecentPoint
	// end set transaction type and merchant id that get point based on tags

	return
}

// createAnniversaryJoinVoucher : function to create voucher for customer, especially have anniversary join date
func (h *Handler) createAnniversaryJoinVoucher(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	o := orm.NewOrm()
	o.Using("read_only")

	ormWriter := orm.NewOrm()

	var (
		currentYear, currentMonth, currentDay                             int
		code, imageUrl, membershipLevel, termCondition, redeemCode, value string
		discAmount                                                        float64
		listMerchantAnniversary                                           []*model.Merchant
		configApp                                                         *model.ConfigApp
	)
	currentTime := time.Now()
	currentYear = currentTime.Year()
	currentMonth = int(currentTime.Month())
	currentDay = currentTime.Day()
	anniversaryDate := ctx.QueryParam("anniversary_date")

	if anniversaryDate != "" {
		anniversaryDateArr := strings.Split(anniversaryDate, "-")
		currentYear, _ = strconv.Atoi(anniversaryDateArr[0])
		currentMonth, _ = strconv.Atoi(anniversaryDateArr[1])
		currentDay, _ = strconv.Atoi(anniversaryDateArr[2])
	}

	lastTimeOfTheCurrentYear := time.Date(currentYear, 12, 31, 23, 59, 59, 0, time.Local)

	// Get eligible membership business type
	if configApp, e = repository.GetConfigApp("attribute", "eligible_membership_business_type"); e != nil {
		return e
	}

	eligibleMembershipArr := strings.Split(configApp.Value, ",")

	for counter := 0; counter < len(eligibleMembershipArr); counter++ {
		value = value + "?, "
	}

	value = strings.TrimSuffix(value, ", ")

	// Get list merchant that anniversary
	query := "SELECT * FROM merchant m WHERE YEAR(m.created_at) < ? AND MONTH(m.created_at)= ? AND DAY(m.created_at)=? AND m.status = 1 AND m.business_type_id IN (" + value + ")"
	if _, e = o.Raw(query, currentYear, currentMonth, currentDay, eligibleMembershipArr).QueryRows(&listMerchantAnniversary); e != nil {
		return e
	}

	// Get Term & Conditions for Anniversary Voucher
	if configApp, e = repository.GetConfigApp("attribute", "term_and_condition_anniversary_voucher"); e != nil {
		return e
	}
	termCondition = configApp.Value
	termCondition = strings.ReplaceAll(termCondition, "#tahun#", fmt.Sprintf("%d", currentYear))

	for _, v := range listMerchantAnniversary {
		ormWriter.Begin()

		if code, e = util.CheckTable("voucher"); e != nil {
			ormWriter.Rollback()
			continue
		}
		if code, e = util.GenerateCode(code, "voucher"); e != nil {
			ormWriter.Rollback()
			continue
		}

		redeemCode = fmt.Sprintf("ANV%d%s", currentYear, code[3:])

		switch v.MembershipLevelID {
		case 2:
			membershipLevel = "juragan"
		case 3:
			membershipLevel = "konglomerat"
		default:
			membershipLevel = "pemula"
		}

		// Get amount anniversary voucher per level
		if configApp, e = repository.GetConfigApp("attribute", "amount_anniversary_voucher_"+membershipLevel); e != nil {
			ormWriter.Rollback()
			continue
		}
		discAmount, _ = strconv.ParseFloat(configApp.Value, 64)

		// Get url image anniversary voucher per level
		if configApp, e = repository.GetConfigApp("attribute", "url_image_anniversary_voucher_"+membershipLevel); e != nil {
			ormWriter.Rollback()
			continue
		}
		imageUrl = configApp.Value

		// set default membership level id for merchant that having membership level id is 0
		if v.MembershipLevelID == 0 {
			q := "SELECT id FROM membership_level WHERE status=1 ORDER BY level LIMIT 1"
			if e = o.Raw(q).QueryRow(&v.MembershipLevelID); e != nil {
				continue
			}
			q = "SELECT id FROM membership_checkpoint WHERE status=1 ORDER BY checkpoint LIMIT 1"
			if e = o.Raw(q).QueryRow(&v.MembershipCheckpointID); e != nil {
				continue
			}
		}

		v := &model.Voucher{
			MerchantID:             v.ID,
			Area:                   &model.Area{ID: 1},
			Archetype:              &model.Archetype{ID: 22},
			Code:                   code,
			RedeemCode:             redeemCode,
			Type:                   4,
			Name:                   "Anniversary Voucher",
			StartTimestamp:         currentTime,
			EndTimestamp:           lastTimeOfTheCurrentYear,
			OverallQuota:           1,
			RemOverallQuota:        1,
			UserQuota:              1,
			MinOrder:               0,
			DiscAmount:             discAmount,
			Note:                   "Auto Create Anniversary Voucher",
			Status:                 int8(1),
			ChannelVoucher:         "2,3",
			MembershipLevelID:      v.MembershipLevelID,
			MembershipCheckpointID: v.MembershipCheckpointID,
		}

		if _, e = ormWriter.Insert(v); e != nil {
			ormWriter.Rollback()
			continue
		}

		vc := &model.VoucherContent{Voucher: v, ImageUrl: imageUrl, TermConditions: termCondition}
		if _, e = ormWriter.Insert(vc); e != nil {
			ormWriter.Rollback()
			continue
		}
		if e = log.AuditLogByUser(&model.Staff{ID: 222}, v.ID, "voucher", "auto_create", "anniversary voucher"); e != nil {
			ormWriter.Rollback()
			continue
		}
		ormWriter.Commit()
	}
	return e
}

func (h *Handler) ExpiredEdenPoint(c echo.Context) (e error) {
	ctx := c.(*cuxs.Context)
	o := orm.NewOrm()
	o.Using("read_only")

	ormWriter := orm.NewOrm()

	var (
		data                                                           []*model.MerchantPointExpiration
		isExistMpe, isExistMps                                         bool
		recentPoint                                                    float64
		customDate                                                     string
		currentMonth                                                   time.Month
		currentYear, monthCurrentPeriod, monthNextPeriod, quartalPhase int
		currentPeriodDate, nextPeriodDate                              time.Time
	)

	currentTime := time.Now()
	customDate = ""
	// customDate : to set custom current date, be used in testing
	if ctx.QueryParam("custom_date") != "" {
		customDate = ctx.QueryParam("custom_date")
		currentTime, _ = time.Parse("2006-01-02", customDate)
	}

	currentMonth = currentTime.Month()
	currentYear = currentTime.Year()

	termExpiredConfig, _ := repository.GetConfigApp("attribute", "eden_point_expiration_term")
	termExpiredEdenpoint, _ := strconv.Atoi(termExpiredConfig.Value)

	quartalPhase = (int(currentMonth) / termExpiredEdenpoint)

	monthCurrentPeriod = ((quartalPhase + 1) * termExpiredEdenpoint)
	monthNextPeriod = ((quartalPhase + 2) * termExpiredEdenpoint)

	// Get current expiration date and next expiration date
	currentPeriodDate = time.Date(currentYear, time.Month(monthCurrentPeriod)+1, 0, 0, 0, 0, 0, currentTime.Location())
	nextPeriodDate = time.Date(currentYear, time.Month(monthNextPeriod)+1, 0, 0, 0, 0, 0, currentTime.Location())

	// Check if there is data with expired date same as current date
	o.Raw("select exists(select merchant_id from merchant_point_expiration where current_period_date = ?)", currentTime.Format("2006-01-02")).QueryRow(&isExistMpe)

	if isExistMpe {
		// Get List Merchant Point Expiration
		o.Raw("Select * from merchant_point_expiration WHERE current_period_date = ?", currentTime.Format("2006-01-02")).QueryRows(&data)

		for _, v := range data {
			ormWriter.Begin()

			if v.CurrentPeriodPoint != 0 {
				o.Raw("SELECT recent_point from merchant_point_log where merchant_id = ? order by id desc limit 1 ", v.ID).QueryRow(&recentPoint)

				// Reduce recent point
				recentPoint -= v.CurrentPeriodPoint

				mpl := &model.MerchantPointLog{
					PointValue:      v.CurrentPeriodPoint,
					RecentPoint:     recentPoint,
					Status:          6,
					TransactionType: 9,
					Note:            "Auto Expired Eden Point",
					Merchant:        &model.Merchant{ID: v.ID},
					CreatedDate:     currentTime,
				}
				// insert new row for adding point back because of expired
				if _, e = ormWriter.Insert(mpl); e != nil {
					ormWriter.Rollback()
					return
				}

				merchant := &model.Merchant{ID: v.ID, TotalPoint: recentPoint}

				if _, e = ormWriter.Update(merchant, "TotalPoint"); e != nil {
					ormWriter.Rollback()
					return
				}

				if e = o.Raw("select exists(select id from merchant_point_summary mps where merchant_id = ? and summary_date = ?)", v.ID, currentTime.Format("2006-01-02")).QueryRow(&isExistMps); e != nil || (e == nil && !isExistMps) {
					if _, e = ormWriter.Raw("insert into merchant_point_summary (merchant_id, summary_date, redeemed_point) values (?, ?, ?)", v.ID, currentTime.Format("2006-01-02"), v.CurrentPeriodPoint).Exec(); e != nil {
						ormWriter.Rollback()
						return
					}
				} else {
					if _, e = ormWriter.Raw("update merchant_point_summary set redeemed_point = redeemed_point + ? where merchant_id = ? and summary_date = ?", v.CurrentPeriodPoint, v.ID, currentTime.Format("2006-01-02")).Exec(); e != nil {
						ormWriter.Rollback()
						return
					}
				}
			}

			merchantPointExpiration := &model.MerchantPointExpiration{
				ID:                 v.ID,
				CurrentPeriodPoint: v.NextPeriodPoint,
				NextPeriodPoint:    0,
				CurrentPeriodDate:  currentPeriodDate,
				NextPeriodDate:     nextPeriodDate,
				LastUpdatedAt:      time.Now(),
			}

			// Update merchant point expiration
			if e = merchantPointExpiration.Save(); e != nil {
				ormWriter.Rollback()
				return
			}

			if e = log.AuditLogByUser(&model.Staff{ID: 222}, v.ID, "merchant_point_expiration", "auto_expired", "Cronjob Auto Expired Eden Point"); e != nil {
				ormWriter.Rollback()
				return
			}

			ormWriter.Commit()
		}
	}

	return ctx.Serve(e)
}
