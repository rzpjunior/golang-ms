package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs"
	"git.edenfarm.id/edenlabs/edenlabs/log"
	"git.edenfarm.id/edenlabs/edenlabs/opt"
	bridgeService "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/global"
	dto "git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/dto"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/model"
	"git.edenfarm.id/project-version3/erp-services/erp-site-service/internal/app/repository"
)

type IPrintLabelService interface {
	Get(ctx context.Context, req *dto.PrintLabelGetRequest) (res *dto.PrintLabelGetResponse, err error)
	GetDO(ctx context.Context, req *dto.PrintLabelGetRequest) (res *dto.PrintLabelGetResponse, err error)
	GetDeliveryKoli(ctx context.Context, req *dto.PrintLabelGetRequest) (res dto.DeliveryKoliRes, err error)
	ReprintLabel(ctx context.Context, req *dto.RePrintLabelGetRequest) (res *dto.PrintLabelGetResponse, err error)
}

type PrintLabelService struct {
	opt                          opt.Options
	RepositoryPickingOrderAssign repository.IPickingOrderAssignRepository
	RepositoryPickingOrder       repository.IPickingOrderRepository
	RepositoryDeliveryKoli       repository.IDeliveryKoliRepository
}

func NewPrintLabelService() IPrintLabelService {
	return &PrintLabelService{
		opt:                          global.Setup.Common,
		RepositoryPickingOrderAssign: repository.NewPickingOrderAssignRepository(),
		RepositoryPickingOrder:       repository.NewPickingOrderRepository(),
		RepositoryDeliveryKoli:       repository.NewDeliveryKoliRepository(),
	}
}

func (s *PrintLabelService) Get(ctx context.Context, req *dto.PrintLabelGetRequest) (res *dto.PrintLabelGetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PrintLabelService.Get")
	defer span.End()

	switch req.TypePrint {
	case "label_picking":
		var pickingOrderAssign *model.PickingOrderAssign
		if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.Condition); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}

		// get picking order assign's sales order information
		var salesOrder *bridgeService.GetSalesOrderGPListResponse
		if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
			Id: pickingOrderAssign.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
			return
		}

		// TODO uncomment wrt
		// get wrt
		var wrt *bridgeService.GetWrtGPResponse
		if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
			Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
			GnlRegion: salesOrder.Data[0].GnL_Region,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
			return
		}

		// count delivery koli
		var (
			deliveryKoli []*model.DeliveryKoli
			totalKoli    float64
		)
		if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
			SopNumber: pickingOrderAssign.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("delivery koli")
			return
		}
		for _, v2 := range deliveryKoli {
			totalKoli += v2.Quantity
		}

		// helper code
		var pickingOrder *model.PickingOrder
		if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}

		request := dto.LabelPickingRequest{
			SalesOrder: struct {
				Code   string "json:\"code\""
				Branch struct {
					Name string "json:\"name\""
				} "json:\"branch\""
				Wrt struct {
					Name string "json:\"name\""
				} "json:\"wrt\""
				OrderType struct {
					Value string "json:\"value\""
				} "json:\"order_type\""
			}{
				Code: pickingOrderAssign.SopNumber,
				Branch: struct {
					Name string "json:\"name\""
				}{
					Name: salesOrder.Data[0].Customer[0].Custname,
				},
				Wrt: struct {
					Name string "json:\"name\""
				}{
					// TODO uncomment
					// Name: "DUMMY WRT 04.00-05.00",
					Name: wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
				},
				OrderType: struct {
					Value string "json:\"value\""
				}{
					Value: salesOrder.Data[0].SoptypE_STRING,
				},
			},
			TotalKoli: int64(totalKoli),
			Helper: struct {
				Code string "json:\"code\""
			}{
				Code: pickingOrder.PickerId,
			},
		}

		// send to service print
		req := make(map[string]interface{})
		req["pls"] = request

		url := s.SendPrint(req, "read/picking_print")

		res = &dto.PrintLabelGetResponse{
			Data: url,
		}

		return
	default:
		var pickingOrderAssign *model.PickingOrderAssign
		if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.Condition); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order assign")
			return
		}

		// get picking order assign's sales order information
		var salesOrder *bridgeService.GetSalesOrderGPListResponse
		if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
			Id: pickingOrderAssign.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
			return
		}

		// TODO uncomment wrt
		// // get wrt
		var wrt *bridgeService.GetWrtGPResponse
		if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
			Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
			GnlRegion: salesOrder.Data[0].GnL_Region,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
			return
		}

		// count delivery koli
		var (
			deliveryKoli []*model.DeliveryKoli
			totalKoli    float64
		)
		if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
			SopNumber: pickingOrderAssign.SopNumber,
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("delivery koli")
			return
		}
		for _, v2 := range deliveryKoli {
			totalKoli += v2.Quantity
		}

		// helper code
		var pickingOrder *model.PickingOrder
		if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorInvalid("picking order")
			return
		}

		request := dto.LabelPickingRequest{
			SalesOrder: struct {
				Code   string "json:\"code\""
				Branch struct {
					Name string "json:\"name\""
				} "json:\"branch\""
				Wrt struct {
					Name string "json:\"name\""
				} "json:\"wrt\""
				OrderType struct {
					Value string "json:\"value\""
				} "json:\"order_type\""
			}{
				Code: pickingOrderAssign.SopNumber,
				Branch: struct {
					Name string "json:\"name\""
				}{
					Name: "DUMMY CUSTOMER NAME",
				},
				Wrt: struct {
					Name string "json:\"name\""
				}{
					// TODO uncomment
					// Name: "DUMMY WRT 04.00-05.00",
					Name: wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
				},
				OrderType: struct {
					Value string "json:\"value\""
				}{
					Value: salesOrder.Data[0].SoptypE_STRING,
				},
			},
			TotalKoli: int64(totalKoli),
			Helper: struct {
				Code string "json:\"code\""
			}{
				Code: pickingOrder.PickerId,
			},
		}

		// send to service print
		req := make(map[string]interface{})
		req["pls"] = request

		url := s.SendPrint(req, "read/picking_print")

		res = &dto.PrintLabelGetResponse{
			Data: url,
		}

		return
	}
}

