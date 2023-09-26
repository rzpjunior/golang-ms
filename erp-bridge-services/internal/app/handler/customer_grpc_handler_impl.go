package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *BridgeGrpcHandler) GetCustomerList(ctx context.Context, req *bridgeService.GetCustomerListRequest) (res *bridgeService.GetCustomerListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerList")
	defer span.End()

	var customers []dto.CustomerResponse
	customers, _, err = h.ServicesCustomer.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.CustomerTypeId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Customer
	for _, customer := range customers {
		data = append(data, &bridgeService.Customer{
			Id:                         customer.ID,
			Code:                       customer.Code,
			Name:                       customer.Name,
			Gender:                     int32(customer.Gender),
			BirthDate:                  timestamppb.New(customer.BirthDate),
			PicName:                    customer.PicName,
			PhoneNumber:                customer.PhoneNumber,
			AltPhoneNumber:             customer.AltPhoneNumber,
			Email:                      customer.Email,
			Password:                   customer.Password,
			BillingAddress:             customer.BillingAddress,
			Note:                       customer.Note,
			ReferenceInfo:              customer.ReferenceInfo,
			TagCustomer:                customer.TagCustomer,
			TagCustomerName:            customer.TagCustomerName,
			Status:                     int32(customer.Status),
			Suspended:                  int32(customer.Suspended),
			UpgradeStatus:              int32(customer.UpgradeStatus),
			CustomerGroup:              int32(customer.CustomerGroup),
			ReferralCode:               customer.ReferralCode,
			ReferrerCode:               customer.ReferrerCode,
			TotalPoint:                 customer.TotalPoint,
			CustomerTypeCreditLimit:    int32(customer.CustomerTypeCreditLimit),
			EarnedPoint:                customer.EarnedPoint,
			RedeemedPoint:              customer.RedeemedPoint,
			CustomCreditLimit:          int32(customer.CustomCreditLimit),
			CreditLimitAmount:          customer.CreditLimitAmount,
			RemainingCreditLimitAmount: customer.RemainingCreditLimitAmount,
			ProfileCode:                customer.ProfileCode,
			AverageSales:               customer.AverageSales,
			RemainingOutstanding:       customer.RemainingOutstanding,
			OverdueDebt:                customer.OverdueDebt,
			KTPPhotosUrl:               customer.KTPPhotosUrl,
			KTPPhotosUrlArr:            customer.KTPPhotosUrlArr,
			MerchantPhotosUrl:          customer.MerchantPhotosUrl,
			MerchantPhotosUrlArr:       customer.MerchantPhotosUrlArr,
			MembershipLevelID:          customer.MembershipLevelID,
			MembershipRewardID:         customer.MembershipRewardID,
			MembershipCheckpointID:     customer.MembershipCheckpointID,
			MembershipRewardAmount:     customer.MembershipRewardAmount,
			CreatedAt:                  timestamppb.New(customer.CreatedAt),
			CreatedBy:                  customer.CreatedBy,
			LastUpdatedAt:              timestamppb.New(customer.LastUpdatedAt),
			LastUpdatedBy:              customer.LastUpdatedBy,
			BirthDateString:            customer.BirthDateString,
			TermPaymentSlsId:           customer.SalesPaymentTermID,
			CustomerTypeId:             customer.CustomerTypeId,
		})
	}

	res = &bridgeService.GetCustomerListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetCustomerDetail(ctx context.Context, req *bridgeService.GetCustomerDetailRequest) (res *bridgeService.GetCustomerDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerDetail")
	defer span.End()

	var customer dto.CustomerResponse
	customer, err = h.ServicesCustomer.GetDetail(ctx, req.Id, req.Code, req.PhoneNumber)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetCustomerDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Customer{
			Id:                         customer.ID,
			Code:                       customer.Code,
			Name:                       customer.Name,
			BirthDate:                  timestamppb.New(customer.BirthDate),
			PicName:                    customer.PicName,
			PhoneNumber:                customer.PhoneNumber,
			AltPhoneNumber:             customer.AltPhoneNumber,
			Email:                      customer.Email,
			Password:                   customer.Password,
			BillingAddress:             customer.BillingAddress,
			Note:                       customer.Note,
			ReferenceInfo:              customer.ReferenceInfo,
			TagCustomer:                customer.TagCustomer,
			TagCustomerName:            customer.TagCustomerName,
			Status:                     int32(customer.Status),
			Suspended:                  int32(customer.Suspended),
			UpgradeStatus:              int32(customer.UpgradeStatus),
			CustomerGroup:              int32(customer.CustomerGroup),
			ReferralCode:               customer.ReferralCode,
			ReferrerCode:               customer.ReferrerCode,
			TotalPoint:                 customer.TotalPoint,
			CustomerTypeCreditLimit:    int32(customer.CustomerTypeCreditLimit),
			EarnedPoint:                customer.EarnedPoint,
			RedeemedPoint:              customer.RedeemedPoint,
			CustomCreditLimit:          int32(customer.CustomCreditLimit),
			CreditLimitAmount:          customer.CreditLimitAmount,
			RemainingCreditLimitAmount: customer.RemainingCreditLimitAmount,
			ProfileCode:                customer.ProfileCode,
			AverageSales:               customer.AverageSales,
			RemainingOutstanding:       customer.RemainingOutstanding,
			OverdueDebt:                customer.OverdueDebt,
			KTPPhotosUrl:               customer.KTPPhotosUrl,
			KTPPhotosUrlArr:            customer.KTPPhotosUrlArr,
			MerchantPhotosUrl:          customer.MerchantPhotosUrl,
			MerchantPhotosUrlArr:       customer.MerchantPhotosUrlArr,
			MembershipLevelID:          customer.MembershipLevelID,
			MembershipRewardID:         customer.MembershipRewardID,
			MembershipCheckpointID:     customer.MembershipCheckpointID,
			MembershipRewardAmount:     customer.MembershipRewardAmount,
			CreatedAt:                  timestamppb.New(customer.CreatedAt),
			CreatedBy:                  customer.CreatedBy,
			LastUpdatedAt:              timestamppb.New(customer.LastUpdatedAt),
			LastUpdatedBy:              customer.LastUpdatedBy,
			BirthDateString:            customer.BirthDateString,
			CustomerTypeId:             customer.CustomerTypeId,
			TermPaymentSlsId:           customer.SalesPaymentTermID,
		},
	}
	return
}

