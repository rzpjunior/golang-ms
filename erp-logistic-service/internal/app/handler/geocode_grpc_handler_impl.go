package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *LogisticGrpcHandler) Geocode(ctx context.Context, req *logisticService.GeocodeAddressRequest) (res *logisticService.GeocodeAddressResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.Geocode")
	defer span.End()

	geocode, err := h.ServicesControlTower.Geocode(ctx, &logisticService.GeocodeAddressRequest{
		SalesOrderId: req.SalesOrderId,
		AddressId:    req.AddressId,
		AddressName:  req.AddressName,
		SubDistrict:  req.SubDistrict,
		City:         req.City,
		Region:       req.Region,
		Zip:          req.Zip,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &logisticService.GeocodeAddressResponse{
		Latitude:  geocode.Latitude,
		Longitude: geocode.Longitude,
	}

	return
}
