package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetBankList
func (h *BridgeGrpcHandler) GetBankList(ctx context.Context, req *bridgeService.GetBankListRequest) (res *bridgeService.GetBankListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetBankList")
	defer span.End()

	var Bankes []dto.BankResponse
	Bankes, _, err = h.ServicesBank.Get(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*bridgeService.Bank
	for _, Bank := range Bankes {
		data = append(data, &bridgeService.Bank{
			Id:              Bank.ID,
			Code:            Bank.Code,
			Name:            Bank.Description,
			Note:            "",
			Status:          int32(Bank.Status),
			Value:           Bank.Value,
			ImageUrl:        Bank.ImageUrl,
			PaymentGuideUrl: Bank.PaymentGuideUrl,
			PublishIva:      int32(Bank.PublishIVA),
			PublishFva:      int32(Bank.PublishFVA),
		})
	}

	res = &bridgeService.GetBankListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *BridgeGrpcHandler) GetBankDetail(ctx context.Context, req *bridgeService.GetBankDetailRequest) (res *bridgeService.GetBankDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetBankDetail")
	defer span.End()

	var Bank dto.BankResponse
	Bank, err = h.ServicesBank.GetDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &bridgeService.GetBankDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &bridgeService.Bank{
			Id:              Bank.ID,
			Code:            Bank.Code,
			Name:            Bank.Description,
			Note:            "",
			Status:          int32(Bank.Status),
			Value:           Bank.Value,
			ImageUrl:        Bank.ImageUrl,
			PaymentGuideUrl: Bank.PaymentGuideUrl,
			PublishIva:      int32(Bank.PublishIVA),
			PublishFva:      int32(Bank.PublishFVA),
		},
	}
	return
}
