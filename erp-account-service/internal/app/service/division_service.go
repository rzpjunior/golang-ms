package service

import (
	"context"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/repository"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	dto "git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/model"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func ServiceDivision() IDivisionService {
	m := new(DivisionService)
	m.opt = global.Setup.Common
	return m
}

type IDivisionService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.DivisionResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.DivisionResponse, err error)
	Create(ctx context.Context, req dto.DivisionRequestCreate) (res dto.DivisionResponse, err error)
	Update(ctx context.Context, req dto.DivisionRequestUpdate, id int64) (res dto.DivisionResponse, err error)
	Archive(ctx context.Context, id int64) (res dto.DivisionResponse, err error)
	GetDivisonByCustomerType(ctx context.Context, customerTypeID string) (res *dto.DivisionResponse, err error)
}

type DivisionService struct {
	opt opt.Options
}

func (s *DivisionService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.DivisionResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DivisionService.Get")
	defer span.End()
	rDivision := repository.RepositoryDivision()

	var divisions []*model.Division

	divisions, total, err = rDivision.Get(ctx, offset, limit, status, search, orderBy)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, division := range divisions {
		res = append(res, &dto.DivisionResponse{
			ID:            division.ID,
			Code:          division.Code,
			Name:          division.Name,
			Note:          division.Note,
			CreatedAt:     division.CreatedAt,
			UpdatedAt:     division.UpdatedAt,
			Status:        division.Status,
			StatusConvert: statusx.ConvertStatusValue(division.Status),
		})
	}

	return
}

func (s *DivisionService) GetDetail(ctx context.Context, id int64, code string) (res dto.DivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DivisionService.GetByID")
	defer span.End()
	rDivision := repository.RepositoryDivision()

	var division *model.Division
	division, err = rDivision.GetDetail(ctx, id, code)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.DivisionResponse{
		ID:            division.ID,
		Code:          division.Code,
		Name:          division.Name,
		Note:          division.Note,
		CreatedAt:     division.CreatedAt,
		UpdatedAt:     division.UpdatedAt,
		Status:        division.Status,
		StatusConvert: statusx.ConvertStatusValue(division.Status),
	}

	return
}

func (s *DivisionService) Create(ctx context.Context, req dto.DivisionRequestCreate) (res dto.DivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DivisionService.Create")
	defer span.End()
	rDivision := repository.RepositoryDivision()

	division := &model.Division{
		Code:      req.Code,
		Name:      req.Name,
		Note:      req.Note,
		CreatedAt: time.Now(),
		Status:    1,
	}

	// validate division name
	var existsDivision *model.Division
	existsDivision, _ = rDivision.GetByName(ctx, req.Name)
	if existsDivision.ID != 0 {
		err = edenlabs.ErrorValidation("name", "The name is already exists")
		return
	}

	span.AddEvent("creating new division")
	err = rDivision.Create(ctx, division)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("division is created", trace.WithAttributes(attribute.Int64("division_id", division.ID)))

	res = dto.DivisionResponse{
		ID:        division.ID,
		Code:      division.Code,
		Name:      division.Name,
		CreatedAt: division.CreatedAt,
		UpdatedAt: division.UpdatedAt,
		Status:    division.Status,
	}

	return
}

func (s *DivisionService) Update(ctx context.Context, req dto.DivisionRequestUpdate, id int64) (res dto.DivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DivisionService.Update")
	defer span.End()
	rDivision := repository.RepositoryDivision()
	division := &model.Division{
		ID:        id,
		Code:      req.Code,
		Name:      req.Name,
		Note:      req.Note,
		UpdatedAt: time.Now(),
	}

	// validate data is exist
	var divisionOld *model.Division
	divisionOld, err = rDivision.GetDetail(ctx, id, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	if divisionOld.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	if divisionOld.Name != req.Name {
		// validate division name
		var existsDivision *model.Division
		existsDivision, _ = rDivision.GetByName(ctx, req.Name)
		if existsDivision.ID != 0 {
			err = edenlabs.ErrorValidation("name", "The name is already exists")
			return
		}
	}

	err = rDivision.Update(ctx, division, "Code", "Name", "Note", "UpdatedAt")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.DivisionResponse{
		ID:        division.ID,
		Code:      division.Code,
		Name:      division.Name,
		CreatedAt: division.CreatedAt,
		UpdatedAt: division.UpdatedAt,
		Status:    division.Status,
	}

	return
}

func (s *DivisionService) Archive(ctx context.Context, id int64) (res dto.DivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DivisionService.Delete")
	defer span.End()

	rDivision := repository.RepositoryDivision()
	rUser := repository.RepositoryUser()
	// validate division is exist
	var divisionOld *model.Division
	divisionOld, err = rDivision.GetDetail(ctx, id, "")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate user is inactive
	var usersActive []*model.User
	usersActive, _ = rUser.GetByDivisionID(ctx, id)
	if usersActive != nil {
		err = edenlabs.ErrorValidation("status", "The division still have active users, please check users active")
		return
	}

	if divisionOld.Status == 2 {
		err = edenlabs.ErrorValidation("status", "The status has been archived")
		return
	}

	if divisionOld.Status != 1 {
		err = edenlabs.ErrorValidation("status", "The status must be active")
		return
	}

	err = rDivision.Archive(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *DivisionService) GetDivisonByCustomerType(ctx context.Context, customerTypeID string) (res *dto.DivisionResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DivisionService.Get")
	defer span.End()
	rDivision := repository.RepositoryDivision()

	var (
		division *model.Division
		config   *configuration_service.GetConfigAppListResponse
	)

	_, err = s.opt.Client.BridgeServiceGrpc.GetCustomerTypeGPDetail(ctx, &bridge_service.GetCustomerTypeGPDetailRequest{
		Id: customerTypeID,
	})

	if err != nil {
		err = edenlabs.ErrorInvalid("customer_type_id")
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Check customer type to config COA
	config, err = s.opt.Client.ConfigurationServiceGrpc.GetConfigAppList(ctx, &configuration_service.GetConfigAppListRequest{
		Offset:    0,
		Limit:     1,
		Value:     customerTypeID,
		Attribute: "customer_type_coa",
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// Set default to others
	if len(config.Data) == 0 {
		res = &dto.DivisionResponse{
			ID:   0,
			Code: "OTH",
			Name: "Others",
		}
		return
	}

	// Get Name Division from config
	divisionName := strings.Split(config.Data[0].Attribute, "_")[3]

	division, err = rDivision.GetByName(ctx, divisionName)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.DivisionResponse{
		ID:            division.ID,
		Code:          division.Code,
		Name:          division.Name,
		Note:          division.Note,
		CreatedAt:     division.CreatedAt,
		UpdatedAt:     division.UpdatedAt,
		Status:        division.Status,
		StatusConvert: statusx.ConvertStatusValue(division.Status),
	}

	return
}
