package service

import (
	"context"
	"time"

	. "github.com/ahmetb/go-linq/v3"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/repository"
)

type IConsolidatedShipmentService interface {
	Get(ctx context.Context, req *dto.ConsolidatedShipmentRequestList) (res []*dto.ConsolidatedShipmentResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res *dto.ConsolidatedShipmentResponse, err error)
	Create(ctx context.Context, req *dto.ConsolidatedShipmentRequestCreate) (res *dto.ConsolidatedShipmentResponse, err error)
	Update(ctx context.Context, req *dto.ConsolidatedShipmentRequestUpdate, id int64) (res *dto.ConsolidatedShipmentResponse, err error)
	Signature(ctx context.Context, req *dto.ConsolidatedShipmentSignatureRequestCreate) (res *dto.ConsolidatedShipmentSignatureResponse, err error)
	Print(ctx context.Context, id int64) (res *dto.ConsolidatedShipmentResponse, err error)
}

type ConsolidatedShipmentService struct {
	opt                                     opt.Options
	RepositoryConsolidatedShipment          repository.IConsolidatedShipmentRepository
	RepositoryConsolidatedShipmentSignature repository.IConsolidatedShipmentSignatureRepository
	RepositoryPurchaseOrder                 repository.IPurchaseOrderRepository
	RepositoryPurchaseOrderItem             repository.IPurchaseOrderItemRepository
}

func NewConsolidatedShipmentService() IConsolidatedShipmentService {
	return &ConsolidatedShipmentService{
		opt:                                     global.Setup.Common,
		RepositoryConsolidatedShipment:          repository.NewConsolidatedShipmentRepository(),
		RepositoryConsolidatedShipmentSignature: repository.NewConsolidatedShipmentSignatureRepository(),
		RepositoryPurchaseOrder:                 repository.NewPurchaseOrderRepository(),
		RepositoryPurchaseOrderItem:             repository.NewPurchaseOrderItemRepository(),
	}
}

func (s *ConsolidatedShipmentService) Get(ctx context.Context, req *dto.ConsolidatedShipmentRequestList) (res []*dto.ConsolidatedShipmentResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConsolidatedShipmentService.Get")
	defer span.End()

	// check site
	var siteCode string
	if req.SiteID != "" {
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

		siteCode = site.Data[0].Locndscr
	}

	var consolidatedShipments []*model.ConsolidatedShipment
	consolidatedShipments, total, err = s.RepositoryConsolidatedShipment.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, siteCode, req.CreatedAtFrom, req.CreatedAtTo, req.CreatedBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, cs := range consolidatedShipments {

		var createdBy *accountService.GetUserDetailResponse
		createdBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
			Id: cs.CreatedBy,
		})

		res = append(res, &dto.ConsolidatedShipmentResponse{
			ID:                cs.ID,
			Code:              cs.Code,
			DriverName:        cs.DriverName,
			VehicleNumber:     cs.VehicleNumber,
			DriverPhoneNumber: cs.DriverPhoneNumber,
			DeltaPrint:        cs.DeltaPrint,
			Status:            cs.Status,
			CreatedAt:         cs.CreatedAt,
			CreatedBy: &dto.UserResponse{
				ID:    createdBy.Data.Id,
				Name:  createdBy.Data.Name,
				Email: createdBy.Data.Email,
			},
			SiteName: cs.SiteName,
			// ConsolidatedShipmentSignatures: cs.ConsolidatedShipmentSignatures,
		})
	}

	return
}