func (s *PrintLabelService) GetDO(ctx context.Context, req *dto.PrintLabelGetRequest) (res *dto.PrintLabelGetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PrintLabelService.GetDO")
	defer span.End()

	// // Get Sales Movement with SO ID
	// var soMovement *bridgeService.GetSalesMovementGPResponse
	// if soMovement, err = s.opt.Client.BridgeServiceGrpc.GetSalesMovementGP(ctx, &bridgeService.GetSalesMovementGPRequest{
	// 	SoNumber: req.Condition,
	// }); err != nil {
	// 	span.RecordError(err)
	// 	s.opt.Logger.AddMessage(log.ErrorLevel, err)
	// 	err = edenlabs.ErrorRpcNotFound("bridge", "delivery order")
	// 	return
	// }

	// vaidation for provide SI or DO while print
	var salesOrder *bridgeService.GetSalesOrderGPListResponse
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: req.Condition,
	}); err != nil || salesOrder.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// get customer
	var customer *bridgeService.GetCustomerGPResponse
	if customer, err = s.opt.Client.BridgeServiceGrpc.GetCustomerGPDetail(ctx, &bridgeService.GetCustomerGPDetailRequest{
		Id: salesOrder.Data[0].Customer[0].Custnmbr,
	}); err != nil || customer.Data == nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "customer")
		return
	}

	// Get Folder PDF AWS Bucket from ENV
	keyS3BucketName := fmt.Sprintf("%s", s.opt.Env.GetString("s3.bucket_name_pdf"))

	if customer.Data[0].Shipcomplete == 1 {
		// Get SI using SO ID
		var salesInvoice *bridgeService.GetSalesInvoiceGPListResponse
		if salesInvoice, err = s.opt.Client.BridgeServiceGrpc.GetSalesInvoiceGPList(ctx, &bridgeService.GetSalesInvoiceGPListRequest{
			SoNumber:  req.Condition,
			DeltaUser: "user_print01",
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "sales invoice")
			return
		}
		res = &dto.PrintLabelGetResponse{
			Data: salesInvoice.Data[0].DataAttachment.DocumentViewBody,
		}

		// Add upload to S3
		// Open the file to upload
		reportingServiceURL := salesInvoice.Data[0].DataAttachment.DocumentViewBody

		parsedURL, _ := url.Parse(reportingServiceURL)
		queryValues := parsedURL.Query()
		fileName := queryValues.Get("SOPNUMBE") + ".pdf"
		// Send an HTTP GET request to the Reporting Service URL
		response, error := http.Get(reportingServiceURL)
		if error != nil {
			fmt.Println("Error:", err)
			return
		}
		defer response.Body.Close()

		// Check if the response status code indicates success (e.g. 200 OK)
		if response.StatusCode != http.StatusOK {
			fmt.Printf("Request failed with status: %s\n", response.Status)
			return
		}

		dir, error := os.Getwd()
		if error != nil {
			fmt.Println("Error directory", error)
			return
		}

		// Create a new file to save the downloaded PDF
		file, error := os.Create(fileName)
		if error != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		// Copy the response content to the file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println("Error copying response:", err)
			return
		}

		fmt.Println("PDF downloaded successfully.")

		// file location to upload
		fileLocation := dir + "/" + fileName

		targetFile, error := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
		if error != nil {
			err = fmt.Errorf("failed to open file | %v", err)
			return
		}
		defer targetFile.Close()

		info, error := s.opt.S3x.UploadPublicFile(ctx, keyS3BucketName, fileName, fileLocation, "sales-invoice")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = fmt.Errorf("failed to upload file | %v", err)
			return
		}
		// remove file pdf
		os.Remove(fileLocation)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("upload_pdf", "Error uploading file to S3")
			return
		}

		res = &dto.PrintLabelGetResponse{
			Data: info,
		}

	} else {
		// Get DO using SO ID
		var deliveryOrder *bridgeService.GetDeliveryOrderGPListResponse
		if deliveryOrder, err = s.opt.Client.BridgeServiceGrpc.GetDeliveryOrderListGP(ctx, &bridgeService.GetDeliveryOrderGPListRequest{
			SopNumbe:  req.Condition,
			DeltaUser: "user_print01",
		}); err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("bridge", "delivery order")
			return
		}
		res = &dto.PrintLabelGetResponse{
			Data: deliveryOrder.Data[0].DataAttachment.DocumentViewBody,
		}

		// Add upload to S3
		// Open the file to upload
		reportingServiceURL := deliveryOrder.Data[0].DataAttachment.DocumentViewBody
		parsedURL, _ := url.Parse(reportingServiceURL)
		queryValues := parsedURL.Query()
		fileName := queryValues.Get("SOPNUMBE") + ".pdf"
		// Send an HTTP GET request to the Reporting Service URL
		response, error := http.Get(reportingServiceURL)
		if error != nil {
			fmt.Println("Error:", err)
			return
		}
		defer response.Body.Close()

		// Check if the response status code indicates success (e.g. 200 OK)
		if response.StatusCode != http.StatusOK {
			fmt.Printf("Request failed with status: %s\n", response.Status)
			return
		}

		dir, error := os.Getwd()
		if error != nil {
			fmt.Println("Error directory", error)
			return
		}

		// Create a new file to save the downloaded PDF
		file, error := os.Create(fileName)
		if error != nil {
			fmt.Println("Error creating file:", err)
			return
		}
		defer file.Close()

		// Copy the response content to the file
		_, err = io.Copy(file, response.Body)
		if err != nil {
			fmt.Println("Error copying response:", err)
			return
		}

		fmt.Println("PDF downloaded successfully.")

		// file location to upload
		fileLocation := dir + "/" + fileName

		targetFile, error := os.OpenFile(fileLocation, os.O_RDWR|os.O_CREATE, 0775)
		if error != nil {
			err = fmt.Errorf("failed to open file | %v", err)
			return
		}
		defer targetFile.Close()

		info, error := s.opt.S3x.UploadPublicFile(ctx, keyS3BucketName, fileName, fileLocation, "delivery-order")
		if err != nil {
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = fmt.Errorf("failed to upload file | %v", err)
			return
		}
		// remove file pdf
		os.Remove(fileLocation)
		if err != nil {
			span.RecordError(err)
			s.opt.Logger.AddMessage(log.ErrorLevel, err)
			err = edenlabs.ErrorRpcNotFound("upload_pdf", "Error uploading file to S3")
			return
		}

		res = &dto.PrintLabelGetResponse{
			Data: info,
		}

	}

	return

}

