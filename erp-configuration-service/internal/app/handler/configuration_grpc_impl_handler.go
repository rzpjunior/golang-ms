package handler

import (
	context "context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	dto "git.edenfarm.id/project-version3/erp-services/erp-configuration-service/internal/app/dto"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (h *ConfigurationGrpcHandler) GetGenerateCode(ctx context.Context, req *pb.GetGenerateCodeRequest) (res *pb.GetGenerateCodeResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetCode")
	defer span.End()

	codeGenerated, err := h.ServiceCodeGenerator.GetGenerateCode(ctx, req.Format, req.Domain, int(req.Length))
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.GetGenerateCodeResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.GenerateCode{
			Code: codeGenerated,
		},
	}
	return
}

func (h *ConfigurationGrpcHandler) GetGenerateCustomerCode(ctx context.Context, req *pb.GetGenerateCodeRequest) (res *pb.GetGenerateCodeResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetGenerateCustomerCode")
	defer span.End()

	codeGenerated, err := h.ServiceCodeGenerator.GenerateCustomerCode(ctx, req.Format, req.Domain)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.GetGenerateCodeResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.GenerateCode{
			Code: codeGenerated,
		},
	}
	return
}

func (h *ConfigurationGrpcHandler) GetGenerateReferralCode(ctx context.Context, req *pb.GetGenerateCodeRequest) (res *pb.GetGenerateCodeResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetGenerateReferralCode")
	defer span.End()

	codeGenerated, err := h.ServiceCodeGenerator.GenerateReferralCode(ctx, int(req.Length))
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.GetGenerateCodeResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.GenerateCode{
			Code: codeGenerated,
		},
	}
	return
}

func (h *ConfigurationGrpcHandler) GetGlossaryList(ctx context.Context, req *pb.GetGlossaryListRequest) (res *pb.GetGlossaryListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetGlossaryList")
	defer span.End()

	glossaries, _, err := h.ServiceGlossary.Get(ctx, 0, 0, req.Table, req.Attribute, int(req.ValueInt), req.ValueName)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*pb.Glossary
	for _, glossary := range glossaries {
		data = append(data, &pb.Glossary{
			Id:        int32(glossary.ID),
			Table:     glossary.Table,
			Attribute: glossary.Attribute,
			ValueInt:  int32(glossary.ValueInt),
			ValueName: glossary.ValueName,
			Note:      glossary.Note,
		})
	}

	res = &pb.GetGlossaryListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *ConfigurationGrpcHandler) GetGlossaryDetail(ctx context.Context, req *pb.GetGlossaryDetailRequest) (res *pb.GetGlossaryDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetGlossaryDetail")
	defer span.End()

	glossary, err := h.ServiceGlossary.GetDetail(ctx, req.Table, req.Attribute, int(req.ValueInt), req.ValueName)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.GetGlossaryDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.Glossary{
			Id:        int32(glossary.ID),
			Table:     glossary.Table,
			Attribute: glossary.Attribute,
			ValueInt:  int32(glossary.ValueInt),
			ValueName: glossary.ValueName,
			Note:      glossary.Note,
		},
	}

	return
}

func (h *ConfigurationGrpcHandler) GetWrtDetail(ctx context.Context, req *pb.GetWrtDetailRequest) (res *pb.GetWrtDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetWrtDetail")
	defer span.End()

	var Wrt dto.WrtResponse

	Wrt, err = h.ServiceWrt.GetDetail(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.Wrt

	data = &pb.Wrt{
		Id:       Wrt.ID,
		RegionId: Wrt.RegionID,
		Code:     Wrt.Code,
		Name:     Wrt.Name,
		Type:     int32(Wrt.Type),
		Note:     Wrt.Note,
	}

	res = &pb.GetWrtDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *ConfigurationGrpcHandler) GetWrtIdGP(ctx context.Context, req *pb.GetWrtDetailRequest) (res *pb.GetWrtDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetWrtDetail")
	defer span.End()

	var Wrt dto.WrtResponse

	Wrt, err = h.ServiceWrt.GetIDDetail(ctx, req.Id)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data *pb.Wrt

	data = &pb.Wrt{
		Id: Wrt.ID,
		// RegionId: Wrt.RegionID,
		Code: Wrt.Code,
		Name: Wrt.Name,
		Type: int32(Wrt.Type),
		Note: Wrt.Note,
	}

	res = &pb.GetWrtDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *ConfigurationGrpcHandler) GetWrtList(ctx context.Context, req *pb.GetWrtListRequest) (res *pb.GetWrtListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetWrtList")
	defer span.End()

	var wrtList []dto.WrtResponse

	wrtList, _, err = h.ServiceWrt.Get(ctx, int(req.Offset), int(req.Limit), int(req.Type), req.RegionId, req.Search, 0)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	var data []*pb.Wrt
	for _, v := range wrtList {
		data = append(data, &pb.Wrt{
			Id:       v.ID,
			RegionId: v.RegionID,
			Code:     v.Code,
			Name:     v.Name,
			Type:     int32(v.Type),
			Note:     v.Note,
		})
	}

	res = &pb.GetWrtListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}
	return
}

