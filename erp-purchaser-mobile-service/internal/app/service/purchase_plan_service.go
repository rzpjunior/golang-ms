package service

import (
	"context"
	"fmt"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/notification_service"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	. "github.com/ahmetb/go-linq/v3"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IPurchasePlanService interface {
	Get(ctx context.Context, req *dto.PurchasePlanListRequest) (res []*dto.PurchasePlanResponse, total int64, err error)
	GetByID(ctx context.Context, id string) (res *dto.PurchasePlanResponse, err error)
	GetSummary(ctx context.Context, req *dto.PurchasePlanListRequest) (res *dto.SummaryPurchasePlanResponse, err error)
	Assign(ctx context.Context, req dto.PurchasePlanRequestAssign, id string) (res *dto.PurchasePlanResponse, err error)
	CancelAssign(ctx context.Context, id string) (res *dto.PurchasePlanResponse, err error)
}

type PurchasePlanService struct {
	opt opt.Options
}

func NewPurchasePlanService() IPurchasePlanService {
	return &PurchasePlanService{
		opt: global.Setup.Common,
	}
}

func (s *PurchasePlanService) Get(ctx context.Context, req *dto.PurchasePlanListRequest) (res []*dto.PurchasePlanResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanService.Get")
	defer span.End()

	var statusGP int32

	switch statusx.ConvertStatusValue(int8(req.Status)) {
	case statusx.Cancelled:
		statusGP = 3
	case statusx.Active:
		statusGP = 4
	// case statusx.Draft:
	// 	statusGP = 1
	case statusx.Finished:
		statusGP = 12
	case statusx.Closed:
		statusGP = 11
	}

	if req.Status == 0 {
		statusGP = 4
	}

	var purchasePlans *bridgeService.GetPurchasePlanGPResponse
	getPurchasePlanListRequest := &bridgeService.GetPurchasePlanGPListRequest{
		Limit:  req.Limit,
		Offset: req.Offset,
		Status: statusGP,
		// OrderBy:             req.OrderBy,
		Locncode:          req.SiteID,
		FieldPurchaser:    req.FieldPurchaser,
		PrpPurchaseplanNo: req.Search,
	}
	if req.Search != "" {
		getPurchasePlanListRequest.PrpPurchaseplanNo = req.Search
		getPurchasePlanListRequest.PrpVendorOrgDesc = req.Search
	}

	// if timex.IsValid(req.RecognitionDateFrom) {
	// 	getPurchasePlanListRequest.PrpPurchaseplanDateFrom = req.RecognitionDateFrom.Format(timex.InFormatDate)
	// }

	// if timex.IsValid(req.RecognitionDateTo) {
	// 	getPurchasePlanListRequest.PrpPurchaseplanDateTo = req.RecognitionDateTo.Format(timex.InFormatDate)
	// }

	if timex.IsValid(req.PurchasePlanDateFrom) {
		getPurchasePlanListRequest.PrpPurchaseplanDateFrom = req.PurchasePlanDateFrom.Format(timex.InFormatDate)
	} else {
		// Get the current date
		req.PurchasePlanDateFrom = time.Now().AddDate(0, 0, -30)
		getPurchasePlanListRequest.PrpPurchaseplanDateFrom = req.PurchasePlanDateFrom.Format(timex.InFormatDate)
	}

	if timex.IsValid(req.PurchasePlanDateTo) {
		getPurchasePlanListRequest.PrpPurchaseplanDateTo = req.PurchasePlanDateTo.Format(timex.InFormatDate)
	} else {
		//set to end of year
		currentTime := time.Now().AddDate(0, 0, 1)
		// currentYear := currentTime.Year()

		// endOfYear := time.Date(currentYear, time.December, 31, 23, 59, 59, 0, currentTime.Location())
		getPurchasePlanListRequest.PrpPurchaseplanDateTo = currentTime.Format(timex.InFormatDate)
	}

	purchasePlans, err = s.opt.Client.BridgeServiceGrpc.GetPurchasePlanGPList(ctx, getPurchasePlanListRequest)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase_plan")
		return
	}

	for _, purchasePlan := range purchasePlans.Data {

		sumPPQtyByUomMaps := map[string]float64{}

		// get tonnage by uom
		for _, ppi := range purchasePlan.Detail {
			// get item
			var item *bridgeService.GetItemGPResponse
			item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
				Id: ppi.Itemnmbr,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "item")
				return
			}

			if item.Data[0].Uomschdl != "" {
				// get uom
				var uom *bridgeService.GetUomGPResponse
				uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
					Id: item.Data[0].Uomschdl,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorRpcNotFound("bridge", "uom")
					return
				}

				sumPPQtyByUomMaps[uom.Data[0].Umschdsc] += ppi.Quantity
			}
		}

		var tonnagePurchasePlanResponse []*dto.TonnagePurchasePlanResponse
		for uomName, totalPPQty := range sumPPQtyByUomMaps {
			tonnagePurchasePlanResponse = append(tonnagePurchasePlanResponse, &dto.TonnagePurchasePlanResponse{
				UomName:     uomName,
				TotalWeight: totalPPQty,
			})
		}

		// // get vendor classification
		// var vendorClassification *bridgeService.GetVendorClassificationGPResponse
		// vendorClassification, err = s.opt.Client.BridgeServiceGrpc.GetVendorClassificationGPDetail(ctx, &bridgeService.GetVendorClassificationGPDetailRequest{
		// 	Id: vendorOrganization.Data[0].c,
		// })
		// if err != nil {
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorRpcNotFound("bridge", "vendor_classification")
		// 	return
		// }

		// get site
		var site *bridgeService.GetSiteGPResponse
		site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
			Id: purchasePlan.Locncode,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}

		var recognitionDate time.Time
		recognitionDate, err = time.Parse("2006-01-02", purchasePlan.Docdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("recognition_date")
			return
		}

		var etaDate time.Time
		etaDate, err = time.Parse("2006-01-02", purchasePlan.PrpPurchaseplanDate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_date")
			return
		}

		var etaTime time.Time
		etaTime, err = time.Parse("15:04:05", purchasePlan.PrpEstimatedarrivalTim)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_time")
			return
		}

		switch purchasePlan.PrStatus {
		case 3:
			statusGP = 3
		case 4:
			statusGP = 1
		case 1:
			statusGP = 5
		case 12:
			statusGP = 2
		case 11:
			statusGP = 33
		}

		purchasePlanResponse := &dto.PurchasePlanResponse{
			ID:   purchasePlan.PrpPurchaseplanNo,
			Code: purchasePlan.PrpPurchaseplanNo,
			Site: &dto.SiteResponse{
				ID:          site.Data[0].Locncode,
				Code:        site.Data[0].Locncode,
				Description: site.Data[0].Locndscr,
			},
			RecognitionDate:     recognitionDate,
			EtaDate:             etaDate,
			EtaTime:             etaTime.Format("15:04"),
			Note:                purchasePlan.PrpNote,
			Status:              statusGP,
			TotalSku:            len(purchasePlan.Detail),
			TonnagePurchasePlan: tonnagePurchasePlanResponse,
		}

		if purchasePlan.PrpVendorOrgId != "" {
			// get vendor organization
			var vendorOrganization *bridgeService.GetVendorOrganizationGPResponse
			vendorOrganization, err = s.opt.Client.BridgeServiceGrpc.GetVendorOrganizationGPDetail(ctx, &bridgeService.GetVendorOrganizationGPDetailRequest{
				Id: purchasePlan.PrpVendorOrgId,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "vendor_organization")
				return
			}

			purchasePlanResponse.VendorOrganization = &dto.VendorOrganizationResponse{
				ID:          vendorOrganization.Data[0].PrpVendorOrgId,
				Code:        vendorOrganization.Data[0].PrpVendorOrgId,
				Description: vendorOrganization.Data[0].PrpVendorOrgDesc,
			}
		}

		res = append(res, purchasePlanResponse)

		var assigneTo *accountService.GetUserDetailResponse
		assigneTo, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
			EmployeeCode: purchasePlan.PrpPurchaseplanUser,
		})
		if assigneTo != nil {
			purchasePlanResponse.AssignedTo = &dto.UserResponse{
				ID:           assigneTo.Data.Id,
				Name:         assigneTo.Data.Name,
				EmployeeCode: assigneTo.Data.EmployeeCode,
			}
		}

		// var assigneBy *accountService.GetUserDetailResponse
		// assigneBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		// 	Id: purchasePlan.AssignedBy,
		// })
		// if assigneBy != nil {
		// 	purchasePlanResponse.AssignedBy = &dto.UserResponse{
		// 		ID:   assigneBy.Data.Id,
		// 		Name: assigneBy.Data.Name,
		// 	}
		// }

		// var createdBy *accountService.GetUserDetailResponse
		// createdBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		// 	Id: purchasePlan.CreatedBy,
		// })
		// if createdBy != nil {
		// 	purchasePlanResponse.CreatedBy = &dto.UserResponse{
		// 		ID:   createdBy.Data.Id,
		// 		Name: createdBy.Data.Name,
		// 	}
		// }
	}

	total = int64(len(res))

	return
}

