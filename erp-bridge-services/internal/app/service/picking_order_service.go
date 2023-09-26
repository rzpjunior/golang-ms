package service

import (
	"context"
	"fmt"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/sirupsen/logrus"
)

type IPickingOrderService interface {
	GetGrpc(ctx context.Context, req *pb.GetPickingOrderGPHeaderRequest) (res *pb.GetPickingOrderGPHeaderResponse, err error)
	GetDetailGrpc(ctx context.Context, req *pb.GetPickingOrderGPDetailRequest) (res *pb.GetPickingOrderGPDetailResponse, err error)
	SubmitPickingChecking(ctx context.Context, req *pb.SubmitPickingCheckingRequest) (res *pb.SubmitPickingCheckingResponse, err error)
}

type PickingOrderService struct {
	opt opt.Options
}

func NewPickingOrderService() IPickingOrderService {
	return &PickingOrderService{
		opt: global.Setup.Common,
	}
}

func (s *PickingOrderService) GetGrpc(ctx context.Context, req *pb.GetPickingOrderGPHeaderRequest) (res *pb.GetPickingOrderGPHeaderResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.GetGrpc")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Locncode != "" {
		params["locncode"] = req.Locncode
	}

	if req.Sopnumbe != "" {
		params["sopnumbe"] = req.Sopnumbe
	}

	if req.Docnumbr != "" {
		params["docnumbr"] = req.Docnumbr
	}

	if req.Itemnmbr != "" {
		params["itemnmbr"] = req.Itemnmbr
	}

	if req.DocdateFrom != "" && req.DocdateTo != "" {
		params["docdate_from"] = req.DocdateFrom
		params["docdate_to"] = req.DocdateTo
	}

	if req.GnlHelperId != "" {
		params["gnl_helper_id"] = req.GnlHelperId
	}

	if req.WmsPickingStatus != 0 {
		params["wms_picking_status"] = strconv.Itoa(int(req.WmsPickingStatus))
	}

	if req.Custname != "" {
		params["custname"] = req.Custname
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "PickingOrder/list", nil, &res, params)

	return
}

func (s *PickingOrderService) GetDetailGrpc(ctx context.Context, req *pb.GetPickingOrderGPDetailRequest) (res *pb.GetPickingOrderGPDetailResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.GetDetailGrpc")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "PickingOrder/detail", nil, &res, params)

	return
}

func (s *PickingOrderService) SubmitPickingChecking(ctx context.Context, req *pb.SubmitPickingCheckingRequest) (res *pb.SubmitPickingCheckingResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PickingOrderService.SubmitPickingChecking")
	defer span.End()

	var (
		requestPickingDetail = []*dto.PickingDetails{}
		requestChecking      = []*dto.Checking{}
	)

	for _, v := range req.Picking.Details {
		requestPickingDetail = append(requestPickingDetail, &dto.PickingDetails{
			Sopnumbe:     v.Sopnumbe,
			Lnitmseq:     v.Lnitmseq,
			IvmQtyPickso: v.IvmQtyPickso,
		})
	}

	for _, v := range req.Checking {
		var checkingDetail []*dto.CheckingDetails
		for _, v2 := range v.Details {
			checkingDetail = append(checkingDetail, &dto.CheckingDetails{
				Sopnumbe:     v2.Sopnumbe,
				Lnitmseq:     v2.Lnitmseq,
				IvmQtyPickso: v2.IvmQtyPickso,
			})
		}

		requestChecking = append(requestChecking, &dto.Checking{
			Docnumbr:     v.Docnumbr,
			Sopnumbe:     v.Sopnumbe,
			Strttime:     v.Strttime,
			Endtime:      v.Endtime,
			WmsPickerId:  v.WmsPickerId,
			IvmKoli:      v.IvmKoli,
			IvmJenisKoli: v.ImvJenisKoli,
			Details:      checkingDetail,
		})
	}
	requestdto := &dto.SubmitPickingCheckingRequest{
		Interid:  global.EnvDatabaseGP,
		Uniqueid: req.Uniqueid,
		Bachnumb: req.Bachnumb,
		Picking: &dto.Picking{
			Docnumbr: req.Picking.Docnumbr,
			Strttime: req.Picking.Strttime,
			Endtime:  req.Picking.Endtime,
			Details:  requestPickingDetail,
		},
		Checking: requestChecking,
	}

	err = global.HttpRestApiToMicrosoftGP("POST", "sales/submitchecking", requestdto, &res, nil)
	if err != nil {
		fmt.Println("err", err)
		logrus.Error(err.Error())
		return
	}

	return
}
