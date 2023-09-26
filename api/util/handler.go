package util

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"
)

type Handler struct{}

func (h *Handler) URLMapping(r *echo.Group) {
	r.POST("", h.funcUploadFile)
	r.POST("/img", h.funcUploadImageFile)
	r.POST("/pdf", h.funcUploadPDFFile)
	r.POST("/img/field-purchaser", h.funcUploadImageFileFieldPurchaser)
}

func (h *Handler) funcUploadFile(r echo.Context) (e error) {
	ctx := r.(*cuxs.Context)
	if err := r.Request().ParseMultipartForm(1024); err != nil {
		return err
	}

	alias := r.FormValue("alias")
	typeRequest := r.FormValue("type")

	handler, err := r.FormFile("file")
	uploadedFile, _ := handler.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	filename := handler.Filename
	// alias = for rename file
	if alias != "" {
		filename = fmt.Sprintf("%s%s", alias, filepath.Ext(handler.Filename))
	}
	fileLocation := filepath.Join(dir, "", filename)
	targetFile, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return err
	}

	if err != nil {
		return err
	}

	fl, err := UploadImageToS3(filename, fileLocation, typeRequest)
	if err != nil {
		return err
	}

	ctx.Data(fl)
	os.Remove(fileLocation)

	return ctx.Serve(err)
}

//funcUploadImageFile : function to upload image to S3 storage
func (h *Handler) funcUploadImageFile(r echo.Context) (e error) {
	ctx := r.(*cuxs.Context)
	if err := r.Request().ParseMultipartForm(2048); err != nil {
		return err
	}

	handler, err := r.FormFile("file")
	typeRequest := r.FormValue("type")
	fileType := handler.Header.Get("Content-Type")

	if fileType != "image/jpeg" && fileType != "image/png" {
		err = errors.New("The provided file format is not allowed. Please upload a JPEG or PNG image")
		return err
	}

	uploadedFile, _ := handler.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fileName := handler.Filename
	fileLocation := filepath.Join(dir, "", fileName)
	targetFile, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return err
	}

	fl, err := UploadImageToS3(fileName, fileLocation, typeRequest)
	if err != nil {
		return err
	}
	ctx.Data(fl)

	os.Remove(fileLocation)

	return ctx.Serve(err)
}

//funcUploadPDFFile : function to upload PDF to S3 storage
func (h *Handler) funcUploadPDFFile(r echo.Context) (e error) {
	ctx := r.(*cuxs.Context)
	if err := r.Request().ParseMultipartForm(1024); err != nil {
		return err
	}

	handler, err := r.FormFile("file")
	typeRequest := r.FormValue("type")
	fileType := handler.Header.Get("Content-Type")

	if fileType != "application/pdf" {
		err = errors.New("The provided file format is not allowed. Please upload a PDF file.")
		return err
	}

	uploadedFile, _ := handler.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fileName := handler.Filename
	fileLocation := filepath.Join(dir, "", fileName)
	targetFile, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return err
	}

	fl, err := UploadImageToS3(fileName, fileLocation, typeRequest)
	if err != nil {
		return err
	}
	ctx.Data(fl)

	os.Remove(fileLocation)

	return ctx.Serve(err)
}

//funcUploadImageFileFieldPurchaser : function to upload image to S3 storage for Field Purchaser
func (h *Handler) funcUploadImageFileFieldPurchaser(r echo.Context) (e error) {
	ctx := r.(*cuxs.Context)
	if err := r.Request().ParseMultipartForm(1024); err != nil {
		return err
	}

	handler, err := r.FormFile("file")
	typeRequest := "field-purchaser"
	fileType := handler.Header.Get("Content-Type")

	if fileType != "image/jpeg" && fileType != "image/png" {
		err = errors.New("the provided file format is not allowed. Please upload a JPEG or PNG image")
		return err
	}

	uploadedFile, _ := handler.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	fileName := handler.Filename
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	validFileExtensions := [...]string{".jpg", ".jpeg", ".png"}
	isFileExtensionValid := ItemExists(validFileExtensions, fileExtension)

	if !isFileExtensionValid {
		err = errors.New("the provided file format is not allowed. Please upload a JPEG or PNG image")
		return err
	}

	fileLocation := filepath.Join(dir, "", fileName)
	targetFile, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	if _, err := io.Copy(targetFile, uploadedFile); err != nil {
		return err
	}

	fl, err := UploadImageToS3FieldPurchaser(fileName, fileLocation, typeRequest, "no")
	if err != nil {
		return err
	}
	ctx.Data(fl)

	os.Remove(fileLocation)

	return ctx.Serve(err)
}
