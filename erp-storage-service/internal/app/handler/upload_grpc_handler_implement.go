package handler

import (
	"bytes"
	"context"
	"fmt"

	// "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/storage_service"

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/storage_service"
	"google.golang.org/grpc/codes"
)

// func (h *StorageGrpcHandler) UploadImageGRPCStream(stream pb.StorageService_UploadImageGRPCStreamServer) (err error) {
// 	var ctx context.Context
// 	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateSalesInvoiceExternal")
// 	defer span.End()
// 	imageData := bytes.Buffer{}
// 	imageSize := 0
// 	tempUpload := &pb.UploadImageGRPCStreamRequest{}
// 	const maxImageSize = 1 << 20
// 	for {
// 		// log.Print("waiting to receive more data")

// 		req, err := stream.Recv()
// 		tempUpload = req
// 		if err == io.EOF {
// 			// log.Print("no more data")
// 			break
// 		}
// 		if err != nil {
// 			err = status.New(codes.NotFound, err.Error()).Err()
// 			h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()

// 			break
// 		}
// 		chunk := req.GetContent()
// 		size := len(chunk)
// 		// log.Printf("received a chunk with size: %d", size)

// 		imageSize += size
// 		if imageSize > maxImageSize {
// 			// return logError(status.Errorf(codes.InvalidArgument, "image is too large: %d > %d", imageSize, maxImageSize))
// 		}
// 		_, err = imageData.Write(chunk)
// 		if err != nil {
// 			// return logError(status.Errorf(codes.Internal, "cannot write chunk data: %v", err))
// 		}
// 	}
// 	if err != nil {
// 		err = status.New(codes.NotFound, err.Error()).Err()
// 		h.Option.Common.Logger.AddMessage(log.ErrorLevel, err).Print()
// 		return
// 	}
// 	url, err := h.ServicesUpload.UploadImageStream(ctx, imageData, tempUpload)
// 	fmt.Println(imageData, imageSize)

// 	fmt.Print(url)
// 	res := &storage_service.UploadImageGRPCStreamResponse{
// 		Code:    int32(codes.OK),
// 		Message: codes.OK.String(),
// 		Url:     url,
// 	}
// 	err = stream.SendAndClose(res)
// 	if err != nil {
// 		// return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
// 	}

// 	return
// }

func (h *StorageGrpcHandler) UploadImageGRPC(ctx context.Context, req *pb.UploadImageGRPCStreamRequest) (res *pb.UploadImageGRPCStreamResponse, err error) {
	ctx, span := h.Option.Common.Trace.Start(ctx, "Grpc.CreateSalesInvoiceExternal")
	defer span.End()
	imageData := bytes.Buffer{}
	imageSize := 0
	// tempUpload := &pb.UploadImageGRPCStreamRequest{}

	imageData = *bytes.NewBuffer(req.Content)
	url, err := h.ServicesUpload.UploadImageStream(ctx, imageData, req)
	fmt.Println(imageData, imageSize)

	fmt.Print(url)
	res = &pb.UploadImageGRPCStreamResponse{
		Code:    int32(codes.OK),
		Message: codes.OK.String(),
		Url:     url,
	}
	if err != nil {
		// return logError(status.Errorf(codes.Unknown, "cannot send response: %v", err))
	}

	return
}
