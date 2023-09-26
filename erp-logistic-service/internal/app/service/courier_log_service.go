package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-logistic-service/internal/app/repository"
	logisticService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/logistic_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ICourierLogService interface {
	Create(ctx context.Context, req *logisticService.CreateCourierLogRequest) (res *dto.CourierLogResponse, err error)
	GetLastCourierLog(ctx context.Context, req *logisticService.GetLastCourierLogRequest) (res *logisticService.GetLastCourierLogResponse, err error)
}

type CourierLogService struct {
	opt                  opt.Options
	RepositoryCourierLog repository.ICourierLogRepository
}

func NewCourierLogService() ICourierLogService {
	return &CourierLogService{
		opt:                  global.Setup.Common,
		RepositoryCourierLog: repository.NewCourierLogRepository(),
	}
}

func (s *CourierLogService) Create(ctx context.Context, req *logisticService.CreateCourierLogRequest) (res *dto.CourierLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierLogService.Create")
	defer span.End()

	model := &dto.CourierLog{
		Latitude:     req.Model.Latitude,
		Longitude:    req.Model.Longitude,
		CreatedAt:    req.Model.CreatedAt.AsTime(),
		CourierID:    req.Model.CourierId,
		SalesOrderID: req.Model.SalesOrderId,
	}

	if err = s.RepositoryCourierLog.Create(ctx, model); err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.CourierLogResponse{
		ID:           model.ID,
		Latitude:     model.Latitude,
		Longitude:    model.Longitude,
		CreatedAt:    model.CreatedAt,
		CourierID:    model.CourierID,
		SalesOrderID: model.SalesOrderID,
	}

	return
}

func (s *CourierLogService) GetLastCourierLog(ctx context.Context, req *logisticService.GetLastCourierLogRequest) (res *logisticService.GetLastCourierLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CourierLogService.GetLastCourierLog")
	defer span.End()
	var courierLog *dto.CourierLog
	courierLog, err = s.RepositoryCourierLog.GetLastCourierLog(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &logisticService.GetLastCourierLogResponse{
		Latitude:  *courierLog.Latitude,
		Longitude: *courierLog.Longitude,
		CourierId: courierLog.CourierID,
		CreatedAt: timestamppb.New(courierLog.CreatedAt),
	}

	return
}
