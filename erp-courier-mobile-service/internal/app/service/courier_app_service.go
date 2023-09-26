package service

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/jwt"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/edenlabs/edenlabs/utils"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-courier-mobile-service/internal/app/dto"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	configurationService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	siteService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/site_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ICourierAppService interface {
	Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error)
	CreateCourierLog(ctx context.Context, req dto.CreateCourierLogRequest) (err error)
	Get(ctx context.Context, req dto.CourierAppGetRequest) (res []dto.CourierAppGetResponse, total int64, err error)
	Detail(ctx context.Context, id int64, courierID string) (res dto.CourierAppDetailResponse, err error)
	ScanDetail(ctx context.Context, req dto.CourierAppScanDetailRequest) (res dto.CourierAppDetailResponse, err error)
	Scan(ctx context.Context, req dto.CourierAppScanRequest) (res dto.CourierAppScanResponse, err error)
	SelfAssign(ctx context.Context, req dto.CourierAppSelfAssignRequest) (res dto.CourierAppSelfAssignResponse, err error)
	StartDelivery(ctx context.Context, req dto.CourierAppStartDeliveryRequest) (res dto.CourierAppStartDeliveryResponse, err error)
	SuccessDelivery(ctx context.Context, req dto.CourierAppSuccessDeliveryRequest) (res dto.CourierAppSuccessDeliveryResponse, err error)
	PostponeDelivery(ctx context.Context, req dto.CourierAppPostponeDeliveryRequest) (res dto.CourierAppPostponeDeliveryResponse, err error)
	FailDelivery(ctx context.Context, req dto.CourierAppFailDeliveryRequest) (res dto.CourierAppFailDeliveryResponse, err error)
	StatusDelivery(ctx context.Context, req dto.CourierAppStatusDeliveryRequest) (res dto.CourierAppStatusDeliveryResponse, err error)
	ActivateEmergency(ctx context.Context, req dto.CourierAppActivateEmergencyRequest) (res dto.CourierAppActivateEmergencyResponse, err error)
	DeactivateEmergency(ctx context.Context, req dto.CourierAppDeactivateEmergencyRequest) (res dto.CourierAppDeactivateEmergencyResponse, err error)
	CreateMerchantDeliveryLog(ctx context.Context, req dto.CourierAppCreateMerchantDeliveryLogRequest) (res dto.CourierAppCreateMerchantDeliveryLogResponse, err error)
	GetGlossary(ctx context.Context, req dto.CourierAppGetGlossaryRequest) (res []dto.CourierAppGetGlossaryResponse, total int64, err error)
}

type CourierAppService struct {
	opt opt.Options
}

func NewServiceCourierApp() ICourierAppService {
	return &CourierAppService{
		opt: global.Setup.Common,
	}
}

func (s *CourierAppService) Login(ctx context.Context, req dto.LoginRequest) (res dto.LoginResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.Login")
	defer span.End()
	var (
		courier        *bridgeService.GetCourierGPResponse
		vehicleProfile *bridgeService.GetVehicleProfileGPResponse
		courierVendor  *bridgeService.GetCourierVendorGPResponse
		site           *bridgeService.GetSiteGPResponse
	)

	// TODO : Login by endpoint login user in API req 1.5, User , Authentication for User
	// get courier by id need search by user id
	if courier, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPDetail(ctx, &bridgeService.GetCourierGPDetailRequest{
		Id: req.CourierCode,
	}); err != nil || !courier.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}

	// check password
	if courier.Data[0].Password != req.Password {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}

	if vehicleProfile, err = s.opt.Client.BridgeServiceGrpc.GetVehicleProfileGPDetail(ctx, &bridgeService.GetVehicleProfileGPDetailRequest{
		Id: courier.Data[0].GnlVehicleProfileId,
	}); err != nil || !vehicleProfile.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vehicle profile")
		return
	}

	if courierVendor, err = s.opt.Client.BridgeServiceGrpc.GetCourierVendorGPDetail(ctx, &bridgeService.GetCourierVendorGPDetailRequest{
		Id: vehicleProfile.Data[0].GnlCourierVendorId,
	}); err != nil || !courierVendor.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier vendor")
		return
	}

	if site, err = s.opt.Client.BridgeServiceGrpc.GetSiteGPDetail(ctx, &bridgeService.GetSiteGPDetailRequest{
		Id: courierVendor.Data[0].Locncode,
	}); err != nil || !site.Succeeded {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "site")
		return
	}

	jwtInit := jwt.NewJWT([]byte(s.opt.Config.Jwt.Key))
	uc := jwt.UserCourierClaim{
		CourierID: courier.Data[0].GnlCourierId,
		SiteID:    site.Data[0].Locncode,
		Timezone:  req.Timezone,
	}

	jwtGenerate, err := jwtInit.Create(uc)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.LoginResponse{
		Courier: &dto.GlobalCourier{
			Id:           courier.Data[0].GnlCourierId,
			Name:         courier.Data[0].GnlCourierName,
			PhoneNumber:  courier.Data[0].Phonname,
			LicensePlate: courier.Data[0].GnlLicensePlate,
			// EmergencyMode: courier.Data[0].EmergencyMode,
			Status: courier.Data[0].Inactive,
			VehicleProfile: &dto.GlobalVehicleProfile{
				RoutingProfile: vehicleProfile.Data[0].GnlRoutingProfile,
			},
		},
		Token: jwtGenerate,
	}

	return
}

func (s *CourierAppService) CreateCourierLog(ctx context.Context, req dto.CreateCourierLogRequest) (err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.CreateCourierLog")
	defer span.End()

	// var deliveryRunSheetItems *logisticService.GetDeliveryRunSheetItemListResponse

	// if deliveryRunSheetItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
	// 	Status:    []int32{2},
	// 	CourierId: []string{req.CourierID},
	// }); err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
	// 	return
	// }

	// if len(deliveryRunSheetItems.Data) > 0 {
	// 	if _, err = s.opt.Client.LogisticServiceGrpc.CreateCourierLog(ctx, &logisticService.CreateCourierLogRequest{
	// 		Model: &logisticService.CourierLog{
	// 			CourierId:    req.CourierID,
	// 			SalesOrderId: deliveryRunSheetItems.Data[0].SalesOrderId,
	// 			Latitude:     &req.Latitude,
	// 			Longitude:    &req.Longitude,
	// 			CreatedAt:    timestamppb.Now(),
	// 		},
	// 	}); err != nil {
	// 		span.RecordError(err)
	// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 		err = edenlabs.ErrorRpcNotFound("logistic", "courier log")
	// 		return
	// 	}
	// } else {
	if _, err = s.opt.Client.LogisticServiceGrpc.CreateCourierLog(ctx, &logisticService.CreateCourierLogRequest{
		Model: &logisticService.CourierLog{
			CourierId: req.CourierID,
			Latitude:  &req.Latitude,
			Longitude: &req.Longitude,
			CreatedAt: timestamppb.Now(),
		},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "courier log")
		return
	}
	// }

	return
}

func (s *CourierAppService) Get(ctx context.Context, req dto.CourierAppGetRequest) (res []dto.CourierAppGetResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.Get")
	defer span.End()

	var (
		deliveryRunSheetItems *logisticService.GetDeliveryRunSheetItemListResponse
		statusInt             []int32
		stepTypeInt           []int32
		SalesOrderID          []string
		salesOrders           *bridgeService.GetSalesOrderGPListResponse
	)

	for _, s := range req.StatusIDs {
		statusInt = append(statusInt, int32(s))
	}

	if req.StepType != 0 {
		stepTypeInt = append(stepTypeInt, int32(req.StepType))
	}

	if req.Search != "" {
		// Get the current time
		now := time.Now()
		// Format the dates directly using the desired layout
		twoDaysAheadStr := now.AddDate(0, 0, 2).Format("2006-01-02")     // 2 days ahead from time now
		threeDaysBehindStr := now.AddDate(0, 0, -3).Format("2006-01-02") // 3 days behind from time now
		salesOrders, _ = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPAll(ctx, &bridgeService.GetSalesOrderGPListRequest{
			Custname:        req.Search,
			ReqShipDateFrom: threeDaysBehindStr,
			ReqShipDateTo:   twoDaysAheadStr,
			Limit:           200,
			Offset:          0,
		})

		if salesOrders != nil && len(salesOrders.Data) > 0 {
			for _, so := range salesOrders.Data {
				SalesOrderID = append(SalesOrderID, so.Sopnumbe)
			}
		}
	}

	if deliveryRunSheetItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
		Limit:        int32(req.Limit),
		Offset:       int32(req.Offset),
		Status:       statusInt,
		OrderBy:      req.OrderBy,
		StepType:     stepTypeInt,
		CourierId:    []string{req.CourierId},
		SalesOrderId: SalesOrderID,
		// Search:       req.SearchSalesOrderCode,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	for _, drsi := range deliveryRunSheetItems.Data {
		var (
			salesOrder       *bridgeService.GetSalesOrderGPListResponse
			customer         *bridgeService.GetCustomerGPResponse
			salesPaymentTerm *bridgeService.GetSalesPaymentTermDetailResponse
			orderType        *bridgeService.GetOrderTypeDetailResponse
			wrt              *bridgeService.GetWrtGPResponse
			address          *bridgeService.GetAddressGPResponse
			courierLog       *logisticService.GetLastCourierLogResponse
			distances        float64
			// customerLatitude  *float64
			// customerLongitude *float64
		)

		if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
			Id: drsi.SalesOrderId,
		}); err != nil || salesOrder.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales order in service courier mobile")
			return
		}

		// get customer
		if customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{
			Id: salesOrder.Data[0].Customer[0].Custnmbr,
		}); err != nil || customer.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "customer")
			return
		}

		// sales payment term
		if salesPaymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentTermDetail(ctx, &bridgeService.GetSalesPaymentTermDetailRequest{
			// Id: customer.Data[0].Pymtrmid,
			Id: 1,
		}); err != nil || salesPaymentTerm.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales payment term")
			return
		}

		// get order type
		if orderType, err = s.opt.Client.BridgeServiceGrpc.GetOrderTypeDetail(ctx, &bridgeService.GetOrderTypeDetailRequest{
			// Id: salesOrder.Data.OrderTypeId,
			Id: 1,
		}); err != nil || orderType.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "order type")
			return
		}

		if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
			Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
			GnlRegion: salesOrder.Data[0].GnL_Region,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
			return
		}

		address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
			Id: salesOrder.Data[0].Address[0].Prstadcd,
		})

		// TODO geocode customer location then get the distance from customer's location to
		// courier latest location from courier log

		// // geocode
		// if customerLatitude, customerLongitude, err = s.geocode(ctx, address); err != nil {
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorValidation("location.invalid", "geocode")
		// 	return
		// }
		// check distance with customer

		courierLog, _ = s.opt.Client.LogisticServiceGrpc.GetLastCourierLog(ctx, &logisticService.GetLastCourierLogRequest{
			CourierId: req.CourierId,
		})

		if courierLog != nil && address.Data[0].GnL_Latitude != 0 && address.Data[0].GnL_Longitude != 0 {
			distances = distance(address.Data[0].GnL_Latitude, address.Data[0].GnL_Longitude, courierLog.Latitude, courierLog.Longitude)
		}

		res = append(res, dto.CourierAppGetResponse{
			DeliveryRunSheetItem: &dto.GlobalDeliveryRunSheetItem{
				Id:       drsi.Id,
				StepType: int8(drsi.StepType),
				Status:   int8(drsi.Status),
				Note:     drsi.Note,
				SalesOrder: &dto.GlobalSalesOrder{
					Code:      salesOrder.Data[0].Sopnumbe,
					OrderDate: salesOrder.Data[0].Docdate,
					Address: &dto.GlobalAddress{
						AddressName: address.Data[0].Custname,
						AdmDivision: &dto.GlobalAdmDivsion{
							SubDistrict: &dto.GlobalSubDistrict{
								Description: address.Data[0].AdministrativeDiv.GnlSubdistrict,
							},
						},
					},
					Customer: &dto.GlobalCustomer{
						Name: customer.Data[0].Custname,
						SalesPaymentTerm: &dto.GlobalSalesPaymentTerm{
							Description: salesPaymentTerm.Data.Description,
						},
					},
					Wrt: &dto.GlobalWrt{
						StartTime: wrt.Data[0].Strttime,
						EndTime:   wrt.Data[0].Endtime,
					},
					OrderType: &dto.GlobalOrderType{
						Description: orderType.Data.Description,
					},
				},
			},
			Distance: distances,
		})
	}

	return
}

