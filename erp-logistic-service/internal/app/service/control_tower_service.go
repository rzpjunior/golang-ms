package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"

	dto "git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/repository"
)

type IControlTowerService interface {
	GetDRS(ctx context.Context, req dto.ControlTowerGetDRSRequest) (res []dto.ControlTowerGetDRSResponse, total int64, err error)
	GetCourier(ctx context.Context, req dto.ControlTowerGetCourierRequest) (res dto.ControlTowerGetCourierResponse, err error)
	GetDRSDetail(ctx context.Context, id int64) (res dto.ControlTowerGetDRSDetailResponse, err error)
	GetCourierDetail(ctx context.Context, id int64) (res dto.ControlTowerGetCourierDetailResponse, err error)
	CancelDRS(ctx context.Context, req dto.ControlTowerCancelDRSRequest) (res dto.ControlTowerCancelDRSResponse, err error)
	CancelItem(ctx context.Context, req dto.ControlTowerCancelItemRequest) (res dto.ControlTowerCancelItemResponse, err error)
	Geocode(ctx context.Context, req *logistic_service.GeocodeAddressRequest) (res *logistic_service.GeocodeAddressResponse, err error)
}

type ControlTowerService struct {
	opt                             opt.Options
	RepositoryCourierLog            repository.ICourierLogRepository
	RepositoryDeliveryRunSheetItem  repository.IDeliveryRunSheetItemRepository
	RepositoryDeliveryRunSheet      repository.IDeliveryRunSheetRepository
	RepositoryDeliveryRunReturn     repository.IDeliveryRunReturnRepository
	RepositoryDeliveryRunReturnItem repository.IDeliveryRunReturnItemRepository
	RepositoryPostponeDeliveryLog   repository.IPostponeDeliveryLogRepository
	RepositoryAddressCoordinateLog  repository.IAddressCoordinateLogRepository
}

func NewServiceControlTower() IControlTowerService {
	return &ControlTowerService{
		opt:                             global.Setup.Common,
		RepositoryCourierLog:            repository.NewCourierLogRepository(),
		RepositoryDeliveryRunSheetItem:  repository.NewDeliveryRunSheetItemRepository(),
		RepositoryDeliveryRunSheet:      repository.NewDeliveryRunSheetRepository(),
		RepositoryDeliveryRunReturn:     repository.NewDeliveryRunReturnRepository(),
		RepositoryDeliveryRunReturnItem: repository.NewDeliveryRunReturnItemRepository(),
		RepositoryPostponeDeliveryLog:   repository.NewPostponeDeliveryLogRepository(),
		RepositoryAddressCoordinateLog:  repository.NewAddressCoordinateLogRepository(),
	}
}

