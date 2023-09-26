package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	siteService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/site_service"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *SiteGrpcHandler) GetSalesOrderDeliveryKoli(ctx context.Context, req *siteService.GetSalesOrderDeliveryKoliRequest) (res *siteService.GetSalesOrderDeliveryKoliResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetSalesOrderDeliveryKoli")
	defer span.End()

	request := &dto.DeliveryKoliGetRequest{
		SopNumber: req.SopNumber,
	}

	deliveryKoli, _, err := h.ServicesDeliveryKoli.Get(ctx, request)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*siteService.DeliveryKoli
	for _, v := range deliveryKoli {
		data = append(data, &siteService.DeliveryKoli{
			Id:             v.Id,
			SalesOrderCode: v.SopNumber,
			KoliId:         v.KoliId,
			Quantity:       v.Quantity,
		})
	}

	res = &siteService.GetSalesOrderDeliveryKoliResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}
