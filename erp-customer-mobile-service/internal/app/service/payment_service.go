package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/settlement_service"
)

type IPaymentService interface {
	GetPaymentMethod(ctx context.Context, req *dto.PaymentMethodRequestGet) (res []*dto.PaymentMethod, err error)
	GetInvoiceXendit(ctx context.Context, req *dto.PaymentRequest) (res *dto.InvoiceXenditModify, err error)
}

type PaymentService struct {
	opt opt.Options
}

func NewPaymentService() IPaymentService {
	return &PaymentService{
		opt: global.Setup.Common,
	}
}

func (s *PaymentService) GetPaymentMethod(ctx context.Context, req *dto.PaymentMethodRequestGet) (res []*dto.PaymentMethod, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentService.Get")
	defer span.End()

	paymentChannel, err := s.opt.Client.SalesServiceGrpc.GetPaymentChannelList(ctx, &sales_service.GetPaymentChannelListRequest{
		Status:     1,
		PublishIva: 1,
		// PublishFva:      1,
		PaymentMethodId: 2,
	})

	paymentGroupComb, err := s.opt.Client.SalesServiceGrpc.GetPaymentGroupCombList(ctx, &sales_service.GetPaymentGroupCombListRequest{
		PaymentGroupSls: "",
		TermPaymentSls:  req.Session.Customer.TermPaymentSls,
	})

	var bankTransferPaymentOptions []*dto.PaymentOption

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Get Bank Transfer options

	for _, paymentChannel := range paymentChannel.Data {
		bankTransferPaymentOptions = append(bankTransferPaymentOptions, &dto.PaymentOption{
			Name:            paymentChannel.Name,
			Value:           paymentChannel.Value,
			ImageURL:        paymentChannel.ImageUrl,
			PaymentGuideURL: paymentChannel.PaymentGuideUrl,
		})
	}

	for _, pgc := range paymentGroupComb.Data {
		if pgc.PaymentGroupSls != "Advance" {

			res = append(res, &dto.PaymentMethod{
				Name:        "Kredit Eden Farm " + pgc.TermPaymentSls,
				Value:       req.Session.Customer.TermPaymentSls,
				Description: "Silakan bayar belanjaan kamu sesuai waktu yang sudah ditentukan",
			})

			// TODO create conditions to check the user is eligible COD / BNS Payment
			res = append(res, &dto.PaymentMethod{
				Name:        "COD",
				Value:       "COD",
				Description: "Silakan bayar belanjaan kamu melalui kurir saat pesanan tiba",
				Note:        "(Bayar saat kurir datang)",
			})

		}

	}

	res = append(res, &dto.PaymentMethod{
		Name:           "Bank Transfer",
		Value:          "PBD",
		Note:           "(Dicek langsung)",
		Description:    "Segera bayar pesanan kamu, untuk menghindari pembatalan otomatis",
		PaymentOptions: bankTransferPaymentOptions,
	})

	return
}

func (s *PaymentService) GetInvoiceXendit(ctx context.Context, req *dto.PaymentRequest) (res *dto.InvoiceXenditModify, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PaymentService.Get")
	defer span.End()

	salesInvoiceExternalXendit, err := s.opt.Client.SettlementGrpc.GetSalesInvoiceExternalXendit(ctx, &settlement_service.GetSalesInvoiceExternalRequest{
		SalesOrderId: req.SalesOrderID,
	})

	paymentChannel, err := s.opt.Client.SalesServiceGrpc.GetPaymentChannelList(ctx, &sales_service.GetPaymentChannelListRequest{
		Status:          1,
		PublishIva:      1,
		PaymentMethodId: 2,
	})

	loc, _ := time.LoadLocation("Asia/Jakarta")
	deadlinePay := time.Unix(salesInvoiceExternalXendit.Data.ExpiryDate, 0)

	for _, pc := range paymentChannel.Data {
		if pc.Value == salesInvoiceExternalXendit.Data.BankCode {
			res = &dto.InvoiceXenditModify{
				ServerTime:      time.Now().In(loc),
				VaNumber:        salesInvoiceExternalXendit.Data.BankAccountNumber,
				PaymentNominal:  salesInvoiceExternalXendit.Data.Amount,
				DeadlinePayment: deadlinePay.In(loc),
				ImageUrl:        pc.ImageUrl,
				Name:            pc.Name,
				PaymentGuideUrl: pc.PaymentGuideUrl,
			}

		}
	}

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}
