package service

import (
	"context"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/promotion_service"
)

type IVoucherService interface {
	Get(ctx context.Context, req *dto.VoucherRequestGet) (res []*dto.VoucherResponse, err error)
	GetDetail(ctx context.Context, req *dto.VoucherRequestGetDetail) (res *dto.VoucherResponse, err error)
	ApplyVoucher(ctx context.Context, req *dto.VoucherRequestApply) (res *dto.VoucherResponse, err error)
	GetVoucherItem(ctx context.Context, req *dto.VoucherRequestGetItemList) (res []*dto.VoucherGetItemResponse, err error)
}

type VoucherService struct {
	opt opt.Options
}

func NewVoucherService() IVoucherService {
	return &VoucherService{
		opt: global.Setup.Common,
	}
}

func (s *VoucherService) Get(ctx context.Context, req *dto.VoucherRequestGet) (res []*dto.VoucherResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.Get")
	defer span.End()

	var membershipLevelID, membershipCheckpointID int64
	// validate region id
	_, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridge_service.GetAdmDivisionGPDetailRequest{
		Id:   req.Data.RegionID,
		Type: "region",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("region_id")
		return
	}

	// validate archetype id
	_, err = s.opt.Client.BridgeServiceGrpc.GetArchetypeGPDetail(ctx, &bridge_service.GetArchetypeGPDetailRequest{
		Id: req.Data.ArchetypeID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("archetype_id")
		return
	}

	// validate customer type id
	_, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
		Id: req.Data.CustomerTypeID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("customer_type_id")
		return
	}

	if req.Data.MembershipLevel != "" {
		var membershipLevel *campaign_service.GetMembershipLevelDetailResponse
		membershipLevel, err = s.opt.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx, &campaign_service.GetMembershipLevelDetailRequest{
			Level: utils.ToInt64(req.Data.MembershipLevel),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "membership level")
			return
		}
		membershipLevelID = membershipLevel.Data.Id
	}

	if req.Data.MembershipCheckpoint == "" {
		membershipCheckpointID = -1
	} else {
		var membershipCheckpoint *campaign_service.GetMembershipCheckpointDetailResponse
		membershipCheckpoint, err = s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointDetail(ctx, &campaign_service.GetMembershipCheckpointDetailRequest{
			Checkpoint: utils.ToInt64(req.Data.MembershipCheckpoint),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "membership checkpoint")
			return
		}

		membershipCheckpointID = membershipCheckpoint.Data.Id
	}

	if !req.Data.IsMembershipOnly || req.Data.MembershipLevel == "" {
		membershipLevelID = utils.ToInt64(req.Session.Customer.MembershipLevelID)
		membershipCheckpointID = utils.ToInt64(req.Session.Customer.MembershipCheckpointID)
	} else if req.Data.IsMembershipOnly || (req.Data.MembershipLevel == req.Session.Customer.MembershipLevelID && req.Data.MembershipCheckpoint == "") {
		membershipCheckpointID = utils.ToInt64(req.Session.Customer.MembershipCheckpointID)
	}

	customerLevelID := utils.ToInt(req.Session.Customer.MembershipLevelID)
	customerID := utils.ToInt(req.Session.Customer.ID)
	vouchers, err := s.opt.Client.PromotionServiceGrpc.GetVoucherMobileList(ctx, &promotion_service.GetVoucherMobileListRequest{
		RegionId:               req.Data.RegionID,
		ArchetypeId:            req.Data.ArchetypeID,
		CustomerTypeId:         req.Data.CustomerTypeID,
		MembershipLevelId:      int32(membershipLevelID),
		MembershipCheckpointId: int32(membershipCheckpointID),
		IsMembershipOnly:       req.Data.IsMembershipOnly,
		CustomerLevelId:        int32(customerLevelID),
		CustomerId:             int64(customerID),
		Limit:                  req.Limit,
		Offset:                 req.Offset,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("promotion", "voucher")
		return
	}

	for _, voucher := range vouchers.Data {
		var (
			membershipLevelResponse      *dto.MembershipLevelResponse
			membershipCheckpointResponse *dto.MembershipCheckpointResponse
		)

		if voucher.MembershipLevelId != 0 {
			var membershipLevel *campaign_service.GetMembershipLevelDetailResponse
			membershipLevel, err = s.opt.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx, &campaign_service.GetMembershipLevelDetailRequest{
				Id: voucher.MembershipLevelId,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("campaign", "membership level")
				return
			}
			membershipLevelResponse = &dto.MembershipLevelResponse{
				ID:       membershipLevel.Data.Id,
				Code:     membershipLevel.Data.Code,
				Level:    int8(membershipLevel.Data.Level),
				Name:     membershipLevel.Data.Name,
				ImageUrl: membershipLevel.Data.ImageUrl,
				Status:   int8(membershipLevel.Data.Status),
			}
			var membershipCheckpoint *campaign_service.GetMembershipCheckpointDetailResponse
			membershipCheckpoint, err = s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointDetail(ctx, &campaign_service.GetMembershipCheckpointDetailRequest{
				Id: voucher.MembershipCheckpointId,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("campaign", "membership checkpoint")
				return
			}
			membershipCheckpointResponse = &dto.MembershipCheckpointResponse{
				ID:                membershipCheckpoint.Data.Id,
				Checkpoint:        int8(membershipCheckpoint.Data.Checkpoint),
				TargetAmount:      membershipCheckpoint.Data.TargetAmount,
				Status:            int8(membershipCheckpoint.Data.Status),
				MembershipLevelID: membershipCheckpoint.Data.MembershipLevelId,
			}
		}
		if voucher.ImageUrl != "" {
			res = append(res, &dto.VoucherResponse{
				ID:                   voucher.Id,
				VoucherName:          voucher.Name,
				RedeemCode:           voucher.RedeemCode,
				ImageUrl:             voucher.ImageUrl,
				MinOrder:             voucher.MinOrder,
				EndTime:              voucher.EndTime.AsTime().Format("02/01/2006"),
				RemUserQuota:         voucher.RemUserQuota,
				VoucherItem:          int64(voucher.VoucherItem),
				MembershipLevel:      membershipLevelResponse,
				MembershipCheckpoint: membershipCheckpointResponse,
			})
		}
	}

	return
}

func (s *VoucherService) GetDetail(ctx context.Context, req *dto.VoucherRequestGetDetail) (res *dto.VoucherResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.GetDetail")
	defer span.End()
	customerID, _ := strconv.Atoi(req.Session.Customer.ID)

	voucher, err := s.opt.Client.PromotionServiceGrpc.GetVoucherMobileDetail(ctx, &promotion_service.GetVoucherMobileDetailRequest{
		RedeemCode: req.Data.RedeemCode,
		CustomerId: int64(customerID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("promotion", "voucher")
		return
	}

	res = &dto.VoucherResponse{
		ID:                     voucher.Data.Id,
		VoucherName:            voucher.Data.Name,
		RedeemCode:             voucher.Data.RedeemCode,
		ImageUrl:               voucher.Data.ImageUrl,
		MinOrder:               voucher.Data.MinOrder,
		EndTime:                voucher.Data.EndTime.AsTime().Format("02/01/2006"),
		RemUserQuota:           voucher.Data.RemUserQuota,
		VoucherItem:            int64(voucher.Data.VoucherItem),
		TermCondition:          voucher.Data.TermConditions,
		MembershipLevelID:      voucher.Data.MembershipLevelId,
		MembershipCheckpointID: voucher.Data.MembershipCheckpointId,
	}

	return
}

func (s *VoucherService) ApplyVoucher(ctx context.Context, req *dto.VoucherRequestApply) (res *dto.VoucherResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.ApplyVoucher")
	defer span.End()

	var deliveryFee float64
	customerID, _ := strconv.Atoi(req.Session.Customer.ID)

	voucher, err := s.opt.Client.PromotionServiceGrpc.GetVoucherMobileDetail(ctx, &promotion_service.GetVoucherMobileDetailRequest{
		RedeemCode: req.Data.RedeemCode,
		CustomerId: int64(customerID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("redeem_code", "Voucher tidak ditemukan")
		return
	}

	// Validate region id
	_, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPDetail(ctx, &bridge_service.GetAdmDivisionGPDetailRequest{
		Id:   req.Data.RegionID,
		Type: "region",
	})
	if err != nil {
		err = edenlabs.ErrorValidation("redeem_code", "region id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if voucher.Data.RegionId != "" && req.Data.RegionID != voucher.Data.RegionId {
		err = edenlabs.ErrorValidation("redeem_code", "Voucher tidak valid untuk area")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Validate address id
	address, err := s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridge_service.GetAddressGPDetailRequest{
		Id: req.Data.AddressID,
	})
	if err != nil {
		err = edenlabs.ErrorValidation("redeem_code", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if address.Data[0].Custnmbr != req.Session.Customer.Code {
		err = edenlabs.ErrorValidation("redeem_code", "address id tidak valid")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
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

	deliveryFeeGp, err := s.opt.Client.BridgeServiceGrpc.GetDeliveryFeeDetail(ctx, &bridge_service.GetDeliveryFeeDetailRequest{
		RegionId:       utils.ToInt64(req.Data.RegionID),
		CustomerTypeId: "1",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "Delivery Fee")
		return
	}
	deliveryFee = deliveryFeeGp.Data.DeliveryFee
	if req.Data.TotalPrice >= deliveryFeeGp.Data.MinimumOrder {
		deliveryFee = 0
	}

	if voucher.Data.Type == 1 { //Type 1 total discount
		if req.Data.TotalPrice < voucher.Data.MinOrder {
			err = edenlabs.ErrorValidation("redeem_code", "Total Order harus sama atau lebih besar dari minimum order")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		} else if req.Data.TotalPrice < voucher.Data.DiscAmount {
			err = edenlabs.ErrorValidation("redeem_code", "Total Order harus sama atau lebih besar dari nominal diskon")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else if voucher.Data.Type == 2 { //Type 2 delivery discount
		if req.Data.TotalPrice < voucher.Data.MinOrder {
			err = edenlabs.ErrorValidation("redeem_code", "Total Order harus sama atau lebih besar dari minimum order")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		} else if deliveryFee < voucher.Data.DiscAmount {
			err = edenlabs.ErrorValidation("redeem_code", "Ongkos Kirim harus sama atau lebih besar dari nominal diskon")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else if voucher.Data.Type == 3 { // Type 3 extra edenpoint
		if req.Data.TotalPrice < voucher.Data.MinOrder {
			err = edenlabs.ErrorValidation("redeem_code", "Total Order harus sama atau lebih besar dari minimum order")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	if voucher.Data.CustomerTypeId != "" && req.Session.Customer.CustomerType != voucher.Data.CustomerTypeId {
		err = edenlabs.ErrorValidation("redeem_code", "Voucher tidak valid untuk customer type")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if voucher.Data.ArchetypeId != "" && address.Data[0].GnL_Archetype_ID != voucher.Data.ArchetypeId {
		err = edenlabs.ErrorValidation("redeem_code", "Voucher tidak valid untuk archetype")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if voucher.Data.RemUserQuota < 1 {
		err = edenlabs.ErrorValidation("redeem_code", "Masa berlaku Voucher sudah habis")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if voucher.Data.VoucherItem == 1 {
		if len(req.Data.VoucherItems) <= 0 {
			err = edenlabs.ErrorValidation("redeem_code", "Silahkan diisi voucher item")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		} else {
			var voucherItems *promotion_service.GetVoucherItemListResponse
			voucherItems, err = s.opt.Client.PromotionServiceGrpc.GetVoucherItemList(ctx, &promotion_service.GetVoucherItemListRequest{
				VoucherId: voucher.Data.Id,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("redeem_code", "Voucher tidak ditemukan")
				return
			}
			productList := make(map[int64]string)
			var listValidProduct []int64
			for _, voucherItem := range req.Data.VoucherItems {
				var itemID int64
				if voucherItem.ItemID == "" {
					err = edenlabs.ErrorValidation("redeem_code", "Silahkan diisi item")
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				}

				itemID = utils.ToInt64(voucherItem.ItemID)

				if _, exist := productList[itemID]; exist {
					err = edenlabs.ErrorValidation("redeem_code", "Produk duplikat. Silahkan masukkan produk lain")
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					return
				} else {
					productList[itemID] = "t"
				}

				_, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
					Id: utils.ToString(itemID),
				})
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorRpcNotFound("catalog", "item")
					return
				}
				for _, itemVoucherDB := range voucherItems.Data {
					if itemID == itemVoucherDB.ItemId {
						if voucherItem.OrderQty < itemVoucherDB.MinQtyDisc {
							err = edenlabs.ErrorValidation("redeem_code", "Sesuaikan pesanan anda dengan syarat dan ketentuan voucher")
							span.RecordError(err)
							s.opt.Logger.AddMessage(log.ErrorLevel, err)
							return
						}
						listValidProduct = append(listValidProduct, itemID) // input to listValidProduct if same value
					}
				}
			}
			if len(voucherItems.Data) != len(listValidProduct) { // check length product from DB and from Request post
				err = edenlabs.ErrorValidation("redeem_code", "Sesuaikan pesanan anda dengan syarat dan ketentuan voucher")
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		}
	}

	if voucher.Data.MembershipLevelId != 0 {

		if utils.ToInt64(req.Session.Customer.MembershipLevelID) == 0 {
			err = edenlabs.ErrorValidation("redeem_code", "Voucher hanya berlaku untuk membership customer")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		var membershipLevel *campaign_service.GetMembershipLevelDetailResponse
		membershipLevel, err = s.opt.Client.CampaignServiceGrpc.GetMembershipLevelDetail(ctx, &campaign_service.GetMembershipLevelDetailRequest{
			Id: utils.ToInt64(req.Session.Customer.MembershipLevelID),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "Membership Level")
			return
		}

		var membershipCheckpoint *campaign_service.GetMembershipCheckpointDetailResponse
		membershipCheckpoint, err = s.opt.Client.CampaignServiceGrpc.GetMembershipCheckpointDetail(ctx, &campaign_service.GetMembershipCheckpointDetailRequest{
			Id: utils.ToInt64(req.Session.Customer.MembershipCheckpointID),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("campaign", "Membership Checkpoint")
			return
		}

		if utils.ToInt64(req.Session.Customer.MembershipLevelID) < voucher.Data.MembershipLevelId {
			err = edenlabs.ErrorValidation("redeem_code", "Voucher tidak berlaku untuk "+membershipLevel.Data.Name)
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		if utils.ToInt64(req.Session.Customer.MembershipCheckpointID) < voucher.Data.MembershipCheckpointId {
			err = edenlabs.ErrorValidation("redeem_code", "Voucher tidak berlaku untuk "+membershipLevel.Data.Name+" Lapak "+utils.ToString(membershipCheckpoint.Data.Checkpoint))
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	if voucher.Data.CustomerId != 0 {
		if utils.ToInt64(req.Session.Customer.ID) != voucher.Data.CustomerId {
			err = edenlabs.ErrorValidation("redeem_code", "Voucher tidak berlaku untuk akun anda")
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	res = &dto.VoucherResponse{
		ID:         voucher.Data.Id,
		RedeemCode: voucher.Data.RedeemCode,
		DiscAmount: voucher.Data.DiscAmount,
		Type:       int8(voucher.Data.Type),
	}

	return
}

func (s *VoucherService) GetVoucherItem(ctx context.Context, req *dto.VoucherRequestGetItemList) (res []*dto.VoucherGetItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VoucherService.GetVoucherItem")
	defer span.End()
	customerID, _ := strconv.Atoi(req.Session.Customer.ID)

	voucher, err := s.opt.Client.PromotionServiceGrpc.GetVoucherMobileDetail(ctx, &promotion_service.GetVoucherMobileDetailRequest{
		RedeemCode: req.Data.RedeemCode,
		CustomerId: int64(customerID),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("promotion", "voucher")
		return
	}

	voucherItems, err := s.opt.Client.PromotionServiceGrpc.GetVoucherItemList(ctx, &promotion_service.GetVoucherItemListRequest{
		VoucherId: voucher.Data.Id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("promotion", "voucher item")
		return
	}

	for _, v := range voucherItems.Data {
		detailItem, err := s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalog_service.GetItemDetailByInternalIdRequest{
			Id: utils.ToString(v.ItemId),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item")
		}

		uomDetail, err := s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridge_service.GetUomGPDetailRequest{
			Id: detailItem.Data.UomId,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "uom")
		}

		res = append(res, &dto.VoucherGetItemResponse{
			ItemID:     detailItem.Data.Id,
			ItemName:   detailItem.Data.Description,
			MinQtyDisc: v.MinQtyDisc,
			UomName:    uomDetail.Data[0].Umschdsc,
		})
	}

	return
}
