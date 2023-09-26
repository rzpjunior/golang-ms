package service

import (
	"context"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/configuration_service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type IDeliveryDateService interface {
	Get(ctx context.Context, req dto.DeliveryDateRequest) (res dto.DeliveryDateResponse, err error)
}

type DeliveryDateService struct {
	opt opt.Options
	//RepositoryDeliveryDate repository.IDeliveryDateRepository
}

func NewDeliveryDateService() IDeliveryDateService {
	return &DeliveryDateService{
		opt: global.Setup.Common,
		//RepositoryDeliveryDate: repository.NewDeliveryDateRepository(),
	}
}

func (s *DeliveryDateService) Get(ctx context.Context, req dto.DeliveryDateRequest) (res dto.DeliveryDateResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DeliveryDateService.GetByID")
	defer span.End()

	//cek data area id ada ga areanya,kalau ada ambil region policy.kalau uda ambil region policy nanti balikin yang ada di depan
	// regionID, _ := strconv.Atoi(req.Data.RegionID)
	Region, err := s.opt.Client.BridgeServiceGrpc.GetAdmDivisionGPList(ctx, &bridge_service.GetAdmDivisionGPListRequest{
		Region: req.Data.RegionID,
		Limit:  1,
		Offset: 0,
	})

	if err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		return
	}
	if Region.Data[0].Code == "" {
		//throw error
	}

	var regionPolicy *model.RegionPolicy
	RegionPol, err := s.opt.Client.ConfigurationServiceGrpc.GetRegionPolicyDetail(ctx, &configuration_service.GetRegionPolicyDetailRequest{
		RegionId: req.Data.RegionID,
	})

	regionPolicy = &model.RegionPolicy{
		ID:                 int64(RegionPol.Data.Id),
		Region:             RegionPol.Data.Region,
		OrderTimeLimit:     RegionPol.Data.OrderTimeLimit,
		MaxDayDeliveryDate: int(RegionPol.Data.MaxDayDeliveryDate),
		WeeklyDayOff:       int(RegionPol.Data.WeeklyDayOff),
	}

	var bunchDate []string
	var weeklyDayOff int
	isDayOff := false

	// in go sunday is 0, replace manual 7 to 0
	if regionPolicy.WeeklyDayOff == 7 {
		weeklyDayOff = 0
	} else {
		weeklyDayOff = regionPolicy.WeeklyDayOff
	}
	date := time.Now().AddDate(0, 0, 1)
	start_date := time.Now().AddDate(0, 0, 1)
	end_date := time.Now().AddDate(0, 0, 30)

	var dayOffs []*model.DayOff

	dayOff, err := s.opt.Client.ConfigurationServiceGrpc.GetDayOffList(ctx, &configuration_service.GetDayOffListRequest{
		StartDate: timestamppb.New(start_date),
		EndDate:   timestamppb.New(end_date),
	})

	for _, do := range dayOff.Data {
		dayOffs = append(dayOffs, &model.DayOff{
			ID:      do.Id,
			OffDate: do.OffDate.AsTime(),
			Note:    do.Note,
			Status:  int8(do.Status),
		})
	}
	i := 0
	for i < regionPolicy.MaxDayDeliveryDate {
		//check if date is weekly day off
		if int(date.Weekday()) != weeklyDayOff {
			//check if date is listed in day off data within 30 days
			for j := 0; j < len(dayOffs); j++ {
				if date.Format("2006-01-02") == dayOffs[j].OffDate.Format("2006-01-02") {
					isDayOff = true
					break
				}
			}

			if isDayOff == false {
				bunchDate = append(bunchDate, date.Format("2006-01-02"))
				i++
			}

			isDayOff = false
		}
		date = date.AddDate(0, 0, 1)

	}
	res.Date = bunchDate

	return
}
