package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	util "git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/settlement_service"
	"google.golang.org/protobuf/types/known/timestamppb"

	re "regexp"
)

type ISalesOrderService interface {
	Create(ctx context.Context, req *dto.CreateRequestSalesOrder) (res dto.SalesOrderDetailResponse, err error)
	UpdateCOD(ctx context.Context, req *dto.UpdateCodRequest) (res dto.SalesOrderDetailResponse, err error)
	GetSalesOrderFeedback(ctx context.Context, req *dto.GetFeedback) (res []dto.SalesOrderFeedback, err error)
	CreateSalesOrderFeedback(ctx context.Context, req *dto.CreateSalesFeedback) (res dto.SalesOrderFeedback, err error)
}

type SalesOrderService struct {
	opt opt.Options
}

func NewSalesOrderService() ISalesOrderService {
	return &SalesOrderService{
		opt: global.Setup.Common,
	}
}

func (s *SalesOrderService) Create(ctx context.Context, req *dto.CreateRequestSalesOrder) (res dto.SalesOrderDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Create")
	defer span.End()
	var (
		//err                                          error
		//filter, exclude               map[string]interface{}
		//addressID, regionID,
		totalAcc int64

		//orderQty                      []float64
		totalQty                         float64
		addressID, regionID, itemListStr string
		//orderChannel                  int8
		//checkitem                                 *model.itemPush
		paymentMethodName string
	)
	//currentTime := time.Now()
	discAmount := 0.0
	itemList := make(map[int64]string)
	layout := "2006-01-02"

	//var soGrpc *sales_service.CreateSalesOrderRequest
	header := &sales_service.SalesOrder{}
	details := []*sales_service.SalesOrderItem{}

	customerID := req.Session.Customer.ID
	// if req.Data.AddressID != req.Session.Address.ID {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }
	// adressID, _ := strconv.Atoi(req.Data.AddressID)
	addressID = req.Data.AddressID
	archetypeID, _ := strconv.Atoi(req.Session.Address.ArchetypeID)
	// address, err := s.opt.Client.BridgeServiceGrpc.GetAddressDetail(ctx, &bridge_service.GetAddressDetailRequest{
	// 	Id: addressID,
	// })
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridge_service.GetAddressGPListRequest{
		CustomerNumber: req.Session.Customer.Code,
		Limit:          1,
		Offset:         0,
	})
	address1, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
		Id: req.Data.AddressID,
	})
	fmt.Print(customerID, addressID, address1)
	if len(address.Data) == 0 {
		return
	}
	req.Data.IsCreateMerchantVa = map[string]int8{"bca": 0, "permata": 0}
	if req.Session.Customer.Suspended == "1" {
		//o.Failure("customer.suspended", util.ErrorSuspended())
	}
	// if address.Data.Status != 1 {

	// }
	// admDivision, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridge_service.GetAdmDivisionDetailRequest{
	// 	Code: ,
	// })

	archetype, err := s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
		Id: req.Session.Address.ArchetypeID,
	})
	// region, err := s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridge_service.GetRegionDetailRequest{
	// 	Id: int64(admDivision.Data.RegionId),
	// })
	regionPolicy, err := s.opt.Client.ConfigurationServiceGrpc.GetRegionPolicyDetail(ctx, &configuration_service.GetRegionPolicyDetailRequest{
		RegionId: address.Data[0].AdministrativeDiv.GnlRegion,
	})

	// customerTypeID, _ := strconv.Atoi(req.Session.Customer.CustomerType)

	customerType, err := s.opt.Client.BridgeServiceGrpc.GetCustomerTypeDetail(ctx, &bridge_service.GetCustomerTypeDetailRequest{
		Id: int64(1),
	})
	// customerTypeGP, err := s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
	// 	Id: req.Session.Customer.CustomerType,
	// })
	if err != nil {
		err = errors.New("Customer Type invalid data.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	deliveryFee, err := s.opt.Client.BridgeServiceGrpc.GetDeliveryFeeGPList(ctx, &bridge_service.GetDeliveryFeeGPListRequest{
		Limit:     100,
		Offset:    0,
		GnlRegion: address.Data[0].AdministrativeDiv.GnlRegion,
	})
	if err != nil {
		err = errors.New("Region Business Policy invalid data.")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	var deliveryFees model.DeliveryFee
	var exist bool
	for _, v := range deliveryFee.Data {
		// if v.CustomerTypeId == customerTypeGP.Data[0].GnL_Cust_Type_ID {
		if v.GnlCustTypeId == strconv.Itoa(int(customerType.Data.Id)) {
			deliveryFees = model.DeliveryFee{
				// ID:          strconv.Itoa(int(v.Id)),
				MinOrder:    float64(v.Minorqty),
				DeliveryFee: v.GnlDeliveryFee,
			}
			exist = true
			break
		}
	}
	if !exist {
		for _, v := range deliveryFee.Data {
			if v.GnlCustTypeId == "" {
				deliveryFees = model.DeliveryFee{
					// ID:          strconv.Itoa(int(v.Id)),
					MinOrder:    float64(v.Minorqty),
					DeliveryFee: v.GnlDeliveryFee,
				}
				break
			}
		}
	}

	if req.Data.OrderTypeID == 0 {
		req.Data.OrderTypeID = 1
	}

	// Validation only order type regular or self pickup can created
	if req.Data.OrderTypeID != 1 && req.Data.OrderTypeID != 6 {
		// o.Failure("order_type_id.invalid", util.ErrorInvalidData("order type"))
	}

	// req.Data.OrderType = &model.OrderType{ID: req.Data.OrderTypeID}
	// if err = req.Data.OrderType.Read("ID"); err != nil {
	// 	// o.Failure("order_type_id.invalid", util.ErrorInvalidData("order type"))
	// }

	// if req.Data.Branch.Salesperson != nil {
	// 	req.Data.Salesperson = req.Data.Branch.Salesperson
	// 	req.Data.Salesperson.Read("ID")
	// }

	// site, err := s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridge_service.GetSiteDetailRequest{
	// 	// Id: address.Data[0].Locncode,
	// 	Id: 1,
	// })
	site, err := s.opt.Client.BridgeServiceGrpc.GetSiteGPList(ctx, &bridge_service.GetSiteGPListRequest{
		// Id: address.Data[0].Locncode,
		Limit:    100,
		Offset:   0,
		Locncode: address.Data[0].Locncode,
	})

	// salesterm, err := s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridge_service.GetSiteDetailRequest{
	// 	Id: address.Data.SiteId,
	// })
	// invoiceterm, err := s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridge_service.GetSiteDetailRequest{
	// 	Id: address.Data.SiteId,
	// })
	// paymentgroup, err := s.opt.Client.BridgeServiceGrpc.GetSiteDetail(ctx, &bridge_service.GetSiteDetailRequest{
	// 		Id: address.Data.SiteId,
	// 	})

	// Special Conditions if order type self pick up
	if req.Data.OrderTypeID == 6 {
		// // Set warehouse that available for self pickup based on area branch
		// if req.Data.Warehouse, err = repository.GetWarehouseSelfPickupByAreaID(req.Data.Branch.Area.ID); err != nil {
		// 	o.Failure("area_id.invalid", util.ErrorInvalidData("area"))
		// }
		// // Set default payment group to advance for self pickup
		// req.Data.PaymentGroup = &model.PaymentGroup{ID: 1}
		// if err = req.Data.PaymentGroup.Read("ID"); err != nil {
		// 	o.Failure("payment_group_id.invalid", util.ErrorInvalidData("payment group"))
		// }
		// // Set default sales term to PBD for self pickup
		// req.Data.SalesTerm = &model.SalesTerm{ID: 11}
		// if err = req.Data.SalesTerm.Read("ID"); err != nil {
		// 	o.Failure("sales_term_id.invalid", util.ErrorInvalidData("sales term"))
		// }
	}

	wrtID, _ := strconv.Atoi(req.Data.WrtID)
	// wrt, err := s.opt.Client.ConfigurationServiceGrpc.GetWrtDetail(ctx, &configuration_service.GetWrtDetailRequest{
	// 	Id: int64(wrtID),
	// })

	wrtIdGP, err := s.opt.Client.ConfigurationServiceGrpc.GetWrtIdGP(ctx, &configuration_service.GetWrtDetailRequest{
		Id: int64(wrtID),
	})

	// wrt, err := s.opt.Client.ConfigurationServiceGrpc.GetWrtList(ctx, &configuration_service.GetWrtListRequest{
	// 	// Id: int64(wrtID),
	// 	Limit:    1,
	// 	Offset:   0,
	// 	RegionId: req.Data.RegionID,
	// 	Search:   wrt,
	// })

	// Validation wrt for self pickup order type
	if wrtIdGP.Data.Type == 2 && req.Data.OrderTypeID != 6 {
		// o.Failure("wrt_id.invalid", util.ErrorOnlyValidFor("wrt", "pemesanan", "ambil sendiri"))
	}

	// Validation wrt for regular order type (delivery order)
	if wrtIdGP.Data.Type == 1 && req.Data.OrderTypeID != 1 {
		// o.Failure("wrt_id.invalid", util.ErrorOnlyValidFor("wrt", "pemesanan", "dengan pengiriman"))
	}

	// Validation for cannot use both redeem point and apply voucher
	if req.Data.RedeemPoint != 0 && req.Data.RedeemCode != "" {
		err = edenlabs.ErrorValidation("redeem_point", "Voucher dan Eden Point tidak dapat digunakan bersamaan")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// untuk check order time limit
	if req.Data.DeliveryDateStr != "" {
		if req.Data.DeliveryDate, err = time.Parse(layout, req.Data.DeliveryDateStr); err != nil {
			// o.Failure("delivery_date.invalid", util.ErrorInvalidData("delivery date"))
		} else {
			req.Data.DeliveryDate = time.Date(req.Data.DeliveryDate.Year(), req.Data.DeliveryDate.Month(), req.Data.DeliveryDate.Day(), 0, 0, 0, 0, time.Now().Local().Location())
			wib, _ := time.LoadLocation("Asia/Jakarta")
			currentTime := time.Now().In(wib)
			var weeklyDayOff int
			// ini check perbedaan tanggal nya brp lama misal tanggal 23 dan 24 maka perbedaannya menjadi totalTime = 1
			// totalTime := util.DaysBetween(util.Date(req.Data.DeliveryDate.Format("2006-01-02")), util.Date(currentTime.Format("2006-01-02")))

			// c.Session.Merchant.PaymentGroup.Read("ID")

			t1, _ := time.Parse("15:04:05", regionPolicy.Data.OrderTimeLimit+":00")
			t2, _ := time.Parse("15:04:05", currentTime.Format("15:04:05"))

			//check payment group data, if merchant is advance payment, check order time limit -30 min validation
			// if req.Session.Customer.PaymentGroup.Name == "Advance Payment" {
			// 	t1, _ := time.Parse("15:04:05", t1.Add(-time.Minute*30).Format("15:04:05"))

			// 	req.Data.Second = t1.Sub(t2).Seconds()
			// 	// if req.Data.Second < 0 && totalTime < 2 {
			// 	// 	o.Failure("delivery_date.invalid", util.ErrorOrderTimeLimit())
			// 	// }
			// } else {
			// 	req.Data.Second = t1.Sub(t2).Seconds()
			// 	// if req.Data.Second < 0 && totalTime < 2 {
			// 	// 	o.Failure("delivery_date.invalid", util.ErrorOrderTimeLimit())
			// 	// }
			// }
			fmt.Print(t1, t2)
			if req.Data.DeliveryDate.Before(time.Now().In(wib)) {
				// o.Failure("delivery_date", util.ErrorInvalidDeliveryDate())
			}

			// in go sunday is 0, replace manual 7 to 0
			if regionPolicy.Data.WeeklyDayOff == 7 {
				weeklyDayOff = 0
			} else {
				weeklyDayOff = int(regionPolicy.Data.WeeklyDayOff)
			}

			if int(req.Data.DeliveryDate.Weekday()) == weeklyDayOff {
				// o.Failure("delivery_date", util.ErrorInvalidDeliveryDate())
			}

			// orm.Raw(
			// 	"SELECT * FROM day_off"+
			// 		" WHERE off_date = ? LIMIT 1",
			// 	req.Data.DeliveryDate.Format("2006-01-02")).QueryRow(&req.Data.DayOff)
			// if req.Data.DayOff != nil {
			// 	o.Failure("delivery_date", util.ErrorInvalidDeliveryDate())
			// }

		}
	}

	// glossaryOC, err := s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
	// 	Attribute: "order_channel",
	// 	ValueName: req.Platform,
	// })

	if len(req.Data.Items) > 0 {
		for i, v := range req.Data.Items {
			var (
				orderMaxQty                 float64
				maxDayDeliveryDate          int64
				maxAvailableDate            time.Time
				dayCount                    int8
				parentName, grandparentName string
			)

			if v.Quantity <= 0 {
				// o.Failure("qty"+strconv.Itoa(i)+".greater", util.ErrorGreater("product quantity", "0"))
			} else {

				if v.UnitPrice < 0 {
					// o.Failure("unit_price"+strconv.Itoa(i)+".equalorgreater", util.ErrorEqualGreater("product unit price", "0"))
				}

				if v.ItemID == "" {
					// o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInputRequired("product"))
				} else {
					if itemID, err := strconv.Atoi(v.ItemID); err == nil {
						itemDetail, err := s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
							Id: utils.ToString(itemID),
						})
						if err != nil {
							fmt.Print(err)
						}
						uom, er := s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridge_service.GetUomDetailRequest{
							Id: utils.ToInt64(itemDetail.Data.UomId),
						})
						if er != nil {
							fmt.Println(er)
						}
						category, e := s.opt.Client.CatalogServiceGrpc.GetItemCategoryDetail(ctx, &catalog_service.GetItemCategoryDetailRequest{
							// Id: itemDetail.Data.ItemCategoryId[0],
							Id: 1,
						})
						if e != nil {
							fmt.Println(e)
						}
						// v.Price = &model.Price{Product: itemDetail.Data, PriceSet: req.Data.Branch.PriceSet}
						// v.Price.Read("Product", "PriceSet")
						// itemDetail.Data.Category.Read("ID")

						// if itemDetail.Data.Category.ParentID != 0 {
						// 	itemDetail.Data.Category.Parent = &model.Category{ID: itemDetail.Data.Category.ParentID}
						// 	itemDetail.Data.Category.Parent.Read("ID")
						// 	parentName = itemDetail.Data.Category.Parent.Name
						// }
						// if itemDetail.Data.Category.GrandParentID != 0 {
						// 	itemDetail.Data.Category.GrandParent = &model.Category{ID: itemDetail.Data.Category.GrandParentID}
						// 	itemDetail.Data.Category.GrandParent.Read("ID")
						// 	grandparentName = itemDetail.Data.Category.GrandParent.Name
						// }

						// price, err := s.opt.Client.BridgeServiceGrpc.getprice(ctx, &bridge_service.GetUomDetailRequest{
						// 	Id: itemDetail.Data.UomId,
						// })
						if itemDetail.Data.OrderChannelRestriction != "" {
							orderChannel := []string{}
							if req.Platform == "orca" {
								orderChannel = append(orderChannel, "2")
							}
							if req.Platform == "mantis" {
								orderChannel = append(orderChannel, "3")
							}
							matched, err := re.MatchString(strings.Join(orderChannel[:], "|"), itemDetail.Data.OrderChannelRestriction)
							if err != nil {
								// o.Failure("product_id_exclude["+itemDetail.DataID+"].invalid", util.ErrorInvalidData("order channel restriction"))
							}
							if matched {
								// o.Failure("product_id_exclude["+itemDetail.DataID+"].invalid", util.ErrorOrderChannelRestriction())
							}
						}

						//get price from price set
						//dummy
						req.Data.Items[i].UnitPrice = 5000
						// req.Data.Items[i].UnitPrice = v.Price.UnitPrice

						v.Subtotal = v.Quantity * float64(5000)
						// v.Subtotal = v.Quantity * float64(v.Price.UnitPrice)
						if _, exist := itemList[int64(itemID)]; exist {
							// o.Failure("product_id.duplicate", util.ErrorDuplicate("produk"))
						} else {
							itemList[int64(itemID)] = "t"
						}

						//Wfilter = map[string]interface{}{"item_id": itemID, "site_id": site.Data.Id, "salable": 1}
						// if _, countStock, err := repository.CheckStockData(filter, exclude); err == nil && countStock == 0 {
						// 	o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorProductMustAvailable())
						// }

						if strings.Contains(","+itemDetail.Data.ExcludeArchetype+",", ","+strconv.Itoa(int(archetypeID))+",") {
							// o.Failure("product_id_exclude["+itemDetail.DataID+"].invalid", "Produk ini tidak dapat dibeli")
							// productListStr += itemDetail.Data.Name + ", "
							continue
						}

						if itemDetail.Data.MaxDayDeliveryDate > 0 && int64(itemDetail.Data.MaxDayDeliveryDate) < int64(regionPolicy.Data.MaxDayDeliveryDate) {
							maxDayDeliveryDate = int64(itemDetail.Data.MaxDayDeliveryDate)
						} else {
							maxDayDeliveryDate = int64(regionPolicy.Data.MaxDayDeliveryDate)
						}

						// get maximum available date based on max day delivery date for each product or max day delivery date on area policy
						for dayCount, maxAvailableDate = 0, time.Now(); ; maxAvailableDate = maxAvailableDate.AddDate(0, 0, 1) {
							if (int(maxAvailableDate.Weekday()) == 0 && int(regionPolicy.Data.WeeklyDayOff) == 7) || (int(maxAvailableDate.Weekday()) == int(regionPolicy.Data.WeeklyDayOff)) {
								continue
							}

							// if err = orm.Raw("SELECT * FROM day_off WHERE off_date = ? LIMIT 1", maxAvailableDate.Format("2006-01-02")).QueryRow(&req.Data.DayOff); err == nil && req.Data.DayOff.ID != 0 {
							// 	continue
							// }

							dayCount++
							if int64(dayCount) > maxDayDeliveryDate {
								break
							}
						}
						maxAvailableDate = time.Date(maxAvailableDate.Year(), maxAvailableDate.Month(), maxAvailableDate.Day(), 0, 0, 0, 0, time.Local)
						if maxAvailableDate.Before(req.Data.DeliveryDate) {
							// o.Failure("product_id_delivery["+itemDetail.DataID+"].invalid", "Maks. pengiriman H+"+strconv.Itoa(int(maxDayDeliveryDate)))
							// o.Failure("product_id.invalid", "Produk kamu ada yang melebihi batas waktu pengiriman yang diperbolehkan")
						}

						salesOrderItem, err := s.opt.Client.SalesServiceGrpc.GetSalesOrderItemDetail(ctx, &sales_service.GetSalesOrderItemDetailRequest{
							//isi delivery date,product id,customer id
						})
						// if _, err = orm.Raw("select soi.order_qty "+
						// 	"from sales_order so "+
						// 	"join sales_order_item soi on so.id = soi.sales_order_id "+
						// 	"join branch b on so.branch_id = b.id "+
						// 	"where so.delivery_date = ? and so.status not in (3,4) and soi.product_id = ? and b.merchant_id = ?", req.Data.DeliveryDate.Format("2006-01-02"), itemDetail.Data.ID, req.Data.Branch.Merchant.ID).QueryRows(&orderQty); err == nil {
						// 	for _, v := range orderQty {
						// 		totalQty = totalQty + v
						// 	}

						// 	if err = orm.Raw("select p.order_max_qty "+
						// 		"from product p "+
						// 		"where p.id = ?", itemDetail.Data.ID).QueryRow(&orderMaxQty); err == nil {
						// 		if orderMaxQty > 0 && totalQty+v.Quantity > orderMaxQty {
						// 			o.Failure("product_id_qty["+itemDetail.DataID+"].invalid", "Maks. "+strconv.Itoa(int(itemDetail.Data.OrderMaxQty))+" "+itemDetail.Data.Uom.Name+" per hari")
						// 			o.Failure("product_id.invalid", "Produk kamu ada yang melebihi batas maksimum pembelian")
						// 		}
						// 	}
						// }

						//prevent modulus a float number
						a := 10000.00
						mod := math.Mod(v.Quantity*a, itemDetail.Data.OrderMinQty*a)

						//check is quantity under minimum order quantity
						if v.Quantity < itemDetail.Data.OrderMinQty {
							// o.Failure("product_id_qty["+itemDetail.DataID+"].invalid", "Minimum "+fmt.Sprintf("%v", itemDetail.Data.OrderMinQty)+" "+itemDetail.Data.Uom.Name)
						}

						//check if the input is input the multiple min quantity
						if mod != 0 {
							// o.Failure("product_id_qty["+itemDetail.DataID+"].invalid", "Hanya bisa kelipatan "+fmt.Sprintf("%v", itemDetail.Data.OrderMinQty)+" "+itemDetail.Data.Uom.Name)
						}

						// // start sku discount
						// if v.SkuDiscountItem, err = repository.GetSkuDiscountData(c.Session.Merchant.ID, req.Data.Branch.PriceSet.ID, itemDetail.Data.ID, 0, orderChannel, currentTime); err == nil {
						// 	var maxDiscQty float64

						// 	orm.LoadRelated(v.SkuDiscountItem, "SkuDiscountItemTiers")
						// 	for _, val := range v.SkuDiscountItem.SkuDiscountItemTiers {
						// 		if v.Quantity < val.MinimumQty {
						// 			break
						// 		} else {
						// 			v.UnitPriceDiscount = val.Amount
						// 		}
						// 	}

						// 	if v.UnitPriceDiscount > 0 {
						// 		maxDiscQty = v.SkuDiscountItem.RemOverallQuota

						// 		if maxDiscQty <= 0 {
						// 			o.Failure("sku_discount.invalid", util.ErrorHasChange("Promo"))
						// 			continue
						// 		}

						// 		if maxDiscQty > v.SkuDiscountItem.RemQuotaPerUser {
						// 			maxDiscQty = v.SkuDiscountItem.RemQuotaPerUser
						// 		}

						// 		if maxDiscQty <= 0 {
						// 			o.Failure("sku_discount.invalid", util.ErrorHasChange("Promo"))
						// 			continue
						// 		}

						// 		if maxDiscQty > v.SkuDiscountItem.RemDailyQuotaPerUser {
						// 			maxDiscQty = v.SkuDiscountItem.RemDailyQuotaPerUser
						// 		}

						// 		if maxDiscQty <= 0 {
						// 			o.Failure("sku_discount.invalid", util.ErrorHasChange("Promo"))
						// 			continue
						// 		}

						// 		if maxDiscQty < v.DiscQty || (maxDiscQty > 0 && v.DiscQty == 0) {
						// 			o.Failure("sku_discount.invalid", util.ErrorHasChange("Promo"))
						// 			continue
						// 		}

						// 		v.DiscAmount = v.DiscQty * v.UnitPriceDiscount
						// 		req.Data.TotalSkuDiscount += v.DiscAmount

						// 		v.Subtotal -= v.DiscAmount

						// 		v.IsUseSkuDiscount = 1
						// 	}
						// } else {
						// 	if v.DiscQty > 0 {
						// 		o.Failure("sku_discount.invalid", util.ErrorHasChange("Promo"))
						// 		continue
						// 	}
						// }
						// // end sku discount

						v.Weight = v.Quantity * itemDetail.Data.UnitWeightConversion
						req.Data.TotalPrice = req.Data.TotalPrice + v.Subtotal
						req.Data.TotalWeight = req.Data.TotalWeight + v.Weight

						// start setup item data to be sent to talon
						unitPrice := v.UnitPrice
						if v.IsUseSkuDiscount == 1 {
							unitPrice = v.UnitPrice - v.UnitPriceDiscount
						}

						detail := &sales_service.SalesOrderItem{
							ItemId:             itemDetail.Data.Id,
							UnitPrice:          unitPrice,
							OrderQty:           v.Quantity,
							Weight:             v.Weight,
							Subtotal:           v.Subtotal,
							Note:               v.Note,
							ShadowPrice:        5000,
							TaxableItem:        int32(v.TaxableItem),
							TaxPercentage:      v.TaxPercentage,
							DiscountQty:        v.DiscQty,
							UnitPriceDiscount:  v.UnitPriceDiscount,
							ItemDiscountAmount: v.DiscAmount,
							ItemIdGp:           itemDetail.Data.Code,
							// ItemDiscountId: 1,
						}
						// itemData := &model.SessionItemData{
						// 	ProductName:  itemDetail.Data.Name,
						// 	ProductCode:  itemDetail.Data.Code,
						// 	CategoryName: itemDetail.Data.Category.Name,
						// 	UnitPrice:    unitPrice,
						// 	OrderQty:     1,
						// 	UnitWeight:   v.Quantity,
						// 	Attributes: map[string]string{
						// 		"parent_category":       parentName,
						// 		"grand_parent_category": grandparentName,
						// 	},
						// }
						details = append(details, detail)
						// req.Data.ItemList = append(req.Data.ItemList, itemData)
						// itemData = nil
						// end setup item data to be sent to talon

						// //check push product
						// itemDetail.DataPush = 2
						// if e := orm.Raw("SELECT * FROM product_push pp "+
						// 	"WHERE pp.product_id = ? AND ? >= pp.start_date AND pp.area_id = ? AND pp.archetype_id = ? AND pp.status = 1",
						// 	productID, currentTime.Format(layout), req.Data.Branch.Area.ID, req.Data.Branch.Archetype.ID).QueryRow(&checkProduct); e != nil {
						// 	continue
						// }

						// if checkProduct != nil && currentTime.Unix() >= checkProduct.StartDate.Unix() {
						// 	itemDetail.DataPush = checkProduct.Status
						// }

						fmt.Println(uom, err, category, salesOrderItem, orderMaxQty, parentName, grandparentName)

						// if v.IsUseSkuDiscount == 1 {
						// 	v.SkuDiscountItem.RemOverallQuota = v.SkuDiscountItem.RemOverallQuota - v.DiscQty
						// 	if v.SkuDiscountItem.IsUseBudget == 1 {
						// 		v.SkuDiscountItem.RemBudget = v.SkuDiscountItem.RemBudget - v.DiscAmount
						// 	}
						// 	// update remaining quota and budget sku discount item
						// 	if _, e = o.Update(v.SkuDiscountItem, "RemOverallQuota", "RemBudget"); e != nil {
						// 		goto CANCEL
						// 	}

						// 	sdl = &model.SkuDiscountLog{
						// 		SalesOrderItem:  soi,
						// 		SkuDiscount:     v.SkuDiscountItem.SkuDiscount,
						// 		SkuDiscountItem: v.SkuDiscountItem,
						// 		Merchant:        r.Data.Branch.Merchant,
						// 		Branch:          r.Data.Branch,
						// 		Product:         v.Product,
						// 		DiscountQty:     v.DiscQty,
						// 		DiscountAmount:  v.DiscAmount,
						// 		CreatedAt:       time.Now(),
						// 		Status:          1,
						// 	}

						// 	if _, e = o.Insert(sdl); e != nil {
						// 		goto CANCEL
						// 	}
						// }
					} else {
						// o.Failure("product_id"+strconv.Itoa(i)+".invalid", util.ErrorInvalidData("product"))
					}

				}
			}
		}

		if itemListStr != "" {
			itemListStr = strings.TrimSuffix(itemListStr, ", ")
			// o.Failure("product_id.invalid", "Kamu tidak dapat membeli produk "+productListStr)
		}
		req.Data.DeliveryFee = deliveryFees.DeliveryFee
		if req.Data.TotalPrice >= deliveryFees.MinOrder || req.Data.OrderTypeID == 6 {
			req.Data.DeliveryFee = 0
		}

	}
	req.Data.TotalCharge = float64(req.Data.TotalPrice) + req.Data.DeliveryFee

	//exclude = map[string]interface{}{"account_number": " "}
	// filter = map[string]interface{}{"customer_id": req.Session.Customer.ID, "payment_channel_id": 6, "account_number__isnull": ""}
	// _, totalAcc, err = repository.CheckMerchantAccNumData(filter, exclude)
	if totalAcc == 0 {
		req.Data.IsCreateMerchantVa["bca"] = 1
	}

	// filter = map[string]interface{}{"customer_id": req.Session.Customer.ID, "payment_channel_id": 7, "account_number__isnull": ""}
	// _, totalAcc, err = repository.CheckMerchantAccNumData(filter, exclude)
	if totalAcc == 0 {
		req.Data.IsCreateMerchantVa["permata"] = 1
	}
	// orm.Raw(
	// 	"SELECT * FROM day_off"+
	// 		" WHERE off_date = ? LIMIT 1",
	// 	req.Data.DeliveryDate.Format("2006-01-02")).QueryRow(&req.Data.DayOff)

	// exclude = map[string]interface{}{}
	if req.Data.RedeemCode != "" {
		// voucher, err := s.opt.Client.CampaignServiceGrpc.getvou(ctx, &bridge_service.GetArchetypeDetailRequest{
		// 	Id: int64(archetypeID),
		// })
		// req.Data.Voucher = &model.Voucher{RedeemCode: req.Data.RedeemCode, Status: 1}
		// if err = req.Data.Voucher.Read("RedeemCode", "Status"); err == nil {
		// 	if req.Data.Voucher.Status != 1 {
		// 		o.Failure("redeem_code.inactive", util.ErrorActive("voucher"))
		// 		return o
		// 	}

		// 	if currentTime.Before(req.Data.Voucher.StartTimestamp) {
		// 		o.Failure("redeem_code.invalid", util.ErrorNotInPeriod("voucher"))
		// 		return o
		// 	}

		// 	if currentTime.After(req.Data.Voucher.EndTimestamp) {
		// 		o.Failure("redeem_code.invalid", util.ErrorOutOfPeriod("voucher"))
		// 		return o
		// 	}

		// 	if req.Data.Voucher.Type == 1 { //type total discount
		// 		if req.Data.Voucher.MinOrder > req.Data.TotalPrice {
		// 			o.Failure("redeem_code.invalid", util.ErrorEqualGreater("total order", "minimum order"))
		// 			return o
		// 		}

		// 		if req.Data.Voucher.VoucherAmount > req.Data.TotalPrice {
		// 			o.Failure("redeem_code.invalid", util.ErrorEqualGreater("total order", "discount amount"))
		// 			return o
		// 		}
		// 	} else if req.Data.Voucher.Type == 2 { // type grand total discount
		// 		if req.Data.Voucher.MinOrder > req.Data.TotalCharge {
		// 			o.Failure("redeem_code.invalid", util.ErrorEqualGreater("grand total order", "minimum order"))
		// 			return o
		// 		}

		// 		if req.Data.Voucher.VoucherAmount > req.Data.TotalCharge {
		// 			o.Failure("redeem_code.invalid", util.ErrorEqualGreater("grand total order", "discount amount"))
		// 			return o
		// 		}
		// 	} else if req.Data.Voucher.Type == 3 { // type delivery discount
		// 		if req.Data.Voucher.MinOrder > req.Data.TotalPrice {
		// 			o.Failure("redeem_code.invalid", util.ErrorEqualGreater("total order", "minimum order"))
		// 			return o
		// 		}

		// 		if req.Data.Voucher.VoucherAmount > req.Data.DeliveryFee {
		// 			o.Failure("redeem_code.invalid", util.ErrorEqualGreater("delivery fee", "discount amount"))
		// 			return o
		// 		}
		// 	} else if req.Data.Voucher.Type == 4 { // type extra edenpoint
		// 		if req.Data.Voucher.MinOrder > req.Data.TotalPrice {
		// 			o.Failure("redeem_code.invalid", util.ErrorEqualGreater("total order", "minimum order"))
		// 			return o
		// 		}
		// 	}

		// 	if req.Data.Voucher.RemOverallQuota < 1 {
		// 		o.Failure("redeem_code.invalid", util.ErrorFullyUsed("voucher"))
		// 		return o
		// 	}

		// 	filter = map[string]interface{}{"merchant_id": c.Session.Merchant.ID, "voucher_id": req.Data.Voucher.ID, "status": 1}
		// 	if _, countVoucherLog, err := repository.CheckVoucherLogData(filter, exclude); err == nil && countVoucherLog >= req.Data.Voucher.UserQuota {
		// 		o.Failure("redeem_code.invalid", util.ErrorFullyUsed("voucher"))
		// 		return o
		// 	}

		// 	for _, v := range strings.Split(c.Session.Merchant.TagCustomer, ",") {
		// 		if strings.Contains(req.Data.Voucher.TagCustomer, v) {
		// 			tagCustomerID, _ := strconv.Atoi(v)
		// 			tagCustomer := &model.TagCustomer{ID: int64(tagCustomerID)}
		// 			tagCustomer.Read("ID")

		// 			req.Data.SameTagCustomer = req.Data.SameTagCustomer + "," + tagCustomer.Name
		// 		}
		// 	}

		// 	if req.Data.Voucher.TagCustomer != "" {
		// 		sameTagCustomer := ""
		// 		for _, v := range strings.Split(c.Session.Merchant.TagCustomer, ",") {
		// 			if strings.Contains(req.Data.Voucher.TagCustomer, v) {
		// 				sameTagCustomer = sameTagCustomer + "," + v
		// 			}
		// 		}

		// 		sameTagCustomer = strings.Trim(sameTagCustomer, ",")
		// 		if sameTagCustomer == "" {
		// 			o.Failure("redeem_code.invalid", util.ErrorNotValidFor("Voucher", "customer tag"))
		// 		}

		// 		req.Data.SameTagCustomer = strings.Trim(req.Data.SameTagCustomer, ",")
		// 		if req.Data.SameTagCustomer == "" {
		// 			o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "customer tag"))
		// 			return o
		// 		}
		// 	}

		// 	if req.Data.Voucher.Area.ID != 1 && req.Data.Branch.Area.ID != req.Data.Voucher.Area.ID {
		// 		o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "area"))
		// 		return o
		// 	}

		// 	err = req.Data.Voucher.Archetype.Read("ID")
		// 	if err != nil {
		// 		o.Failure("redeem_code.invalid", util.ErrorNotFound("Archetype"))
		// 		return o
		// 	}

		// 	err = req.Data.Voucher.Archetype.BusinessType.Read("ID")
		// 	if err != nil {
		// 		o.Failure("redeem_code.invalid", util.ErrorNotFound("Business type"))
		// 		return o
		// 	}

		// 	if req.Data.Voucher.Archetype.BusinessType.AuxData != 1 {
		// 		if req.Data.Voucher.Archetype.AuxData != 1 {
		// 			if req.Data.Branch.Archetype.ID != req.Data.Voucher.Archetype.ID {
		// 				o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "archetype"))
		// 				return o
		// 			}
		// 		} else {
		// 			if req.Data.Branch.Archetype.BusinessType.ID != req.Data.Voucher.Archetype.BusinessType.ID {
		// 				o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "business type"))
		// 				return o
		// 			}
		// 		}
		// 	}

		// 	// start voucher membership validation
		// 	if req.Data.Voucher.MembershipLevelID != 0 {
		// 		c.Session.Merchant.MembershipLevel = &model.MembershipLevel{ID: c.Session.Merchant.MembershipLevelID}
		// 		if err = c.Session.Merchant.MembershipLevel.Read("ID"); err != nil {
		// 			o.Failure("redeem_code.invalid", util.ErrorInvalidData("voucher"))
		// 			return o
		// 		}

		// 		c.Session.Merchant.MembershipCheckpoint = &model.MembershipCheckpoint{ID: c.Session.Merchant.MembershipCheckpointID}
		// 		if err = c.Session.Merchant.MembershipCheckpoint.Read("ID"); err != nil {
		// 			o.Failure("redeem_code.invalid", util.ErrorInvalidData("voucher"))
		// 			return o
		// 		}

		// 		req.Data.Voucher.MembershipLevel = &model.MembershipLevel{ID: req.Data.Voucher.MembershipLevelID}
		// 		if err = req.Data.Voucher.MembershipLevel.Read("ID"); err != nil {
		// 			o.Failure("redeem_code.invalid", util.ErrorInvalidData("voucher"))
		// 			return o
		// 		}

		// 		req.Data.Voucher.MembershipCheckpoint = &model.MembershipCheckpoint{ID: req.Data.Voucher.MembershipCheckpointID}
		// 		if err = req.Data.Voucher.MembershipCheckpoint.Read("ID"); err != nil {
		// 			o.Failure("redeem_code.invalid", util.ErrorInvalidData("voucher"))
		// 			return o
		// 		}

		// 		if c.Session.Merchant.MembershipLevel.Level < req.Data.Voucher.MembershipLevel.Level {
		// 			o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "level "+c.Session.Merchant.MembershipLevel.Name))
		// 			return o
		// 		}

		// 		if c.Session.Merchant.MembershipCheckpoint.Checkpoint < req.Data.Voucher.MembershipCheckpoint.Checkpoint {
		// 			o.Failure("redeem_code.invalid", util.ErrorNotValidFor("voucher", "level "+c.Session.Merchant.MembershipLevel.Name+" lapak "+strconv.Itoa(int(c.Session.Merchant.MembershipCheckpoint.Checkpoint))))
		// 			return o
		// 		}
		// 	}
		// 	// end voucher membership validation

		// 	// set discount amount if type is not 4
		// 	if req.Data.Voucher.Type != 4 {
		// 		discAmount = req.Data.Voucher.VoucherAmount
		// 	}
		// } else {
		// 	o.Failure("redeem_code.invalid", util.ErrorNotFound("voucher"))
		// 	return o
		// }

		// add discount to total charge
		req.Data.TotalCharge = req.Data.TotalCharge - discAmount
	}

	// start set up customer profile talon
	// if err = c.Session.Merchant.FinanceArea.Read("ID"); err != nil {
	// 	o.Failure("branch.invalid", util.ErrorInvalidData("Merchant area"))
	// 	return o
	// }

	// if c.Session.Merchant.PaymentMethod != nil {
	// 	c.Session.Merchant.PaymentMethod.Read("ID")
	// }
	if req.Session.Customer.ProfileCode == "" {
		req.Data.IsInitCustomerTalonPoints = true
	}
	req.Session.Customer.ProfileCode = req.Session.Customer.Code
	// end set up customer profile talon

	// start set up advocate customer profile talon, if customer have referrer
	if req.Session.Customer.ReferrerCode != "" {
		customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerDetail(ctx, &bridge_service.GetCustomerDetailRequest{
			//referal code
		})
		fmt.Println(customer, err)
		// if err = req.Session.Customer.ReferrerCustomer.FinanceArea.Read("ID"); err != nil {
		// 	o.Failure("branch.invalid", util.ErrorInvalidData("Merchant area"))
		// 	return o
		// }

		// if err = req.Session.Customer.ReferrerCustomer.CustomerTypeID.Read("ID"); err != nil {
		// 	o.Failure("branch.invalid", util.ErrorInvalidData("Merchant business type"))
		// 	return o
		// }
		// if c.Session.Merchant.ReferrerMerchant.PaymentMethod != nil {
		// 	c.Session.Merchant.ReferrerMerchant.PaymentMethod.Read("ID")
		// }
		// if c.Session.Merchant.ReferrerMerchant.ProfileCode == "" {
		// 	req.Data.IsInitReferrerTalonPoints = true
		// }
		// c.Session.Merchant.ReferrerMerchant.ProfileCode = c.Session.Merchant.ReferrerMerchant.Code
		// req.Data.ReferrerData = []string{c.Session.Merchant.ReferrerMerchant.ProfileCode, c.Session.Merchant.ReferrerMerchant.ReferralCode}
	}
	// end set up advocate customer profile talon

	// start set up customer session talon
	req.Data.IntegrationCode = strings.ReplaceAll(time.Now().Format("20060102150405.99"), ".", "") + req.Session.Customer.Code
	// end set up customer session talon

	//validasi untuk redeem point
	if req.Data.RedeemPoint != 0 {
		// var isCanRedeem bool
		// if err = orm.Raw("select exists(select id from config_app where `attribute` = 'redeem_edenpoint_business_type' and find_in_set(?, value) > 0)", req.Session.Customer.CustomerType).QueryRow(&isCanRedeem); err != nil || (err == nil && !isCanRedeem) {
		// 	// o.Failure("redeem_point.invalid", util.ErrorRedeemUnallowed())
		// 	// return o
		// }

		if req.Data.RedeemPoint <= 0 {
			// o.Failure("redeem_point.invalid", util.ErrorGreater("Point", "0"))
			// return o
		}

		var IDPointLog int64

		var mPoint float64

		// orm.Raw("SELECT id from merchant_point_log mpl where status = 2 AND created_date = CURRENT_DATE AND merchant_id = ?;", c.Session.Merchant.ID).QueryRow(&IDPointLog)

		if IDPointLog != 0 {
			// o.Failure("redeem_point.invalid", "Penggunaan point sudah mencapai batas maksimum.")
			// return o
		}

		// orm.Raw("SELECT total_point from merchant where id = ?;", c.Session.Merchant.ID).QueryRow(&mPoint)

		if mPoint < req.Data.RedeemPoint {

			// o.Failure("redeem_point.invalid", "Point tidak mencukupi.")
			// return o
		}

		if req.Data.RedeemPoint > req.Data.TotalCharge*0.5 {
			// o.Failure("redeem_point.invalid", "Penggunaan point tidak sesuai syarat dan ketentuan")
			// return o
		}

		// orm.Raw("select recent_point from merchant_point_log where merchant_id = ? order by id desc limit 1;", c.Session.Merchant.ID).QueryRow(&req.Data.RecentPoint)

		// add discount to total charge
		req.Data.TotalCharge = req.Data.TotalCharge - req.Data.RedeemPoint

		// req.Data.MerchantPointExpiration = &model.MerchantPointExpiration{ID: c.Session.Merchant.ID}
		// if err = req.Data.MerchantPointExpiration.Read("ID"); err != nil {
		// 	o.Failure("merchant_point_expiration.invalid", util.ErrorInvalidData("merchant point expiration"))
		// 	return o
		// }

		req.Data.CurrentPointUsed = req.Data.RedeemPoint
		// Calculate to get current period point and next period point
		// req.Data.CurrentPeriodPoint = req.Data.MerchantPointExpiration.CurrentPeriodPoint - req.Data.RedeemPoint
		// req.Data.NextPeriodPoint = req.Data.MerchantPointExpiration.NextPeriodPoint

		// if current period point didn't can to cover total point redeem, use next period point
		if req.Data.CurrentPeriodPoint < 0 {
			// req.Data.NextPointUsed -= req.Data.CurrentPeriodPoint
			// req.Data.NextPeriodPoint = req.Data.MerchantPointExpiration.NextPeriodPoint + req.Data.CurrentPeriodPoint
			// req.Data.CurrentPeriodPoint = 0
			// req.Data.CurrentPointUsed = req.Data.MerchantPointExpiration.CurrentPeriodPoint
		}
	}

	if req.Data.TotalCharge < 10000 {
		// o.Failure("grand_total.invalid", util.ErrorEqualGreater("grand total", "10000"))
	}

	// Validate total charge not greater than credit limit remaining
	remainingCreditLimitAmount, _ := strconv.ParseFloat(req.Session.Customer.RemainingCreditLimitAmount, 64)

	req.Data.CreditLimitBefore = remainingCreditLimitAmount
	if req.Session.Customer.CreditLimitAmount == "0" {
		if req.Data.CreditLimitBefore < 0 {
			// o.Failure("credit_limit.invalid", util.ErrorCreditLimitExceeded())
		}
	} else {
		req.Data.CreditLimitAfter = req.Data.CreditLimitBefore - req.Data.TotalCharge
		if req.Data.CreditLimitAfter < 0 {
			// o.Failure("credit_limit.invalid", util.ErrorCreditLimitExceeded())
		}
		req.Data.IsCreateCreditLimitLog = true
	}

	header = &sales_service.SalesOrder{
		SalesOrderNumber: "",
		DocNumber:        "",
		AddressId:        addressID,
		CustomerId:       customerID,
		SalespersonId:    0,
		WrtId:            req.Data.WrtID,
		Application:      0,
		Status:           1,
		OrderTypeId:      req.Data.OrderTypeID,
		OrderDate:        timestamppb.New(time.Now()),
		Total:            totalQty,
		CreatedBy:        0,
		CreatedAt:        timestamppb.New(time.Now()),
		SiteId:           "1",
		SiteCode:         site.Data[0].Locncode,
		SalesGroupId:     0,
		// SubDistrictId:       address.Data[0].AdministrativeDiv.GnlSubdistrict,
		RegionId:            regionID,
		VoucherId:           0,
		PriceLevelId:        0,
		PaymentGroupSlsId:   0,
		ArchetypeId:         "0", //archetype.Data[0].GnL_Archetype_ID,
		OrderTypeSlsId:      0,
		IntegrationCode:     "",
		RecognitionDate:     timestamppb.New(time.Now()),
		RequestDeliveryDate: timestamppb.New(req.Data.DeliveryDate),
		BillingAddress:      req.Session.Customer.BillingAddress,
		// ShippingAddress:     address.Data.ShippingAddress,
		ShippingAddressNote: "",
		DeliveryFee:         deliveryFees.DeliveryFee,
		VouRedeemCode:       "",
		VouDiscAmount:       discAmount,
		PointRedeemAmount:   req.Data.RedeemPoint,
		PointRedeemId:       0,
		EdenPointCampaignId: 0,
		TotalSkuDiscAmount:  discAmount,
		TotalPrice:          req.Data.TotalPrice,
		TotalCharge:         req.Data.TotalCharge,
		TotalWeight:         req.Data.TotalWeight,
		Note:                "",
		PaymentReminder:     0,
		CancelType:          0,
	}

	//price level here

	// res.AddressID = req.Session.Address.ID
	// res.Type = req.Data.Type
	// res.Category = req.Data.Category
	soNumber := ""
	totCharge := 0.00
	// TODO: add validation order time limit
	if req.Data.Payment != nil {

		if req.Data.Payment.PaymentMethod == "PBD" {
			so, err := s.opt.Client.SalesServiceGrpc.CreateSalesOrder(ctx, &sales_service.CreateSalesOrderRequest{
				Data:     header,
				Dataitem: details,
			})
			fmt.Println(so, err)
			salesInvoiceExternalXendit, err := s.opt.Client.SettlementGrpc.CreateSalesInvoiceExternal(ctx, &settlement_service.CreateSalesInvoiceExternalRequest{
				SalesOrderId:   so.Data.Id,
				Email:          "Testing@edenfarm.id",
				OrderTimeLimit: regionPolicy.Data.OrderTimeLimit,
				DeliveryDate:   req.Data.DeliveryDate.Format("2006-01-02"),
				SalesOrderCode: so.Data.SalesOrderNumber,
				TotalCharge:    so.Data.TotalCharge,
				PaymentMethod:  req.Data.Payment.PaymentChannel,
			})

			if err != nil {
				fmt.Println("=====salesInvoiceExternalXendit========ERROR=====", err)

			}
			fmt.Println("=====salesInvoiceExternalXendit========", salesInvoiceExternalXendit)
			soNumber = so.Data.SalesOrderNumber
			totCharge = so.Data.TotalCharge
		} else {
			tempDetails := []*sales_service.CreateSalesOrderGPRequest_DetailItem{}

			// var tempDetails *[]sales_service.CreateSalesOrderGPRequest_DetailItem
			for _, v := range details {
				item, _ := s.opt.Client.CatalogServiceGrpc.GetItemDetail(ctx, &catalog_service.GetItemDetailRequest{
					Id: v.ItemIdGp,
				})
				detail := &sales_service.CreateSalesOrderGPRequest_DetailItem{
					// Sopnumbe: ,
					// Itemnmbr: v.ItemIdGp,
					Itemnmbr: "100XLG",
					// Itemdesc: v.,
					Quantity:   0,
					Unitprce:   0,
					Xtndprce:   0,
					GnL_Weight: int32(v.Weight),
					Pricelvl:   "RETAIL",
					// Uofm:       item.Data.UomId,
					Uofm:     "Each",
					Locncode: site.Data[0].Locncode,
					// Quantity:   int32(v.OrderQty),
					// Unitprce:   int32(v.UnitPrice),
					// Xtndprce:   int32(v.Subtotal),
					// GnL_Weight: int32(v.Weight),
					// Pricelvl:   "HIGHGJ",
					// Uofm:       item.Data.UomId,
				}
				fmt.Print(item)
				tempDetails = append(tempDetails, detail)
			}
			tempHeader := *&sales_service.CreateSalesOrderGPRequest{
				// Interid:            "",
				// Sopnumbe:           "",
				Docid:    "STDORD",
				Docdate:  header.RecognitionDate.AsTime().Format(layout),
				Custnmbr: req.Session.Customer.Code,
				Custname: req.Session.Customer.Name,
				Prstadcd: req.Session.Address.Code,
				Curncyid: "",
				// Subtotal:           int32(header.Total),
				Subtotal: 0,
				Trdisamt: 0,
				// Freight:            int32(header.DeliveryFee),
				Freight:            0,
				Miscamnt:           0,
				Taxamnt:            0,
				Docamnt:            0,
				GnlRequestShipDate: header.RequestDeliveryDate.AsTime().Format(layout),
				GnlRegion:          address.Data[0].AdministrativeDiv.GnlRegion,
				GnlWrtId:           wrtIdGP.Data.Code,
				GnlArchetypeId:     req.Session.Address.ArchetypeID,
				GnlOrderChannel:    "",
				GnlSoCodeApps:      "",
				GnlTotalweight:     int32(header.TotalWeight),
				Userid:             "",
				//Detailitems:        []*sales_service.CreateSalesOrderGPRequest_DetailItem{},
			}
			tempHeader.Detailitems = append(tempHeader.Detailitems, tempDetails...)

			soGP, err := s.opt.Client.SalesServiceGrpc.CreateSalesOrderGP(ctx, &tempHeader)
			if err != nil {
				fmt.Print(err)
			}
			soNumber = soGP.Data.DocNumber
		}

	}

	// start update customer profile talon of merchant
	// if paymentMethodID != 0 { // TODO: get paymentmethod from object payment and search to SalesServiceGrpc.GetPaymentGroupCombList
	// paymentMethod, err := s.opt.Client.BridgeServiceGrpc.GetSalesPaymentTermDetail(ctx, &bridge_service.GetSalesPaymentTermDetailRequest{
	// 	Id: int64(paymentMethodID),
	// })
	// if err != nil {

	// }
	// req.Session.Customer.PaymentMethod = &model.PaymentMethod{
	// 	Name: paymentMethod.Data.Description,
	// 	ID:   paymentMethod.Data.Id,
	// 	Code: paymentMethod.Data.Code,
	// }
	// paymentMethodName = req.Session.Customer.PaymentMethod.Name
	// paymentMethodName = ""
	// }
	//payment method

	// if e = talon.UpdateCustomerProfileTalon(r.Session.Merchant.ProfileCode, r.Session.Merchant.TagCustomer, r.Session.Merchant.FinanceArea.Name, r.Session.Merchant.BusinessType.Name, paymentMethodName, r.Session.Merchant.CreatedAt.Format("2006-01-02T15:04:05Z07:00"), r.Data.ReferrerData...); e != nil {
	// 	o.Rollback()
	// 	return nil, e
	// }
	// if _, e = o.Update(r.Session.Merchant, "ProfileCode"); e != nil {
	// 	o.Rollback()
	// 	return nil, e
	// }
	// if r.Data.IsInitCustomerTalonPoints {
	// 	e = talon.ChangeTalonPoints("add_points", "init points", r.Session.Merchant.ProfileCode, r.Session.Merchant.TotalPoint)
	// }
	// end update customer profile talon of merchant

	// // start update customer profile talon of merchant referrer and set up csv file for referral code talon.one
	// if r.Session.Merchant.ReferrerCode != "" {
	// 	paymentMethodName = ""
	// 	if r.Session.Merchant.ReferrerMerchant.PaymentMethod != nil {
	// 		paymentMethodName = r.Session.Merchant.ReferrerMerchant.PaymentMethod.Name
	// 	}
	// 	if e = talon.UpdateCustomerProfileTalon(r.Session.Merchant.ReferrerMerchant.ProfileCode, r.Session.Merchant.ReferrerMerchant.TagCustomer, r.Session.Merchant.ReferrerMerchant.FinanceArea.Name, r.Session.Merchant.ReferrerMerchant.BusinessType.Name, paymentMethodName, r.Session.Merchant.ReferrerMerchant.CreatedAt.Format("2006-01-02T15:04:05Z07:00")); e != nil {
	// 		o.Rollback()
	// 		return nil, e
	// 	}
	// 	if _, e = o.Update(r.Session.Merchant.ReferrerMerchant, "ProfileCode"); e != nil {
	// 		o.Rollback()
	// 		return nil, e
	// 	}
	// 	if r.Data.IsInitReferrerTalonPoints {
	// 		e = talon.ChangeTalonPoints("add_points", "init points", r.Session.Merchant.ReferrerMerchant.ProfileCode, r.Session.Merchant.ReferrerMerchant.TotalPoint)
	// 	}
	// 	if e = talon.SetUpCsvFileForReferral(r.Session.Merchant.ReferrerCode, r.Session.Merchant.ReferrerMerchant.ProfileCode); e != nil {
	// 		o.Rollback()
	// 		return nil, e
	// 	}
	// }
	// // end update customer profile talon of merchant referrer and set up csv file for referral code talon.one

	// // start hit talon's customer session api
	// voucherAmount := 0.00
	// if so.VouDiscAmount > 0 && r.Data.Voucher != nil && r.Data.Voucher.Type != 3 {
	// 	voucherAmount = so.VouDiscAmount
	// }
	// if r.Data.SessionResponse, e = talon.UpdateCustomerSessionTalon("open", "false", r.Data.IntegrationCode, r.Session.Merchant.ProfileCode, r.Data.Branch.Archetype.Name, r.Data.Branch.PriceSet.Name, r.Session.Merchant.ReferrerCode, r.Data.ItemList, isUsePoint, voucherAmount, r.Data.OrderType.Name); e != nil {
	// 	o.Rollback()
	// 	return nil, e
	// }
	// // end hit talon's customer session api

	// if r.Data.IsCreateCreditLimitLog {
	// 	if e = log.CreditLimitLogByMerchant(r.Session.Merchant, so.ID, "sales_order", r.Data.CreditLimitBefore, r.Data.CreditLimitAfter, "create sales order"); e != nil {
	// 		o.Rollback()
	// 		return nil, e/project-version3/erp-infra/env-erp/-/wikis/home
	// 	}
	// 	r.Session.Merchant.RemainingCreditLimitAmount = r.Data.CreditLimitAfter
	// 	if _, e = o.Update(r.Session.Merchant, "credit_limit_remaining"); e != nil {
	// 		o.Rollback()
	// 		return nil, e
	// 	}
	// }

	// if log.AuditLogByUser(&model.Staff{ID: 222}, so.ID, "sales_order", "create", ""); e != nil {
	// 	goto CANCEL
	// }

	// if r.Data.Voucher != nil {
	// 	r.Data.Voucher.RemOverallQuota = r.Data.Voucher.RemOverallQuota - 1
	// 	if _, e = o.Update(r.Data.Voucher, "rem_overall_quota"); e != nil {
	// 		goto CANCEL
	// 	}

	// 	vl := &model.VoucherLog{
	// 		Voucher:           r.Data.Voucher,
	// 		Merchant:          r.Data.Branch.Merchant,
	// 		Branch:            r.Data.Branch,
	// 		SalesOrder:        so,
	// 		TagCustomer:       r.Data.SameTagCustomer,
	// 		VoucherDiscAmount: r.Data.Voucher.VoucherAmount,
	// 		Timestamp:         time.Now(),
	// 		Status:            int8(1),
	// 	}

	// 	if _, e = o.Insert(vl); e != nil {
	// 		goto CANCEL
	// 	}

	// }

	// if r.Data.RedeemPoint != 0 {
	// 	var isExist bool

	// 	currentDate := time.Now()
	// 	mpl := &model.MerchantPointLog{
	// 		SalesOrder:       so,
	// 		PointValue:       r.Data.RedeemPoint,
	// 		RecentPoint:      r.Data.RecentPoint - r.Data.RedeemPoint,
	// 		CreatedDate:      currentDate,
	// 		Status:           2,
	// 		CurrentPointUsed: r.Data.CurrentPointUsed,
	// 		NextPointUsed:    r.Data.NextPointUsed,
	// 		ExpiredDate:      r.Data.MerchantPointExpiration.CurrentPeriodDate,
	// 		TransactionType:  5,
	// 	}

	// 	if so.Branch.Merchant.ID != 0 {
	// 		mpl.Merchant = so.Branch.Merchant
	// 	}

	// 	if _, e = o.Insert(mpl); e != nil {
	// 		goto CANCEL
	// 	}

	// 	so.PointRedeemID = mpl.ID

	// 	if _, e = o.Update(so, "point_redeem_id"); e != nil {
	// 		goto CANCEL
	// 	}

	// 	so.Branch.Merchant.TotalPoint = mpl.RecentPoint

	// 	if _, e = o.Update(so.Branch.Merchant, "TotalPoint"); e != nil {
	// 		goto CANCEL
	// 	}

	// 	// start create or update merchant point summary
	// 	if e = o.Raw("select exists(select id from merchant_point_summary mps where merchant_id = ? and summary_date = ?)", so.Branch.Merchant.ID, currentDate.Format("2006-01-02")).QueryRow(&isExist); e != nil || (e == nil && !isExist) {
	// 		if _, e = o.Raw("insert into merchant_point_summary (merchant_id, summary_date, earned_point, redeemed_point) values (?, ?, 0, ?)", so.Branch.Merchant.ID, currentDate.Format("2006-01-02"), r.Data.RedeemPoint).Exec(); e != nil {
	// 			goto CANCEL
	// 		}
	// 	} else {
	// 		if _, e = o.Raw("update merchant_point_summary set redeemed_point = redeemed_point + ? where merchant_id = ? and summary_date = ?", r.Data.RedeemPoint, so.Branch.Merchant.ID, currentDate.Format("2006-01-02")).Exec(); e != nil {
	// 			goto CANCEL
	// 		}
	// 	}
	// 	// end create or update merchant point summary

	// 	r.Data.MerchantPointExpiration.CurrentPeriodPoint = r.Data.CurrentPeriodPoint
	// 	r.Data.MerchantPointExpiration.NextPeriodPoint = r.Data.NextPeriodPoint
	// 	if _, e = o.Update(r.Data.MerchantPointExpiration, "CurrentPeriodPoint", "NextPeriodPoint"); e != nil {
	// 		goto CANCEL
	// 	}

	// }
	// CANCEL:
	// 	if !isCommited {
	// 		_, e = talon.UpdateCustomerSessionTalon("cancelled", "false", r.Data.IntegrationCode, r.Session.Merchant.ProfileCode, r.Data.Branch.Archetype.Name, r.Data.Branch.PriceSet.Name, r.Session.Merchant.ReferrerCode, r.Data.ItemList, isUsePoint, voucherAmount, r.Data.OrderType.Name)
	// 		o.Rollback()
	// 		return nil, e
	// 	}

	// 	orm := orm.NewOrm()
	// 	orm.Using("read_only")
	// 	messageNotif := &util.MessageNotification{}

	// 	if r.Session.Merchant.PaymentGroup.ID != 1 {
	// 		orm.Raw("SELECT message, title FROM notification WHERE code= 'NOT0001'").QueryRow(&messageNotif)
	// 	} else if r.Session.Merchant.PaymentGroup.ID == 1 {
	// 		orm.Raw("SELECT message, title FROM notification WHERE code= 'NOT0008'").QueryRow(&messageNotif)
	// 	}
	// Send Notification Transaction
	_, err = s.opt.Client.NotificationServiceGrpc.SendNotificationTransaction(ctx, &notification_service.SendNotificationTransactionRequest{
		CustomerId: req.Session.Customer.ID,
		RefId:      soNumber,
		Type:       "1",
		SendTo:     req.Session.Customer.UserCustomer.FirebaseToken,
		NotifCode:  "NOT0001",
		RefCode:    soNumber,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// 	if r.Data.IsCreateMerchantVa["bca"] == 1 {
	// 		xendit.BCAXenditFixedVA(r.Data.Branch.Merchant)
	// 	}

	// 	if r.Data.IsCreateMerchantVa["permata"] == 1 {
	// 		xendit.PermataXenditFixedVA(r.Data.Branch.Merchant)
	// 	}

	// 	// cahnge finish cart log
	// 	cart.Finish(r.Session.Merchant, r.Data.Branch, salesOrderItems)

	fmt.Println(archetype, paymentMethodName)
	res = dto.SalesOrderDetailResponse{
		ID:          soNumber,
		TotalCharge: strconv.Itoa(int(totCharge)),
		Code:        soNumber,
	}
	return
}

func (s *SalesOrderService) UpdateCOD(ctx context.Context, req *dto.UpdateCodRequest) (res dto.SalesOrderDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.UpdateCOD")
	defer span.End()

	soID, _ := strconv.Atoi(req.Data.SalesOrderID)
	addressID, _ := strconv.Atoi(req.Session.Address.ID)
	so, err := s.opt.Client.SalesServiceGrpc.GetSalesOrderDetail(ctx, &sales_service.GetSalesOrderDetailRequest{
		Id: int64(soID),
	})
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressDetail(ctx, &bridge_service.GetAddressDetailRequest{
		Id: int64(addressID),
	})
	// // check is eligible to update to COD and Exclude validation for order type 'self pickup'
	// if req.Session.Customer.PaymentGroup.ID != 3 && so.Data.OrderTypeId != 6 {
	// o.Failure("sales_order_id", util.ErrorCannotUpdatePayment())
	// return o
	// }

	// //check is SO already choose advance payment
	// if so.Data.HasExtInvoice == 1 {
	// 	o.Failure("sales_order_id", util.ErrorCannotUpdatePayment())
	// 	return o
	// }

	if so.Data.Status == 0 || so.Data.Status == 2 || so.Data.Status == 3 || so.Data.Status == 4 {
		// o.Failure("sales_order_id", util.ErrorCannotUpdatePayment())
		// return o
	}

	so.Data.PaymentGroupSlsId = 2
	so.Data.TermInvoiceSlsId = 2
	so.Data.TermPaymentSlsId = 10

	_, err = s.opt.Client.SalesServiceGrpc.UpdateSalesOrderHeader(ctx, &sales_service.UpdateSalesOrderHeaderRequest{
		Data: &sales_service.SalesOrder{
			Id:                int64(soID),
			PaymentGroupSlsId: so.Data.PaymentGroupSlsId,
			TermInvoiceSlsId:  so.Data.TermInvoiceSlsId,
			TermPaymentSlsId:  so.Data.TermPaymentSlsId,
		},
	})
	// log.AuditLogByUser(&model.Staff{ID: 222}, r.Data.SalesOrder.ID, "sales_order", "update", "Update Payment Data")

	fmt.Print(address)
	return
}

func (s *SalesOrderService) GetSalesOrderFeedback(ctx context.Context, req *dto.GetFeedback) (res []dto.SalesOrderFeedback, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.UpdateCOD")
	defer span.End()
	customerID, _ := strconv.Atoi(req.Session.Customer.ID)

	sof, err := s.opt.Client.SalesServiceGrpc.GetSalesOrderFeedbackList(ctx, &sales_service.GetSalesOrderFeedbackListRequest{
		FeedbackType: req.Data.FeedbackType,
		CustomerId:   int64(customerID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range sof.Data {
		res = append(res, dto.SalesOrderFeedback{
			SalesOrderCode: v.SalesOrderCode,
			DeliveryDate:   v.DeliveryDate,
			RatingScore:    strconv.Itoa(int(v.RatingScore)),
			Tags:           v.Tags,
			Description:    v.Description,
			TotalCharge:    strconv.FormatFloat(v.TotalCharge, 'f', 1, 64),
			SalesOrder:     strconv.Itoa(int(v.SalesOrderId)),
		})
	}
	return
}

func (s *SalesOrderService) CreateSalesOrderFeedback(ctx context.Context, req *dto.CreateSalesFeedback) (res dto.SalesOrderFeedback, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.UpdateCOD")
	defer span.End()

	customerID, _ := strconv.Atoi(req.Session.Customer.ID)
	soID, _ := strconv.Atoi(req.Data.SalesOrderId)

	so, err := s.opt.Client.BridgeServiceGrpc.GetSalesOrderDetail(ctx, &bridge_service.GetSalesOrderDetailRequest{
		Id: int64(soID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	_, err = s.opt.Client.BridgeServiceGrpc.GetAddressDetail(ctx, &bridge_service.GetAddressDetailRequest{
		Id: utils.ToInt64(so.Data.AddressId),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	customerSO, err := s.opt.Client.BridgeServiceGrpc.GetCustomerDetail(ctx, &bridge_service.GetCustomerDetailRequest{
		Id: utils.ToInt64(so.Data.CustomerId),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// _, err = s.opt.Client.SalesServiceGrpc.GetSalesOrderFeedbackList(ctx, &sales_service.GetSalesOrderFeedbackListRequest{
	// 	SalesOrderId: so.Data.Id,
	// })
	// if err == nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	if req.Data.RatingScore > 10 {
		// o.Failure("rating_score", "rating tidak boleh lebih dari 10")
	}

	if req.Data.RatingScore < 0 {
		// o.Failure("rating_score", "rating tidak boleh kurang dari 0")
	}

	//array tags from mobile that still hardcoded
	tags := [...]string{"Pengiriman terlambat",
		"Kualitas produk buruk",
		"Timbangan kurang",
		"Tidak sesuai spesifikasi",
		"Produk kurang/tertinggal",
		"Tidak sesuai SKU",
		"Produk kosong",
		"Kurir bermasalah",
		"Return bermasalah",
		"Pengiriman cepat",
		"Kualitas produk bagus",
		"Harga murah",
	}

	//compare between data from mobile with array tags. if it doesnt exist doesnt save it
	for i := 0; i < len(req.Data.Tags); i++ {
		if !util.ItemExists(tags, req.Data.Tags[i]) {
			continue
		}
		req.Data.ExistTags = req.Data.ExistTags + req.Data.Tags[i] + ";"
	}

	if customerID != int(customerSO.Data.Id) {
		// o.Failure("merchant.invalid", "sales order tidak cocok dengan merchant")
	}
	layout := "2006-01-02"

	data := &sales_service.SalesOrderFeedback{
		SalesOrderCode: so.Data.Code,
		DeliveryDate:   so.Data.CreatedAt.AsTime().Format(layout),
		RatingScore:    int32(req.Data.RatingScore),
		Tags:           req.Data.ExistTags,
		Description:    req.Data.Description,
		ToBeContacted:  int32(req.Data.ToBeContacted),
		CreatedAt:      timestamppb.New(time.Now()),
		TotalCharge:    so.Data.Total,
		SalesOrderId:   int64(soID),
		CustomerId:     int64(customerID),
	}
	_, err = s.opt.Client.SalesServiceGrpc.CreateSalesOrderFeedback(ctx, &sales_service.CreateSalesOrderFeedbackRequest{
		Data: data,
	})
	res = dto.SalesOrderFeedback{
		SalesOrderCode: data.SalesOrderCode,
	}
	return
}
