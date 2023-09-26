package service

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ICashReceiptService interface {
	Get(ctx context.Context, req *dto.RequestGetHistoryTransaction) (res dto.DataGetHistoryTransaction, err error)
}

type CashReceiptService struct {
	opt opt.Options
}

func NewCashReceiptService() ICashReceiptService {
	return &CashReceiptService{
		opt: global.Setup.Common,
	}
}

func (s *CashReceiptService) Get(ctx context.Context, req *dto.RequestGetHistoryTransaction) (res dto.DataGetHistoryTransaction, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CashReceiptService.Get")
	defer span.End()

	cashReceipt, err := s.opt.Client.BridgeServiceGrpc.GetCashReceiptList(ctx, &bridge_service.GetCashReceiptListRequest{
		Limit:    100,
		Offset:   0,
		Custnmbr: req.Session.Customer.Code,
	})
	fmt.Println(cashReceipt)

	return
}
