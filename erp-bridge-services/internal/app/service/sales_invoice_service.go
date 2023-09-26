package service

import (
	"context"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ISalesInvoiceService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, addressID int64, customerID int64, salespersonID int64, orderDateFrom time.Time, orderDateTo time.Time) (res []dto.SalesInvoiceResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.SalesInvoiceResponse, err error)
	GetGP(ctx context.Context, req *bridge_service.GetSalesInvoiceGPListRequest) (res *bridge_service.GetSalesInvoiceGPListResponse, err error)
	CreateGP(ctx context.Context, req *pb.CreateSalesInvoiceGPRequest) (res *pb.CreateSalesInvoiceGPResponse, err error)
	GetDetailGP(ctx context.Context, req *bridge_service.GetSalesInvoiceGPDetailRequest) (res *bridge_service.GetSalesInvoiceGPDetailResponse, err error)
}

type SalesInvoiceService struct {
	opt                    opt.Options
	RepositorySalesInvoice repository.ISalesInvoiceRepository
}

func NewSalesInvoiceService() ISalesInvoiceService {
	return &SalesInvoiceService{
		opt:                    global.Setup.Common,
		RepositorySalesInvoice: repository.NewSalesInvoiceRepository(),
	}
}

func (s *SalesInvoiceService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, addressID int64, customerID int64, salespersonID int64, orderDateFrom time.Time, orderDateTo time.Time) (res []dto.SalesInvoiceResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceService.Get")
	defer span.End()

	var salesInvoices []*model.SalesInvoice
	salesInvoices, total, err = s.RepositorySalesInvoice.Get(ctx, offset, limit, status, search, orderBy, addressID, customerID, salespersonID, orderDateFrom, orderDateTo)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesInvoice := range salesInvoices {
		res = append(res, dto.SalesInvoiceResponse{
			ID:            salesInvoice.ID,
			Code:          salesInvoice.Code,
			Status:        salesInvoice.Status,
			DeliveryFee:   salesInvoice.DeliveryFee,
			VouDiscAmount: salesInvoice.VouDiscAmount,
			TotalPrice:    salesInvoice.TotalPrice,
			TotalCharge:   salesInvoice.TotalCharge,
		})
	}

	return
}

func (s *SalesInvoiceService) GetDetail(ctx context.Context, id int64, code string) (res dto.SalesInvoiceResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceService.GetDetail")
	defer span.End()

	var salesInvoice *model.SalesInvoice
	salesInvoice, err = s.RepositorySalesInvoice.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesInvoiceResponse{
		ID:            salesInvoice.ID,
		Code:          salesInvoice.Code,
		Status:        salesInvoice.Status,
		DeliveryFee:   salesInvoice.DeliveryFee,
		VouDiscAmount: salesInvoice.VouDiscAmount,
		TotalPrice:    salesInvoice.TotalPrice,
		TotalCharge:   salesInvoice.TotalCharge,
	}

	return
}

