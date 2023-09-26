package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/sirupsen/logrus"
)

type IDeliveryOrderService interface {
	GetDetail(ctx context.Context, req *bridge_service.GetDeliveryOrderDetailRequest) (res dto.DeliveryOrderResponse, err error)
	GetListGP(ctx context.Context, req *bridge_service.GetDeliveryOrderGPListRequest) (res *bridge_service.GetDeliveryOrderGPListResponse, err error)
	Create(ctx context.Context, req *bridge_service.CreateDeliveryOrderRequest) (res *bridge_service.CreateDeliveryOrderResponse, err error)
}

type DeliveryOrderService struct {
	opt                     opt.Options
	RepositoryDeliveryOrder repository.IDeliveryOrderRepository
}

func NewDeliveryOrderService() IDeliveryOrderService {
	return &DeliveryOrderService{
		opt:                     global.Setup.Common,
		RepositoryDeliveryOrder: repository.NewDeliveryOrderRepository(),
	}
}

func (s *DeliveryOrderService) GetListGP(ctx context.Context, req *bridge_service.GetDeliveryOrderGPListRequest) (res *bridge_service.GetDeliveryOrderGPListResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryOrderService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.DocNumber != "" {
		params["docnumbr"] = req.DocNumber
	}
	if req.SopNumbe != "" {
		params["sopnumbe"] = req.SopNumbe
	}
	if req.DeltaUser != "" {
		params["delta_user"] = req.DeltaUser
	}
	err = global.HttpRestApiToMicrosoftGP("GET", "DeliveryOrder/list", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *DeliveryOrderService) GetDetail(ctx context.Context, req *bridge_service.GetDeliveryOrderDetailRequest) (res dto.DeliveryOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryOrderService.GetDetail")
	defer span.End()

	var DeliveryOrder *model.DeliveryOrder
	DeliveryOrder, err = s.RepositoryDeliveryOrder.GetDetail(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.DeliveryOrderResponse{
		ID:              DeliveryOrder.ID,
		Code:            DeliveryOrder.Code,
		CustomerID:      DeliveryOrder.CustomerID,
		WrtID:           DeliveryOrder.WrtID,
		SiteID:          DeliveryOrder.SiteID,
		Status:          DeliveryOrder.Status,
		RecognitionDate: timex.ToLocTime(ctx, DeliveryOrder.RecognitionDate),
		StatusConvert:   statusx.ConvertStatusValue(DeliveryOrder.Status),
		CreatedDate:     timex.ToLocTime(ctx, DeliveryOrder.CreatedDate),
	}

	return
}

func (s *DeliveryOrderService) Create(ctx context.Context, req *bridge_service.CreateDeliveryOrderRequest) (res *bridge_service.CreateDeliveryOrderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryOrderService.Create")
	defer span.End()

	var (
		requestDetailOrder []*dto.DeliveryOrderDetailOrder
		requestDetailItem  []*dto.DeliveryOrderDetailItem
	)

	for _, v := range req.DetailOrder {
		requestDetailOrder = append(requestDetailOrder, &dto.DeliveryOrderDetailOrder{
			IvmCb:      v.IvmCb,
			SopNumber:  v.Sopnumbe,
			QtyOrder:   v.Qtyorder,
			IvmQtyPack: v.IvmQtyPack,
		})
	}

	for _, v := range req.DetailItem {
		requestDetailItem = append(requestDetailItem, &dto.DeliveryOrderDetailItem{
			Lnseqnbr:   v.Lnseqnbr,
			ItemNumber: v.Itemnmbr,
			Uofm:       v.Uofm,
			QtyOrder:   float64(v.Qtyorder),
			IvmQtyPack: float64(v.IvmQtyPack),
			Locncode:   v.Locncode,
		})
	}

	requestdto := &dto.CreateDeliveryOrderRequest{
		Interid:      global.EnvDatabaseGP,
		Docnumber:    req.Docnumbr,
		Docdate:      req.Docdate,
		Custnumber:   req.Custnmbr,
		Custname:     req.Custname,
		GnlCourierId: req.GnLCourierId,
		DetailOrder:  requestDetailOrder,
		DetailItem:   requestDetailItem,
	}

	err = global.HttpRestApiToMicrosoftGP("POST", "DeliveryOrder/create", requestdto, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		return
	}

	return
}
