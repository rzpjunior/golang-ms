package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-crm-service/internal/app/repository"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	catalogService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/catalog_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
)

type ICustomerAcquisitionService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (res []*dto.CustomerAcquisitionResponse, total int64, err error)
	GetSubmissions(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (res []*dto.CustomerAcquisitionResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.CustomerAcquisitionResponse, err error)
	CheckActiveTask(ctx context.Context, salesPersonID string) (exist bool, err error)
	SubmitTask(ctx context.Context, req dto.SubmitTaskCustomerAcqRequest) (res *dto.CustomerAcquisitionResponse, err error)
	GetWithExcludedIds(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time, excludedIds []int64) (res []*dto.CustomerAcquisitionResponse, total int64, err error)
	CountCustomerAcq(ctx context.Context, salespersonID int64, submitDateFrom time.Time, submitDateTo time.Time) (total int64, err error)
}

type CustomerAcquisitionService struct {
	opt                               opt.Options
	RepositoryCustomerAcquisition     repository.ICustomerAcquisitionRepository
	RepositoryCustomerAcquisitionItem repository.ICustomerAcquisitionItemRepository
	RepoSalesAssignmentItem           repository.ISalesAssignmentItemRepository
}

func NewCustomerAcquisitionService() ICustomerAcquisitionService {
	return &CustomerAcquisitionService{
		opt:                               global.Setup.Common,
		RepositoryCustomerAcquisition:     repository.NewCustomerAcquisitionRepository(),
		RepositoryCustomerAcquisitionItem: repository.NewCustomerAcquisitionItemRepository(),
		RepoSalesAssignmentItem:           repository.NewSalesAssignmentItemRepository(),
	}
}

func (s *CustomerAcquisitionService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (res []*dto.CustomerAcquisitionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerAcquisitionService.Get")
	defer span.End()

	var customerAcquisitions []*model.CustomerAcquisition
	customerAcquisitions, total, err = s.RepositoryCustomerAcquisition.Get(ctx, offset, limit, status, search, orderBy, territoryID, salespersonID, submitDateFrom, submitDateTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, customerAcquisition := range customerAcquisitions {

		var territory *bridgeService.GetSalesTerritoryGPResponse
		territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: customerAcquisition.TerritoryIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "territory")
			return
		}

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

		res = append(res, &dto.CustomerAcquisitionResponse{
			ID:               customerAcquisition.ID,
			Task:             customerAcquisition.Task,
			Name:             customerAcquisition.Name,
			PhoneNumber:      customerAcquisition.PhoneNumber,
			Latitude:         customerAcquisition.Latitude,
			Longitude:        customerAcquisition.Longitude,
			AddressName:      customerAcquisition.AddressName,
			FoodApp:          customerAcquisition.FoodApp,
			PotentialRevenue: customerAcquisition.PotentialRevenue,
			TaskImageUrl:     customerAcquisition.TaskImageUrl,
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
			FinishDate:    customerAcquisition.FinishDate,
			SubmitDate:    customerAcquisition.SubmitDate,
			CreatedAt:     customerAcquisition.CreatedAt,
			UpdatedAt:     customerAcquisition.UpdatedAt,
			Status:        customerAcquisition.Status,
			StatusConvert: statusx.ConvertStatusValue(customerAcquisition.Status),
		})
	}

	return
}