func (s *CourierAppService) Detail(ctx context.Context, id int64, courierID string) (res dto.CourierAppDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.Detail")
	defer span.End()

	var (
		deliveryRunSheetItem *logisticService.GetDeliveryRunSheetItemDetailResponse

		salesOrder *bridgeService.GetSalesOrderGPListResponse
		customer   *bridgeService.GetCustomerGPResponse
		// salesPaymentTerm *bridgeService.GetSalesPaymentTermDetailResponse
		courierLog     *logisticService.GetLastCourierLogResponse
		address        *bridgeService.GetAddressGPResponse
		admDivision    *bridgeService.GetAdmDivisionGPResponse
		wrt            *bridgeService.GetWrtGPResponse
		orderType      *bridgeService.GetOrderTypeDetailResponse
		salesInvoice   *bridgeService.GetSalesInvoiceGPListResponse
		items          []*dto.GlobalSalesOrderItem
		deliveryOrders *bridgeService.GetDeliveryOrderGPListResponse
		salesInvoices  *bridgeService.GetSalesInvoiceGPListResponse

		deliveryRunReturn      *logisticService.GetDeliveryRunReturnDetailResponse
		deliveryRunReturnItems *logisticService.GetDeliveryRunReturnItemListResponse
		respDRRI               []*dto.GlobalDeliveryRunReturnItem
		dRRITemp               *dto.GlobalDeliveryRunReturnItem
		glossary               *configurationService.GetGlossaryDetailResponse
		distances              float64
	)

	// get delivery run sheet item
	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: id,
	}); err != nil || deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// get sales order
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: deliveryRunSheetItem.Data.SalesOrderId,
	}); err != nil || salesOrder.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// get customer
	if customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{
		Id: salesOrder.Data[0].Customer[0].Custnmbr,
	}); err != nil || customer.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	// TODO CHANGE DUMMY SALES PAYMENT TERM
	// sales payment term
	// if salesPaymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetPaymentTermGPDetail(ctx, &bridgeService.GetPaymentTermGPDetailRequest{
	// 	Id: customer.Data[0].Pymtrmid[0].Pymtrmid,
	// }); err != nil || salesPaymentTerm.Data == nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "sales payment term")
	// 	return
	// }
	// TODO CHANGE WRT DUMMY
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
		GnlRegion: salesOrder.Data[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	// TODO CHANGE ADDRESS DUMMY
	if address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
		Id: salesOrder.Data[0].Address[0].Prstadcd,
	}); err != nil || address.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		return
	}

	// TODO CHANGE ADM DIVSIION DUMMY
	// get address's adm division
	admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
		AdmDivisionCode: address.Data[0].GnL_Administrative_Code,
	})
	if err != nil || admDivision.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
		return
	}

	// TODO CHANGE ORDER TYPE DUMMY
	// get order type
	if orderType, err = s.opt.Client.BridgeServiceGrpc.GetOrderTypeDetail(ctx, &bridgeService.GetOrderTypeDetailRequest{
		// Id: salesOrder.Data.OrderTypeId,
		Id: 1,
	}); err != nil || orderType.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "order type")
		return
	}

	// prepare res items
	for _, soi := range salesOrder.Data[0].Details {
		var (
			item *bridgeService.GetItemGPResponse
			uom  *bridgeService.GetUomGPResponse
			// deliveryOrderItem *bridgeService.GetDeliveryOrderItemResponse
		)

		if item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: soi.Itemnmbr,
		}); err != nil || item.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item")
			return
		}

		if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
			Id: item.Data[0].Uomschdl,
		}); err != nil || uom.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			return
		}

		// TODO get delivery order item
		// if deliveryOrderItem, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderItemDetail(ctx, &bridgeService.GetDeliveryOrderItemRequest{
		// 	salesOrderItemId: soi.Id,
		// }); err != nil || deliveryOrderItem.Data == nil {
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorRpcNotFound("bridge", "uom")
		// 	return
		// }

		courierLog, err = s.opt.Client.LogisticServiceGrpc.GetLastCourierLog(ctx, &logisticService.GetLastCourierLogRequest{
			CourierId: courierID,
		})

		if courierLog != nil && address.Data[0].GnL_Latitude != 0 && address.Data[0].GnL_Longitude != 0 {
			distances = distance(address.Data[0].GnL_Latitude, address.Data[0].GnL_Longitude, courierLog.Latitude, courierLog.Longitude)
		}

		deliverQuantity := 9999999.0
		items = append(items, &dto.GlobalSalesOrderItem{
			OrderQty:  &soi.Quantity,
			UnitPrice: &soi.Unitprce,
			Subtotal:  &soi.Xtndprce,
			Item: &dto.GlobalItem{
				Description: item.Data[0].Itemdesc,
				Uom: &dto.GlobalUom{
					Description: uom.Data[0].Umschdsc,
				},
			},
			DeliveryOrderItem: &dto.GlobalDeliveryOrderItem{
				DeliverQty: &deliverQuantity,
			},
		})
	}

	// get drsi delivery run return
	deliveryRunReturn, _ = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnDetail(ctx, &logisticService.GetDeliveryRunReturnDetailRequest{
		DeliveryRunSheetItemId: deliveryRunSheetItem.Data.Id,
	})
	if salesInvoice, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		Limit:    1,
		SoNumber: salesOrder.Data[0].Sopnumbe,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice list")
		return
	}

	if deliveryOrders, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderListGP(ctx, &bridgeService.GetDeliveryOrderGPListRequest{
		Limit:    1,
		Offset:   1,
		SopNumbe: salesOrder.Data[0].Sopnumbe,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "delivery order")
		return
	}

	if salesInvoices, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		SoNumber: salesOrder.Data[0].Sopnumbe,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	if len(deliveryOrders.Data) == 0 && len(salesInvoices.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("courier", "sales order ini belum memiliki delivery order/ sales invoice")
		return
	}

	var salesOrderItems []*dto.GlobalSalesOrderItem

	for _, soi := range salesOrder.Data[0].Details {
		for _, sii := range salesInvoices.Data[0].Details {
			if soi.Itemnmbr == sii.Itemnmbr {
				salesOrderItems = append(salesOrderItems, &dto.GlobalSalesOrderItem{
					ItemNumber:  soi.Itemnmbr,
					ItemDesc:    soi.Itemdesc,
					OrderQty:    &soi.Quantity,
					UnitPrice:   &soi.Unitprce,
					Subtotal:    &soi.Xtndprce,
					Uom:         soi.Uofm,
					DeliveryQty: &sii.Quantity,
				})
			}
		}
	}

	if deliveryRunReturn != nil && deliveryRunReturn.Data != nil {
		// get drr delivery run return item
		if deliveryRunReturnItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnItemList(ctx, &logisticService.GetDeliveryRunReturnItemListRequest{
			DeliveryRunReturnId: []int64{deliveryRunReturn.Data.Id},
		}); err != nil || len(deliveryRunReturnItems.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run return item")
			return
		}

		// TODO DELIVERY RETURN ITEM DATA
		// response SubControlTowerGetCourierDetailDRRI
		for _, drri := range deliveryRunReturnItems.Data {
			var (
				// deliveryOrderItem *bridgeService.GetDeliveryOrderItemResponse
				item *bridgeService.GetItemGPResponse
				uom  *bridgeService.GetUomGPResponse
			)

			// 	// if deliveryOrderItem, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderItemDetail(ctx, &bridgeService.GetDeliveryOrderItemRequest{
			// 	// 	Id:drri.deliveryOrderItemID
			// 	// }); err != nil || deliveryOrderItem.Data==nil{
			// 	// 	span.RecordError(err)
			// 	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// 	// 	err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			// 	// return
			// 	// }

			// get configuration for delivery return item
			if glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
				Table:     "delivery_run_return_item",
				Attribute: "item_return_reason",
				ValueInt:  int32(drri.ReturnReason),
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
				return
			}

			if item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
				Id: drri.DeliveryOrderItemId,
			}); err != nil || item.Data == nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "item")
				return
			}

			if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
				Id: item.Data[0].Uomschdl,
			}); err != nil || uom.Data == nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("bridge", "uom")
				return
			}

			// 	dummyDeliverQty := 999.0
			dRRITemp = &dto.GlobalDeliveryRunReturnItem{
				ID:                  drri.Id,
				ReceiveQty:          &drri.ReceiveQty,
				ReturnReason:        int8(drri.ReturnReason),
				ReturnEvidence:      drri.ReturnEvidence,
				ReturnReasonValue:   glossary.Data.ValueName,
				Subtotal:            &drri.Subtotal,
				DeliveryOrderItemId: drri.DeliveryOrderItemId,
				// DeliveryOrderItem: &dto.GlobalDeliveryOrderItem{
				// 	DeliverQty:     &dummyDeliverQty,
				// 	SalesOrderItem: salesOrderItems,
				// },
			}

			for _, v := range salesInvoice.Data[0].Details {
				if v.Itemnmbr == item.Data[0].Itemnmbr {
					dRRITemp.DeliverQty = &v.Quantity
				}
			}
			respDRRI = append(respDRRI, dRRITemp)
		}
	}

	arrivalTime := deliveryRunSheetItem.Data.ArrivalTime.AsTime()
	startedAt := deliveryRunSheetItem.Data.StartedAt.AsTime()
	finishedAt := deliveryRunSheetItem.Data.FinishedAt.AsTime()
	res = dto.CourierAppDetailResponse{
		DeliveryRunSheetItem: &dto.GlobalDeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.Data.Id,
			StepType:                    int8(deliveryRunSheetItem.Data.StepType),
			Status:                      int8(deliveryRunSheetItem.Data.Status),
			Note:                        deliveryRunSheetItem.Data.Note,
			RecipientName:               deliveryRunSheetItem.Data.RecipientName,
			MoneyReceived:               &deliveryRunSheetItem.Data.MoneyReceived,
			DeliveryEvidenceImageURL:    deliveryRunSheetItem.Data.DeliveryEvidenceImageUrl,
			TransactionEvidenceImageURL: deliveryRunSheetItem.Data.TransactionEvidenceImageUrl,
			ArrivalTime:                 &arrivalTime,
			UnpunctualReason:            int8(deliveryRunSheetItem.Data.UnpunctualReason),
			UnpunctualDetail:            int8(deliveryRunSheetItem.Data.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.Data.FarDeliveryReason,
			StartedAt:                   &startedAt,
			FinishedAt:                  &finishedAt,
		},
		Distance: distances,
		SalesOrder: &dto.SubCourierAppDetailSO{
			Code:                    salesOrder.Data[0].Sopnumbe,
			DeliveryDate:            salesOrder.Data[0].ReqShipDate,
			DeliveryFee:             salesOrder.Data[0].Frtamnt,
			VoucherDiscountAmount:   salesOrder.Data[0].Trdisamt,
			PointRedeemAmount:       0, // POINT REDEEM AMOUNT BELUM ADA
			SalesInvoiceTotalCharge: salesInvoice.Data[0].Ordocamt,
			SalesPayment:            salesOrder.Data[0].Pymtrmid,
			Customer: &dto.GlobalCustomer{
				Name: customer.Data[0].Custname,
				// SalesPaymentTerm: &dto.GlobalSalesPaymentTerm{
				// 	Description: salesOrder.Data[0].Pymtrmid,
				// },
			},
			Address: &dto.GlobalAddress{
				AddressName:     address.Data[0].Custname,
				Phone_1:         address.Data[0].PhonE1,
				ShippingAddress: address.Data[0].AddresS1,
				Note:            address.Data[0].GnL_Address_Note,
				Latitude:        &address.Data[0].GnL_Latitude,
				Longitude:       &address.Data[0].GnL_Longitude,
				AdmDivision: &dto.GlobalAdmDivsion{
					PostalCode: admDivision.Data[0].Code,
					SubDistrict: &dto.GlobalSubDistrict{
						Description: address.Data[0].AdministrativeDiv.GnlSubdistrict,
					},
				},
			},
			Wrt: &dto.GlobalWrt{
				StartTime: wrt.Data[0].Strttime,
				EndTime:   wrt.Data[0].Endtime,
			},
			OrderType: &dto.GlobalOrderType{
				Description: orderType.Data.Description,
			},
			Item: salesOrderItems,
		},
	}
	if deliveryRunReturn != nil && deliveryRunReturn.Data != nil {
		res.DeliveryRunReturn = &dto.GlobalDeliveryRunReturn{
			ID:                    deliveryRunReturn.Data.Id,
			Code:                  deliveryRunReturn.Data.Code,
			TotalPrice:            &deliveryRunReturn.Data.TotalPrice,
			TotalCharge:           &deliveryRunReturn.Data.TotalCharge,
			DeliveryRunReturnItem: respDRRI,
		}
	}
	return
}

