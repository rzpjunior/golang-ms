package service

import (
	"context"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/repository"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ISalesAssignmentSubmissionService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time, task int, outOfRoute int) (res []*dto.SalesAssignmentSubmissionResponse, total int64, err error)
}

type SalesAssignmentSubmissionService struct {
	opt                                opt.Options
	RepositorySalesAssignment          repository.ISalesAssignmentRepository
	RepositorySalesAssignmentItem      repository.ISalesAssignmentItemRepository
	RepositorySalesAssignmentObjective repository.ISalesAssignmentObjectiveRepository
	RepositoryCustomerAcquisition      repository.ICustomerAcquisitionRepository
}

func NewSalesAssignmentSubmissionService() ISalesAssignmentSubmissionService {
	return &SalesAssignmentSubmissionService{
		opt:                                global.Setup.Common,
		RepositorySalesAssignment:          repository.NewSalesAssignmentRepository(),
		RepositorySalesAssignmentItem:      repository.NewSalesAssignmentItemRepository(),
		RepositorySalesAssignmentObjective: repository.NewSalesAssignmentObjectiveRepository(),
		RepositoryCustomerAcquisition:      repository.NewCustomerAcquisitionRepository(),
	}
}

func (s *SalesAssignmentSubmissionService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time, task int, outOfRoute int) (res []*dto.SalesAssignmentSubmissionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentSubmissionService.Get")
	defer span.End()

	var salesAssignmentItems []*model.SalesAssignmentItem
	salesAssignmentItems, total, err = s.RepositorySalesAssignmentItem.GetSubmissions(ctx, offset, limit, status, search, orderBy, territoryID, salespersonID, submitDateFrom, submitDateTo, task, outOfRoute)
	for _, salesAssignmentItem := range salesAssignmentItems {

		var salesAssignment *model.SalesAssignment
		salesAssignment, err = s.RepositorySalesAssignment.GetByID(ctx, *salesAssignmentItem.SalesAssignmentID)
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

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

		var salesPerson *bridgeService.GetSalesPersonGPResponse
		salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: salesAssignmentItem.SalesPersonIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "salesperson")
			return
		}

		submission := &dto.SalesAssignmentSubmissionResponse{
			ID:                    salesAssignmentItem.ID,
			SalesAssignmentID:     salesAssignmentItem.SalesAssignmentID,
			AddressID:             salesAssignmentItem.AddressIDGP,
			CustomerAcquisitionID: salesAssignmentItem.CustomerAcquisitionID,
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
			Latitude:       salesAssignmentItem.Latitude,
			Longitude:      salesAssignmentItem.Longitude,
			Task:           salesAssignmentItem.Task,
			CustomerType:   salesAssignmentItem.CustomerType,
			ObjectiveCodes: salesAssignmentItem.ObjectiveCodes,
			ActualDistance: salesAssignmentItem.ActualDistance,
			OutOfRoute:     salesAssignmentItem.OutOfRoute,
			StartDate:      salesAssignmentItem.StartDate,
			EndDate:        salesAssignmentItem.EndDate,
			FinishDate:     salesAssignmentItem.FinishDate,
			SubmitDate:     salesAssignmentItem.SubmitDate,
			TaskImageUrls:  strings.Split(salesAssignmentItem.TaskImageUrl, ","),
			TaskAnswer:     salesAssignmentItem.TaskAnswer,
			Status:         salesAssignmentItem.Status,
			StatusConvert:  statusx.ConvertStatusValue(salesAssignmentItem.Status),
		}

		if salesAssignmentItem.AddressID != 0 {
			var address *bridgeService.GetAddressGPResponse
			address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
				Id: salesAssignmentItem.AddressIDGP,
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

		if salesAssignmentItem.CustomerAcquisitionID != 0 {
			var ca *model.CustomerAcquisition
			ca, err = s.RepositoryCustomerAcquisition.GetByID(ctx, salesAssignmentItem.CustomerAcquisitionID)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "address")
				return
			}
			submission.CustomerAcquisition = &dto.CustomerAcquisitionResponse{
				ID:   ca.ID,
				Code: ca.Code,
				Name: ca.Name,
			}
		}

		if salesAssignmentItem.ObjectiveCodes != "" {
			var objectiveCodes []*model.SalesAssignmentObjective
			codes := strings.Split(salesAssignmentItem.ObjectiveCodes, ",")
			objectiveCodes, _, err = s.RepositorySalesAssignmentObjective.Get(ctx, 0, 10, 0, "", codes, "")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			for _, val := range objectiveCodes {
				submission.ObjectiveValues = append(submission.ObjectiveValues, &dto.SalesAssignmentObjectiveResponse{
					ID:         val.ID,
					Code:       val.Code,
					Name:       val.Name,
					Objective:  val.Objective,
					SurveyLink: val.SurveyLink,
					Status:     val.Status,
					CreatedAt:  val.CreatedAt,
					CreatedBy: &dto.CreatedByResponse{
						ID: val.CreatedBy,
					},
					UpdatedAt: val.UpdatedAt,
				})
			}
		}

		res = append(res, submission)
	}

	return
}
