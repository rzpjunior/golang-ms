package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/repository"
)

type IPackingOrderPackService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID string, itemID string, packingDateFrom time.Time, packingDateTo time.Time) (res []*dto.PackingOrderItemResponse, total int64, err error)
	GetDetail(ctx context.Context, packingOrderID int64, packType float64, itemID string) (res *dto.PackingOrderItemResponse, err error)
	Update(ctx context.Context, req dto.PackingOrderItemRequestUpdate, packingOrderID int64) (res *dto.PackingOrderItemResponse, err error)
	Print(ctx context.Context, req dto.PackingOrderItemRequestPrint, packingOrderID int64) (res *dto.PackingOrderItemBarcodeResponse, err error)
	Dispose(ctx context.Context, req dto.PackingOrderItemRequestDispose, packingOrderID int64) (res *dto.PackingOrderItemResponse, err error)
}

type PackingOrderPackService struct {
	opt                        opt.Options
	RepositoryPackingOrder     repository.IPackingOrderRepository
	RepositoryPackingOrderItem repository.IPackingOrderItemRepository
}

func NewPackingOrderPackService() IPackingOrderPackService {
	return &PackingOrderPackService{
		opt:                        global.Setup.Common,
		RepositoryPackingOrder:     repository.NewPackingOrderRepository(),
		RepositoryPackingOrderItem: repository.NewPackingOrderItemRepository(),
	}
}

