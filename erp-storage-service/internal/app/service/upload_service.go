package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	"git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/storage_service"
	"git.edenfarm.id/project-version3/erp-services/erp-storage-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-storage-service/internal/app/dto"
)

type IUploadService interface {
	Image(ctx context.Context, fileHeader multipart.FileHeader, fileType string) (res dto.UploadResponse, err error)
	UploadImageStream(ctx context.Context, imageByte bytes.Buffer, req *storage_service.UploadImageGRPCStreamRequest) (res string, err error)
	// UploadImage(ctx context.Context, imageByte bytes.Buffer, req *storage_service.UploadImageGRPCStreamRequest) (res string, err error)
	File(ctx context.Context, fileHeader multipart.FileHeader, fileType string) (res dto.UploadResponse, err error)
}

type UploadService struct {
	opt opt.Options
}

func NewUploadService() IUploadService {
	return &UploadService{
		opt: global.Setup.Common,
	}
}

func (s *UploadService) Image(ctx context.Context, fileHeader multipart.FileHeader, fileType string) (res dto.UploadResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UploadService.Image")
	defer span.End()

	uploadedFile, _ := fileHeader.Open()
	if err != nil {
		return
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	fileName := fileHeader.Filename
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	var validFileExt bool

	switch fileExtension {
	case ".jpg":
		validFileExt = true
	case ".jpeg":
		validFileExt = true
	case ".png":
		validFileExt = true
	default:
		validFileExt = false
	}

	if !validFileExt {
		err = edenlabs.ErrorValidation("file", "The file format is not allowed, allowed only for jpeg and png")
		return
	}

	fileLocation := filepath.Join(dir, "", fileName)
	targetFile, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		err = fmt.Errorf("failed to open file | %v", err)
		return
	}
	defer targetFile.Close()

	if _, err = io.Copy(targetFile, uploadedFile); err != nil {
		err = fmt.Errorf("failed to copy file | %v", err)
		return
	}

	info, err := s.opt.S3x.UploadPublicFile(ctx, s.opt.Config.S3.BucketName, fileName, fileLocation, fileType)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = fmt.Errorf("failed to upload file | %v", err)
		return
	}

	os.Remove(fileLocation)

	res = dto.UploadResponse{
		Url: info,
	}

	return
}

func (s *UploadService) UploadImageStream(ctx context.Context, imageByte bytes.Buffer, req *storage_service.UploadImageGRPCStreamRequest) (res string, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UploadService.Image")
	defer span.End()

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	fileName := req.FileName
	fileExtension := strings.ToLower(req.FileType)
	var validFileExt bool

	switch fileExtension {
	case ".jpg":
		validFileExt = true
	case ".jpeg":
		validFileExt = true
	case ".png":
		validFileExt = true
	default:
		validFileExt = false
	}

	if !validFileExt {
		err = edenlabs.ErrorValidation("file", "The file format is not allowed, allowed only for jpeg and png")
		return
	}

	fileLocation := filepath.Join(dir, "", fileName)
	targetFile, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		err = fmt.Errorf("failed to open file | %v", err)
		return
	}
	defer targetFile.Close()

	if _, err = io.Copy(targetFile, &imageByte); err != nil {
		err = fmt.Errorf("failed to copy file | %v", err)
		return
	}

	info, err := s.opt.S3x.UploadPublicFile(ctx, s.opt.Config.S3.BucketName, fileName, fileLocation, req.FileType)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = fmt.Errorf("failed to upload file | %v", err)
		return
	}

	os.Remove(fileLocation)

	// res = dto.UploadResponse{
	// 	Url: info,
	// }

	res = info

	return
}

func (s *UploadService) File(ctx context.Context, fileHeader multipart.FileHeader, fileType string) (res dto.UploadResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "UploadService.Image")
	defer span.End()

	uploadedFile, _ := fileHeader.Open()
	if err != nil {
		return
	}
	defer uploadedFile.Close()

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	fileName := fileHeader.Filename
	fileExtension := strings.ToLower(filepath.Ext(fileName))
	var validFileExt bool

	switch fileExtension {
	case ".pdf":
		validFileExt = true
	default:
		validFileExt = false
	}

	if !validFileExt {
		err = edenlabs.ErrorValidation("file", "The file format is not allowed, allowed only for pdf")
		return
	}

	fileLocation := filepath.Join(dir, "", fileName)
	targetFile, err := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
	if err != nil {
		err = fmt.Errorf("failed to open file | %v", err)
		return
	}
	defer targetFile.Close()

	if _, err = io.Copy(targetFile, uploadedFile); err != nil {
		err = fmt.Errorf("failed to copy file | %v", err)
		return
	}

	info, err := s.opt.S3x.UploadPublicFile(ctx, s.opt.Config.S3.BucketName, fileName, fileLocation, fileType)
	if err != nil {
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = fmt.Errorf("failed to upload file | %v", err)
		return
	}

	os.Remove(fileLocation)

	res = dto.UploadResponse{
		Url: info,
	}

	return
}
