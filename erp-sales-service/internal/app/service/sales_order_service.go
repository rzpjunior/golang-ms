package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/promotion_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	salesService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/internal/app/repository"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ISalesOrderService interface {
	GetListGRPC(ctx context.Context, req *salesService.GetSalesOrderListRequest) (res []dto.SalesOrderResponse, total int64, err error)
	GetDetailGRPC(ctx context.Context, req *salesService.GetSalesOrderDetailRequest) (res dto.SalesOrderResponse, err error)
	GetListItemGRPC(ctx context.Context, req *salesService.GetSalesOrderItemListRequest) (res []dto.SalesOrderItemResponse, total int64, err error)
	GetDetailItemGRPC(ctx context.Context, req *salesService.GetSalesOrderItemDetailRequest) (res dto.SalesOrderItemResponse, err error)
	CreateSalesOrder(ctx context.Context, req *salesService.CreateSalesOrderRequest) (res dto.SalesOrderResponse, err error)
	UpdateSalesOrder(ctx context.Context, req *salesService.UpdateSalesOrderRequest) (res int64, err error)
	GetListGRPCMobile(ctx context.Context, req *salesService.GetSalesOrderListRequest) (res []dto.SalesOrderResponse, total int64, err error)
	GetListSalesOrderFeedbackGRPC(ctx context.Context, req *salesService.GetSalesOrderFeedbackListRequest) (res []dto.SalesOrderFeedback, total int64, err error)
	CreateSalesOrderFeedback(ctx context.Context, req *salesService.CreateSalesOrderFeedbackRequest) (res dto.SalesOrderFeedback, err error)
	GetSalesOrderListCronJob(ctx context.Context, req *salesService.GetSalesOrderListCronjobRequest) (res []*model.SalesOrder, err error)
	UpdateSalesOrderRemindPayment(ctx context.Context, req *salesService.UpdateSalesOrderRemindPaymentRequest) (res *salesService.UpdateSalesOrderRemindPaymentResponse, err error)
	ExpiredSalesOrder(ctx context.Context, req *salesService.ExpiredSalesOrderRequest) (res *salesService.ExpiredSalesOrderResponse, err error)
	CreateSalesOrderPaid(ctx context.Context, req *salesService.CreateSalesOrderPaidRequest) (res *salesService.CreateSalesOrderPaidResponse, err error)
}

type SalesOrderService struct {
	opt                         opt.Options
	RepositorySalesOrder        repository.ISalesOrderRepository
	RepositroySalesOrderVoucher repository.ISalesOrderVoucherRepository
}

func NewSalesOrderService() ISalesOrderService {
	return &SalesOrderService{
		opt:                         global.Setup.Common,
		RepositorySalesOrder:        repository.NewSalesOrderRepository(),
		RepositroySalesOrderVoucher: repository.NewSalesOrderVoucherRepository(),
	}
}

