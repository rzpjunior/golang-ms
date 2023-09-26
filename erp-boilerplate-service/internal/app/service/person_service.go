package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-boilerplate-service/internal/app/repository"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type IPersonService interface {
	Get(ctx context.Context, offset int, limit int, search string) (res []dto.PersonResponseGet, total int64, err error)
	GetByID(ctx context.Context, id int64) (res dto.PersonResponseGet, err error)
	Create(ctx context.Context, req dto.PersonRequestCreate) (res dto.PersonResponseCreate, err error)
	Update(ctx context.Context, req dto.PersonRequestUpdate) (res dto.PersonResponseUpdate, err error)
	Delete(ctx context.Context, req dto.PersonRequestDelete) (res dto.PersonResponseDelete, err error)
}

type PersonService struct {
	opt              opt.Options
	RepositoryPerson repository.IPersonRepository
}

func NewPersonService() IPersonService {
	return &PersonService{
		opt:              global.Setup.Common,
		RepositoryPerson: repository.NewPersonRepository(),
	}
}

func (s *PersonService) Get(ctx context.Context, offset int, limit int, search string) (res []dto.PersonResponseGet, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PersonService.Get")
	defer span.End()

	var persons []*model.Person
	persons, total, err = s.RepositoryPerson.Get(ctx, offset, limit, search)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, person := range persons {
		res = append(res, dto.PersonResponseGet{
			ID:      person.ID,
			Name:    person.Name,
			City:    person.City,
			Country: person.Country,
		})
	}

	return
}

func (s *PersonService) GetByID(ctx context.Context, id int64) (res dto.PersonResponseGet, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PersonService.GetByID")
	defer span.End()

	var person *model.Person
	person, err = s.RepositoryPerson.GetByID(ctx, id)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.PersonResponseGet{
		ID:      person.ID,
		Name:    person.Name,
		City:    person.City,
		Country: person.Country,
	}

	return
}

func (s *PersonService) Create(ctx context.Context, req dto.PersonRequestCreate) (res dto.PersonResponseCreate, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PersonService.Create")
	defer span.End()

	person := &model.Person{
		Name:      req.Name,
		City:      req.City,
		Country:   req.Country,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	span.AddEvent("creating new person")
	err = s.RepositoryPerson.Create(ctx, person)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("person is created", trace.WithAttributes(attribute.Int64("person_id", person.ID)))

	span.AddEvent("publish new person")
	err = s.opt.Producer.KafkaProducer.PublishMessage(ctx, person)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("person is published")

	span.AddEvent("person is created", trace.WithAttributes(attribute.Int64("person_id", person.ID)))
	res = dto.PersonResponseCreate{
		ID:        person.ID,
		Name:      person.Name,
		City:      person.City,
		Country:   person.Country,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}

	return
}

func (s *PersonService) Update(ctx context.Context, req dto.PersonRequestUpdate) (res dto.PersonResponseUpdate, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PersonService.Update")
	defer span.End()

	person := &model.Person{
		ID:        req.ID,
		Name:      req.Name,
		City:      req.City,
		Country:   req.Country,
		UpdatedAt: time.Now(),
	}

	// validate data is exist
	_, err = s.RepositoryPerson.GetByID(ctx, req.ID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = s.RepositoryPerson.Update(ctx, person)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.PersonResponseUpdate{
		ID:        person.ID,
		Name:      person.Name,
		City:      person.City,
		Country:   person.Country,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}

	return
}

func (s *PersonService) Delete(ctx context.Context, req dto.PersonRequestDelete) (res dto.PersonResponseDelete, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PersonService.Delete")
	defer span.End()

	person := &model.Person{
		ID: req.ID,
	}

	// validate data is exist
	_, err = s.RepositoryPerson.GetByID(ctx, req.ID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	err = s.RepositoryPerson.Delete(ctx, person)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.PersonResponseDelete{
		Note: "",
	}

	return
}