func (s *PrintLabelService) SendPrint(req map[string]interface{}, url string) string {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	var client = &http.Client{Transport: tr}

	jsonReq, _ := json.Marshal(req)

	request, _ := http.NewRequest("POST", s.opt.Config.PrintService.Url+url, bytes.NewBuffer(jsonReq))

	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	defer response.Body.Close() // MUST CLOSED THIS

	var bodyBytes []byte
	if response.Body != nil {
		bodyBytes, _ = ioutil.ReadAll(response.Body)
	}

	var res dto.PrintLabelGetResponse

	json.Unmarshal(bodyBytes, &res)
	response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

	_, err = io.Copy(ioutil.Discard, response.Body) // WE READ THE BODY
	if err != nil {
		return "read the body"
	}

	return res.Data
}

func (s *PrintLabelService) GetDeliveryKoli(ctx context.Context, req *dto.PrintLabelGetRequest) (res dto.DeliveryKoliRes, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PrintLabelService.GetDO")
	defer span.End()

	var (
		deliveryKoli []*model.DeliveryKoli
		totalKoli    int64
	)

	// if not exist it should be return error
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: req.Condition,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}

	for idx, v := range deliveryKoli {
		res.DeliveryKoli = append(res.DeliveryKoli, &dto.DeliveryKoli{
			Id:        v.Id,
			SopNumber: v.SopNumber,
			KoliId:    v.KoliId,
			Quantity:  v.Quantity,
			Increment: int64(idx) + 1,
		})
		totalKoli += 1
	}
	res.Total = totalKoli
	return
}