func (s *SalesOrderService) GetListGRPC(ctx context.Context, req *salesService.GetSalesOrderListRequest) (res []dto.SalesOrderResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Get")
	defer span.End()

	var SalesOrderes []*model.SalesOrder
	SalesOrderes, total, err = s.RepositorySalesOrder.GetListGRPC(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesOrder := range SalesOrderes {
		res = append(res, dto.SalesOrderResponse{
			ID:                  salesOrder.ID,
			SalesOrderNumber:    salesOrder.SalesOrderNumber,
			SalesOrderNumberGP:  salesOrder.SalesOrderNumber,
			AddressIDGP:         salesOrder.AddressIDGP,
			CustomerIDGP:        salesOrder.CustomerIDGP,
			WrtIDGP:             salesOrder.WrtIDGP,
			TermPaymentSlsIDGP:  salesOrder.TermPaymentSlsIDGP,
			SiteIDGP:            salesOrder.SiteIDGP,
			SubDistrictIDGP:     salesOrder.SubDistrictIDGP,
			RegionIDGP:          salesOrder.RegionIDGP,
			PaymentGroupSlsID:   int32(salesOrder.PaymentGroupSlsID),
			ArchetypeIDGP:       salesOrder.ArchetypeIDGP,
			RecognitionDate:     time.Now(),
			RequestsShipDate:    salesOrder.RequestsShipDate,
			BillingAddress:      salesOrder.BillingAddress,
			ShippingAddress:     salesOrder.ShippingAddress,
			DeliveryFee:         salesOrder.DeliveryFee,
			VouDiscAmount:       salesOrder.VouDiscAmount,
			CustomerPointLogID:  salesOrder.CustomerPointLogID,
			TotalPrice:          salesOrder.TotalPrice,
			TotalCharge:         salesOrder.TotalCharge,
			TotalWeight:         salesOrder.TotalWeight,
			Note:                salesOrder.Note,
			ShippingAddressNote: salesOrder.ShippingAddressNote,
			Status:              salesOrder.Status,
			CreatedAt:           time.Now(),
			CreatedBy:           salesOrder.CreatedBy,
			EdenPointCampaignID: salesOrder.EdenPointCampaignID,
			IntegrationCode:     salesOrder.IntegrationCode,
			PaymentReminder:     salesOrder.PaymentReminder,
			CancelType:          int8(salesOrder.CancelType),
			PriceLevelIDGP:      salesOrder.PriceLevelIDGP,
			ShippingMethodIDGP:  salesOrder.ShippingMethodIDGP,
			CustomerNameGP:      salesOrder.CustomerNameGP,
		})
	}

	return
}

func (s *SalesOrderService) GetDetailGRPC(ctx context.Context, req *salesService.GetSalesOrderDetailRequest) (res dto.SalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Get")
	defer span.End()

	var salesOrder *model.SalesOrder

	salesOrder, err = s.RepositorySalesOrder.GetDetailGRPC(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesOrderResponse{
		ID:                  salesOrder.ID,
		SalesOrderNumber:    salesOrder.SalesOrderNumber,
		SalesOrderNumberGP:  salesOrder.SalesOrderNumber,
		AddressIDGP:         salesOrder.AddressIDGP,
		CustomerIDGP:        salesOrder.CustomerIDGP,
		WrtIDGP:             salesOrder.WrtIDGP,
		TermPaymentSlsIDGP:  salesOrder.TermPaymentSlsIDGP,
		SiteIDGP:            salesOrder.SiteIDGP,
		SubDistrictIDGP:     salesOrder.SubDistrictIDGP,
		RegionIDGP:          salesOrder.RegionIDGP,
		PaymentGroupSlsID:   int32(salesOrder.PaymentGroupSlsID),
		ArchetypeIDGP:       salesOrder.ArchetypeIDGP,
		RecognitionDate:     time.Now(),
		RequestsShipDate:    salesOrder.RequestsShipDate,
		BillingAddress:      salesOrder.BillingAddress,
		ShippingAddress:     salesOrder.ShippingAddress,
		DeliveryFee:         salesOrder.DeliveryFee,
		VouDiscAmount:       salesOrder.VouDiscAmount,
		CustomerPointLogID:  salesOrder.CustomerPointLogID,
		TotalPrice:          salesOrder.TotalPrice,
		TotalCharge:         salesOrder.TotalCharge,
		TotalWeight:         salesOrder.TotalWeight,
		Note:                salesOrder.Note,
		ShippingAddressNote: salesOrder.ShippingAddressNote,
		Status:              salesOrder.Status,
		CreatedAt:           time.Now(),
		CreatedBy:           salesOrder.CreatedBy,
		EdenPointCampaignID: salesOrder.EdenPointCampaignID,
		IntegrationCode:     salesOrder.IntegrationCode,
		PaymentReminder:     salesOrder.PaymentReminder,
		CancelType:          int8(salesOrder.CancelType),
		PriceLevelIDGP:      salesOrder.PriceLevelIDGP,
		ShippingMethodIDGP:  salesOrder.ShippingMethodIDGP,
		CustomerNameGP:      salesOrder.CustomerNameGP,
	}

	// Get item list of SO
	var SalesOrderItems []*model.SalesOrderItem
	req2 := &sales_service.GetSalesOrderItemListRequest{
		SalesOrderId: salesOrder.ID,
	}
	SalesOrderItems, _, err = s.RepositorySalesOrder.GetListItemGRPC(ctx, req2)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range SalesOrderItems {
		itemDetail, _ := s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
			ItemIdGp: v.ItemIDGP,
		})
		res.SalesOrderItem = append(res.SalesOrderItem, &dto.SalesOrderItemResponse{
			ID:               v.ID,
			SalesOrderID:     v.SalesOrderID,
			ItemIDGP:         itemDetail.Data.Code,
			ItemName:         itemDetail.Data.Description,
			OrderQty:         v.OrderQty,
			UnitPrice:        v.UnitPrice,
			Subtotal:         v.Subtotal,
			Weight:           float64(v.Weight),
			UomName:          itemDetail.Data.UomName,
			UomIDGP:          v.UomIDGP,
			ImageUrl:         itemDetail.Data.ItemImage[0].ImageUrl,
			PriceTieringIDGP: v.PriceTieringIDGP,
		})
	}

	// Get item list voucher
	var SalesOrderVouchers []*model.SalesOrderVoucher
	SalesOrderVouchers, _, err = s.RepositroySalesOrderVoucher.GetList(ctx,
		&dto.GetSalesOrderVoucherListRequest{SalesOrderID: salesOrder.ID})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, v := range SalesOrderVouchers {
		res.SalesOrderVoucher = append(res.SalesOrderVoucher, &dto.SalesOrderVoucherResponse{
			ID:           v.ID,
			SalesOrderID: v.SalesOrderID,
			VoucherIDGP:  v.VoucherIDGP,
			DiscAmount:   v.DiscAmount,
			VoucherType:  v.VoucherType,
		})
	}

	return
}

