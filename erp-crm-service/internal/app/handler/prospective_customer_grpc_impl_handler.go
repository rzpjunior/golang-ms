package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	crmService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *CrmGrpcHandler) GetProspectiveCustomerList(ctx context.Context, req *crmService.GetProspectiveCustomerListRequest) (res *crmService.GetProspectiveCustomerListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetProspectiveCustomerList")
	defer span.End()

	param := &dto.ProspectiveCustomerGetRequest{
		Search:         req.Search,
		Offset:         int64(req.Offset),
		Limit:          int64(req.Limit),
		Status:         int8(req.Status),
		ArchetypeID:    req.ArchetypeId,
		CustomerID:     req.CustomerId,
		CustomerTypeID: req.CustomerTypeId,
		SalesPersonID:  req.SalespersonId,
		RequestBy:      req.RequestedBy,
		RegionID:       req.RegionId,
		OrderBy:        req.OrderBy,
	}

	var pros_cust []*dto.ProspectiveCustomerResponse
	pros_cust, _, err = h.ServicesProspectiveCustomer.Get(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*crmService.ProspectiveCustomer
	for _, prosCust := range pros_cust {
		data = append(data, &crmService.ProspectiveCustomer{
			Id:        prosCust.ID,
			Code:      prosCust.Code,
			CreatedAt: timestamppb.New(prosCust.CreatedAt),
			UpdatedAt: timestamppb.New(prosCust.UpdatedAt),
		})
	}

	res = &crmService.GetProspectiveCustomerListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CrmGrpcHandler) GetProspectiveCustomerDetail(ctx context.Context, req *crmService.GetProspectiveCustomerDetailRequest) (res *crmService.GetProspectiveCustomerDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetProspectiveCustomerDetail")
	defer span.End()

	var pros_cust *dto.ProspectiveCustomerResponse
	pros_cust, err = h.ServicesProspectiveCustomer.GetDetail(ctx, &dto.ProspectiveCustomerGetDetailRequest{ID: req.Id, Code: req.Code})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &crmService.GetProspectiveCustomerDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &crmService.ProspectiveCustomer{
			Id:   pros_cust.ID,
			Code: pros_cust.Code,
		},
	}
	return
}

func (h *CrmGrpcHandler) DeleteProspectiveCustomer(ctx context.Context, req *crmService.DeleteProspectiveCustomerRequest) (res *crmService.DeleteProspectiveCustomerResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetProspectiveCustomerDetail")
	defer span.End()

	_, err = h.ServicesProspectiveCustomer.Delete(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &crmService.DeleteProspectiveCustomerResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}

func (h *CrmGrpcHandler) CreateProspectiveCustomer(ctx context.Context, req *crmService.CreateProspectiveCustomerRequest) (res *crmService.CreateProspectiveCustomerResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetProspectiveCustomerDetail")
	defer span.End()

	_, err = h.ServicesProspectiveCustomer.Create(ctx, &dto.ProspectiveCustomerCreateRequest{})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &crmService.CreateProspectiveCustomerResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
	}
	return
}
