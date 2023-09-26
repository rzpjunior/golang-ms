package service

import (
	"context"
	"errors"
	"math"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/constants"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/repository"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ISalesAssignmentItemService interface {
	Get(ctx context.Context, offset, limit, status int, search, orderBy string, territoryID string, salesPersonID string, submitDateFrom, submitDateTo, startDateFrom, startDateTo, endDateFrom, endDateTo time.Time, task, outOfRoute int, customerType int32) (res []*dto.SalesAssignmentItemResponse, count int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.SalesAssignmentItemResponse, err error)
	CheckActiveTask(ctx context.Context, salesPersonID string) (exist bool, err error)
	SubmitTaskVisitFU(ctx context.Context, req dto.UpdateSubmitTaskVisitFURequest) (res *dto.SalesAssignmentItemResponse, err error)
	CheckoutTaskVisitFU(ctx context.Context, req dto.CheckoutTaskRequest) (err error)
	BulkCheckoutTaskVisitFU(ctx context.Context, req dto.BulkCheckoutTaskRequest) (err error)
	Create(ctx context.Context, req *dto.SalesAssignmentItemRequest) (res *dto.SalesAssignmentItemResponse, err error)
}

type SalesAssignmentItemService struct {
	opt                                opt.Options
	RepositorySalesAssignmentItem      repository.ISalesAssignmentItemRepository
	RepositorySalesAssignmentObjective repository.ISalesAssignmentObjectiveRepository
	RepositoryCustomerAcquisition      repository.ICustomerAcquisitionRepository
}

func NewSalesAssignmentItemService() ISalesAssignmentItemService {
	return &SalesAssignmentItemService{
		opt:                                global.Setup.Common,
		RepositorySalesAssignmentItem:      repository.NewSalesAssignmentItemRepository(),
		RepositorySalesAssignmentObjective: repository.NewSalesAssignmentObjectiveRepository(),
		RepositoryCustomerAcquisition:      repository.NewCustomerAcquisitionRepository(),
	}
}

func (s *SalesAssignmentItemService) Get(ctx context.Context, offset, limit, status int, search, orderBy string, territoryID string, salesPersonID string, submitDateFrom, submitDateTo, startDateFrom, startDateTo, endDateFrom, endDateTo time.Time, task, outOfRoute int, customerType int32) (res []*dto.SalesAssignmentItemResponse, count int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentItemService.Get")
	defer span.End()

	var salesAssignmentItems []*model.SalesAssignmentItem
	salesAssignmentItems, count, err = s.RepositorySalesAssignmentItem.Get(ctx, offset, limit, status, search, orderBy, territoryID, salesPersonID, submitDateFrom, submitDateTo, startDateFrom, startDateTo, endDateFrom, endDateTo, task, outOfRoute, customerType)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesAssignmentItem := range salesAssignmentItems {
		var (
			salesPerson *bridgeService.GetSalesPersonGPResponse
			addressObj  *dto.AddressResponse
			caObj       *dto.CustomerAcquisitionResponse
		)
		salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
			Id: salesAssignmentItem.SalesPersonIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "salesperson")
			return
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
			addressObj = &dto.AddressResponse{
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
			caObj = &dto.CustomerAcquisitionResponse{
				ID:   ca.ID,
				Code: ca.Code,
				Name: ca.Name,
			}
		}

		item := dto.SalesAssignmentItemResponse{
			ID:                salesAssignmentItem.ID,
			SalesAssignmentID: salesAssignmentItem.SalesAssignmentID,
			SalesPerson: &dto.SalespersonResponse{
				ID:   salesPerson.Data[0].Slprsnid,
				Code: salesPerson.Data[0].Slprsnid,
				Name: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
			},
			AddressID:             salesAssignmentItem.AddressIDGP,
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
			Address:               addressObj,
			CustomerAcquisition:   caObj,
		}

		if salesAssignmentItem.ObjectiveCodes != "" {
			var objectiveCodes []*model.SalesAssignmentObjective
			codes := strings.Split(salesAssignmentItem.ObjectiveCodes, ",")
			objectiveCodes, count, err = s.RepositorySalesAssignmentObjective.Get(ctx, 0, 10, 0, "", codes, "")
			if err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

			for _, val := range objectiveCodes {
				item.ObjectiveValues = append(item.ObjectiveValues, &dto.SalesAssignmentObjectiveResponse{
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
		res = append(res, &item)
	}
	return
}

func (s *SalesAssignmentItemService) GetByID(ctx context.Context, id int64) (res dto.SalesAssignmentItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentItemService.GetByID")
	defer span.End()

	var salesAssignmentItem *model.SalesAssignmentItem
	salesAssignmentItem, err = s.RepositorySalesAssignmentItem.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var (
		salesPerson *bridgeService.GetSalesPersonGPResponse
		addressObj  *dto.AddressResponse
		caObj       *dto.CustomerAcquisitionResponse
	)
	salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
		Id: salesAssignmentItem.SalesPersonIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "salesperson")
		return
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
		addressObj = &dto.AddressResponse{
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
		caObj = &dto.CustomerAcquisitionResponse{
			ID:   ca.ID,
			Code: ca.Code,
			Name: ca.Name,
		}
	}

	// get sales assignment item
	res = dto.SalesAssignmentItemResponse{
		ID:                salesAssignmentItem.ID,
		SalesAssignmentID: salesAssignmentItem.SalesAssignmentID,
		SalesPerson: &dto.SalespersonResponse{
			ID:   salesPerson.Data[0].Slprsnid,
			Code: salesPerson.Data[0].Slprsnid,
			Name: salesPerson.Data[0].Slprsnfn + " " + salesPerson.Data[0].Sprsnsmn + " " + salesPerson.Data[0].Sprsnsln,
		},
		AddressID:             salesAssignmentItem.AddressIDGP,
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
		Address:               addressObj,
		CustomerAcquisition:   caObj,
	}

	if salesAssignmentItem.ObjectiveCodes != "" {
		var objectiveCodes []*model.SalesAssignmentObjective
		objectiveCodes, _, err = s.RepositorySalesAssignmentObjective.Get(ctx, 0, 10, 0, "", strings.Split(salesAssignmentItem.ObjectiveCodes, ","), "")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		for _, val := range objectiveCodes {
			res.ObjectiveValues = append(res.ObjectiveValues, &dto.SalesAssignmentObjectiveResponse{
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

	return
}

func (s *SalesAssignmentItemService) CheckActiveTask(ctx context.Context, salesPersonID string) (exist bool, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentItemService.GetByID")
	defer span.End()

	exist, err = s.RepositorySalesAssignmentItem.GetSingleActiveTask(ctx, salesPersonID)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = nil
	return
}

func (s *SalesAssignmentItemService) SubmitTaskVisitFU(ctx context.Context, req dto.UpdateSubmitTaskVisitFURequest) (res *dto.SalesAssignmentItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentItemService.SubmitTaskVisitFU")
	defer span.End()

	// validate sales assignment object is exist
	var salesAssignmentItemOld *model.SalesAssignmentItem
	salesAssignmentItemOld, err = s.RepositorySalesAssignmentItem.GetByID(ctx, req.ID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if salesAssignmentItemOld.Status != statusx.ConvertStatusName(statusx.Active) {
		err = edenlabs.ErrorMustActive("status")
		return
	}

	if salesAssignmentItemOld.Task == 1 && (len(req.TaskImageUrls) > 7 && len(req.TaskImageUrls) < 1) {
		err = edenlabs.ErrorMustEqualOrLess("task_image_urls", "7")
		return
	} else if salesAssignmentItemOld.Task == 2 && (len(req.TaskImageUrls) > 3 && len(req.TaskImageUrls) < 1) {
		err = edenlabs.ErrorMustEqualOrLess("task_image_urls", "3")
		return
	}

	var (
		addressLatitude, addressLongitude float64
		address                           *bridgeService.GetAddressGPResponse
	)
	if req.CustomerType == 1 {
		address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
			Id: req.Address.ID,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "address")
			return
		}
		addressLatitude = address.Data[0].GnL_Latitude
		addressLongitude = address.Data[0].GnL_Longitude
	} else if req.CustomerType == 2 {
		_, err = s.RepositoryCustomerAcquisition.GetByID(ctx, req.CustomerAcquisitionID)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	// check first if branch has coordinate
	if addressLatitude != 0 && addressLongitude != 0 {
		// test comparing multiple coordinates to get distance
		var PI = 3.141592653589793
		var lat1 = req.Latitude
		var lng1 = req.Longitude
		var lat2 = &addressLatitude
		var lng2 = &addressLongitude

		radlat1 := float64(PI * lat1 / 180)
		radlat2 := float64(PI * *lat2 / 180)

		theta := float64(lng1 - *lng2)
		radtheta := float64(PI * theta / 180)

		dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

		if dist > 1 {
			dist = 1
		}

		dist = math.Acos(dist)
		dist = dist * 180 / PI
		dist = dist * 60 * 1.1515

		// convert into Ms because default is in miles
		dist = (dist * 1.609344) * 1000

		// get glossary from configuration for validation task answer
		var configApp *configurationService.GetConfigAppDetailResponse
		configApp, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppDetail(ctx, &configurationService.GetConfigAppDetailRequest{
			Application: 4,
			Attribute:   "max_distance_task",
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("configuration", "config app")
			return
		}

		maxDistance, _ := strconv.ParseFloat(configApp.Data.Value, 64)

		// if config app value == 0 don't validate
		if maxDistance > 0 {
			if dist > maxDistance {
				edenlabs.ErrorValidation("id", "Failed. You are exceed from allowed distance")
			}
		}
	}

	// get glossary from configuration for validation task answer
	_, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryList(ctx, &configurationService.GetGlossaryListRequest{
		Table:     "sales_assignment_item",
		Attribute: "task_answer",
		ValueInt:  int32(req.TaskAnswer),
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}

	salesAssignmentItem := &model.SalesAssignmentItem{
		ID:             req.ID,
		TaskImageUrl:   strings.Join(req.TaskImageUrls, ","),
		TaskAnswer:     req.TaskAnswer,
		Latitude:       req.Latitude,
		Longitude:      req.Longitude,
		ActualDistance: req.ActualDistance,
		SubmitDate:     time.Now(),
		OutOfRoute:     2,
		Status:         2,
	}

	err = s.RepositorySalesAssignmentItem.Update(ctx, salesAssignmentItem, "TaskImageUrl", "TaskAnswer", "Latitude", "Longitude", "ActualDistance", "SubmitDate", "OutOfRoute", "Status")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.SalesAssignmentItemResponse{
		ID:                    salesAssignmentItem.ID,
		SalesAssignmentID:     salesAssignmentItem.SalesAssignmentID,
		AddressID:             salesAssignmentItem.AddressIDGP,
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

	return
}

func (s *SalesAssignmentItemService) CheckoutTaskVisitFU(ctx context.Context, req dto.CheckoutTaskRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentItemService.SubmitTaskVisitFU")
	defer span.End()

	// validate sales assignment object is exist
	var (
		salesAssignmentItem *model.SalesAssignmentItem
		customerAcquisition *model.CustomerAcquisition
	)

	if !req.CustomerAcquisition {
		_, err = s.RepositorySalesAssignmentItem.GetByID(ctx, req.Id)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		timeNow := time.Now()
		salesAssignmentItem = &model.SalesAssignmentItem{
			ID:         req.Id,
			FinishDate: &timeNow,
		}

		err = s.RepositorySalesAssignmentItem.Update(ctx, salesAssignmentItem, "FinishDate")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	} else {
		_, err = s.RepositoryCustomerAcquisition.GetByID(ctx, req.Id)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		customerAcquisition = &model.CustomerAcquisition{
			ID:         req.Id,
			FinishDate: time.Now(),
		}

		err = s.RepositoryCustomerAcquisition.Update(ctx, customerAcquisition, "FinishDate")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	if req.Task < 1 && req.Task > 3 {
		err = edenlabs.ErrorInvalid("task")
	}

	return
}

func (s *SalesAssignmentItemService) BulkCheckoutTaskVisitFU(ctx context.Context, req dto.BulkCheckoutTaskRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentItemService.SubmitTaskVisitFU")
	defer span.End()

	var (
		salesAssignmentItems []*model.SalesAssignmentItem
		customerAcquisitions []*model.CustomerAcquisition
	)

	salesAssignmentItems, _, err = s.RepositorySalesAssignmentItem.GetMultiTaskActive(ctx, req.SalesPersonID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	customerAcquisitions, _, err = s.RepositoryCustomerAcquisition.GetMultiTaskActive(ctx, req.SalesPersonID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	timeNow := time.Now()
	for _, sai := range salesAssignmentItems {
		param := model.SalesAssignmentItem{
			ID:         sai.ID,
			FinishDate: &timeNow,
		}
		err = s.RepositorySalesAssignmentItem.Update(ctx, &param, "FinishDate")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
		var auditFunction string
		if sai.Task == 1 {
			auditFunction = "checkout task visit"
		} else if sai.Task == 2 {
			auditFunction = "checkout task follow up"
		}

		userID := ctx.Value(constants.KeyUserID).(int64)

		_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
			Log: &auditService.Log{
				UserId:      userID,
				ReferenceId: utils.ToString(sai.ID),
				Type:        "sales_assignment_iteme",
				Function:    auditFunction,
				CreatedAt:   timestamppb.New(time.Now()),
			},
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpc("audit")
			return
		}
	}

	for _, ca := range customerAcquisitions {
		param := model.CustomerAcquisition{
			ID:         ca.ID,
			FinishDate: time.Now(),
		}
		err = s.RepositoryCustomerAcquisition.Update(ctx, &param, "FinishDate")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		userID := ctx.Value(constants.KeyUserID).(int64)

		_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
			Log: &auditService.Log{
				UserId:      userID,
				ReferenceId: utils.ToString(ca.ID),
				Type:        "customer_acquisition",
				Function:    "checkout task customer acquisition",
				CreatedAt:   timestamppb.New(time.Now()),
			},
		})
	}

	return
}

func (s *SalesAssignmentItemService) Create(ctx context.Context, req *dto.SalesAssignmentItemRequest) (res *dto.SalesAssignmentItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentItemService.Create")
	defer span.End()

	var sai = model.SalesAssignmentItem{
		SalesAssignmentID:     req.SalesAssignmentID,
		SalesPersonIDGP:       req.SalesPersonID,
		AddressIDGP:           req.AddressID,
		CustomerAcquisitionID: req.CustomerAcquisitionID,
		Latitude:              req.Latitude,
		Longitude:             req.Longitude,
		Task:                  req.Task,
		CustomerType:          req.CustomerType,
		ObjectiveCodes:        req.ObjectiveCodes,
		ActualDistance:        req.ActualDistance,
		OutOfRoute:            req.OutOfRoute,
		StartDate:             req.StartDate,
		EndDate:               req.EndDate,
		SubmitDate:            req.SubmitDate,
		TaskImageUrl:          strings.Join(req.TaskImageUrls, ","),
		TaskAnswer:            req.TaskAnswer,
		Status:                req.Status,
		EffectiveCall:         req.EffectiveCall,
	}
	err = s.RepositorySalesAssignmentItem.Create(ctx, &sai)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.SalesAssignmentItemResponse{
		ID:                sai.ID,
		SalesAssignmentID: sai.SalesAssignmentID,
		SalesPerson: &dto.SalespersonResponse{
			ID: sai.SalesPersonIDGP,
		},
		AddressID:             sai.AddressIDGP,
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
		EffectiveCall:         sai.EffectiveCall,
	}

	return
}
