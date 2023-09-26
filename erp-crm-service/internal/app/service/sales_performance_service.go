package service

import (
	"context"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/repository"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ISalesPerformanceService interface {
	Get(ctx context.Context, offset int, limit int, territoryID string, salespersonID string, startDateFrom time.Time, startDateTo time.Time) (res []*dto.SalesPerformanceResponse, total int64, err error)
	GetByID(ctx context.Context, id string, status int, startDateFrom time.Time, startDateTo time.Time, task int) (res *dto.SalesPerformanceDetailResponse, err error)
}

type SalesPerformanceService struct {
	opt                                     opt.Options
	RepositorySalesAssignment               repository.ISalesAssignmentRepository
	RepositorySalesAssignmentItemRepository repository.ISalesAssignmentItemRepository
	RepositoryCustomerAcquisition           repository.ICustomerAcquisitionRepository
}

func NewSalesPerformanceService() ISalesPerformanceService {
	return &SalesPerformanceService{
		opt:                                     global.Setup.Common,
		RepositorySalesAssignment:               repository.NewSalesAssignmentRepository(),
		RepositorySalesAssignmentItemRepository: repository.NewSalesAssignmentItemRepository(),
		RepositoryCustomerAcquisition:           repository.NewCustomerAcquisitionRepository(),
	}
}

