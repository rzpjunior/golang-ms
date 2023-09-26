package service

import (
	"context"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/reportx"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/repository"
)

type IPackingOrderService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID string, deliveryDateFrom time.Time, deliveryDateTo time.Time) (res []*dto.PackingOrderResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res *dto.PackingOrderResponse, err error)
	Export(ctx context.Context, id int64) (res *dto.PackingOrderResponseExport, err error)
	Generate(ctx context.Context, req dto.PackingOrderRequestGenerate) (err error)
}

type PackingOrderService struct {
	opt                        opt.Options
	RepositoryPackingOrder     repository.IPackingOrderRepository
	RepositoryPackingOrderItem repository.IPackingOrderItemRepository
}

func NewPackingOrderService() IPackingOrderService {
	return &PackingOrderService{
		opt:                        global.Setup.Common,
		RepositoryPackingOrder:     repository.NewPackingOrderRepository(),
		RepositoryPackingOrderItem: repository.NewPackingOrderItemRepository(),
	}
}

func (s *PackingOrderService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID string, deliveryDateFrom time.Time, deliveryDateTo time.Time) (res []*dto.PackingOrderResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderService.Get")
	defer span.End()

	var packingOrders []*model.PackingOrder
	packingOrders, total, err = s.RepositoryPackingOrder.Get(ctx, offset, limit, status, search, orderBy, siteID, deliveryDateFrom, deliveryDateTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, packingOrder := range packingOrders {
		var site *bridgeService.GetSiteGPResponse
		site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
			Id: packingOrder.SiteIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "site")
			return
		}

		res = append(res, &dto.PackingOrderResponse{
			ID:            packingOrder.ID,
			Code:          packingOrder.Code,
			Note:          packingOrder.Note,
			DeliveryDate:  packingOrder.DeliveryDate,
			Status:        packingOrder.Status,
			StatusConvert: statusx.ConvertStatusValue(packingOrder.Status),
			Site: &dto.SiteResponse{
				ID:   site.Data[0].Locncode,
				Code: site.Data[0].Locncode,
				Name: site.Data[0].Locndscr,
			},
		})
	}
	return
}

