package service

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/reportx"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/repository"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type ISalesAssignmentService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, startDateFrom time.Time, startDateTo time.Time, endDateFrom time.Time, endDateTo time.Time) (res []*dto.SalesAssignmentResponse, total int64, err error)
	GetByID(ctx context.Context, id int64, status int, search string, taskType int, finishDateFrom time.Time, finishDateTo time.Time) (res dto.SalesAssignmentResponse, err error)
	CancelBatch(ctx context.Context, id int64) (res dto.SalesAssignmentResponse, err error)
	CancelItem(ctx context.Context, id int64) (res dto.SalesAssignmentResponse, err error)
	Export(ctx context.Context, territoryID string) (res dto.SalesAssignmentExportResponse, err error)
	Import(ctx context.Context, req dto.SalesAssignmentImportRequest) (err error)
}

type SalesAssignmentService struct {
	opt                                opt.Options
	RepositorySalesAssignment          repository.ISalesAssignmentRepository
	RepositorySalesAssignmentItem      repository.ISalesAssignmentItemRepository
	RepositoryCustomerAcquisition      repository.ICustomerAcquisitionRepository
	RepositorySalesAssignmentObjective repository.ISalesAssignmentObjectiveRepository
}

func NewSalesAssignmentService() ISalesAssignmentService {
	return &SalesAssignmentService{
		opt:                                global.Setup.Common,
		RepositorySalesAssignment:          repository.NewSalesAssignmentRepository(),
		RepositorySalesAssignmentItem:      repository.NewSalesAssignmentItemRepository(),
		RepositoryCustomerAcquisition:      repository.NewCustomerAcquisitionRepository(),
		RepositorySalesAssignmentObjective: repository.NewSalesAssignmentObjectiveRepository(),
	}
}

func (s *SalesAssignmentService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, startDateFrom time.Time, startDateTo time.Time, endDateFrom time.Time, endDateTo time.Time) (res []*dto.SalesAssignmentResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentService.Get")
	defer span.End()

	var salesAssignments []*model.SalesAssignment
	salesAssignments, total, err = s.RepositorySalesAssignment.Get(ctx, offset, limit, status, search, orderBy, territoryID, startDateFrom, startDateTo, endDateFrom, endDateTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesAssignment := range salesAssignments {

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

		var statusActiveIsExist bool
		var statusCancelIsExist bool
		var statusFailedIsExist bool
		var salesAssignmentItems []*model.SalesAssignmentItem
		salesAssignmentItems, err = s.RepositorySalesAssignmentItem.GetBySalesAssignmentItemID(ctx, salesAssignment.ID, 0, 0, time.Time{}, time.Time{})
		for _, salesAssignmentItem := range salesAssignmentItems {
			if salesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Active) {
				statusActiveIsExist = true
			}
			if salesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Cancelled) {
				statusCancelIsExist = true
			}
			if salesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Failed) {
				statusFailedIsExist = true
			}
		}

		var salesAssignmentStatus int8
		if statusActiveIsExist {
			salesAssignmentStatus = statusx.ConvertStatusName(statusx.Active)
		} else if statusCancelIsExist {
			if statusFailedIsExist {
				salesAssignmentStatus = statusx.ConvertStatusName(statusx.Finished)
			} else {
				salesAssignmentStatus = statusx.ConvertStatusName(statusx.Cancelled)
			}
		} else {
			salesAssignmentStatus = statusx.ConvertStatusName(statusx.Finished)
		}

		if status != 0 {
			if status == int(salesAssignmentStatus) {
				res = append(res, &dto.SalesAssignmentResponse{
					ID:   salesAssignment.ID,
					Code: salesAssignment.Code,
					Territory: &dto.TerritoryResponse{
						ID:          territory.Data[0].Salsterr,
						Code:        territory.Data[0].Salsterr,
						Description: territory.Data[0].Slterdsc,
					},
					StartDate:     salesAssignment.StartDate,
					EndDate:       salesAssignment.EndDate,
					Status:        salesAssignmentStatus,
					StatusConvert: statusx.ConvertStatusValue(salesAssignmentStatus),
				})
			}
		} else {
			res = append(res, &dto.SalesAssignmentResponse{
				ID:   salesAssignment.ID,
				Code: salesAssignment.Code,
				Territory: &dto.TerritoryResponse{
					ID:          territory.Data[0].Salsterr,
					Code:        territory.Data[0].Salsterr,
					Description: territory.Data[0].Slterdsc,
				},
				StartDate:     salesAssignment.StartDate,
				EndDate:       salesAssignment.EndDate,
				Status:        salesAssignmentStatus,
				StatusConvert: statusx.ConvertStatusValue(salesAssignmentStatus),
			})
		}
	}

	return
}

