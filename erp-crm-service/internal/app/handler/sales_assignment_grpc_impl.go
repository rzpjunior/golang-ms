package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	crmService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *CrmGrpcHandler) GetSalesAssignmentList(ctx context.Context, req *crmService.GetSalesAssignmentListRequest) (res *crmService.GetSalesAssignmentListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetSalesAssignmentList")
	defer span.End()

	// var assignments []*dto.SalesAssignmentResponse
	// assignments, _, err = h.ServicesSalesAssignment.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, req.TeritoryId, req.StartDateFrom.AsTime(), req.StartDateTo.AsTime(), req.EndDateFrom.AsTime(), req.StartDateTo.AsTime())
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// var data []*crmService.SalesAssignment
	// for _, assignment := range assignments {
	// data = append(data, &crmService.SalesAssignment{
	// 	Id:          assignment.ID,
	// 	Code:        assignment.Code,
	// 	TerritoryId: assignment.Territory.ID,
	// 	StartDate:   timestamppb.New(assignment.StartDate),
	// 	EndDate:     timestamppb.New(assignment.EndDate),
	// 	Status:      int32(assignment.Status),
	// })
	// }

	// res = &crmService.GetSalesAssignmentListResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data:    data,
	// }
	return
}

func (h *CrmGrpcHandler) GetSalesAssignmentDetail(ctx context.Context, req *crmService.GetSalesAssignmentDetailRequest) (res *crmService.GetSalesAssignmentDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetSalesAssignmentDetail")
	defer span.End()

	var assignment dto.SalesAssignmentResponse
	assignment, err = h.ServicesSalesAssignment.GetByID(ctx, req.Id, int(req.Status), req.Search, int(req.TaskType), req.FinishDateFrom.AsTime(), req.FinishDateTo.AsTime())
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &crmService.GetSalesAssignmentDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &crmService.SalesAssignment{
			Id:   assignment.ID,
			Code: assignment.Code,
		},
	}
	return
}