func (s *ControlTowerService) GetDRS(ctx context.Context, req dto.ControlTowerGetDRSRequest) (res []dto.ControlTowerGetDRSResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ControlTowerService.GetDRS")
	defer span.End()

	// get sales orders from siteID and start delivery date until end delivery date
	var (
		arrSalesOrderIDs []string
		salesOrders      *bridgeService.GetSalesOrderGPListResponse
	)

	if salesOrders, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPAll(ctx, &bridgeService.GetSalesOrderGPListRequest{
		Locncode:        req.SiteID,
		Limit:           10000,
		Offset:          0,
		ReqShipDateFrom: req.StartDeliveryDate.Format("2006-01-02"),
		ReqShipDateTo:   req.EndDeliveryDate.Format("2006-01-02"),
	}); err != nil {
		fmt.Println("Errrrrrrrrrr", err)
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	for _, salesOrder := range salesOrders.Data {
		arrSalesOrderIDs = append(arrSalesOrderIDs, salesOrder.Sopnumbe)
	}
	// if no sales order in date range, return nil
	if len(arrSalesOrderIDs) == 0 {
		return
	}

	// OPTIONAL FIELDS
	// courier vendor parameter
	var arrCourierVendorCourierIDs []string
	var arrVehicleProfile *bridgeService.GetVehicleProfileGPResponse
	// get from vehicle profile
	arrVehicleProfile, err = s.opt.Client.BridgeServiceGrpc.GetVehicleProfileGPList(ctx, &bridgeService.GetVehicleProfileGPListRequest{
		Limit:              10000,
		Offset:             0,
		GnlCourierVendorId: req.CourierVendorID,
		Orderby:            "desc",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vehicle profile")
		return
	}

	// get courier that has associated vehicle profile
	for _, vehicleProfile := range arrVehicleProfile.Data {
		var courier *bridgeService.GetCourierGPResponse

		fmt.Println(vehicleProfile)
		courier, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPList(ctx, &bridgeService.GetCourierGPListRequest{
			Limit:              10000,
			Offset:             0,
			GnlCourierVendorId: req.CourierVendorID,
			Orderby:            "desc",
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		}
		fmt.Println(courier.Data[0])
		for _, courier := range courier.Data {
			arrCourierVendorCourierIDs = append(arrCourierVendorCourierIDs, courier.GnlCourierId)
		}
	}

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
	}

	// get grouped delivery run sheet item
	var deliveryRunSheetItems []*model.DeliveryRunSheetItem
	deliveryRunSheetItems, total, err = s.RepositoryDeliveryRunSheetItem.GetAllGroupedDeliveryRunSheetItem(ctx, dto.GroupedDeliveryRunSheetItemGetRequest{
		Offset:  req.Offset,
		Limit:   req.Limit,
		OrderBy: req.OrderBy,
		// Status:                      req.StatusIDs,
		ArrSalesOrderIDs:            arrSalesOrderIDs,
		ArrCourierVendorsCourierIDs: arrCourierVendorCourierIDs,
		CourierID:                   req.CourierID,
		SearchSalesOrderID:          req.Search,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("delivery_run_sheet_item", "Invalid delivery run sheet item")
		return
	}
	// prepare res
	for _, v := range deliveryRunSheetItems {
		// get total completed sales order and total sales order
		var totalFinished, totalSalesOrder int64

		_, totalFinished, err = s.RepositoryDeliveryRunSheetItem.Get(ctx, dto.DeliveryRunSheetItemGetRequest{
			StepType:            []int{2},
			Status:              []int{3},
			DeliveryRunSheetIDs: []int64{v.DeliveryRunSheetID},
			CourierIDs:          []string{v.CourierID},
		})
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("delivery_run_sheet_item", "Invalid delivery run sheet item")
			return
		}

		_, totalSalesOrder, err = s.RepositoryDeliveryRunSheetItem.Get(ctx, dto.DeliveryRunSheetItemGetRequest{
			StepType:            []int{2},
			Status:              []int{},
			DeliveryRunSheetIDs: []int64{v.DeliveryRunSheetID},
			CourierIDs:          []string{v.CourierID},
		})
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("delivery_run_sheet_item", "Invalid delivery run sheet item")
			return
		}

		// get delivery run sheet detail
		var deliveryRunSheetDetail *model.DeliveryRunSheet
		deliveryRunSheetDetail, err = s.RepositoryDeliveryRunSheet.GetByID(ctx, v.DeliveryRunSheetID, "", req.StatusIDs)

		// get courier detail
		var courierDetail *bridgeService.GetCourierGPResponse
		courierDetail, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPDetail(ctx, &bridgeService.GetCourierGPDetailRequest{
			Id: v.CourierID,
		})
		if err != nil || !courierDetail.Succeeded {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "courier")
			return
		}
		if deliveryRunSheetDetail != nil {
			res = append(res, dto.ControlTowerGetDRSResponse{
				CompletedSalesOrder: totalFinished,
				TotalSalesOrder:     totalSalesOrder,
				Courier: &dto.SubControlTowerGetDRSCourier{
					Code: courierDetail.Data[0].GnlCourierId,
					Name: courierDetail.Data[0].GnlCourierName,
				},
				DeliveryRunSheet: &dto.SubControlTowerGetDRSDeliveryRunSheet{
					ID:           deliveryRunSheetDetail.ID,
					Code:         deliveryRunSheetDetail.Code,
					DeliveryDate: deliveryRunSheetDetail.DeliveryDate,
					Status:       deliveryRunSheetDetail.Status,
				},
			})
		}

	}

	return
}

