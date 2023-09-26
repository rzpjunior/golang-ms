package service

import (
	"context"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/repository"
)

type IUserCustomerService interface {
	GetDetail(ctx context.Context, req *dto.GetDetailUserCustomerRequest) (res *dto.UserCustomerResponse, err error)
}

type UserCustomerService struct {
	opt                    opt.Options
	RepositoryUserCustomer repository.IUserCustomerRepository
}

func NewUserCustomerService() IUserCustomerService {
	return &UserCustomerService{
		opt:                    global.Setup.Common,
		RepositoryUserCustomer: repository.NewUserCustomerRepository(),
	}
}

func (s *UserCustomerService) GetDetail(ctx context.Context, req *dto.GetDetailUserCustomerRequest) (res *dto.UserCustomerResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "CustomerService.Get")
	defer span.End()

	userCustomer := &model.UserCustomer{
		CustomerID:   req.CustomerID,
		CustomerIDGP: strconv.Itoa(int(req.CustomerID)),
	}

	userCustomer, err = s.RepositoryUserCustomer.GetDetail(ctx, userCustomer)

	res = &dto.UserCustomerResponse{
		ID:            userCustomer.ID,
		CustomerID:    userCustomer.CustomerID,
		FirebaseToken: userCustomer.FirebaseToken,
	}
	return
}
