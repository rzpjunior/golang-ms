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

func (h *CrmGrpcHandler) GetSalesAssignmentObjectiveList(ctx context.Context, req *crmService.GetSalesAssignmentObjectiveListRequest) (res *crmService.GetSalesAssignmentObjectiveListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetSalesAssignmentObjectiveList")
	defer span.End()

	var objectives []*dto.SalesAssignmentObjectiveResponse
	objectives, _, err = h.ServicesSalesAssignmentObjective.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.Codes, req.OrderBy)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*crmService.SalesAssignmentObjective
	for _, objective := range objectives {
		data = append(data, &crmService.SalesAssignmentObjective{
			Id:         objective.ID,
			Code:       objective.Code,
			Name:       objective.Name,
			Objective:  objective.Objective,
			SurveyLink: objective.SurveyLink,
			CreatedAt:  timestamppb.New(objective.CreatedAt),
			CreatedBy:  objective.CreatedBy.ID,
			UpdatedAt:  timestamppb.New(objective.UpdatedAt),
			Status:     int32(objective.Status),
		})
	}

	res = &crmService.GetSalesAssignmentObjectiveListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *CrmGrpcHandler) GetSalesAssignmentObjectiveDetail(ctx context.Context, req *crmService.GetSalesAssignmentObjectiveDetailRequest) (res *crmService.GetSalesAssignmentObjectiveDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetSalesAssignmentObjectiveDetail")
	defer span.End()

	var objective dto.SalesAssignmentObjectiveResponse
	objective, err = h.ServicesSalesAssignmentObjective.GetByID(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &crmService.GetSalesAssignmentObjectiveDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &crmService.SalesAssignmentObjective{
			Id:         objective.ID,
			Code:       objective.Code,
			Name:       objective.Name,
			Objective:  objective.Objective,
			SurveyLink: objective.SurveyLink,
			CreatedAt:  timestamppb.New(objective.CreatedAt),
			CreatedBy:  objective.CreatedBy.ID,
			UpdatedAt:  timestamppb.New(objective.UpdatedAt),
			Status:     int32(objective.Status),
		},
	}
	return
}