func (s *SalesAssignmentService) GetByID(ctx context.Context, id int64, status int, search string, taskType int, finishDateFrom time.Time, finishDateTo time.Time) (res dto.SalesAssignmentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentService.GetByID")
	defer span.End()

	var salesAssignment *model.SalesAssignment
	salesAssignment, err = s.RepositorySalesAssignment.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var salesAssignmentItemResponses []*dto.SalesAssignmentItemResponse

	var statusActiveIsExist bool
	var statusCancelIsExist bool
	var statusFailedIsExist bool
	var salesAssignmentItems []*model.SalesAssignmentItem
	salesAssignmentItems, err = s.RepositorySalesAssignmentItem.GetBySalesAssignmentItemID(ctx, id, status, taskType, finishDateFrom, finishDateTo)
	for _, salesAssignmentItem := range salesAssignmentItems {
		if salesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Active) {
			statusActiveIsExist = true
		}
		if salesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Cancelled) {
			statusCancelIsExist = true
		}
		if salesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Failed) {
			statusFailedIsExist = true
		}

		SalesAssignmentItemResponse := &dto.SalesAssignmentItemResponse{
			ID:                    salesAssignmentItem.ID,
			SalesAssignmentID:     &salesAssignment.ID,
			CustomerAcquisitionID: salesAssignmentItem.CustomerAcquisitionID,
			Latitude:              salesAssignmentItem.Latitude,
			Longitude:             salesAssignmentItem.Longitude,
			Task:                  salesAssignmentItem.Task,
			CustomerType:          salesAssignmentItem.CustomerType,
			ObjectiveCodes:        salesAssignmentItem.ObjectiveCodes,
			ActualDistance:        salesAssignmentItem.ActualDistance,
			OutOfRoute:            salesAssignmentItem.OutOfRoute,
			StartDate:             salesAssignmentItem.StartDate,
			EndDate:               salesAssignmentItem.EndDate,
			FinishDate:            salesAssignmentItem.FinishDate,
			SubmitDate:            salesAssignmentItem.SubmitDate,
			TaskImageUrls:         strings.Split(salesAssignmentItem.TaskImageUrl, ","),
			TaskAnswer:            salesAssignmentItem.TaskAnswer,
			Status:                salesAssignmentItem.Status,
			StatusConvert:         statusx.ConvertStatusValue(salesAssignmentItem.Status),
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

			SalesAssignmentItemResponse.Address = &dto.AddressResponse{
				ID:   address.Data[0].Custnmbr,
				Code: address.Data[0].Custnmbr,
				Name: address.Data[0].Custname,
			}
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

		salespersonName := salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln
		if strings.Contains(salespersonName, search) {
			SalesAssignmentItemResponse.SalesPerson = &dto.SalespersonResponse{
				ID:   salesPerson.Data[0].Slprsnid,
				Code: salesPerson.Data[0].Slprsnid,
				Name: salespersonName,
			}

			salesAssignmentItemResponses = append(salesAssignmentItemResponses, SalesAssignmentItemResponse)
		}
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

	var salesAssignmentStatus int8
	if statusActiveIsExist {
		salesAssignmentStatus = statusx.ConvertStatusName(statusx.Active)
	} else if statusCancelIsExist {
		if statusFailedIsExist {
			salesAssignmentStatus = statusx.ConvertStatusName(statusx.Finished)
		} else {
			salesAssignmentStatus = statusx.ConvertStatusName(statusx.Cancelled)
		}
	} else {
		salesAssignmentStatus = statusx.ConvertStatusName(statusx.Finished)
	}

	res = dto.SalesAssignmentResponse{
		ID:   salesAssignment.ID,
		Code: salesAssignment.Code,
		Territory: &dto.TerritoryResponse{
			ID:          territory.Data[0].Salsterr,
			Code:        territory.Data[0].Salsterr,
			Description: territory.Data[0].Slterdsc,
		},
		StartDate:           salesAssignment.StartDate,
		EndDate:             salesAssignment.EndDate,
		Status:              salesAssignmentStatus,
		StatusConvert:       statusx.ConvertStatusValue(salesAssignmentStatus),
		SalesAssignmentItem: salesAssignmentItemResponses,
	}

	return
}