func (s *ControlTowerService) GetCourier(ctx context.Context, req dto.ControlTowerGetCourierRequest) (res dto.ControlTowerGetCourierResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ControlTowerService.GetCourier")
	defer span.End()

	var (
		vehicleProfile        *bridgeService.GetVehicleProfileGPResponse
		requestVehicleProfile bridgeService.GetVehicleProfileGPListRequest
		requestCourierList    bridgeService.GetCourierGPListRequest
		courier               *bridgeService.GetCourierGPResponse
		mapVehicleProfile     map[string]*bridgeService.VehicleProfileGP
		arrCourierID          []string
		courierLogs           []*dto.CourierLog
		mapCourierLogs        map[string]*dto.CourierLog
		countEmergencyMode    int64
	)

	// If the SiteID is empty, return without further processing
	if req.SiteID == "" {
		return
	}

	// Set up request parameters for retrieving vehicle profiles
	requestVehicleProfile.Locncode = req.SiteID
	requestVehicleProfile.Limit = 10000
	requestVehicleProfile.Offset = 0

	// Set up request parameters for retrieving couriers
	requestCourierList.Limit = 10000
	requestCourierList.Offset = 0

	// Apply filters for CourierVendorID and CourierID if provided
	if req.CourierVendorID != "" {
		requestVehicleProfile.GnlCourierVendorId = req.CourierVendorID
		requestCourierList.GnlCourierVendorId = req.CourierVendorID
	}

	if req.CourierID != "" {
		requestCourierList.GnlCourierId = req.CourierID
	}

	// Retrieve vehicle profiles based on the request parameters
	vehicleProfile, err = s.opt.Client.BridgeServiceGrpc.GetVehicleProfileGPList(ctx, &requestVehicleProfile)
	if vehicleProfile != nil && vehicleProfile.Data != nil {
		mapVehicleProfile = make(map[string]*bridgeService.VehicleProfileGP)
		for _, vpData := range vehicleProfile.Data {
			mapVehicleProfile[vpData.GnlVehicleProfileId] = vpData
		}
	}

	// Retrieve couriers based on the request parameters
	courier, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPList(ctx, &requestCourierList)

	// Collect the IDs of the retrieved couriers
	for _, courierData := range courier.Data {
		arrCourierID = append(arrCourierID, courierData.GnlCourierId)
	}

	// Retrieve courier location logs for the collected courier IDs
	courierLogs, err = s.RepositoryCourierLog.GetAllCourierLocation(ctx, arrCourierID)
	mapCourierLogs = make(map[string]*dto.CourierLog)
	for _, courierLog := range courierLogs {
		mapCourierLogs[courierLog.CourierID] = courierLog
	}

	// Process each courier and populate the response
	for _, courierData := range courier.Data {
		tempRes := dto.SubControlTowerGetCourierCourier{
			ID:            courierData.GnlCourierId,
			Name:          courierData.GnlCourierName,
			PhoneNumber:   courierData.Phonname,
			LicensePlate:  courierData.GnlLicensePlate,
			EmergencyMode: int8(*courierData.GnlEmergencymode),
		}

		if mapCourierLogs[courierData.GnlCourierId] != nil && mapCourierLogs[courierData.GnlCourierId].CourierID == courierData.GnlCourierId {
			tempRes.Latitude = *mapCourierLogs[courierData.GnlCourierId].Latitude
			tempRes.Longitude = *mapCourierLogs[courierData.GnlCourierId].Longitude
			tempRes.LastUpdated = mapCourierLogs[courierData.GnlCourierId].CreatedAt
		}

		// Check if vehicle profile information is available and populate it
		if mapVehicleProfile[courierData.GnlVehicleProfileId] != nil && mapVehicleProfile[courierData.GnlVehicleProfileId].GnlVehicleProfileId == courierData.GnlVehicleProfileId {
			tempRes.VehicleProfileType = mapVehicleProfile[courierData.GnlVehicleProfileId].GnlRoutingProfile
			tempRes.VendorCourierCode = mapVehicleProfile[courierData.GnlVehicleProfileId].GnlCourierVendorId
		}

		res.Couriers = append(res.Couriers, &tempRes)
		// Calculate the number of couriers with emergency status
		if tempRes.EmergencyMode == 1 {
			countEmergencyMode++
		}
	}

	res.OnEmergency = countEmergencyMode

	return
}

