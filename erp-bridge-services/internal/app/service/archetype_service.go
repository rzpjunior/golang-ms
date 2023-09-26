package service

import (
	"context"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
)

type IArchetypeService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, businesssTypeID int64) (res []dto.ArchetypeResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.ArchetypeResponse, err error)
	GetGP(ctx context.Context, req *pb.GetArchetypeGPListRequest) (res *pb.GetArchetypeGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetArchetypeGPDetailRequest) (res *pb.GetArchetypeGPResponse, err error)
}

type ArchetypeService struct {
	opt                    opt.Options
	RepositoryArchetype    repository.IArchetypeRepository
	RepositoryCustomerType repository.ICustomerTypeRepository
}

func NewArchetypeService() IArchetypeService {
	return &ArchetypeService{
		opt:                    global.Setup.Common,
		RepositoryArchetype:    repository.NewArchetypeRepository(),
		RepositoryCustomerType: repository.NewCustomerTypeRepository(),
	}
}

func (s *ArchetypeService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, businesssTypeID int64) (res []dto.ArchetypeResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ArchetypeService.Get")
	defer span.End()

	var archetypes []*model.Archetype
	archetypes, total, err = s.RepositoryArchetype.Get(ctx, offset, limit, status, search, orderBy, businesssTypeID)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, archetype := range archetypes {
		var customerType *model.CustomerType
		customerType, err = s.RepositoryCustomerType.GetDetail(ctx, archetype.CustomerTypeID, "")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}

		customerTypeResponse := &dto.CustomerTypeResponse{
			ID:            customerType.ID,
			Code:          customerType.Code,
			Description:   customerType.Description,
			GroupType:     customerType.GroupType,
			Abbreviation:  customerType.Abbreviation,
			Status:        customerType.Status,
			StatusConvert: statusx.ConvertStatusValue(customerType.Status),
			CreatedAt:     customerType.CreatedAt,
			UpdatedAt:     customerType.UpdatedAt,
		}
		res = append(res, dto.ArchetypeResponse{
			ID:             archetype.ID,
			Code:           archetype.Code,
			CustomerTypeID: archetype.CustomerTypeID,
			Description:    archetype.Description,
			Status:         archetype.Status,
			StatusConvert:  statusx.ConvertStatusValue(archetype.Status),
			CreatedAt:      timex.ToLocTime(ctx, archetype.CreatedAt),
			UpdatedAt:      timex.ToLocTime(ctx, archetype.UpdatedAt),
			CustomerType:   customerTypeResponse,
		})
	}

	return
}

func (s *ArchetypeService) GetDetail(ctx context.Context, id int64, code string) (res dto.ArchetypeResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ArchetypeService.GetDetail")
	defer span.End()

	var archetype *model.Archetype
	archetype, err = s.RepositoryArchetype.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	var customerType *model.CustomerType
	customerType, err = s.RepositoryCustomerType.GetDetail(ctx, archetype.CustomerTypeID, "")
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	customerTypeResponse := &dto.CustomerTypeResponse{
		ID:            customerType.ID,
		Code:          customerType.Code,
		Description:   customerType.Description,
		GroupType:     customerType.GroupType,
		Abbreviation:  customerType.Abbreviation,
		Status:        customerType.Status,
		StatusConvert: statusx.ConvertStatusValue(customerType.Status),
		CreatedAt:     customerType.CreatedAt,
		UpdatedAt:     customerType.UpdatedAt,
	}

	res = dto.ArchetypeResponse{
		ID:             archetype.ID,
		Code:           archetype.Code,
		CustomerTypeID: archetype.CustomerTypeID,
		Description:    archetype.Description,
		Status:         archetype.Status,
		StatusConvert:  statusx.ConvertStatusValue(archetype.Status),
		CreatedAt:      timex.ToLocTime(ctx, archetype.CreatedAt),
		UpdatedAt:      timex.ToLocTime(ctx, archetype.UpdatedAt),
		CustomerType:   customerTypeResponse,
	}

	return
}

func (s *ArchetypeService) GetGP(ctx context.Context, req *pb.GetArchetypeGPListRequest) (res *pb.GetArchetypeGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ArchetypeService.GetGP")
	defer span.End()

	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	// Filter param
	if req.GnlArchetypeId != "" {
		req.GnlArchetypeId = url.PathEscape(req.GnlArchetypeId)
		params["gnl_archetype_id"] = req.GnlArchetypeId
	}

	if req.GnlArchetypedescription != "" {
		req.GnlArchetypedescription = url.PathEscape(req.GnlArchetypedescription)
		params["gnl_archetypedescription"] = req.GnlArchetypedescription
	}

	if req.GnlCustTypeId != "" {
		req.GnlCustTypeId = url.PathEscape(req.GnlCustTypeId)
		params["gnl_cust_type_id"] = req.GnlCustTypeId
	}

	if req.Inactive != "" {
		params["inactive"] = req.Inactive
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "archetype/getall", nil, &res, params)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	return
}

func (s *ArchetypeService) GetDetailGP(ctx context.Context, req *pb.GetArchetypeGPDetailRequest) (res *pb.GetArchetypeGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "ArchetypeService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "archetype/getbyid", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if len(res.Data) == 0 || len(res.Data) > 1 {
		err = edenlabs.ErrorNotFound("archetype")
	}

	return
}
