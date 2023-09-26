package service

import (
	"context"
	"time"

	. "github.com/ahmetb/go-linq/v3"
	"google.golang.org/protobuf/types/known/timestamppb"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"

	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/repository"
)

type IPurchaseOrderService interface {
	Get(ctx context.Context, req *dto.PurchaseOrderListRequest) (res []*dto.PurchaseOrderResponse, total int64, err error)
	GetByID(ctx context.Context, id string) (res *dto.PurchaseOrderResponse, err error)
	Create(ctx context.Context, req dto.PurchaseOrderRequestCreate) (res *dto.PurchaseOrderResponse, err error)
	Update(ctx context.Context, req dto.PurchaseOrderRequestUpdate, id string) (res *dto.PurchaseOrderResponse, err error)
	Assign(ctx context.Context, req dto.PurchaseOrderRequestAssign, id string) (res *dto.PurchaseOrderResponse, err error)
	Signature(ctx context.Context, req dto.PurchaseOrderRequestSignature, id string) (err error)
	Print(ctx context.Context, id string) (res *dto.PurchaseOrderResponse, err error)
	Cancel(ctx context.Context, req dto.PurchaseOrderRequestCancel, id string) (res *dto.PurchaseOrderResponse, err error)
}

type PurchaseOrderService struct {
	opt                              opt.Options
	RepositoryPurchaseOrder          repository.IPurchaseOrderRepository
	RepositoryPurchaseOrderItem      repository.IPurchaseOrderItemRepository
	RepositoryPurchaseOrderSignature repository.IPurchaseOrderSignatureRepository
	RepositoryConsolidatedShipment   repository.IConsolidatedShipmentRepository
	RepositoryPurchaseOrderImage     repository.IPurchaseOrderImageRepository
}

func NewPurchaseOrderService() IPurchaseOrderService {
	return &PurchaseOrderService{
		opt:                              global.Setup.Common,
		RepositoryPurchaseOrder:          repository.NewPurchaseOrderRepository(),
		RepositoryPurchaseOrderItem:      repository.NewPurchaseOrderItemRepository(),
		RepositoryPurchaseOrderSignature: repository.NewPurchaseOrderSignatureRepository(),
		RepositoryConsolidatedShipment:   repository.NewConsolidatedShipmentRepository(),
		RepositoryPurchaseOrderImage:     repository.NewPurchaseOrderImageRepository(),
	}
}

