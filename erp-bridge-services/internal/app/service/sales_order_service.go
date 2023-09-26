package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/sirupsen/logrus"
)

type ISalesOrderService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, addressID int64, customerID int64, salespersonID int64, orderDateFrom time.Time, orderDateTo time.Time) (res []dto.SalesOrderResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.SalesOrderResponse, err error)
	CreateSalesOrder(ctx context.Context, req *bridge_service.CreateSalesOrderRequest) (res dto.SalesOrderResponse, err error)
	CreateSalesOrderGP(ctx context.Context, req *bridge_service.CreateSalesOrderGPRequest) (res dto.CommonGPResponse, err error)
	GetGP(ctx context.Context, req *bridge_service.GetSalesOrderGPListRequest) (res *bridge_service.GetSalesOrderGPListResponse, err error)
	GetGPByID(ctx context.Context, req *bridge_service.GetSalesOrderGPListByIDRequest) (res *bridge_service.GetSalesOrderGPListResponse, err error)
	GetSalesMovementGP(ctx context.Context, req *bridge_service.GetSalesMovementGPRequest) (res *bridge_service.GetSalesMovementGPResponse, err error)
}

type SalesOrderService struct {
	opt                  opt.Options
	RepositorySalesOrder repository.ISalesOrderRepository
}

func NewSalesOrderService() ISalesOrderService {
	return &SalesOrderService{
		opt:                  global.Setup.Common,
		RepositorySalesOrder: repository.NewSalesOrderRepository(),
	}
}

