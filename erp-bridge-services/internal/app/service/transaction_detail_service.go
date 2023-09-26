package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ITransactionDetailService interface {
	GetGP(ctx context.Context, req *pb.GetTransactionDetailGPListRequest) (res *pb.GetTransactionDetailGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetTransactionDetailGPDetailRequest) (res *pb.GetTransactionDetailGPResponse, err error)
}

type TransactionDetailService struct {
	opt opt.Options
}

func NewTransactionDetailService() ITransactionDetailService {
	return &TransactionDetailService{
		opt: global.Setup.Common,
	}
}

func (s *TransactionDetailService) GetGP(ctx context.Context, req *pb.GetTransactionDetailGPListRequest) (res *pb.GetTransactionDetailGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransactionDetailService.GetGP")
	defer span.End()

	// params := map[string]string{
	// 	"interid":    global.EnvDatabaseGP,
	// 	"PageNumber": strconv.Itoa(int(req.Offset)),
	// 	"PageSize":   strconv.Itoa(int(req.Limit)),
	// }

	// err = global.HttpRestApiToMicrosoftGP("GET", "transactionDetail/getall", nil, &res, params)

	// if err != nil {
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	res = &pb.GetTransactionDetailGPResponse{
		PageNumber:   1,
		PageSize:     1,
		TotalPages:   1,
		TotalRecords: 1,
		Data: []*pb.TransactionDetailGP{
			{
				SopNumber:      "DUMMY SOP NUMBER",
				ItemNumber:     "DUMMY ITEM NUMBER",
				UomId:          "DUMMY UOM ID",
				OrderQuantity:  1,
				UnitPrice:      10000,
				ExtendedPrice:  0,
				ItemDesc:       "DUMMY ITEM DESCRIPTION",
				Markdown:       0,
				UnitCost:       5000,
				ReqShipDate:    "01-01-2023",
				DateShipped:    "01-01-2023",
				ProductPush:    false,
				QtyOrdered:     1,
				QtyFulfilled:   1,
				QtyCanceled:    0,
				QtyBackOrder:   0,
				TotalWeight:    1,
				SiteId:         "DUMMY SITE ID",
				PriceLevelId:   "DUMMY PRICE LEVEL ID",
				ShipToAddress:  "01-01-2023",
				ShippingMethod: "DUMMY SHIPPING METHOD",
				QtyAvailable:   5,
			},
		},
		Succeeded: true,
		Errors:    []string{},
		Message:   "",
	}

	return
}

func (s *TransactionDetailService) GetDetailGP(ctx context.Context, req *pb.GetTransactionDetailGPDetailRequest) (res *pb.GetTransactionDetailGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransactionDetailService.GetDetailGP")
	defer span.End()

	// params := map[string]string{
	// 	"interid": global.EnvDatabaseGP,
	// 	"type":    req.Type,
	// 	"id":      req.Id,
	// }

	// err = global.HttpRestApiToMicrosoftGP("GET", "transactionDetail/getbyid", nil, &res, params)

	// if err != nil {
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	res = &pb.GetTransactionDetailGPResponse{
		PageNumber:   1,
		PageSize:     1,
		TotalPages:   1,
		TotalRecords: 1,
		Data: []*pb.TransactionDetailGP{
			{
				SopNumber:      "DUMMY SOP NUMBER",
				ItemNumber:     "DUMMY ITEM NUMBER",
				UomId:          "DUMMY UOM ID",
				OrderQuantity:  1,
				UnitPrice:      10000,
				ExtendedPrice:  0,
				ItemDesc:       "DUMMY ITEM DESCRIPTION",
				Markdown:       0,
				UnitCost:       5000,
				ReqShipDate:    "01-01-2023",
				DateShipped:    "01-01-2023",
				ProductPush:    false,
				QtyOrdered:     1,
				QtyFulfilled:   1,
				QtyCanceled:    0,
				QtyBackOrder:   0,
				TotalWeight:    1,
				SiteId:         "DUMMY SITE ID",
				PriceLevelId:   "DUMMY PRICE LEVEL ID",
				ShipToAddress:  "01-01-2023",
				ShippingMethod: "DUMMY SHIPPING METHOD",
				QtyAvailable:   5,
			},
		},
		Succeeded: true,
		Errors:    []string{},
		Message:   "",
	}

	return
}
