package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/promotion_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IEdenPointService interface {
	IsUsedEdenPoint(ctx context.Context, req *dto.IsPointUsedRequest) (isPointUsed bool)
	GetPotentialEdenPoint(ctx context.Context, req *dto.GetPotentialEdenPointRequest) (res *dto.GetItemPotentialEdenPointResponse, err error)
	GetPointHistoryMobile(ctx context.Context, req *dto.RequestGetPointHistory) (res []*dto.PointHistoryList, err error)
	GetCustomerPointExpiration(ctx context.Context, req *dto.GetCustomerPointExpirationRequest) (res *dto.GetCustomerPointExpirationResponse, err error)
}

type EdenPointService struct {
	opt opt.Options
}

func NewEdenPointService() IEdenPointService {
	return &EdenPointService{
		opt: global.Setup.Common,
	}
}

func (s *EdenPointService) IsUsedEdenPoint(ctx context.Context, req *dto.IsPointUsedRequest) (isPointUsed bool) {
	ctx, span := s.opt.Trace.Start(ctx, "EdenPointService.IsUsedEdenPoint")
	defer span.End()

	_, err := s.opt.Client.CampaignServiceGrpc.GetCustomerPointLogDetail(ctx, &campaign_service.GetCustomerPointLogDetailRequest{
		CustomerId:  utils.ToInt64(req.Session.Customer.ID),
		Status:      int32(statusx.ConvertStatusName("Used")),
		CreatedDate: time.Now().Format("2006-01-02"),
	})

	if err != nil {
		return false
	}
	return true
}