func (s *PurchaseOrderService) Get(ctx context.Context, req *dto.PurchaseOrderListRequest) (res []*dto.PurchaseOrderResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Get")
	defer span.End()

	var statusGP int32

	switch statusx.ConvertStatusValue(int8(req.Status)) {
	// status convert
	case statusx.Active:
		statusGP = 4
	case statusx.Cancelled: // when cancel PO, GP return 8 as cancelled
		statusGP = 6
	case statusx.Finished:
		statusGP = 12
	case statusx.Closed:
		statusGP = 8 // change 11 ito 8
	}

	getPurchaseOrdersRequest := &bridgeService.GetPurchaseOrderGPListRequest{
		Limit:               req.Limit,
		Offset:              req.Offset,
		Status:              statusGP,
		Code:                req.Search,
		PurchasePlanId:      req.PurchasePlanID,
		IsNotConsolidated:   req.IsNotConsolidated,
		IsPurchasePlan:      "yes",
		Orderby:             "desc",
		PrpPurchaseplanUser: req.EmployeeCode,
		Locncode:            req.Site,
		PrpCsNo:             req.PrpCsNo,
		// OrderBy: req.OrderBy,
		// Locncode: req.SiteID,
		// RecognitionDateFrom: req.RecognitionDateFrom,
		// RecognitionDateTo:   req.RecognitionDateTo,
	}

	if timex.IsValid(req.RecognitionDateFrom) {
		getPurchaseOrdersRequest.ReqdateFrom = req.RecognitionDateFrom.Format(timex.InFormatDate)
	}

	if timex.IsValid(req.RecognitionDateTo) {
		getPurchaseOrdersRequest.ReqdateTo = req.RecognitionDateTo.Format(timex.InFormatDate)
	}

	var purchaseOrders *bridgeService.GetPurchaseOrderGPResponse
	purchaseOrders, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPList(ctx, getPurchaseOrdersRequest)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase_order")
		return
	}

	for _, purchaseOrder := range purchaseOrders.Data {
		var purchaseOrderItemsResponse []*dto.PurchaseOrderItemResponse
		var poiResponse *dto.PurchaseOrderItemResponse
		var vendor *bridgeService.GetVendorGPResponse
		vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
			Id: purchaseOrder.Vendorid,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			return
		}

		purchaseOrderItems := purchaseOrder.PoDetail

		sumPPQtyByUomMaps := map[string]float64{}

		// var uomLabel string
		// get item PO
		for _, poi := range purchaseOrderItems {
			poiResponse = &dto.PurchaseOrderItemResponse{
				ID:              poi.Ponumber,
				PurchaseOrderID: poi.Ponumber,
				Item: &dto.ItemResponse{
					ID:          poi.Itemnmbr,
					Code:        poi.Itemnmbr,
					Description: poi.Itemdesc,
					Uom: &dto.UomResponse{
						ID:   poi.Uofm,
						Code: poi.Uofm,
						Name: poi.Uofm,
					},
				},
				OrderQty:    poi.Qtyorder,
				UnitPrice:   poi.Unitcost,
				Subtotal:    poi.Qtyorder * poi.Unitcost,
				PurchaseQty: poi.Qtyorder,
			}
			sumPPQtyByUomMaps[poi.Uofm] += poi.Qtyorder
			// append all item
			purchaseOrderItemsResponse = append(purchaseOrderItemsResponse, poiResponse)
		}

		var tonnagePurchaseOrderResponse []*dto.TonnagePurchaseOrderResponse
		for uomName, totalPPQty := range sumPPQtyByUomMaps {
			tonnagePurchaseOrderResponse = append(tonnagePurchaseOrderResponse, &dto.TonnagePurchaseOrderResponse{
				UomName:     uomName,
				TotalWeight: totalPPQty,
			})
		}

		var docDate time.Time
		docDate, err = time.Parse("2006-01-02T15:04:05", purchaseOrder.Docdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("doc_date")
			return
		}

		var recognitionDate time.Time
		recognitionDate, err = time.Parse("2006-01-02T15:04:05", purchaseOrder.Reqdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("recognition_date")
			return
		}

		var etaDate time.Time
		etaDate, err = time.Parse("2006-01-02T15:04:05", purchaseOrder.PrpEstimatedarrivalDat)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_date")
			return
		}

		var etaTime time.Time
		etaTime, err = time.Parse("2006-01-02T15:04:05", purchaseOrder.PrpEstimatedarrivalTim)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_time")
			return
		}

		switch purchaseOrder.Postatus {
		case 4:
			statusGP = 1
		case 1:
			statusGP = 5
		case 8:
			statusGP = 33 //change 3 to 33
		case 12:
			statusGP = 2
		case 11:
			statusGP = 33
		}

		var accountName *account_service.GetUserDetailResponse

		accountName, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &account_service.GetUserDetailRequest{
			EmployeeCode: purchaseOrder.PrpPurchaseplanUser,
		})

		purchaseOrderResponse := &dto.PurchaseOrderResponse{
			ID:   purchaseOrder.Ponumber,
			Code: purchaseOrder.Ponumber,
			Vendor: &dto.VendorResponse{
				ID:   vendor.Data[0].VENDORID,
				Code: vendor.Data[0].VENDORID,
				VendorOrganization: &dto.VendorOrganizationResponse{
					ID:          vendor.Data[0].Organization.PRP_Vendor_Org_ID,
					Code:        vendor.Data[0].Organization.PRP_Vendor_Org_ID,
					Description: vendor.Data[0].Organization.PRP_Vendor_Org_Desc,
				},
				Name: vendor.Data[0].VENDNAME,
			},
			TermPaymentPur: &dto.PurchaseTermResponse{},
			VendorClassification: &dto.VendorClassificationResponse{
				ID:          vendor.Data[0].Classification.PRP_Vendor_CLASF_ID,
				Code:        vendor.Data[0].Classification.PRP_Vendor_CLASF_ID,
				Description: vendor.Data[0].Classification.PRP_Vendor_CLASF_Desc,
			},
			Status:               statusGP,
			DocDate:              docDate,
			RecognitionDate:      recognitionDate,
			EtaDate:              etaDate,
			SiteAddress:          purchaseOrder.PrpLocncode,
			EtaTime:              etaTime.Format("15:04"),
			TaxPct:               purchaseOrder.Obtaxamt,
			TotalSku:             len(purchaseOrderItems),
			TonnagePurchaseOrder: tonnagePurchaseOrderResponse,
			PurchaseOrderItems:   purchaseOrderItemsResponse,
			CreatedBy: &dto.UserResponse{
				ID:           accountName.Data.Id,
				Name:         accountName.Data.Name,
				EmployeeCode: purchaseOrder.PrpPurchaseplanUser},
			AssignedTo: &dto.UserResponse{
				ID:           accountName.Data.Id,
				Name:         accountName.Data.Name,
				EmployeeCode: purchaseOrder.PrpPurchaseplanUser},
		}

		switch purchaseOrder.Postatus {
		case 1:
			purchaseOrderResponse.Status = int32(statusx.ConvertStatusName(statusx.Draft))
		case 4:
			purchaseOrderResponse.Status = int32(statusx.ConvertStatusName(statusx.Active))
		case 8:
			purchaseOrderResponse.Status = int32(statusx.ConvertStatusName(statusx.Closed)) // change Cancelled into Closed
		case 12:
			purchaseOrderResponse.Status = int32(statusx.ConvertStatusName(statusx.Finished))
		case 11:
			purchaseOrderResponse.Status = int32(statusx.ConvertStatusName(statusx.Closed))
		}

		if purchaseOrder.PrpPpReference != "" {
			var purchasePlan *bridgeService.GetPurchasePlanGPResponse
			purchasePlan, err = s.opt.Client.BridgeServiceGrpc.GetPurchasePlanGPDetail(ctx, &bridgeService.GetPurchasePlanGPDetailRequest{
				Id: purchaseOrder.PrpPpReference,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "purchasePlan")
				return
			}

			purchaseOrderResponse.PurchasePlan = &dto.PurchasePlanResponse{
				ID:   purchasePlan.Data[0].PrpPurchaseplanNo,
				Code: purchasePlan.Data[0].PrpPurchaseplanNo,
			}
		}

		// get site
		if purchaseOrder.PrpLocncode != "" {
			var site *bridgeService.GetSiteGPResponse
			site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
				Id: purchaseOrder.PrpLocncode,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "site")
				return
			}

			purchaseOrderResponse.Site = &dto.SiteResponse{
				ID:          site.Data[0].Locncode,
				Code:        site.Data[0].Locncode,
				Description: site.Data[0].Locndscr,
			}
		}

		// add flag is consolidated
		var consolidatedShipment *model.ConsolidatedShipment
		if purchaseOrder.PrpCsReference != "" {
			consolidatedShipment, err = s.RepositoryConsolidatedShipment.GetByPurchaseOrderID(ctx, purchaseOrder.PrpCsReference)
			if err != nil {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
			purchaseOrderResponse.ConsolidatedShipment = &dto.ConsolidatedShipmentResponse{
				ID:                consolidatedShipment.ID,
				Code:              consolidatedShipment.Code,
				DriverName:        consolidatedShipment.DriverName,
				VehicleNumber:     consolidatedShipment.VehicleNumber,
				DriverPhoneNumber: consolidatedShipment.DriverPhoneNumber,
				DeltaPrint:        consolidatedShipment.DeltaPrint,
				Status:            consolidatedShipment.Status,
				CreatedAt:         consolidatedShipment.CreatedAt,
			}
		} else {
			consolidatedShipment = nil
		}

		res = append(res, purchaseOrderResponse)
	}
	total = int64(len(res))

	return
}

