package service

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-purchaser-mobile-service/internal/app/dto"
)

type IVendorService interface {
	Get(ctx context.Context, req *dto.VendorListRequest) (res []*dto.VendorResponse, total int64, err error)
	GetByID(ctx context.Context, id string) (res *dto.VendorResponse, err error)
	Create(ctx context.Context, req *dto.VendorRequestCreate) (res *dto.VendorRequestCreateResponse, err error)
}

type VendorService struct {
	opt opt.Options
}

func NewVendorService() IVendorService {
	return &VendorService{
		opt: global.Setup.Common,
	}
}

func (s *VendorService) Get(ctx context.Context, req *dto.VendorListRequest) (res []*dto.VendorResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.Get")
	defer span.End()

	var statusGP int32
	switch statusx.ConvertStatusValue(int8(req.Status)) {
	case statusx.Archived:
		statusGP = 2
	case statusx.Active:
		statusGP = 1
	}

	var vendors *bridgeService.GetVendorGPResponse
	vendors, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPList(ctx, &bridgeService.GetVendorGPListRequest{
		Limit:     req.Limit,
		Offset:    req.Offset,
		Status:    statusGP,
		Search:    req.Search,
		Orderby:   "desc",
		VendorOrg: req.VendorOrg,
		State:     req.State,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
		return
	}

	if len(vendors.Data) == 0 {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
		return
	}

	switch vendors.Data[0].VENDSTTS {
	case 2:
		statusGP = 7
	case 1:
		statusGP = 1
	}

	for _, vendor := range vendors.Data {

		res = append(res, &dto.VendorResponse{
			ID:   vendor.VENDORID,
			Code: vendor.VENDORID,
			VendorOrganization: &dto.VendorOrganizationResponse{
				ID:          vendor.Organization.PRP_Vendor_Org_ID,
				Code:        vendor.Organization.PRP_Vendor_Org_ID,
				Description: vendor.Organization.PRP_Vendor_Org_Desc,
			},
			VendorClassification: &dto.VendorClassificationResponse{
				ID:          vendor.Classification.PRP_Vendor_CLASF_ID,
				Code:        vendor.Classification.PRP_Vendor_CLASF_ID,
				Description: vendor.Classification.PRP_Vendor_CLASF_Desc,
			},
			AdmDivision: &dto.AdmDivisionResponse{
				Code:        vendor.Vaddcdpr.PrpAdministrativeCode.GnlAdministrativeCode,
				Province:    vendor.Vaddcdpr.PrpAdministrativeCode.GnlProvince,
				Region:      vendor.Vaddcdpr.PrpAdministrativeCode.GnlRegion,
				City:        vendor.Vaddcdpr.PrpAdministrativeCode.GnlCity,
				District:    vendor.Vaddcdpr.PrpAdministrativeCode.GnlDistrict,
				SubDistrict: vendor.Vaddcdpr.PrpAdministrativeCode.GnlSubdistrict,
			},

			// add payment term
			PaymentTerm: &dto.PaymentTermResponse{
				Id:        vendor.PYMTRMID.PYMTRMID,
				Code:      vendor.PYMTRMID.PYMTRMID,
				DaysValue: int64(vendor.PYMTRMID.Calculatedatefromdays),
			},
			// add payment method
			PaymentMethod: &dto.PaymentMethodResponse{
				ID:     vendor.PaymentMethod.PRP_Payment_Method,
				Code:   vendor.PaymentMethod.PRP_Payment_Method,
				Name:   vendor.PaymentMethod.PRP_Payment_Method_Desc,
				Status: int8(vendor.PaymentMethod.INACTIVE), // show active = 0 / inactive = 1
			},
			Name:           vendor.VENDNAME,
			PicName:        vendor.VNDCNTCT,
			PhoneNumber:    vendor.PHNUMBR1,
			PhoneNumberAlt: vendor.PHNUMBR2,
			// Rejectable:           1,
			// Returnable:           1,
			Address: vendor.ADDRESS1,
			// Note:                 vendor.ADDRESS1,
			Status: statusGP,
			// Latitude:  vendor.Latitude,
			// Longitude: vendor.Longitude,
			CreatedAt: time.Now(),
			CreatedBy: 1,
		})
	}
	total = int64(len(res))

	return
}

func (s *VendorService) GetByID(ctx context.Context, id string) (res *dto.VendorResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.GetByID")
	defer span.End()

	var statusGP int32

	var vendor *bridgeService.GetVendorGPResponse
	vendor, err = s.opt.Client.BridgeServiceGrpc.GetVendorGPDetail(ctx, &bridgeService.GetVendorGPDetailRequest{
		Id: id,
	})
	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "vendor")
		return
	}
	switch vendor.Data[0].VENDSTTS {
	case 2:
		statusGP = 7
	case 1:
		statusGP = 1
	}

	res = &dto.VendorResponse{
		ID:   vendor.Data[0].VENDORID,
		Code: vendor.Data[0].VENDORID,
		VendorOrganization: &dto.VendorOrganizationResponse{
			ID:          vendor.Data[0].Organization.PRP_Vendor_Org_ID,
			Code:        vendor.Data[0].Organization.PRP_Vendor_Org_ID,
			Description: vendor.Data[0].Organization.PRP_Vendor_Org_Desc,
		},
		VendorClassification: &dto.VendorClassificationResponse{
			ID:          vendor.Data[0].Classification.PRP_Vendor_CLASF_ID,
			Code:        vendor.Data[0].Classification.PRP_Vendor_CLASF_ID,
			Description: vendor.Data[0].Classification.PRP_Vendor_CLASF_Desc,
		},
		AdmDivision: &dto.AdmDivisionResponse{
			Code:        vendor.Data[0].Vaddcdpr.PrpAdministrativeCode.GnlAdministrativeCode,
			Province:    vendor.Data[0].Vaddcdpr.PrpAdministrativeCode.GnlProvince,
			Region:      vendor.Data[0].Vaddcdpr.PrpAdministrativeCode.GnlRegion,
			City:        vendor.Data[0].Vaddcdpr.PrpAdministrativeCode.GnlCity,
			District:    vendor.Data[0].Vaddcdpr.PrpAdministrativeCode.GnlDistrict,
			SubDistrict: vendor.Data[0].Vaddcdpr.PrpAdministrativeCode.GnlSubdistrict,
		},
		// add payment term
		PaymentTerm: &dto.PaymentTermResponse{
			Id:        vendor.Data[0].PYMTRMID.PYMTRMID,
			Code:      vendor.Data[0].PYMTRMID.PYMTRMID,
			DaysValue: int64(vendor.Data[0].PYMTRMID.Calculatedatefromdays),
		},
		// add payment method
		PaymentMethod: &dto.PaymentMethodResponse{
			ID:     vendor.Data[0].PaymentMethod.PRP_Payment_Method,
			Code:   vendor.Data[0].PaymentMethod.PRP_Payment_Method,
			Name:   vendor.Data[0].PaymentMethod.PRP_Payment_Method_Desc,
			Status: int8(vendor.Data[0].PaymentMethod.INACTIVE), // show active = 0 / inactive = 1
		},
		Name:           vendor.Data[0].VENDNAME,
		PicName:        vendor.Data[0].VNDCNTCT,
		PhoneNumber:    vendor.Data[0].PHNUMBR1,
		PhoneNumberAlt: vendor.Data[0].PHNUMBR2,
		// Rejectable:           1,
		// Returnable:           1,
		Address: vendor.Data[0].ADDRESS1,
		Note:    vendor.Data[0].COMMENT1,
		Status:  statusGP,
		// Latitude:  vendor.Latitude,
		// Longitude: vendor.Longitude,
		CreatedAt: time.Now(),
		CreatedBy: 1,
	}

	return
}

