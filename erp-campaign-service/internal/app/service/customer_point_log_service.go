package service

import (
	"context"
	"strconv"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-campaign-service/internal/app/repository"
	auditService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/audit_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/campaign_service"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ICustomerPointLogService interface {
	Get(ctx context.Context, req *dto.CustomerPointLogRequestGet) (res []*dto.CustomerPointLogResponse, total int64, err error)
	GetDetail(ctx context.Context, req *dto.CustomerPointLogRequestGetDetail) (res *dto.CustomerPointLogResponse, err error)
	Create(ctx context.Context, req *dto.CustomerPointLogRequestCreate) (customerPointLogID int64, err error)
	GetDetailHistoryMobile(ctx context.Context, req *dto.CustomerPointLogRequestGetDetail) (res *dto.PointHistoryList, err error)
	GetReferralDataCustomer(ctx context.Context, req *campaign_service.GetReferralHistoryRequest) (referralHistoryReturn *dto.ReferralHistoryReturn, err error)
	Cancel(ctx context.Context, req *dto.CustomerPointLogRequestCancel) (customerPointLogID int64, err error)
}

type CustomerPointLogService struct {
	opt                        opt.Options
	RepositoryCustomerPointLog repository.ICustomerPointLogRepository
}

func NewCustomerPointLogService() ICustomerPointLogService {
	return &CustomerPointLogService{
		opt:                        global.Setup.Common,
		RepositoryCustomerPointLog: repository.NewCustomerPointLogRepository(),
	}
}

func (s *CustomerPointLogService) Get(ctx context.Context, req *dto.CustomerPointLogRequestGet) (res []*dto.CustomerPointLogResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointLogService.Get")
	defer span.End()

	var CustomerPointLogs []*model.CustomerPointLog
	CustomerPointLogs, total, err = s.RepositoryCustomerPointLog.Get(ctx, req)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, customerPointLog := range CustomerPointLogs {
		res = append(res, &dto.CustomerPointLogResponse{
			ID:               customerPointLog.ID,
			CustomerID:       customerPointLog.CustomerID,
			SalesOrderID:     customerPointLog.SalesOrderID,
			EPCampaignID:     customerPointLog.EPCampaignID,
			CurrentPointUsed: customerPointLog.CurrentPointUsed,
			NextPointUsed:    customerPointLog.NextPointUsed,
			PointValue:       customerPointLog.PointValue,
			RecentPoint:      customerPointLog.RecentPoint,
			Status:           customerPointLog.Status,
			TransactionType:  customerPointLog.TransactionType,
			CreatedDate:      customerPointLog.CreatedDate,
			ExpiredDate:      customerPointLog.ExpiredDate,
		})
	}

	return
}

