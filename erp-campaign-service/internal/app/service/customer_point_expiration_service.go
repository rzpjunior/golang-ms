package service

import (
	"context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/repository"
)

type ICustomerPointExpirationService interface {
	GetDetail(ctx context.Context, id, customerID int64) (res *dto.CustomerPointExpirationResponse, err error)
	// Create(ctx context.Context, req *dto.CustomerPointExpirationRequestCreate) (err error)
	// Update(ctx context.Context, req *dto.CustomerPointExpirationRequestUpdate) (err error)
}

type CustomerPointExpirationService struct {
	opt                               opt.Options
	RepositoryCustomerPointExpiration repository.ICustomerPointExpirationRepository
}

func NewCustomerPointExpirationService() ICustomerPointExpirationService {
	return &CustomerPointExpirationService{
		opt:                               global.Setup.Common,
		RepositoryCustomerPointExpiration: repository.NewCustomerPointExpirationRepository(),
	}
}

func (s *CustomerPointExpirationService) GetDetail(ctx context.Context, id, customerID int64) (res *dto.CustomerPointExpirationResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointExpirationService.GetDetail")
	defer span.End()

	var customerPointExpiration *model.CustomerPointExpiration
	customerPointExpiration, err = s.RepositoryCustomerPointExpiration.GetDetail(ctx, id, customerID)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.CustomerPointExpirationResponse{
		ID:                 customerPointExpiration.ID,
		CustomerID:         customerPointExpiration.CustomerID,
		CurrentPeriodPoint: customerPointExpiration.CurrentPeriodPoint,
		NextPeriodPoint:    customerPointExpiration.NextPeriodPoint,
		CurrentPeriodDate:  customerPointExpiration.CurrentPeriodDate,
		NextPeriodDate:     customerPointExpiration.NextPeriodDate,
		LastUpdatedAt:      customerPointExpiration.LastUpdatedAt,
	}

	return
}

// func (s *CustomerPointExpirationService) Create(ctx context.Context, req *dto.CustomerPointExpirationRequestCreate) (err error) {
// 	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointExpirationService.Create")
// 	defer span.End()

// 	customerPointLog := &model.CustomerPointExpiration{
// 		CustomerID:     req.CustomerID,
// 		EarnedPoint:    req.EarnedPoint,
// 		RedeemedPoint:  req.RedeemedPoint,
// 		ExpirationDate: req.ExpirationDate,
// 	}

// 	span.AddEvent("creating new customer point Expiration")
// 	err = s.RepositoryCustomerPointExpiration.Create(ctx, customerPointLog)
// 	if err != nil {
// 		span.RecordError(err)
// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
// 		return
// 	}
// 	span.AddEvent("customer point Expiration is created", trace.WithAttributes(attribute.Int64("customer_point_Expiration_id", customerPointLog.ID)))

// 	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
// 		Log: &auditService.Log{
// 			UserId:      22,
// 			ReferenceId: customerPointLog.ID,
// 			Type:        "customer_point_Expiration",
// 			Function:    "create",
// 			CreatedAt:   timestamppb.New(time.Now()),
// 		},
// 	})
// 	if err != nil {
// 		span.RecordError(err)
// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
// 		return
// 	}

// 	return
// }

// func (s *CustomerPointExpirationService) Update(ctx context.Context, req *dto.CustomerPointExpirationRequestUpdate) (err error) {
// 	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointExpirationService.Update")
// 	defer span.End()

// 	customerPointLog := &model.CustomerPointExpiration{
// 		ID:            req.ID,
// 		CustomerID:    req.CustomerID,
// 		EarnedPoint:   req.EarnedPoint,
// 		RedeemedPoint: req.RedeemedPoint,
// 	}

// 	span.AddEvent("updating new customer point Expiration")
// 	err = s.RepositoryCustomerPointExpiration.Update(ctx, customerPointLog, req.FieldUpdate...)
// 	if err != nil {
// 		span.RecordError(err)
// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
// 		return
// 	}
// 	span.AddEvent("customer point Expiration is updated", trace.WithAttributes(attribute.Int64("customer_point_Expiration_id", customerPointLog.ID)))

// 	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
// 		Log: &auditService.Log{
// 			UserId:      22,
// 			ReferenceId: customerPointLog.ID,
// 			Type:        "customer_point_Expiration",
// 			Function:    "update",
// 			CreatedAt:   timestamppb.New(time.Now()),
// 		},
// 	})
// 	if err != nil {
// 		span.RecordError(err)
// 		s.opt.Logger.AddMessage(log.ErrorLevel, err)
// 		return
// 	}

// 	return
// }