func (s *ControlTowerService) GetDRSDetail(ctx context.Context, id int64) (res dto.ControlTowerGetDRSDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ControlTowerService.GetDRSDetail")
	defer span.End()

	// get delivery run sheet
	var deliveryRunSheet *model.DeliveryRunSheet
	deliveryRunSheet, err = s.RepositoryDeliveryRunSheet.GetByID(ctx, id, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// get courier detail
	var courierDetail *bridgeService.GetCourierGPResponse
	courierDetail, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPDetail(ctx, &bridgeService.GetCourierGPDetailRequest{
		Id: deliveryRunSheet.CourierID,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}

	// get courier vehicle profile
	var vehicleProfileDetail *bridgeService.GetVehicleProfileGPResponse
	vehicleProfileDetail, err = s.opt.Client.BridgeServiceGrpc.GetVehicleProfileGPDetail(ctx, &bridgeService.GetVehicleProfileGPDetailRequest{
		Id: courierDetail.Data[0].GnlVehicleProfileId,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vehicle profile")
		return
	}

	// get courier courier vendor
	var courierVendorDetail *bridgeService.GetCourierVendorGPResponse
	courierVendorDetail, err = s.opt.Client.BridgeServiceGrpc.GetCourierVendorGPDetail(ctx, &bridgeService.GetCourierVendorGPDetailRequest{
		Id: vehicleProfileDetail.Data[0].GnlCourierVendorId,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier vendor")
		return
	}

	// get courier site
	var siteDetail *bridgeService.GetSiteGPResponse
	siteDetail, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: courierVendorDetail.Data[0].Locncode,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	// prepare res
	res = dto.ControlTowerGetDRSDetailResponse{
		Id:           deliveryRunSheet.ID,
		Code:         deliveryRunSheet.Code,
		DeliveryDate: deliveryRunSheet.DeliveryDate,
		Courier: &dto.SubControlTowerGetDRSDetailCourier{
			Code:                 courierDetail.Data[0].GnlCourierId,
			Name:                 courierDetail.Data[0].GnlCourierName,
			CourierPhoneNumber:   courierDetail.Data[0].Phonname,
			CourierVendorName:    courierVendorDetail.Data[0].GnlCourierVendorName,
			CourierVehicleName:   vehicleProfileDetail.Data[0].GnlDescription100,
			LicensePlate:         courierDetail.Data[0].GnlLicensePlate,
			CourierWarehouseName: siteDetail.Data[0].Locndscr,
		},
		StartedAt:         deliveryRunSheet.StartedAt,
		FinishedAt:        deliveryRunSheet.FinishedAt,
		StartingLatitude:  deliveryRunSheet.StartingLatitude,
		StartingLongitude: deliveryRunSheet.StartingLongitude,
		FinishedLatitude:  deliveryRunSheet.FinishedLatitude,
		FinishedLongitude: deliveryRunSheet.FinishedLongitude,
		Status:            deliveryRunSheet.Status,
	}

	return
}

func (s *ControlTowerService) GetCourierDetail(ctx context.Context, id int64) (res dto.ControlTowerGetCourierDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ControlTowerService.GetCourierDetail")
	defer span.End()

	// get delivery run sheet
	var deliveryRunSheet *model.DeliveryRunSheet
	deliveryRunSheet, err = s.RepositoryDeliveryRunSheet.GetByID(ctx, id, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// get courier log
	var courierLog *dto.CourierLog
	courierLog, err = s.RepositoryCourierLog.GetLastCourierLog(ctx, &logistic_service.GetLastCourierLogRequest{
		CourierId: deliveryRunSheet.CourierID,
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// prepare courier detail res
	// get courier detail
	var courierDetail *bridgeService.GetCourierGPResponse
	courierDetail, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPDetail(ctx, &bridgeService.GetCourierGPDetailRequest{
		Id: courierLog.CourierID,
	})
	emergencyTime, err := time.Parse("2006-01-02 15:04:05", courierDetail.Data[0].GnlEmergencydate+" "+courierDetail.Data[0].GnlEmergencytime)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}
	// get courier vehicle profile
	var vehicleProfileDetail *bridgeService.GetVehicleProfileGPResponse
	vehicleProfileDetail, err = s.opt.Client.BridgeServiceGrpc.GetVehicleProfileGPDetail(ctx, &bridgeService.GetVehicleProfileGPDetailRequest{
		Id: courierDetail.Data[0].GnlVehicleProfileId,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vehicle profile")
		return
	}

	// get all drsi
	var deliveryRunSheetItems []*model.DeliveryRunSheetItem
	deliveryRunSheetItems, _, err = s.RepositoryDeliveryRunSheetItem.Get(ctx, dto.DeliveryRunSheetItemGetRequest{
		StepType:            []int{2},
		DeliveryRunSheetIDs: []int64{deliveryRunSheet.ID},
	})
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("delivery_run_sheet_item", "Invalid delivery run sheet item")
		return
	}

	var resTotalSelfPickup int64
	var resTotalDeliveryReturn int64
	var deliveryRunSheetItemResponse []*dto.SubControlTowerGetCourierDetailDRSI

	// prepare delivery run sheet item res
	for _, deliveryRunSheetItem := range deliveryRunSheetItems {
		var (
			salesOrder        *bridgeService.GetSalesOrderGPListResponse
			salesInvoice      *bridgeService.GetSalesInvoiceGPListResponse
			addressDetail     *bridgeService.GetAddressGPResponse
			admDivisionDetail *bridgeService.GetAdmDivisionGPResponse
			// subDistrictDetail *bridgeService.GetSubDistrictDetailResponse
			wrtDetail *bridgeService.GetWrtGPResponse
			// termPayment *bridgeService.GetPaymentTermDetailResponse
			// salesInvoiceDetail *bridgeService.GetSalesInvoiceDetailResponse
			// customerLatitude  *float64
			// customerLongitude *float64

			deliveryRunReturn *model.DeliveryRunReturn

			respDRRI               []*dto.SubControlTowerGetCourierDetailDRRI
			deliveryRunReturnItems []*model.DeliveryRunReturnItem
			// deliveryOrderItem *bridgeService.GetDeliveryOrderItemResponse
			productDetail        *bridgeService.GetItemGPResponse
			uomDetail            *bridgeService.GetUomGPResponse
			deliveryRunReturnRes *dto.SubControlTowerGetCourierDetailDRR

			respPostponeDeliveryLog []*dto.SubControlTowerGetCourierDetailPostponeDeliveryLog
			postponeDeliveryLogs    []*model.PostponeDeliveryLog
			salesInvoiceExist       *dto.SalesInvoice
			glossary                *configuration_service.GetGlossaryDetailResponse
			glossaryTemp            *configuration_service.Glossary
			unpunctualReason        string
		)

		// get drsi sales order
		if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
			Id: deliveryRunSheetItem.SalesOrderID,
		}); err != nil || salesOrder.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
			return
		}

		salesInvoice, _ = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
			SoNumber: salesOrder.Data[0].Sopnumbe,
		})

		if salesInvoice != nil && salesInvoice.Data != nil {
			salesInvoiceExist = &dto.SalesInvoice{
				TotalCharge: salesInvoice.Data[0].Ordocamt,
			}
		}

		// get sales order's address
		addressDetail, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
			Id: salesOrder.Data[0].Address[0].Prstadcd,
		})
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "address")
			return
		}

		if wrtDetail, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
			Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
			GnlRegion: salesOrder.Data[0].GnL_Region,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
			return
		}

		// get address's adm division
		admDivisionDetail, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
			AdmDivisionCode: addressDetail.Data[0].GnL_Administrative_Code,
		})
		if err != nil || admDivisionDetail == nil || admDivisionDetail.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
			return
		}

		// // get adm division's sub district
		// subDistrictDetail, err = s.opt.Client.BridgeServiceGrpc.GetSubDistrictDetail(ctx, &bridgeService.GetSubDistrictDetailRequest{
		// 	Id: admDivisionDetail.Data[0].District,
		// })
		// if err != nil {
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorRpcNotFound("bridge", "sub district")
		// 	return
		// }

		// get sales order's term payment
		// get sales order's sales invoice

		// TODO CHANGE GEOCODE CODE TO ACCEPT THE RIGHT DATAMODEL
		// geocode customer location
		// if location, err = s.Geocode(ctx, addressDetail); err != nil {
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorValidation("location", "Invalid location")
		// 	return
		// }

		// get drsi delivery run return
		deliveryRunReturn, err = s.RepositoryDeliveryRunReturn.GetByID(ctx, 0, "", deliveryRunSheetItem.ID)
		// if err != nil {
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorValidation("delivery_run_return", "Invalid delivery run return")
		// 	return
		// }
		// get drr delivery run return item

		if deliveryRunReturn != nil && deliveryRunReturn.ID != 0 {

			deliveryRunReturnRes = &dto.SubControlTowerGetCourierDetailDRR{
				TotalPrice:  deliveryRunReturn.TotalPrice,
				TotalCharge: deliveryRunReturn.TotalCharge,
			}
			deliveryRunReturnItems, _, err = s.RepositoryDeliveryRunReturnItem.Get(ctx, dto.DeliveryRunReturnItemGetRequest{
				ArrDeliveryRunReturnIDs: []int64{deliveryRunReturn.ID},
			})
			if err != nil {
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorValidation("delivery_run_return_item", "Invalid delivery run return item")
				return
			}
			// response SubControlTowerGetCourierDetailDRRI
			for _, deliveryRunReturnItem := range deliveryRunReturnItems {
				for _, v := range salesInvoice.Data[0].Details {
					if v.Itemnmbr == deliveryRunReturnItems[0].DeliveryOrderItemID {

						if deliveryRunReturnItem.ReturnReason == 0 {
							// this for handing value int
							// we dont need to remove omitempty while getting glossary data
							glossaryTemp = &configuration_service.Glossary{
								ValueName: "",
							}
							glossary = &configuration_service.GetGlossaryDetailResponse{
								Data: glossaryTemp,
							}
						} else {
							// get configuration for delivery return item
							if glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
								Table:     "delivery_run_return_item",
								Attribute: "item_return_reason",
								ValueInt:  int32(deliveryRunReturnItem.ReturnReason),
							}); err != nil {
								span.RecordError(err)
								s.opt.Logger.AddMessage(log.ErrorLevel, err)
								err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
								return
							}
						}

						// get item
						productDetail, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
							Id: v.Itemnmbr,
						})
						if err != nil {
							span.RecordError(err)
							s.opt.Logger.AddMessage(log.ErrorLevel, err)
							err = edenlabs.ErrorRpcNotFound("bridge", "item")
							return
						}

						// get item
						uomDetail, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
							Id: v.Uofm,
						})
						if err != nil {
							span.RecordError(err)
							s.opt.Logger.AddMessage(log.ErrorLevel, err)
							err = edenlabs.ErrorRpcNotFound("bridge", "item")
							return
						}

						respDRRI = append(respDRRI, &dto.SubControlTowerGetCourierDetailDRRI{
							DeliveryQty:       v.Quantity,
							ReceiveQty:        deliveryRunReturnItem.ReceiveQty,
							Subtotal:          deliveryRunReturnItem.Subtotal,
							ReturnReason:      deliveryRunReturnItem.ReturnReason,
							ReturnReasonValue: glossary.Data.ValueName,
							ReturnEvidence:    deliveryRunReturnItem.ReturnEvidence,
							ProductName:       productDetail.Data[0].Itemdesc,
							UOMName:           uomDetail.Data[0].Uofm,
						})
					}
				}
			}
			deliveryRunReturnRes.DeliveryRunReturnItem = respDRRI
			resTotalDeliveryReturn++
		}
		if deliveryRunSheetItem.UnpunctualDetail == 1 {
			if glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
				Table:     "delivery_run_sheet_item",
				Attribute: "early_delivery_reason",
				ValueInt:  int32(deliveryRunSheetItem.UnpunctualReason),
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}

		}

		if deliveryRunSheetItem.UnpunctualDetail == 2 {
			if glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configuration_service.GetGlossaryDetailRequest{
				Table:     "delivery_run_sheet_item",
				Attribute: "late_delivery_reason",
				ValueInt:  int32(deliveryRunSheetItem.UnpunctualReason),
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				return
			}
		}
		if glossary != nil && glossary.Data != nil {
			unpunctualReason = glossary.Data.ValueName
		}
		// get postpone delivery log
		postponeDeliveryLogs, err = s.RepositoryPostponeDeliveryLog.Get(ctx, deliveryRunSheetItem.ID)
		// if err != nil {
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorValidation("delivery_run_return", "Invalid delivery run return")
		// 	return
		// }

		if len(postponeDeliveryLogs) > 0 {
			// response postpone delivery log
			for _, postponeDeliveryLog := range postponeDeliveryLogs {
				respPostponeDeliveryLog = append(respPostponeDeliveryLog, &dto.SubControlTowerGetCourierDetailPostponeDeliveryLog{
					PostponeReason:   postponeDeliveryLog.PostponeReason,
					PostponeEvidence: postponeDeliveryLog.PostponeEvidence,
					StartedAt:        time.Unix(postponeDeliveryLog.StartedAtUnix, 0),
					PostponedAt:      time.Unix(postponeDeliveryLog.PostponedAtUnix, 0),
				})
			}

		}
		// append to deliveryResponse
		deliveryRunSheetItemResponse = append(deliveryRunSheetItemResponse, &dto.SubControlTowerGetCourierDetailDRSI{
			Id:                          deliveryRunSheetItem.ID,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               deliveryRunSheetItem.MoneyReceived,
			StartTime:                   deliveryRunSheetItem.StartedAt,
			ArrivalTime:                 deliveryRunSheetItem.ArrivalTime,
			FinishTime:                  deliveryRunSheetItem.FinishedAt,
			UnpunctualDetail:            deliveryRunSheetItem.UnpunctualDetail,
			UnpunctualReason:            deliveryRunSheetItem.UnpunctualReason,
			UnpunctualReasonValue:       unpunctualReason,
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			DeliveryEvidenceImageURL:    deliveryRunSheetItem.DeliveryEvidenceImageURL,
			TransactionEvidenceImageURL: deliveryRunSheetItem.TransactionEvidenceImageURL,
			Note:                        deliveryRunSheetItem.Note,
			Status:                      deliveryRunSheetItem.Status,
			SalesOrder: &dto.SubControlTowerGetCourierDetailSO{
				Code:                  salesOrder.Data[0].Sopnumbe,
				DeliveryDate:          salesOrder.Data[0].ReqShipDate,
				DeliveryFee:           salesOrder.Data[0].Frtamnt,
				VoucherDiscountAmount: salesOrder.Data[0].Trdisamt,
				PointRedeemAmount:     0,                                       // POINT REDEEM AMOUNT BELUM ADA
				CustomerName:          salesOrder.Data[0].Customer[0].Custname, // CONFIRM LAGI CUSTOMER NAME AMBIL DARI ADDRESS DETAIL
				// CustomerLatitude:        *customerLatitude,                        // DARI GEOCODE
				// CustomerLongitude:       *customerLongitude,                       // DARI GEOCODE
				AddressName:        addressDetail.Data[0].ShipToName,
				AddressPhoneNumber: addressDetail.Data[0].PhonE1, // PHONE NUMBER CONFIRM LAGI
				ShippingAddress:    addressDetail.Data[0].AddresS1 + addressDetail.Data[0].AddresS2 + addressDetail.Data[0].AddresS3,
				SubDistrictDetail:  admDivisionDetail.Data[0].Subdistrict,
				PostalCode:         admDivisionDetail.Data[0].Zipcode,
				WrtName:            wrtDetail.Data[0].Strttime + "-" + wrtDetail.Data[0].Endtime,
				PaymentTypeName:    salesOrder.Data[0].Pymtrmid,
				SalesInvoice:       salesInvoiceExist,
			},
			DeliveryRunReturn: deliveryRunReturnRes,
			PostponeLog:       respPostponeDeliveryLog,
		})
	}

	res = dto.ControlTowerGetCourierDetailResponse{
		TotalSalesOrder:     int64(len(deliveryRunSheetItems)),
		TotalSelfPickup:     resTotalSelfPickup,
		TotalDeliveryReturn: resTotalDeliveryReturn,
		Courier: &dto.SubControlTowerGetCourierDetailCourier{
			Latitude:           *courierLog.Latitude,
			Longitude:          *courierLog.Longitude,
			EmergencyMode:      int8(*courierDetail.Data[0].GnlEmergencymode),
			LastEmergencyTime:  emergencyTime,
			LastUpdated:        emergencyTime,
			VehicleProfileType: vehicleProfileDetail.Data[0].GnlRoutingProfile,
		},
		DeliveryRunSheetItem: deliveryRunSheetItemResponse,
	}
	return
}