func (s *SalesPerformanceService) Get(ctx context.Context, offset int, limit int, territoryID string, salespersonID string, startDateFrom time.Time, startDateTo time.Time) (res []*dto.SalesPerformanceResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPerformanceService.Get")
	defer span.End()

	// get sales assignment item group by salesperson
	var salesPerformances []*model.SalesAssignmentItem
	salesPerformances, total, err = s.RepositorySalesAssignmentItemRepository.GetGroupBySalespersonID(ctx, territoryID, salespersonID, startDateFrom, startDateTo, []int{1, 2})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesPerformance := range salesPerformances {
		var persentageVisit, persentageFollowUp, persentageEffectiveCall float64

		// get salesperson
		var salesPerson *bridgeService.GetSalesPersonGPResponse
		salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: salesPerformance.SalesPersonIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "salesperson")
			return
		}

		// get visit
		var saiPlanVisit []*model.SalesAssignmentItem
		saiPlanVisit, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Active)), int(statusx.ConvertStatusName(statusx.Finished)), int(statusx.ConvertStatusName(statusx.Failed))}, territoryID, salesPerformance.SalesPersonIDGP, startDateFrom, startDateTo, 1)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		planVisit := len(saiPlanVisit)

		var saiActualVisit []*model.SalesAssignmentItem
		saiActualVisit, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Finished))}, territoryID, salesPerformance.SalesPersonIDGP, startDateFrom, startDateTo, 1)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		actualVisit := len(saiActualVisit)

		if actualVisit != 0 {
			persentageVisit = float64(actualVisit) / float64(planVisit) * 100
		}

		// get follow up
		statusPlan := []int{int(statusx.ConvertStatusName(statusx.Active)), int(statusx.ConvertStatusName(statusx.Finished)), int(statusx.ConvertStatusName(statusx.Failed))}
		var saiPlanFollowUp []*model.SalesAssignmentItem
		saiPlanFollowUp, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, statusPlan, territoryID, salesPerformance.SalesPersonIDGP, startDateFrom, startDateTo, 2)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		planFollowUp := len(saiPlanFollowUp)

		statusActual := []int{int(statusx.ConvertStatusName(statusx.Finished))}
		var saiActualFollowUp []*model.SalesAssignmentItem
		saiActualFollowUp, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, statusActual, territoryID, salesPerformance.SalesPersonIDGP, startDateFrom, startDateTo, 2)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		actualFollowUp := len(saiActualFollowUp)

		if actualFollowUp != 0 {
			persentageFollowUp = float64(actualFollowUp) / float64(planFollowUp) * 100
		}

		// get customer acquisition
		var customerAcquisitions []*model.CustomerAcquisition
		customerAcquisitions, _, err = s.RepositoryCustomerAcquisition.GetPerformances(ctx, territoryID, salespersonID, startDateFrom, startDateTo)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		totalCa := len(customerAcquisitions)

		var glossaryCustomerType *configurationService.GetGlossaryDetailResponse
		glossaryCustomerType, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "sales_assignment_item",
			Attribute: "customer_type",
			ValueName: "Existing Customer",
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
			return
		}

		// get revenue effective call
		var revenueEffectiveCall float64
		tempAddressEffectiveCall := map[int64]bool{}
		var totalEffectiveCall int

		for _, actualVisitSalesAssignmentItem := range saiActualVisit {
			if actualVisitSalesAssignmentItem.CustomerType == int8(glossaryCustomerType.Data.ValueInt) && (actualVisitSalesAssignmentItem.Task == 1 || actualVisitSalesAssignmentItem.Task == 2) && actualVisitSalesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Finished) {
				var (
					salesOrders                *bridgeService.GetSalesOrderListResponse
					orderDateFrom, orderDateTo *timestamppb.Timestamp
				)
				if actualVisitSalesAssignmentItem.FinishDate != nil {
					orderDateFrom = timestamppb.New(timex.ToStartTime(*actualVisitSalesAssignmentItem.FinishDate))
					orderDateTo = timestamppb.New(timex.ToLastTime(*actualVisitSalesAssignmentItem.FinishDate))
				}
				salesOrders, _ = s.opt.Client.BridgeServiceGrpc.GetSalesOrderList(ctx, &bridgeService.GetSalesOrderListRequest{
					AddressId:     actualVisitSalesAssignmentItem.AddressID,
					SalespersonId: actualVisitSalesAssignmentItem.SalesPersonID,
					OrderDateFrom: orderDateFrom,
					OrderDateTo:   orderDateTo,
				})
				if !tempAddressEffectiveCall[actualVisitSalesAssignmentItem.AddressID] && len(salesOrders.Data) > 0 {
					tempAddressEffectiveCall[actualVisitSalesAssignmentItem.AddressID] = true
					totalEffectiveCall += 1
				}

				for _, salesOrder := range salesOrders.Data {
					if salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Cancelled)) && salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Deleted)) {
						revenueEffectiveCall += salesOrder.Total
					}
				}
			}
		}

		for _, actualFollowUpSalesAssignmentItem := range saiActualFollowUp {
			if actualFollowUpSalesAssignmentItem.CustomerType == int8(glossaryCustomerType.Data.ValueInt) && (actualFollowUpSalesAssignmentItem.Task == 1 || actualFollowUpSalesAssignmentItem.Task == 2) && actualFollowUpSalesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Finished) {
				var (
					salesOrders                *bridgeService.GetSalesOrderListResponse
					orderDateFrom, orderDateTo *timestamppb.Timestamp
				)
				if actualFollowUpSalesAssignmentItem.FinishDate != nil {
					orderDateFrom = timestamppb.New(timex.ToStartTime(*actualFollowUpSalesAssignmentItem.FinishDate))
					orderDateTo = timestamppb.New(timex.ToLastTime(*actualFollowUpSalesAssignmentItem.FinishDate))
				}
				salesOrders, _ = s.opt.Client.BridgeServiceGrpc.GetSalesOrderList(ctx, &bridgeService.GetSalesOrderListRequest{
					AddressId:     actualFollowUpSalesAssignmentItem.AddressID,
					SalespersonId: actualFollowUpSalesAssignmentItem.SalesPersonID,
					OrderDateFrom: orderDateFrom,
					OrderDateTo:   orderDateTo,
				})
				if !tempAddressEffectiveCall[actualFollowUpSalesAssignmentItem.AddressID] && len(salesOrders.Data) > 0 {
					tempAddressEffectiveCall[actualFollowUpSalesAssignmentItem.AddressID] = true
					totalEffectiveCall += 1
				}

				for _, salesOrder := range salesOrders.Data {
					if salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Cancelled)) && salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Deleted)) {
						revenueEffectiveCall += salesOrder.Total
					}
				}
			}
		}
		if actualVisit != 0 || actualFollowUp != 0 {
			persentageEffectiveCall = float64(totalEffectiveCall) / (float64(actualVisit) + float64(actualFollowUp)) * 100
		}

		var revenueTotal float64
		var salesOrders *bridgeService.GetSalesOrderListResponse
		salesOrders, _ = s.opt.Client.BridgeServiceGrpc.GetSalesOrderList(ctx, &bridgeService.GetSalesOrderListRequest{
			SalespersonId: salesPerformance.SalesPersonID,
			OrderDateFrom: timestamppb.New(timex.ToStartTime(startDateFrom)),
			OrderDateTo:   timestamppb.New(timex.ToLastTime(startDateTo)),
		})
		for _, salesOrder := range salesOrders.Data {
			if salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Cancelled)) && salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Deleted)) {
				revenueTotal += salesOrder.Total
			}
		}

		res = append(res, &dto.SalesPerformanceResponse{
			PlanVisit:                planVisit,
			VisitActual:              actualVisit,
			VisitPercentage:          persentageVisit,
			PlanFollowUp:             planFollowUp,
			FollowUpActual:           actualFollowUp,
			FollowUpPercentage:       persentageFollowUp,
			EffectiveCall:            totalEffectiveCall,
			EffectiveCallPercentage:  persentageEffectiveCall,
			RevenueEffectiveCall:     revenueEffectiveCall,
			RevenueTotal:             revenueTotal,
			TotalCustomerAcquisition: totalCa,
			Salesperson: &dto.SalespersonResponse{
				ID:   salesPerson.Data[0].Slprsnid,
				Code: salesPerson.Data[0].Slprsnid,
				Name: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
			},
		})
	}

	return
}