func (s *PrintLabelService) ReprintLabel(ctx context.Context, req *dto.RePrintLabelGetRequest) (res *dto.PrintLabelGetResponse, err error) {
	ctx, span := s.opt.Trace.Start(ctx, "PrintLabelService.Get")
	defer span.End()
	var pickingOrderAssign *model.PickingOrderAssign
	if pickingOrderAssign, err = s.RepositoryPickingOrderAssign.GetByID(ctx, 0, req.SalesOrderCode); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order assign")
		return
	}

	// get picking order assign's sales order information
	var salesOrder *bridgeService.GetSalesOrderGPListResponse
	if salesOrder, err = s.opt.Client.BridgeServiceGrpc.GetSalesOrderListGPByID(ctx, &bridgeService.GetSalesOrderGPListByIDRequest{
		Id: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "sales order")
		return
	}

	// TODO uncomment wrt
	// get wrt
	var wrt *bridgeService.GetWrtGPResponse
	if wrt, err = s.opt.Client.BridgeServiceGrpc.GetWrtGPList(ctx, &bridgeService.GetWrtGPListRequest{
		Search:    salesOrder.Data[0].Wrt[0].GnL_WRT_ID,
		GnlRegion: salesOrder.Data[0].GnL_Region,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorRpcNotFound("bridge", "wrt")
		return
	}

	// count delivery koli
	var (
		deliveryKoli []*model.DeliveryKoli
		// totalKoli    float64
	)
	if deliveryKoli, _, err = s.RepositoryDeliveryKoli.Get(ctx, &dto.DeliveryKoliGetRequest{
		SopNumber: pickingOrderAssign.SopNumber,
	}); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("delivery koli")
		return
	}
	// for _, v2 := range deliveryKoli {
	// 	_ = v2
	// 	totalKoli += 1
	// }

	// helper code
	var pickingOrder *model.PickingOrder
	if pickingOrder, err = s.RepositoryPickingOrder.GetDetail(ctx, pickingOrderAssign.PickingOrderId, ""); err != nil {
		span.RecordError(err)
		s.opt.Logger.AddMessage(log.ErrorLevel, err)
		err = edenlabs.ErrorInvalid("picking order")
		return
	}

	var requests []*dto.LabelPickingReprintRequest
	for _, v := range req.Increments {
		req := &dto.LabelPickingReprintRequest{
			SalesOrder: struct {
				Code   string "json:\"code\""
				Branch struct {
					Name string "json:\"name\""
				} "json:\"branch\""
				Wrt struct {
					Name string "json:\"name\""
				} "json:\"wrt\""
				OrderType struct {
					Value string "json:\"value\""
				} "json:\"order_type\""
			}{
				Code: pickingOrderAssign.SopNumber,
				Branch: struct {
					Name string "json:\"name\""
				}{
					Name: salesOrder.Data[0].Customer[0].Custname,
				},
				Wrt: struct {
					Name string "json:\"name\""
				}{
					// TODO uncomment
					// Name: "DUMMY WRT 04.00-05.00",
					Name: wrt.Data[0].Strttime + "-" + wrt.Data[0].Endtime,
				},
				OrderType: struct {
					Value string "json:\"value\""
				}{
					Value: salesOrder.Data[0].SoptypE_STRING,
				},
			},
			TotalKoli:  int64(len(deliveryKoli)),
			Increments: v,
			Helper: struct {
				Code string "json:\"code\""
			}{
				Code: pickingOrder.PickerId,
			},
		}
		requests = append(requests, req)
	}

	// send to service print
	reqPrint := make(map[string]interface{})
	reqPrint["plis"] = requests

	url := s.SendPrint(reqPrint, "read/label_reprint")

	res = &dto.PrintLabelGetResponse{
		Data: url,
	}

	return
}
