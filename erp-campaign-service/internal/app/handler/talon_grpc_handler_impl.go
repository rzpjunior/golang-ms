package handler

import (
	context "context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *CampaignGrpcHandler) UpdateCustomerProfileTalon(ctx context.Context, req *pb.UpdateCustomerProfileTalonRequest) (res *pb.UpdateCustomerProfileTalonResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateCustomerProfileTalon")
	defer span.End()

	customerProfile := &dto.TalonRequestUpdateCustomerProfile{
		ProfileCode:  req.ProfileCode,
		Region:       req.Region,
		CustomerType: req.CustomerType,
		CreatedDate:  req.CreatedDate.AsTime(),
		ReferrerData: req.ReferrerData,
	}

	err = h.ServiceTalon.UpdateCustomerProfileTalon(ctx, customerProfile)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	res = &pb.UpdateCustomerProfileTalonResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *CampaignGrpcHandler) UpdateCustomerSessionTalon(ctx context.Context, req *pb.UpdateCustomerSessionTalonRequest) (res *pb.UpdateCustomerSessionTalonResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateCustomerSessionTalon")
	defer span.End()

	var itemList []*dto.SessionItemData
	for _, v := range req.ItemList {
		itemList = append(itemList, &dto.SessionItemData{
			ItemName:   v.ItemName,
			ItemCode:   v.ItemCode,
			ClassName:  v.ClassName,
			UnitPrice:  v.UnitPrice,
			OrderQty:   v.OrderQty,
			UnitWeight: v.UnitWeight,
			Attributes: v.Attributes,
		})
	}

	customerSession := &dto.TalonRequestUpdateCustomerSession{
		IntegrationCode: req.IntegrationCode,
		ProfileCode:     req.ProfileCode,
		Status:          req.Status,
		IsDry:           req.IsDry,
		Archetype:       req.Archetype,
		PriceSet:        req.PriceSet,
		ReferralCode:    req.ReferralCode,
		OrderType:       req.OrderType,
		IsUsePoint:      req.IsUsePoint,
		VouDiscAmount:   req.VouDiscAmount,
		ItemList:        itemList,
	}

	var customerSessionReturn *dto.CustomerSessionReturn
	customerSessionReturn, err = h.ServiceTalon.UpdateCustomerSessionTalon(ctx, customerSession)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var totalPoints float64
	var effects []*pb.Effect
	for _, v := range customerSessionReturn.Effects {
		switch tes := v.Props.Value.(type) {
		case string:
			fmt.Println(tes)
			continue
		}
		effects = append(effects, &pb.Effect{
			CampaignId:             int32(v.CampaignID),
			EffectType:             v.EffectType,
			Name:                   v.Props.Name,
			Value:                  v.Props.Value.(float64),
			RecipientIntegrationId: v.Props.RecipientIntegrationID,
			SubledgerId:            v.Props.SubLedgerID,
		})
		if v.EffectType == "addLoyaltyPoints" && v.Props.RecipientIntegrationID == customerSession.ProfileCode {
			totalPoints += v.Props.Value.(float64)
		}
	}

	customerSessionReturnPB := &pb.CustomerSessionReturn{
		CustomerSession: &pb.CustomerSession{
			ID:               int32(customerSessionReturn.CustomerSession.ID),
			CreatedDate:      timestamppb.New(customerSessionReturn.CustomerSession.CreatedDate),
			IntegrationCode:  customerSessionReturn.CustomerSession.IntegrationCode,
			ApplicationID:    int32(customerSessionReturn.CustomerSession.ApplicationID),
			ProfileCode:      customerSessionReturn.CustomerSession.ProfileCode,
			PointEarned:      customerSessionReturn.CustomerSession.Attributes.PointEarned,
			CountGetCampaign: int32(customerSessionReturn.CustomerSession.Attributes.CountGetCampaign),
			TotalCharge:      customerSessionReturn.CustomerSession.TotalCharge,
			Subtotal:         customerSessionReturn.CustomerSession.Subtotal,
			AdditionalFee:    customerSessionReturn.CustomerSession.AdditionalFee,
		},
		CustomerProfile: &pb.Profile{
			ID:                int32(customerSessionReturn.CustomerProfile.ID),
			CreatedDate:       timestamppb.New(customerSessionReturn.CustomerProfile.CreatedDate),
			IntegrationID:     customerSessionReturn.CustomerProfile.IntegrationID,
			AccountID:         int32(customerSessionReturn.CustomerProfile.AccountID),
			ClosedSessions:    int32(customerSessionReturn.CustomerProfile.ClosedSessions),
			TotalSales:        int32(customerSessionReturn.CustomerProfile.TotalSales),
			LoyaltyMembership: customerSessionReturn.CustomerProfile.LoyaltyMembership,
			LastActivity:      timestamppb.New(customerSessionReturn.CustomerProfile.LastActivity),
		},
		Effect:      effects,
		TotalPoints: totalPoints,
	}

	res = &pb.UpdateCustomerSessionTalonResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    customerSessionReturnPB,
	}
	return
}