func (s *CustomerAcquisitionService) GetSubmissions(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time) (res []*dto.CustomerAcquisitionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerAcquisitionService.Get")
	defer span.End()

	var customerAcquisitions []*model.CustomerAcquisition
	customerAcquisitions, total, err = s.RepositoryCustomerAcquisition.GetSubmissions(ctx, offset, limit, status, search, orderBy, territoryID, salespersonID, submitDateFrom, submitDateTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, customerAcquisition := range customerAcquisitions {

		var territory *bridgeService.GetSalesTerritoryGPResponse
		territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: customerAcquisition.TerritoryIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "territory")
			return
		}

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
		res = append(res, &dto.CustomerAcquisitionResponse{
			ID:               customerAcquisition.ID,
			Task:             customerAcquisition.Task,
			Name:             customerAcquisition.Name,
			PhoneNumber:      customerAcquisition.PhoneNumber,
			Latitude:         customerAcquisition.Latitude,
			Longitude:        customerAcquisition.Longitude,
			AddressName:      customerAcquisition.AddressName,
			FoodApp:          customerAcquisition.FoodApp,
			PotentialRevenue: customerAcquisition.PotentialRevenue,
			TaskImageUrl:     customerAcquisition.TaskImageUrl,
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
			FinishDate:    customerAcquisition.FinishDate,
			SubmitDate:    customerAcquisition.SubmitDate,
			CreatedAt:     customerAcquisition.CreatedAt,
			UpdatedAt:     customerAcquisition.UpdatedAt,
			Status:        customerAcquisition.Status,
			StatusConvert: statusx.ConvertStatusValue(customerAcquisition.Status),
		})
	}

	return
}

func (s *CustomerAcquisitionService) GetByID(ctx context.Context, id int64) (res dto.CustomerAcquisitionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerAcquisitionService.GetByID")
	defer span.End()

	var customerAcquisition *model.CustomerAcquisition
	customerAcquisition, err = s.RepositoryCustomerAcquisition.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// get customer acquisition item
	var customerAcquisitionItemsResponse []*dto.CustomerAcquisitionItemResponse
	var customerAcquisitionItems []*model.CustomerAcquisitionItem
	customerAcquisitionItems, _, _ = s.RepositoryCustomerAcquisitionItem.GetByCustomerAcquisitionID(ctx, customerAcquisition.ID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, customerAcquisitionItem := range customerAcquisitionItems {
		var item *bridgeService.GetItemGPResponse
		item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: customerAcquisitionItem.ItemIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("item_id")
			err = edenlabs.ErrorRpcNotFound("bridge", "item")
			return
		}

		customerAcquisitionItemsResponse = append(customerAcquisitionItemsResponse, &dto.CustomerAcquisitionItemResponse{
			ID:                    customerAcquisitionItem.ID,
			CustomerAcquisitionID: customerAcquisitionItem.CustomerAcquisitionID,
			Item: &dto.ItemResponse{
				ID:          item.Data[0].Itemnmbr,
				Code:        item.Data[0].Itemnmbr,
				Description: item.Data[0].Itmgedsc,
			},
			IsTop:     customerAcquisitionItem.IsTop,
			CreatedAt: customerAcquisitionItem.CreatedAt,
			UpdatedAt: customerAcquisitionItem.UpdatedAt,
		})
	}

	var territory *bridgeService.GetSalesTerritoryGPResponse
	territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
		Id: customerAcquisition.TerritoryIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "territory")
		return
	}

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

	res = dto.CustomerAcquisitionResponse{
		ID:               customerAcquisition.ID,
		Task:             customerAcquisition.Task,
		Name:             customerAcquisition.Name,
		PhoneNumber:      customerAcquisition.PhoneNumber,
		Latitude:         customerAcquisition.Latitude,
		Longitude:        customerAcquisition.Longitude,
		AddressName:      customerAcquisition.AddressName,
		FoodApp:          customerAcquisition.FoodApp,
		PotentialRevenue: customerAcquisition.PotentialRevenue,
		TaskImageUrl:     customerAcquisition.TaskImageUrl,
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
		FinishDate:               customerAcquisition.FinishDate,
		SubmitDate:               customerAcquisition.SubmitDate,
		CreatedAt:                customerAcquisition.CreatedAt,
		UpdatedAt:                customerAcquisition.UpdatedAt,
		Status:                   customerAcquisition.Status,
		StatusConvert:            statusx.ConvertStatusValue(customerAcquisition.Status),
		CustomerAcquisitionItems: customerAcquisitionItemsResponse,
	}

	return
}

