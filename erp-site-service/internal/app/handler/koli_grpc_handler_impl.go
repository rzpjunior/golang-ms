package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	siteService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/site_service"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *SiteGrpcHandler) GetKoliList(ctx context.Context, req *siteService.GetKoliListRequest) (res *siteService.GetKoliListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetKoliList")
	defer span.End()

	request := dto.KoliGetRequest{
		Offset:  int(req.Offset),
		Limit:   int(req.Limit),
		OrderBy: req.OrderBy,
		Status:  int(req.Status),
	}

	kolis, _, err := h.ServicesKoli.Get(ctx, &request)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*siteService.Koli
	for _, koli := range kolis {
		data = append(data, &siteService.Koli{
			Id:     koli.Id,
			Code:   koli.Code,
			Value:  koli.Value,
			Name:   koli.Name,
			Note:   koli.Note,
			Status: int32(koli.Status),
		})
	}

	res = &siteService.GetKoliListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *SiteGrpcHandler) GetKoliDetail(ctx context.Context, req *siteService.GetKoliDetailRequest) (res *siteService.GetKoliDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetKoliDetail")
	defer span.End()

	koli, err := h.ServicesKoli.GetDetail(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &siteService.GetKoliDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &siteService.Koli{
			Id:     koli.Id,
			Code:   koli.Code,
			Value:  koli.Value,
			Name:   koli.Name,
			Note:   koli.Note,
			Status: int32(koli.Status),
		},
	}

	return
}