func (s *CourierAppService) ScanDetail(ctx context.Context, req dto.CourierAppScanDetailRequest) (res dto.CourierAppDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.ScanDetail")
	defer span.End()

	var (
		deliveryRunSheetItems *logisticService.GetDeliveryRunSheetItemListResponse
		deliveryRunSheetItem  *logisticService.DeliveryRunSheetItem

		salesOrder       *bridgeService.GetSalesOrderGPListResponse
		customer         *bridgeService.GetCustomerGPResponse
		salesPaymentTerm *bridgeService.GetSalesPaymentTermDetailResponse
		address          *bridgeService.GetAddressGPResponse
		admDivision      *bridgeService.GetAdmDivisionDetailResponse
		subDistrict      *bridgeService.GetSubDistrictDetailResponse
		wrt              *bridgeService.GetWrtGPResponse
		orderType        *bridgeService.GetOrderTypeGPResponse
		// salesInvoice *bridgeService.GetSalesInvoiceDetailResponse
		items []*dto.GlobalSalesOrderItem

		deliveryRunReturn *logisticService.GetDeliveryRunReturnDetailResponse

		deliveryRunReturnItems *logisticService.GetDeliveryRunReturnItemListResponse
		respDRRI               []*dto.GlobalDeliveryRunReturnItem
		deliveryKoli           *siteService.GetSalesOrderDeliveryKoliResponse
		// salesInvoice *bridgeService.GetSalesInvoiceDetailResponse
		totalKoli int64
	)

	// get sales order
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: req.Code,
	}); err != nil || salesOrder.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	if deliveryRunSheetItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
		Status:       []int32{1, 2, 3, 4},
		StepType:     []int32{2},
		CourierId:    []string{req.CourierId},
		SalesOrderId: []string{salesOrder.Data[0].Sopnumbe},
	}); err != nil || len(deliveryRunSheetItems.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// get delivery run sheet item
	deliveryRunSheetItem = deliveryRunSheetItems.Data[0]

	// get customer
	if customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{
		Id: salesOrder.Data[0].Customer[0].Custnmbr,
	}); err != nil || customer.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	// TODO CHANGE SALES PAYMENT TERM DUMMY
	// sales payment term
	if salesPaymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentTermDetail(ctx, &bridgeService.GetSalesPaymentTermDetailRequest{
		// Id: customer.Data[0].Pymtrmid,
		Id: 1,
	}); err != nil || salesPaymentTerm.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales payment term")
		return
	}

	// TODO CHANGE ADDRESS DUMMY
	// get address
	if address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
		Id: salesOrder.Data[0].Address[0].Prstadcd,
	}); err != nil || address.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		return
	}

	// TODO ADM DIVISION DUMMY
	// get adm division
	admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionDetail(ctx, &bridgeService.GetAdmDivisionDetailRequest{
		// Id: address.Data[0].GnL_Administrative_Code,
		Id: 1,
	})
	if err != nil || admDivision.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
		return
	}

	// get sub district
	if subDistrict, err = s.opt.Client.BridgeServiceGrpc.GetSubDistrictDetail(ctx, &bridgeService.GetSubDistrictDetailRequest{
		Id: admDivision.Data.SubDistrictId,
	}); err != nil || subDistrict.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sub district")
		return
	}

	// TODO CHANGE WRT DUMMY
	// get wrt
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
		GnlRegion: salesOrder.Data[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	// TODO ORDER TYPE DUMMY (?)
	// get order type
	if orderType, err = s.opt.Client.BridgeServiceGrpc.GetOrderTypeGPDetail(ctx, &bridgeService.GetOrderTypeGPDetailRequest{
		Id: salesOrder.Data[0].Docid,
	}); err != nil || orderType.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "order type")
		return
	}

	// prepare res items
	for _, soi := range salesOrder.Data[0].Details {
		var (
			item *bridgeService.GetItemGPResponse
			uom  *bridgeService.GetUomGPResponse
			// deliveryOrderItem *bridgeService.GetDeliveryOrderItemResponse
		)

		if item, err = s.opt.Client.BridgeServiceGrpc.GetItemGPDetail(ctx, &bridgeService.GetItemGPDetailRequest{
			Id: soi.Itemnmbr,
		}); err != nil || item.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "item")
			return
		}

		if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomGPDetail(ctx, &bridgeService.GetUomGPDetailRequest{
			Id: item.Data[0].Uomschdl,
		}); err != nil || uom.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			return
		}

		// TODO GET DELIVERY ORDER ITEM
		// if deliveryOrderItem, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderItemDetail(ctx, &bridgeService.GetDeliveryOrderItemRequest{
		// 	salesOrderItemId:soi.Id
		// }); err != nil || deliveryOrderItem.Data == nil{
		// 	span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorRpcNotFound("bridge", "uom")
		// return
		// }

		deliverQuantity := 9999999.0
		items = append(items, &dto.GlobalSalesOrderItem{
			OrderQty:  &soi.Quantity,
			UnitPrice: &soi.Unitprce,
			Subtotal:  &soi.Xtndprce,
			Item: &dto.GlobalItem{
				Description: item.Data[0].Itemdesc,
				Uom: &dto.GlobalUom{
					Description: uom.Data[0].Umschdsc,
				},
			},
			DeliveryOrderItem: &dto.GlobalDeliveryOrderItem{
				DeliverQty: &deliverQuantity,
			},
		})
	}

	// get drsi delivery run return
	deliveryRunReturn, _ = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnDetail(ctx, &logisticService.GetDeliveryRunReturnDetailRequest{
		DeliveryRunSheetItemId: deliveryRunSheetItem.Id,
	})

	if deliveryRunReturn.Data != nil {
		// get drr delivery run return item
		if deliveryRunReturnItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnItemList(ctx, &logisticService.GetDeliveryRunReturnItemListRequest{
			DeliveryRunReturnId: []int64{deliveryRunReturn.Data.Id},
		}); err != nil || len(deliveryRunReturnItems.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run return item")
			return
		}

		// TODO GET DELIVERY RUN RETURN ITEM
		// response SubControlTowerGetCourierDetailDRRI
		for _, drri := range deliveryRunReturnItems.Data {
			var (
			// deliveryOrderItem *bridgeService.GetDeliveryOrderItemResponse
			// item              *bridgeService.GetItemDetailResponse
			// uom               *bridgeService.GetUomDetailResponse
			)

			// if deliveryOrderItem, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderItemDetail(ctx, &bridgeService.GetDeliveryOrderItemRequest{
			// 	Id:drri.deliveryOrderItemID
			// }); err != nil || deliveryOrderItem.Data==nil{
			// 	span.RecordError(err)
			// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// 	err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			// return
			// }

			// if item, err = s.opt.Client.BridgeServiceGrpc.GetItemDetail(ctx, &bridgeService.GetItemDetailRequest{
			// 	Id: deliveryOrderItem.ItemId,
			// }); err != nil || item.Data==nil {
			// 	span.RecordError(err)
			// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// 	err = edenlabs.ErrorRpcNotFound("bridge", "item")
			// return
			// }

			// if uom, err = s.opt.Client.BridgeServiceGrpc.GetUomDetail(ctx, &bridgeService.GetUomDetailRequest{
			// 	Id: item.Data.UomId,
			// }); err != nil || uom.Data==nil{
			// 	span.RecordError(err)
			// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
			// 	err = edenlabs.ErrorRpcNotFound("bridge", "uom")
			// return
			// }

			dummyDeliverQty := 999.0
			respDRRI = append(respDRRI, &dto.GlobalDeliveryRunReturnItem{
				ID:             drri.Id,
				ReceiveQty:     &drri.ReceiveQty,
				ReturnReason:   int8(drri.ReturnReason),
				ReturnEvidence: drri.ReturnEvidence,
				Subtotal:       &drri.Subtotal,
				DeliveryOrderItem: &dto.GlobalDeliveryOrderItem{
					DeliverQty: &dummyDeliverQty,
					// SalesOrderItem: &dto.GlobalSalesOrderItem{
					// 	Item: &dto.GlobalItem{
					// 		Description: "DUMMY NAME BELUM ADA",
					// 		Uom: &dto.GlobalUom{
					// 			Description: "DUMMY UOM BELUM ADA",
					// 		},
					// 	},
					// },
				},
			})
		}
	}

	// check if there's a delivery for the SO
	if deliveryKoli, err = s.opt.Client.SiteServiceGrpc.GetSalesOrderDeliveryKoli(ctx, &siteService.GetSalesOrderDeliveryKoliRequest{
		SopNumber: req.Code,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "delivery koli")
		return
	}

	// get total koli in 1 so
	for _, v := range deliveryKoli.Data {
		totalKoli += int64(v.Quantity)
	}

	arrivalTime := deliveryRunSheetItem.ArrivalTime.AsTime()
	startedAt := deliveryRunSheetItem.StartedAt.AsTime()
	finishedAt := deliveryRunSheetItem.FinishedAt.AsTime()
	res = dto.CourierAppDetailResponse{
		DeliveryRunSheetItem: &dto.GlobalDeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.Id,
			StepType:                    int8(deliveryRunSheetItem.StepType),
			Status:                      int8(deliveryRunSheetItem.Status),
			Note:                        deliveryRunSheetItem.Note,
			RecipientName:               deliveryRunSheetItem.RecipientName,
			MoneyReceived:               &deliveryRunSheetItem.MoneyReceived,
			DeliveryEvidenceImageURL:    deliveryRunSheetItem.DeliveryEvidenceImageUrl,
			TransactionEvidenceImageURL: deliveryRunSheetItem.TransactionEvidenceImageUrl,
			ArrivalTime:                 &arrivalTime,
			UnpunctualReason:            int8(deliveryRunSheetItem.UnpunctualReason),
			UnpunctualDetail:            int8(deliveryRunSheetItem.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.FarDeliveryReason,
			StartedAt:                   &startedAt,
			FinishedAt:                  &finishedAt,
		},
		Distance: 0, // BELUM ADA
		SalesOrder: &dto.SubCourierAppDetailSO{
			Code:                    salesOrder.Data[0].Sopnumbe,
			DeliveryDate:            salesOrder.Data[0].ReqShipDate,
			DeliveryFee:             999999,  // BELUM ADA
			VoucherDiscountAmount:   0,       // VOUCHER DISCOUNT AMOUNT BELUM ADA
			PointRedeemAmount:       0,       // POINT REDEEM AMOUNT BELUM ADA
			SalesInvoiceTotalCharge: 9999999, // BELUM ADA SALES INVOICE
			Customer: &dto.GlobalCustomer{
				Name: customer.Data[0].Custname,
				SalesPaymentTerm: &dto.GlobalSalesPaymentTerm{
					Description: salesPaymentTerm.Data.Description,
				},
			},
			Address: &dto.GlobalAddress{
				AddressName:     address.Data[0].Custname,
				Phone_1:         address.Data[0].PhonE1,
				ShippingAddress: address.Data[0].AddresS1,
				Note:            address.Data[0].GnL_Address_Note,
				AdmDivision: &dto.GlobalAdmDivsion{
					PostalCode: admDivision.Data.PostalCode,
					SubDistrict: &dto.GlobalSubDistrict{
						Description: address.Data[0].AdministrativeDiv.GnlSubdistrict,
					},
				},
			},
			Wrt: &dto.GlobalWrt{
				StartTime: wrt.Data[0].Strttime,
				EndTime:   wrt.Data[0].Endtime,
			},
			OrderType: &dto.GlobalOrderType{
				Description: orderType.Data[0].Docid,
			},
			Item: items,
		},
		DeliveryRunReturn: &dto.GlobalDeliveryRunReturn{
			ID:                    deliveryRunReturn.Data.Id,
			Code:                  deliveryRunReturn.Data.Code,
			TotalPrice:            &deliveryRunReturn.Data.TotalPrice,
			TotalCharge:           &deliveryRunReturn.Data.TotalCharge,
			DeliveryRunReturnItem: respDRRI,
		},
		Koli: totalKoli,
	}

	return
}