func (s *ControlTowerService) CancelDRS(ctx context.Context, req dto.ControlTowerCancelDRSRequest) (res dto.ControlTowerCancelDRSResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ControlTowerService.CancelDRS")
	defer span.End()

	var (
		deliveryRunSheet              *model.DeliveryRunSheet
		deliveryRunSheetItemsDelivery []*model.DeliveryRunSheetItem
		deliveryRunSheetItemsPickup   []*model.DeliveryRunSheetItem
		salesOrderIDs                 []string
	)

	deliveryRunSheet, err = s.RepositoryDeliveryRunSheet.GetByID(ctx, req.DeliveryRunSheetID, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorNotFound("delivery run sheet")
		return
	}

	if deliveryRunSheet.Status != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorMustActive("status")
		return
	}

	deliveryRunSheetItemsDelivery, _, err = s.RepositoryDeliveryRunSheetItem.Get(ctx, dto.DeliveryRunSheetItemGetRequest{
		StepType:            []int{2},
		Status:              []int{1, 2, 4},
		DeliveryRunSheetIDs: []int64{deliveryRunSheet.ID},
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorNotFound("delivery run sheet item")
		return
	}

	for _, drsi := range deliveryRunSheetItemsDelivery {
		salesOrderIDs = append(salesOrderIDs, drsi.SalesOrderID)
	}

	if len(salesOrderIDs) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorNotFound("delivery run sheet item")
		return
	}

	deliveryRunSheetItemsPickup, _, err = s.RepositoryDeliveryRunSheetItem.Get(ctx, dto.DeliveryRunSheetItemGetRequest{
		Offset:              0,
		Limit:               1000,
		Status:              []int{1, 2, 4},
		DeliveryRunSheetIDs: []int64{req.DeliveryRunSheetID},
		StepType:            []int{2},
		// CourierIDs:          []string{},
		// ArrSalesOrderIDs:    salesOrderIDs,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorNotFound("delivery run sheet item")
		return
	}

	// update delivery run sheet
	updatedDRS := &model.DeliveryRunSheet{
		ID:         deliveryRunSheet.ID,
		Status:     3,
		FinishedAt: time.Now(),
	}
	err = s.RepositoryDeliveryRunSheet.Update(ctx, updatedDRS, "status", "finished_at")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("internal_error", "error updating drs")
		return
	}

	// update delivery run sheet items
	for _, drsiPickup := range deliveryRunSheetItemsPickup {
		updatedDRSI := &model.DeliveryRunSheetItem{
			ID:         drsiPickup.ID,
			Status:     5,
			Note:       req.Note,
			FinishedAt: time.Now(),
		}

		err = s.RepositoryDeliveryRunSheetItem.Update(ctx, updatedDRSI, "status", "note", "finished_at")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("internal_error", "error updating drsi")
			return
		}
	}

	for _, drsiDelivery := range deliveryRunSheetItemsDelivery {
		updatedDRSI := &model.DeliveryRunSheetItem{
			ID:         drsiDelivery.ID,
			Status:     5,
			Note:       req.Note,
			FinishedAt: time.Now(),
		}

		err = s.RepositoryDeliveryRunSheetItem.Update(ctx, updatedDRSI, "status", "note", "finished_at")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("internal_error", "error updating drsi")
			return
		}
	}

	res = dto.ControlTowerCancelDRSResponse{
		ID:     updatedDRS.ID,
		Code:   updatedDRS.Code,
		Status: updatedDRS.Status,
	}

	return
}

