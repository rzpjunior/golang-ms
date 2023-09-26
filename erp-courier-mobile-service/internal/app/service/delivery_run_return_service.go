package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/dto"
	util "git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/utils"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IDeliveryRunReturnService interface {
	// Get(ctx context.Context, req dto.DeliveryRunReturnGetRequest) (res []*dto.DeliveryRunReturnResponse, total int64, err error)
	// GetDetail(ctx context.Context, id int64, code string, deliveryRunSheetItemId int64) (res *dto.DeliveryRunReturnResponse, err error)
	Create(ctx context.Context, req dto.DeliveryReturnRequest) (err error)
	Update(ctx context.Context, req dto.DeliveryReturnRequest) (err error)
	Delete(ctx context.Context, req dto.DeleteDeliveryReturnRequest) (err error)
}

type DeliveryRunReturnService struct {
	opt opt.Options
}

func NewDeliveryRunReturnService() IDeliveryRunReturnService {
	return &DeliveryRunReturnService{
		opt: global.Setup.Common,
	}
}

func (s *DeliveryRunReturnService) Create(ctx context.Context, req dto.DeliveryReturnRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnService.Create")
	defer span.End()

	var (
		deliveryRunSheetItem *logisticService.GetDeliveryRunSheetItemDetailResponse
		salesOrder           *bridgeService.GetSalesOrderGPListResponse
		// deliveryOrders       *bridgeService.GetDeliveryOrderGPListResponse
		salesInvoiceList       *bridgeService.GetSalesInvoiceGPListResponse
		uom                    *bridgeService.GetUomGPResponse
		deliveryReturnDetail   *logisticService.GetDeliveryRunReturnDetailResponse
		deliveryReturn         *logisticService.DeliveryRunReturn
		deliveryReturnItemTemp *logisticService.DeliveryRunReturnItem
		deliveryReturnItem     []*logisticService.DeliveryRunReturnItem
		deliveryReturnRes      *logisticService.CreateDeliveryRunReturnResponse
	)
	mapItem := map[string]bool{}

	// validation cek delivery run return exist
	deliveryReturnDetail, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnDetail(ctx, &logisticService.GetDeliveryRunReturnDetailRequest{
		DeliveryRunSheetItemId: req.ID,
	})
	if deliveryReturnDetail != nil {
		err = edenlabs.ErrorExists("delivery run sheet item")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: req.ID,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// route has to be the courier job
	if deliveryRunSheetItem.Data.CourierId != req.CourierID {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", util.ErrorJobCourierInd())
		return
	}
	// step type has to be delivery type
	if deliveryRunSheetItem.Data.StepType != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("step_type.invalid", util.ErrorMustBeDeliveryInd())
		return
	}
	// status must be on progress
	if deliveryRunSheetItem.Data.Status != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("status.invalid", util.ErrorOnProgressInd("status"))
		return
	}

	// TODO GET DELIVERY ORDER ITEM

	// check if the delivery run return exist
	// filter = map[string]interface{}{"delivery_run_sheet_item_id": r.DeliveryRunSheetItem.ID}
	// _, countDeliveryRunReturn, err := repository.CheckDeliveryRunReturn(filter, exclude)
	// if err != nil {
	// 	o.Failure("delivery_run_sheet.invalid", util.ErrorInvalidDataInd("delivery run sheet"))
	// 	return o
	// }
	// if countDeliveryRunReturn > 0 {
	// 	o.Failure("delivery_run_sheet.invalid", util.ReturnExistInd("membuat pengembalian"))
	// 	return o
	// }

	// if deliveryOrderItem, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderItemDetail(ctx, &bridgeService.GetDeliveryOrderItemRequest{
	// 	salesOrderItemId:soi.Id
	// }); err != nil || deliveryOrderItem.Data == nil{
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "uom")
	// return
	// }

	// check if the length data of existing DOI and items returned are the same
	// if len(r.DeliveryOrder.DeliveryOrderItems) != len(r.Items) {
	// 	o.Failure("length.invalid", util.ErrorMustBeSameInd("jumlah produk pengembalian", "delivery order"))
	// 	return o
	// }

	// it use for get delivery fee or FrtAmount on TotalCharge calculation case
	salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: deliveryRunSheetItem.Data.SalesOrderId,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	if salesInvoiceList, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		Limit:    1,
		SoNumber: deliveryRunSheetItem.Data.SalesOrderId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice list")
		return
	}
	// return
	// if len(salesOrder.Data[0].Details) != len(req.Items) {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorValidation("length.invalid", util.ErrorMustBeSameInd("jumlah produk pengembalian", "delivery order"))
	// 	return
	// }

	// request for delivery run return
	deliveryReturn = &logisticService.DeliveryRunReturn{
		DeliveryRunSheetItemId: req.ID,
		CreatedAt:              timestamppb.New(time.Now()),
		TotalCharge:            salesOrder.Data[0].Frtamnt - salesOrder.Data[0].Trdisamt,
	}
	// operation for every item
	for _, item := range req.Items {
		if _, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "delivery_run_return_item",
			Attribute: "item_return_reason",
			ValueInt:  int32(item.ReturnReason),
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
			return
		}
		// check if a duplicate items
		if mapItem[item.ItemNumber] {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("sales_order_id.invalid", "item number "+item.ItemNumber+" duplikat")
			return
		}
		mapItem[item.ItemNumber] = true

		if item.ReceiveQty < 0 || item.ReceiveQty > item.DeliveryQty {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("receive_quantity.invalid", util.ErrorLessNotZeroInd("jumlah penerimaan", "jumlah pengiriman"))
			return
		}

		// this one for handling return data
		if item.ReceiveQty != item.DeliveryQty {
			req.ReturnedSomething = true

			for i, sii := range salesInvoiceList.Data[0].Details {
				if sii.Itemnmbr == item.ItemNumber {
					if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
						Id: salesInvoiceList.Data[0].Details[i].Uofm,
					}); err != nil {
						span.RecordError(err)
						s.opt.Logger.AddMessage(log.ErrorLevel, err)
						err = edenlabs.ErrorRpcNotFound("bridge", "uom")
						return
					}
					if uom.Data[0].Umdpqtys != 3 && (sii.Uofm == uom.Data[0].Uofm) {
						if math.Mod(item.ReceiveQty, 1) != 0 {
							span.RecordError(err)
							s.opt.Logger.AddMessage(log.ErrorLevel, err)
							err = edenlabs.ErrorValidation("receive_qty.invalid", util.ErrorNotAllowedForInd("angka decimal", "jumlah yang diterima"))
							return
						}

					}

					if item.ReceiveQty < sii.Quantity {
						if item.ReturnEvidence == "" {
							span.RecordError(err)
							s.opt.Logger.AddMessage(log.ErrorLevel, err)
							err = edenlabs.ErrorValidation("return_evidence.invalid", util.ErrorInputRequiredIndo("bukti pengembalian"))
							return
						}

						if item.ReturnReason == 0 {
							span.RecordError(err)
							s.opt.Logger.AddMessage(log.ErrorLevel, err)
							err = edenlabs.ErrorValidation("return_reason.invalid", util.ErrorInputRequiredIndo("alasan pengembalian"))
							return
						}
						//TODO: CHECK PRICE TIERING DISCOUNT
						// 	pricePerUnit := v.DeliveryOrderItem.SalesOrderItem.UnitPrice - v.DeliveryOrderItem.SalesOrderItem.UnitPriceDiscount
						item.Subtotal = item.ReceiveQty * sii.Unitprce
						deliveryReturnItemTemp = &logisticService.DeliveryRunReturnItem{
							Subtotal:            item.Subtotal,
							DeliveryOrderItemId: item.ItemNumber,
						}
						deliveryReturn.TotalPrice += item.Subtotal
						deliveryReturn.TotalCharge += item.Subtotal // it have to be change after discount from SI already seatled
					}

					// if v.ReceiveQty < v.DeliveryOrderItem.SalesOrderItem.OrderQty {
					// 	// calculate price per unit after price discount
					// 	pricePerUnit := v.DeliveryOrderItem.SalesOrderItem.UnitPrice - v.DeliveryOrderItem.SalesOrderItem.UnitPriceDiscount

					// 	// recalculate the price
					// 	v.Subtotal = v.ReceiveQty * pricePerUnit
					// } else {
					// 	// if receive qty is higher than SO order qty but is lower than the DOI deliver qty, will take the sales invoice item subtotal / SO subtotal
					// 	if !existSalesInvoice {
					// 		v.Subtotal = v.DeliveryOrderItem.SalesOrderItem.Subtotal
					// 	} else {
					// 		v.Subtotal = v.DeliveryOrderItem.SalesOrderItem.SalesInvoiceItem.Subtotal
					// 	}
					// }

					deliveryReturnItemTemp.ReceiveQty = item.ReceiveQty
					deliveryReturnItemTemp.ReturnReason = int32(item.ReturnReason)
					deliveryReturnItemTemp.ReturnEvidence = item.ReturnEvidence
					deliveryReturnItem = append(deliveryReturnItem, deliveryReturnItemTemp)
				}
			}

		} else { // this one for handling item not returned but store in database, so it will come it edit enpoint, when get detail drsi
			for _, sii := range salesInvoiceList.Data[0].Details {
				if sii.Itemnmbr == item.ItemNumber {
					item.Subtotal = item.ReceiveQty * sii.Unitprce
					deliveryReturnItemTemp = &logisticService.DeliveryRunReturnItem{
						Subtotal:            item.Subtotal,
						DeliveryOrderItemId: item.ItemNumber,
						ReceiveQty:          item.ReceiveQty,
						ReturnReason:        int32(item.ReturnReason),
						ReturnEvidence:      item.ReturnEvidence,
					}
					deliveryReturn.TotalPrice += item.Subtotal
					deliveryReturn.TotalCharge += item.Subtotal
					deliveryReturnItem = append(deliveryReturnItem, deliveryReturnItemTemp)
				}
			}
		}
	}
	if !req.ReturnedSomething {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("return.invalid", util.RequiredDataInd("sebuah pengembalian"))
		return
	}

	// insert into delivery run return header
	if deliveryReturnRes, err = s.opt.Client.LogisticServiceGrpc.CreateDeliveryRunReturn(ctx, &logisticService.CreateDeliveryRunReturnRequest{
		Model: deliveryReturn,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "create delivery run return")
		return
	}

	// insert into delivery run return item
	for _, v := range deliveryReturnItem {
		v.DeliveryRunReturnId = deliveryReturnRes.Data.Id
		if _, err = s.opt.Client.LogisticServiceGrpc.CreateDeliveryRunReturnItem(ctx, &logisticService.CreateDeliveryRunReturnItemRequest{
			Model: v,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "create delivery run return item")
			return
		}
	}

	return
}

