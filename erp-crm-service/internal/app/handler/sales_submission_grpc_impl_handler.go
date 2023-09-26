package handler

import (
	context "context"

	crmService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/crm_service"
)

func (h *CrmGrpcHandler) GetSalesSubmissionList(ctx context.Context, req *crmService.GetSalesSubmissionListRequest) (res *crmService.GetSalesAssignmentItemListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "CrmGrpcHandler.GetSalesSubmissionList")
	defer span.End()

	// var assignments []*dto.SalesAssignmentSubmissionResponse
	// assignments, _, err = h.ServicesSalesSubmission.Get(ctx, int(req.Offset), int(req.Limit), int(req.Status), req.Search, req.OrderBy, int(req.TerritoryId), int(req.SalespersonId), req.SubmitDateFrom.AsTime(), req.SubmitDateTo.AsTime(), int(req.Task), int(req.OutOfRoute))
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// var data []*crmService.SalesAssignmentItem
	// for _, assignment := range assignments {
	// 	var (
	// 		addressId  string
	// 		finishDate *timestamppb.Timestamp
	// 	)
	// 	if assignment.Address != nil {
	// 		addressId = assignment.Address.ID
	// 	}
	// 	if assignment.FinishDate != nil {
	// 		finishDate = timestamppb.New(*assignment.FinishDate)
	// 	}
	// 	var (
	// 		objectiveValues []*crmService.SalesAssignmentObjective
	// 		address         *crmService.Address
	// 		ca              *crmService.CustomerAcquisitionResponse
	// 	)
	// 	for _, val := range assignment.ObjectiveValues {
	// 		objectiveValues = append(objectiveValues, &crmService.SalesAssignmentObjective{
	// 			Id:         val.ID,
	// 			Code:       val.Code,
	// 			Name:       val.Name,
	// 			Objective:  val.Objective,
	// 			SurveyLink: val.SurveyLink,
	// 			Status:     int32(val.Status),
	// 			CreatedAt:  timestamppb.New(val.CreatedAt),
	// 			CreatedBy:  val.CreatedBy.ID,
	// 			UpdatedAt:  timestamppb.New(val.UpdatedAt),
	// 		})
	// 	}
	// 	if assignment.AddressID != 0 {
	// 		address = &crmService.Address{
	// 			Id:           assignment.Address.ID,
	// 			Code:         assignment.Address.Code,
	// 			CustomerName: assignment.Address.Name,
	// 		}
	// 	}
	// 	if assignment.CustomerAcquisitionID != 0 {
	// 		ca = &crmService.CustomerAcquisitionResponse{
	// 			Id:   assignment.CustomerAcquisition.ID,
	// 			Code: assignment.CustomerAcquisition.Code,
	// 			Name: assignment.CustomerAcquisition.Name,
	// 		}
	// 	}
	// 	data = append(data, &crmService.SalesAssignmentItem{
	// 		Id:                    assignment.ID,
	// 		SalesAssignmentId:     assignment.SalesAssignmentID,
	// 		SalesPersonId:         assignment.SalesPerson.ID,
	// 		AddressId:             addressId,
	// 		CustomerAcquisitionId: assignment.CustomerAcquisitionID,
	// 		Latitude:              assignment.Latitude,
	// 		Longitude:             assignment.Longitude,
	// 		Task:                  int32(assignment.Task),
	// 		CustomerType:          int32(assignment.CustomerType),
	// 		ObjectiveCodes:        assignment.ObjectiveCodes,
	// 		ActualDistance:        assignment.ActualDistance,
	// 		OutOfRoute:            int32(assignment.OutOfRoute),
	// 		ObectiveValues:        objectiveValues,
	// 		FinishDate:            finishDate,
	// 		SubmitDate:            timestamppb.New(assignment.SubmitDate),
	// 		EffectiveCall:         int32(assignment.EffectiveCall),
	// 		StartDate:             timestamppb.New(assignment.StartDate),
	// 		EndDate:               timestamppb.New(assignment.EndDate),
	// 		Address:               address,
	// 		CustomerAcquisition:   ca,
	// 		Status:                int32(assignment.Status),
	// 	})
	// }

	// res = &crmService.GetSalesAssignmentItemListResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data:    data,
	// }
	return
}