func (s *ConsolidatedShipmentService) GetByID(ctx context.Context, id int64) (res *dto.ConsolidatedShipmentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConsolidatedShipmentService.GetByID")
	defer span.End()

	var consolidatedShipment *model.ConsolidatedShipment
	consolidatedShipment, err = s.RepositoryConsolidatedShipment.GetByID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var createdBy *accountService.GetUserDetailResponse
	createdBy, _ = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		Id: consolidatedShipment.CreatedBy,
	})

	var skuSummaries []*dto.SkuSummaryResponse
	var purchaseOrdersResponse []*dto.PurchaseOrderResponse

	var purchaseOrdersConsolidated []*model.PurchaseOrder
	purchaseOrdersConsolidated, _, err = s.RepositoryPurchaseOrder.GetByConsolidatedShipmentID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, purchaseOrder := range purchaseOrdersConsolidated {
		var purchaseOrderItemsResponse []*dto.PurchaseOrderItemResponse

		var purchaseOrderGP *bridgeService.GetPurchaseOrderGPResponse
		purchaseOrderGP, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPDetail(ctx, &bridgeService.GetPurchaseOrderGPDetailRequest{
			Id: purchaseOrder.PurchaseOrderIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "purchase_order")
			return
		}

		purchaseOrderItems := purchaseOrderGP.Data[0].PoDetail

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

			purchaseOrderItemsResponse = append(purchaseOrderItemsResponse, &dto.PurchaseOrderItemResponse{
				ID:              poi.Ponumber,
				PurchaseOrderID: poi.Ponumber,
				// PurchasePlanItemID: poi.PurchasePlanItemID,
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
			})
			exist := From(skuSummaries).WhereT(
				func(f *dto.SkuSummaryResponse) bool {
					return (f.ItemID == item.Data[0].Itemnmbr)
				},
			).Count()
			if exist == 0 {
				temppo := &dto.PurchaseOrderSKUSummary{
					PurchaseOrderCode: poi.Ponumber,
					Qty:               poi.Qtyorder,
				}
				tempSKUSummary := &dto.SkuSummaryResponse{
					ItemID:   item.Data[0].Itemnmbr,
					ItemName: item.Data[0].Itemdesc,
					UomName:  uom.Data[0].Umschdsc,
					TotalQty: poi.Qtyorder,
				}
				tempSKUSummary.PurchaseOrders = append(tempSKUSummary.PurchaseOrders, temppo)
				skuSummaries = append(skuSummaries, tempSKUSummary)
			} else {
				for i := range skuSummaries {
					if skuSummaries[i].ItemID == item.Data[0].Itemnmbr {
						skuSummaries[i].TotalQty = skuSummaries[i].TotalQty + poi.Qtyorder // Update the qty to the desired value
						skuSummaries[i].PurchaseOrders = append(skuSummaries[i].PurchaseOrders, &dto.PurchaseOrderSKUSummary{
							PurchaseOrderCode: poi.Ponumber,
							Qty:               poi.Qtyorder,
						})
						break
					}
				}
			}
		}

		var recognitionDate time.Time
		recognitionDate, err = time.Parse("2006-01-02T15:04:05", purchaseOrderGP.Data[0].Reqdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("recognition_date")
			return
		}

		var etaDate time.Time
		etaDate, err = time.Parse("2006-01-02T15:04:05", purchaseOrderGP.Data[0].PrpEstimatedarrivalDat)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_date")
			return
		}

		var etaTime time.Time
		etaTime, err = time.Parse("2006-01-02T15:04:05", purchaseOrderGP.Data[0].PrpEstimatedarrivalTim)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_time")
			return
		}
		var vendor *bridgeService.GetVendorGPResponse

		vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
			Id: purchaseOrderGP.Data[0].Vendorid,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			return
		}

		if len(vendor.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			return
		}

		purchaseOrdersResponse = append(purchaseOrdersResponse, &dto.PurchaseOrderResponse{
			ID:              purchaseOrderGP.Data[0].Ponumber,
			Code:            purchaseOrderGP.Data[0].Ponumber,
			RecognitionDate: recognitionDate,
			EtaDate:         etaDate,
			SiteAddress:     purchaseOrderGP.Data[0].PrpLocncode,
			EtaTime:         etaTime.Format("15:04"),
			TaxPct:          purchaseOrderGP.Data[0].Obtaxamt,
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
			// DeliveryFee:     purchaseOrder.DeliveryFee,
			// TotalPrice:      purchaseOrder.TotalPrice,
			// TaxAmount:       purchaseOrder.TaxAmount,
			// TotalCharge:        purchaseOrder.TotalCharge,
			// TotalInvoice:       purchaseOrder.TotalInvoice,
			// TotalWeight:        purchaseOrder.TotalWeight,
			// Note:               purchaseOrder.Note,
			// DeltaPrint:         purchaseOrder.DeltaPrint,
			// Latitude:           purchaseOrder.Latitude,
			// Longitude:          purchaseOrder.Longitude,
			PurchaseOrderItems: purchaseOrderItemsResponse,
		})
	}

	var consolidatedShipmentSignatures []*model.ConsolidatedShipmentSignature
	consolidatedShipmentSignatures, _, _ = s.RepositoryConsolidatedShipmentSignature.GetByConsolidatedShipmentID(ctx, consolidatedShipment.ID)

	var consolidatedShipmentSignaturesResponse []*dto.ConsolidatedShipmentSignatureResponse

	for _, css := range consolidatedShipmentSignatures {
		consolidatedShipmentSignaturesResponse = append(consolidatedShipmentSignaturesResponse, &dto.ConsolidatedShipmentSignatureResponse{
			ID:           css.ID,
			JobFunction:  css.JobFunction,
			Name:         css.Name,
			SignatureURL: css.SignatureURL,
			CreatedAt:    css.CreatedAt,
			CreatedBy:    css.CreatedBy,
		})
	}

	res = &dto.ConsolidatedShipmentResponse{
		ID:                consolidatedShipment.ID,
		Code:              consolidatedShipment.Code,
		DriverName:        consolidatedShipment.DriverName,
		VehicleNumber:     consolidatedShipment.VehicleNumber,
		DriverPhoneNumber: consolidatedShipment.DriverPhoneNumber,
		DeltaPrint:        consolidatedShipment.DeltaPrint,
		Status:            consolidatedShipment.Status,
		CreatedAt:         consolidatedShipment.CreatedAt,
		CreatedBy: &dto.UserResponse{
			ID:    createdBy.Data.Id,
			Name:  createdBy.Data.Name,
			Email: createdBy.Data.Email,
		},
		SiteName:                       consolidatedShipment.SiteName,
		ConsolidatedShipmentSignatures: consolidatedShipmentSignaturesResponse,
		PurchaseOrders:                 purchaseOrdersResponse,
		SkuSummaries:                   skuSummaries,
	}

	return
}