func (s *EdenPointService) GetPotentialEdenPoint(ctx context.Context, req *dto.GetPotentialEdenPointRequest) (res *dto.GetItemPotentialEdenPointResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "EdenPointService.IsUsedEdenPoint")
	defer span.End()

	var (
		addressID, orderTypeID, itemID int
		voucherAmount                  float64
		integrationCode                string

		orderType             *bridge_service.GetOrderTypeDetailResponse
		voucher               *promotion_service.GetVoucherMobileDetailResponse
		sessionItemList       []*campaign_service.SessionItemData
		itemDetail            *bridge_service.GetItemDetailResponse
		classDetail           *bridge_service.GetClassDetailResponse
		customerSessionReturn *campaign_service.UpdateCustomerSessionTalonResponse
	)

	addressID, err = strconv.Atoi(req.Data.AddressID)
	if err != nil {
		err = edenlabs.ErrorValidation("address_id", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	_, err = s.opt.Client.BridgeServiceGrpc.GetAddressDetail(ctx, &bridge_service.GetAddressDetailRequest{
		Id: int64(addressID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("address_id", "address id tidak valid")
		return
	}

	orderTypeID, err = strconv.Atoi(req.Data.OrderTypeID)
	if err != nil {
		err = edenlabs.ErrorValidation("order_type_id", "order type id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// Set Default order type regular, in case on existing mobile didn't send value of order type
	if orderTypeID == 0 {
		orderTypeID = 1
	}

	// Validation only order type regular or self pickup can created
	if orderTypeID != 1 && orderTypeID != 6 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("order_type_id", "order type id tidak valid")
		return
	}
	orderType, err = s.opt.Client.BridgeServiceGrpc.GetOrderTypeDetail(ctx, &bridge_service.GetOrderTypeDetailRequest{
		Id: int64(orderTypeID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("address_id", "address id tidak valid")
		return
	}
	if req.Data.RedeemCode != "" {
		voucher, err = s.opt.Client.PromotionServiceGrpc.GetVoucherMobileDetail(ctx, &promotion_service.GetVoucherMobileDetailRequest{
			RedeemCode: req.Data.RedeemCode,
			CustomerId: utils.ToInt64(req.Session.Customer.ID),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("redeem_code", "Voucher tidak ditemukan")
			return
		}

		currentTime := time.Now()

		if currentTime.Before(voucher.Data.StartTime.AsTime()) {
			err = edenlabs.ErrorValidation("redeem_code", "Voucher belum dapat digunakan")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		} else if currentTime.After(voucher.Data.EndTime.AsTime()) {
			err = edenlabs.ErrorValidation("redeem_code", "Masa berlaku Voucher sudah habis")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		} else if voucher.Data.RemOverallQuota < 1 {
			err = edenlabs.ErrorValidation("redeem_code", "Voucher ini sudah habis")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if voucher.Data.Type == 1 {
			voucherAmount = voucher.Data.DiscAmount
		}
	}

	if req.Session.Customer.ProfileCode == "" {
		_, err = s.opt.Client.CampaignServiceGrpc.UpdateCustomerProfileTalon(ctx, &campaign_service.UpdateCustomerProfileTalonRequest{
			ProfileCode:  req.Session.Customer.Code,
			Region:       "Jakarta",
			CustomerType: "Dummy CustomerType 1",
			CreatedDate:  timestamppb.New(req.Session.Customer.CreatedAt),
			// ReferrerData: []string{req.Session.Customer.ReferrerCode},
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "UpdateCustomerProfileTalon")
			return
		}
		_, err = s.opt.Client.CrmServiceGrpc.UpdateCustomer(ctx, &crm_service.UpdateCustomerRequest{
			CustomerIdGp: req.Session.Customer.Code,
			ProfileCode:  req.Session.Customer.Code,
			FieldUpdate:  []string{"ProfileCode"},
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("crm", "UpdateCustomer")
			return
		}
		// Change Talon Points not yet
		// Refferer Campaign not yet
	}
	integrationCode = strings.ReplaceAll(time.Now().Format("20060102150405.99"), ".", "") + req.Session.Customer.Code

	for _, item := range req.Data.Items {
		itemID, err = strconv.Atoi(item.ID)
		if err != nil {
			err = edenlabs.ErrorValidation("item_id", "item id tidak valid")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		itemDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridge_service.GetItemDetailRequest{
			Id: int64(itemID),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "GetItemDetail")
			return
		}

		classDetail, err = s.opt.Client.BridgeServiceGrpc.GetClassDetail(ctx, &bridge_service.GetClassDetailRequest{
			Id: int64(itemDetail.Data.ClassId),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "GetClassDetail")
			return
		}

		sessionItemList = append(sessionItemList, &campaign_service.SessionItemData{
			ItemName:   itemDetail.Data.Description,
			ItemCode:   itemDetail.Data.Code,
			ClassName:  classDetail.Data.Description,
			UnitPrice:  item.Price,
			OrderQty:   1,
			UnitWeight: item.TotalWeight,
		})
	}
	customerSessionReturn, err = s.opt.Client.CampaignServiceGrpc.UpdateCustomerSessionTalon(ctx, &campaign_service.UpdateCustomerSessionTalonRequest{
		IntegrationCode: integrationCode,
		ProfileCode:     req.Session.Customer.Code,
		Status:          "closed",
		IsDry:           "true",
		Archetype:       "Dummy Archetype",
		PriceSet:        "Dummy PriceSet",
		ReferralCode:    req.Session.Customer.ReferralCode,
		OrderType:       orderType.Data.Description,
		IsUsePoint:      false,
		VouDiscAmount:   voucherAmount,
		ItemList:        sessionItemList,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("campaign", "UpdateCustomerSessionTalon")
		return
	}

	res = &dto.GetItemPotentialEdenPointResponse{
		Points: customerSessionReturn.Data.TotalPoints,
	}

	return
}

func (s *EdenPointService) GetPointHistoryMobile(ctx context.Context, req *dto.RequestGetPointHistory) (res []*dto.PointHistoryList, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.UpdateCOD")
	defer span.End()
	customerID, _ := strconv.Atoi(req.Session.Customer.ID)

	so, err := s.opt.Client.BridgeServiceGrpc.GetSalesOrderList(ctx, &bridge_service.GetSalesOrderListRequest{
		CustomerId: int64(customerID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	for _, v := range so.Data {
		eph, err := s.opt.Client.CampaignServiceGrpc.GetCustomerPointLogDetailHistoryMobile(ctx, &campaign_service.GetCustomerPointLogDetailRequest{
			CustomerId:   int64(customerID),
			SalesOrderId: v.Id,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// return nil, err
			continue
		}
		res = append(res, &dto.PointHistoryList{
			SalesOrderCode: v.Code,
			CreatedDate:    eph.Data.CreatedDate,
			PointValue:     strconv.FormatFloat(eph.Data.PointValue, 'f', 1, 64),
			StatusType:     eph.Data.StatusType,
			Status:         strconv.Itoa(int(eph.Data.Status)),
		})
	}

	return
}

func (s *EdenPointService) GetCustomerPointExpiration(ctx context.Context, req *dto.GetCustomerPointExpirationRequest) (res *dto.GetCustomerPointExpirationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.UpdateCOD")
	defer span.End()

	res = &dto.GetCustomerPointExpirationResponse{}

	customerPointExpiration, err := s.opt.Client.CampaignServiceGrpc.GetCustomerPointExpirationDetail(ctx, &campaign_service.GetCustomerPointExpirationDetailRequest{
		CustomerId: utils.ToInt64(req.Session.Customer.ID),
	})
	if err != nil {
		res.IsHavePointExpiration = false
		err = nil
		return
	}

	res.ID = customerPointExpiration.Data.Id
	res.CustomerID = customerPointExpiration.Data.CustomerId
	res.CurrentPeriodPoint = customerPointExpiration.Data.CurrentPeriodPoint
	res.NextPeriodPoint = customerPointExpiration.Data.NextPeriodPoint
	res.CurrentPeriodDate = customerPointExpiration.Data.CurrentPeriodDate.AsTime()
	res.NextPeriodDate = customerPointExpiration.Data.NextPeriodDate.AsTime()
	res.LastUpdatedAt = customerPointExpiration.Data.LastUpdatedAt.AsTime()
	res.IsHavePointExpiration = true

	return
}