func (s *CustomerAcquisitionService) CheckActiveTask(ctx context.Context, salesPersonID string) (exist bool, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesAssignmentItemService.CheckActiveTask")
	defer span.End()

	exist, err = s.RepositoryCustomerAcquisition.GetSingleActiveTask(ctx, salesPersonID)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = nil
	return
}

func (s *CustomerAcquisitionService) SubmitTask(ctx context.Context, req dto.SubmitTaskCustomerAcqRequest) (res *dto.CustomerAcquisitionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerAcquisitionService.SubmitTask")
	defer span.End()

	var exist bool
	exist, err = s.RepositoryCustomerAcquisition.GetSingleActiveTask(ctx, req.SalesPersonID)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	if exist {
		err = edenlabs.ErrorExists("task")
	}

	exist, err = s.RepoSalesAssignmentItem.GetSingleActiveTask(ctx, req.SalesPersonID)
	if err != nil && !errors.Is(err, orm.ErrNoRows) {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	if exist {
		err = edenlabs.ErrorExists("task")
		return
	}

	if len(req.CustomerName) > 20 {
		err = edenlabs.ErrorMustEqualOrLess("customer name", "20")
		return
	}

	if len(req.PhoneNumber) < 8 {
		err = edenlabs.ErrorMustEqualOrGreater("phone number", "8")
		return
	}

	if len(req.PhoneNumber) > 15 {
		err = edenlabs.ErrorMustEqualOrLess("phone number", "15")
		return
	}

	if len(req.AddressDetail) > 250 {
		err = edenlabs.ErrorMustEqualOrLess("phone number", "250")
		return
	}

	customerAcqPhotoArr := strings.Split(req.CustomerAcquisitionPhoto, ",")

	if len(customerAcqPhotoArr) > 3 {
		err = edenlabs.ErrorMustEqualOrLess("phone number", "3")
		return
	}

	var products []*model.CustomerAcquisitionItem
	for _, product := range req.Product {
		// valid item/product
		_, err = s.opt.Client.CatalogServiceGrpc.GetItemDetailByInternalId(ctx, &catalogService.GetItemDetailByInternalIdRequest{
			Id: utils.ToString(product.Id),
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("catalog", "item detail")
			return
		}
		products = append(products, &model.CustomerAcquisitionItem{
			ItemID:    product.Id,
			IsTop:     product.Top,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		})
	}

	var codeGenerator *configurationService.GetGenerateCodeResponse
	codeGenerator, err = s.opt.Client.ConfigurationServiceGrpc.GetGenerateCode(ctx, &configurationService.GetGenerateCodeRequest{
		Format: "CA",
		Domain: "customer_acquisition",
		Length: 6,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "code_generator")
		return
	}
	code := codeGenerator.Data.Code

	ca := &model.CustomerAcquisition{
		Code:             code,
		Task:             3,
		Name:             req.CustomerName,
		PhoneNumber:      req.PhoneNumber,
		Latitude:         req.UserLatitude,
		Longitude:        req.UserLongitude,
		AddressName:      req.AddressDetail,
		FoodApp:          req.FoodApp,
		PotentialRevenue: req.PotentialRevenue,
		TaskImageUrl:     req.CustomerAcquisitionPhoto,
		SalespersonIDGP:  req.SalesPersonID,
		// TerritoryID:      salesPerson.Data[0].Salsterr,
		SubmitDate: time.Now(),
		UpdatedAt:  time.Now(),
		Status:     statusx.ConvertStatusName(statusx.Active),
		CreatedAt:  time.Now(),
	}

	caId, err := s.RepositoryCustomerAcquisition.CreateWithItem(ctx, ca, products)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	caItems, _, err := s.RepositoryCustomerAcquisitionItem.GetByCustomerAcquisitionID(ctx, caId)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var items []*dto.CustomerAcquisitionItemResponse
	for _, cai := range caItems {
		items = append(items, &dto.CustomerAcquisitionItemResponse{
			ID:                    cai.ID,
			CustomerAcquisitionID: cai.CustomerAcquisitionID,
			Item: &dto.ItemResponse{
				ID: cai.ItemIDGP,
			},
			IsTop:     cai.IsTop,
			CreatedAt: cai.CreatedAt,
			UpdatedAt: cai.UpdatedAt,
		})
	}

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

	var salesPerson *bridgeService.GetSalesPersonGPResponse
	salesPerson, err = s.opt.Client.BridgeServiceGrpc.GetSalesPersonGPDetail(ctx, &bridgeService.GetSalesPersonGPDetailRequest{
		Id: ca.SalespersonIDGP,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "salesperson")
		return
	}

	res = &dto.CustomerAcquisitionResponse{
		ID:               ca.ID,
		Code:             ca.Code,
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
		StatusConvert:            statusx.ConvertStatusValue(ca.Status),
		CreatedAt:                timex.ToLocTime(ctx, ca.CreatedAt),
		UpdatedAt:                timex.ToLocTime(ctx, ca.UpdatedAt),
		CustomerAcquisitionItems: items,
	}

	return
}

func (s *CustomerAcquisitionService) GetWithExcludedIds(ctx context.Context, offset int, limit int, status int, search string, orderBy string, territoryID string, salespersonID string, submitDateFrom time.Time, submitDateTo time.Time, excludedIds []int64) (res []*dto.CustomerAcquisitionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerAcquisitionService.GetWithExcludedIds")
	defer span.End()

	var customerAcquisitions []*model.CustomerAcquisition
	customerAcquisitions, total, err = s.RepositoryCustomerAcquisition.GetWithExcludedIds(ctx, offset, limit, status, search, orderBy, territoryID, salespersonID, submitDateFrom, submitDateTo, excludedIds)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, customerAcquisition := range customerAcquisitions {

		var territory *bridgeService.GetSalesTerritoryGPResponse
		territory, err = s.opt.Client.BridgeServiceGrpc.GetSalesTerritoryGPDetail(ctx, &bridgeService.GetSalesTerritoryGPDetailRequest{
			Id: customerAcquisition.TerritoryIDGP,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "territory")
			return
		}

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

		res = append(res, &dto.CustomerAcquisitionResponse{
			ID:               customerAcquisition.ID,
			Task:             customerAcquisition.Task,
			Name:             customerAcquisition.Name,
			PhoneNumber:      customerAcquisition.PhoneNumber,
			Latitude:         customerAcquisition.Latitude,
			Longitude:        customerAcquisition.Longitude,
			AddressName:      customerAcquisition.AddressName,
			FoodApp:          customerAcquisition.FoodApp,
			PotentialRevenue: customerAcquisition.PotentialRevenue,
			TaskImageUrl:     customerAcquisition.TaskImageUrl,
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
			FinishDate:    customerAcquisition.FinishDate,
			SubmitDate:    customerAcquisition.SubmitDate,
			CreatedAt:     customerAcquisition.CreatedAt,
			UpdatedAt:     customerAcquisition.UpdatedAt,
			Status:        customerAcquisition.Status,
			StatusConvert: statusx.ConvertStatusValue(customerAcquisition.Status),
		})
	}

	return
}

func (s *CustomerAcquisitionService) CountCustomerAcq(ctx context.Context, salespersonID int64, submitDateFrom time.Time, submitDateTo time.Time) (total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerAcquisitionService.CountCustomerAcq")
	defer span.End()

	if salespersonID == 0 {
		err = edenlabs.ErrorInvalid("salesperson_id")
		return
	}

	total, err = s.RepositoryCustomerAcquisition.CountCustomerAcquisition(ctx, salespersonID, submitDateFrom, submitDateTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