func (s *ConsolidatedShipmentService) Create(ctx context.Context, req *dto.ConsolidatedShipmentRequestCreate) (res *dto.ConsolidatedShipmentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConsolidatedShipmentService.Create")
	defer span.End()

	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "CS",
		Domain: "consolidated_shipment",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "generate_code")
		return
	}

	var siteID string
	var vendorOrg string
	for i, po := range req.PurchaseOrders {
		// check purchase order id
		var purchaseOrder *model.PurchaseOrder
		purchaseOrder, err = s.RepositoryPurchaseOrder.GetByID(ctx, po.PurchaseOrderID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRowInvalid("purchase_orders", i, "purchase_order_id")
			return
		}
		var poGP *bridgeService.GetPurchaseOrderGPResponse
		poGP, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPDetail(ctx, &bridgeService.GetPurchaseOrderGPDetailRequest{
			Id: purchaseOrder.PurchaseOrderIDGP,
		})

		// check vendor id
		var vendor *bridgeService.GetVendorGPResponse
		vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
			Id: poGP.Data[0].Vendorid,
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

		if vendorOrg == "" {
			vendorOrg = vendor.Data[0].Organization.PRP_Vendor_Org_ID
		} else {
			if vendorOrg != vendor.Data[0].Organization.PRP_Vendor_Org_ID {
				err = edenlabs.ErrorValidation("vendor_organization", "vendor organization must be same in one consolidated shipment.")
				return
			}
		}

		siteID = purchaseOrder.SiteIDGP
	}

	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: siteID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	consolidatedShipment := &model.ConsolidatedShipment{
		Code:              codeGenerator.Data.Code,
		DriverName:        req.DriverName,
		VehicleNumber:     req.VehicleNumber,
		DriverPhoneNumber: req.DriverPhoneNumber,
		Status:            statusx.ConvertStatusName(statusx.Active),
		CreatedAt:         time.Now(),
		CreatedBy:         userID,
		SiteName:          site.Data[0].Locndscr,
	}

	err = s.RepositoryConsolidatedShipment.Create(ctx, consolidatedShipment)

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	res = &dto.ConsolidatedShipmentResponse{
		ID:                consolidatedShipment.ID,
		Code:              consolidatedShipment.Code,
		DriverName:        consolidatedShipment.DriverName,
		VehicleNumber:     consolidatedShipment.VehicleNumber,
		DriverPhoneNumber: consolidatedShipment.DriverPhoneNumber,
		DeltaPrint:        consolidatedShipment.DeltaPrint,
		Status:            consolidatedShipment.Status,
		CreatedAt:         consolidatedShipment.CreatedAt,
		CreatedBy: &dto.UserResponse{
			ID: userID,
		},
		SiteName: consolidatedShipment.SiteName,
	}

	var poList []*bridgeService.PurchaseOrderConsolidatedShipment

	for _, po := range req.PurchaseOrders {
		var purchaseOrder *model.PurchaseOrder

		purchaseOrder, err = s.RepositoryPurchaseOrder.GetByID(ctx, po.PurchaseOrderID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		err = s.RepositoryPurchaseOrder.Update(ctx, &model.PurchaseOrder{
			ID:                     purchaseOrder.ID,
			ConsolidatedShipmentID: consolidatedShipment.ID,
		}, "ConsolidatedShipmentID")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// Update CS to GP
		poList = append(poList, &bridgeService.PurchaseOrderConsolidatedShipment{
			Ponumber: po.PurchaseOrderID,
		})

	}

	csCreateReq := &bridgeService.CreateConsolidatedShipmentGPRequest{
		PrpCsNo:         codeGenerator.Data.Code,
		PrpDriverName:   req.DriverName,
		PrVehicleNumber: req.VehicleNumber,
		Phonname:        req.DriverPhoneNumber,
		PurchaseOrder:   poList,
	}

	_, err = s.opt.Client.BridgeServiceGrpc.CreateConsolidatedShipmentGP(ctx, csCreateReq)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "create_consolidated_shipment")
		return
	}

	return
}