func (s *PurchasePlanService) GetByID(ctx context.Context, id string) (res *dto.PurchasePlanResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanService.GetByID")
	defer span.End()

	var purchasePlan *bridgeService.GetPurchasePlanGPResponse
	var totalPurchasePlanQty, totalPurchaseQty float64
	purchasePlan, err = s.opt.Client.BridgeServiceGrpc.GetPurchasePlanGPDetail(ctx, &bridgeService.GetPurchasePlanGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchasePlan")
		return
	}
	purchaseOrder, err := s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPList(ctx, &bridgeService.GetPurchaseOrderGPListRequest{
		Limit:          1000,
		Offset:         0,
		PurchasePlanId: id,
		Status:         4,
	})
	// fmt.Println(purchaseOrder)
	var purchasePlanItemsResponse []*dto.PurchasePlanItemResponse
	var TotalPrice float64

	for i, ppi := range purchasePlan.Data[0].Detail {
		// get item
		var item *bridgeService.GetItemGPResponse
		item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: ppi.Itemnmbr,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item")
			return
		}

		itemResponse := &dto.ItemResponse{
			ID:    item.Data[0].Itemnmbr,
			Code:  item.Data[0].Itemnmbr,
			Class: &dto.ClassResponse{
				// ID: item.Data.ClassId,
			},
			Description:          item.Data[0].Itemdesc,
			UnitWeightConversion: item.Data[0].GnlWeighttolerance,
			OrderMinQty:          item.Data[0].Minorqty,
			OrderMaxQty:          item.Data[0].Maxordqty,
			ItemType:             item.Data[0].ItemTypeDesc,
			Capitalize:           item.Data[0].Itemdesc,
			UnitPrice:            item.Data[0].Currcost,
			// MaxDayDeliveryDate:   int8(item.Data.MaxDayDeliveryDate),
			// Taxable:              item.Data.Taxable,
			// Note:                 item.Data.Note,
			// Status:               int8(item.Data.Status),
		}
		totalPurchasePlanQty += ppi.Quantity
		if item.Data[0].Uomschdl != "" {
			var uom *bridgeService.GetUomGPResponse
			uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
				Id: item.Data[0].Uomschdl,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "uom")
				return
			}
			decimalEnabled := 0
			if uom.Data[0].Umdpqtys == 3 {
				decimalEnabled = 1
			}

			itemResponse.Uom = &dto.UomResponse{
				ID:             uom.Data[0].Uomschdl,
				Code:           uom.Data[0].Uomschdl,
				Name:           uom.Data[0].Umschdsc,
				DecimalEnabled: decimalEnabled,
			}
		}

		purchasePlanItemsResponse = append(purchasePlanItemsResponse, &dto.PurchasePlanItemResponse{
			ID:             ppi.Itemnmbr,
			PurchasePlanID: purchasePlan.Data[0].PrpPurchaseplanNo,
			Item:           itemResponse,
			OrderQty:       ppi.Quantity,
			PurchaseQty:    0,
			UnitPrice:      ppi.Unitcost,
			Subtotal:       (ppi.Unitcost * ppi.Quantity),
			// Weight:      ppi.Weight,
		},
		)
		TotalPrice += purchasePlanItemsResponse[i].Subtotal

		var tempPurchOrderList []dto.PurchaseOrderItem
		tempQty := 0.0
		for _, v := range purchaseOrder.Data {
			// fmt.Print(tempQty)
			exist := From(v.PoDetail).WhereT(
				func(f *bridgeService.PurchaseOrderGP_PoDetail) bool {
					return (f.Itemnmbr == ppi.Itemnmbr)
				},
			).Count()
			if exist > 0 {
				qtyPOItem := From(v.PoDetail).WhereT(
					func(f *bridgeService.PurchaseOrderGP_PoDetail) bool {
						return (f.Itemnmbr == ppi.Itemnmbr)
					},
				).SelectT(func(f *bridgeService.PurchaseOrderGP_PoDetail) float64 {
					return f.Qtyorder
				}).SumFloats()
				fmt.Println("PO=" + v.Ponumber)
				fmt.Println("QTY=" + utils.ToString(qtyPOItem))
				tempQty += qtyPOItem

				tempPurchOrderList = append(tempPurchOrderList, dto.PurchaseOrderItem{
					PurchaseQty:   qtyPOItem,
					PurchaseOrder: v.Ponumber,
				})
			}
			// fmt.Print(exist)
		}
		purchasePlanItemsResponse[i].PurchaseOrderID = tempPurchOrderList
		purchasePlanItemsResponse[i].PurchaseQty = tempQty
		totalPurchaseQty += tempQty
	}

	// get vendor classification
	// var vendorClassification *bridgeService.GetVendorClassificationDetailResponse
	// vendorClassification, err = s.opt.Client.BridgeServiceGrpc.GetVendorClassificationDetail(ctx, &bridgeService.GetVendorClassificationDetailRequest{
	// 	Id: vendorOrganization.Data.VendorClassificationId,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "vendor_classification")
	// 	return
	// }

	// get site
	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: purchasePlan.Data[0].SiteId[0].Locncode,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	// get admDivision
	// var admDivision *bridgeService.GetAdmDivisionDetailResponse
	// admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridgeService.GetAdmDivisionDetailRequest{
	// 	SubDistrictId: vendorOrganization.Data.SubDistrictId,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "adm_division")
	// 	return
	// }

	var recognitionDate time.Time
	recognitionDate, err = time.Parse("2006-01-02", purchasePlan.Data[0].Docdate)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("recognition_date")
		return
	}

	var etaDate time.Time
	etaDate, err = time.Parse("2006-01-02", purchasePlan.Data[0].PrpPurchaseplanDate)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("eta_date")
		return
	}

	var etaTime time.Time
	etaTime, err = time.Parse("15:04:05", purchasePlan.Data[0].PrpEstimatedarrivalTim)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("eta_time")
		return
	}

	var statusGP int32
	switch purchasePlan.Data[0].PrStatus {
	case 3:
		statusGP = 3
	case 4:
		statusGP = 1
	case 1:
		statusGP = 5
	case 11:
		statusGP = 33
	case 12:
		statusGP = 2
	}

	res = &dto.PurchasePlanResponse{
		ID:   purchasePlan.Data[0].PrpPurchaseplanNo,
		Code: purchasePlan.Data[0].PrpPurchaseplanNo,
		Site: &dto.SiteResponse{
			ID:          site.Data[0].Locncode,
			Code:        site.Data[0].Locncode,
			Description: site.Data[0].Locndscr,
		},
		RecognitionDate: recognitionDate,
		EtaDate:         etaDate,
		EtaTime:         etaTime.Format("15:04"),
		TotalPrice:      TotalPrice,
		// TotalWeight:          purchasePlan.TotalWeight,
		TotalPurchasePlanQty: totalPurchasePlanQty,
		TotalPurchaseQty:     totalPurchaseQty,
		Status:               statusGP,
		Note:                 purchasePlan.Data[0].PrpNote,
		// AssignedAt:           purchasePlan.AssignedAt.AsTime(),
		// CreatedAt:            purchasePlan.CreatedAt.AsTime(),
		TotalSku:          len(purchasePlan.Data[0].Detail),
		PurchasePlanItems: purchasePlanItemsResponse,
	}

	// get vendor organization
	if purchasePlan.Data[0].VendorOrganization[0].PrpVendorOrgId != "" {
		var vendorOrganization *bridgeService.GetVendorOrganizationGPResponse
		vendorOrganization, err = s.opt.Client.BridgeServiceGrpc.GetVendorOrganizationGPDetail(ctx, &bridgeService.GetVendorOrganizationGPDetailRequest{
			Id: purchasePlan.Data[0].VendorOrganization[0].PrpVendorOrgId,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor_organization")
			return
		}

		res.VendorOrganization = &dto.VendorOrganizationResponse{
			ID:          vendorOrganization.Data[0].PrpVendorOrgId,
			Code:        vendorOrganization.Data[0].PrpVendorOrgId,
			Description: vendorOrganization.Data[0].PrpVendorOrgDesc,
		}
	}

	var assigneTo *accountService.GetUserDetailResponse
	assigneTo, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		EmployeeCode: purchasePlan.Data[0].AssignedTo.PrpPurchaseplanUser,
	})
	if assigneTo != nil {
		res.AssignedTo = &dto.UserResponse{
			ID:           assigneTo.Data.Id,
			Name:         assigneTo.Data.Name,
			EmployeeCode: assigneTo.Data.EmployeeCode,
		}
	}

	// var assigneBy *accountService.GetUserDetailResponse
	// assigneBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	Id: purchasePlan.Data.AssignedBy,
	// })
	// if assigneBy != nil {
	// 	res.AssignedBy = &dto.UserResponse{
	// 		ID:   assigneBy.Data.Id,
	// 		Name: assigneBy.Data.Name,
	// 	}
	// }

	// var createdBy *accountService.GetUserDetailResponse
	// createdBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	Id: purchasePlan.Data.CreatedBy,
	// })
	// if createdBy != nil {
	// 	res.CreatedBy = &dto.UserResponse{
	// 		ID:   createdBy.Data.Id,
	// 		Name: createdBy.Data.Name,
	// 	}
	// }

	return
}

