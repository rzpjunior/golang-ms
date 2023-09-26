package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type IDayOffService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, startDate time.Time, endDate time.Time) (res []dto.DayOffResponse, total int64, err error)
	Create(ctx context.Context, req dto.DayOffRequestCreate) (res dto.DayOffResponse, err error)
	Archive(ctx context.Context, id int64) (res dto.DayOffResponse, err error)
	UnArchive(ctx context.Context, id int64) (res dto.DayOffResponse, err error)
}

type DayOffService struct {
	opt              opt.Options
	RepositoryDayOff repository.IDayOffRepository
}

func NewDayOffService() IDayOffService {
	return &DayOffService{
		opt:              global.Setup.Common,
		RepositoryDayOff: repository.NewDayOffRepository(),
	}
}

func (s *DayOffService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string, startDate time.Time, endDate time.Time) (res []dto.DayOffResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DayOffService.Get")
	defer span.End()

	var dayOffs []*model.DayOff
	dayOffs, total, err = s.RepositoryDayOff.Get(ctx, offset, limit, status, search, orderBy, startDate, endDate)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, dayOff := range dayOffs {
		res = append(res, dto.DayOffResponse{
			ID:            dayOff.ID,
			OffDate:       dayOff.OffDate,
			Note:          dayOff.Note,
			Status:        dayOff.Status,
			StatusConvert: statusx.ConvertStatusValue(dayOff.Status),
		})
	}

	return
}

func (s *DayOffService) Create(ctx context.Context, req dto.DayOffRequestCreate) (res dto.DayOffResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DayOffService.Create")
	defer span.End()

	dayOff := &model.DayOff{
		OffDate: req.OffDate,
		Note:    req.Note,
		Status:  1,
	}

	// validate day off, off date and status
	dayOffValidate, _ := s.RepositoryDayOff.GetByOffDate(ctx, req.OffDate.Format("2006-01-02"))
	if dayOffValidate.ID != 0 {
		err = edenlabs.ErrorValidation("day_off", "The day off already exists")
		return
	}

	span.AddEvent("creating new dayOff")
	err = s.RepositoryDayOff.Create(ctx, dayOff)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("dayOff is created", trace.WithAttributes(attribute.Int64("dayOff_id", dayOff.ID)))

	res = dto.DayOffResponse{
		ID:      dayOff.ID,
		OffDate: dayOff.OffDate,
		Note:    dayOff.Note,
		Status:  dayOff.Status,
	}

	return
}

func (s *DayOffService) Archive(ctx context.Context, id int64) (res dto.DayOffResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DayOffService.Archive")
	defer span.End()

	dayOff := &model.DayOff{
		ID:     id,
		Status: 2,
	}

	// validate data is exist
	_, err = s.RepositoryDayOff.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = s.RepositoryDayOff.Archive(ctx, dayOff, "Status")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.DayOffResponse{
		ID:      dayOff.ID,
		OffDate: dayOff.OffDate,
		Note:    dayOff.Note,
		Status:  dayOff.Status,
	}

	return
}

func (s *DayOffService) UnArchive(ctx context.Context, id int64) (res dto.DayOffResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DayOffService.UnArchive")
	defer span.End()

	dayOff := &model.DayOff{
		ID:     id,
		Status: 1,
	}

	// validate data is exist
	dOff, err := s.RepositoryDayOff.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	// validate day off, off date and status
	dayOffValidate, _ := s.RepositoryDayOff.GetByOffDate(ctx, dOff.OffDate.Format("2006-01-02"))
	if dayOffValidate.ID != 0 {
		err = edenlabs.ErrorValidation("day_off", "The day off you want to unarchive is already exists with status active")
		return
	}

	err = s.RepositoryDayOff.UnArchive(ctx, dayOff, "Status")
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.DayOffResponse{
		ID:      dayOff.ID,
		OffDate: dayOff.OffDate,
		Note:    dayOff.Note,
		Status:  dayOff.Status,
	}

	return
}
