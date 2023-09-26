package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
)

type IFieldPurchaserService interface {
	Get(ctx context.Context, req *dto.FieldPurchaserListRequest) (res []*dto.FieldPurchaserResponse, total int64, err error)
	GetByID(ctx context.Context, id int64) (res *dto.FieldPurchaserResponse, err error)
}

type FieldPurchaserService struct {
	opt opt.Options
}

func NewFieldPurchaserService() IFieldPurchaserService {
	return &FieldPurchaserService{
		opt: global.Setup.Common,
	}
}

func (s *FieldPurchaserService) Get(ctx context.Context, req *dto.FieldPurchaserListRequest) (res []*dto.FieldPurchaserResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "FieldPurchaserService.Get")
	defer span.End()

	var fieldPurchasers *accountService.GetUserListResponse
	fieldPurchasers, err = s.opt.Client.AccountServiceGrpc.GetUserList(ctx, &accountService.GetUserListRequest{
		Limit:    req.Limit,
		Offset:   req.Offset,
		Status:   req.Status,
		Search:   req.Search,
		OrderBy:  req.OrderBy,
		SiteIdGp: req.SiteIDGp,
		// SiteId:  int64(req.SiteID),
		Apps: "mob-purchaser",
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "field_purchaser")
		return
	}

	for _, fp := range fieldPurchasers.Data {
		if fp.MainRole == "Field Purchaser" {
			res = append(res, &dto.FieldPurchaserResponse{
				ID:           fp.Id,
				Name:         fp.Name,
				Nickname:     fp.Nickname,
				EmployeeCode: fp.EmployeeCode,
				PhoneNumber:  fp.PhoneNumber,
				Division:     fp.Division,
				MainRole:     fp.MainRole,
			})
		}
	}

	total = int64(len(res))

	return
}

func (s *FieldPurchaserService) GetByID(ctx context.Context, id int64) (res *dto.FieldPurchaserResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "FieldPurchaserService.GetByID")
	defer span.End()

	var fieldPurchaser *accountService.GetUserDetailResponse
	fieldPurchaser, err = s.opt.Client.AccountServiceGrpc.GetUserDetail(ctx, &accountService.GetUserDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("account", "field_purchaser")
		return
	}

	res = &dto.FieldPurchaserResponse{
		ID:           fieldPurchaser.Data.Id,
		Name:         fieldPurchaser.Data.Name,
		Nickname:     fieldPurchaser.Data.Nickname,
		EmployeeCode: fieldPurchaser.Data.EmployeeCode,
		PhoneNumber:  fieldPurchaser.Data.PhoneNumber,
		Division:     fieldPurchaser.Data.Division,
		MainRole:     fieldPurchaser.Data.MainRole,
	}

	return
}