func (s *SalesOrderService) GetListItemGRPC(ctx context.Context, req *salesService.GetSalesOrderItemListRequest) (res []dto.SalesOrderItemResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Get")
	defer span.End()

	var SalesOrderItems []*model.SalesOrderItem
	SalesOrderItems, total, err = s.RepositorySalesOrder.GetListItemGRPC(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesOrderItem := range SalesOrderItems {
		res = append(res, dto.SalesOrderItemResponse{
			ID:               salesOrderItem.ID,
			ItemIDGP:         salesOrderItem.ItemIDGP,
			OrderQty:         salesOrderItem.OrderQty,
			UnitPrice:        salesOrderItem.UnitPrice,
			Subtotal:         salesOrderItem.Subtotal,
			Weight:           salesOrderItem.Weight,
			PriceTieringIDGP: salesOrderItem.PriceTieringIDGP,
		})
	}

	return
}

func (s *SalesOrderService) GetDetailItemGRPC(ctx context.Context, req *salesService.GetSalesOrderItemDetailRequest) (res dto.SalesOrderItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Get")
	defer span.End()

	var salesOrderItem *model.SalesOrderItem
	salesOrderItem, err = s.RepositorySalesOrder.GetDetailItemGRPC(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesOrderItemResponse{
		ID:               salesOrderItem.ID,
		ItemIDGP:         salesOrderItem.ItemIDGP,
		OrderQty:         salesOrderItem.OrderQty,
		UnitPrice:        salesOrderItem.UnitPrice,
		Subtotal:         salesOrderItem.Subtotal,
		Weight:           salesOrderItem.Weight,
		PriceTieringIDGP: salesOrderItem.PriceTieringIDGP,
	}

	return
}

func (s *SalesOrderService) CreateSalesOrder(ctx context.Context, req *salesService.CreateSalesOrderRequest) (res dto.SalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Get")
	defer span.End()
	var codeGenerator *configuration_service.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configuration_service.GetGenerateCodeRequest{
		Format: "SO",
		Domain: "sales_order",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "generate_code")
		return
	}
	req.Data.SalesOrderNumber = codeGenerator.Data.Code
	//var salesOrder *model.SalesOrder
	so, err := s.RepositorySalesOrder.CreateSalesOrder(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	res = dto.SalesOrderResponse{
		ID:               so.ID,
		SalesOrderNumber: so.SalesOrderNumber,
		TotalCharge:      so.TotalCharge,
		RecognitionDate:  so.RequestsShipDate,
	}
	return
}

func (s *SalesOrderService) UpdateSalesOrder(ctx context.Context, req *salesService.UpdateSalesOrderRequest) (res int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.UpdateSalesOrderHeader")
	defer span.End()

	SalesOrder := &model.SalesOrder{
		ID:                  req.Data.Id,
		SalesOrderNumber:    req.Data.SalesOrderNumber,
		SalesOrderNumberGP:  req.Data.SalesOrderNumberGp,
		AddressIDGP:         req.Data.AddressIdGp,
		CustomerIDGP:        req.Data.CustomerIdGp,
		WrtIDGP:             req.Data.WrtIdGp,
		TermPaymentSlsIDGP:  req.Data.TermPaymentSlsIdGp,
		SiteIDGP:            req.Data.SiteIdGp,
		SubDistrictIDGP:     req.Data.SubDistrictIdGp,
		RegionIDGP:          req.Data.RegionIdGp,
		PaymentGroupSlsID:   int32(req.Data.PaymentGroupSlsId),
		ArchetypeIDGP:       req.Data.ArchetypeIdGp,
		BillingAddress:      req.Data.BillingAddress,
		ShippingAddress:     req.Data.ShippingAddress,
		DeliveryFee:         req.Data.DeliveryFee,
		VouDiscAmount:       req.Data.VouDiscAmount,
		CustomerPointLogID:  req.Data.CustomerPointLogId,
		TotalPrice:          req.Data.TotalPrice,
		TotalCharge:         req.Data.TotalCharge,
		TotalWeight:         req.Data.TotalWeight,
		Note:                req.Data.Note,
		ShippingAddressNote: req.Data.ShippingAddressNote,
		Status:              int8(req.Data.Status),
		EdenPointCampaignID: req.Data.EdenPointCampaignId,
		IntegrationCode:     req.Data.IntegrationCode,
		CancelType:          int8(req.Data.CancelType),
		PriceLevelIDGP:      req.Data.PriceLevelIdGp,
		ShippingMethodIDGP:  req.Data.ShippingMethodIdGp,
		CustomerNameGP:      req.Data.CustomerNameGp,
	}

	//var salesOrder *model.SalesOrder
	_, err = s.RepositorySalesOrder.UpdateSalesOrder(ctx, SalesOrder, req.FieldUpdate...)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesOrderService) GetListGRPCMobile(ctx context.Context, req *salesService.GetSalesOrderListRequest) (res []dto.SalesOrderResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Get")
	defer span.End()

	var SalesOrderes []*model.SalesOrder
	SalesOrderes, total, err = s.RepositorySalesOrder.GetListGRPC(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	tempRes := []dto.SalesOrderResponse{}
	for _, salesOrder := range SalesOrderes {
		//tembak ke gp dengan id tersebut. kalau ada di gp cek statusnya dari gp
		//logic dari gp apakah sudah posted atau belom,kalau sudah rubah status di db lokal

		tempRes = append(tempRes, dto.SalesOrderResponse{
			ID:                  salesOrder.ID,
			SalesOrderNumber:    salesOrder.SalesOrderNumber,
			SalesOrderNumberGP:  salesOrder.SalesOrderNumberGP,
			AddressIDGP:         salesOrder.AddressIDGP,
			CustomerIDGP:        salesOrder.CustomerIDGP,
			WrtIDGP:             salesOrder.WrtIDGP,
			TermPaymentSlsIDGP:  salesOrder.TermPaymentSlsIDGP,
			SiteIDGP:            salesOrder.SiteIDGP,
			SubDistrictIDGP:     salesOrder.SubDistrictIDGP,
			RegionIDGP:          salesOrder.RegionIDGP,
			PaymentGroupSlsID:   int32(salesOrder.PaymentGroupSlsID),
			ArchetypeIDGP:       salesOrder.ArchetypeIDGP,
			RecognitionDate:     time.Now(),
			RequestsShipDate:    salesOrder.RequestsShipDate,
			BillingAddress:      salesOrder.BillingAddress,
			ShippingAddress:     salesOrder.ShippingAddress,
			DeliveryFee:         salesOrder.DeliveryFee,
			VouDiscAmount:       salesOrder.VouDiscAmount,
			CustomerPointLogID:  salesOrder.CustomerPointLogID,
			TotalPrice:          salesOrder.TotalPrice,
			TotalCharge:         salesOrder.TotalCharge,
			TotalWeight:         salesOrder.TotalWeight,
			Note:                salesOrder.Note,
			ShippingAddressNote: salesOrder.ShippingAddressNote,
			Status:              salesOrder.Status,
			CreatedAt:           time.Now(),
			CreatedBy:           salesOrder.CreatedBy,
			EdenPointCampaignID: salesOrder.EdenPointCampaignID,
			IntegrationCode:     salesOrder.IntegrationCode,
			PaymentReminder:     salesOrder.PaymentReminder,
			CancelType:          int8(salesOrder.CancelType),
			PriceLevelIDGP:      salesOrder.PriceLevelIDGP,
		})
	}

	return
}

func (s *SalesOrderService) GetListSalesOrderFeedbackGRPC(ctx context.Context, req *salesService.GetSalesOrderFeedbackListRequest) (res []dto.SalesOrderFeedback, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Get")
	defer span.End()

	var salesOrderFeedbacks []*model.SalesOrderFeedback
	if req.FeedbackType == 2 {
		salesOrderFeedbacks, total, err = s.RepositorySalesOrder.GetListFeedbackGRPC(ctx, req)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else if req.FeedbackType == 1 || req.FeedbackType == 0 {
		//get from gp first
		soGP, e := s.opt.Client.BridgeServiceGrpc.GetSalesOrderList(ctx, &bridge_service.GetSalesOrderListRequest{
			CustomerId:    req.CustomerId,
			OrderDateFrom: timestamppb.New(time.Now().AddDate(0, 0, -35)),
		})
		if e != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		var tempSO []*salesService.SalesOrder
		for _, v := range soGP.Data {
			tempSO = append(tempSO, &salesService.SalesOrder{
				Id:               v.Id,
				SalesOrderNumber: v.Code,
				TotalCharge:      v.Total,
			})
		}
		salesOrderFeedbacks, total, err = s.RepositorySalesOrder.GetListUnreviewedGRPC(ctx, tempSO)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}
	for _, v := range salesOrderFeedbacks {
		res = append(res, dto.SalesOrderFeedback{
			SalesOrderCode: v.SalesOrderCode,
			DeliveryDate:   v.DeliveryDate,
			TotalCharge:    v.TotalCharge,
			RatingScore:    v.RatingScore,
			Tags:           v.Tags,
			Description:    v.Description,
			SalesOrderID:   v.SalesOrder,
		})
	}
	return
}

func (s *SalesOrderService) CreateSalesOrderFeedback(ctx context.Context, req *salesService.CreateSalesOrderFeedbackRequest) (res dto.SalesOrderFeedback, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.CreateSalesOrderFeedback")
	defer span.End()

	_, err = s.RepositorySalesOrder.CreateSalesOrderFeedback(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesOrderService) GetSalesOrderListCronJob(ctx context.Context, req *salesService.GetSalesOrderListCronjobRequest) (res []*model.SalesOrder, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.GetSalesOrderListCronJob")
	defer span.End()

	res, err = s.RepositorySalesOrder.GetSalesOrderListCronJob(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesOrderService) UpdateSalesOrderRemindPayment(ctx context.Context, req *salesService.UpdateSalesOrderRemindPaymentRequest) (res *salesService.UpdateSalesOrderRemindPaymentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.UpdateSalesOrderRemindPayment")
	defer span.End()

	//var salesOrder *model.SalesOrder
	res, err = s.RepositorySalesOrder.UpdateSalesOrderRemindPayment(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// fmt.Println("============constants.KeyUserID============", ctx.Value(constants.KeyUserID).(int64))
	// userID := ctx.Value(constants.KeyUserID).(int64)

	// _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
	// 	Log: &auditService.Log{
	// 		UserId: userID,
	// 		// ReferenceId: sai.ID,
	// 		Type:      "update_sales_order",
	// 		Function:  "UpdateSalesOrderRemindPayment",
	// 		CreatedAt: timestamppb.New(time.Now()),
	// 	},
	// })
	return
}

func (s *SalesOrderService) ExpiredSalesOrder(ctx context.Context, req *salesService.ExpiredSalesOrderRequest) (res *salesService.ExpiredSalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.ExpiredSalesOrder")
	defer span.End()

	var salesOrder *model.SalesOrder
	salesOrder, err = s.RepositorySalesOrder.GetDetailGRPC(ctx, &salesService.GetSalesOrderDetailRequest{
		Code: req.SalesOrderCode,
		// PaymentReminder: 1, //sementara di comment
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var customerInternal *crm_service.GetCustomerDetailResponse
	customerInternal, err = s.opt.Client.CrmServiceGrpc.GetCustomerDetail(ctx, &crm_service.GetCustomerDetailRequest{
		CustomerIdGp: salesOrder.CustomerIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if salesOrder != nil && salesOrder.Status == 1 {
		salesOrder.Status = statusx.ConvertStatusName(statusx.Cancelled)
		_, err = s.RepositorySalesOrder.UpdateSalesOrder(ctx, salesOrder, "status")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		var salesOrderVoucherList []*model.SalesOrderVoucher
		salesOrderVoucherList, _, err = s.RepositroySalesOrderVoucher.GetList(ctx, &dto.GetSalesOrderVoucherListRequest{
			SalesOrderID: salesOrder.ID,
		})

		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		for _, salesOrderVoucher := range salesOrderVoucherList {
			// Exclude voucher type redeem point, because this voucher not processed in internal voucher
			if salesOrderVoucher.VoucherType != 3 {
				var voucherDetail *promotion_service.GetVoucherMobileDetailResponse
				voucherDetail, err = s.opt.Client.PromotionServiceGrpc.GetVoucherMobileDetail(ctx, &promotion_service.GetVoucherMobileDetailRequest{
					Code: salesOrderVoucher.VoucherIDGP,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

				_, err = s.opt.Client.PromotionServiceGrpc.CancelVoucherLog(ctx, &promotion_service.CancelVoucherLogRequest{
					SalesOrderIdGp: req.SalesOrderCode,
					CustomerId:     customerInternal.Data.Id,
					VoucherId:      voucherDetail.Data.Id,
					AddressIdGp:    salesOrder.AddressIDGP,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

				_, err = s.opt.Client.PromotionServiceGrpc.UpdateVoucher(ctx, &promotion_service.UpdateVoucherRequest{
					VoucherId:       voucherDetail.Data.Id,
					RemOverallQuota: voucherDetail.Data.RemOverallQuota + 1,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}
			} else {
				// Handling voucher redeem point
				currentDate := time.Now().Format("2006-01-02")
				// Cancel customer point log
				_, err = s.opt.Client.CampaignServiceGrpc.CancelCustomerPointLog(ctx, &campaign_service.CancelCustomerPointLogRequest{
					CustomerId:   customerInternal.Data.Id,
					SalesOrderId: salesOrder.ID,
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

				// Create Point Log
				_, err = s.opt.Client.CampaignServiceGrpc.CreateCustomerPointLog(ctx, &campaign_service.CreateCustomerPointLogRequest{
					CustomerId:      customerInternal.Data.Id,
					SalesOrderId:    salesOrder.ID,
					PointValue:      salesOrderVoucher.DiscAmount,
					RecentPoint:     float64(customerInternal.Data.TotalPoint) + salesOrderVoucher.DiscAmount,
					CreatedDate:     currentDate,
					Status:          1,
					Note:            "Point Issued From Cancellation Redeem",
					TransactionType: 7,
					// CurrentPointUsed: req.Data.CurrentPointUsed,
					// NextPointUsed:    req.Data.NextPointUsed,
					ExpiredDate: "2023-12-31",
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

				// Update total point in customer
				_, err = s.opt.Client.CrmServiceGrpc.UpdateCustomer(ctx, &crm_service.UpdateCustomerRequest{
					Id:          customerInternal.Data.Id,
					TotalPoint:  customerInternal.Data.TotalPoint + int64(salesOrderVoucher.DiscAmount),
					FieldUpdate: []string{"total_point"},
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

				// Check existing data on customer point summary
				var customerPointSummary *campaign_service.GetCustomerPointSummaryDetailResponse
				customerPointSummary, err = s.opt.Client.CampaignServiceGrpc.GetCustomerPointSummaryDetail(ctx, &campaign_service.GetCustomerPointSummaryRequestDetail{
					CustomerId:  customerInternal.Data.Id,
					SummaryDate: currentDate,
				})
				// start create or update merchant point summary
				// if point summary is not exist, insert new point summary
				if err != nil {
					_, err = s.opt.Client.CampaignServiceGrpc.CreateCustomerPointSummary(ctx, &campaign_service.CreateCustomerPointSummaryRequest{
						CustomerId:    customerInternal.Data.Id,
						RedeemedPoint: float64(customerInternal.Data.TotalPoint) + salesOrderVoucher.DiscAmount,
						SummaryDate:   currentDate,
					})
				} else {
					// if point summary is exist, update point summary
					_, err = s.opt.Client.CampaignServiceGrpc.UpdateCustomerPointSummary(ctx, &campaign_service.UpdateCustomerPointSummaryRequest{
						Id:          customerPointSummary.Data.Id,
						EarnedPoint: customerPointSummary.Data.EarnedPoint + salesOrderVoucher.DiscAmount,
						FieldUpdate: []string{"earned_point"},
					})
				}
				// end create or update merchant point summary
			}
		}
	}

	salesOrder.Status = statusx.ConvertStatusName(statusx.Cancelled)

	_, err = s.RepositorySalesOrder.UpdateSalesOrder(ctx, salesOrder, "status")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &salesService.ExpiredSalesOrderResponse{
		CustomerIdGp: salesOrder.CustomerIDGP,
	}

	return
}

func (s *SalesOrderService) CreateSalesOrderPaid(ctx context.Context, req *salesService.CreateSalesOrderPaidRequest) (res *salesService.CreateSalesOrderPaidResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Get")
	defer span.End()

	var (
		salesOrderInternal           *model.SalesOrder
		salesOrderItems              []*model.SalesOrderItem
		salesOrderItemBodyRequest    []*bridge_service.CreateSalesOrderGPRequest_DetailItem
		salesOrderVouchers           []*model.SalesOrderVoucher
		salesOrderVoucherBodyRequest []*bridge_service.CreateSalesOrderGPRequest_VoucherApply
	)
	salesOrderInternal, err = s.RepositorySalesOrder.GetDetailGRPC(ctx, &salesService.GetSalesOrderDetailRequest{Code: req.SoCodePaidXendit})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Get Sales Order Item
	salesOrderItems, _, err = s.RepositorySalesOrder.GetListItemGRPC(ctx, &sales_service.GetSalesOrderItemListRequest{
		SalesOrderId: salesOrderInternal.ID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Set body items
	for _, item := range salesOrderItems {
		salesOrderItemBodyRequest = append(salesOrderItemBodyRequest, &bridge_service.CreateSalesOrderGPRequest_DetailItem{
			Itemnmbr:   item.ItemIDGP,
			Quantity:   item.OrderQty,
			Unitprce:   item.UnitPrice,
			Xtndprce:   item.Subtotal,
			GnL_Weight: item.Weight,
			Pricelvl:   salesOrderInternal.PriceLevelIDGP,
			Uofm:       item.UomIDGP,
			Locncode:   salesOrderInternal.SiteIDGP,
		})
	}

	// Get Sales Order Voucher
	salesOrderVouchers, _, err = s.RepositroySalesOrderVoucher.GetList(ctx, &dto.GetSalesOrderVoucherListRequest{
		SalesOrderID: salesOrderInternal.ID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Set body apply voucher
	for _, voucher := range salesOrderVouchers {
		// Exclude voucher extra eden point on body
		if voucher.VoucherType != 4 {
			salesOrderVoucherBodyRequest = append(salesOrderVoucherBodyRequest, &bridge_service.CreateSalesOrderGPRequest_VoucherApply{
				GnlVoucherType: int32(voucher.VoucherType),
				GnlVoucherId:   voucher.VoucherIDGP,
				Ordocamt:       voucher.DiscAmount,
			})
		}
	}

	layout := "2006-01-02"

	createSalesOrderGPRequest := &bridge_service.CreateSalesOrderGPRequest{
		Docid:              "SOR",
		Docdate:            salesOrderInternal.RecognitionDate.Format(layout),
		Custnmbr:           salesOrderInternal.CustomerIDGP,
		Custname:           salesOrderInternal.CustomerNameGP,
		Prstadcd:           salesOrderInternal.AddressIDGP,
		Subtotal:           salesOrderInternal.TotalPrice,
		Trdisamt:           salesOrderInternal.VouDiscAmount,
		Freight:            salesOrderInternal.DeliveryFee,
		Docamnt:            salesOrderInternal.TotalCharge,
		GnlRequestShipDate: salesOrderInternal.RequestsShipDate.Format(layout),
		GnlRegion:          salesOrderInternal.RegionIDGP,
		GnlWrtId:           salesOrderInternal.WrtIDGP,
		GnlArchetypeId:     salesOrderInternal.WrtIDGP,
		GnlOrderChannel:    "CUSTOMER_APP",
		GnlSoCodeApps:      req.SoCodePaidXendit,
		GnlTotalweight:     salesOrderInternal.TotalWeight,
		Locncode:           salesOrderInternal.SiteIDGP,
		Shipmthd:           salesOrderInternal.ShippingMethodIDGP,
		Pymtrmid:           salesOrderInternal.TermPaymentSlsIDGP,
		Note:               salesOrderInternal.Note,
		VoucherApply:       salesOrderVoucherBodyRequest,
		Detailitems:        salesOrderItemBodyRequest,
	}

	var (
		salesOrderGpResponse *bridge_service.CreateSalesOrderGPResponse
		cashReceiptGPRes     *bridge_service.CreateCashReceiptResponse
		salesOrderPayments   []*model.SalesOrderPayment
	)

	salesOrderGpResponse, err = s.opt.Client.BridgeServiceGrpc.CreateSalesOrderGP(ctx, createSalesOrderGPRequest)
	if err != nil {
		err = edenlabs.ErrorRpcNotFound("sales_order_gp", err.Error())
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	salesOrderInternal.SalesOrderNumberGP = salesOrderGpResponse.Sopnumbe
	// sementara makek status yg ini ntar tolong diganti sesuai yg seharusnya
	salesOrderInternal.Status = statusx.ConvertStatusName(statusx.InProgress)

	_, err = s.RepositorySalesOrder.UpdateSalesOrder(ctx, salesOrderInternal, "sales_order_number_gp", "status")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	cashReceiptGPRes, err = s.opt.Client.BridgeServiceGrpc.CreateCashReceipt(ctx, &bridge_service.CreateCashReceiptRequest{
		Sopnumbe:       salesOrderGpResponse.Sopnumbe,
		AmountReceived: req.Amount,
		Chekbkid:       "BCA - 0230", // Masih hardcode
		Docdate:        time.Now().Format("2006-01-02"),
	})
	if err != nil {
		err = edenlabs.ErrorRpcNotFound("cash receipt", err.Error())
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// save order payments of gp into table sales_order_payment
	salesOrderPayments = append(salesOrderPayments, &model.SalesOrderPayment{SalesOrderID: salesOrderInternal.ID, CashReceiptIdGP: cashReceiptGPRes.Docnumbr})
	_, err = s.RepositorySalesOrder.CreateSalesOrderPayment(ctx, salesOrderPayments)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &salesService.CreateSalesOrderPaidResponse{
		CustomerIdGp: salesOrderInternal.CustomerIDGP,
	}

	return
}