func (s *DeliveryRunReturnService) Update(ctx context.Context, req dto.DeliveryReturnRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnService.Create")
	defer span.End()

	var (
		deliveryRunSheetItem *logisticService.GetDeliveryRunSheetItemDetailResponse
		salesOrder           *bridgeService.GetSalesOrderGPListResponse
		// deliveryOrders       *bridgeService.GetDeliveryOrderGPListResponse
		salesInvoiceList          *bridgeService.GetSalesInvoiceGPListResponse
		uom                       *bridgeService.GetUomGPResponse
		deliveryReturnDetail      *logisticService.GetDeliveryRunReturnDetailResponse
		deliveryReturn            *logisticService.DeliveryRunReturn
		deliveryReturnItemTemp    *logisticService.DeliveryRunReturnItem
		deliveryReturnItemTempRes *logisticService.GetDeliveryRunReturnItemDetailResponse
		deliveryReturnItem        []*logisticService.DeliveryRunReturnItem
		deliveryReturnRes         *logisticService.UpdateDeliveryRunReturnResponse
	)
	mapItem := map[string]bool{}

	// validation cek delivery run return exist
	deliveryReturnDetail, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnDetail(ctx, &logisticService.GetDeliveryRunReturnDetailRequest{
		Id: req.ID,
	})
	if err != nil {
		err = edenlabs.ErrorNotFound("delivery run return")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: deliveryReturnDetail.Data.DeliveryRunSheetItemId,
	}); err != nil && deliveryRunSheetItem == nil && deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// route has to be the courier job
	if deliveryRunSheetItem.Data.CourierId != req.CourierID {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", util.ErrorJobCourierInd())
		return
	}
	// step type has to be delivery type
	if deliveryRunSheetItem.Data.StepType != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("step_type.invalid", util.ErrorMustBeDeliveryInd())
		return
	}
	// status must be on progress
	if deliveryRunSheetItem.Data.Status != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("status.invalid", util.ErrorOnProgressInd("status"))
		return
	}

	// it use for get delivery fee or FrtAmount on TotalCharge calculation case
	salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: deliveryRunSheetItem.Data.SalesOrderId,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	if salesInvoiceList, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		Limit:    1,
		SoNumber: deliveryRunSheetItem.Data.SalesOrderId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice list")
		return
	}

	// request for delivery run return
	deliveryReturn = &logisticService.DeliveryRunReturn{
		Id:          req.ID,
		CreatedAt:   timestamppb.New(time.Now()),
		TotalCharge: salesOrder.Data[0].Frtamnt - salesOrder.Data[0].Trdisamt,
	}
	// operation for every item
	for _, item := range req.Items {
		if deliveryReturnItemTempRes, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnItemDetail(ctx, &logisticService.GetDeliveryRunReturnItemDetailRequest{
			Id: item.DeliveryRunReturnItemID,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "drri")
			return
		}
		if _, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "delivery_run_return_item",
			Attribute: "item_return_reason",
			ValueInt:  int32(item.ReturnReason),
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
			return
		}
		// check if a duplicate items
		if mapItem[deliveryReturnItemTempRes.Data.DeliveryOrderItemId] {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("sales_order_id.invalid", "item number "+deliveryReturnItemTempRes.Data.DeliveryOrderItemId+" duplikat")
			return
		}
		mapItem[deliveryReturnItemTempRes.Data.DeliveryOrderItemId] = true

		if item.ReceiveQty < 0 || item.ReceiveQty > item.DeliveryQty {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("receive_quantity.invalid", util.ErrorLessNotZeroInd("jumlah penerimaan", "jumlah pengiriman"))
			return
		}

		// this one for handling return data
		if item.ReceiveQty != item.DeliveryQty {
			req.ReturnedSomething = true

			for i, sii := range salesInvoiceList.Data[0].Details {
				if sii.Itemnmbr == deliveryReturnItemTempRes.Data.DeliveryOrderItemId {
					if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
						Id: salesInvoiceList.Data[0].Details[i].Uofm,
					}); err != nil {
						span.RecordError(err)
						s.opt.Logger.AddMessage(log.ErrorLevel, err)
						err = edenlabs.ErrorRpcNotFound("bridge", "uom")
						return
					}
					if uom.Data[0].Umdpqtys != 3 && (sii.Uofm == uom.Data[0].Uofm) {
						if math.Mod(item.ReceiveQty, 1) != 0 {
							span.RecordError(err)
							s.opt.Logger.AddMessage(log.ErrorLevel, err)
							err = edenlabs.ErrorValidation("receive_qty.invalid", util.ErrorNotAllowedForInd("angka decimal", "jumlah yang diterima"))
							return
						}

					}

					if item.ReceiveQty < sii.Quantity {
						if item.ReturnEvidence == "" {
							span.RecordError(err)
							s.opt.Logger.AddMessage(log.ErrorLevel, err)
							err = edenlabs.ErrorValidation("return_evidence.invalid", util.ErrorInputRequiredIndo("bukti pengembalian"))
							return
						}

						if item.ReturnReason == 0 {
							span.RecordError(err)
							s.opt.Logger.AddMessage(log.ErrorLevel, err)
							err = edenlabs.ErrorValidation("return_reason.invalid", util.ErrorInputRequiredIndo("alasan pengembalian"))
							return
						}
						//TODO: CHECK PRICE TIERING DISCOUNT
						item.Subtotal = item.ReceiveQty * sii.Unitprce
						deliveryReturnItemTemp = &logisticService.DeliveryRunReturnItem{
							Id:                  item.DeliveryRunReturnItemID,
							Subtotal:            item.Subtotal,
							DeliveryOrderItemId: deliveryReturnItemTempRes.Data.DeliveryOrderItemId,
						}
						deliveryReturn.TotalPrice += item.Subtotal
						deliveryReturn.TotalCharge += item.Subtotal // it have to be change after discount from SI already seatled
					}

					deliveryReturnItemTemp.ReceiveQty = item.ReceiveQty
					deliveryReturnItemTemp.ReturnReason = int32(item.ReturnReason)
					deliveryReturnItemTemp.ReturnEvidence = item.ReturnEvidence
					deliveryReturnItem = append(deliveryReturnItem, deliveryReturnItemTemp)
				}
			}

		} else { // this one for handling item not returned but store in database, so it will come it edit enpoint, when get detail drsi
			for _, sii := range salesInvoiceList.Data[0].Details {
				if sii.Itemnmbr == deliveryReturnItemTempRes.Data.DeliveryOrderItemId {
					item.Subtotal = item.ReceiveQty * sii.Unitprce
					deliveryReturnItemTemp = &logisticService.DeliveryRunReturnItem{
						Id:                  item.DeliveryRunReturnItemID,
						Subtotal:            item.Subtotal,
						DeliveryOrderItemId: deliveryReturnItemTempRes.Data.DeliveryOrderItemId,
						ReceiveQty:          item.ReceiveQty,
						ReturnReason:        int32(item.ReturnReason),
						ReturnEvidence:      item.ReturnEvidence,
					}
					deliveryReturn.TotalPrice += item.Subtotal
					deliveryReturn.TotalCharge += item.Subtotal
					deliveryReturnItem = append(deliveryReturnItem, deliveryReturnItemTemp)
				}
			}
		}
	}
	if !req.ReturnedSomething {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("return.invalid", util.RequiredDataInd("sebuah pengembalian"))
		return
	}

	// update into delivery run return header
	fmt.Println(deliveryReturnDetail, "Lv Header ------------------- ,", deliveryReturn.Id)
	if deliveryReturnRes, err = s.opt.Client.LogisticServiceGrpc.UpdateDeliveryRunReturn(ctx, &logisticService.UpdateDeliveryRunReturnRequest{
		Id:          deliveryReturn.Id,
		TotalPrice:  deliveryReturn.TotalPrice,
		TotalCharge: deliveryReturn.TotalCharge,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "create delivery run return")
		return
	}

	// update into delivery run return item
	for _, v := range deliveryReturnItem {

		fmt.Println(v, "Lv ITEM ------------------- ,", v.Id)
		v.DeliveryRunReturnId = deliveryReturnRes.Data.Id
		if _, err = s.opt.Client.LogisticServiceGrpc.UpdateDeliveryRunReturnItem(ctx, &logisticService.UpdateDeliveryRunReturnItemRequest{
			Id:             v.Id,
			ReceiveQty:     v.ReceiveQty,
			ReturnReason:   v.ReturnReason,
			ReturnEvidence: v.ReturnEvidence,
			Subtotal:       v.Subtotal,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "create delivery run return item")
			return
		}
	}

	return
}