func (s *CourierAppService) Scan(ctx context.Context, req dto.CourierAppScanRequest) (res dto.CourierAppScanResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.Scan")
	defer span.End()

	var (
		countDeliveryRunSheetItem *logisticService.GetDeliveryRunSheetItemListResponse
		salesOrder                *bridgeService.GetSalesOrderGPListResponse
		customer                  *bridgeService.GetCustomerGPResponse
		address                   *bridgeService.GetAddressGPResponse
		admDivision               *bridgeService.GetAdmDivisionGPResponse
		wrt                       *bridgeService.GetWrtGPResponse
		deliveryOrders            *bridgeService.GetDeliveryOrderGPListResponse
		salesInvoices             *bridgeService.GetSalesInvoiceGPListResponse
		// salesPaymentTerm          *bridgeService.GetSalesPaymentTermDetailResponse
		orderType    *bridgeService.GetOrderTypeGPResponse
		deliveryKoli *siteService.GetSalesOrderDeliveryKoliResponse
		// salesInvoice *bridgeService.GetSalesInvoiceDetailResponse
		totalKoli             int64
		deliveryRunSheetItems *logisticService.GetDeliveryRunSheetItemListResponse
	)

	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: req.Code,
	}); err != nil || salesOrder.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// check sales order and courier site if match
	if salesOrder.Data[0].Site[0].Locncode != req.CourierSiteId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", "Anda tidak dapat mengambil sales order ini")
		return
	}

	// check if there's a delivery for the SO
	if deliveryRunSheetItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
		Status:       []int32{1, 2, 3, 4},
		SalesOrderId: []string{salesOrder.Data[0].Sopnumbe},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}
	if len(deliveryRunSheetItems.Data) != 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("sales_order.invalid", "Sales order ini sudah pernah dipindai")
		return
	}

	// TODO UNCOMMENT VALIDATION
	// so status has to be on delivery/ invoiced on delivery / paid on delivery
	// if r.SalesOrder.Status != 7 && r.SalesOrder.Status != 10 && r.SalesOrder.Status != 13 {
	// 	o.Failure("status.invalid", util.ErrorStatusNotAcceptableInd("sales order"))
	// 	return o
	// }

	// check sales order and courier site if match
	if salesOrder.Data[0].Site[0].Locncode != req.CourierSiteId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", "Anda tidak dapat mengambil sales order ini")
		return
	}

	// sales order should not exist in delivery run sheet item
	if countDeliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
		Status:       []int32{1, 2, 3, 4},
		SalesOrderId: []string{salesOrder.Data[0].Sopnumbe},
	}); err != nil || len(countDeliveryRunSheetItem.Data) != 0 {
		// TODO UINCOMMENT VALIDATION
		// span.RecordError(err)
		// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
		// 	err = edenlabs.ErrorValidation("sales_order.invalid", "Sales order sudah dipindai")
		// 	return
	}

	// get customer
	if customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{
		Id: salesOrder.Data[0].Customer[0].Custnmbr,
	}); err != nil || customer.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	// TODO SALES PAYMENT TERM DUMMY
	// sales payment term
	// if salesPaymentTerm, err = s.opt.Client.BridgeServiceGrpc.GetSalesPaymentTermDetail(ctx, &bridgeService.GetSalesPaymentTermDetailRequest{
	// 	// Id: customer.Data[0].Pymtrmid,
	// 	Id: 1,
	// }); err != nil || salesPaymentTerm.Data == nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "sales payment term")
	// 	return
	// }

	// TODO ADDRESS DUMMY
	if address, err = s.opt.Client.BridgeServiceGrpc.GetAddressGPDetail(ctx, &bridgeService.GetAddressGPDetailRequest{
		Id: salesOrder.Data[0].Address[0].Prstadcd,
	}); err != nil || address.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "address")
		return
	}

	// TODO ADM DIVISION DUMMY
	// get adm division
	admDivision, err = s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridgeService.GetAdmDivisionGPListRequest{
		// Id: address.Data[0].GnL_Administrative_Code,
		AdmDivisionCode: address.Data[0].AdministrativeDiv.GnlAdministrativeCode,
	})
	if err != nil || admDivision.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "adm division")
		return
	}

	// if subDistrict, err = s.opt.Client.BridgeServiceGrpc.GetSubDistrictDetail(ctx, &bridgeService.GetSubDistrictDetailRequest{
	// 	Id: admDivision.Data.SubDistrictId,
	// }); err != nil || subDistrict.Data == nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "sub district")
	// 	return
	// }

	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
		GnlRegion: salesOrder.Data[0].Wrt[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
	}

	// TODO ORDER TYPE DUMMY (?)
	// get order type
	if orderType, err = s.opt.Client.BridgeServiceGrpc.GetOrderTypeGPDetail(ctx, &bridgeService.GetOrderTypeGPDetailRequest{
		Id: salesOrder.Data[0].Docid,
	}); err != nil || orderType.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "order type")
		return
	}

	// TODO GET SALES INVOICE
	// sales invoice
	if orderType, err = s.opt.Client.BridgeServiceGrpc.GetOrderTypeGPDetail(ctx, &bridgeService.GetOrderTypeGPDetailRequest{
		Id: salesOrder.Data[0].Docid,
	}); err != nil || orderType.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "order type")
		return
	}

	// check if there's a delivery for the SO
	if deliveryKoli, err = s.opt.Client.SiteServiceGrpc.GetSalesOrderDeliveryKoli(ctx, &siteService.GetSalesOrderDeliveryKoliRequest{
		SopNumber: req.Code,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("site", "delivery koli")
		return
	}

	// get total koli in 1 so
	for _, v := range deliveryKoli.Data {
		totalKoli += int64(v.Quantity)
	}

	if deliveryOrders, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderListGP(ctx, &bridgeService.GetDeliveryOrderGPListRequest{
		Limit:    1,
		Offset:   1,
		SopNumbe: salesOrder.Data[0].Sopnumbe,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "delivery order")
		return
	}
	if salesInvoices, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
		SoNumber: salesOrder.Data[0].Sopnumbe,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
		return
	}

	if len(deliveryOrders.Data) == 0 && len(salesInvoices.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("courier", "sales order ini belum memiliki delivery order/ sales invoice")
		return
	}
	var salesOrderItems []*dto.GlobalSalesOrderItem

	for _, soi := range salesOrder.Data[0].Details {
		for _, sii := range salesInvoices.Data[0].Details {
			if soi.Itemnmbr == sii.Itemnmbr {
				salesOrderItems = append(salesOrderItems, &dto.GlobalSalesOrderItem{
					ItemNumber:  soi.Itemnmbr,
					ItemDesc:    soi.Itemdesc,
					OrderQty:    &soi.Quantity,
					UnitPrice:   &soi.Unitprce,
					Subtotal:    &soi.Xtndprce,
					Uom:         soi.Uofm,
					DeliveryQty: &sii.Quantity,
				})
			}
		}
	}

	res = dto.CourierAppScanResponse{
		SalesOrder: &dto.GlobalSalesOrder{
			Code:      salesOrder.Data[0].Sopnumbe,
			Status:    salesOrder.Data[0].Status,
			OrderDate: salesOrder.Data[0].Docdate,
			Address: &dto.GlobalAddress{
				AddressName: customer.Data[0].AddresS1 + ", " + customer.Data[0].AddresS2,
				AdmDivision: &dto.GlobalAdmDivsion{
					District: &dto.GlobalDistrict{
						Description: admDivision.Data[0].District,
					},
				},
			},
			SalesOrderItem: salesOrderItems,
			Customer: &dto.GlobalCustomer{
				Name: customer.Data[0].Custname,
				SalesPaymentTerm: &dto.GlobalSalesPaymentTerm{
					Description: salesOrder.Data[0].Pymtrmid,
				},
			},
			Wrt: &dto.GlobalWrt{
				StartTime: wrt.Data[0].Strttime,
				EndTime:   wrt.Data[0].Endtime,
			},
			OrderType: &dto.GlobalOrderType{
				Description: orderType.Data[0].Docid,
			},
			// DeliveryOrder: &dto.GlobalDeliveryOrder{
			// 	Id: 1,
			// 	DeliveryOrderItem: []*dto.GlobalDeliveryOrderItem{ // FULL DUMMY
			// 		{
			// 			Id:         1,
			// 			DeliverQty: &dummyQty,
			// 			SalesOrderItem: &dto.GlobalSalesOrderItem{
			// 				OrderQty: &dummyQty,
			// 			},
			// 		},
			// 		{
			// 			Id:         2,
			// 			DeliverQty: &dummyQty,
			// 			SalesOrderItem: &dto.GlobalSalesOrderItem{
			// 				OrderQty: &dummyQty,
			// 			},
			// 		},
			// 	},
			// },
		},
		Koli: totalKoli,
	}

	return
}

