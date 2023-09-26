package handler

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type StorageHandler struct {
	Option         global.HandlerOptions
	ServiceStorage service.IStorageService
}

// URLMapping declare endpoint with handler function.
func (h *StorageHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServiceStorage = service.NewStorageService()
	cMiddleware := middleware.NewMiddleware()
	r.POST("/image", h.Image, cMiddleware.Authorized("public"))

}

func (h StorageHandler) Image(c echo.Context) (err error) {
	ctx := c.(*edenlabs.Context)

	if err := ctx.Request().ParseMultipartForm(1024); err != nil {
		err = edenlabs.ErrorValidation("file", "The file is invalid")
		return err
	}

	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		err = edenlabs.ErrorValidation("file", "The file is invalid")
		return
	}

	fileType := ctx.FormValue("type")

	contentType := fileHeader.Header.Get("Content-Type")
	if contentType != "image/jpeg" && contentType != "image/png" {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		err = edenlabs.ErrorValidation("file", "The file format is not allowed, allowed only for jpeg and png")
		return err
	}

	uploadedFile, _ := fileHeader.Open()
	if err != nil {
		return
	}
	defer uploadedFile.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, uploadedFile); err != nil {
		// return nil, err
	}
	fileName := fileHeader.Filename
	fileExtension := filepath.Ext(fileName)
	fmt.Print(fileType)
	ctx.ResponseData, err = h.ServiceStorage.UploadImage(ctx, buf.Bytes(), fileExtension, fileName)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}
