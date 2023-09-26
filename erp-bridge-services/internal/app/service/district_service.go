package service

import (
	"context"
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/edenlabs/edenlabs/statusx"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/model"
)

type IDistrictService interface {
	Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.DistrictResponse, total int64, err error)
	GetDetail(ctx context.Context, id int64, code string) (res *dto.DistrictResponse, err error)
	GetInIds(ctx context.Context, offset int, limit int, ids []int64, status int, search string, orderBy string) (res []*dto.DistrictResponse, total int64, err error)
}

type DistrictService struct {
	opt opt.Options
}

func NewDistrictService() IDistrictService {
	return &DistrictService{
		opt: global.Setup.Common,
	}
}

func (s *DistrictService) Get(ctx context.Context, offset int, limit int, status int, search string, orderBy string) (res []*dto.DistrictResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DistrictService.Get")
	defer span.End()

	for _, district := range s.MockDatas(100) {
		res = append(res, &dto.DistrictResponse{
			ID:            district.ID,
			Code:          district.Code,
			Value:         district.Value,
			Name:          district.Name,
			Note:          district.Note,
			Status:        district.Status,
			StatusConvert: district.StatusConvert,
		})
	}

	res = res[offset : offset+limit]
	return
}

func (s *DistrictService) GetDetail(ctx context.Context, id int64, code string) (res *dto.DistrictResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DistrictService.GetDetail")
	defer span.End()

	district := s.MockDatas(100)[id-1]
	res = &dto.DistrictResponse{
		ID:            district.ID,
		Code:          district.Code,
		Value:         district.Value,
		Name:          district.Name,
		Note:          district.Note,
		Status:        district.Status,
		StatusConvert: district.StatusConvert,
	}

	return
}

func (s *DistrictService) GetInIds(ctx context.Context, offset int, limit int, ids []int64, status int, search string, orderBy string) (res []*dto.DistrictResponse, total int64, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "DistrictService.Get")
	defer span.End()

	for _, district := range s.MockDatas(100) {

		res = append(res, &dto.DistrictResponse{
			ID:            district.ID,
			Code:          district.Code,
			Value:         district.Value,
			Name:          district.Name,
			Note:          district.Note,
			Status:        district.Status,
			StatusConvert: district.StatusConvert,
		})
	}

	res = res[offset : offset+limit]
	return
}

func (r *DistrictService) MockDatas(total int) (mockDatas []*model.District) {
	for i := 1; i <= total; i++ {
		mockDatas = append(mockDatas,
			&model.District{
				ID:            int64(i),
				Code:          fmt.Sprintf("DIS%d", i),
				Value:         fmt.Sprintf("DUMMY Value District %d", i),
				Name:          fmt.Sprintf("DUMMY Value Name %d", i),
				Note:          fmt.Sprintf("DUMMY Note District %d", i),
				Status:        1,
				StatusConvert: statusx.ConvertStatusValue(1),
			})
	}
	return
}