func (s *CourierAppService) SelfAssign(ctx context.Context, req dto.CourierAppSelfAssignRequest) (res dto.CourierAppSelfAssignResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.SelfAssign")
	defer span.End()

	var (
		salesOrder            *bridgeService.GetSalesOrderGPListResponse
		deliveryRunSheets     *logisticService.GetDeliveryRunSheetListResponse
		deliveryRunSheetItems *logisticService.GetDeliveryRunSheetItemListResponse
		delivery              *logisticService.DeliveryRunSheetItem
		supportiveData        string // this variable for save the linking for all auditlog related to this process
		responseDrsi          *logisticService.CreateDeliveryRunSheetItemResponse
	)

	supportiveData = utils.GenerateUnixTime() + utils.GenerateUUID()

	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: req.SopNumber,
	}); err != nil || salesOrder.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// check sales order and courier site if match
	if salesOrder.Data[0].Site[0].Locncode != req.CourierSiteId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", "Anda tidak dapat mengambil sales order ini")
		return
	}

	// check if there's a delivery for the SO
	if deliveryRunSheetItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
		Status:       []int32{1, 2, 3, 4},
		SalesOrderId: []string{salesOrder.Data[0].Sopnumbe},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}
	if len(deliveryRunSheetItems.Data) != 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("sales_order.invalid", "Sales order ini sudah pernah dipindai")
		return
	}

	// check if drs exist
	deliveryRunSheets, _ = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetList(ctx, &logisticService.GetDeliveryRunSheetListRequest{
		Status:    []int32{2},
		CourierId: []string{req.CourierId},
	})

	if len(deliveryRunSheets.Data) != 0 {
		// insert pickup
		if responseDrsi, err = s.opt.Client.LogisticServiceGrpc.CreateDeliveryRunSheetItemPickup(ctx, &logisticService.CreateDeliveryRunSheetItemRequest{
			Model: &logisticService.DeliveryRunSheetItem{
				DeliveryRunSheetId: deliveryRunSheets.Data[0].Id,
				CourierId:          req.CourierId,
				SalesOrderId:       salesOrder.Data[0].Sopnumbe,
				Latitude:           &req.Latitude,
				Longitude:          &req.Longitude,
				ArrivalTime:        timestamppb.Now(),
				CreatedAt:          timestamppb.Now(),
				StartedAt:          timestamppb.Now(),
				FinishedAt:         timestamppb.Now(),
			},
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
			return
		}
		// create log
		if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
			Log: &auditService.Log{
				UserIdGp:       req.CourierId,
				ReferenceId:    strconv.Itoa(int(responseDrsi.Data.Id)),
				Type:           "delivery_run_sheet_item_pickup",
				Function:       "create",
				CreatedAt:      timestamppb.New(time.Now()),
				SupportiveData: supportiveData,
			},
		}); err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			span.RecordError(err)
		}

		delivery = &logisticService.DeliveryRunSheetItem{
			DeliveryRunSheetId: deliveryRunSheets.Data[0].Id,
			CourierId:          req.CourierId,
			SalesOrderId:       salesOrder.Data[0].Sopnumbe,
			CreatedAt:          timestamppb.Now(),
		}

		// insert delivery
		if responseDrsi, err = s.opt.Client.LogisticServiceGrpc.CreateDeliveryRunSheetItemDelivery(ctx, &logisticService.CreateDeliveryRunSheetItemRequest{
			Model: delivery,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
			return
		}

		if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
			Log: &auditService.Log{
				UserIdGp:       req.CourierId,
				ReferenceId:    strconv.Itoa(int(responseDrsi.Data.Id)),
				Type:           "delivery_run_sheet_item_delivery",
				Function:       "create",
				CreatedAt:      timestamppb.New(time.Now()),
				SupportiveData: supportiveData,
			},
		}); err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			span.RecordError(err)
		}

	} else {
		var deliveryRunSheet *logisticService.CreateDeliveryRunSheetResponse

		// create a new drs
		if deliveryRunSheet, err = s.opt.Client.LogisticServiceGrpc.CreateDeliveryRunSheet(ctx, &logisticService.CreateDeliveryRunSheetRequest{
			Model: &logisticService.DeliveryRunSheet{
				CourierId:         req.CourierId,
				DeliveryDate:      timestamppb.Now(),
				StartedAt:         timestamppb.Now(),
				StartingLatitude:  &req.Latitude,
				StartingLongitude: &req.Longitude,
				Status:            2,
			},
		}); err != nil || deliveryRunSheet.Data == nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet")
			return
		}
		if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
			Log: &auditService.Log{
				UserIdGp:       req.CourierId,
				ReferenceId:    strconv.Itoa(int(deliveryRunSheet.Data.Id)),
				Type:           "delivery_run_sheet",
				Function:       "create",
				CreatedAt:      timestamppb.New(time.Now()),
				SupportiveData: supportiveData,
			},
		}); err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			span.RecordError(err)
		}
		// insert pickup
		if responseDrsi, err = s.opt.Client.LogisticServiceGrpc.CreateDeliveryRunSheetItemPickup(ctx, &logisticService.CreateDeliveryRunSheetItemRequest{
			Model: &logisticService.DeliveryRunSheetItem{
				DeliveryRunSheetId: deliveryRunSheet.Data.Id,
				CourierId:          req.CourierId,
				SalesOrderId:       salesOrder.Data[0].Sopnumbe,
				StepType:           1,
				Latitude:           &req.Latitude,
				Longitude:          &req.Longitude,
				Status:             3,
				ArrivalTime:        timestamppb.Now(),
				CreatedAt:          timestamppb.Now(),
				StartedAt:          timestamppb.Now(),
				FinishedAt:         timestamppb.Now(),
			},
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
			return
		}
		if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
			Log: &auditService.Log{
				UserIdGp:       req.CourierId,
				ReferenceId:    strconv.Itoa(int(responseDrsi.Data.Id)),
				Type:           "delivery_run_sheet_item_pickup",
				Function:       "create",
				CreatedAt:      timestamppb.New(time.Now()),
				SupportiveData: supportiveData,
			},
		}); err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			span.RecordError(err)
		}

		delivery = &logisticService.DeliveryRunSheetItem{
			DeliveryRunSheetId: deliveryRunSheet.Data.Id,
			CourierId:          req.CourierId,
			SalesOrderId:       salesOrder.Data[0].Sopnumbe,
			StepType:           2,
			Status:             1,
			CreatedAt:          timestamppb.Now(),
		}

		// insert delivery
		if responseDrsi, err = s.opt.Client.LogisticServiceGrpc.CreateDeliveryRunSheetItemDelivery(ctx, &logisticService.CreateDeliveryRunSheetItemRequest{
			Model: delivery,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
			return
		}

		if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
			Log: &auditService.Log{
				UserIdGp:       req.CourierId,
				ReferenceId:    strconv.Itoa(int(responseDrsi.Data.Id)),
				Type:           "delivery_run_sheet_item_delivery",
				Function:       "create",
				CreatedAt:      timestamppb.New(time.Now()),
				SupportiveData: supportiveData,
			},
		}); err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			span.RecordError(err)
		}

	}

	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserIdGp:       req.CourierId,
			ReferenceId:    strconv.Itoa(int(delivery.DeliveryRunSheetId)),
			Type:           "self_assign",
			Function:       "create",
			CreatedAt:      timestamppb.New(time.Now()),
			SupportiveData: supportiveData,
		},
	}); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	createdAt := delivery.CreatedAt.AsTime()
	res = dto.CourierAppSelfAssignResponse{
		DeliveryRunSheetItem: &dto.GlobalDeliveryRunSheetItem{
			StepType:           int8(delivery.StepType),
			Status:             int8(delivery.Status),
			CreatedAt:          &createdAt,
			DeliveryRunSheetID: delivery.DeliveryRunSheetId,
			CourierID:          delivery.CourierId,
			SalesOrderID:       delivery.SalesOrderId,
		},
	}

	return
}