func (s *SalesAssignmentService) Export(ctx context.Context, territoryID string) (res dto.SalesAssignmentExportResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentService.Export")
	defer span.End()

	// validation territory
	var territory *bridgeService.GetSalesTerritoryGPResponse
	territory, _ = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
		Id: territoryID,
	})

	var ex reportx.Excelx
	header := []string{
		"Territory_Code",
		"Territory_Name",
		"Customer_Type",
		"Customer_Code",
		"Customer_Name",
		"Sub_District",
		"District",
		"Salesperson_Code",
		"Salesperson_Name",
		"Task",
		"Visit_Date",
		"Objective_Codes",
	}

	var cells []interface{}

	var addresses *bridgeService.GetAddressGPResponse
	addresses, _ = s.opt.Client.BridgeServiceGrpc.GetAddressGPList(ctx, &bridgeService.GetAddressGPListRequest{
		Limit:  100,
		Offset: 0,
		// TerritoryId: territoryID,
		// Status:      int32(statusx.ConvertStatusName(statusx.Active)),
	})

	for _, address := range addresses.Data {

		dataCells := dto.SalesAssignmentTemplate{
			TerritoryCode: territory.Data[0].Salsterr,
			TerritoryName: territory.Data[0].Slterdsc,
			CustomerType:  "Existing Customer",
			CustomerCode:  address.Custnmbr,
			CustomerName:  address.Custname,
			SubDistrict:   "Dummy Sub District",
			District:      "Dummy District",
		}

		if address.Slprsnid != "" {
			var salesPerson *bridgeService.GetSalesPersonGPResponse
			salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
				Id: address.Slprsnid,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "salesperson")
				return
			}

			dataCells.SalespersonCode = salesPerson.Data[0].Slprsnid
			dataCells.SalespersonName = salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln
		}
		cells = append(cells, dataCells)
	}

	var customerAcquisitions []*model.CustomerAcquisition
	customerAcquisitions, _, _ = s.RepositoryCustomerAcquisition.GetByTerritoryID(ctx, territoryID)
	for _, customerAcquisition := range customerAcquisitions {
		var salesPerson *bridgeService.GetSalesPersonGPResponse
		salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: customerAcquisition.SalespersonIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "salesperson")
			return
		}

		dataCells := dto.SalesAssignmentTemplate{
			TerritoryCode:   territory.Data[0].Salsterr,
			TerritoryName:   territory.Data[0].Slterdsc,
			CustomerType:    "Customer Acquisition",
			CustomerCode:    customerAcquisition.Code,
			CustomerName:    customerAcquisition.Name,
			SubDistrict:     "Dummy SubDistrict",
			District:        "Dummy District",
			SalespersonCode: salesPerson.Data[0].Slprsnid,
			SalespersonName: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
		}
		cells = append(cells, dataCells)
	}

	ex.Sheets = append(ex.Sheets, reportx.Sheet{
		WithNumbering: true,
		Name:          "Sheet1",
		Headers:       header,
		Bodys:         cells,
	})

	fileName := fmt.Sprintf("TaskAssignment_%s_%s.xlsx", time.Now().Format(timex.InFormatDate), utils.GenerateRandomDoc(5))
	fileLocation, err := reportx.GenerateXlsx(fileName, ex)

	info, err := s.opt.S3x.UploadPrivateFile(ctx, s.opt.Config.S3.BucketName, fileName, fileLocation, "application/xlsx")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = fmt.Errorf("failed to upload file | %v", err)
		return
	}

	os.Remove(fileLocation)

	res = dto.SalesAssignmentExportResponse{
		Url: info,
	}

	return
}