func (s *SalesInvoiceService) GetGP(ctx context.Context, req *bridge_service.GetSalesInvoiceGPListRequest) (res *bridge_service.GetSalesInvoiceGPListResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Custnumber != "" {
		params["custnmbr"] = req.Custnumber
	}
	// if req.Custname != "" {
	// 	params["custname"] = req.Custname
	// }
	//
	if req.SoNumber != "" {
		params["orignumbr"] = req.SoNumber
	}
	if req.DeltaUser != "" {
		params["delta_user"] = req.DeltaUser
	}
	if req.SiNumber != "" {
		params["sopnumbe"] = req.SiNumber
	}
	if req.DocdateFrom != "" {
		params["docdate_from"] = req.DocdateFrom
	}
	if req.DocdateTo != "" {
		params["docdate_to"] = req.DocdateTo
	}
	if req.SoDocdateFrom != "" {
		params["so_docdate_from"] = req.SoDocdateFrom
	}
	if req.SoDocdateTo != "" {
		params["so_docdate_to"] = req.SoDocdateTo
	}
	if req.Region != "" {
		params["region"] = req.Region
	}
	if req.Ordertype != "" {
		params["ordertype"] = req.Ordertype
	}
	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}
	if req.Status != "" {
		params["sopstatus"] = req.Status
	}
	if req.OrderBy != "" {
		params["orderby"] = req.OrderBy
	}
	if req.GnlCustTypeId != "" {
		params["gnl_cust_type_id"] = req.GnlCustTypeId
	}
	if req.RemainingAmountFlag != "" {
		params["remaining_amount_flag"] = req.RemainingAmountFlag
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "SalesInvoice/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesInvoiceService) GetDetailGP(ctx context.Context, req *bridge_service.GetSalesInvoiceGPDetailRequest) (res *bridge_service.GetSalesInvoiceGPDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "AddressService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "SalesInvoice/detail", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesInvoiceService) CreateGP(ctx context.Context, req *pb.CreateSalesInvoiceGPRequest) (res *pb.CreateSalesInvoiceGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesInvoiceService.CreateGP")
	defer span.End()

	req.Interid = global.EnvDatabaseGP
	var payloadGP *dto.CreateSalesInvoiceGPRequest
	payloadGP = &dto.CreateSalesInvoiceGPRequest{
		Interid:  global.EnvDatabaseGP,
		Orignumb: req.Orignumb,
		Sopnumbe: req.Sopnumbe,
		Docid:    req.Docid,
		Docdate:  req.Docdate,
		Custnmbr: req.Custnmbr,
		Custname: req.Custname,
		Prstadcd: req.Prstadcd,
		Curncyid: req.Curncyid,
		Subtotal: req.Subtotal,
		Trdisamt: req.Trdisamt,
		Freight:  req.Freight,
		Miscamnt: req.Miscamnt,
		Taxamnt:  req.Taxamnt,
		Docamnt:  req.Docamnt,
		AmountReceived: &dto.CreateSalesInvoiceGPRequest_AmountReceived{
			Amount:   req.AmountReceived.Amount,
			Chekbkid: req.AmountReceived.Chekbkid,
		},
		GnlRequestShipDate: req.GnlRequestShipDate,
		GnlRegion:          req.GnlRegion,
		GnlWrtId:           req.GnlWrtId,
		GnlArchetypeId:     req.GnlArchetypeId,
		GnlOrderChannel:    req.GnlOrderChannel,
		GnlSoCodeApps:      req.GnlSoCodeApps,
		GnlTotalweight:     req.GnlTotalweight,
		Userid:             req.Userid,
		Shipmthd:           req.Shipmthd,
		Locncode:           req.Locncode,
		Pymtrmid:           req.Pymtrmid,
	}
	if len(req.VoucherApply) == 0 {
		payloadGP.VoucherApply = []*dto.CreateSalesInvoiceGPRequest_VoucherApply{}
	} else {
		for _, v := range req.VoucherApply {
			payloadGP.VoucherApply = append(payloadGP.VoucherApply, &dto.CreateSalesInvoiceGPRequest_VoucherApply{
				GnlVoucherType: v.GnlVoucherType,
				GnlVoucherId:   v.GnlVoucherId,
				Ordocamt:       v.Ordocamt,
			})
		}

	}

	for _, v := range req.Detailitems {
		payloadGP.Detailitems = append(payloadGP.Detailitems, &dto.CreateSalesInvoiceGPRequest_DetailItem{
			Lnitmseq:   v.Lnitmseq,
			Itemnmbr:   v.Itemnmbr,
			Locncode:   v.Locncode,
			Uofm:       v.Uofm,
			Pricelvl:   v.Pricelvl,
			Quantity:   v.Quantity,
			Unitprce:   v.Unitprce,
			Xtndprce:   v.Xtndprce,
			GnL_Weight: v.GnL_Weight,
		})
	}
	err = global.HttpRestApiToMicrosoftGP("POST", "sales/createinvoice", payloadGP, &res, nil)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