func (s *CustomerPointLogService) GetDetail(ctx context.Context, req *dto.CustomerPointLogRequestGetDetail) (res *dto.CustomerPointLogResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointLogService.GetByID")
	defer span.End()

	var customerPointLog *model.CustomerPointLog
	customerPointLog, err = s.RepositoryCustomerPointLog.GetDetail(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.CustomerPointLogResponse{
		ID:               customerPointLog.ID,
		CustomerID:       customerPointLog.CustomerID,
		SalesOrderID:     customerPointLog.SalesOrderID,
		EPCampaignID:     customerPointLog.EPCampaignID,
		CurrentPointUsed: customerPointLog.CurrentPointUsed,
		NextPointUsed:    customerPointLog.NextPointUsed,
		PointValue:       customerPointLog.PointValue,
		RecentPoint:      customerPointLog.RecentPoint,
		Status:           customerPointLog.Status,
		TransactionType:  customerPointLog.TransactionType,
		CreatedDate:      customerPointLog.CreatedDate,
		ExpiredDate:      customerPointLog.ExpiredDate,
	}

	return
}

func (s *CustomerPointLogService) GetDetailHistoryMobile(ctx context.Context, req *dto.CustomerPointLogRequestGetDetail) (res *dto.PointHistoryList, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointLogService.GetByID")
	defer span.End()

	var customerPointLog *model.PointHistoryList
	customerPointLog, err = s.RepositoryCustomerPointLog.GetDetailHistoryMobile(ctx, req)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = &dto.PointHistoryList{
		CreatedDate: customerPointLog.CreatedDate,
		PointValue:  customerPointLog.PointValue,
		StatusType:  customerPointLog.StatusType,
		Status:      customerPointLog.Status,
	}

	return
}

func (s *CustomerPointLogService) Create(ctx context.Context, req *dto.CustomerPointLogRequestCreate) (customerPointLogID int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointLogService.Create")
	defer span.End()

	customerPointLog := &model.CustomerPointLog{
		CustomerID:       req.CustomerID,
		SalesOrderID:     req.SalesOrderID,
		EPCampaignID:     req.EPCampaignID,
		CurrentPointUsed: req.CurrentPointUsed,
		NextPointUsed:    req.NextPointUsed,
		PointValue:       req.PointValue,
		RecentPoint:      req.RecentPoint,
		Status:           req.Status,
		TransactionType:  req.TransactionType,
		CreatedDate:      req.CreatedDate,
		ExpiredDate:      req.ExpiredDate,
		Note:             req.Note,
	}

	span.AddEvent("creating new customer point log")
	err = s.RepositoryCustomerPointLog.Create(ctx, customerPointLog)
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	span.AddEvent("customer point log is created", trace.WithAttributes(attribute.Int64("customer_point_log_id", customerPointLog.ID)))

	_, err = s.opt.Client.AuditServiceGrpc.CreateLog(ctx, &auditService.CreateLogRequest{
		Log: &auditService.Log{
			UserId:      22,
			ReferenceId: strconv.Itoa(int(customerPointLog.ID)),
			Type:        "customer_point_log",
			Function:    "create",
			CreatedAt:   timestamppb.New(time.Now()),
			Note:        req.Note,
		},
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	customerPointLogID = customerPointLog.ID

	return
}

func (s *CustomerPointLogService) GetReferralDataCustomer(ctx context.Context, req *campaign_service.GetReferralHistoryRequest) (referralHistoryReturn *dto.ReferralHistoryReturn, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointLogRepository.Update")
	defer span.End()

	referralHistoryReturn = &dto.ReferralHistoryReturn{}

	customer, err := s.opt.Client.BridgeServiceGrpc.GetCustomerDetail(ctx, &bridge_service.GetCustomerDetailRequest{
		Id: req.ReferrerId,
	})
	//get from db first
	referralHistoryReturn, err = s.RepositoryCustomerPointLog.GetReferralDataCustomer(ctx, req.ReferrerId)
	for i, _ := range referralHistoryReturn.Detail.ReferralPointList {
		referralHistoryReturn.Detail.ReferralPointList[i].Name = customer.Data.Name
	}
	//get customer referral
	customers, err := s.opt.Client.BridgeServiceGrpc.GetCustomerList(ctx, &bridge_service.GetCustomerListRequest{
		ReferrerId: req.ReferrerId,
	})
	//layout := "2006-01-02 15:04:05"

	for _, v := range customers.Data {
		referralHistoryReturn.Detail.ReferralList = append(referralHistoryReturn.Detail.ReferralList, &model.ReferralList{
			Name:      v.Name,
			CreatedAt: v.CreatedAt.AsTime(),
		})
	}
	referralHistoryReturn.Summary.TotalReferral = int64(len(referralHistoryReturn.Detail.ReferralList))

	return
}

func (s *CustomerPointLogService) Cancel(ctx context.Context, req *dto.CustomerPointLogRequestCancel) (customerPointLogID int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerPointLogService.Create")
	defer span.End()

	var customerPointLogs []*model.CustomerPointLog

	param := &dto.CustomerPointLogRequestGet{
		CustomerID:   req.CustomerID,
		SalesOrderID: req.SalesOrderID,
		Limit:        1,
		Status:       statusx.ConvertStatusName(statusx.Used),
	}

	customerPointLogs, _, err = s.RepositoryCustomerPointLog.Get(ctx, param)

	for _, v := range customerPointLogs {
		customerPointLog := &model.CustomerPointLog{
			ID:     v.ID,
			Status: statusx.ConvertStatusName(statusx.Cancelled),
			Note:   "Cancellation due to cancel sales order",
		}
		err = s.RepositoryCustomerPointLog.Update(ctx, customerPointLog, "Status")
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			return
		}
	}

	return
}