func (s *SalesAssignmentService) Import(ctx context.Context, req dto.SalesAssignmentImportRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentService.Import")
	defer span.End()

	// validate territory
	var territory *bridgeService.GetSalesTerritoryGPResponse
	territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
		Id: req.TerritoryCode,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "territory")
		return
	}

	for i, assignment := range req.Assignments {
		// validate salesperson id
		var salesperson *bridgeService.GetSalesPersonGPResponse
		salesperson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: assignment.SalespersonCode,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRowInvalid("assignments", i, "salesperson_code")
			return
		}
		assignment.SalespersonID = salesperson.Data[0].Slprsnid

		// if salesperson.Data[0].Status != int32(statusx.ConvertStatusName(statusx.Active)) {
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorRowMustActive("assignments", i, "salesperson_id")
		// 	return
		// }

		if assignment.CustomerType == "Existing Customer" {
			// validate address
			var address *bridgeService.GetAddressGPResponse
			address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
				Id: assignment.CustomerCode,
			})
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRowInvalid("assignments", i, "customer_code")
				return
			}
			assignment.AddressID = address.Data[0].Custnmbr

			if salesperson.Data[0].Slprsnid != address.Data[0].Slprsnid {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRowMustSame("assignments", i, "salesperson_id", "salesperson_id")
				return
			}
		} else if assignment.CustomerType == "Customer Acquisition" {
			// validate customer acquisition
			var customerAcquisition *model.CustomerAcquisition
			customerAcquisition, err = s.RepositoryCustomerAcquisition.GetByCode(ctx, assignment.CustomerCode)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRowInvalid("assignments", i, "customer_code")
				return
			}
			assignment.CustomerAcquisitionID = customerAcquisition.ID

			if salesperson.Data[0].Slprsnid != customerAcquisition.SalespersonIDGP {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRowMustSame("assignments", i, "salesperson_id", "salesperson_id")
				return
			}
		} else {
			err = edenlabs.ErrorRowInvalid("assignments", i, "customer_type")
			return
		}

		currentTime := time.Now()
		var visitDate time.Time
		visitDate, err = time.Parse(timex.InFormatDate, assignment.VisitDate)
		if err != nil {
			err = edenlabs.ErrorRowInvalid("assignments", i, "visit_date")
			return
		}

		if visitDate.Before(currentTime) {
			err = edenlabs.ErrorRowMustEqualOrGreater("assignments", i, "visit_date", "current time")
			return
		}

		if i == 0 {
			req.StartDate = visitDate
			req.EndDate = visitDate
		}

		if visitDate.Before(req.StartDate) {
			req.StartDate = visitDate
		}

		if visitDate.After(req.EndDate) {
			req.EndDate = visitDate
		}

		if assignment.ObjectiveCodes != "" {
			objectiveCodes := strings.Split(assignment.ObjectiveCodes, ",")
			for _, objectiveCode := range objectiveCodes {
				var salesAssignmentObjective *model.SalesAssignmentObjective
				salesAssignmentObjective, err = s.RepositorySalesAssignmentObjective.GetByCode(ctx, objectiveCode)
				if err != nil {
					span.RecordError(err)
					s.opt.Logger.AddMessage(log.ErrorLevel, err)
					err = edenlabs.ErrorRowInvalid("assignments", i, "sales_assignment_code")
					return
				}
				if salesAssignmentObjective.Status != statusx.ConvertStatusName(statusx.Active) {
					err = edenlabs.ErrorRowMustActive("assignments", i, "sales_assignment_code")
					return
				}
			}
		}

		var glossaryTask *configurationService.GetGlossaryDetailResponse
		glossaryTask, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "sales_assignment_item",
			Attribute: "task",
			ValueName: assignment.Task,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRowInvalid("assignments", i, "task")
			return
		}
		assignment.TaskValue = int8(glossaryTask.Data.ValueInt)

		var glossaryCustomerType *configurationService.GetGlossaryDetailResponse
		glossaryCustomerType, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
			Table:     "sales_assignment_item",
			Attribute: "customer_type",
			ValueName: assignment.CustomerType,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRowInvalid("assignments", i, "customer_type")
			return
		}
		assignment.CustomerTypeValue = int8(glossaryCustomerType.Data.ValueInt)

	}

	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "SLA",
		Domain: "sales_assignment",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "sales_assignment")
		return
	}
	code := codeGenerator.Data.Code

	salesAssignment := &model.SalesAssignment{
		Code:          code,
		TerritoryIDGP: territory.Data[0].Salsterr,
		StartDate:     req.StartDate,
		EndDate:       req.EndDate,
		Status:        statusx.ConvertStatusName(statusx.Active),
	}

	err = s.RepositorySalesAssignment.Create(ctx, salesAssignment)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, assignment := range req.Assignments {
		salesAssignmentItem := &model.SalesAssignmentItem{
			SalesAssignmentID:     &salesAssignment.ID,
			SalesPersonIDGP:       assignment.SalespersonID,
			AddressIDGP:           assignment.AddressID,
			CustomerAcquisitionID: assignment.CustomerAcquisitionID,
			Task:                  assignment.TaskValue,
			CustomerType:          assignment.CustomerTypeValue,
			StartDate:             req.StartDate,
			EndDate:               req.EndDate,
			Status:                statusx.ConvertStatusName(statusx.Active),
			ObjectiveCodes:        assignment.ObjectiveCodes,
		}
		err = s.RepositorySalesAssignmentItem.Create(ctx, salesAssignmentItem)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}
	return
}