func (s *CourierAppService) StartDelivery(ctx context.Context, req dto.CourierAppStartDeliveryRequest) (res dto.CourierAppStartDeliveryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.StartDelivery")
	defer span.End()

	var (
		deliveryRunSheetItem  *logisticService.GetDeliveryRunSheetItemDetailResponse
		deliveryRunSheetItems *logisticService.GetDeliveryRunSheetItemListResponse
		updateDRSI            *logisticService.StartDeliveryRunSheetItemResponse
	)

	// get delivery run sheet item
	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: req.Id,
	}); err != nil || deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// check if DRSI is courier's job
	if deliveryRunSheetItem.Data.CourierId != req.CourierId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", "Tugas ini bukan tugas anda")
		return
	}

	// check if DRSI is a delivery type
	if deliveryRunSheetItem.Data.StepType != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("step_type.invalid", "Tugas harus berupa pengantaran")
		return
	}

	// check if DRSI status is active or postponed
	if deliveryRunSheetItem.Data.Status != 1 && deliveryRunSheetItem.Data.Status != 4 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("status.invalid", "Status harus dalam keadaan tidak aktif")
		return
	}

	// check if there's another job on going
	if deliveryRunSheetItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
		Status:    []int32{2},
		CourierId: []string{req.CourierId},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}
	if len(deliveryRunSheetItems.Data) != 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("multiple_job.invalid", "Anda memiliki pekerjaan yang sedang berjalan")
		return
	}

	if updateDRSI, err = s.opt.Client.LogisticServiceGrpc.StartDeliveryRunSheetItem(ctx, &logisticService.StartDeliveryRunSheetItemRequest{
		Id:        deliveryRunSheetItem.Data.Id,
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserIdGp:    req.CourierId,
			ReferenceId: strconv.Itoa(int(deliveryRunSheetItem.Data.DeliveryRunSheetId)),
			Type:        "delivery run sheet",
			Function:    "Start Delivery",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	}); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	startedAt := timex.ToLocTime(ctx, updateDRSI.Data.StartedAt.AsTime())
	res = dto.CourierAppStartDeliveryResponse{
		DeliveryRunSheetItem: &dto.GlobalDeliveryRunSheetItem{
			Id:        updateDRSI.Data.Id,
			Status:    int8(updateDRSI.Data.Status),
			StartedAt: &startedAt,
		},
	}

	return
}

func (s *CourierAppService) SuccessDelivery(ctx context.Context, req dto.CourierAppSuccessDeliveryRequest) (res dto.CourierAppSuccessDeliveryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.SuccessDelivery")
	defer span.End()

	var (
		deliveryRunSheetItem *logisticService.GetDeliveryRunSheetItemDetailResponse
		salesOrder           *bridgeService.GetSalesOrderGPListResponse
		// customer             *bridgeService.GetCustomerGPResponse
		// salesPaymentTerm *bridgeService.GetSalesPaymentTermDetailResponse
		// orderType                  *bridgeService.GetOrderTypeDetailResponse
		wrt                        *bridgeService.GetWrtGPResponse
		countDeliveryRunSheetItems *logisticService.GetDeliveryRunSheetItemListResponse
		merchantDeliveryLog        *logisticService.GetFirstMerchantDeliveryLogResponse
		deliveryRunReturn          *logisticService.GetDeliveryRunReturnDetailResponse
		deliveryRunReturnItems     *logisticService.GetDeliveryRunReturnItemListResponse
		updateDRSI                 *logisticService.SuccessDeliveryRunSheetItemResponse
	)

	// get delivery run sheet item
	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: req.Id,
	}); err != nil || deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// check if DRSI is courier's job
	if deliveryRunSheetItem.Data.CourierId != req.CourierId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", "Tugas ini bukan tugas anda")
		return
	}

	// check if DRSI is a delivery type
	if deliveryRunSheetItem.Data.StepType != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("step_type.invalid", "Tugas harus berupa pengantaran")
		return
	}

	// check if DRSI status on progress
	if deliveryRunSheetItem.Data.Status != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("status.invalid", "Status harus dalam keadaan aktif")
		return
	}

	// get sales order detail
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: deliveryRunSheetItem.Data.SalesOrderId,
	}); err != nil || salesOrder.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	if salesOrder.Data[0].Pymtrmid == "COD" {
		if req.MoneyReceived <= 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorValidation("money_received.required", "Uang yang diterima dibutuhkan")
			return
		}
	}

	// TODO ORDER TYPE DUMMY
	// get order type
	// if orderType, err = s.opt.Client.BridgeServiceGrpc.GetOrderTypeGPDetail(ctx, &bridgeService.GetOrderTypeDetailRequest{
	// 	// Id: salesOrder.Data.OrderTypeId,
	// 	Id: 1,
	// }); err != nil || orderType.Data == nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "order type")
	// 	return
	// }

	// get wrt detail
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
		GnlRegion: salesOrder.Data[0].Wrt[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
	}

	// check if the errand is the last delivery run sheet item
	if countDeliveryRunSheetItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
		Status:             []int32{1, 2, 4},
		DeliveryRunSheetId: []int64{deliveryRunSheetItem.Data.DeliveryRunSheetId},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	if len(countDeliveryRunSheetItems.Data) == 1 {
		if _, err = s.opt.Client.LogisticServiceGrpc.FinishDeliveryRunSheet(ctx, &logisticService.FinishDeliveryRunSheetRequest{
			Id:        countDeliveryRunSheetItems.Data[0].DeliveryRunSheetId,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
			return
		}
	}

	// // self pickup will never be late/early and away from location
	// if orderType.Data.Id != 6 {
	var (
		timeStartLimit  time.Time
		timeEndLimit    time.Time
		drsiArrivalTime time.Time
	)

	// Punctuality test
	wrtStartClock := wrt.Data[0].Strttime[0:2]
	wrtStartClockInt, _ := strconv.Atoi(wrtStartClock)
	wrtEndClock := wrt.Data[0].Endtime[0:2]
	wrtEndClockInt, _ := strconv.Atoi(wrtEndClock)

	// get delivery date combined with wrt
	layout1 := "2006-01-02"
	layout := "2006-01-02T15:04:05"
	dateTime, err := time.Parse(layout1, salesOrder.Data[0].ReqShipDate)
	outputString := dateTime.Format(layout)

	if salesOrder.Data[0].Docdate == "" {
		outputString = time.Now().Format("2006-01-02T15:04:05")
	}
	var docDelivDate time.Time
	docDelivDate, err = time.Parse(layout, outputString)

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	timeStartLimit = docDelivDate.Add(time.Duration(wrtStartClockInt) * time.Hour)
	timeEndLimit = docDelivDate.Add(time.Duration(wrtEndClockInt) * time.Hour)
	drsiArrivalTime = timex.ToLocTime(ctx, deliveryRunSheetItem.Data.ArrivalTime.AsTime())
	drsiArrivalTime, _ = time.Parse("2006-01-02 15:04:05", drsiArrivalTime.Format("2006-01-02 15:04:05"))

	// compare time with range of time limit
	if drsiArrivalTime.Before(timeStartLimit) || drsiArrivalTime.After(timeEndLimit) {
		// if it's not punctual then unpunctual_reason is required, early = 1 & late =2
		if drsiArrivalTime.Before(timeStartLimit) {
			req.UnpunctualDetail = 1
			if _, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
				Table:     "delivery_run_sheet_item",
				Attribute: "early_delivery_reason",
				ValueInt:  int32(req.UnpunctualReason),
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
				return
			}
		} else {
			req.UnpunctualDetail = 2
			if _, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryDetail(ctx, &configurationService.GetGlossaryDetailRequest{
				Table:     "delivery_run_sheet_item",
				Attribute: "late_delivery_reason",
				ValueInt:  int32(req.UnpunctualReason),
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
				return
			}
		}
	} else {
		req.UnpunctualReason = 0
	}

	// Proximity test
	if merchantDeliveryLog, err = s.opt.Client.LogisticServiceGrpc.GetFirstMerchantDeliveryLog(ctx, &logisticService.GetFirstMerchantDeliveryLogRequest{
		DeliveryRunSheetItemId: deliveryRunSheetItem.Data.Id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "merchant delivery log")
		return
	}
	if merchantDeliveryLog.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("merchant_delivery_log.required", "Foto pengiriman dibutuhkan")
		return
	}

	distance := distance(*merchantDeliveryLog.Data.Latitude, *merchantDeliveryLog.Data.Longitude, req.Latitude, req.Longitude)

	// if more than 100 meters and no far delivery reason
	if distance > 100 && req.FarDeliveryReason == "" {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("far_delivery_reason.required", "Alasan jauh dari lokasi pelanggan dibutuhkan")
		return
	}
	// }

	// get drsi delivery run return
	deliveryRunReturn, _ = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnDetail(ctx, &logisticService.GetDeliveryRunReturnDetailRequest{
		DeliveryRunSheetItemId: deliveryRunSheetItem.Data.Id,
	})

	if deliveryRunReturn != nil && deliveryRunReturn.Data != nil {
		// get drr delivery run return item
		if deliveryRunReturnItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnItemList(ctx, &logisticService.GetDeliveryRunReturnItemListRequest{
			DeliveryRunReturnId: []int64{deliveryRunReturn.Data.Id},
		}); err != nil && deliveryRunReturnItems != nil && len(deliveryRunReturnItems.Data) == 0 {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run return item")
			return
		}

		for _, drri := range deliveryRunReturnItems.Data {
			// TODO
			fmt.Println(drri)
		}
	}

	// update drsi
	if updateDRSI, err = s.opt.Client.LogisticServiceGrpc.SuccessDeliveryRunSheetItem(ctx, &logisticService.SuccessDeliveryRunSheetItemRequest{
		Id:                          deliveryRunSheetItem.Data.Id,
		Latitude:                    &req.Latitude,
		Longitude:                   &req.Longitude,
		Note:                        req.Note,
		RecipientName:               req.RecipientName,
		MoneyReceived:               req.MoneyReceived,
		DeliveryEvidenceImageUrl:    req.DeliveryEvidence,
		TransactionEvidenceImageUrl: req.TransactionEvidence,
		UnpunctualReason:            int32(req.UnpunctualReason),
		UnpunctualDetail:            int32(req.UnpunctualDetail),
		FarDeliveryReason:           req.FarDeliveryReason,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserIdGp:    req.CourierId,
			ReferenceId: strconv.Itoa(int(deliveryRunSheetItem.Data.DeliveryRunSheetId)),
			Type:        "delivery run sheet",
			Function:    "Success Delivery",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	}); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	arrivalTime := updateDRSI.Data.ArrivalTime.AsTime()
	createdAt := updateDRSI.Data.CreatedAt.AsTime()
	startedAt := updateDRSI.Data.StartedAt.AsTime()
	finishedAt := updateDRSI.Data.FinishedAt.AsTime()
	res = dto.CourierAppSuccessDeliveryResponse{
		DeliveryRunSheetItem: &dto.GlobalDeliveryRunSheetItem{
			Id:                          updateDRSI.Data.Id,
			StepType:                    int8(updateDRSI.Data.StepType),
			Latitude:                    updateDRSI.Data.Latitude,
			Longitude:                   updateDRSI.Data.Longitude,
			Status:                      int8(updateDRSI.Data.Status),
			Note:                        updateDRSI.Data.Note,
			RecipientName:               updateDRSI.Data.RecipientName,
			MoneyReceived:               &updateDRSI.Data.MoneyReceived,
			DeliveryEvidenceImageURL:    updateDRSI.Data.DeliveryEvidenceImageUrl,
			TransactionEvidenceImageURL: updateDRSI.Data.TransactionEvidenceImageUrl,
			ArrivalTime:                 &arrivalTime,
			UnpunctualReason:            int8(updateDRSI.Data.UnpunctualReason),
			UnpunctualDetail:            int8(updateDRSI.Data.UnpunctualDetail),
			FarDeliveryReason:           updateDRSI.Data.FarDeliveryReason,
			CreatedAt:                   &createdAt,
			StartedAt:                   &startedAt,
			FinishedAt:                  &finishedAt,
			DeliveryRunSheetID:          updateDRSI.Data.DeliveryRunSheetId,
			CourierID:                   updateDRSI.Data.CourierId,
			SalesOrderID:                updateDRSI.Data.SalesOrderId,
		},
	}

	return
}

