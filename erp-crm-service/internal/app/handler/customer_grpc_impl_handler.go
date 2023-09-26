package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	crmService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CrmGrpcHandler) GetCustomerDetail(ctx context.Context, req *crmService.GetCustomerDetailRequest) (res *crmService.GetCustomerDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetCustomerDetail")
	defer span.End()

	param := &dto.CustomerRequestGetDetail{
		ID:           req.Id,
		CustomerIDGP: req.CustomerIdGp,
		Email:        req.Email,
		ReferrerCode: req.ReferrerCode,
	}

	var customer *dto.CustomerResponseGet
	customer, err = h.ServicesCustomer.GetDetail(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &crmService.GetCustomerDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &crmService.Customer{
			Id:                      customer.ID,
			CustomerIdGp:            customer.Code,
			ProspectiveCustomerId:   customer.ProspectiveCustomerID,
			MembershipLevelId:       customer.MembershipLevelID,
			MembershipCheckpointId:  customer.MembershipCheckpointID,
			TotalPoint:              customer.TotalPoint,
			ProfileCode:             customer.ProfileCode,
			Email:                   customer.Email,
			ReferenceInfo:           customer.ReferenceInfo,
			UpgradeStatus:           int32(customer.UpgradeStatus),
			KtpPhotosUrl:            customer.KtpPhotosUrl,
			CustomerPhotosUrl:       customer.CustomerPhotosUrl,
			CustomerSelfieUrl:       customer.CustomerSelfieUrl,
			MembershipRewardId:      customer.MembershipRewardID,
			MembershipRewardAmmount: customer.MembershipRewardAmmount,
			ReferralCode:            customer.ReferralCode,
			ReferrerCode:            customer.ReferrerCode,
		},
	}
	// fmt.Print(res)
	return
}

func (h *CrmGrpcHandler) UpdateCustomer(ctx context.Context, req *crmService.UpdateCustomerRequest) (res *crmService.UpdateCustomerResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetCustomerDetail")
	defer span.End()

	updateCustomer := &dto.CustomerRequestUpdate{
		ID:                     req.Id,
		CustomerIDGP:           req.CustomerIdGp,
		ProspectiveCustomerID:  req.ProspectiveCustomerId,
		MembershipLevelID:      req.MembershipLevelId,
		MembershipCheckpointID: req.MembershipCheckpointId,
		TotalPoint:             req.TotalPoint,
		ProfileCode:            req.ProfileCode,
		ReferenceInfo:          req.ReferenceInfo,
		UpgradeStatus:          int8(req.UpgradeStatus),
		FieldUpdate:            req.FieldUpdate,
	}

	err = h.ServicesCustomer.Update(ctx, updateCustomer)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &crmService.UpdateCustomerResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *CrmGrpcHandler) CreateCustomer(ctx context.Context, req *crmService.CreateCustomerRequest) (res *crmService.CreateCustomerResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetCustomerDetail")
	defer span.End()

	customer, err := h.ServicesCustomer.Create(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &crmService.CreateCustomerResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &crmService.Customer{
			Id: customer.ID,
		},
	}
	return
}

func (h *CrmGrpcHandler) GetCustomerID(ctx context.Context, req *crmService.GetCustomerIDRequest) (res *crmService.GetCustomerIDResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetCustomerID")
	defer span.End()

	var customer []*model.Customer
	var data []*crmService.CustomerID
	customer, err = h.ServicesCustomer.GetCustomerID(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	for _, cust := range customer {
		data = append(data, &crmService.CustomerID{
			CustomerId:   cust.ID,
			CustomerIdGp: cust.CustomerIDGP,
		})
	}
	res = &crmService.GetCustomerIDResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}