func (s *SalesPerformanceService) GetByID(ctx context.Context, id string, status int, startDateFrom time.Time, startDateTo time.Time, task int) (res *dto.SalesPerformanceDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPerformanceService.GetByID")
	defer span.End()

	var _ int64
	var persentageEffectiveCall float64

	// get salesperson
	var salesPerson *bridgeService.GetSalesPersonGPResponse
	salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// get territory
	var territory *bridgeService.GetSalesTerritoryGPResponse
	territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
		Id: salesPerson.Data[0].Slprsnid,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "territory")
		return
	}

	// get visit
	var planVisitSalesAssignmentItems []*model.SalesAssignmentItem
	planVisitSalesAssignmentItems, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Active)), int(statusx.ConvertStatusName(statusx.Finished)), int(statusx.ConvertStatusName(statusx.Failed))}, territory.Data[0].Salsterr, id, startDateFrom, startDateTo, 1)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	totalPlanVisit := len(planVisitSalesAssignmentItems)

	// total actual
	var actualVisitSalesAssignmentItems []*model.SalesAssignmentItem
	actualVisitSalesAssignmentItems, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Finished))}, territory.Data[0].Salsterr, id, startDateFrom, startDateTo, 1)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	totalActualVisit := len(actualVisitSalesAssignmentItems)

	// total failed
	var failedVisitSalesAssignmentItems []*model.SalesAssignmentItem
	failedVisitSalesAssignmentItems, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Failed))}, territory.Data[0].Salsterr, id, startDateFrom, startDateTo, 1)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	totalFailedVisit := len(failedVisitSalesAssignmentItems)

	// total cancelled
	var cancelledVisitSalesAssignmentItems []*model.SalesAssignmentItem
	cancelledVisitSalesAssignmentItems, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Cancelled))}, territory.Data[0].Salsterr, id, startDateFrom, startDateTo, 1)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	totalCancelledVisit := len(cancelledVisitSalesAssignmentItems)

	// get follow up
	var planFollowUpSalesAssignmentItems []*model.SalesAssignmentItem
	planFollowUpSalesAssignmentItems, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Active)), int(statusx.ConvertStatusName(statusx.Finished)), int(statusx.ConvertStatusName(statusx.Failed))}, territory.Data[0].Salsterr, id, startDateFrom, startDateTo, 2)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	totalPlanFollowUp := len(planFollowUpSalesAssignmentItems)

	// total actual
	var actualFollowUpSalesAssignmentItems []*model.SalesAssignmentItem
	actualFollowUpSalesAssignmentItems, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Finished))}, territory.Data[0].Salsterr, id, startDateFrom, startDateTo, 2)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	totalActualFollowUp := len(actualFollowUpSalesAssignmentItems)

	// total failed
	var failedFollowUpSalesAssignmentItems []*model.SalesAssignmentItem
	failedFollowUpSalesAssignmentItems, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Failed))}, territory.Data[0].Salsterr, id, startDateFrom, startDateTo, 2)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	totalFailedFollowUp := len(failedFollowUpSalesAssignmentItems)

	// total cancelled
	var cancelledFollowUpSalesAssignmentItems []*model.SalesAssignmentItem
	cancelledFollowUpSalesAssignmentItems, _, err = s.RepositorySalesAssignmentItemRepository.GetByTask(ctx, []int{int(statusx.ConvertStatusName(statusx.Cancelled))}, territory.Data[0].Salsterr, id, startDateFrom, startDateTo, 2)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	totalCancelledFollowUp := len(cancelledFollowUpSalesAssignmentItems)

	var glossaryCustomerType *configurationService.GetGlossaryDetailResponse
	glossaryCustomerType, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
		Table:     "sales_assignment_item",
		Attribute: "customer_type",
		ValueName: "Existing Customer",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}

	// get customer acquisition
	var customerAcquisitions []*model.CustomerAcquisition
	customerAcquisitions, _, err = s.RepositoryCustomerAcquisition.GetPerformances(ctx, territory.Data[0].Salsterr, id, startDateFrom, startDateTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// get revenue effective call
	var revenueEffectiveCall float64
	tempAddressEffectiveCall := map[int64]bool{}
	var totalEffectiveCall int64
	var totalOutOfRouteVisit, totalOutOfRouteFollowUp int

	for _, actualVisitSalesAssignmentItem := range actualVisitSalesAssignmentItems {
		if actualVisitSalesAssignmentItem.CustomerType == int8(glossaryCustomerType.Data.ValueInt) && (actualVisitSalesAssignmentItem.Task == 1 || actualVisitSalesAssignmentItem.Task == 2) && actualVisitSalesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Finished) {
			if actualVisitSalesAssignmentItem.OutOfRoute == 2 {
				totalOutOfRouteVisit += 1
			}

			var (
				salesOrders                *bridgeService.GetSalesOrderListResponse
				orderDateFrom, orderDateTo *timestamppb.Timestamp
			)
			if actualVisitSalesAssignmentItem.FinishDate != nil {
				orderDateFrom = timestamppb.New(timex.ToStartTime(*actualVisitSalesAssignmentItem.FinishDate))
				orderDateTo = timestamppb.New(timex.ToLastTime(*actualVisitSalesAssignmentItem.FinishDate))
			}
			salesOrders, _ = s.opt.Client.BridgeServiceGrpc.GetSalesOrderList(ctx, &bridgeService.GetSalesOrderListRequest{
				AddressId:     actualVisitSalesAssignmentItem.AddressID,
				SalespersonId: actualVisitSalesAssignmentItem.SalesPersonID,
				OrderDateFrom: orderDateFrom,
				OrderDateTo:   orderDateTo,
			})
			if !tempAddressEffectiveCall[actualVisitSalesAssignmentItem.AddressID] && len(salesOrders.Data) > 0 {
				tempAddressEffectiveCall[actualVisitSalesAssignmentItem.AddressID] = true
				totalEffectiveCall += 1
				actualVisitSalesAssignmentItem.EffectiveCall = 1
			}

			for _, salesOrder := range salesOrders.Data {
				if salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Cancelled)) && salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Deleted)) {
					revenueEffectiveCall += salesOrder.Total
				}

			}
		}
	}

	for _, actualFollowUpSalesAssignmentItem := range actualFollowUpSalesAssignmentItems {
		if actualFollowUpSalesAssignmentItem.CustomerType == int8(glossaryCustomerType.Data.ValueInt) && (actualFollowUpSalesAssignmentItem.Task == 1 || actualFollowUpSalesAssignmentItem.Task == 2) && actualFollowUpSalesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Finished) {
			if actualFollowUpSalesAssignmentItem.OutOfRoute == 2 {
				totalOutOfRouteFollowUp += 1
			}

			var (
				salesOrders                *bridgeService.GetSalesOrderListResponse
				orderDateFrom, orderDateTo *timestamppb.Timestamp
			)
			if actualFollowUpSalesAssignmentItem.FinishDate != nil {
				orderDateFrom = timestamppb.New(timex.ToStartTime(*actualFollowUpSalesAssignmentItem.FinishDate))
				orderDateTo = timestamppb.New(timex.ToLastTime(*actualFollowUpSalesAssignmentItem.FinishDate))
			}
			salesOrders, _ = s.opt.Client.BridgeServiceGrpc.GetSalesOrderList(ctx, &bridgeService.GetSalesOrderListRequest{
				AddressId:     actualFollowUpSalesAssignmentItem.AddressID,
				SalespersonId: actualFollowUpSalesAssignmentItem.SalesPersonID,
				OrderDateFrom: orderDateFrom,
				OrderDateTo:   orderDateTo,
			})
			if !tempAddressEffectiveCall[actualFollowUpSalesAssignmentItem.AddressID] && len(salesOrders.Data) > 0 {
				tempAddressEffectiveCall[actualFollowUpSalesAssignmentItem.AddressID] = true
				totalEffectiveCall += 1
				actualFollowUpSalesAssignmentItem.EffectiveCall = 1
			}

			for _, salesOrder := range salesOrders.Data {
				if salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Cancelled)) && salesOrder.Status != int32(statusx.ConvertStatusName(statusx.Deleted)) {
					revenueEffectiveCall += salesOrder.Total
				}

			}
		}
	}
	if totalEffectiveCall != 0 || totalActualVisit != 0 || totalActualFollowUp != 0 {
		persentageEffectiveCall = float64(totalEffectiveCall) / (float64(totalActualVisit) + float64(totalActualFollowUp)) * 100
	}

	// generate sales assignment response
	var salesAssignmentItemsResponse []*dto.SalesAssignmentSubmissionResponse

	for _, sai := range planVisitSalesAssignmentItems {
		var salesAssignment *model.SalesAssignment
		salesAssignment, err = s.RepositorySalesAssignment.GetByID(ctx, *sai.SalesAssignmentID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		// get territory
		var territory *bridgeService.GetSalesTerritoryGPResponse
		territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: salesAssignment.TerritoryIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "territory")
			return
		}

		// plan
		submission := &dto.SalesAssignmentSubmissionResponse{
			ID:                sai.ID,
			SalesAssignmentID: sai.SalesAssignmentID,
			SalesPerson: &dto.SalespersonResponse{
				ID:   salesPerson.Data[0].Slprsnid,
				Code: salesPerson.Data[0].Slprsnid,
				Name: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
			},
			Territory: &dto.TerritoryResponse{
				ID:          territory.Data[0].Salsterr,
				Code:        territory.Data[0].Salsterr,
				Description: territory.Data[0].Slterdsc,
			},
			CustomerAcquisitionID: sai.CustomerAcquisitionID,
			Latitude:              sai.Latitude,
			Longitude:             sai.Longitude,
			Task:                  sai.Task,
			CustomerType:          sai.CustomerType,
			ObjectiveCodes:        sai.ObjectiveCodes,
			ActualDistance:        sai.ActualDistance,
			OutOfRoute:            sai.OutOfRoute,
			StartDate:             sai.StartDate,
			EndDate:               sai.EndDate,
			FinishDate:            sai.FinishDate,
			SubmitDate:            sai.SubmitDate,
			TaskImageUrls:         strings.Split(sai.TaskImageUrl, ","),
			TaskAnswer:            sai.TaskAnswer,
			Status:                sai.Status,
			StatusConvert:         statusx.ConvertStatusValue(sai.Status),
		}

		statusEffectiveCall := 2
		// check effective call
		for idTempAddressEC, statusTempAddressEC := range tempAddressEffectiveCall {
			if sai.AddressID == idTempAddressEC && statusTempAddressEC {
				statusEffectiveCall = 1
			}
		}

		submission.EffectiveCall = int8(statusEffectiveCall)

		// get address
		if sai.AddressIDGP != "" {
			var address *bridgeService.GetAddressGPResponse
			address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
				Id: sai.AddressIDGP,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "address")
				return
			}
			submission.Address = &dto.AddressResponse{
				ID:   address.Data[0].Custnmbr,
				Code: address.Data[0].Custnmbr,
				Name: address.Data[0].Custname,
			}
		}
		salesAssignmentItemsResponse = append(salesAssignmentItemsResponse, submission)
	}
	// cancelled
	for _, sai := range cancelledVisitSalesAssignmentItems {
		var salesAssignment *model.SalesAssignment
		salesAssignment, err = s.RepositorySalesAssignment.GetByID(ctx, *sai.SalesAssignmentID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		// get territory
		var territory *bridgeService.GetSalesTerritoryGPResponse
		territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: salesAssignment.TerritoryIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "territory")
			return
		}

		submission := &dto.SalesAssignmentSubmissionResponse{
			ID:                sai.ID,
			SalesAssignmentID: sai.SalesAssignmentID,
			SalesPerson: &dto.SalespersonResponse{
				ID:   salesPerson.Data[0].Slprsnid,
				Code: salesPerson.Data[0].Slprsnid,
				Name: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
			},
			Territory: &dto.TerritoryResponse{
				ID:          territory.Data[0].Salsterr,
				Code:        territory.Data[0].Salsterr,
				Description: territory.Data[0].Slterdsc,
			},
			CustomerAcquisitionID: sai.CustomerAcquisitionID,
			Latitude:              sai.Latitude,
			Longitude:             sai.Longitude,
			Task:                  sai.Task,
			CustomerType:          sai.CustomerType,
			ObjectiveCodes:        sai.ObjectiveCodes,
			ActualDistance:        sai.ActualDistance,
			OutOfRoute:            sai.OutOfRoute,
			StartDate:             sai.StartDate,
			EndDate:               sai.EndDate,
			FinishDate:            sai.FinishDate,
			SubmitDate:            sai.SubmitDate,
			TaskImageUrls:         strings.Split(sai.TaskImageUrl, ","),
			TaskAnswer:            sai.TaskAnswer,
			Status:                sai.Status,
			StatusConvert:         statusx.ConvertStatusValue(sai.Status),
			EffectiveCall:         sai.EffectiveCall,
		}
		if sai.AddressID != 0 {
			var address *bridgeService.GetAddressGPResponse
			address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
				Id: sai.AddressIDGP,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "address")
				return
			}
			submission.Address = &dto.AddressResponse{
				ID:   address.Data[0].Custnmbr,
				Code: address.Data[0].Custnmbr,
				Name: address.Data[0].Custname,
			}
		}
		salesAssignmentItemsResponse = append(salesAssignmentItemsResponse, submission)
	}

	// plan
	for _, sai := range planFollowUpSalesAssignmentItems {
		var salesAssignment *model.SalesAssignment
		salesAssignment, err = s.RepositorySalesAssignment.GetByID(ctx, *sai.SalesAssignmentID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		// get territory
		var territory *bridgeService.GetSalesTerritoryGPResponse
		territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: salesAssignment.TerritoryIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "territory")
			return
		}

		// plan
		submission := &dto.SalesAssignmentSubmissionResponse{
			ID:                sai.ID,
			SalesAssignmentID: sai.SalesAssignmentID,
			SalesPerson: &dto.SalespersonResponse{
				ID:   salesPerson.Data[0].Slprsnid,
				Code: salesPerson.Data[0].Slprsnid,
				Name: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
			},
			Territory: &dto.TerritoryResponse{
				ID:          territory.Data[0].Salsterr,
				Code:        territory.Data[0].Salsterr,
				Description: territory.Data[0].Slterdsc,
			},
			CustomerAcquisitionID: sai.CustomerAcquisitionID,
			Latitude:              sai.Latitude,
			Longitude:             sai.Longitude,
			Task:                  sai.Task,
			CustomerType:          sai.CustomerType,
			ObjectiveCodes:        sai.ObjectiveCodes,
			ActualDistance:        sai.ActualDistance,
			OutOfRoute:            sai.OutOfRoute,
			StartDate:             sai.StartDate,
			EndDate:               sai.EndDate,
			FinishDate:            sai.FinishDate,
			SubmitDate:            sai.SubmitDate,
			TaskImageUrls:         strings.Split(sai.TaskImageUrl, ","),
			TaskAnswer:            sai.TaskAnswer,
			Status:                sai.Status,
			StatusConvert:         statusx.ConvertStatusValue(sai.Status),
			EffectiveCall:         sai.EffectiveCall,
		}
		if sai.AddressID != 0 {
			var address *bridgeService.GetAddressGPResponse
			address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
				Id: sai.AddressIDGP,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "address")
				return
			}
			submission.Address = &dto.AddressResponse{
				ID:   address.Data[0].Custnmbr,
				Code: address.Data[0].Custnmbr,
				Name: address.Data[0].Custname,
			}
		}
		salesAssignmentItemsResponse = append(salesAssignmentItemsResponse, submission)
	}
	// cancelled
	for _, sai := range cancelledFollowUpSalesAssignmentItems {
		var salesAssignment *model.SalesAssignment
		salesAssignment, err = s.RepositorySalesAssignment.GetByID(ctx, *sai.SalesAssignmentID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		// get territory
		var territory *bridgeService.GetSalesTerritoryGPResponse
		territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: salesAssignment.TerritoryIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "territory")
			return
		}

		submission := &dto.SalesAssignmentSubmissionResponse{
			ID:                sai.ID,
			SalesAssignmentID: sai.SalesAssignmentID,
			SalesPerson: &dto.SalespersonResponse{
				ID:   salesPerson.Data[0].Slprsnid,
				Code: salesPerson.Data[0].Slprsnid,
				Name: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
			},
			Territory: &dto.TerritoryResponse{
				ID:          territory.Data[0].Salsterr,
				Code:        territory.Data[0].Salsterr,
				Description: territory.Data[0].Slterdsc,
			},
			CustomerAcquisitionID: sai.CustomerAcquisitionID,
			Latitude:              sai.Latitude,
			Longitude:             sai.Longitude,
			Task:                  sai.Task,
			CustomerType:          sai.CustomerType,
			ObjectiveCodes:        sai.ObjectiveCodes,
			ActualDistance:        sai.ActualDistance,
			OutOfRoute:            sai.OutOfRoute,
			StartDate:             sai.StartDate,
			EndDate:               sai.EndDate,
			FinishDate:            sai.FinishDate,
			SubmitDate:            sai.SubmitDate,
			TaskImageUrls:         strings.Split(sai.TaskImageUrl, ","),
			TaskAnswer:            sai.TaskAnswer,
			Status:                sai.Status,
			StatusConvert:         statusx.ConvertStatusValue(sai.Status),
			EffectiveCall:         sai.EffectiveCall,
		}
		if sai.AddressID != 0 {
			var address *bridgeService.GetAddressGPResponse
			address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
				Id: sai.AddressIDGP,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "address")
				return
			}
			submission.Address = &dto.AddressResponse{
				ID:   address.Data[0].Custnmbr,
				Code: address.Data[0].Custnmbr,
				Name: address.Data[0].Custname,
			}
		}
		salesAssignmentItemsResponse = append(salesAssignmentItemsResponse, submission)
	}

	var customerAcquisitionsResponse []*dto.CustomerAcquisitionResponse
	for _, ca := range customerAcquisitions {
		var territory *bridgeService.GetSalesTerritoryGPResponse
		territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: ca.TerritoryIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "territory")
			return
		}
		customerAcquisitionsResponse = append(customerAcquisitionsResponse, &dto.CustomerAcquisitionResponse{
			ID:               ca.ID,
			Task:             ca.Task,
			Name:             ca.Name,
			PhoneNumber:      ca.PhoneNumber,
			Latitude:         ca.Latitude,
			Longitude:        ca.Longitude,
			AddressName:      ca.AddressName,
			FoodApp:          ca.FoodApp,
			PotentialRevenue: ca.PotentialRevenue,
			TaskImageUrl:     ca.TaskImageUrl,
			Salesperson: &dto.SalespersonResponse{
				ID:   salesPerson.Data[0].Slprsnid,
				Code: salesPerson.Data[0].Slprsnid,
				Name: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
			},
			Territory: &dto.TerritoryResponse{
				ID:          territory.Data[0].Salsterr,
				Code:        territory.Data[0].Salsterr,
				Description: territory.Data[0].Slterdsc,
			},
			FinishDate:    ca.FinishDate,
			SubmitDate:    ca.SubmitDate,
			CreatedAt:     ca.CreatedAt,
			UpdatedAt:     ca.UpdatedAt,
			Status:        ca.Status,
			StatusConvert: statusx.ConvertStatusValue(ca.Status),
		})
	}

	res = &dto.SalesPerformanceDetailResponse{
		VisitTracker: &dto.PerformanceTrackerResponse{
			TotalPlan:       int(totalPlanVisit),
			TotalFinished:   int(totalActualVisit),
			TotalFailed:     int(totalFailedVisit),
			TotalCancelled:  int(totalCancelledVisit),
			TotalOutOfRoute: int(totalOutOfRouteVisit),
		},
		FollowUpTracker: &dto.PerformanceTrackerResponse{
			TotalPlan:       int(totalPlanFollowUp),
			TotalFinished:   int(totalActualFollowUp),
			TotalFailed:     int(totalFailedFollowUp),
			TotalCancelled:  int(totalCancelledFollowUp),
			TotalOutOfRoute: int(totalOutOfRouteFollowUp),
		},
		EffectiveCallPercentage:   persentageEffectiveCall,
		SalesAssignmentSubmission: salesAssignmentItemsResponse,
		CustomerAcquisition:       customerAcquisitionsResponse,
	}

	return
}
