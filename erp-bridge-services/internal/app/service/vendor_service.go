package service

import (
	"context"
	"errors"
	"net/url"
	"strconv"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/repository"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"github.com/sirupsen/logrus"
)

type IVendorService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.VendorResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res dto.VendorResponse, err error)
	GetGP(ctx context.Context, req *pb.GetVendorGPListRequest) (res *pb.GetVendorGPResponse, err error)
	GetDetailGP(ctx context.Context, req *pb.GetVendorGPDetailRequest) (res *pb.GetVendorGPResponse, err error)
	CreateGP(ctx context.Context, req *dto.CreateVendorGPRequest) (res dto.CreateVendorGPResponse, err error)
}

type VendorService struct {
	opt              opt.Options
	RepositoryVendor repository.IVendorRepository
}

func NewVendorService() IVendorService {
	return &VendorService{
		opt: global.Setup.Common,
		// RepositoryVendor: repository.NewVendorRepository(),
	}
}

func (s *VendorService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []dto.VendorResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.Get")
	defer span.End()

	var vendors []*model.Vendor
	vendors, total, err = s.RepositoryVendor.Get(ctx, offset, limit, status, search, orderBy)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	for _, vendor := range vendors {
		res = append(res, dto.VendorResponse{
			ID:                     vendor.ID,
			Code:                   vendor.Code,
			VendorOrganizationID:   vendor.VendorOrganizationID,
			VendorClassificationID: vendor.VendorClassificationID,
			SubDistrictID:          vendor.SubDistrictID,
			PicName:                vendor.PicName,
			Email:                  vendor.Email,
			PhoneNumber:            vendor.PhoneNumber,
			PaymentTermID:          vendor.PaymentTermID,
			Rejectable:             vendor.Rejectable,
			Returnable:             vendor.Returnable,
			Address:                vendor.Address,
			Note:                   vendor.Note,
			Status:                 vendor.Status,
			Latitude:               vendor.Latitude,
			Longitude:              vendor.Longitude,
			CreatedAt:              vendor.CreatedAt,
			CreatedBy:              vendor.CreatedBy,
		})
	}

	return
}

func (s *VendorService) GetDetail(ctx context.Context, id int64, code string) (res dto.VendorResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.GetDetail")
	defer span.End()

	var vendor *model.Vendor
	vendor, err = s.RepositoryVendor.GetDetail(ctx, id, code)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	res = dto.VendorResponse{
		ID:                     vendor.ID,
		Code:                   vendor.Code,
		VendorOrganizationID:   vendor.VendorOrganizationID,
		VendorClassificationID: vendor.VendorClassificationID,
		SubDistrictID:          vendor.SubDistrictID,
		PicName:                vendor.PicName,
		Email:                  vendor.Email,
		PhoneNumber:            vendor.PhoneNumber,
		PaymentTermID:          vendor.PaymentTermID,
		Rejectable:             vendor.Rejectable,
		Returnable:             vendor.Returnable,
		Address:                vendor.Address,
		Note:                   vendor.Note,
		Status:                 vendor.Status,
		Latitude:               vendor.Latitude,
		Longitude:              vendor.Longitude,
		CreatedAt:              vendor.CreatedAt,
		CreatedBy:              vendor.CreatedBy,
	}

	return
}

