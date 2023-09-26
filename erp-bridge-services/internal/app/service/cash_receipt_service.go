package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ICashReceiptService interface {
	GetGP(ctx context.Context, req *bridge_service.GetCashReceiptListRequest) (res *bridge_service.GetCashReceiptListResponse, err error)
	Create(ctx context.Context, req *bridge_service.CreateCashReceiptRequest) (res *bridge_service.CreateCashReceiptResponse, err error)
	// GetGPByID(ctx context.Context, req *bridge_service.GetSalesOrderGPListByIDRequest) (res *bridge_service.GetSalesOrderGPListResponse, err error)
}

type CashReceiptService struct {
	opt opt.Options
	// RepositoryCashReceipt repository.ICashReceiptRepository
}

func NewCashReceiptService() ICashReceiptService {
	return &CashReceiptService{
		opt: global.Setup.Common,
		// RepositoryCashReceipt: repository.NewCashReceiptRepository(),
	}
}

func (s *CashReceiptService) GetGP(ctx context.Context, req *bridge_service.GetCashReceiptListRequest) (res *bridge_service.GetCashReceiptListResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CashReceiptServices.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Custnmbr != "" {
		params["custnmbr"] = req.Custnmbr
	}

	if req.Sopnumbe != "" { //(this is si number)
		params["sopnumbe"] = req.Sopnumbe
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
	if req.SiDocdateFrom != "" {
		params["si_docdate_from"] = req.SiDocdateFrom
	}
	if req.SiDocdateTo != "" {
		params["si_docdate_to"] = req.SiDocdateTo
	}
	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}
	if req.Docnumber != "" {
		params["docnumbr"] = req.Docnumber
	}
	if req.Status != "" {
		params["status"] = req.Status
	}
	if req.GnlRegion != "" {
		params["gnl_region"] = req.GnlRegion
	}
	err = global.HttpRestApiToMicrosoftGP("GET", "CashReceipt/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *CashReceiptService) Create(ctx context.Context, req *bridge_service.CreateCashReceiptRequest) (res *bridge_service.CreateCashReceiptResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CashReceiptServices.GetGP")
	defer span.End()

	reqCreateCashReceipt := &dto.CreateCashReceiptRequest{
		Interid:        global.EnvDatabaseGP,
		Sopnumbe:       req.Sopnumbe,
		AmountReceived: req.AmountReceived,
		Chekbkid:       req.Chekbkid,
		Docdate:        req.Docdate,
	}
	// add return payment code in response
	err = global.HttpRestApiToMicrosoftGP("POST", "sales/createpayment", reqCreateCashReceipt, &res, nil)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