func (s *PurchaseOrderService) GetByID(ctx context.Context, id string) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.GetByID")
	defer span.End()

	var purchaseOrder *bridgeService.GetPurchaseOrderGPResponse
	purchaseOrder, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPDetail(ctx, &bridgeService.GetPurchaseOrderGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase_order")
		return
	}

	var vendor *bridgeService.GetVendorGPResponse
	vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
		Id: purchaseOrder.Data[0].Vendorid,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
		return
	}

	// get vendor organization purchase plan
	// var vendorOrganizationPurchasePlan *bridgeService.GetVendorOrganizationDetailResponse
	// vendorOrganizationPurchasePlan, err = s.opt.Client.BridgeServiceGrpc.GetVendorOrganizationDetail(ctx, &bridgeService.GetVendorOrganizationDetailRequest{
	// 	Id: purchaseOrder.PurchasePlanID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "vendor_organization")
	// 	return
	// }

	var purchaseOrderItemsResponse []*dto.PurchaseOrderItemResponse

	purchaseOrderItems := purchaseOrder.Data[0].PoDetail

	var totalPrice float64
	for _, poi := range purchaseOrderItems {

		var item *bridgeService.GetItemGPResponse
		item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: poi.Itemnmbr,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item")
			return
		}

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

		poiResponse := &dto.PurchaseOrderItemResponse{
			ID:              poi.Ponumber,
			PurchaseOrderID: poi.Ponumber,
			Item: &dto.ItemResponse{
				ID:                   item.Data[0].Itemnmbr,
				Code:                 item.Data[0].Itemnmbr,
				Description:          item.Data[0].Itemdesc,
				UnitWeightConversion: item.Data[0].GnlWeighttolerance,
				OrderMinQty:          item.Data[0].Minorqty,
				OrderMaxQty:          item.Data[0].Maxordqty,
				ItemType:             item.Data[0].ItemTypeDesc,
				Uom: &dto.UomResponse{
					ID:   uom.Data[0].Uomschdl,
					Code: uom.Data[0].Uomschdl,
					Name: uom.Data[0].Umschdsc,
				},
			},
			OrderQty:  poi.Qtyorder,
			UnitPrice: poi.Unitcost,
			// TaxableItem:   poi.TaxableItem,
			// IncludeTax:    poi.IncludeTax,
			// TaxPercentage: poi.TaxPercentage,
			// TaxAmount:     poi.TaxAmount,
			// UnitPriceTax:  poi.UnitPriceTax,
			Subtotal: poi.Qtyorder * poi.Unitcost,
			// Weight:      poi.Weight,
			// Note:        poi.Note,
			PurchaseQty: poi.Qtyorder,
		}
		totalPrice += poiResponse.Subtotal

		if purchaseOrder.Data[0].PrpPpReference != "" {
			var purchasePlan *bridgeService.GetPurchasePlanGPResponse
			purchasePlan, err = s.opt.Client.BridgeServiceGrpc.GetPurchasePlanGPDetail(ctx, &bridgeService.GetPurchasePlanGPDetailRequest{
				Id: purchaseOrder.Data[0].PrpPpReference,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "purchasePlan")
				return
			}

			var ppQty float64

			for _, ppItem := range purchasePlan.Data[0].Detail {
				ppQty += ppItem.Quantity
			}

			poiResponse.PurchasePlanItem = &dto.PurchasePlanItemResponse{
				OrderQty:    ppQty,
				PurchaseQty: ppQty,
			}
		}

		purchaseOrderItemsResponse = append(purchaseOrderItemsResponse, poiResponse)
	}

	// // get vendor organization
	// var vendorOrganization *bridgeService.GetVendorOrganizationDetailResponse
	// vendorOrganization, err = s.opt.Client.BridgeServiceGrpc.GetVendorOrganizationDetail(ctx, &bridgeService.GetVendorOrganizationDetailRequest{
	// 	Id: vendor.Data.VendorOrganizationId,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "vendor_organization")
	// 	return
	// }

	// read delta print from purchase order
	var purchaseOrderPrint *model.PurchaseOrder
	purchaseOrderPrint, err = s.RepositoryPurchaseOrder.GetByID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var poSignatures []*model.PurchaseOrderSignature
	poSignatures, _, err = s.RepositoryPurchaseOrderSignature.GetSignatureByPurchaseOrderID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var purchaseOrderSignaturesResponse []*dto.PurchaseOrderSignatureResponse
	for _, pos := range poSignatures {
		purchaseOrderSignaturesResponse = append(purchaseOrderSignaturesResponse, &dto.PurchaseOrderSignatureResponse{
			ID:           pos.ID,
			JobFunction:  pos.JobFunction,
			Name:         pos.Name,
			SignatureURL: pos.SignatureURL,
			CreatedAt:    pos.CreatedAt,
		})
	}

	var poImages []*model.PurchaseOrderImage
	poImages, _, err = s.RepositoryPurchaseOrderImage.GetImageByPurchaseOrderID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var purchaseOrderImagesResponse []*dto.PurchaseOrderImageResponse
	for _, poImage := range poImages {
		purchaseOrderImagesResponse = append(purchaseOrderImagesResponse, &dto.PurchaseOrderImageResponse{
			ID:        poImage.ID,
			ImageURL:  poImage.ImageURL,
			CreatedAt: poImage.CreatedAt,
		})
	}

	var recognitionDate time.Time
	recognitionDate, err = time.Parse("2006-01-02T15:04:05", purchaseOrder.Data[0].Reqdate)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("recognition_date")
		return
	}

	var etaDate time.Time
	etaDate, err = time.Parse("2006-01-02T15:04:05", purchaseOrder.Data[0].PrpEstimatedarrivalDat)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("eta_date")
		return
	}

	var etaTime time.Time
	etaTime, err = time.Parse("2006-01-02T15:04:05", purchaseOrder.Data[0].PrpEstimatedarrivalTim)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("eta_time")
		return
	}

	var statusGP int32
	switch purchaseOrder.Data[0].Postatus {
	case 4:
		statusGP = 1
	case 1:
		statusGP = 5
	case 8:
		statusGP = 33 // change 3 into 33
	case 12:
		statusGP = 2
	case 11:
		statusGP = 33
	}
	accountName, err := s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &account_service.GetUserDetailRequest{
		EmployeeCode: purchaseOrder.Data[0].PrpPurchaseplanUser,
	})
	res = &dto.PurchaseOrderResponse{
		ID:   purchaseOrder.Data[0].Ponumber,
		Code: purchaseOrder.Data[0].Ponumber,
		Vendor: &dto.VendorResponse{
			ID:   vendor.Data[0].VENDORID,
			Code: vendor.Data[0].VENDORID,
			VendorOrganization: &dto.VendorOrganizationResponse{
				ID:          vendor.Data[0].Organization.PRP_Vendor_Org_ID,
				Code:        vendor.Data[0].Organization.PRP_Vendor_Org_ID,
				Description: vendor.Data[0].Organization.PRP_Vendor_Org_Desc,
			},
			VendorClassification: &dto.VendorClassificationResponse{
				ID:          vendor.Data[0].Classification.PRP_Vendor_CLASF_ID,
				Code:        vendor.Data[0].Classification.PRP_Vendor_CLASF_ID,
				Description: vendor.Data[0].Classification.PRP_Vendor_CLASF_Desc,
			},
			AdmDivision:    &dto.AdmDivisionResponse{},
			PaymentTerm:    &dto.PaymentTermResponse{},
			PaymentMethod:  &dto.PaymentMethodResponse{},
			Name:           vendor.Data[0].VENDNAME,
			PicName:        vendor.Data[0].VNDCNTCT,
			PhoneNumber:    vendor.Data[0].PHNUMBR1,
			PhoneNumberAlt: vendor.Data[0].PHNUMBR2,
		},
		DeltaPrint:           purchaseOrderPrint.DeltaPrint,
		TermPaymentPur:       &dto.PurchaseTermResponse{},
		VendorClassification: &dto.VendorClassificationResponse{},
		Status:               statusGP,
		RecognitionDate:      recognitionDate,
		EtaDate:              etaDate,
		SiteAddress:          purchaseOrder.Data[0].PrpLocncode,
		EtaTime:              etaTime.Format("15:04"),
		TaxPct:               purchaseOrder.Data[0].Obtaxamt,
		// DeliveryFee:            purchaseOrder.Data[0].DeliveryFee,
		TotalPrice: totalPrice,
		// TaxAmount:              purchaseOrder.Data[0].TaxAmount,
		// TotalCharge:            purchaseOrder.Data[0].TotalCharge,
		// TotalInvoice:           purchaseOrder.Data[0].TotalInvoice,
		// TotalWeight:            purchaseOrder.Data[0].TotalWeight,
		Note: purchaseOrder.Data[0].Commntid,
		// DeltaPrint:             purchaseOrder.Data[0].DeltaPrint,
		// Latitude:               purchaseOrder.Data[0].Latitude,
		// Longitude:              purchaseOrder.Data[0].Longitude,
		// CreatedFrom:            purchaseOrder.Data[0].CreatedFrom,
		// HasFinishedGr:          purchaseOrder.Data[0].HasFinishedGr,
		// CreatedAt:              purchaseOrder.Data[0].CreatedAt,
		// CommittedAt:            purchaseOrder.Data[0].CommittedAt,
		// AssignedAt:             purchaseOrder.Data[0].AssignedAt,
		// UpdatedAt:              purchaseOrder.Data[0].UpdatedAt,
		// Locked:                 purchaseOrder.Data[0].Locked,
		CreatedBy: &dto.UserResponse{
			ID:           accountName.Data.Id,
			Name:         accountName.Data.Name,
			EmployeeCode: purchaseOrder.Data[0].PrpPurchaseplanUser},
		AssignedTo: &dto.UserResponse{
			ID:           accountName.Data.Id,
			Name:         accountName.Data.Name,
			EmployeeCode: purchaseOrder.Data[0].PrpPurchaseplanUser},
		PurchaseOrderItems:     purchaseOrderItemsResponse,
		PurchaseOrderSignature: purchaseOrderSignaturesResponse,
		PurchaseOrderImage:     purchaseOrderImagesResponse,
	}

	if purchaseOrder.Data[0].PrpPpReference != "" {
		var purchasePlan *bridgeService.GetPurchasePlanGPResponse
		purchasePlan, err = s.opt.Client.BridgeServiceGrpc.GetPurchasePlanGPDetail(ctx, &bridgeService.GetPurchasePlanGPDetailRequest{
			Id: purchaseOrder.Data[0].PrpPpReference,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "purchasePlan")
			return
		}

		res.PurchasePlan = &dto.PurchasePlanResponse{
			ID:   purchasePlan.Data[0].PrpPurchaseplanNo,
			Code: purchasePlan.Data[0].PrpPurchaseplanNo,
		}
	}

	// get site
	if purchaseOrder.Data[0].PrpLocncode != "" {
		var site *bridgeService.GetSiteGPResponse
		site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
			Id: purchaseOrder.Data[0].PrpLocncode,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}

		res.Site = &dto.SiteResponse{
			ID:          site.Data[0].Locncode,
			Code:        site.Data[0].Locncode,
			Description: site.Data[0].Locndscr,
		}
	}

	var consolidatedSipment *model.ConsolidatedShipment
	consolidatedSipment, _ = s.RepositoryConsolidatedShipment.GetByPurchaseOrderID(ctx, purchaseOrder.Data[0].Ponumber)
	if consolidatedSipment.ID != 0 {
		res.ConsolidatedShipment = &dto.ConsolidatedShipmentResponse{
			ID:                consolidatedSipment.ID,
			Code:              consolidatedSipment.Code,
			DriverName:        consolidatedSipment.DriverName,
			VehicleNumber:     consolidatedSipment.VehicleNumber,
			DriverPhoneNumber: consolidatedSipment.DriverPhoneNumber,
			DeltaPrint:        consolidatedSipment.DeltaPrint,
			Status:            consolidatedSipment.Status,
			CreatedAt:         consolidatedSipment.CreatedAt,
		}
	}

	// var createdBy *accountService.GetUserDetailResponse
	// createdBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	Id: purchaseOrder.CreatedBy,
	// })
	// if createdBy != nil {
	// 	res.CreatedBy = &dto.UserResponse{
	// 		ID:   createdBy.Data.Id,
	// 		Name: createdBy.Data.Name,
	// 	}
	// }

	// var committedBy *accountService.GetUserDetailResponse
	// committedBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	Id: purchaseOrder.CommittedBy,
	// })
	// if committedBy != nil {
	// 	res.CommittedBy = &dto.UserResponse{
	// 		ID:   committedBy.Data.Id,
	// 		Name: committedBy.Data.Name,
	// 	}
	// }

	// var assignedTo *accountService.GetUserDetailResponse
	// assignedTo, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	Id: purchaseOrder.AssignedTo,
	// })
	// if assignedTo != nil {
	// 	res.AssignedTo = &dto.UserResponse{
	// 		ID:   assignedTo.Data.Id,
	// 		Name: assignedTo.Data.Name,
	// 	}
	// }

	// var assignedBy *accountService.GetUserDetailResponse
	// assignedBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	Id: purchaseOrder.AssignedBy,
	// })
	// if assignedBy != nil {
	// 	res.AssignedBy = &dto.UserResponse{
	// 		ID:   assignedBy.Data.Id,
	// 		Name: assignedBy.Data.Name,
	// 	}
	// }

	// var updatedBy *accountService.GetUserDetailResponse
	// updatedBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	Id: purchaseOrder.UpdatedBy,
	// })
	// if updatedBy != nil {
	// 	res.UpdatedBy = &dto.UserResponse{
	// 		ID:   updatedBy.Data.Id,
	// 		Name: updatedBy.Data.Name,
	// 	}
	// }

	// var lockedBy *accountService.GetUserDetailResponse
	// lockedBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	Id: purchaseOrder.LockedBy,
	// })
	// if lockedBy != nil {
	// 	res.LockedBy = &dto.UserResponse{
	// 		ID:   lockedBy.Data.Id,
	// 		Name: lockedBy.Data.Name,
	// 	}
	// }

	return
}

