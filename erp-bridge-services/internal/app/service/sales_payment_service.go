package service

import (
	"context"
	"errors"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/sirupsen/logrus"
)

type ISalesPaymentService interface {
	Get(ctx context.Context, salesInvoiceID int64) (res []dto.SalesPaymentResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.SalesPaymentResponse, err error)
	GetGP(ctx context.Context, req *pb.GetSalesPaymentGPListRequest) (res *pb.GetSalesPaymentGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetSalesPaymentGPDetailRequest) (res *pb.GetSalesPaymentGPResponse, err error)

	// GP Integrated
	CreateGP(ctx context.Context, req *pb.CreateSalesPaymentGPRequest) (res *pb.CreateSalesInvoiceGPResponse, err error)
	CreateGPnonPBD(ctx context.Context, req *pb.CreateSalesPaymentGPnonPBDRequest) (res *pb.CreateSalesPaymentGPnonPBDResponse, err error)
}

type SalesPaymentService struct {
	opt                    opt.Options
	RepositorySalesPayment repository.ISalesPaymentRepository
}

func NewSalesPaymentService() ISalesPaymentService {
	return &SalesPaymentService{
		opt:                    global.Setup.Common,
		RepositorySalesPayment: repository.NewSalesPaymentRepository(),
	}
}

func (s *SalesPaymentService) Get(ctx context.Context, salesInvoiceID int64) (res []dto.SalesPaymentResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.Get")
	defer span.End()

	var salesPayments []*model.SalesPayment
	salesPayments, total, err = s.RepositorySalesPayment.Get(ctx, salesInvoiceID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, salesPayment := range salesPayments {
		res = append(res, dto.SalesPaymentResponse{
			ID:              salesPayment.ID,
			Code:            salesPayment.Code,
			Status:          salesPayment.Status,
			Amount:          salesPayment.Amount,
			RecognitionDate: salesPayment.RecognitionDate,
			ReceivedDate:    salesPayment.ReceivedDate,
		})
	}

	return
}

func (s *SalesPaymentService) GetDetail(ctx context.Context, id int64, code string) (res dto.SalesPaymentResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.GetDetail")
	defer span.End()

	var salesPayment *model.SalesPayment
	salesPayment, err = s.RepositorySalesPayment.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.SalesPaymentResponse{
		ID:              salesPayment.ID,
		Code:            salesPayment.Code,
		Status:          salesPayment.Status,
		Amount:          salesPayment.Amount,
		RecognitionDate: salesPayment.RecognitionDate,
		ReceivedDate:    salesPayment.ReceivedDate,
	}

	return
}

func (s *SalesPaymentService) CreateGP(ctx context.Context, req *pb.CreateSalesPaymentGPRequest) (res *pb.CreateSalesInvoiceGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.CreateGP")
	defer span.End()
	var tempRes *dto.CreateSalesPaymentGPResponse
	req.Interid = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("POST", "sales/createpayment", req, &tempRes, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}
	res = &pb.CreateSalesInvoiceGPResponse{
		Code:     tempRes.Code,
		Message:  tempRes.Message,
		Sopnumbe: tempRes.Docnumber,
	}

	return
}

func (s *SalesPaymentService) CreateGPnonPBD(ctx context.Context, req *pb.CreateSalesPaymentGPnonPBDRequest) (res *pb.CreateSalesPaymentGPnonPBDResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.CreateGP")
	defer span.End()

	var tempRes *dto.CreateSalesPaymentGPnonPBDResponse
	var reqAPI *dto.CreateSalesPaymentGPnonPBDRequest
	var applyTo []*dto.CreateSalesPaymentGPnonPBDRequest_ApplyTo
	applyTo = append(applyTo, &dto.CreateSalesPaymentGPnonPBDRequest_ApplyTo{
		Sopnumbe:    req.ApplyTo[0].Sopnumbe,
		ApplyAmount: req.ApplyTo[0].ApplyAmount,
	})
	// reqAPI.Interid = global.EnvDatabaseGP
	reqAPI = &dto.CreateSalesPaymentGPnonPBDRequest{
		Interid:  global.EnvDatabaseGP,
		Bachnumb: "",
		Docnumbr: "",
		Custnmbr: req.Custnmbr,
		Docdate:  req.Docdate,
		Cshrctyp: req.Cshrctyp,
		Curncyid: "",
		Chekbkid: req.Chekbkid,
		Ortrxamt: req.Ortrxamt,
		Trxdscrn: "",
		ApplyTo:  applyTo,
	}
	err = global.HttpRestApiToMicrosoftGP("POST", "cashreceipt/create", reqAPI, &tempRes, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}
	res = &pb.CreateSalesPaymentGPnonPBDResponse{
		Code:          tempRes.Code,
		Message:       tempRes.Message,
		Paymentnumber: tempRes.Paymentnumber,
	}

	return
}

func (s *SalesPaymentService) GetGP(ctx context.Context, req *pb.GetSalesPaymentGPListRequest) (res *pb.GetSalesPaymentGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Custnmbr != "" {
		params["custnmbr"] = req.Custnmbr
	}

	if req.Docnumbr != "" {
		params["docnumbr"] = req.Docnumbr
	}

	if req.Status != "" {
		params["dcstatus"] = req.Status
	}

	if req.Sopnumbe != "" {
		params["sopnumbe"] = req.Sopnumbe
	}

	if req.GnlRegion != "" {
		params["gnl_region"] = req.GnlRegion
	}

	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}

	if req.DocdateFrom != "" {
		params["docdate_from"] = req.DocdateFrom
	}

	if req.DocdateTo != "" {
		params["docdate_to"] = req.DocdateTo
	}

	if req.SiDocdateFrom != "" {
		params["si_docdate_from"] = req.SiDocdateFrom
	}

	if req.SiDocdateTo != "" {
		params["si_docdate_to"] = req.SiDocdateTo
	}

	if req.SoDocdateFrom != "" {
		params["so_docdate_from"] = req.SoDocdateFrom
	}

	if req.SoDocdateTo != "" {
		params["so_docdate_to"] = req.SoDocdateTo
	}

	if req.OrderBy != "" {
		params["orderby"] = req.OrderBy
	}
	if req.GnlCustTypeId != "" {
		params["gnl_cust_type_id"] = req.GnlCustTypeId
	}
	err = global.HttpRestApiToMicrosoftGP("GET", "CashReceipt/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *SalesPaymentService) GetDetailGP(ctx context.Context, req *pb.GetSalesPaymentGPDetailRequest) (res *pb.GetSalesPaymentGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPaymentService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "CashReceipt/detail", nil, &res, params)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