func (s *PurchasePlanService) GetSummary(ctx context.Context, req *dto.PurchasePlanListRequest) (res *dto.SummaryPurchasePlanResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanService.GetSummary")
	defer span.End()
	datefrom := time.Now().AddDate(0, 0, -30).Format(timex.InFormatDate)
	dateto := time.Now().AddDate(0, 0, 1).Format(timex.InFormatDate)
	var purchasePlans *bridgeService.GetPurchasePlanGPResponse
	if req.FieldPurchaser != "" {
		// Get the current date

		getPurchasePlanListRequest := &bridgeService.GetPurchasePlanGPListRequest{
			Limit:                   1,
			Status:                  4,
			FieldPurchaser:          req.FieldPurchaser,
			PrpPurchaseplanDateFrom: datefrom,
			PrpPurchaseplanDateTo:   dateto,
			Locncode:                req.SiteID,
		}

		purchasePlans, err = s.opt.Client.BridgeServiceGrpc.GetPurchasePlanGPList(ctx, getPurchasePlanListRequest)
		fmt.Println(err)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "purchase_plan")
			return
		}

		res = &dto.SummaryPurchasePlanResponse{
			TotalActive: int64(purchasePlans.TotalRecords),
			// TotalAssigned: int64(purchasePlans.TotalRecords),
		}
	} else {
		getPurchasePlanListRequest := &bridgeService.GetPurchasePlanGPListRequest{
			Limit:                   1,
			PrpPurchaseplanDateFrom: datefrom,
			PrpPurchaseplanDateTo:   dateto,
			Status:                  4,
			Locncode:                req.SiteID,
		}

		purchasePlans, err = s.opt.Client.BridgeServiceGrpc.GetPurchasePlanGPList(ctx, getPurchasePlanListRequest)
		fmt.Println(err)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "purchase_plan")
			return
		}
		var totalAssigned int64
		for _, v := range purchasePlans.Data {
			if v.PrpPurchaseplanUser != "" {
				totalAssigned += 1
			}
		}

		res = &dto.SummaryPurchasePlanResponse{
			TotalActive: int64(purchasePlans.TotalRecords),
			// TotalAssigned: totalAssigned,
		}
	}
	//tes
	return
}