func (s *ControlTowerService) CancelItem(ctx context.Context, req dto.ControlTowerCancelItemRequest) (res dto.ControlTowerCancelItemResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ControlTowerService.CancelItem")
	defer span.End()

	var (
		deliveryRunSheetItemDelivery *model.DeliveryRunSheetItem
		deliveryRunSheetItemPickup   *model.DeliveryRunSheetItem
		finishedDeliveryRunSheet     bool
		// finishedDeliveryOrder        bool
		// deliveryOrders *bridgeService.GetDeliveryOrderListResponse
	)

	deliveryRunSheetItemDelivery, err = s.RepositoryDeliveryRunSheetItem.GetByID(ctx, req.DeliveryRunSheetItemID)
	if err != nil {
		err = edenlabs.ErrorNotFound("delivery run sheet item")
		return
	}

	if deliveryRunSheetItemDelivery.StepType != 2 {
		err = edenlabs.ErrorInvalid("step type")
		return
	}

	if deliveryRunSheetItemDelivery.Status == 3 || deliveryRunSheetItemDelivery.Status == 5 {
		err = edenlabs.ErrorMustActive("status")
		return
	}

	deliveryRunSheetItemPickup, err = s.RepositoryDeliveryRunSheetItem.GetByID(ctx, deliveryRunSheetItemDelivery.ID-1)
	if err != nil {
		err = edenlabs.ErrorNotFound("delivery run sheet item")
		return
	}

	// // check if the so has delivery order with status of active / finished / on progress(from postponed) / failed
	// filter = map[string]interface{}{"sales_order_id": r.DeliveryRunSheetItem.SalesOrder.ID, "status__in": []int64{1, 2, 5, 7}}
	// deliveryOrders, countDeliveryOrder, err := repository.CheckDeliveryOrder(filter, exclude)
	// if err != nil || countDeliveryOrder == 0 {
	// 	o.Failure("delivery_order.invalid", util.ErrorInvalidData("delivery order"))
	// 	return o
	// }
	// // if DO status finished (canceling for bug SO), this will not change the DO status
	// if r.DeliveryOrder.Status == 2 {
	// 	r.FinishedDeliveryOrder = true
	// }

	// check if the errand is the last delivery run sheet item
	_, countDRSI, err := s.RepositoryDeliveryRunSheetItem.Get(ctx, dto.DeliveryRunSheetItemGetRequest{
		Status:              []int{1, 2, 4},
		DeliveryRunSheetIDs: []int64{deliveryRunSheetItemDelivery.DeliveryRunSheetID},
	})
	if err != nil {
		err = edenlabs.ErrorNotFound("delivery run sheet item")
		return
	}
	if countDRSI == 1 {
		finishedDeliveryRunSheet = true
	}

	// update delivery run sheet items
	updatedPickupDRSI := &model.DeliveryRunSheetItem{
		ID:         deliveryRunSheetItemPickup.ID,
		Status:     5,
		Note:       req.Note,
		FinishedAt: time.Now(),
	}
	err = s.RepositoryDeliveryRunSheetItem.Update(ctx, updatedPickupDRSI, "status", "note", "finished_at")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("internal_error", "error updating drsi")
		return
	}

	updatedDeliveryDRSI := &model.DeliveryRunSheetItem{
		ID:         deliveryRunSheetItemDelivery.ID,
		Status:     5,
		Note:       req.Note,
		FinishedAt: time.Now(),
	}
	err = s.RepositoryDeliveryRunSheetItem.Update(ctx, updatedDeliveryDRSI, "status", "note", "finished_at")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("internal_error", "error updating drsi")
		return
	}

	// update delivery run sheet
	if finishedDeliveryRunSheet {
		updatedDRS := &model.DeliveryRunSheet{
			ID:         deliveryRunSheetItemDelivery.DeliveryRunSheetID,
			Status:     3,
			FinishedAt: time.Now(),
		}
		err = s.RepositoryDeliveryRunSheet.Update(ctx, updatedDRS, "status", "finished_at")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("internal_error", "error updating drs")
			return
		}
	}

	// update delivery orders

	res = dto.ControlTowerCancelItemResponse{
		ID:     updatedDeliveryDRSI.ID,
		Status: updatedDeliveryDRSI.Status,
	}

	return
}