func (s *ConsolidatedShipmentService) Update(ctx context.Context, req *dto.ConsolidatedShipmentRequestUpdate, id int64) (res *dto.ConsolidatedShipmentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConsolidatedShipmentService.Update")
	defer span.End()

	var csInfo *model.ConsolidatedShipment
	// check existed id
	csInfo, err = s.RepositoryConsolidatedShipment.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	var siteID string
	var vendorOrg string
	for i, po := range req.PurchaseOrders {
		// check purchase order id
		var purchaseOrder *model.PurchaseOrder
		purchaseOrder, err = s.RepositoryPurchaseOrder.GetByID(ctx, po.PurchaseOrderID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRowInvalid("purchase_orders", i, "purchase_order_id")
			return
		}

		var poGP *bridgeService.GetPurchaseOrderGPResponse
		poGP, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPDetail(ctx, &bridgeService.GetPurchaseOrderGPDetailRequest{
			Id: purchaseOrder.PurchaseOrderIDGP,
		})

		// check vendor id
		var vendor *bridgeService.GetVendorGPResponse
		vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
			Id: poGP.Data[0].Vendorid,
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

		if vendorOrg == "" {
			vendorOrg = vendor.Data[0].Organization.PRP_Vendor_Org_ID
		} else {
			if vendorOrg != vendor.Data[0].Organization.PRP_Vendor_Org_ID {
				err = edenlabs.ErrorValidation("vendor_organization", "vendor organization must be same in one consolidated shipment.")
				return
			}
		}
		siteID = purchaseOrder.SiteIDGP
	}

	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: siteID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	consolidatedShipment := &model.ConsolidatedShipment{
		ID:                id,
		Code:              req.CSNo,
		DriverName:        req.DriverName,
		VehicleNumber:     req.VehicleNumber,
		DriverPhoneNumber: req.DriverPhoneNumber,
		Status:            statusx.ConvertStatusName(statusx.Active),
		CreatedAt:         time.Now(),
		CreatedBy:         userID,
		SiteName:          site.Data[0].Locndscr,
	}

	err = s.RepositoryConsolidatedShipment.Update(ctx, consolidatedShipment, "DriverName", "VehicleNumber", "DriverPhoneNumber", "status", "CreatedAt", "CreatedBy", "SiteName")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	var pOs []*model.PurchaseOrder
	pOs, _, err = s.RepositoryPurchaseOrder.GetByConsolidatedShipmentID(ctx, id)

	err = s.RepositoryPurchaseOrder.DeleteConsolidatedShipmentID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, po := range req.PurchaseOrders {
		var purchaseOrder *model.PurchaseOrder
		purchaseOrder, err = s.RepositoryPurchaseOrder.GetByID(ctx, po.PurchaseOrderID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		err = s.RepositoryPurchaseOrder.Update(ctx, &model.PurchaseOrder{
			ID:                     purchaseOrder.ID,
			ConsolidatedShipmentID: id,
		}, "ConsolidatedShipmentID")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	res = &dto.ConsolidatedShipmentResponse{
		ID:                consolidatedShipment.ID,
		Code:              consolidatedShipment.Code,
		DriverName:        consolidatedShipment.DriverName,
		VehicleNumber:     consolidatedShipment.VehicleNumber,
		DriverPhoneNumber: consolidatedShipment.DriverPhoneNumber,
		DeltaPrint:        consolidatedShipment.DeltaPrint,
		Status:            consolidatedShipment.Status,
		CreatedAt:         consolidatedShipment.CreatedAt,
		CreatedBy: &dto.UserResponse{
			ID: userID,
		},
		SiteName: consolidatedShipment.SiteName,
	}

	mPo := make(map[string]int64)
	var poList, poListDel []*bridgeService.PurchaseOrderConsolidatedShipment

	for _, po := range req.PurchaseOrders {

		// ADD INFO CS TO PO GP (UPDATE GP)
		_, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPDetail(ctx, &bridgeService.GetPurchaseOrderGPDetailRequest{
			Id: po.PurchaseOrderID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "purchase order")
			return
		}

		// list out all po belongs to cs
		mPo[po.PurchaseOrderID] = 0

		// Update CS to GP
		poList = append(poList, &bridgeService.PurchaseOrderConsolidatedShipment{
			Ponumber: po.PurchaseOrderID,
		})
	}

	// Update CS to GP
	csUpdateReq := &bridgeService.UpdateConsolidatedShipmentGPRequest{
		PrpCsNo:         csInfo.Code,
		PrpDriverName:   csInfo.DriverName,
		PrVehicleNumber: csInfo.VehicleNumber,
		Phonname:        csInfo.DriverPhoneNumber,
		PurchaseOrder:   poList,
	}

	_, err = s.opt.Client.BridgeServiceGrpc.UpdateConsolidatedShipmentGP(ctx, csUpdateReq)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "update_consolidated_shipment")
		return
	}

	// remove PO deleted from CS to GP
	for _, poDel := range pOs {
		// If the key exists
		_, ok := mPo[poDel.PurchaseOrderIDGP]
		if !ok {
			// Append all deleted CS
			poListDel = append(poListDel, &bridgeService.PurchaseOrderConsolidatedShipment{
				Ponumber: poDel.PurchaseOrderIDGP,
			})
		}
	}

	// Update  CS to GP
	// Check if any deleted PO
	if poListDel != nil {
		csUpdateDelReq := &bridgeService.UpdateConsolidatedShipmentGPRequest{
			PrpCsNo:         "",
			PrpDriverName:   "",
			PrVehicleNumber: "",
			Phonname:        "",
			PurchaseOrder:   poListDel,
		}

		_, err = s.opt.Client.BridgeServiceGrpc.UpdateConsolidatedShipmentGP(ctx, csUpdateDelReq)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "update_consolidated_shipment")
			return
		}
	}
	return
}