func (s *PurchaseOrderService) Create(ctx context.Context, req dto.PurchaseOrderRequestCreate) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Create")
	defer span.End()

	// check purchase plan id

	// check vendor id
	var vendor *bridgeService.GetVendorGPResponse
	vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
		Id: req.VendorID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
		return
	}

	if len(vendor.Data) == 0 {
		err = edenlabs.ErrorInvalid("vendor")
		return
	}

	// check purchase term id

	// check site id
	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: req.SiteID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	if len(site.Data) == 0 {
		err = edenlabs.ErrorInvalid("site")
		return
	}

	if req.RecognitionDate == "" {
		err = edenlabs.ErrorInvalid("recognition_date")
		return
	}

	if req.EtaDate == "" {
		err = edenlabs.ErrorInvalid("eta_date")
		return
	}

	if req.PRStatus == 0 {
		err = edenlabs.ErrorInvalid("pr status")
		return
	}

	// userID := ctx.Value(constants.KeyUserID).(int64)

	var poDetail []*bridgeService.CreatePurchaseOrderGPRequest_Podtl
	var totalAll float64

	ppDetail, _ := s.opt.Client.BridgeServiceGrpc.GetPurchasePlanGPDetail(ctx, &bridgeService.GetPurchasePlanGPDetailRequest{
		Id: req.PurchasePlanID,
	})

	if len(ppDetail.Data[0].Detail) < len(req.PurchaseOrderItems) {
		err = edenlabs.ErrorValidation("requested_item", "requested purchase order item more than requested purchase plan item.")
	}

	purchaseOrder, err := s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPList(ctx, &bridgeService.GetPurchaseOrderGPListRequest{
		Limit:          1000,
		Offset:         0,
		PurchasePlanId: req.PurchasePlanID,
		Status:         4,
	})

	var tempPurchasePlanItem []*dto.PurchasePlanItemResponse
	for _, ppi := range ppDetail.Data[0].Detail {
		for _, v := range purchaseOrder.Data {
			for _, x := range v.PoDetail {
				if x.Itemnmbr == ppi.Itemnmbr {
					tempPurchasePlanItem = append(tempPurchasePlanItem, &dto.PurchasePlanItemResponse{
						Item: &dto.ItemResponse{
							Code: ppi.Itemnmbr,
						},
						OrderQty: x.Qtyorder,
					})
				}
			}
		}
	}
	for _, poi := range req.PurchaseOrderItems {
		qtyPOItem := From(tempPurchasePlanItem).WhereT(
			func(f *dto.PurchasePlanItemResponse) bool {
				return (f.Item.Code == poi.ItemID)
			},
		).SelectT(func(f *dto.PurchasePlanItemResponse) float64 {
			return f.OrderQty
		}).SumFloats()
		orderPO := poi.OrderQty + qtyPOItem

		existItem := From(ppDetail.Data[0].Detail).WhereT(
			func(f *bridgeService.PurchasePlanGP_Detail) bool {
				return (f.Itemnmbr == poi.ItemID)
			},
		).Count()

		// continue
		exist := From(ppDetail.Data[0].Detail).WhereT(
			func(f *bridgeService.PurchasePlanGP_Detail) bool {
				return (f.Itemnmbr == poi.ItemID) && orderPO <= f.Quantity
			},
		).Count()
		// get item
		var item *bridgeService.GetItemGPResponse
		item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: poi.ItemID,
		})

		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}
		if exist == 0 {
			err = edenlabs.ErrorValidation("item_order_qty", "item "+item.Data[0].Itemdesc+" order quantity more than quantity in purchase plan.")
			return
			//break
		}

		if existItem == 0 {
			err = edenlabs.ErrorValidation("item_order_qty", "item "+item.Data[0].Itemdesc+" does not exist in purchase plan.")
			return
			//break
		}

		if poi.OrderQty <= 0 {
			err = edenlabs.ErrorInvalid("order_qty")
		}

		if poi.UnitPrice < 0 {
			err = edenlabs.ErrorInvalid("unit_price")
		}

		if poi.TaxPercentage < 0 {
			err = edenlabs.ErrorInvalid("tax_percentage")
		}

		// taxAmount := math.Round((poi.UnitPrice * poi.TaxPercentage / 100) * poi.OrderQty)
		// unitPriceTax := math.Round(poi.UnitPrice * (100 + poi.TaxPercentage) / 100)

		// subtotal := poi.OrderQty * item.Data[0].Currcost
		subtotal := poi.OrderQty * poi.UnitPrice

		totalAll += subtotal
		// if poi.PurchasePlanItemID != 0 {
		// 	poi.PurchaseQty = poi.OrderQty
		// 	poi.IncludeTax = 2
		// 	poi.TaxableItem = 2
		// 	poi.TaxPercentage = 0
		// }

		poDetail = append(poDetail, &bridgeService.CreatePurchaseOrderGPRequest_Podtl{
			Ord:      0,
			Itemnmbr: item.Data[0].Itemnmbr,
			Uofm:     item.Data[0].Uomschdl,
			Qtyorder: poi.OrderQty,
			Qtycance: 0,
			Unitcost: poi.UnitPrice,
			Notetext: poi.Note,
		})

	}

	poCreateReq := &bridgeService.CreatePurchaseOrderGPRequest{
		Potype:                  1,
		Ponumber:                "",
		Docdate:                 req.RecognitionDate,
		Buyerid:                 "",
		Vendorid:                vendor.Data[0].VENDORID,
		Curncyid:                "",
		Deprtmnt:                "PUR",
		Locncode:                site.Data[0].Locncode,
		Taxschid:                "",
		Subtotal:                totalAll,
		Trdisamt:                0,
		Frtamnt:                 0,
		Miscamnt:                0,
		Taxamnt:                 0,
		PrpPurchaseplanNo:       req.PurchasePlanID,
		PrpPaymentMethod:        vendor.Data[0].PaymentMethod.PRP_Payment_Method,
		Pymtrmid:                req.PaymentTermID,
		Duedate:                 "1900-01-01",
		PrpRegion:               req.RegionID,
		PrpEstimatedarrivalDate: req.EtaDate,
		Notetext:                req.Note,
		Detail:                  poDetail,
		PrStatus:                req.PRStatus,
	}

	var createPurchaseOrderGPResponse *bridgeService.CreatePurchaseOrderGPResponse
	createPurchaseOrderGPResponse, err = s.opt.Client.BridgeServiceGrpc.CreatePurchaseOrderGP(ctx, poCreateReq)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "create_purchase_order")
		return
	}

	// Check if PO is already in Eden table with 2 times loop
	for i := 0; i < 2; i++ {
		var purchaseOrder *model.PurchaseOrder

		purchaseOrder, err = s.RepositoryPurchaseOrder.GetByID(ctx, createPurchaseOrderGPResponse.Ponumber)
		if err != nil {
			poCreate := &model.PurchaseOrder{
				PurchaseOrderIDGP: purchaseOrder.PurchaseOrderIDGP,
				SiteIDGP:          site.Data[0].Locncode,
			}

			err = s.RepositoryPurchaseOrder.Create(ctx, poCreate)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			for _, image := range req.Images {
				err = s.RepositoryPurchaseOrderImage.Create(ctx, &model.PurchaseOrderImage{
					PurchaseOrderIDGP: createPurchaseOrderGPResponse.Ponumber,
					ImageURL:          image,
					CreatedAt:         time.Now(),
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
			}
		} else {
			break
		}
	}

	res = &dto.PurchaseOrderResponse{
		ID:   createPurchaseOrderGPResponse.Ponumber,
		Code: createPurchaseOrderGPResponse.Ponumber,
		Vendor: &dto.VendorResponse{
			ID:             vendor.Data[0].VENDORID,
			Code:           vendor.Data[0].VENDORID,
			Name:           vendor.Data[0].VENDNAME,
			PicName:        vendor.Data[0].VNDCNTCT,
			PhoneNumber:    vendor.Data[0].PHNUMBR1,
			PhoneNumberAlt: vendor.Data[0].PHNUMBR2,
		},
		Site: &dto.SiteResponse{
			ID:          site.Data[0].Locncode,
			Code:        site.Data[0].Locncode,
			Description: site.Data[0].Locndscr,
		},
		TermPaymentPur: &dto.PurchaseTermResponse{},
		PurchasePlan:   &dto.PurchasePlanResponse{},
	}

	return
}

func (s *PurchaseOrderService) Update(ctx context.Context, req dto.PurchaseOrderRequestUpdate, id string) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Update")
	defer span.End()

	// check purchase order id
	var purchaseOrder *bridgeService.GetPurchaseOrderGPResponse
	purchaseOrder, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPDetail(ctx, &bridgeService.GetPurchaseOrderGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchase_order")
		return
	}

	// check site id
	// var site *bridgeService.GetSiteGPResponse
	// site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
	// 	Id: purchaseOrder.Data[0].PrpLocncode,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "site")
	// 	return
	// }

	// userID := ctx.Value(constants.KeyUserID).(int64)

	var purchaseOrderItemsResponse []*dto.PurchaseOrderItemResponse

	var poDetail []*bridgeService.UpdatePurchaseOrderGPRequest_Podtl
	var totalAll float64

	for _, poi := range req.PurchaseOrderItems {

		// get item
		var item *bridgeService.GetItemGPResponse
		item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: poi.ItemID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}

		if poi.OrderQty <= 0 {
			err = edenlabs.ErrorInvalid("order_qty")
		}

		if poi.UnitPrice < 0 {
			err = edenlabs.ErrorInvalid("unit_price")
		}

		if poi.TaxPercentage < 0 {
			err = edenlabs.ErrorInvalid("tax_persentage")
		}

		// taxAmount := math.Round((poi.UnitPrice * poi.TaxPercentage / 100) * poi.OrderQty)
		// unitPriceTax := math.Round(poi.UnitPrice * (100 + poi.TaxPercentage) / 100)

		// subtotal := poi.OrderQty * poi.UnitPrice
		subtotal := poi.OrderQty * item.Data[0].Currcost

		totalAll += subtotal

		poi.PurchaseQty = poi.OrderQty
		poi.IncludeTax = 2
		poi.TaxableItem = 2
		poi.TaxPercentage = 0

		poDetail = append(poDetail, &bridgeService.UpdatePurchaseOrderGPRequest_Podtl{
			Ord:      0,
			Itemnmbr: item.Data[0].Itemnmbr,
			Uofm:     item.Data[0].Uomschdl,
			Qtyorder: poi.OrderQty,
			Qtycance: 0,
			Unitcost: item.Data[0].Currcost,
			Notetext: poi.Note,
		})

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

		purchaseOrderItemsResponse = append(purchaseOrderItemsResponse, &dto.PurchaseOrderItemResponse{
			ID:              id,
			PurchaseOrderID: id,
			Item: &dto.ItemResponse{
				ID:                   item.Data[0].Itemnmbr,
				Code:                 item.Data[0].Itemnmbr,
				Description:          item.Data[0].Itemdesc,
				UnitWeightConversion: item.Data[0].GnlWeighttolerance,
				OrderMinQty:          item.Data[0].Minorqty,
				OrderMaxQty:          item.Data[0].Maxordqty,
				ItemType:             item.Data[0].ItemTypeDesc,
				Uom: &dto.UomResponse{
					ID:   uom.Data[0].Uomschdl,
					Code: uom.Data[0].Uomschdl,
					Name: uom.Data[0].Umschdsc,
				},
			},
			OrderQty:    0,
			UnitPrice:   0,
			Subtotal:    0,
			PurchaseQty: 0,
		})
	}

	poUpdateReq := &bridgeService.UpdatePurchaseOrderGPRequest{
		Potype:                  1,
		Ponumber:                purchaseOrder.Data[0].Ponumber,
		Docdate:                 purchaseOrder.Data[0].Reqdate,
		Buyerid:                 purchaseOrder.Data[0].Buyerid,
		Vendorid:                purchaseOrder.Data[0].Vendorid,
		Curncyid:                "",
		Deprtmnt:                purchaseOrder.Data[0].Deprtmnt,
		Locncode:                purchaseOrder.Data[0].PrpLocncode,
		Taxschid:                "",
		Subtotal:                totalAll,
		Trdisamt:                0,
		Frtamnt:                 0,
		Miscamnt:                0,
		Taxamnt:                 0,
		PrpPurchaseplanNo:       purchaseOrder.Data[0].PrpPpReference,
		PrpPaymentMethod:        purchaseOrder.Data[0].PrpPaymentMethod,
		PrpRegion:               purchaseOrder.Data[0].PrpRegion,
		PrpEstimatedarrivalDate: purchaseOrder.Data[0].PrpEstimatedarrivalDat,
		Notetext:                purchaseOrder.Data[0].Commntid,
		Detail:                  poDetail,
	}

	var updatePurchaseOrderGPResponse *bridgeService.CreatePurchaseOrderGPResponse
	updatePurchaseOrderGPResponse, err = s.opt.Client.BridgeServiceGrpc.UpdatePurchaseOrderGP(ctx, poUpdateReq)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "update_purchase_order")
		return
	}

	err = s.RepositoryPurchaseOrderImage.Delete(ctx, updatePurchaseOrderGPResponse.Ponumber)
	for _, image := range req.Images {
		err = s.RepositoryPurchaseOrderImage.Create(ctx, &model.PurchaseOrderImage{
			PurchaseOrderIDGP: updatePurchaseOrderGPResponse.Ponumber,
			ImageURL:          image,
			CreatedAt:         time.Now(),
		})
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	res = &dto.PurchaseOrderResponse{
		Site:               &dto.SiteResponse{},
		TermPaymentPur:     &dto.PurchaseTermResponse{},
		PurchasePlan:       &dto.PurchasePlanResponse{},
		PurchaseOrderItems: purchaseOrderItemsResponse,
	}

	return
}

