package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *BridgeGrpcHandler) GetVoucherGPList(ctx context.Context, req *bridgeService.GetVoucherGPListRequest) (res *bridgeService.GetVoucherGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetWrtGPList")
	defer span.End()

	res, err = h.ServicesVoucher.Get(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}

func (h *BridgeGrpcHandler) CreateVoucherGP(ctx context.Context, req *bridgeService.CreateVoucherGPRequest) (res *bridgeService.CreateVoucherGPResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetWrtGPDetail")
	defer span.End()

	body := &dto.CreateVoucherGPRequest{
		GnlVoucherID:      req.GnlVoucherId,
		GnlChannel:        req.GnlChannel,
		GnlVoucherType:    req.GnlVoucherType,
		GnlVoucherName:    req.GnlVoucherName,
		GnlExpenseAccount: req.GnlExpenseAccount,
		GnlVoucherCode:    req.GnlVoucherCode,
		GnlMinimumOrder:   req.GnlMinimumOrder,
		GnlDiscountAmount: req.GnlDiscountAmount,
		GnlVoucherStatus:  req.GnlVoucherStatus,
		Inactive:          req.Inactive,
		Restriction: &dto.Restriction{
			GnlRegion:      req.Restriction.GnlRegion,
			GnlCustTypeID:  req.Restriction.GnlCustTypeId,
			GnlArchetypeID: req.Restriction.GnlArchetypeId,
			DefaultCB:      req.Restriction.DefaultCb,
		},
		AdvancedProperties: &dto.AdvancedProperties{
			Custnmbr:              req.AdvancedProperties.Custnmbr,
			GnlStartPeriod:        req.AdvancedProperties.GnlStartPeriod,
			GnlEndPeriod:          req.AdvancedProperties.GnlEndPeriod,
			GnlTotalQuotaCount:    req.AdvancedProperties.GnlTotalQuotaCount,
			GnlTotalQuotaCountPE:  req.AdvancedProperties.GnlTotalQuotaCountPe,
			GnlRemainingOverallQu: req.AdvancedProperties.GnlRemainingOverallQu,
			GnlMobileVoucher:      req.AdvancedProperties.GnlMobileVoucher,
		},
	}

	res, err = h.ServicesVoucher.Create(ctx, body)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	return
}