func (s *PackingOrderService) GetByID(ctx context.Context, id int64) (res *dto.PackingOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderService.GetByID")
	defer span.End()

	var packingOrder *model.PackingOrder
	packingOrder, err = s.RepositoryPackingOrder.GetByID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var packingOrderItems []*model.PackingOrderItem
	packingOrderItems, _, err = s.RepositoryPackingOrderItem.GetByPackingOrderID(ctx, packingOrder.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var packingRecommendations []*dto.PackingRecommendationResponse
	// if packingOrder.Item != "" {
	for _, itemID := range packingOrderItems {
		var packingRecommendation *dto.PackingRecommendationResponse

		var item *bridgeService.GetItemGPResponse
		item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: itemID.ItemIDGP,
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

		packingRecommendation = &dto.PackingRecommendationResponse{
			PackingOrder: &dto.PackingOrderResponse{
				ID:            packingOrder.ID,
				Code:          packingOrder.Code,
				Note:          packingOrder.Note,
				DeliveryDate:  packingOrder.DeliveryDate,
				Status:        packingOrder.Status,
				StatusConvert: statusx.ConvertStatusValue(packingOrder.Status),
			},
			ItemID: itemID.ItemIDGP,
			Item: &dto.ItemResponse{
				ID:          item.Data[0].Itemnmbr,
				Code:        item.Data[0].Itemnmbr,
				Name:        item.Data[0].Itmshnam,
				Description: item.Data[0].Itemdesc,
				UnitWeight:  float64(item.Data[0].Itemshwt),
				OrderMinQty: 1,
				OrderMaxQty: 100,
				Uom: &dto.UomResponse{
					ID:   uom.Data[0].Uomschdl,
					Code: uom.Data[0].Uomschdl,
					Name: uom.Data[0].Umschdsc,
				},
			},
		}

		var itemPacksResponse []*dto.ItemPackResponse
		// for _, packingOrderItem := range packingOrderItems {
		// if packingOrderItem.ItemIDGP == itemID {
		itemPacksResponse = append(itemPacksResponse, &dto.ItemPackResponse{
			PackType:          itemID.PackType,
			ExpectedTotalPack: itemID.ExpectedTotalPack,
			ActualTotalPack:   itemID.ActualTotalPack,
		})
		// }
		// }
		packingRecommendation.ItemPack = itemPacksResponse
		packingRecommendations = append(packingRecommendations, packingRecommendation)
	}
	// }

	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: packingOrder.SiteIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	var region *bridgeService.GetAdmDivisionGPResponse
	region, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridgeService.GetAdmDivisionGPDetailRequest{
		Id:   packingOrder.RegionIDGP,
		Type: "region",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	res = &dto.PackingOrderResponse{
		ID:            packingOrder.ID,
		Code:          packingOrder.Code,
		Note:          packingOrder.Note,
		DeliveryDate:  packingOrder.DeliveryDate,
		Status:        packingOrder.Status,
		StatusConvert: statusx.ConvertStatusValue(packingOrder.Status),
		Site: &dto.SiteResponse{
			ID:   site.Data[0].Locncode,
			Code: site.Data[0].Locncode,
			Name: site.Data[0].Locndscr,
		},
		Region: &dto.RegionResponse{
			ID:   region.Data[0].Code,
			Code: region.Data[0].Code,
			Name: packingOrder.RegionIDGP,
		},
		PackingRecommendation: packingRecommendations,
	}

	return
}

func (s *PackingOrderService) Export(ctx context.Context, id int64) (res *dto.PackingOrderResponseExport, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderService.GetByID")
	defer span.End()

	var packingOrder *model.PackingOrder
	packingOrder, err = s.RepositoryPackingOrder.GetByID(ctx, id)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var packingOrderItems []*model.PackingOrderItem
	packingOrderItems, _, err = s.RepositoryPackingOrderItem.GetByPackingOrderID(ctx, packingOrder.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var ex reportx.Excelx
	header := []string{
		"Packing_Order_Code",
		"Site",
		"Packing_Date",
		"Product_Code",
		"Product_Name",
		"UOM",
	}

	// get pack size
	var glossaryPackSizes *configurationService.GetGlossaryListResponse
	glossaryPackSizes, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryList(ctx, &configurationService.GetGlossaryListRequest{
		Table:     "packing_order",
		Attribute: "pack_size",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}

	var packSizes []float64
	for _, packSize := range glossaryPackSizes.Data {
		var size float64
		size, err = strconv.ParseFloat(packSize.ValueName, 32)
		if err != nil {
			err = edenlabs.ErrorInvalid("glossary_pack_size")
		}
		packSizes = append(packSizes, size)
	}

	for _, packSize := range packSizes {
		header = append(header, fmt.Sprintf("Pack(%.2f)", packSize))
	}

	var cells []interface{}

	var packingRecommendations []*dto.PackingRecommendationResponse
	for _, itemID := range utils.StringToStringArray(packingOrder.Item) {
		var packingRecommendation *dto.PackingRecommendationResponse

		var item *bridgeService.GetItemGPResponse
		item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: itemID,
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

		packingRecommendation = &dto.PackingRecommendationResponse{
			PackingOrder: &dto.PackingOrderResponse{
				ID:            packingOrder.ID,
				Code:          packingOrder.Code,
				Note:          packingOrder.Note,
				DeliveryDate:  packingOrder.DeliveryDate,
				Status:        packingOrder.Status,
				StatusConvert: statusx.ConvertStatusValue(packingOrder.Status),
			},
			ItemID: itemID,
			Item: &dto.ItemResponse{
				ID:          item.Data[0].Itemnmbr,
				Code:        item.Data[0].Itemnmbr,
				Name:        item.Data[0].Itmshnam,
				Description: item.Data[0].Itemdesc,
				UnitWeight:  float64(item.Data[0].Itemshwt),
				OrderMinQty: 1,
				OrderMaxQty: 100,
				Uom: &dto.UomResponse{
					ID:   uom.Data[0].Uomschdl,
					Code: uom.Data[0].Uomschdl,
					Name: uom.Data[0].Umschdsc,
				},
			},
		}

		var itemPacksResponse []*dto.ItemPackResponse
		for _, packingOrderItem := range packingOrderItems {
			itemPacksResponse = append(itemPacksResponse, &dto.ItemPackResponse{
				PackType:          packingOrderItem.PackType,
				ExpectedTotalPack: packingOrderItem.ExpectedTotalPack,
				ActualTotalPack:   packingOrderItem.ActualTotalPack,
			})
		}
		packingRecommendation.ItemPack = itemPacksResponse
		packingRecommendations = append(packingRecommendations, packingRecommendation)
	}

	var site *bridgeService.GetSiteGPResponse
	site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: packingOrder.SiteIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	for _, packingRecommendation := range packingRecommendations {
		dataCells := []interface{}{
			packingOrder.Code,
			site.Data[0].Locndscr,
			packingOrder.DeliveryDate,
			packingRecommendation.Item.Code,
			packingRecommendation.Item.Description,
			packingRecommendation.Item.Uom.Name,
		}

		for _, itemPack := range packingRecommendation.ItemPack {
			dataCells = append(dataCells, itemPack.ExpectedTotalPack)
		}

		cells = append(cells, dataCells)
	}

	ex.Sheets = append(ex.Sheets, reportx.Sheet{
		WithNumbering: true,
		Name:          "Sheet1",
		Headers:       header,
		Bodys:         cells,
	})

	fileName := fmt.Sprintf("PackingOrder_%s_%s.xlsx", time.Now().Format(timex.InFormatDate), utils.GenerateRandomDoc(5))
	fileLocation, err := reportx.GenerateXlsx(fileName, ex)

	info, err := s.opt.S3x.UploadPrivateFile(ctx, s.opt.Config.S3.BucketName, fileName, fileLocation, "application/xlsx")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = fmt.Errorf("failed to upload file | %v", err)
		return
	}

	os.Remove(fileLocation)

	res = &dto.PackingOrderResponseExport{
		Url: info,
	}

	return
}

func (s *PackingOrderService) Generate(ctx context.Context, req dto.PackingOrderRequestGenerate) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderService.Generate")
	defer span.End()

	var deliveryDate time.Time
	deliveryDate, err = time.Parse(timex.InFormatDate, req.DeliveryDate)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery_date")
		return
	}

	// validate region
	// var region *bridgeService.GetRegionDetailResponse
	_, err = s.opt.Client.BridgeServiceGrpc.GetRegionDetail(ctx, &bridgeService.GetRegionDetailRequest{
		Id: utils.ToInt64(req.SiteID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "region")
		return
	}

	// validate site
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

	// get pack size
	var glossaryPackSizes *configurationService.GetGlossaryListResponse
	glossaryPackSizes, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryList(ctx, &configurationService.GetGlossaryListRequest{
		Table:     "packing_order",
		Attribute: "pack_size",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}

	var packSizes []float64
	for _, packSize := range glossaryPackSizes.Data {
		var size float64
		size, err = strconv.ParseFloat(packSize.ValueName, 32)
		if err != nil {
			err = edenlabs.ErrorInvalid("glossary_pack_size")
		}
		packSizes = append(packSizes, size)
	}
	sort.Sort(sort.Reverse(sort.Float64Slice(packSizes)))

	// get sales order
	var salesOrders *bridgeService.GetSalesOrderGPListResponse
	salesOrders, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPAll(ctx, &bridgeService.GetSalesOrderGPListRequest{
		DocdateFrom: deliveryDate.Format("2006-01-02"),
		DocdateTo:   deliveryDate.Format("2006-01-02"),
		Limit:       1000,
		Offset:      0,
		Locncode:    req.SiteID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	fmt.Println("Amount of SO fetched: ", len(salesOrders.Data))
	fmt.Println("SO data: ", salesOrders.Data)

	// check packing order existed
	var packingOrderExisted *model.PackingOrder
	packingOrderExisted, _ = s.RepositoryPackingOrder.CheckExisted(ctx, site.Data[0].Locncode, deliveryDate)

	if packingOrderExisted.ID != 0 {
		// generate packing order item
		var packItemMaps []*model.PackingOrderPack
		var wg sync.WaitGroup

		resultChan := make(chan []*model.PackingOrderPack, len(salesOrders.Data))
		for idx, salesOrder := range salesOrders.Data {
			fmt.Println(salesOrder.Sopnumbe, "SO yang diproses", idx)
			wg.Add(1)
			// time.Sleep(1 * time.Millisecond)
			go s.getSalesOrderDetails(ctx, &wg, idx, salesOrder, packingOrderExisted.ID, packSizes, resultChan)
		}

		go func() {
			wg.Wait()
			close(resultChan)
		}()

		for result := range resultChan {
			packItemMaps = append(packItemMaps, result...)
		}
		// Collect result and append into array of slice
		var itemIDs []string
		var resultData []*model.PackingOrderPack

		// Create a map to store the total ExpectedTotalPack for each PackType within each item_id_gp
		packTypeTotal := make(map[string]map[float64]float64)
		itemConfig := make(map[string]*model.ConfigItem)

		// Iterate through the data and calculate the total ExpectedTotalPack for each PackType within each item_id_gp
		for _, item := range packItemMaps {
			itemIDGP := item.ItemIDGP
			if _, ok := packTypeTotal[itemIDGP]; !ok {
				packTypeTotal[itemIDGP] = make(map[float64]float64)
			}

			for _, itemPack := range item.ItemPack {
				packTypeTotal[itemIDGP][itemPack.PackType] += itemPack.ExpectedTotalPack
			}
			itemConfig[item.ItemIDGP] = &model.ConfigItem{
				ItemName: item.ItemName,
				ItemID:   item.ItemID,
				Uom:      item.Uom,
				UomIDGP:  item.UomIDGP,
			}
		}

		// Filter and create a new result array with the desired format
		// var resultData []ItemData
		for itemIDGP, packTypeMap := range packTypeTotal {
			var itemPacks []*model.ItemPack
			for packType, total := range packTypeMap {
				itemPacks = append(itemPacks, &model.ItemPack{
					PackType:          packType,
					ExpectedTotalPack: total,
					ActualTotalPack:   0,
					Status:            statusx.ConvertStatusName(statusx.Active),
				})
			}
			resultData = append(resultData, &model.PackingOrderPack{
				PackingOrderID:     packingOrderExisted.ID,
				ItemID:             itemConfig[itemIDGP].ItemID,
				ItemIDGP:           itemIDGP,
				ItemName:           itemConfig[itemIDGP].ItemName,
				UomIDGP:            itemConfig[itemIDGP].UomIDGP,
				Uom:                itemConfig[itemIDGP].Uom,
				OrderMinQty:        1,
				WeightPack:         0,
				ProgressPercentage: 0,
				ExcessPercentage:   0,
				TotalOrderWeight:   0,
				ItemPack:           itemPacks,
			})
		}
		var packingOrderItems []*model.PackingOrderItem
		for _, packingItem := range resultData {
			itemIDs = append(itemIDs, packingItem.ItemIDGP)
			for _, itemPack := range packingItem.ItemPack {
				packingOrderItems = append(packingOrderItems, &model.PackingOrderItem{
					ID:                 packingItem.ID,
					PackingOrderID:     packingItem.PackingOrderID,
					ItemID:             packingItem.ItemID,
					ItemIDGP:           packingItem.ItemIDGP,
					ItemName:           packingItem.ItemName,
					UomID:              packingItem.UomID,
					UomIDGP:            packingItem.UomIDGP,
					Uom:                packingItem.Uom,
					OrderMinQty:        packingItem.OrderMinQty,
					WeightScale:        0,
					ProgressPercentage: packingItem.ProgressPercentage,
					ExcessPercentage:   packingItem.ExcessPercentage,
					TotalOrderWeight:   packingItem.TotalOrderWeight,
					PackType:           itemPack.PackType,
					ExpectedTotalPack:  itemPack.ExpectedTotalPack,
					ActualTotalPack:    itemPack.ActualTotalPack,
					Status:             itemPack.Status,
				})
			}
		}

		if len(itemIDs) != 0 {
			packingOrderExisted.Item = utils.ArrayStringToString(itemIDs)

			err = s.RepositoryPackingOrder.Update(ctx, packingOrderExisted, "Item")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			// get existed item packs
			var packingOrderItemsExisted []*model.PackingOrderItem
			packingOrderItemsExisted, _, err = s.RepositoryPackingOrderItem.GetByPackingOrderID(ctx, packingOrderExisted.ID)
			if err != nil {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			for _, poi := range packingOrderItems {
				// var existed bool
				for _, poiExisted := range packingOrderItemsExisted {
					if poi.ItemIDGP == poiExisted.ItemIDGP && poi.PackType == poiExisted.PackType {
						poi.ActualTotalPack = poiExisted.ActualTotalPack
						// existed = true

						// poiExisted.ExpectedTotalPack = poi.ExpectedTotalPack

						// // update existed
						// err = s.RepositoryPackingOrderItem.Update(ctx, poiExisted, packingOrderExisted.ID, poi.PackType, poi.ItemIDGP)
						// if err != nil {
						// 	span.RecordError(err)
						// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
						// 	return
						// }
					}
				}

				// if !existed {
				// 	// insert new
				// 	err = s.RepositoryPackingOrderItem.Create(ctx, poi)
				// 	if err != nil {
				// 		span.RecordError(err)
				// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
				// 		return
				// 	}
				// }

			}

			err = s.RepositoryPackingOrder.Update(ctx, packingOrderExisted, "Item")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			// delete all records in mongodb, so it can be replaced with new one
			_, err = s.RepositoryPackingOrderItem.DeleteMany(ctx, packingOrderExisted.ID, nil)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			err = s.RepositoryPackingOrderItem.CreateMany(ctx, packingOrderItems)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		}
	} else {
		// create packing order
		var codeGenerator *configurationService.GetGenerateCodeResponse
		codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
			Format: "PC-" + site.Data[0].Locncode + "-",
			Domain: "packing_order",
			Length: 6,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("configuration", "code_generator")
			return
		}
		code := codeGenerator.Data.Code

		packingOrder := &model.PackingOrder{
			Code:         code,
			SiteIDGP:     req.SiteID,
			RegionIDGP:   req.RegionID,
			DeliveryDate: deliveryDate,
			Note:         req.Note,
			Status:       statusx.ConvertStatusName(statusx.Active),
			CreatedAt:    time.Now(),
		}

		err = s.RepositoryPackingOrder.Create(ctx, packingOrder)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		// generate packing order item
		var packItemMaps []*model.PackingOrderPack
		var wg sync.WaitGroup

		resultChan := make(chan []*model.PackingOrderPack, len(salesOrders.Data))
		for idx, salesOrder := range salesOrders.Data {
			wg.Add(1)
			go s.getSalesOrderDetails(ctx, &wg, idx, salesOrder, packingOrderExisted.ID, packSizes, resultChan)
		}

		go func() {
			wg.Wait()
			close(resultChan)
		}()

		for result := range resultChan {
			packItemMaps = append(packItemMaps, result...)
		}
		// Collect result and append into array of slice
		var itemIDs []string
		var resultData []*model.PackingOrderPack

		// Create a map to store the total ExpectedTotalPack for each PackType within each item_id_gp
		packTypeTotal := make(map[string]map[float64]float64)
		itemConfig := make(map[string]*model.ConfigItem)

		// Iterate through the data and calculate the total ExpectedTotalPack for each PackType within each item_id_gp
		for _, item := range packItemMaps {
			itemIDGP := item.ItemIDGP
			if _, ok := packTypeTotal[itemIDGP]; !ok {
				packTypeTotal[itemIDGP] = make(map[float64]float64)
			}

			for _, itemPack := range item.ItemPack {
				packTypeTotal[itemIDGP][itemPack.PackType] += itemPack.ExpectedTotalPack
			}
			itemConfig[item.ItemIDGP] = &model.ConfigItem{
				ItemName: item.ItemName,
				ItemID:   item.ItemID,
				Uom:      item.Uom,
				UomIDGP:  item.UomIDGP,
			}
		}

		// Filter and create a new result array with the desired format
		// var resultData []ItemData
		for itemIDGP, packTypeMap := range packTypeTotal {
			var itemPacks []*model.ItemPack
			for packType, total := range packTypeMap {
				itemPacks = append(itemPacks, &model.ItemPack{
					PackType:          packType,
					ExpectedTotalPack: total,
					ActualTotalPack:   0,
					Status:            statusx.ConvertStatusName(statusx.Active),
				})
			}
			resultData = append(resultData, &model.PackingOrderPack{
				PackingOrderID:     packingOrder.ID,
				ItemID:             itemConfig[itemIDGP].ItemID,
				ItemIDGP:           itemIDGP,
				ItemName:           itemConfig[itemIDGP].ItemName,
				UomIDGP:            itemConfig[itemIDGP].UomIDGP,
				Uom:                itemConfig[itemIDGP].Uom,
				OrderMinQty:        1,
				WeightPack:         0,
				ProgressPercentage: 0,
				ExcessPercentage:   0,
				TotalOrderWeight:   0,
				ItemPack:           itemPacks,
			})
		}
		var packingOrderItems []*model.PackingOrderItem
		for _, packingItem := range resultData {
			itemIDs = append(itemIDs, packingItem.ItemIDGP)
			for _, itemPack := range packingItem.ItemPack {
				packingOrderItems = append(packingOrderItems, &model.PackingOrderItem{
					ID:                 packingItem.ID,
					PackingOrderID:     packingItem.PackingOrderID,
					ItemID:             packingItem.ItemID,
					ItemIDGP:           packingItem.ItemIDGP,
					ItemName:           packingItem.ItemName,
					UomID:              packingItem.UomID,
					UomIDGP:            packingItem.UomIDGP,
					Uom:                packingItem.Uom,
					OrderMinQty:        packingItem.OrderMinQty,
					WeightScale:        0,
					ProgressPercentage: packingItem.ProgressPercentage,
					ExcessPercentage:   packingItem.ExcessPercentage,
					TotalOrderWeight:   packingItem.TotalOrderWeight,
					PackType:           itemPack.PackType,
					ExpectedTotalPack:  itemPack.ExpectedTotalPack,
					ActualTotalPack:    itemPack.ActualTotalPack,
					Status:             itemPack.Status,
				})
			}
		}

		if len(itemIDs) != 0 {
			packingOrder.Item = utils.ArrayStringToString(itemIDs)

			err = s.RepositoryPackingOrder.Update(ctx, packingOrder, "Item")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			err = s.RepositoryPackingOrderItem.CreateMany(ctx, packingOrderItems)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		}
	}

	return
}

func (s *PackingOrderService) getSalesOrderDetails(ctx context.Context, wg *sync.WaitGroup, idx int, salesOrder *bridgeService.SalesOrderGP, packingOrderID int64, packSizes []float64, resultChan chan<- []*model.PackingOrderPack) {
	defer wg.Done()
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderService.Generate")
	defer span.End()

	var (
		salesOrdersDetail *bridgeService.GetSalesOrderGPListResponse
		err               error
		packItemMaps      []*model.PackingOrderPack
	)
	salesOrdersDetail, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: salesOrder.Sopnumbe,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order item")
		return
	}

	salesOrder.Details = salesOrdersDetail.Data[0].Details

	if salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Cancelled)) {
		for _, soi := range salesOrder.Details {
			// skip for non packability
			var itemInternal *catalog_service.GetItemDetailByInternalIdResponse
			itemInternal, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
				ItemIdGp: soi.Itemnmbr,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("catalog", "item")
				return
			}
			if itemInternal.Data.Packability != "1" && itemInternal.Data.Packability != "packable" {
				continue
			}

			var packItemSizes []*model.ItemPack
			itemQty := math.Round(soi.Quantity*100) / 100
			for _, size := range packSizes {
				var totalPack float64
				if itemQty >= size {
					totalPack = math.Floor(itemQty / size)
					itemQty = itemQty - (totalPack * size)
				}
				itemPack := &model.ItemPack{
					PackType:          size,
					ExpectedTotalPack: totalPack,
					ActualTotalPack:   0,
					Status:            statusx.ConvertStatusName(statusx.Active),
				}
				packItemSizes = append(packItemSizes, itemPack)
			}

			// check exists in struct
			var existsItem bool
			for i, pim := range packItemMaps {
				if pim.ItemIDGP == soi.Itemnmbr {
					existsItem = true
					for i2 := range packItemSizes {
						packItemMaps[i].ItemPack[i2].ExpectedTotalPack = packItemMaps[i].ItemPack[i2].ExpectedTotalPack + packItemSizes[i2].ExpectedTotalPack

					}
				}
			}

			if !existsItem {
				var item *bridgeService.GetItemGPResponse
				item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
					Id: soi.Itemnmbr,
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

				var itemInternal *catalog_service.GetItemDetailByInternalIdResponse
				itemInternal, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
					ItemIdGp: soi.Itemnmbr,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorRpcNotFound("catalog", "item")
					return
				}

				packItemMaps = append(packItemMaps, &model.PackingOrderPack{
					PackingOrderID:     packingOrderID,
					ItemID:             itemInternal.Data.Id,
					ItemIDGP:           soi.Itemnmbr,
					ItemName:           soi.Itemdesc,
					UomIDGP:            uom.Data[0].Uomschdl,
					Uom:                uom.Data[0].Umschdsc,
					OrderMinQty:        1,
					WeightPack:         0,
					ProgressPercentage: 0,
					ExcessPercentage:   0,
					TotalOrderWeight:   0,
					ItemPack:           packItemSizes,
				})
			}
		}
	}
	resultChan <- packItemMaps
}