// FOR GRPC PART
func (s *VendorService) GetGP(ctx context.Context, req *pb.GetVendorGPListRequest) (res *pb.GetVendorGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.GetGP")
	defer span.End()

	var tempVendorGP *dto.GetVendorGPList
	params := map[string]string{
		"interid":    global.EnvDatabaseGP,
		"PageNumber": strconv.Itoa(int(req.Offset)),
		"PageSize":   strconv.Itoa(int(req.Limit)),
	}

	if req.Search != "" {
		params["VENDNAME"] = req.Search
	}
	if req.Status != 0 {
		params["STATUS"] = strconv.Itoa(int(req.Status))
	}

	if req.VendorOrg != "" {
		req.VendorOrg = url.PathEscape(req.VendorOrg)
		params["ORGANIZATION"] = req.VendorOrg
	}

	if req.Status != 0 {
		params["STATUS"] = strconv.Itoa(int(req.Status))
	}

	if req.Orderby != "" {
		params["orderby"] = req.Orderby
	}

	if req.State != "" {
		req.State = url.PathEscape(req.State)
		params["state"] = req.State
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "vendor/list", nil, &tempVendorGP, params)

	var tempData []*pb.VendorGP
	for _, v := range tempVendorGP.Data {
		tempData = append(tempData, &pb.VendorGP{
			VENDORID:           v.VENDORID,
			VENDNAME:           v.VENDNAME,
			VNDCLSID:           v.VNDCLSID,
			VNDCNTCT:           v.VNDCNTCT,
			ADDRESS1:           v.ADDRESS1,
			ADDRESS2:           v.ADDRESS2,
			ADDRESS3:           v.ADDRESS3,
			CITY:               v.CITY,
			STATE:              v.STATE,
			ZIPCODE:            v.ZIPCODE,
			COUNTRY:            v.COUNTRY,
			PHNUMBR1:           v.PHNUMBR1,
			PHNUMBR2:           v.PHNUMBR2,
			PHONE3:             v.PHONE3,
			FAXNUMBR:           v.FAXNUMBR,
			UPSZONE:            v.UPSZONE,
			SHIPMTHD:           v.SHIPMTHD,
			TAXSCHID:           v.TAXSCHID,
			ACNMVNDR:           v.ACNMVNDR,
			TXIDNMBR:           v.TXIDNMBR,
			VENDSTTS:           v.VENDSTTS,
			CREATDDT:           v.CREATDDT,
			CURNCYID:           v.CURNCYID,
			TXRGNNUM:           v.TXRGNNUM,
			TRDDISCT:           v.TRDDISCT,
			MINORDER:           v.MINORDER,
			PYMTRMID:           v.PYMTRMID,
			COMMENT1:           v.COMMENT1,
			COMMENT2:           v.COMMENT2,
			USERDEF1:           v.USERDEF1,
			USERDEF2:           v.USERDEF2,
			PYMNTPRI:           v.PYMNTPRI,
			Organization:       v.Organization,
			Classification:     v.Classification,
			PaymentMethod:      v.PaymentMethod,
			LatestGoodsReceipt: v.LatestGoodsReceipt,
			Vaddcdpr: &pb.VendorGP_VVendorAddress{
				Vendorid: v.Vaddcdpr.Vendorid,
				Adrscode: v.Vaddcdpr.Adrscode,
				PrpAdministrativeCode: &pb.VendorGP_AdministrativeCode{
					GnlAdministrativeCode: v.Vaddcdpr.PrpAdministrativeCode,
				},
			},
		})
	}
	res = &pb.GetVendorGPResponse{
		PageNumber:   tempVendorGP.PageNumber,
		PageSize:     tempVendorGP.PageSize,
		TotalPages:   tempVendorGP.TotalPages,
		TotalRecords: tempVendorGP.TotalRecords,
		Succeeded:    tempVendorGP.Succeeded,
		Errors:       tempVendorGP.Errors,
		Message:      tempVendorGP.Message,
		Data:         tempData,
	}
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	// fmt.Println(res)
	return
}

func (s *VendorService) GetDetailGP(ctx context.Context, req *pb.GetVendorGPDetailRequest) (res *pb.GetVendorGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.GetDetailGP")
	defer span.End()

	params := map[string]string{
		"interid": global.EnvDatabaseGP,
		"id":      req.Id,
	}

	err = global.HttpRestApiToMicrosoftGP("GET", "vendor/detail", nil, &res, params)

	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}

	return
}

func (s *VendorService) CreateGP(ctx context.Context, req *dto.CreateVendorGPRequest) (res dto.CreateVendorGPResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.Create")
	defer span.End()

	req.InterID = global.EnvDatabaseGP

	err = global.HttpRestApiToMicrosoftGP("POST", "vendor/create", req, &res, nil)
	if err != nil {
		logrus.Error(err.Error())
		err = errors.New("Connection to the server could not be established")
		return
	}

	if res.Code != 200 {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = errors.New("Connection to the server could not be established")
		return
	}

	return
}