func (h *BridgeGrpcHandler) UpdateCustomer(ctx context.Context, req *bridgeService.UpdateCustomerRequest) (res *bridgeService.UpdateCustomerResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateCustomer")
	defer span.End()

	var customer dto.CustomerResponse
	customer, err = h.ServicesCustomer.UpdateCustomer(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.UpdateCustomerResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Customer{
			Id:                         customer.ID,
			Code:                       customer.Code,
			Name:                       customer.Name,
			BirthDate:                  timestamppb.New(customer.BirthDate),
			PicName:                    customer.PicName,
			PhoneNumber:                customer.PhoneNumber,
			AltPhoneNumber:             customer.AltPhoneNumber,
			Email:                      customer.Email,
			Password:                   customer.Password,
			BillingAddress:             customer.BillingAddress,
			Note:                       customer.Note,
			ReferenceInfo:              customer.ReferenceInfo,
			TagCustomer:                customer.TagCustomer,
			TagCustomerName:            customer.TagCustomerName,
			Status:                     int32(customer.Status),
			Suspended:                  int32(customer.Suspended),
			UpgradeStatus:              int32(customer.UpgradeStatus),
			CustomerGroup:              int32(customer.CustomerGroup),
			ReferralCode:               customer.ReferralCode,
			ReferrerCode:               customer.ReferrerCode,
			TotalPoint:                 customer.TotalPoint,
			CustomerTypeCreditLimit:    int32(customer.CustomerTypeCreditLimit),
			EarnedPoint:                customer.EarnedPoint,
			RedeemedPoint:              customer.RedeemedPoint,
			CustomCreditLimit:          int32(customer.CustomCreditLimit),
			CreditLimitAmount:          customer.CreditLimitAmount,
			RemainingCreditLimitAmount: customer.RemainingCreditLimitAmount,
			ProfileCode:                customer.ProfileCode,
			AverageSales:               customer.AverageSales,
			RemainingOutstanding:       customer.RemainingOutstanding,
			OverdueDebt:                customer.OverdueDebt,
			KTPPhotosUrl:               customer.KTPPhotosUrl,
			KTPPhotosUrlArr:            customer.KTPPhotosUrlArr,
			MerchantPhotosUrl:          customer.MerchantPhotosUrl,
			MerchantPhotosUrlArr:       customer.MerchantPhotosUrlArr,
			MembershipLevelID:          customer.MembershipLevelID,
			MembershipRewardID:         customer.MembershipRewardID,
			MembershipCheckpointID:     customer.MembershipCheckpointID,
			MembershipRewardAmount:     customer.MembershipRewardAmount,
			CreatedAt:                  timestamppb.New(customer.CreatedAt),
			CreatedBy:                  customer.CreatedBy,
			LastUpdatedAt:              timestamppb.New(customer.LastUpdatedAt),
			LastUpdatedBy:              customer.LastUpdatedBy,
			BirthDateString:            customer.BirthDateString,
			CustomerTypeId:             customer.CustomerTypeId,
			TermPaymentSlsId:           customer.SalesPaymentTermID,
		},
	}
	return
}

// get customer list from gp
func (h *BridgeGrpcHandler) GetCustomerGPList(ctx context.Context, req *bridgeService.GetCustomerGPListRequest) (res *bridgeService.GetCustomerGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerGPList")
	defer span.End()

	res, err = h.ServicesCustomer.GetGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) GetCustomerGPDetail(ctx context.Context, req *bridgeService.GetCustomerGPDetailRequest) (res *bridgeService.GetCustomerGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCustomerGPDetail")
	defer span.End()

	res, err = h.ServicesCustomer.GetDetailGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) CreateCustomerGP(ctx context.Context, req *bridgeService.CreateCustomerGPRequest) (res *bridgeService.CreateCustomerGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateCustomer")
	defer span.End()

	// var customer dto.CustomerResponse
	var response dto.CreateCustomerGPResponse
	response, err = h.ServicesCustomer.CreateCustomerGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.CreateCustomerGPResponse{
		Code:            int32(codes.OK),
		Message:         codes.OK.String(),
		Custnmbr:        response.CustNmbr,
		GnLReferralCode: response.GnlReferralCode,
		GnlReferrerCode: response.GnlReferrerCode,
		Shipcomplete:    utils.ToString(response.ShipComplete),
	}
	return
}

func (h *BridgeGrpcHandler) UpdateCustomerGP(ctx context.Context, req *bridgeService.UpdateCustomerGPRequest) (res *bridgeService.UpdateCustomerGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateCustomer")
	defer span.End()

	// var customer dto.CustomerResponse
	_, err = h.ServicesCustomer.UpdateCustomerGP(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.UpdateCustomerGPResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *BridgeGrpcHandler) UpdateFixedVa(ctx context.Context, req *bridgeService.UpdateFixedVaRequest) (res *bridgeService.UpdateFixedVaResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.UpdateCustomer")
	defer span.End()


	_, err = h.ServicesCustomer.UpdateFixedVa(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.UpdateFixedVaResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}