func (h *ConfigurationGrpcHandler) GetConfigAppList(ctx context.Context, req *pb.GetConfigAppListRequest) (res *pb.GetConfigAppListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetConfigAppList")
	defer span.End()

	fmt.Println("REQ", req)

	param := &dto.ApplicationConfigRequestGet{
		Application: int8(req.Application),
		Attribute:   req.Attribute,
		Search:      req.Field,
		Limit:       req.Limit,
		Offset:      req.Offset,
		Value:       req.Value,
	}

	configApps, _, err := h.ServiceConfigApp.Get(ctx, param)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	var data []*pb.ConfigApp
	for _, configApp := range configApps {
		data = append(data, &pb.ConfigApp{
			Id:          int32(configApp.ID),
			Application: int32(configApp.Application),
			Field:       configApp.Field,
			Attribute:   configApp.Attribute,
			Value:       configApp.Value,
		})
	}

	res = &pb.GetConfigAppListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}
func (h *ConfigurationGrpcHandler) GetConfigAppDetail(ctx context.Context, req *pb.GetConfigAppDetailRequest) (res *pb.GetConfigAppDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetConfigAppDetail")
	defer span.End()

	configApp, err := h.ServiceConfigApp.GetDetail(ctx, req)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.GetConfigAppDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.ConfigApp{
			Id:          int32(configApp.ID),
			Application: int32(configApp.Application),
			Field:       configApp.Field,
			Attribute:   configApp.Attribute,
			Value:       configApp.Value,
		},
	}

	return
}

func (h *ConfigurationGrpcHandler) GetRegionPolicyDetail(ctx context.Context, req *pb.GetRegionPolicyDetailRequest) (res *pb.GetRegionPolicyDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetRegionPolicyDetail")
	defer span.End()

	regionPolicy, err := h.ServiceRegionPolicy.GetDetail(ctx, req.Id, req.Code, req.RegionId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	res = &pb.GetRegionPolicyDetailResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data: &pb.RegionPolicy{
			Id:                 regionPolicy.ID,
			Region:             regionPolicy.Region.Description,
			RegionId:           regionPolicy.Region.ID,
			OrderTimeLimit:     regionPolicy.OrderTimeLimit,
			MaxDayDeliveryDate: int32(regionPolicy.MaxDayDeliveryDate),
			WeeklyDayOff:       int32(regionPolicy.WeeklyDayOff),
			CsPhoneNumber:      regionPolicy.CSPhoneNumber,
		},
	}

	return
}

func (h *ConfigurationGrpcHandler) GetRegionPolicyList(ctx context.Context, req *pb.GetRegionPolicyListRequest) (res *pb.GetRegionPolicyListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetRegionPolicyList")
	defer span.End()

	regionPolicy, _, err := h.ServiceRegionPolicy.Get(ctx, int(req.Offset), int(req.Limit), req.Search, "", req.RegionId)
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	//fmt.Print(total)
	var data []*pb.RegionPolicy
	for _, rp := range regionPolicy {
		data = append(data, &pb.RegionPolicy{
			Id:                 rp.ID,
			Region:             rp.Region.Description,
			RegionId:           rp.Region.ID,
			OrderTimeLimit:     rp.OrderTimeLimit,
			MaxDayDeliveryDate: int32(rp.MaxDayDeliveryDate),
			WeeklyDayOff:       int32(rp.WeeklyDayOff),
			CsPhoneNumber:      rp.CSPhoneNumber,
		})
	}

	res = &pb.GetRegionPolicyListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}

func (h *ConfigurationGrpcHandler) GetDayOffDetail(ctx context.Context, req *pb.GetDayOffDetailRequest) (res *pb.GetDayOffDetailResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDayOffDetail")
	defer span.End()

	// regionPolicy, err := h.ServiceRegionPolicy.GetByID(ctx, req.Id)
	// if err != nil {
	// 	err = status.New(codes.NotFound, err.Error()).Err()
	// 	h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
	// 	return
	// }

	// res = &pb.GetRegionPolicyDetailResponse{
	// 	Code:    int32(codes.OK),
	// 	Message: codes.OK.String(),
	// 	Data: &pb.RegionPolicy{
	// 		Id:                 regionPolicy.ID,
	// 		Region:             regionPolicy.Region.Description,
	// 		RegionId:           regionPolicy.Region.ID,
	// 		OrderTimeLimit:     regionPolicy.OrderTimeLimit,
	// 		MaxDayDeliveryDate: int32(regionPolicy.MaxDayDeliveryDate),
	// 		WeeklyDayOff:       int32(regionPolicy.WeeklyDayOff),
	// 	},
	// }

	return
}

func (h *ConfigurationGrpcHandler) GetDayOffList(ctx context.Context, req *pb.GetDayOffListRequest) (res *pb.GetDayOffListResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.GetDayOffList")
	defer span.End()

	dayOff, _, err := h.ServiceDayOff.Get(ctx, int(req.Offset), int(req.Limit), 1, req.Search, "", req.StartDate.AsTime(), req.EndDate.AsTime())
	if err != nil {
		err = status.New(codes.NotFound, err.Error()).Err()
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}
	// //fmt.Print(total)
	var data []*pb.DayOff
	for _, rp := range dayOff {
		data = append(data, &pb.DayOff{
			Id:            rp.ID,
			OffDate:       timestamppb.New(rp.OffDate),
			Note:          rp.Note,
			Status:        int32(rp.Status),
			StatusConvert: rp.StatusConvert,
		})
	}

	res = &pb.GetDayOffListResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Data:    data,
	}

	return
}
