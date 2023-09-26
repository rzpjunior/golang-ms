package handler

import (
	context "context"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-account-service/internal/app/service"
	accountService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/account_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *AccountGrpcHandler) GetUserList(ctx context.Context, req *accountService.GetUserListRequest) (res *accountService.GetUserListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressList")
	defer span.End()

	sUser := service.ServiceUser()

	request := &dto.GetUserRequest{
		Offset:     int(req.Offset),
		Limit:      int(req.Limit),
		Status:     int(req.Status),
		Search:     req.Search,
		OrderBy:    req.OrderBy,
		SiteID:     req.SiteIdGp,
		DivisionID: req.DivisionId,
		RoleID:     req.RoleId,
		Apps:       req.Apps,
	}

	users, _, err := sUser.GetList(ctx, request)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*accountService.User
	for _, user := range users {
		data = append(data, &accountService.User{
			Id:       user.ID,
			Email:    user.Email,
			ParentId: user.ParentID,
			// RegionId:     user.Region.ID,
			// SiteId:       user.Site.ID,
			// TerritoryId:  user.Territory.ID,
			EmployeeCode: user.EmployeeCode,
			Name:         user.Name,
			Nickname:     user.Nickname,
			PhoneNumber:  user.PhoneNumber,
			Status:       int32(user.Status),
			CreatedAt:    timestamppb.New(user.CreatedAt),
			UpdatedAt:    timestamppb.New(user.UpdatedAt),
			MainRole:     user.MainRole.Name,
			Division:     user.MainRole.Division.Name,
			SiteIdGp:     request.SiteID,
		})
	}

	res = &accountService.GetUserListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *AccountGrpcHandler) GetUserDetail(ctx context.Context, req *accountService.GetUserDetailRequest) (res *accountService.GetUserDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressList")
	defer span.End()

	sUser := service.ServiceUser()

	var user dto.UserResponse
	user, err = sUser.GetDetail(ctx, req.Id, req.Email, req.EmployeeCode)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &accountService.User{
		Id:         user.ID,
		Email:      user.Email,
		ParentId:   user.ParentID,
		RegionIdGp: user.Region.ID,
		SiteIdGp:   user.Site.ID,
		// TerritoryId:  user.Territory.ID,
		EmployeeCode:           user.EmployeeCode,
		Name:                   user.Name,
		Nickname:               user.Nickname,
		PhoneNumber:            user.PhoneNumber,
		Status:                 int32(user.Status),
		CreatedAt:              timestamppb.New(user.CreatedAt),
		UpdatedAt:              timestamppb.New(user.UpdatedAt),
		MainRole:               user.MainRole.Name,
		Division:               user.MainRole.Division.Name,
		PurchaserappNotifToken: user.PurchaserAppNotifToken,
	}

	res = &accountService.GetUserDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *AccountGrpcHandler) UpdateUserSalesAppToken(ctx context.Context, req *accountService.UpdateUserSalesAppTokenRequest) (res *accountService.GetUserDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressList")
	defer span.End()

	sUser := service.ServiceUser()

	var user dto.ProfileResponse
	user, err = sUser.UpdateSalesAppToken(ctx, dto.UpdateSalesAppTokenRequest{
		Id:                 req.Id,
		ForceLogout:        req.ForceLogout,
		SalesAppLoginToken: req.SalesappLoginToken,
		SalesAppNotifToken: req.SalesappNotifToken,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &accountService.User{
		Id:           user.ID,
		Email:        user.Email,
		ParentId:     user.ParentID,
		EmployeeCode: user.EmployeeCode,
		Name:         user.Name,
		Nickname:     user.Nickname,
		PhoneNumber:  user.PhoneNumber,
		Status:       int32(user.Status),
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
	}

	res = &accountService.GetUserDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *AccountGrpcHandler) GetUserBySalesAppLoginToken(ctx context.Context, req *accountService.GetUserBySalesAppLoginTokenRequest) (res *accountService.GetUserDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressList")
	defer span.End()

	sUser := service.ServiceUser()

	var user dto.UserResponse
	user, err = sUser.GetBySalesAppLoginToken(ctx, req.SalesappLoginToken)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &accountService.User{
		Id:           user.ID,
		Email:        user.Email,
		ParentId:     user.ParentID,
		EmployeeCode: user.EmployeeCode,
		Name:         user.Name,
		Nickname:     user.Nickname,
		PhoneNumber:  user.PhoneNumber,
		Status:       int32(user.Status),
		ForceLogout:  int32(user.ForceLogout),
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
	}

	res = &accountService.GetUserDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *AccountGrpcHandler) UpdateUserEdnAppToken(ctx context.Context, req *accountService.UpdateUserEdnAppTokenRequest) (res *accountService.GetUserDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressList")
	defer span.End()

	sUser := service.ServiceUser()

	var user dto.ProfileResponse
	user, err = sUser.UpdateEdnAppToken(ctx, dto.UpdateEdnAppTokenRequest{
		Id:               req.Id,
		ForceLogout:      req.ForceLogout,
		EdnAppLoginToken: req.EdnappLoginToken,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &accountService.User{
		Id:           user.ID,
		Email:        user.Email,
		ParentId:     user.ParentID,
		EmployeeCode: user.EmployeeCode,
		Name:         user.Name,
		Nickname:     user.Nickname,
		PhoneNumber:  user.PhoneNumber,
		Status:       int32(user.Status),
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
	}

	res = &accountService.GetUserDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *AccountGrpcHandler) UpdateUserPurchaserAppToken(ctx context.Context, req *accountService.UpdateUserPurchaserAppTokenRequest) (res *accountService.GetUserDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressList")
	defer span.End()

	sUser := service.ServiceUser()

	var user dto.ProfileResponse
	user, err = sUser.UpdatePurchaserAppToken(ctx, dto.UpdatePurchaserAppTokenRequest{
		Id:                     req.Id,
		ForceLogout:            req.ForceLogout,
		PurchaserAppNotifToken: req.PurchaserappNotifToken,
	})
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &accountService.User{
		Id:           user.ID,
		Email:        user.Email,
		ParentId:     user.ParentID,
		EmployeeCode: user.EmployeeCode,
		Name:         user.Name,
		Nickname:     user.Nickname,
		PhoneNumber:  user.PhoneNumber,
		Status:       int32(user.Status),
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
	}

	res = &accountService.GetUserDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *AccountGrpcHandler) GetUserByEdnAppLoginToken(ctx context.Context, req *accountService.GetUserByEdnAppLoginTokenRequest) (res *accountService.GetUserDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetAddressList")
	defer span.End()

	sUser := service.ServiceUser()

	var user dto.UserResponse
	user, err = sUser.GetByEdnAppLoginToken(ctx, req.EdnappLoginToken)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	data := &accountService.User{
		Id:           user.ID,
		Email:        user.Email,
		ParentId:     user.ParentID,
		EmployeeCode: user.EmployeeCode,
		Name:         user.Name,
		Nickname:     user.Nickname,
		PhoneNumber:  user.PhoneNumber,
		Status:       int32(user.Status),
		ForceLogout:  int32(user.ForceLogout),
		CreatedAt:    timestamppb.New(user.CreatedAt),
		UpdatedAt:    timestamppb.New(user.UpdatedAt),
	}

	res = &accountService.GetUserDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}