func (s *PackingOrderPackService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, siteID string, itemID string, packingDateFrom time.Time, packingDateTo time.Time) (res []*dto.PackingOrderItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderPackService.Get")
	defer span.End()

	var packingOrders []*model.PackingOrder
	packingOrders, _, err = s.RepositoryPackingOrder.Get(ctx, offset, limit, status, "", orderBy, siteID, packingDateFrom, packingDateTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var packingOrderIDs []int64
	for _, packingOrder := range packingOrders {
		packingOrderIDs = append(packingOrderIDs, packingOrder.ID)
	}

	var packingOrderItemBarcodes []*model.PackingOrderItemBarcode
	packingOrderItemBarcodes, _ = s.RepositoryPackingOrderItem.GetListBarcode(ctx, packingOrderIDs, itemID)

	for _, packingOrderItemBarcode := range packingOrderItemBarcodes {

		var po *model.PackingOrder
		po, err = s.RepositoryPackingOrder.GetByID(ctx, packingOrderItemBarcode.PackingOrderID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		var item *bridgeService.GetItemGPResponse
		item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: packingOrderItemBarcode.ItemIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item")
			return
		}

		if strings.Contains(po.Code, search) || strings.Contains(packingOrderItemBarcode.Code, search) || strings.Contains(item.Data[0].Itemnmbr, search) || strings.Contains(item.Data[0].Itmgedsc, search) {

			var site *bridgeService.GetSiteGPResponse
			site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
				Id: po.SiteIDGP,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "site")
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

			res = append(res, &dto.PackingOrderItemResponse{
				PackingOrder: &dto.PackingOrderResponse{
					ID:            po.ID,
					Code:          po.Code,
					Note:          po.Note,
					DeliveryDate:  po.DeliveryDate,
					Status:        po.Status,
					StatusConvert: statusx.ConvertStatusValue(po.Status),
				},
				Site: &dto.SiteResponse{
					ID:   site.Data[0].Locncode,
					Code: site.Data[0].Locncode,
					Name: site.Data[0].Locndscr,
				},
				PackType:      packingOrderItemBarcode.PackType,
				Status:        int8(packingOrderItemBarcode.Status),
				StatusConvert: statusx.ConvertStatusValue(int8(packingOrderItemBarcode.Status)),
				WeightScale:   packingOrderItemBarcode.WeightScale,
				Code:          packingOrderItemBarcode.Code,
				ItemID:        packingOrderItemBarcode.ItemIDGP,
				Item: &dto.ItemResponse{
					ID:          item.Data[0].Itemnmbr,
					Code:        item.Data[0].ItemTypeDesc,
					Name:        item.Data[0].Itemdesc,
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
			})
		}
	}

	return
}

func (s *PackingOrderPackService) GetDetail(ctx context.Context, packingOrderID int64, packType float64, itemID string) (res *dto.PackingOrderItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderPackService.GetByPackingOrderID")
	defer span.End()

	// validate packing order
	var packingOrder *model.PackingOrder
	packingOrder, err = s.RepositoryPackingOrder.GetByID(ctx, packingOrderID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var packingOrderItem *model.PackingOrderItem
	packingOrderItem, err = s.RepositoryPackingOrderItem.GetDetailPack(ctx, packingOrderID, packType, itemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var item *bridgeService.GetItemGPResponse
	item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: packingOrderItem.ItemIDGP,
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

	var itemImage *catalogService.GetItemImageDetailResponse
	itemImage, _ = s.opt.Client.CatalogServiceGrpc.GetItemImageDetail(ctx, &catalogService.GetItemImageDetailRequest{
		ItemId: packingOrderItem.ItemID,
	})

	res = &dto.PackingOrderItemResponse{
		PackingOrder: &dto.PackingOrderResponse{
			ID:            packingOrder.ID,
			Code:          packingOrder.Code,
			Note:          packingOrder.Note,
			DeliveryDate:  packingOrder.DeliveryDate,
			Status:        packingOrder.Status,
			StatusConvert: statusx.ConvertStatusValue(packingOrder.Status),
		},
		PackType:          packingOrderItem.PackType,
		ActualTotalPack:   packingOrderItem.ActualTotalPack,
		ExpectedTotalPack: packingOrderItem.ExpectedTotalPack,
		Status:            packingOrderItem.Status,
		WeightScale:       packingOrderItem.WeightScale,
		ItemID:            packingOrderItem.ItemIDGP,
		Item: &dto.ItemResponse{
			ID:          item.Data[0].Itemnmbr,
			Code:        item.Data[0].ItemTypeDesc,
			Name:        item.Data[0].Itmgedsc,
			Description: item.Data[0].Itmgedsc,
			UnitWeight:  float64(item.Data[0].Itemshwt),
			OrderMinQty: 1,
			OrderMaxQty: 100,
			Uom: &dto.UomResponse{
				ID:   uom.Data[0].Uomschdl,
				Code: uom.Data[0].Uomschdl,
				Name: uom.Data[0].Umschdsc,
			},
			ItemImage: &dto.ItemImageResponse{
				ID:        itemImage.Data.Id,
				ImageUrl:  itemImage.Data.ImageUrl,
				MainImage: int8(itemImage.Data.MainImage),
			},
		},
	}

	return
}

func (s *PackingOrderPackService) Update(ctx context.Context, req dto.PackingOrderItemRequestUpdate, packingOrderID int64) (res *dto.PackingOrderItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderPackService.Update")
	defer span.End()

	// validate packing order
	var packingOrder *model.PackingOrder
	packingOrder, err = s.RepositoryPackingOrder.GetByID(ctx, packingOrderID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var packingOrderItem *model.PackingOrderItem
	packingOrderItem, err = s.RepositoryPackingOrderItem.GetDetail(ctx, packingOrderID, req.PackType, req.ItemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate max actual
	if packingOrderItem.ActualTotalPack == packingOrderItem.ExpectedTotalPack {
		err = edenlabs.ErrorMustSame("actual_total", "expected_total")
		return
	}

	packingOrderItem.ActualTotalPack += 1

	err = s.RepositoryPackingOrderItem.Update(ctx, packingOrderItem, packingOrderID, req.PackType, req.ItemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var item *bridgeService.GetItemGPResponse
	item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: packingOrderItem.ItemIDGP,
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

	var itemImage *catalogService.GetItemImageDetailResponse
	itemImage, _ = s.opt.Client.CatalogServiceGrpc.GetItemImageDetail(ctx, &catalogService.GetItemImageDetailRequest{
		ItemId: 1,
	})

	res = &dto.PackingOrderItemResponse{
		PackingOrder: &dto.PackingOrderResponse{
			ID:            packingOrder.ID,
			Code:          packingOrder.Code,
			Note:          packingOrder.Note,
			DeliveryDate:  packingOrder.DeliveryDate,
			Status:        packingOrder.Status,
			StatusConvert: statusx.ConvertStatusValue(packingOrder.Status),
		},
		PackType:          packingOrderItem.PackType,
		ActualTotalPack:   packingOrderItem.ActualTotalPack,
		ExpectedTotalPack: packingOrderItem.ExpectedTotalPack,
		Status:            packingOrderItem.Status,
		WeightScale:       packingOrderItem.WeightScale,
		ItemID:            packingOrderItem.ItemIDGP,
		Item: &dto.ItemResponse{
			ID:          item.Data[0].Itemnmbr,
			Code:        item.Data[0].ItemTypeDesc,
			Name:        item.Data[0].Itmgedsc,
			Description: item.Data[0].Itmgedsc,
			UnitWeight:  float64(item.Data[0].Itemshwt),
			OrderMinQty: 1,
			OrderMaxQty: 100,
			Uom: &dto.UomResponse{
				ID:   uom.Data[0].Uomschdl,
				Code: uom.Data[0].Uomschdl,
				Name: uom.Data[0].Umschdsc,
			},
			ItemImage: &dto.ItemImageResponse{
				ID:        itemImage.Data.Id,
				ImageUrl:  itemImage.Data.ImageUrl,
				MainImage: int8(itemImage.Data.MainImage),
			},
		},
	}

	return
}

func (s *PackingOrderPackService) Print(ctx context.Context, req dto.PackingOrderItemRequestPrint, packingOrderID int64) (res *dto.PackingOrderItemBarcodeResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderPackService.Print")
	defer span.End()

	// validate packing order
	var packingOrder *model.PackingOrder
	packingOrder, err = s.RepositoryPackingOrder.GetByID(ctx, packingOrderID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var packingOrderItem *model.PackingOrderItem
	packingOrderItem, err = s.RepositoryPackingOrderItem.GetDetail(ctx, packingOrderID, req.PackType, req.ItemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// create packing order item barcode
	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "PK-" + time.Now().Format("20210929") + "-",
		Domain: "packing_order_item_barcode",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "code_generator")
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	var code string
	var packingOrderItemBarcode *model.PackingOrderItemBarcode

	if req.TypePrint == 1 {
		code = codeGenerator.Data.Code
		packingOrderItemBarcode = &model.PackingOrderItemBarcode{
			Code:           code,
			PackingOrderID: packingOrder.ID,
			ItemID:         packingOrderItem.ItemID,
			ItemIDGP:       packingOrderItem.ItemIDGP,
			PackType:       packingOrderItem.PackType,
			WeightScale:    req.WeightScale,
			Status:         int(statusx.ConvertStatusName(statusx.Active)),
			DeltaPrint:     1,
			CreatedAt:      time.Now(),
			CreatedBy:      userID,
		}
		err = s.RepositoryPackingOrderItem.CreateBarcode(ctx, packingOrderItemBarcode)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else {
		packingOrderItemBarcode, err = s.RepositoryPackingOrderItem.GetBarcode(ctx, packingOrderID, req.PackType, req.ItemID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		code = packingOrderItemBarcode.Code

		packingOrderItemBarcode.DeltaPrint += 1
		packingOrderItemBarcode.CreatedAt = time.Now()
		packingOrderItemBarcode.CreatedBy = userID

		err = s.RepositoryPackingOrderItem.UpdateBarcode(ctx, packingOrderItemBarcode, packingOrderID, req.PackType, req.ItemID, code)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	var client = &http.Client{Transport: &http.Transport{
		MaxIdleConns:       s.opt.Env.GetInt("client.print_service_api.max_idle_conns"),
		IdleConnTimeout:    s.opt.Env.GetDuration("client.print_service_api.timeout"),
		DisableCompression: true,
	}}

	var item *bridgeService.GetItemGPResponse
	item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: packingOrderItem.ItemIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "item")
		return
	}

	printReq := &dto.PrintPackingOrderRequest{
		Pk: dto.PackingOrderItemResponse{
			ID: packingOrderItem.ID.String(),
			PackingOrder: &dto.PackingOrderResponse{
				ID:           packingOrder.ID,
				Code:         packingOrder.Code,
				Note:         packingOrder.Note,
				DeliveryDate: packingOrder.DeliveryDate,
				Status:       packingOrder.Status,
			},
			ItemID: packingOrderItem.ItemIDGP,
			Item: &dto.ItemResponse{
				ID:          item.Data[0].Itemnmbr,
				Code:        item.Data[0].ItemTypeDesc,
				Name:        item.Data[0].Itmgedsc,
				Description: item.Data[0].Itmgedsc,
				UnitWeight:  float64(item.Data[0].Itemshwt),
				OrderMinQty: 1,
				OrderMaxQty: 100,
			},
			PackType:          packingOrderItem.PackType,
			ExpectedTotalPack: packingOrderItem.ExpectedTotalPack,
			ActualTotalPack:   packingOrderItem.ActualTotalPack,
			WeightScale:       packingOrderItemBarcode.WeightScale,
			Code:              packingOrderItemBarcode.Code,
			Status:            packingOrderItem.Status,
		},
	}

	jsonReq, _ := json.Marshal(printReq)

	host := fmt.Sprintf("%s:%d", s.opt.Env.GetString("client.print_service_api.host"), s.opt.Env.GetInt("client.print_service_api.port"))
	request, err := http.NewRequest("POST", host+"/api/read/label_packing", bytes.NewBuffer(jsonReq))
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	request.Header.Set("Content-Type", "application/json")
	response, err := client.Do(request)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	defer response.Body.Close()

	var bodyBytes []byte
	if response.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(response.Body)
	}

	var printResponse *dto.PrintResponse
	json.Unmarshal(bodyBytes, &printResponse)
	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	_, err = io.Copy(ioutil.Discard, response.Body)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("print", "error for generating doc print")
		return
	}

	res = &dto.PackingOrderItemBarcodeResponse{
		LinkPrint:         printResponse.Data,
		Code:              code,
		ActualTotalPack:   packingOrderItem.ActualTotalPack,
		ExpectedTotalPack: packingOrderItem.ExpectedTotalPack,
	}

	return
}

func (s *PackingOrderPackService) Dispose(ctx context.Context, req dto.PackingOrderItemRequestDispose, packingOrderID int64) (res *dto.PackingOrderItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PackingOrderPackService.Dispose")
	defer span.End()

	// validate packing order
	var packingOrder *model.PackingOrder
	packingOrder, err = s.RepositoryPackingOrder.GetByID(ctx, packingOrderID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var packingOrderItem *model.PackingOrderItem
	packingOrderItem, err = s.RepositoryPackingOrderItem.GetDetail(ctx, packingOrderID, req.PackType, req.ItemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate max actual
	if packingOrderItem.ActualTotalPack == 0 {
		err = edenlabs.ErrorMustSame("actual_total", "expected_total")
		return
	}

	packingOrderItem.ActualTotalPack -= 1

	err = s.RepositoryPackingOrderItem.Update(ctx, packingOrderItem, packingOrderID, req.PackType, req.ItemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var packingOrderItemBarcode *model.PackingOrderItemBarcode
	packingOrderItemBarcode, err = s.RepositoryPackingOrderItem.GetBarcode(ctx, packingOrderID, req.PackType, req.ItemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	userID := ctx.Value(constants.KeyUserID).(int64)

	packingOrderItemBarcode.Status = int(statusx.ConvertStatusName(statusx.Cancelled))
	packingOrderItemBarcode.DeletedAt = time.Now()
	packingOrderItemBarcode.DeletedBy = userID

	err = s.RepositoryPackingOrderItem.Update(ctx, packingOrderItem, packingOrderID, req.PackType, req.ItemID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = s.RepositoryPackingOrderItem.UpdateBarcode(ctx, packingOrderItemBarcode, packingOrderID, req.PackType, req.ItemID, packingOrderItemBarcode.Code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var item *bridgeService.GetItemGPResponse
	item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
		Id: packingOrderItem.ItemIDGP,
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

	var itemImage *catalogService.GetItemImageDetailResponse
	itemImage, _ = s.opt.Client.CatalogServiceGrpc.GetItemImageDetail(ctx, &catalogService.GetItemImageDetailRequest{
		ItemId: 1,
	})

	res = &dto.PackingOrderItemResponse{
		PackingOrder: &dto.PackingOrderResponse{
			ID:            packingOrder.ID,
			Code:          packingOrder.Code,
			Note:          packingOrder.Note,
			DeliveryDate:  packingOrder.DeliveryDate,
			Status:        packingOrder.Status,
			StatusConvert: statusx.ConvertStatusValue(packingOrder.Status),
		},
		ActualTotalPack:   packingOrderItem.ActualTotalPack,
		ExpectedTotalPack: packingOrderItem.ExpectedTotalPack,
		Code:              packingOrderItemBarcode.Code,
		PackType:          packingOrderItemBarcode.PackType,
		Status:            int8(packingOrderItemBarcode.Status),
		WeightScale:       packingOrderItemBarcode.WeightScale,
		ItemID:            packingOrderItemBarcode.ItemIDGP,
		Item: &dto.ItemResponse{
			ID:          item.Data[0].Itemnmbr,
			Code:        item.Data[0].ItemTypeDesc,
			Name:        item.Data[0].Itmgedsc,
			Description: item.Data[0].Itmgedsc,
			UnitWeight:  float64(item.Data[0].Itemshwt),
			OrderMinQty: 1,
			OrderMaxQty: 100,
			Uom: &dto.UomResponse{
				ID:   uom.Data[0].Uomschdl,
				Code: uom.Data[0].Uomschdl,
				Name: uom.Data[0].Umschdsc,
			},
			ItemImage: &dto.ItemImageResponse{
				ID:        itemImage.Data.Id,
				ImageUrl:  itemImage.Data.ImageUrl,
				MainImage: int8(itemImage.Data.MainImage),
			},
		},
	}
	return
}
