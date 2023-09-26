package upload

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"git.edenfarm.id/cuxs/cuxs"
	"github.com/labstack/echo/v4"

	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
)

type Handler struct{}

func (h *Handler) URLMapping(r *echo.Group) {

	r.POST("", h.funcUploadFile, auth.Authorized("filter_rdl"))
	r.POST("/img", h.funcUploadImageFile, auth.Authorized("filter_rdl"))
	r.POST("/pdf", h.funcUploadPDFFile, auth.Authorized("filter_rdl"))
	r.POST("/img/field-purchaser", h.funcUploadImageFileFieldPurchaser, auth.AuthorizedFieldPurchaserMobile())
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

	values := map[string]io.Reader{
		"file": util.OpenFile(handler.Filename), // this file
		"type": strings.NewReader(typeRequest),  // there's have value (/product/, /product-tag/, /report/,/magang/) this for directory target to save
	}

	if err != nil {
		return err
	}
	client := &http.Client{}
	if err = util.PostFormValue(client, values); err == nil {
		if typeRequest == "product_tag" {
			typeRequest = strings.ReplaceAll(typeRequest, "_", "-")
		}
		ctx.ResponseData, err = util.ResponseURLFromPostUpload+typeRequest+"/"+filename, err
	}
	os.Remove(fileLocation)

	return ctx.Serve(err)
}

//funcUploadImageFile : function to upload image to S3 storage
func (h *Handler) funcUploadImageFile(r echo.Context) (e error) {
	ctx := r.(*cuxs.Context)
	if err := r.Request().ParseMultipartForm(1024); err != nil {
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
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	validFileExtensions := [...]string{".jpg", ".jpeg", ".png"}
	isFileExtensionValid := util.ItemExists(validFileExtensions, fileExtension)

	if !isFileExtensionValid {
		err = errors.New("The provided file format is not allowed. Please upload a JPEG or PNG image")
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
	//fileNameExtension := filepath.Ext(fileName)
	//nameFormat := strconv.Itoa(int(time.Now().UnixNano()))
	//switch typeRequest {
	//case "product":
	//	fileName = "product_img_" + nameFormat + fileNameExtension
	//case "mobile-apps":
	//	fileName = "mobile_apps_img_" + nameFormat + fileNameExtension
	//case "product-tag":
	//	fileName = "product_tag_img_" + nameFormat + fileNameExtension
	//default:
	//	fileName = "eden_img_" + nameFormat + fileNameExtension
	//}
	fl, err := util.UploadImageToS3(fileName, fileLocation, typeRequest)
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

	fl, err := util.UploadImageToS3(fileName, fileLocation, typeRequest)
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
	isFileExtensionValid := util.ItemExists(validFileExtensions, fileExtension)

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

	fl, err := util.UploadImageToS3FieldPurchaser(fileName, fileLocation, typeRequest, "no")
	if err != nil {
		return err
	}
	ctx.Data(fl)

	os.Remove(fileLocation)

	return ctx.Serve(err)
}