func (s *VendorService) Create(ctx context.Context, req *dto.VendorRequestCreate) (res *dto.VendorRequestCreateResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "VendorService.Create")
	defer span.End()

	randomNumber := rand.Int31()

	if len(req.Address) > 60 {
		err = edenlabs.ErrorValidation("address", "Address can not more than 60 character")
		return
	}
	if len(req.BlockNumber) > 60 {
		err = edenlabs.ErrorValidation("block_number", "Block number can not more than 60 character")
		return
	}
	vendorCreateReq := &bridgeService.CreateVendorRequest{
		Vendorid:            fmt.Sprintf("VEND%d", randomNumber),
		Vendname:            req.Name,
		Vendshnm:            req.Name,
		Vndchknm:            req.Name,
		Vendstts:            "1",
		PrP_Vendor_Org_ID:   req.VendorOrganizationID,
		PrP_Vendor_CLASF_ID: req.VendorClassificationID,
		Vndcntct:            req.PicName,
		AddresS1:            req.Address,
		AddresS2:            req.BlockNumber,
		PhnumbR1:            req.PhoneNumber,
		PhnumbR2:            req.AltPhoneNumber,
		PRP_Payment_Method:  req.PaymentMethodID,
		Vaddcdpr:            "PRIMARY",
		Vadcdpad:            "PRIMARY",
		Vadcdsfr:            "WAREHOUSE",
		Vadcdtro:            "PRIMARY",
		Pymtrmid:            req.PaymentTermID,
	}

	_, err = s.opt.Client.BridgeServiceGrpc.CreateVendor(ctx, vendorCreateReq)

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "create_vendor")
		return
	}

	return
}
