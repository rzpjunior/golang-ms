package service

import (
	"fmt"

	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/storage_service"
	"github.com/labstack/echo/v4"
)

func ServiceStorage() IStorageService {
	m := new(StorageService)
	m.opt = global.Setup.Common
	return m
}

type IStorageService interface {
	UploadImage(ctx echo.Context, image []byte, fileType string, fileName string) (url string, err error)
}

type StorageService struct {
	opt opt.Options
}

func NewStorageService() IStorageService {
	return &StorageService{
		opt: global.Setup.Common,
	}
}

func (s *StorageService) UploadImage(ctx echo.Context, image []byte, fileType string, fileName string) (url string, err error) {
	c := ctx.Request().Context()
	c, span := s.opt.Trace.Start(c, "StorageService.DeleteAccount")
	fmt.Print(span)

	urlImage, err := s.opt.Client.StorageServiceGrpc.UploadImageGRPC(ctx.Request().Context(), &storage_service.UploadImageGRPCStreamRequest{
		FileType: fileType,
		Content:  image,
		FileName: fileName,
	})
	url = urlImage.Url
	fmt.Println(url, err)
	return
}
