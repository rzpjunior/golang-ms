package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type ITransactionListService interface {
	GetGP(ctx context.Context, req *pb.GetTransactionListGPListRequest) (res *pb.GetTransactionListGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetTransactionListGPDetailRequest) (res *pb.GetTransactionListGPResponse, err error)
}

type TransactionListService struct {
	opt opt.Options
}

func NewTransactionListService() ITransactionListService {
	return &TransactionListService{
		opt: global.Setup.Common,
	}
}

func (s *TransactionListService) GetGP(ctx context.Context, req *pb.GetTransactionListGPListRequest) (res *pb.GetTransactionListGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransactionListService.GetGP")
	defer span.End()

	// params := map[string]string{
	// 	"interid":    global.EnvDatabaseGP,
	// 	"PageNumber": strconv.Itoa(int(req.Offset)),
	// 	"PageSize":   strconv.Itoa(int(req.Limit)),
	// }

	// err = global.HttpRestApiToMicrosoftGP("GET", "transactionlist/getall", nil, &res, params)

	// if err != nil {
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }

	res = &pb.GetTransactionListGPResponse{
		PageNumber:   1,
		PageSize:     1,
		TotalPages:   1,
		TotalRecords: 1,
		Data: []*pb.TransactionListGP{
			{
				SopTypeId:        "DUMMY SOP TYPE ID",
				SopTypeDesc:      "DUMY SOP TYPE DESC",
				SopNumber:        "DUMMY SOP NUMBER",
				CustomerId:       "DUMMY CUSTOMER ID",
				CustomerName:     "DUMMY CUSTOMER NAME",
				AddressId:        "DUMMY ADDRESS ID",
				RegionDesc:       "DUMMY REGION DESC",
				DocDate:          "01-01-2023",
				RequestShipDate:  "01-01-2023",
				WrtId:            "DUMMY WRT ID",
				ArchetypeId:      "DUMMY ARCHETYPE ID",
				OrderChannelId:   "DUMMY ORDER CHANNEL ID",
				SalesOrderCode:   "DUMMY SALES ORDER CODE",
				BatchId:          "DUMMY BATCH ID",
				DefaultSiteId:    "DUMMY SITE",
				CustomerPoNumber: "DUMMY CUSTOMER PO NUMBER",
				CurrencyId:       "DUMMY CURRENCY ID",
				AmountReceived:   10000,
				TermsDiscTaken:   0,
				OnAccount:        0,
				CommentId:        "DUMMY COMMENT ID",
				TotalWeight:      1,
				Subtotal:         10000,
				TradeDisc:        0,
				Freight:          5000,
				Miscellaneous:    0,
				Tax:              1000,
				Total:            14000,
			},
		},
		Succeeded: true,
		Errors:    []string{},
		Message:   "",
	}

	return
}

func (s *TransactionListService) GetDetailGP(ctx context.Context, req *pb.GetTransactionListGPDetailRequest) (res *pb.GetTransactionListGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "TransactionListService.GetDetailGP")
	defer span.End()

	// params := map[string]string{
	// 	"interid": global.EnvDatabaseGP,
	// 	"type":    req.Type,
	// 	"id":      req.Id,
	// }

	// err = global.HttpRestApiToMicrosoftGP("GET", "transactionlist/getbyid", nil, &res, params)

	// if err != nil {
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	return
	// }
	res = &pb.GetTransactionListGPResponse{
		PageNumber:   1,
		PageSize:     1,
		TotalPages:   1,
		TotalRecords: 1,
		Data: []*pb.TransactionListGP{
			{
				SopTypeId:        "DUMMY SOP TYPE ID",
				SopTypeDesc:      "DUMY SOP TYPE DESC",
				SopNumber:        "DUMMY SOP NUMBER",
				CustomerId:       "DUMMY CUSTOMER ID",
				CustomerName:     "DUMMY CUSTOMER NAME",
				AddressId:        "DUMMY ADDRESS ID",
				RegionDesc:       "DUMMY REGION DESC",
				DocDate:          "01-01-2023",
				RequestShipDate:  "01-01-2023",
				WrtId:            "DUMMY WRT ID",
				ArchetypeId:      "DUMMY ARCHETYPE ID",
				OrderChannelId:   "DUMMY ORDER CHANNEL ID",
				SalesOrderCode:   "DUMMY SALES ORDER CODE",
				BatchId:          "DUMMY BATCH ID",
				DefaultSiteId:    "DUMMY SITE",
				CustomerPoNumber: "DUMMY CUSTOMER PO NUMBER",
				CurrencyId:       "DUMMY CURRENCY ID",
				AmountReceived:   10000,
				TermsDiscTaken:   0,
				OnAccount:        0,
				CommentId:        "DUMMY COMMENT ID",
				TotalWeight:      1,
				Subtotal:         10000,
				TradeDisc:        0,
				Freight:          5000,
				Miscellaneous:    0,
				Tax:              1000,
				Total:            14000,
			},
		},
		Succeeded: true,
		Errors:    []string{},
		Message:   "",
	}

	return
}