func (s *ConsolidatedShipmentService) Signature(ctx context.Context, req *dto.ConsolidatedShipmentSignatureRequestCreate) (res *dto.ConsolidatedShipmentSignatureResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConsolidatedShipmentService.Signature")
	defer span.End()

	// check existed id
	_, err = s.RepositoryConsolidatedShipment.GetByID(ctx, req.ConsolidatedShipmentID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	var totalSignature int64
	_, totalSignature, err = s.RepositoryConsolidatedShipmentSignature.GetByConsolidatedShipmentID(ctx, req.ConsolidatedShipmentID)
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

	// check if role already sign
	var isRoleAlreadySigned bool
	isRoleAlreadySigned, _ = s.RepositoryConsolidatedShipmentSignature.CheckAlreadySigned(ctx, req.ConsolidatedShipmentID, req.JobFunction)
	if !isRoleAlreadySigned {
		err = edenlabs.ErrorValidation("consolidate_shipment_id", "Job function "+req.JobFunction+" already signed by "+req.Name)
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	consolidatedShipmentSignature := &model.ConsolidatedShipmentSignature{
		ConsolidatedShipmentID: req.ConsolidatedShipmentID,
		JobFunction:            req.JobFunction,
		Name:                   req.Name,
		SignatureURL:           req.SignatureURL,
		CreatedAt:              time.Now(),
		CreatedBy:              userID,
	}

	err = s.RepositoryConsolidatedShipmentSignature.Create(ctx, consolidatedShipmentSignature)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	res = &dto.ConsolidatedShipmentSignatureResponse{
		ID:           consolidatedShipmentSignature.ID,
		JobFunction:  consolidatedShipmentSignature.JobFunction,
		Name:         consolidatedShipmentSignature.Name,
		SignatureURL: consolidatedShipmentSignature.SignatureURL,
		CreatedAt:    consolidatedShipmentSignature.CreatedAt,
		CreatedBy:    consolidatedShipmentSignature.CreatedBy,
	}

	return
}

func (s *ConsolidatedShipmentService) Print(ctx context.Context, id int64) (res *dto.ConsolidatedShipmentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ConsolidatedShipmentService.Print")
	defer span.End()

	// check existed id
	var consolidatedShipment *model.ConsolidatedShipment
	consolidatedShipment, err = s.RepositoryConsolidatedShipment.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	//tambahan data
	var skuSummaries []*dto.SkuSummaryResponse

	var purchaseOrdersResponse []*dto.PurchaseOrderResponse

	var purchaseOrdersConsolidated []*model.PurchaseOrder
	purchaseOrdersConsolidated, _, err = s.RepositoryPurchaseOrder.GetByConsolidatedShipmentID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, purchaseOrder := range purchaseOrdersConsolidated {
		var purchaseOrderItemsResponse []*dto.PurchaseOrderItemResponse

		var purchaseOrderGP *bridgeService.GetPurchaseOrderGPResponse
		purchaseOrderGP, err = s.opt.Client.BridgeServiceGrpc.GetPurchaseOrderGPDetail(ctx, &bridgeService.GetPurchaseOrderGPDetailRequest{
			Id: purchaseOrder.PurchaseOrderIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "purchase_order")
			return
		}

		purchaseOrderItems := purchaseOrderGP.Data[0].PoDetail

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

			purchaseOrderItemsResponse = append(purchaseOrderItemsResponse, &dto.PurchaseOrderItemResponse{
				ID:              poi.Ponumber,
				PurchaseOrderID: poi.Ponumber,
				// PurchasePlanItemID: poi.PurchasePlanItemID,
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
			})
			exist := From(skuSummaries).WhereT(
				func(f *dto.SkuSummaryResponse) bool {
					return (f.ItemID == item.Data[0].Itemnmbr)
				},
			).Count()
			if exist == 0 {
				temppo := &dto.PurchaseOrderSKUSummary{
					PurchaseOrderCode: poi.Ponumber,
					Qty:               poi.Qtyorder,
				}
				tempSKUSummary := &dto.SkuSummaryResponse{
					ItemID:   item.Data[0].Itemnmbr,
					ItemName: item.Data[0].Itemdesc,
					UomName:  uom.Data[0].Umschdsc,
					TotalQty: poi.Qtyorder,
				}
				tempSKUSummary.PurchaseOrders = append(tempSKUSummary.PurchaseOrders, temppo)
				skuSummaries = append(skuSummaries, tempSKUSummary)
			} else {
				for i := range skuSummaries {
					if skuSummaries[i].ItemID == item.Data[0].Itemnmbr {
						skuSummaries[i].TotalQty = skuSummaries[i].TotalQty + poi.Qtyorder // Update the qty to the desired value
						skuSummaries[i].PurchaseOrders = append(skuSummaries[i].PurchaseOrders, &dto.PurchaseOrderSKUSummary{
							PurchaseOrderCode: poi.Ponumber,
							Qty:               poi.Qtyorder,
						})
						break
					}
				}
			}
		}

		var recognitionDate time.Time
		recognitionDate, err = time.Parse("2006-01-02T15:04:05", purchaseOrderGP.Data[0].Reqdate)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("recognition_date")
			return
		}

		var etaDate time.Time
		etaDate, err = time.Parse("2006-01-02T15:04:05", purchaseOrderGP.Data[0].PrpEstimatedarrivalDat)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_date")
			return
		}

		var etaTime time.Time
		etaTime, err = time.Parse("2006-01-02T15:04:05", purchaseOrderGP.Data[0].PrpEstimatedarrivalTim)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("eta_time")
			return
		}
		var vendor *bridgeService.GetVendorGPResponse

		vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
			Id: purchaseOrderGP.Data[0].Vendorid,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			return
		}

		if len(vendor.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
			return
		}
		purchaseOrdersResponse = append(purchaseOrdersResponse, &dto.PurchaseOrderResponse{
			ID:              purchaseOrderGP.Data[0].Ponumber,
			Code:            purchaseOrderGP.Data[0].Ponumber,
			RecognitionDate: recognitionDate,
			EtaDate:         etaDate,
			SiteAddress:     purchaseOrderGP.Data[0].PrpLocncode,
			EtaTime:         etaTime.Format("15:04"),
			TaxPct:          purchaseOrderGP.Data[0].Obtaxamt,
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
			// DeliveryFee:     purchaseOrder.DeliveryFee,
			// TotalPrice:      purchaseOrder.TotalPrice,
			// TaxAmount:       purchaseOrder.TaxAmount,
			// TotalCharge:        purchaseOrder.TotalCharge,
			// TotalInvoice:       purchaseOrder.TotalInvoice,
			// TotalWeight:        purchaseOrder.TotalWeight,
			// Note:               purchaseOrder.Note,
			// DeltaPrint:         purchaseOrder.DeltaPrint,
			// Latitude:           purchaseOrder.Latitude,
			// Longitude:          purchaseOrder.Longitude,
			PurchaseOrderItems: purchaseOrderItemsResponse,
		})

	}

	var consolidatedShipmentSignatures []*model.ConsolidatedShipmentSignature
	consolidatedShipmentSignatures, _, _ = s.RepositoryConsolidatedShipmentSignature.GetByConsolidatedShipmentID(ctx, consolidatedShipment.ID)

	var consolidatedShipmentSignaturesResponse []*dto.ConsolidatedShipmentSignatureResponse

	for _, css := range consolidatedShipmentSignatures {
		consolidatedShipmentSignaturesResponse = append(consolidatedShipmentSignaturesResponse, &dto.ConsolidatedShipmentSignatureResponse{
			ID:           css.ID,
			JobFunction:  css.JobFunction,
			Name:         css.Name,
			SignatureURL: css.SignatureURL,
			CreatedAt:    css.CreatedAt,
			CreatedBy:    css.CreatedBy,
		})
	}

	consolidatedShipment.DeltaPrint += 1
	consolidatedShipment.Status = 2

	err = s.RepositoryConsolidatedShipment.Update(ctx, consolidatedShipment, "DeltaPrint", "Status")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("id")
		return
	}

	res = &dto.ConsolidatedShipmentResponse{
		ID:                             consolidatedShipment.ID,
		Code:                           consolidatedShipment.Code,
		DriverName:                     consolidatedShipment.DriverName,
		VehicleNumber:                  consolidatedShipment.VehicleNumber,
		DriverPhoneNumber:              consolidatedShipment.DriverPhoneNumber,
		DeltaPrint:                     consolidatedShipment.DeltaPrint,
		Status:                         consolidatedShipment.Status,
		SiteName:                       consolidatedShipment.SiteName,
		ConsolidatedShipmentSignatures: consolidatedShipmentSignaturesResponse,
		SkuSummaries:                   skuSummaries,
		CreatedAt:                      consolidatedShipment.CreatedAt,
		PurchaseOrders:                 purchaseOrdersResponse,
	}

	return
}
