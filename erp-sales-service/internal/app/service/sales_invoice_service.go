package service

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	salesService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/sales_service"
	"git.edenfarm.id/project-version3/erp-services/erp-sales-service/global"
)

type ISalesInvoiceService interface {
	GetSalesInvoiceGPMobileList(ctx context.Context, req *salesService.GetSalesInvoiceGPMobileListRequest) (res salesService.GetSalesInvoiceGPMobileListResponse, err error)
}

type SalesInvoiceService struct {
	opt opt.Options
}

func NewSalesInvoiceService() ISalesInvoiceService {
	return &SalesInvoiceService{
		opt: global.Setup.Common,
	}
}

func (s *SalesInvoiceService) GetSalesInvoiceGPMobileList(ctx context.Context, req *salesService.GetSalesInvoiceGPMobileListRequest) (res salesService.GetSalesInvoiceGPMobileListResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "SalesPersonService.Get")
	defer span.End()

	siGP, err := s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridge_service.GetSalesInvoiceGPListRequest{
		Limit:      req.Limit,
		Offset:     req.Offset,
		Custnumber: req.Custnumber,
		SoNumber:   req.SoNumber,
	})
	if err != nil {

	}
	// var res []dto.ResponseSIMobile

	for _, v := range siGP.Data {
		fmt.Print(v)
		siDetail, err := s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPDetail(ctx, &bridge_service.GetSalesInvoiceGPDetailRequest{
			Id: v.Sopnumbe,
		})
		fmt.Println(siDetail, err)
		SalesInvoice := &sales_service.SalesInvoice{
			InvoiceId:         v.Sopnumbe,
			Id:                v.Sopnumbe,
			OrderCode:         v.SalesOrder[0].Orignumb,
			OrderDate:         "2000-01-01",
			InvoiceDate:       v.Docdate,
			TotalPrice:        siDetail.Data[0].Subtotal,
			DeliveryFee:       siDetail.Data[0].Frtamnt,
			InvoiceCode:       v.Sopnumbe,
			VoucherAmount:     siDetail.Data[0].Trdisamt,
			TotalCharge:       siDetail.Data[0].Ordocamt,
			PointRedeemAmount: 0,
			AdjustmentAmount:  0,
		}
		var tempItem []*salesService.InvoiceItem
		for _, w := range v.Details {
			fmt.Print(w)
			tempItem = append(tempItem, &salesService.InvoiceItem{
				ItemId:             w.Itemnmbr,
				ItemName:           w.Itemdesc,
				InvoiceQty:         w.Quantity,
				UomName:            w.Uofm,
				UnitPrice:          w.Unitprce,
				Subtotal:           w.Unitprce * w.Quantity,
				ItemDiscountAmount: w.Xtndprce,
			})
		}
		fmt.Println(SalesInvoice, tempItem)
		res.Data = append(res.Data, &salesService.SalesInvoiceMobile{
			SalesInvoice: SalesInvoice,
			InvoiceItem:  tempItem,
		})

	}
	return
}