func (s *CourierAppService) PostponeDelivery(ctx context.Context, req dto.CourierAppPostponeDeliveryRequest) (res dto.CourierAppPostponeDeliveryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.PostponeDelivery")
	defer span.End()

	var (
		deliveryRunSheetItem *logisticService.GetDeliveryRunSheetItemDetailResponse
		deliveryRunReturn    *logisticService.GetDeliveryRunReturnDetailResponse
		updateDRSI           *logisticService.PostponeDeliveryRunSheetItemResponse
	)

	// get delivery run sheet item
	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: req.Id,
	}); err != nil || deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}
	// check if DRSI is courier's job
	if deliveryRunSheetItem.Data.CourierId != req.CourierId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", "Tugas ini bukan tugas anda")
		return
	}

	// check if DRSI is a delivery type
	if deliveryRunSheetItem.Data.StepType != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("step_type.invalid", "Tugas harus berupa pengantaran")
		return
	}

	// check if DRSI status on progress
	if deliveryRunSheetItem.Data.Status != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("status.invalid", "Status harus dalam keadaan aktif")
		return
	}

	// TODO UNCOMMENT VALIDATION
	// // get drsi delivery run return
	deliveryRunReturn, _ = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnDetail(ctx, &logisticService.GetDeliveryRunReturnDetailRequest{
		DeliveryRunSheetItemId: deliveryRunSheetItem.Data.Id,
	})

	if deliveryRunReturn != nil && deliveryRunReturn.Data != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("delivery_run_return.exist", "Terdapat pengambalian pengiriman")
		return
	}

	// update drsi
	if updateDRSI, err = s.opt.Client.LogisticServiceGrpc.PostponeDeliveryRunSheetItem(ctx, &logisticService.PostponeDeliveryRunSheetItemRequest{
		Id:   deliveryRunSheetItem.Data.Id,
		Note: req.Note,
	}); err != nil || updateDRSI.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// create postpone log
	if _, err = s.opt.Client.LogisticServiceGrpc.CreatePostponeDeliveryLog(ctx, &logisticService.CreatePostponeDeliveryLogRequest{
		Model: &logisticService.PostponeDeliveryLog{
			DeliveryRunSheetItemId: deliveryRunSheetItem.Data.Id,
			PostponeReason:         req.Note,
			StartedAtUnix:          deliveryRunSheetItem.Data.StartedAt.AsTime().Unix(),
			PostponedAtUnix:        time.Now().Unix(),
			PostponeEvidence:       req.PostponeDeliveryEvidence,
		},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "postpone delivery log")
		return
	}

	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserIdGp:    req.CourierId,
			ReferenceId: strconv.Itoa(int(deliveryRunSheetItem.Data.DeliveryRunSheetId)),
			Type:        "delivery run sheet",
			Function:    "Postpone Delivery",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	}); err != nil {

		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	arrivalTime := timex.ToLocTime(ctx, deliveryRunSheetItem.Data.ArrivalTime.AsTime())
	createdAt := timex.ToLocTime(ctx, deliveryRunSheetItem.Data.CreatedAt.AsTime())
	startedAt := timex.ToLocTime(ctx, deliveryRunSheetItem.Data.StartedAt.AsTime())
	finishedAt := timex.ToLocTime(ctx, deliveryRunSheetItem.Data.FinishedAt.AsTime())
	res = dto.CourierAppPostponeDeliveryResponse{
		DeliveryRunSheetItem: &dto.GlobalDeliveryRunSheetItem{
			Id:                          deliveryRunSheetItem.Data.Id,
			StepType:                    int8(deliveryRunSheetItem.Data.StepType),
			Latitude:                    deliveryRunSheetItem.Data.Latitude,
			Longitude:                   deliveryRunSheetItem.Data.Longitude,
			Status:                      int8(updateDRSI.Data.Status),
			Note:                        updateDRSI.Data.Note,
			RecipientName:               deliveryRunSheetItem.Data.RecipientName,
			MoneyReceived:               &deliveryRunSheetItem.Data.MoneyReceived,
			DeliveryEvidenceImageURL:    deliveryRunSheetItem.Data.DeliveryEvidenceImageUrl,
			TransactionEvidenceImageURL: deliveryRunSheetItem.Data.TransactionEvidenceImageUrl,
			ArrivalTime:                 &arrivalTime,
			UnpunctualReason:            int8(deliveryRunSheetItem.Data.UnpunctualReason),
			UnpunctualDetail:            int8(deliveryRunSheetItem.Data.UnpunctualDetail),
			FarDeliveryReason:           deliveryRunSheetItem.Data.FarDeliveryReason,
			CreatedAt:                   &createdAt,
			StartedAt:                   &startedAt,
			FinishedAt:                  &finishedAt,
			DeliveryRunSheetID:          deliveryRunSheetItem.Data.DeliveryRunSheetId,
			CourierID:                   deliveryRunSheetItem.Data.CourierId,
			SalesOrderID:                deliveryRunSheetItem.Data.SalesOrderId,
		},
	}

	return
}

func (s *CourierAppService) FailDelivery(ctx context.Context, req dto.CourierAppFailDeliveryRequest) (res dto.CourierAppFailDeliveryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.FailDelivery")
	defer span.End()

	var (
		deliveryRunSheetItem       *logisticService.GetDeliveryRunSheetItemDetailResponse
		deliveryRunSheetItems      *logisticService.GetDeliveryRunSheetItemListResponse
		deliveryRunReturn          *logisticService.GetDeliveryRunReturnDetailResponse
		countDeliveryRunSheetItems *logisticService.GetDeliveryRunSheetItemListResponse
		updateDRSI                 *logisticService.FailDeliveryDeliveryRunSheetItemResponse
	)

	// get delivery run sheet item
	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: req.Id,
	}); err != nil || deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// TODO UNCOMMENT VALIDATION
	// // check if DRSI is courier's job
	if deliveryRunSheetItem.Data.CourierId != req.CourierId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", "Tugas ini bukan tugas anda")
		return
	}

	// check if DRSI is a delivery type
	if deliveryRunSheetItem.Data.StepType != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("step_type.invalid", "Tugas harus berupa pengantaran")
		return
	}

	// the courier must only have 2 delivery run sheet items for that sales order with a status of != fail
	if deliveryRunSheetItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
		Status:       []int32{1, 2, 3, 4},
		CourierId:    []string{req.CourierId},
		SalesOrderId: []string{deliveryRunSheetItem.Data.SalesOrderId},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}
	if len(deliveryRunSheetItems.Data) != 2 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("delivery_run_sheet_item.invalid", "delivery run sheet item")
		return
	}

	// TODO UNCOMMENT VALIDATION
	// // get drsi delivery run return
	deliveryRunReturn, _ = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunReturnDetail(ctx, &logisticService.GetDeliveryRunReturnDetailRequest{
		DeliveryRunSheetItemId: deliveryRunSheetItem.Data.Id,
	})
	if deliveryRunReturn != nil && deliveryRunReturn.Data != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("delivery_run_return.exist", "Terdapat pengambalian pengiriman")
		return
	}

	// check if the errand is the last delivery run sheet item
	if countDeliveryRunSheetItems, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemList(ctx, &logisticService.GetDeliveryRunSheetItemListRequest{
		Status:             []int32{1, 2, 4},
		DeliveryRunSheetId: []int64{deliveryRunSheetItem.Data.DeliveryRunSheetId},
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	if len(countDeliveryRunSheetItems.Data) == 1 {
		if _, err = s.opt.Client.LogisticServiceGrpc.FinishDeliveryRunSheet(ctx, &logisticService.FinishDeliveryRunSheetRequest{
			Id:        countDeliveryRunSheetItems.Data[0].DeliveryRunSheetId,
			Latitude:  req.Latitude,
			Longitude: req.Longitude,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
			return
		}
	}
	for _, drsi := range deliveryRunSheetItems.Data {
		if drsi.StepType == 1 { // Pickup
			if _, err = s.opt.Client.LogisticServiceGrpc.FailPickupDeliveryRunSheetItem(ctx, &logisticService.FailPickupDeliveryRunSheetItemRequest{
				Id:   drsi.Id,
				Note: req.Note,
			}); err != nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
				return
			}

		} else { // Delivery
			if updateDRSI, err = s.opt.Client.LogisticServiceGrpc.FailDeliveryDeliveryRunSheetItem(ctx, &logisticService.FailDeliveryDeliveryRunSheetItemRequest{
				Id:        drsi.Id,
				Latitude:  &req.Latitude,
				Longitude: &req.Longitude,
				Note:      req.Note,
			}); err != nil || updateDRSI.Data == nil {
				span.RecordError(err)
				s.opt.Logger.AddMessage(log.ErrorLevel, err)
				err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
				return
			}
		}
	}

	if _, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserIdGp:    req.CourierId,
			ReferenceId: strconv.Itoa(int(deliveryRunSheetItem.Data.DeliveryRunSheetId)),
			Type:        "delivery run sheet",
			Function:    "Fail Delivery",
			CreatedAt:   timestamppb.New(time.Now()),
		},
	}); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		span.RecordError(err)
		err = edenlabs.ErrorRpcNotFound("audit", "audit_log")
		return
	}

	// arrivalTime := updateDRSI.Data.ArrivalTime.AsTime()
	// createdAt := updateDRSI.Data.CreatedAt.AsTime()
	// startedAt := updateDRSI.Data.StartedAt.AsTime()
	// finishedAt := updateDRSI.Data.FinishedAt.AsTime()
	res = dto.CourierAppFailDeliveryResponse{
		DeliveryRunSheetItem: &dto.GlobalDeliveryRunSheetItem{
			Id:        updateDRSI.Data.Id,
			StepType:  int8(updateDRSI.Data.StepType),
			Latitude:  updateDRSI.Data.Latitude,
			Longitude: updateDRSI.Data.Longitude,
			Status:    int8(updateDRSI.Data.Status),
			// Note:                        updateDRSI.Data.Note,
			// RecipientName:               updateDRSI.Data.RecipientName,
			// MoneyReceived:               &updateDRSI.Data.MoneyReceived,
			// DeliveryEvidenceImageURL:    updateDRSI.Data.DeliveryEvidenceImageUrl,
			// TransactionEvidenceImageURL: updateDRSI.Data.TransactionEvidenceImageUrl,
			// ArrivalTime:                 &arrivalTime,
			// UnpunctualReason:            int8(updateDRSI.Data.UnpunctualReason),
			// UnpunctualDetail:            int8(updateDRSI.Data.UnpunctualDetail),
			// FarDeliveryReason:           updateDRSI.Data.FarDeliveryReason,
			// CreatedAt:                   &createdAt,
			// StartedAt:                   &startedAt,
			// FinishedAt:                  &finishedAt,
			// DeliveryRunSheetID:          updateDRSI.Data.DeliveryRunSheetId,
			// CourierID:                   updateDRSI.Data.CourierId,
			// SalesOrderID:                updateDRSI.Data.SalesOrderId,
		},
	}

	return
}

