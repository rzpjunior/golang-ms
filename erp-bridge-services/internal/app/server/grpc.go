package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/handler"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/service"

	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func StartGrpcServer() (err error) {
	url := fmt.Sprintf("%s:%d", global.Setup.Common.Config.Grpc.Host, global.Setup.Common.Config.Grpc.Port)
	listen, err := net.Listen("tcp", url)
	if err != nil {
		logrus.Errorf("[GRPC] Fail to start listen and server: %v", err)
		return
	}

	opts := []grpc_logrus.Option{
		grpc_logrus.WithDecider(func(method string, err error) bool {
			if global.Setup.Common.Config.App.Debug && err == nil {
				return false
			}
			return true
		}),
	}

	logger := logrus.NewEntry(global.Setup.Common.Logger.Logger())

	// init grpc server
	srv := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_logrus.UnaryServerInterceptor(logger, opts...),
			grpc_recovery.UnaryServerInterceptor(),
			otelgrpc.UnaryServerInterceptor(),
		)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_logrus.StreamServerInterceptor(logger, opts...),
			grpc_recovery.StreamServerInterceptor(),
			otelgrpc.StreamServerInterceptor(),
		)),
	)

	// setup handler
	bridgeGrpcHandler := &handler.BridgeGrpcHandler{
		Option:                       global.Setup,
		ServicesAddress:              service.NewAddressService(),
		ServicesAdmDivision:          service.NewAdmDivisionService(),
		ServicesArchetype:            service.NewArchetypeService(),
		ServicesCashReceipt:          service.NewCashReceiptService(),
		ServicesCustomerType:         service.NewCustomerTypeService(),
		ServicesClass:                service.NewClassService(),
		ServicesItem:                 service.NewItemService(),
		ServicesRegion:               service.NewRegionService(),
		ServicesSalesperson:          service.NewSalespersonService(),
		ServicesSite:                 service.NewSiteService(),
		ServicesSubDistrict:          service.NewSubDistrictService(),
		ServicesTerritory:            service.NewTerritoryService(),
		ServicesUom:                  service.NewUomService(),
		ServicesSalesOrder:           service.NewSalesOrderService(),
		ServicesSalesOrderItem:       service.NewSalesOrderItemService(),
		ServicesCourier:              service.NewCourierService(),
		ServicesCourierVendor:        service.NewCourierVendorService(),
		ServicesVehicleProfile:       service.NewVehicleProfileService(),
		ServicesWrt:                  service.NewWrtService(),
		ServicesCustomer:             service.NewCustomerService(),
		ServicesOrderType:            service.NewOrderTypeService(),
		ServicesSalesPaymentTerm:     service.NewSalesPaymentTermService(),
		ServicesDivision:             service.NewDivisionService(),
		ServicesDistrict:             service.NewDistrictService(),
		ServicesBank:                 service.NewBankService(),
		ServicesDeliveryFee:          service.NewDeliveryFeeService(),
		ServicesVendor:               service.NewVendorService(),
		ServicesHelper:               service.NewHelperService(),
		ServicesVendorOrganization:   service.NewVendorOrganizationService(),
		ServicesVendorClassification: service.NewVendorClassificationService(),
		ServicesPurchasePlan:         service.NewPurchasePlanService(),
		ServicesPurchasePlanItem:     service.NewPurchasePlanItemService(),
		ServicesPurchaseOrder:        service.NewPurchaseOrderService(),
		ServicesPurchaseOrderItem:    service.NewPurchaseOrderItemService(),
		ServiceAdmDivisionCoverage:   service.NewAdmDivisionCoverageService(),
		ServicePickingOrder:          service.NewPickingOrderService(),
		ServicePurchaseOrder:         service.NewPurchaseOrderService(),
		ServiceItemTransfer:          service.NewItemTransferService(),
		ServiceSalesInvoice:          service.NewSalesInvoiceService(),
		ServiceSalesInvoiceItem:      service.NewSalesInvoiceItemService(),
		ServiceSalesPayment:          service.NewSalesPaymentService(),
		ServiceDeliveryOrder:         service.NewDeliveryOrderService(),
		ServicesPaymentMethod:        service.NewPaymentMethodService(),
		ServicesItemClass:            service.NewItemClassService(),
		ServicesPaymentTerm:          service.NewPaymentTermService(),
		ServicesSalesPerson:          service.NewSalesPersonService(),
		ServicesSalesTerritory:       service.NewSalesTerritoryService(),
		ServicesItemTransfer:         service.NewItemTransferService(),
		ServicesItemTransferItem:     service.NewItemTransferItemService(),
		ServicesTransactionList:      service.NewTransactionListService(),
		ServicesTransactionDetail:    service.NewTransactionDetailService(),
		ServiceReceiving:             service.NewReceivingService(),
		ServicesCustomerClass:        service.NewCustomerClassService(),
		ServicesSalesPriceLevel:      service.NewSalesPriceLevelService(),
		ServicesShippingMethod:       service.NewShippingMethodService(),
		ServicesVoucher:              service.NewVoucherService(),
	}
	pb.RegisterBridgeServiceServer(srv, bridgeGrpcHandler)
	grpc_health_v1.RegisterHealthServer(srv, health.NewServer())

	reflection.Register(srv)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			logrus.Infoln("[GRPC] Server is shutting down")
			srv.GracefulStop()
			logrus.Info("[GRPC] Bye")
			return
		}
	}()

	logrus.Infof("[GRPC] GRPC serve at %s\n", url)
	srv.Serve(listen)

	return
}
