package handler

import (
	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/global"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/middleware"
	"git.edenfarm.id/project-version3/erp-services/erp-edn-mobile-service/internal/app/service"
	"github.com/labstack/echo/v4"
)

type UploadHandler struct {
	Option         global.HandlerOptions
	ServicesUpload service.IUploadService
}

// URLMapping implements router.RouteHandlers
func (h *UploadHandler) URLMapping(r *echo.Group) {
	h.Option = global.Setup
	h.ServicesUpload = service.NewUploadService()

	cMiddleware := middleware.NewMiddleware()
	r.POST("/imageKTP", h.ImageKTP, cMiddleware.Authorized())
	r.POST("/imageAddress", h.ImageAddress, cMiddleware.Authorized())
}

func (h UploadHandler) ImageKTP(c echo.Context) (err error) {
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

	ctx.ResponseData, err = h.ServicesUpload.ImageKTP(ctx.Request().Context(), *fileHeader, fileType)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}

func (h UploadHandler) ImageAddress(c echo.Context) (err error) {
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

	ctx.ResponseData, err = h.ServicesUpload.ImageAddress(ctx.Request().Context(), *fileHeader, fileType)
	if err != nil {
		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
		return
	}

	return ctx.Serve(err)
}