// Geocode
func (s *ControlTowerService) Geocode(ctx context.Context, req *logistic_service.GeocodeAddressRequest) (res *logistic_service.GeocodeAddressResponse, err error) {
	res = &logistic_service.GeocodeAddressResponse{}

	// check Address coordinate log
	var addressCoordinateLog *model.AddressCoordinateLog
	addressCoordinateLog, _ = s.RepositoryAddressCoordinateLog.GetMostTrusted(ctx, req.AddressId)
	if addressCoordinateLog != nil {
		res.Latitude = &addressCoordinateLog.Latitude
		res.Longitude = &addressCoordinateLog.Longitude
		return
	}

	// geocode with address
	var (
		geoCodeData            dto.GoogleGeocode
		topLevelDomainCode     = "region=id&"
		addressString          = "address="
		subDistrictInformation = (req.SubDistrict +
			req.City +
			req.Region +
			req.Zip)
		re           = regexp.MustCompile(`[^0-9a-zA-Z]+`)
		addressArray = re.ReplaceAllString(req.AddressName, " ")
		client       = &http.Client{}
	)

	addressSplit := strings.Split(addressArray, " ")
	for _, v := range addressSplit {
		addressString += v + "+"
	}

	request, err := http.NewRequest("GET", s.opt.Config.Google.GeocodeURL+topLevelDomainCode+addressString+subDistrictInformation+s.opt.Config.Google.MapKey, nil)
	if err != nil {
		return
	}

	response, err := client.Do(request)
	if err != nil {
		return
	}

	defer response.Body.Close()

	if err = json.NewDecoder(response.Body).Decode(&geoCodeData.Data); err != nil {
		return
	}

	switch geoCodeData.Data.Status {
	case "OK":

		// create address coordinate log
		if err = s.RepositoryAddressCoordinateLog.Create(ctx, &model.AddressCoordinateLog{
			AddressID:      req.AddressId,
			SalesOrderID:   req.SalesOrderId,
			Latitude:       geoCodeData.Data.Results[0].Geometry.Location.Lat,
			Longitude:      geoCodeData.Data.Results[0].Geometry.Location.Lng,
			LogChannelID:   6,
			MainCoordinate: 0,
			CreatedAt:      time.Now(),
		}); err != nil {
			return
		}

		res.Latitude = &geoCodeData.Data.Results[0].Geometry.Location.Lat
		res.Longitude = &geoCodeData.Data.Results[0].Geometry.Location.Lng
		return

	default: // if fail using address, then geocode only using the sub district information

		var geoCodeData dto.GoogleGeocode
		request, err = http.NewRequest("GET", s.opt.Config.Google.GeocodeURL+topLevelDomainCode+"address="+subDistrictInformation+s.opt.Config.Google.MapKey, nil)
		if err != nil {
			return
		}

		response, err = client.Do(request)
		if err != nil {
			return
		}

		defer response.Body.Close()

		err = json.NewDecoder(response.Body).Decode(&geoCodeData.Data)
		if err != nil {
			return
		}

		if geoCodeData.Data.Status != "OK" {
			return
		}

		// create address coordinate log
		if err = s.RepositoryAddressCoordinateLog.Create(ctx, &model.AddressCoordinateLog{
			AddressID:      req.AddressId,
			SalesOrderID:   req.SalesOrderId,
			Latitude:       geoCodeData.Data.Results[0].Geometry.Location.Lat,
			Longitude:      geoCodeData.Data.Results[0].Geometry.Location.Lng,
			LogChannelID:   6,
			MainCoordinate: 0,
			CreatedAt:      time.Now(),
		}); err != nil {
			return
		}

		res.Latitude = &geoCodeData.Data.Results[0].Geometry.Location.Lat
		res.Longitude = &geoCodeData.Data.Results[0].Geometry.Location.Lng
	}

	return
}