func (s *PurchasePlanService) Assign(ctx context.Context, req dto.PurchasePlanRequestAssign, id string) (res *dto.PurchasePlanResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanService.Assign")
	defer span.End()

	var purchasePlan *bridgeService.GetPurchasePlanGPResponse
	purchasePlan, err = s.opt.Client.BridgeServiceGrpc.GetPurchasePlanGPDetail(ctx, &bridgeService.GetPurchasePlanGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchasePlan")
		return
	}
	var oldPurchaserUser *accountService.GetUserDetailResponse
	oldPurchaserUser, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		EmployeeCode: purchasePlan.Data[0].AssignedTo.PrpPurchaseplanUser,
	})

	var newPurchaserUser *accountService.GetUserDetailResponse
	newPurchaserUser, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		EmployeeCode: req.FieldPurchaserID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user")
		return
	}

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user")
		return
	}
	fmt.Println("")
	fmt.Println(oldPurchaserUser.Data.EmployeeCode)
	fmt.Println(newPurchaserUser.Data.EmployeeCode)

	var _ *bridgeService.AssignPurchasePlanGPResponse
	_, err = s.opt.Client.BridgeServiceGrpc.AssignPurchasePlanGP(ctx, &bridgeService.AssignPurchasePlanGPRequest{
		PrpPurchaseplanNo:   id,
		PrpPurchaseplanUser: req.FieldPurchaserID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "assign_purchase_plan")
		return
	}

	// Notification for old field purchaser NOT0022
	_, err = s.opt.Client.NotificationServiceGrpc.SendNotificationPurchaser(ctx, &notification_service.SendNotificationPurchaserRequest{
		SendTo:    oldPurchaserUser.Data.PurchaserappNotifToken,
		Type:      "6",
		NotifCode: "NOT0022",
		RefId:     id,
		StaffId:   req.Session.EmployeeCode,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Notification for new field purchaser NOT0021
	_, err = s.opt.Client.NotificationServiceGrpc.SendNotificationPurchaser(ctx, &notification_service.SendNotificationPurchaserRequest{
		SendTo:    newPurchaserUser.Data.PurchaserappNotifToken,
		Type:      "6",
		NotifCode: "NOT0021",
		RefId:     id,
		StaffId:   req.Session.EmployeeCode,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// audit log
	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId: userID,
			// ReferenceId:,
			Type:      "assign_sales_order",
			Function:  "PurchasePlanService.Assign",
			CreatedAt: timestamppb.New(time.Now()),
		},
	})

	return
}

func (s *PurchasePlanService) CancelAssign(ctx context.Context, id string) (res *dto.PurchasePlanResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchasePlanService.CancelAssign")
	defer span.End()

	var _ *bridgeService.CancelAssignPurchasePlanResponse
	_, err = s.opt.Client.BridgeServiceGrpc.CancelAssignPurchasePlan(ctx, &bridgeService.CancelAssignPurchasePlanRequest{})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "cancel_assign_purchase_plan")
		return
	}

	return
}