func (s *PurchaseOrderService) Assign(ctx context.Context, req dto.PurchaseOrderRequestAssign, id string) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Assign")
	defer span.End()

	// check purchase order id
	var purchaseOrder *model.PurchaseOrder
	// purchaseOrder, err = s.RepositoryPurchaseOrder.GetByID(ctx, id)
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorInvalid("id")
	// 	return
	// }

	// validate status must draft
	// if purchaseOrder.Status != 5 {
	// 	err = edenlabs.ErrorMustDraft("status")
	// 	return
	// }

	// var fieldPurchaser *accountService.GetUserDetailResponse
	// fieldPurchaser, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	Id: req.FieldPurchaserID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorInvalid("field_purchaser_id")
	// 	return
	// }

	// userID := ctx.Value(constants.KeyUserID).(int64)

	purchaseOrder = &model.PurchaseOrder{
		// ID:         id,
		// AssignedTo: fieldPurchaser.Data.Id,
		// AssignedBy: userID,
		// AssignedAt: time.Now(),
	}

	err = s.RepositoryPurchaseOrder.Update(ctx, purchaseOrder, "AssignedTo", "AssignedBy", "AssignedAt")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// send notification

	res = &dto.PurchaseOrderResponse{
		// ID: purchaseOrder.ID,
	}

	return
}