func (s *CourierAppService) StatusDelivery(ctx context.Context, req dto.CourierAppStatusDeliveryRequest) (res dto.CourierAppStatusDeliveryResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.StatusDelivery")
	defer span.End()

	var (
		deliveryRunSheetItem *logisticService.GetDeliveryRunSheetItemDetailResponse
		salesOrder           *bridgeService.GetSalesOrderGPListResponse
		// orderType            *bridgeService.GetOrderTypeDetailResponse
		wrt                 *bridgeService.GetWrtGPResponse
		merchantDeliveryLog *logisticService.GetFirstMerchantDeliveryLogResponse
		timeStartLimit      time.Time
		timeEndLimit        time.Time
		drsiArrivalTime     time.Time
	)

	// get delivery run sheet item
	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: req.Id,
	}); err != nil || deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// get sales order detail
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: deliveryRunSheetItem.Data.SalesOrderId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// get order type
	// if orderType, err = s.opt.Client.BridgeServiceGrpc.GetOrderTypeDetail(ctx, &bridgeService.GetOrderTypeDetailRequest{
	// 	Id: salesOrder.Data.OrderTypeId,
	// }); err != nil || orderType.Data == nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "order type")
	// 	return
	// }

	// self pickup will never be late/early and away from location
	// if orderType.Data.Id == 6 {
	// 	res = dto.CourierAppStatusDeliveryResponse{
	// 		Punctual: true,
	// 		Nearby:   true,
	// 	}
	// } else {
	// get wrt detail
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
		GnlRegion: salesOrder.Data[0].Wrt[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
	}

	// Punctuality test
	wrtStartClock := wrt.Data[0].Strttime[0:2]
	wrtStartClockInt, _ := strconv.Atoi(wrtStartClock)
	wrtEndClock := wrt.Data[0].Endtime[0:2]
	wrtEndClockInt, _ := strconv.Atoi(wrtEndClock)

	// get delivery date combined with wrt
	var layout string
	layout = "2006-01-02"
	deliveryDate, err := time.Parse(layout, salesOrder.Data[0].ReqShipDate)

	timeStartLimit = deliveryDate.Add(time.Duration(wrtStartClockInt) * time.Hour)
	timeEndLimit = deliveryDate.Add(time.Duration(wrtEndClockInt) * time.Hour)

	drsiArrivalTime = timex.ToLocTime(ctx, deliveryRunSheetItem.Data.ArrivalTime.AsTime())
	drsiArrivalTime, _ = time.Parse("2006-01-02 15:04:05", drsiArrivalTime.Format("2006-01-02 15:04:05"))

	// compare time with range of time limit
	if drsiArrivalTime.After(timeStartLimit) && drsiArrivalTime.Before(timeEndLimit) {
		res.Punctual = true
	} else {
		if drsiArrivalTime.Before(timeStartLimit) {
			res.Earlier = true
		}
	}

	// Proximity test
	if merchantDeliveryLog, err = s.opt.Client.LogisticServiceGrpc.GetFirstMerchantDeliveryLog(ctx, &logisticService.GetFirstMerchantDeliveryLogRequest{
		DeliveryRunSheetItemId: deliveryRunSheetItem.Data.Id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "merchant delivery log")
		return
	}

	if merchantDeliveryLog != nil && merchantDeliveryLog.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("merchant_delivery_log.required", "Foto pengiriman dibutuhkan")
		return
	}

	distance := distance(*merchantDeliveryLog.Data.Latitude, *merchantDeliveryLog.Data.Longitude, req.Latitude, req.Longitude)

	if distance < 100 {
		res.Nearby = true
	}

	return
}

func (s *CourierAppService) ActivateEmergency(ctx context.Context, req dto.CourierAppActivateEmergencyRequest) (res dto.CourierAppActivateEmergencyResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.ActivateEmergency")
	defer span.End()

	var (
		courier *bridgeService.GetCourierGPResponse
	)

	if courier, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPDetail(ctx, &bridgeService.GetCourierGPDetailRequest{
		Id: req.CourierId,
	}); err != nil || courier.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}

	if *courier.Data[0].GnlEmergencymode == 1 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("emergency_mode.invalid", "Anda sudah dalam keadaan emergency")
		return
	}

	if _, err = s.opt.Client.BridgeServiceGrpc.ActivateEmergencyCourier(ctx, &bridgeService.EmergencyCourierRequest{
		Id: req.CourierId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}

	return
}

func (s *CourierAppService) DeactivateEmergency(ctx context.Context, req dto.CourierAppDeactivateEmergencyRequest) (res dto.CourierAppDeactivateEmergencyResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.DeactivateEmergency")
	defer span.End()

	var (
		courier *bridgeService.GetCourierGPResponse
	)

	if courier, err = s.opt.Client.BridgeServiceGrpc.GetCourierGPDetail(ctx, &bridgeService.GetCourierGPDetailRequest{
		Id: req.CourierId,
	}); err != nil || courier.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}

	if *courier.Data[0].GnlEmergencymode == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("emergency_mode.invalid", "Anda sudah tidak dalam keadaan emergency")
		return
	}

	if _, err = s.opt.Client.BridgeServiceGrpc.DeactivateEmergencyCourier(ctx, &bridgeService.EmergencyCourierRequest{
		Id: req.CourierId,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "courier")
		return
	}

	return
}

func (s *CourierAppService) CreateMerchantDeliveryLog(ctx context.Context, req dto.CourierAppCreateMerchantDeliveryLogRequest) (res dto.CourierAppCreateMerchantDeliveryLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.CreateMerchantDeliveryLog")
	defer span.End()

	var (
		deliveryRunSheetItem      *logisticService.GetDeliveryRunSheetItemDetailResponse
		createMerchantDeliveryLog *logisticService.CreateMerchantDeliveryLogResponse
	)

	// get delivery run sheet item
	if deliveryRunSheetItem, err = s.opt.Client.LogisticServiceGrpc.GetDeliveryRunSheetItemDetail(ctx, &logisticService.GetDeliveryRunSheetItemDetailRequest{
		Id: req.Id,
	}); err != nil || deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	// check if DRSI is courier's job
	if deliveryRunSheetItem.Data.CourierId != req.CourierId {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorValidation("job.invalid", "Tugas ini bukan tugas anda")
		return
	}

	if createMerchantDeliveryLog, err = s.opt.Client.LogisticServiceGrpc.CreateMerchantDeliveryLog(ctx, &logisticService.CreateMerchantDeliveryLogRequest{
		Model: &logisticService.MerchantDeliveryLog{
			DeliveryRunSheetItemId: deliveryRunSheetItem.Data.Id,
			Latitude:               &req.Latitude,
			Longitude:              &req.Longitude,
			CreatedAt:              timestamppb.Now(),
		},
	}); err != nil || deliveryRunSheetItem.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "merchant delivery log")
		return
	}

	// update drsi
	if _, err = s.opt.Client.LogisticServiceGrpc.ArrivedDeliveryRunSheetItem(ctx, &logisticService.ArrivedDeliveryRunSheetItemRequest{
		Id: deliveryRunSheetItem.Data.Id,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("logistic", "delivery run sheet item")
		return
	}

	createdAt := createMerchantDeliveryLog.Data.CreatedAt.AsTime()
	res = dto.CourierAppCreateMerchantDeliveryLogResponse{
		MerchantDeliveryLog: &dto.GlobalMerchantDeliveryLog{
			Id:                     createMerchantDeliveryLog.Data.Id,
			Latitude:               createMerchantDeliveryLog.Data.Latitude,
			Longitude:              createMerchantDeliveryLog.Data.Longitude,
			CreatedAt:              &createdAt,
			DeliveryRunSheetItemId: createMerchantDeliveryLog.Data.DeliveryRunSheetItemId,
		},
	}

	return
}

func (s *CourierAppService) GetGlossary(ctx context.Context, req dto.CourierAppGetGlossaryRequest) (res []dto.CourierAppGetGlossaryResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierAppService.GetGlossary")
	defer span.End()

	var glossary *configurationService.GetGlossaryListResponse

	if glossary, err = s.opt.Client.ConfigurationServiceGrpc.GetGlossaryList(ctx, &configuration_service.GetGlossaryListRequest{
		Table:     req.Table,
		Attribute: req.Attribute,
		ValueInt:  int32(req.ValueInt),
		ValueName: req.ValueName,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("configuration", "glossary")
		return
	}

	for _, glossary := range glossary.Data {
		res = append(res, dto.CourierAppGetGlossaryResponse{
			Glossary: &dto.GlobalGlossary{
				ID:        int64(glossary.Id),
				Table:     glossary.Table,
				Attribute: glossary.Attribute,
				ValueInt:  int8(glossary.ValueInt),
				ValueName: glossary.ValueName,
				Note:      glossary.Note,
			},
		})
	}

	total = int64(len(glossary.Data))

	return
}

// Distance calculator
func distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}