func (s *SalesAssignmentService) CancelBatch(ctx context.Context, id int64) (res dto.SalesAssignmentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentService.Delete")
	defer span.End()

	// validate SalesAssignment is exist
	var salesAssignmentOld *model.SalesAssignment
	salesAssignmentOld, err = s.RepositorySalesAssignment.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var salesAssignmentItemOld []*model.SalesAssignmentItem
	salesAssignmentItemOld, err = s.RepositorySalesAssignmentItem.GetBySalesAssignmentItemID(ctx, salesAssignmentOld.ID, 0, 0, time.Time{}, time.Time{})

	for _, salesAssignmentItem := range salesAssignmentItemOld {
		if salesAssignmentItem.Status == statusx.ConvertStatusName(statusx.Active) {
			err = s.RepositorySalesAssignmentItem.Cancel(ctx, salesAssignmentItem.ID)
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		}
	}

	// validate Status
	if salesAssignmentOld.Status != statusx.ConvertStatusName(statusx.Active) {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	err = s.RepositorySalesAssignment.Cancel(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesAssignmentService) CancelItem(ctx context.Context, id int64) (res dto.SalesAssignmentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentService.Delete")
	defer span.End()

	// validate SalesAssignment is exist
	var SalesAssignmentItemOld *model.SalesAssignmentItem
	SalesAssignmentItemOld, err = s.RepositorySalesAssignmentItem.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate Status
	if SalesAssignmentItemOld.Status != statusx.ConvertStatusName(statusx.Active) {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	err = s.RepositorySalesAssignmentItem.Cancel(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