func (s *DeliveryRunReturnService) Delete(ctx context.Context, req dto.DeleteDeliveryReturnRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryRunReturnService.Create")
	defer span.End()

	var (
		deliveryRunSheetItem      *logisticService.GetDeliveryRunSheetItemDetailResponse
		deliveryReturnDetail      *logisticService.GetDeliveryRunReturnDetailResponse
		deliveryReturnItemTempRes *logisticService.GetDeliveryRunReturnItemListResponse
	)

	// validation cek delivery run return exist
	deliveryReturnDetail, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnDetail(ctx, &logisticService.GetDeliveryRunReturnDetailRequest{
		Id: req.ID,
	})
	if err != nil {
		err = edenlabs.ErrorNotFound("delivery run return")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: deliveryReturnDetail.Data.DeliveryRunSheetItemId,
	}); err != nil && deliveryRunSheetItem == nil && deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// route has to be the courier job
	if deliveryRunSheetItem.Data.CourierId != req.CourierID {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", util.ErrorJobCourierInd())
		return
	}
	// step type has to be delivery type
	if deliveryRunSheetItem.Data.StepType != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("step_type.invalid", util.ErrorMustBeDeliveryInd())
		return
	}
	// status must be on progress
	if deliveryRunSheetItem.Data.Status != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("status.invalid", util.ErrorOnProgressInd("status"))
		return
	}

	if deliveryReturnItemTempRes, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnItemList(ctx, &logisticService.GetDeliveryRunReturnItemListRequest{
		Limit:               1000,
		DeliveryRunReturnId: []int64{deliveryReturnDetail.Data.Id},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "drri")
		return
	}

	// delete delivery run return item
	for _, v := range deliveryReturnItemTempRes.Data {
		if _, err = s.opt.Client.LogisticServiceGrpc.DeleteDeliveryRunReturnItem(ctx, &logisticService.DeleteDeliveryRunReturnItemRequest{
			Id: v.Id,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delete delivery run return item")
			return
		}
	}

	// delete delivery run return header
	if _, err = s.opt.Client.LogisticServiceGrpc.DeleteDeliveryRunReturn(ctx, &logisticService.DeleteDeliveryRunReturnRequest{
		Id: deliveryReturnDetail.Data.Id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delete delivery run return")
		return
	}

	return
}