func (s *SalesOrderService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, addressID int64, customerID int64, salespersonID int64, orderDateFrom time.Time, orderDateTo time.Time) (res []dto.SalesOrderResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.Get")
	defer span.End()

	var salesOrders []*model.SalesOrder
	salesOrders, total, err = s.RepositorySalesOrder.Get(ctx, offset, limit, status, search, orderBy, addressID, customerID, salespersonID, orderDateFrom, orderDateTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesOrder := range salesOrders {
		res = append(res, dto.SalesOrderResponse{
			ID:            salesOrder.ID,
			Code:          salesOrder.Code,
			DocNumber:     salesOrder.DocNumber,
			AddressID:     salesOrder.AddressID,
			CustomerID:    salesOrder.CustomerID,
			SalespersonID: salesOrder.SalespersonID,
			WrtID:         salesOrder.WrtID,
			OrderTypeID:   salesOrder.OrderTypeID,
			SiteID:        salesOrder.SiteID,
			Application:   salesOrder.Application,
			Status:        salesOrder.Status,
			OrderDate:     timex.ToLocTime(ctx, salesOrder.OrderDate),
			Total:         salesOrder.Total,
			StatusConvert: statusx.ConvertStatusValue(salesOrder.Status),
			CreatedDate:   timex.ToLocTime(ctx, salesOrder.CreatedDate),
			ModifiedDate:  timex.ToLocTime(ctx, salesOrder.ModifiedDate),
			FinishedDate:  timex.ToLocTime(ctx, salesOrder.FinishedDate),
			CreatedAt:     timex.ToLocTime(ctx, salesOrder.CreatedAt),
			UpdatedAt:     timex.ToLocTime(ctx, salesOrder.UpdatedAt),
		})
	}

	return
}

func (s *SalesOrderService) GetDetail(ctx context.Context, id int64, code string) (res dto.SalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.GetDetail")
	defer span.End()

	var salesOrder *model.SalesOrder
	salesOrder, err = s.RepositorySalesOrder.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesOrderResponse{
		ID:            salesOrder.ID,
		Code:          salesOrder.Code,
		DocNumber:     salesOrder.DocNumber,
		AddressID:     salesOrder.AddressID,
		CustomerID:    salesOrder.CustomerID,
		SalespersonID: salesOrder.SalespersonID,
		WrtID:         salesOrder.WrtID,
		OrderTypeID:   salesOrder.OrderTypeID,
		SiteID:        salesOrder.SiteID,
		Application:   salesOrder.Application,
		Status:        salesOrder.Status,
		OrderDate:     timex.ToLocTime(ctx, salesOrder.OrderDate),
		Total:         salesOrder.Total,
		StatusConvert: statusx.ConvertStatusValue(salesOrder.Status),
		CreatedDate:   timex.ToLocTime(ctx, salesOrder.CreatedDate),
		ModifiedDate:  timex.ToLocTime(ctx, salesOrder.ModifiedDate),
		FinishedDate:  timex.ToLocTime(ctx, salesOrder.FinishedDate),
		CreatedAt:     timex.ToLocTime(ctx, salesOrder.CreatedAt),
		UpdatedAt:     timex.ToLocTime(ctx, salesOrder.UpdatedAt),
	}

	return
}

func (s *SalesOrderService) CreateSalesOrder(ctx context.Context, req *bridge_service.CreateSalesOrderRequest) (res dto.SalesOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.CreateSalesOrder")
	defer span.End()

	//var salesOrder *model.SalesOrder
	reqRest := dto.SalesOrderResponse{
		ID:        req.Data.Id,
		Code:      req.Data.Code,
		DocNumber: req.Data.DocNumber,
		// AddressID:     req.Data.AddressID,
		// CustomerID:    req.Data.CustomerID,
		// SalespersonID: req.Data.SalespersonID,
		// WrtID:         req.Data.WrtID,
		// OrderTypeID:   req.Data.OrderTypeID,
		// SiteID:        req.Data.SiteID,
		// Application:   req.Data.Application,
		// Status:        req.Data.Status,
		// OrderDate:     timex.ToLocTime(ctx, req.Data.OrderDate),
		Total: req.Data.Total,
		// StatusConvert: statusx.ConvertStatusValue(req.Data.Status),
		// CreatedDate:   timex.ToLocTime(ctx, req.Data.CreatedDate),
		// ModifiedDate:  timex.ToLocTime(ctx, req.Data.ModifiedDate),
		// FinishedDate:  timex.ToLocTime(ctx, req.Data.FinishedDate),
		// CreatedAt:     timex.ToLocTime(ctx, req.Data.CreatedAt),
		// UpdatedAt:     timex.ToLocTime(ctx, req.Data.UpdatedAt),
	}
	payload, _ := json.Marshal(reqRest)
	fmt.Println("payload address : ", payload)

	res = reqRest

	return
}

func (s *SalesOrderService) GetGP(ctx context.Context, req *bridge_service.GetSalesOrderGPListRequest) (res *bridge_service.GetSalesOrderGPListResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Custnumber != "" {
		params["custnumber"] = req.Custnumber
	}
	if req.Custname != "" {
		params["custname"] = url.PathEscape(req.Custname)
	}
	if req.SoNumber != "" {
		params["sopnumbe"] = req.SoNumber
	}
	if req.DocdateFrom != "" {
		params["docdate_from"] = req.DocdateFrom
	}
	if req.DocdateTo != "" {
		params["docdate_to"] = req.DocdateTo
	}
	if req.ReqShipDateFrom != "" {
		params["reqshipdate_from"] = req.ReqShipDateFrom
	}
	if req.ReqShipDateTo != "" {
		params["reqshipdate_to"] = req.ReqShipDateTo
	}
	if req.Region != "" {
		params["gnl_region"] = req.Region
	}
	if req.OrderChannel != "" {
		params["gnl_order_channel"] = req.OrderChannel
	}
	if req.SoCodeApps != "" {
		params["gnl_so_code_apps"] = req.SoCodeApps
	}
	if req.WrtId != "" {
		params["gnl_wrt_id"] = req.WrtId
	}
	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}
	if req.Docid != "" {
		params["docid"] = req.Docid
	}
	if req.SalespersonId != "" {
		params["slprsnid"] = req.SalespersonId
	}
	if req.Status != "" {
		params["status"] = req.Status
	}

	if req.Orderby != "" {
		params["orderby"] = req.Orderby
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "SalesOrder/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesOrderService) GetGPByID(ctx context.Context, req *bridge_service.GetSalesOrderGPListByIDRequest) (res *bridge_service.GetSalesOrderGPListResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "SalesOrder/detail", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesOrderService) CreateSalesOrderGP(ctx context.Context, req *bridge_service.CreateSalesOrderGPRequest) (res dto.CommonGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesOrderService.CreateSalesOrderGP")
	defer span.End()
	reqCreateSOGP := &dto.CreateSalesOrderGPRequest{
		Interid:            global.EnvDatabaseGP,
		Sopnumbe:           req.Sopnumbe,
		Docid:              req.Docid,
		Docdate:            req.Docdate,
		Custnmbr:           req.Custnmbr,
		Custname:           req.Custname,
		Prstadcd:           req.Prstadcd,
		Curncyid:           req.Curncyid,
		Subtotal:           req.Subtotal,
		Trdisamt:           req.Trdisamt,
		Freight:            req.Freight,
		Miscamnt:           req.Miscamnt,
		Taxamnt:            req.Taxamnt,
		Docamnt:            req.Docamnt,
		GnlRequestShipDate: req.GnlRequestShipDate,
		GnlRegion:          req.GnlRegion,
		GnlWrtId:           req.GnlWrtId,
		GnlArchetypeId:     req.GnlArchetypeId,
		GnlOrderChannel:    req.GnlOrderChannel,
		GnlSoCodeApps:      req.GnlSoCodeApps,
		GnlTotalweight:     req.GnlTotalweight,
		Userid:             req.Userid,
		Locncode:           req.Locncode,
		Shipmthd:           req.Shipmthd,
		Pymtrmid:           req.Pymtrmid,
		Note:               req.Note,
	}
	for _, v := range req.Detailitems {
		reqCreateSOGP.Detailitems = append(reqCreateSOGP.Detailitems, &dto.CreateSalesOrderGPRequest_DetailItem{
			Sopnumbe:   v.Sopnumbe,
			Itemnmbr:   v.Itemnmbr,
			Itemdesc:   v.Itemdesc,
			Locncode:   v.Locncode,
			Uofm:       v.Uofm,
			Pricelvl:   v.Pricelvl,
			Quantity:   v.Quantity,
			Unitprce:   v.Unitprce,
			Xtndprce:   v.Xtndprce,
			GnL_Weight: v.GnL_Weight,
		})
	}

	for _, v := range req.VoucherApply {
		reqCreateSOGP.VoucherApply = append(reqCreateSOGP.VoucherApply, &dto.VoucherApplyRequest{
			GnlVoucherType: int8(v.GnlVoucherType),
			GnlVoucherID:   v.GnlVoucherId,
			Ordocamt:       v.Ordocamt,
		})
	}

	err = global.HttpRestApiToMicrosoftGP("POST", "sales/createorder", reqCreateSOGP, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		logrus.Error("Error Login: " + res.Message)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}

func (s *SalesOrderService) GetSalesMovementGP(ctx context.Context, req *bridge_service.GetSalesMovementGPRequest) (res *bridge_service.GetSalesMovementGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.SoNumber != "" {
		params["sopnumbe"] = req.SoNumber
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "SalesOrder/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