func (s *PurchaseOrderService) Signature(ctx context.Context, req dto.PurchaseOrderRequestSignature, id string) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Signature")
	defer span.End()

	var totalSignature int64
	_, totalSignature, err = s.RepositoryPurchaseOrderSignature.GetSignatureByPurchaseOrderID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	if totalSignature >= 4 {
		err = edenlabs.ErrorMustEqualOrLess("signature", "4 times")
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	purchaseOrderSignature := &model.PurchaseOrderSignature{
		PurchaseOrderIDGP: req.PurchaseOrderID,
		JobFunction:       req.JobFunction,
		Name:              req.Name,
		SignatureURL:      req.SignatureURL,
		CreatedAt:         time.Now(),
		CreatedBy:         userID,
	}

	err = s.RepositoryPurchaseOrderSignature.Create(ctx, purchaseOrderSignature)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PurchaseOrderService) Print(ctx context.Context, id string) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Print")
	defer span.End()

	// check purchase order id
	var purchaseOrder *model.PurchaseOrder
	purchaseOrder, err = s.RepositoryPurchaseOrder.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	purchaseOrder.DeltaPrint += 1

	err = s.RepositoryPurchaseOrder.Update(ctx, purchaseOrder, "DeltaPrint")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *PurchaseOrderService) Cancel(ctx context.Context, req dto.PurchaseOrderRequestCancel, id string) (res *dto.PurchaseOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PurchaseOrderService.Cancel")
	defer span.End()

	var purchaseOrder *bridgeService.GetPurchaseOrderGPResponse
	purchaseOrder, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPDetail(ctx, &bridgeService.GetPurchaseOrderGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "purchaseOrder")
		return
	}
	if purchaseOrder.Data[0].PrpCsReference != "" {
		err = edenlabs.ErrorValidation("consolidated_shipment", "Purchase order can not be cancel after consolidated.")
		return
	}
	po, err := s.RepositoryPurchaseOrder.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}
	if po.DeltaPrint > 1 {
		err = edenlabs.ErrorValidation("purchase_order", "Purchase order can not be cancel after printed.")
		return

	}
	// var newPurchaserUser *accountService.GetUserDetailResponse
	// newPurchaserUser, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
	// 	EmployeeCode: req.FieldPurchaserID,
	// })
	// if err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("account", "user")
	// 	return
	// }

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "user")
		return
	}

	var _ *bridgeService.CancelPurchaseOrderGPResponse
	_, err = s.opt.Client.BridgeServiceGrpc.CancelPurchaseOrderGP(ctx, &bridgeService.CancelPurchaseOrderGPRequest{
		PoNumber: id,
		UserId:   "",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "cancel_purchase_order")
		return
	}

	// audit log
	userID := ctx.Value(constants.KeyUserID).(int64)

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      userID,
			ReferenceId: purchaseOrder.Data[0].Ponumber,
			Type:        "cancel_purchase_order",
			Function:    "PurchaseOrderService.Cancel",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	})

	// fmt.Print(purchaseOrder)

	return
}
