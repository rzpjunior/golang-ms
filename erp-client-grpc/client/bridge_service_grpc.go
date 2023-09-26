package client

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/env"
	"git.edenfarm.id/project-version3/erp-pkg/erp-client-grpc/cirbreax"
	pb "git.edenfarm.id/project-version3/erp-services/erp-protobuf/gen/proto/bridge_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	bridgeServiceGrpcCommandName = "bridge.service.grpc"
)

type IBridgeServiceGrpc interface {
	// GP Integrated
	// Wrt
	GetWrtGPList(ctx context.Context, req *pb.GetWrtGPListRequest) (res *pb.GetWrtGPResponse, err error)
	GetWrtGPDetail(ctx context.Context, req *pb.GetWrtGPDetailRequest) (res *pb.GetWrtGPResponse, err error)
	// Item
	GetItemGPList(ctx context.Context, req *pb.GetItemGPListRequest) (res *pb.GetItemGPResponse, err error)
	GetItemMasterComplexGP(ctx context.Context, req *pb.GetItemMasterComplexGPListRequest) (res *pb.GetItemMasterComplexGPListResponse, err error)
	GetItemGPDetail(ctx context.Context, req *pb.GetItemGPDetailRequest) (res *pb.GetItemGPResponse, err error)
	// Site
	GetSiteGPList(ctx context.Context, req *pb.GetSiteGPListRequest) (res *pb.GetSiteGPResponse, err error)
	GetSiteGPDetail(ctx context.Context, req *pb.GetSiteGPDetailRequest) (res *pb.GetSiteGPResponse, err error)
	// Courier
	GetCourierGPList(ctx context.Context, req *pb.GetCourierGPListRequest) (res *pb.GetCourierGPResponse, err error)
	GetCourierGPDetail(ctx context.Context, req *pb.GetCourierGPDetailRequest) (res *pb.GetCourierGPResponse, err error)
	// Courier Vendor
	GetCourierVendorGPList(ctx context.Context, req *pb.GetCourierVendorGPListRequest) (res *pb.GetCourierVendorGPResponse, err error)
	GetCourierVendorGPDetail(ctx context.Context, req *pb.GetCourierVendorGPDetailRequest) (res *pb.GetCourierVendorGPResponse, err error)
	// Vehicle Profile
	GetVehicleProfileGPList(ctx context.Context, req *pb.GetVehicleProfileGPListRequest) (res *pb.GetVehicleProfileGPResponse, err error)
	GetVehicleProfileGPDetail(ctx context.Context, req *pb.GetVehicleProfileGPDetailRequest) (res *pb.GetVehicleProfileGPResponse, err error)
	// Adm Division
	GetAdmDivisionGPList(ctx context.Context, req *pb.GetAdmDivisionGPListRequest) (res *pb.GetAdmDivisionGPResponse, err error)
	GetAdmDivisionGPDetail(ctx context.Context, req *pb.GetAdmDivisionGPDetailRequest) (res *pb.GetAdmDivisionGPResponse, err error)
	// Adm Division Coverage
	GetAdmDivisionCoverageGPList(ctx context.Context, req *pb.GetAdmDivisionCoverageGPListRequest) (res *pb.GetAdmDivisionCoverageGPResponse, err error)
	GetAdmDivisionCoverageGPDetail(ctx context.Context, req *pb.GetAdmDivisionCoverageGPDetailRequest) (res *pb.GetAdmDivisionCoverageGPResponse, err error)
	// Payment Method
	GetPaymentMethodGPList(ctx context.Context, req *pb.GetPaymentMethodGPListRequest) (res *pb.GetPaymentMethodGPResponse, err error)
	GetPaymentMethodGPDetail(ctx context.Context, req *pb.GetPaymentMethodGPDetailRequest) (res *pb.GetPaymentMethodGPResponse, err error)
	// Payment Term
	GetPaymentTermGPList(ctx context.Context, req *pb.GetPaymentTermGPListRequest) (res *pb.GetPaymentTermGPResponse, err error)
	GetPaymentTermGPDetail(ctx context.Context, req *pb.GetPaymentTermGPDetailRequest) (res *pb.GetPaymentTermGPResponse, err error)
	// Sales Person
	GetSalesPersonGPList(ctx context.Context, req *pb.GetSalesPersonGPListRequest) (res *pb.GetSalesPersonGPResponse, err error)
	GetSalesPersonGPDetail(ctx context.Context, req *pb.GetSalesPersonGPDetailRequest) (res *pb.GetSalesPersonGPResponse, err error)
	// Customer Type
	GetCustomerTypeGPList(ctx context.Context, req *pb.GetCustomerTypeGPListRequest) (res *pb.GetCustomerTypeGPResponse, err error)
	GetCustomerTypeGPDetail(ctx context.Context, req *pb.GetCustomerTypeGPDetailRequest) (res *pb.GetCustomerTypeGPResponse, err error)
	// Archetype
	GetArchetypeGPList(ctx context.Context, req *pb.GetArchetypeGPListRequest) (res *pb.GetArchetypeGPResponse, err error)
	GetArchetypeGPDetail(ctx context.Context, req *pb.GetArchetypeGPDetailRequest) (res *pb.GetArchetypeGPResponse, err error)
	// Order Type
	GetOrderTypeGPList(ctx context.Context, req *pb.GetOrderTypeGPListRequest) (res *pb.GetOrderTypeGPResponse, err error)
	GetOrderTypeGPDetail(ctx context.Context, req *pb.GetOrderTypeGPDetailRequest) (res *pb.GetOrderTypeGPResponse, err error)
	// Vendor Classification
	GetVendorClassificationGPList(ctx context.Context, req *pb.GetVendorClassificationGPListRequest) (res *pb.GetVendorClassificationGPResponse, err error)
	GetVendorClassificationGPDetail(ctx context.Context, req *pb.GetVendorClassificationGPDetailRequest) (res *pb.GetVendorClassificationGPResponse, err error)
	// Uom
	GetUomGPList(ctx context.Context, req *pb.GetUomGPListRequest) (res *pb.GetUomGPResponse, err error)
	GetUomGPDetail(ctx context.Context, req *pb.GetUomGPDetailRequest) (res *pb.GetUomGPResponse, err error)
	// Vendor
	GetVendorGPList(ctx context.Context, req *pb.GetVendorGPListRequest) (res *pb.GetVendorGPResponse, err error)
	GetVendorGPDetail(ctx context.Context, req *pb.GetVendorGPDetailRequest) (res *pb.GetVendorGPResponse, err error)
	// Helper
	GetHelperGPList(ctx context.Context, req *pb.GetHelperGPListRequest) (res *pb.GetHelperGPResponse, err error)
	GetHelperGPDetail(ctx context.Context, req *pb.GetHelperGPDetailRequest) (res *pb.GetHelperGPResponse, err error)
	// SalesOrder
	GetSalesOrderListGPAll(ctx context.Context, req *pb.GetSalesOrderGPListRequest) (res *pb.GetSalesOrderGPListResponse, err error)
	GetSalesOrderListGPByID(ctx context.Context, req *pb.GetSalesOrderGPListByIDRequest) (res *pb.GetSalesOrderGPListResponse, err error)

	CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (res *pb.CreateAddressResponse, err error)
	UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (res *pb.UpdateAddressResponse, err error)
	SetDefaultAddress(ctx context.Context, req *pb.SetDefaultAddressRequest) (res *pb.SetDefaultAddressResponse, err error)

	// Item Class
	GetItemClassGPList(ctx context.Context, req *pb.GetItemClassGPListRequest) (res *pb.GetItemClassGPResponse, err error)
	GetItemClassGPDetail(ctx context.Context, req *pb.GetItemClassGPDetailRequest) (res *pb.GetItemClassGPResponse, err error)
	// Sales Territory
	GetSalesTerritoryGPList(ctx context.Context, req *pb.GetSalesTerritoryGPListRequest) (res *pb.GetSalesTerritoryGPResponse, err error)
	GetSalesTerritoryGPDetail(ctx context.Context, req *pb.GetSalesTerritoryGPDetailRequest) (res *pb.GetSalesTerritoryGPResponse, err error)
	// Customer
	GetCustomerGPList(ctx context.Context, req *pb.GetCustomerGPListRequest) (res *pb.GetCustomerGPResponse, err error)
	GetCustomerGPDetail(ctx context.Context, req *pb.GetCustomerGPDetailRequest) (res *pb.GetCustomerGPResponse, err error)
	// Address
	GetAddressGPList(ctx context.Context, req *pb.GetAddressGPListRequest) (res *pb.GetAddressGPResponse, err error)
	GetAddressGPDetail(ctx context.Context, req *pb.GetAddressGPDetailRequest) (res *pb.GetAddressGPResponse, err error)
	// Transaction List (Sales Order & Invoice)
	GetTransactionListGPList(ctx context.Context, req *pb.GetTransactionListGPListRequest) (res *pb.GetTransactionListGPResponse, err error)
	GetTransactionListGPDetail(ctx context.Context, req *pb.GetTransactionListGPDetailRequest) (res *pb.GetTransactionListGPResponse, err error)
	// Transaction Detail (Sales Order & Invoice Item)
	GetTransactionDetailGPList(ctx context.Context, req *pb.GetTransactionDetailGPListRequest) (res *pb.GetTransactionDetailGPResponse, err error)
	GetTransactionDetailGPDetail(ctx context.Context, req *pb.GetTransactionDetailGPDetailRequest) (res *pb.GetTransactionDetailGPResponse, err error)
	// Picking Order
	GetPickingOrderGPHeader(ctx context.Context, req *pb.GetPickingOrderGPHeaderRequest) (res *pb.GetPickingOrderGPHeaderResponse, err error)
	GetPickingOrderGPDetail(ctx context.Context, req *pb.GetPickingOrderGPDetailRequest) (res *pb.GetPickingOrderGPDetailResponse, err error)
	SubmitPickingCheckingPickingOrder(ctx context.Context, req *pb.SubmitPickingCheckingRequest) (res *pb.SubmitPickingCheckingResponse, err error)
	// login helper
	LoginHelper(ctx context.Context, req *pb.LoginHelperRequest) (res *pb.LoginHelperResponse, err error)

	UpdateFixedVa(ctx context.Context, req *pb.UpdateFixedVaRequest) (res *pb.UpdateFixedVaResponse, err error)
	GetAddressList(ctx context.Context, req *pb.GetAddressListRequest) (res *pb.GetAddressListResponse, err error)
	GetAddressListWithExcludedIds(ctx context.Context, req *pb.GetAddressListWithExcludedIdsRequest) (res *pb.GetAddressListResponse, err error)
	GetAddressDetail(ctx context.Context, req *pb.GetAddressDetailRequest) (res *pb.GetAddressDetailResponse, err error)
	DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (res *pb.DeleteAddressResponse, err error)
	GetAdmDivisionList(ctx context.Context, req *pb.GetAdmDivisionListRequest) (res *pb.GetAdmDivisionListResponse, err error)
	GetAdmDivisionDetail(ctx context.Context, req *pb.GetAdmDivisionDetailRequest) (res *pb.GetAdmDivisionDetailResponse, err error)
	GetArchetypeList(ctx context.Context, req *pb.GetArchetypeListRequest) (res *pb.GetArchetypeListResponse, err error)
	GetArchetypeDetail(ctx context.Context, req *pb.GetArchetypeDetailRequest) (res *pb.GetArchetypeDetailResponse, err error)
	GetCustomerTypeList(ctx context.Context, req *pb.GetCustomerTypeListRequest) (res *pb.GetCustomerTypeListResponse, err error)
	GetCustomerTypeDetail(ctx context.Context, req *pb.GetCustomerTypeDetailRequest) (res *pb.GetCustomerTypeDetailResponse, err error)
	GetClassList(ctx context.Context, req *pb.GetClassListRequest) (res *pb.GetClassListResponse, err error)
	GetClassDetail(ctx context.Context, req *pb.GetClassDetailRequest) (res *pb.GetClassDetailResponse, err error)
	GetItemList(ctx context.Context, req *pb.GetItemListRequest) (res *pb.GetItemListResponse, err error)
	GetItemDetail(ctx context.Context, req *pb.GetItemDetailRequest) (res *pb.GetItemDetailResponse, err error)
	GetRegionList(ctx context.Context, req *pb.GetRegionListRequest) (res *pb.GetRegionListResponse, err error)
	UpdateItemPackable(ctx context.Context, req *pb.UpdateItemPackableRequest) (res *pb.UpdateItemPackableResponse, err error)
	UpdateItemFragile(ctx context.Context, req *pb.UpdateItemFragileRequest) (res *pb.UpdateItemFragileResponse, err error)
	GetRegionDetail(ctx context.Context, req *pb.GetRegionDetailRequest) (res *pb.GetRegionDetailResponse, err error)
	GetSalespersonList(ctx context.Context, req *pb.GetSalespersonListRequest) (res *pb.GetSalespersonListResponse, err error)
	GetSalespersonDetail(ctx context.Context, req *pb.GetSalespersonDetailRequest) (res *pb.GetSalespersonDetailResponse, err error)
	GetSiteList(ctx context.Context, req *pb.GetSiteListRequest) (res *pb.GetSiteListResponse, err error)
	GetSiteInIdsList(ctx context.Context, req *pb.GetSiteInIdsListRequest) (res *pb.GetSiteListResponse, err error)
	GetSiteDetail(ctx context.Context, req *pb.GetSiteDetailRequest) (res *pb.GetSiteDetailResponse, err error)
	GetSubDistrictList(ctx context.Context, req *pb.GetSubDistrictListRequest) (res *pb.GetSubDistrictListResponse, err error)
	GetSubDistrictDetail(ctx context.Context, req *pb.GetSubDistrictDetailRequest) (res *pb.GetSubDistrictDetailResponse, err error)
	GetTerritoryList(ctx context.Context, req *pb.GetTerritoryListRequest) (res *pb.GetTerritoryListResponse, err error)
	GetTerritoryDetail(ctx context.Context, req *pb.GetTerritoryDetailRequest) (res *pb.GetTerritoryDetailResponse, err error)
	GetUomList(ctx context.Context, req *pb.GetUomListRequest) (res *pb.GetUomListResponse, err error)
	GetUomDetail(ctx context.Context, req *pb.GetUomDetailRequest) (res *pb.GetUomDetailResponse, err error)
	GetSalesOrderList(ctx context.Context, req *pb.GetSalesOrderListRequest) (res *pb.GetSalesOrderListResponse, err error)
	GetSalesOrderDetail(ctx context.Context, req *pb.GetSalesOrderDetailRequest) (res *pb.GetSalesOrderDetailResponse, err error)
	GetSalesOrderItemList(ctx context.Context, req *pb.GetSalesOrderItemListRequest) (res *pb.GetSalesOrderItemListResponse, err error)
	GetSalesOrderItemDetail(ctx context.Context, req *pb.GetSalesOrderItemDetailRequest) (res *pb.GetSalesOrderItemDetailResponse, err error)
	CreateSalesOrder(ctx context.Context, req *pb.CreateSalesOrderRequest) (res *pb.CreateSalesOrderResponse, err error)
	CreateSalesOrderGP(ctx context.Context, req *pb.CreateSalesOrderGPRequest) (res *pb.CreateSalesOrderGPResponse, err error)
	// Courier
	GetCourierList(ctx context.Context, req *pb.GetCourierListRequest) (res *pb.GetCourierListResponse, err error)
	GetCourierDetail(ctx context.Context, req *pb.GetCourierDetailRequest) (res *pb.GetCourierDetailResponse, err error)
	ActivateEmergencyCourier(ctx context.Context, req *pb.EmergencyCourierRequest) (res *pb.EmergencyCourierResponse, err error)
	DeactivateEmergencyCourier(ctx context.Context, req *pb.EmergencyCourierRequest) (res *pb.EmergencyCourierResponse, err error)
	// Courier Vendor
	GetCourierVendorList(ctx context.Context, req *pb.GetCourierVendorListRequest) (res *pb.GetCourierVendorListResponse, err error)
	GetCourierVendorDetail(ctx context.Context, req *pb.GetCourierVendorDetailRequest) (res *pb.GetCourierVendorDetailResponse, err error)
	// Vehicle Profile
	GetVehicleProfileList(ctx context.Context, req *pb.GetVehicleProfileListRequest) (res *pb.GetVehicleProfileListResponse, err error)
	GetVehicleProfileDetail(ctx context.Context, req *pb.GetVehicleProfileDetailRequest) (res *pb.GetVehicleProfileDetailResponse, err error)

	GetWrtList(ctx context.Context, req *pb.GetWrtListRequest) (res *pb.GetWrtListResponse, err error)
	GetCustomerList(ctx context.Context, req *pb.GetCustomerListRequest) (res *pb.GetCustomerListResponse, err error)
	GetCustomerDetail(ctx context.Context, req *pb.GetCustomerDetailRequest) (res *pb.GetCustomerDetailResponse, err error)
	UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (res *pb.UpdateCustomerResponse, err error)
	CreateCustomerGP(ctx context.Context, req *pb.CreateCustomerGPRequest) (res *pb.CreateCustomerGPResponse, err error)
	UpdateCustomerGP(ctx context.Context, req *pb.UpdateCustomerGPRequest) (res *pb.UpdateCustomerGPResponse, err error)

	GetWrtDetail(ctx context.Context, req *pb.GetWrtDetailRequest) (res *pb.GetWrtDetailResponse, err error)
	GetOrderTypeList(ctx context.Context, req *pb.GetOrderTypeListRequest) (res *pb.GetOrderTypeListResponse, err error)
	GetOrderTypeDetail(ctx context.Context, req *pb.GetOrderTypeDetailRequest) (res *pb.GetOrderTypeDetailResponse, err error)
	GetSalesPaymentTermList(ctx context.Context, req *pb.GetSalesPaymentTermListRequest) (res *pb.GetSalesPaymentTermListResponse, err error)
	GetSalesPaymentTermDetail(ctx context.Context, req *pb.GetSalesPaymentTermDetailRequest) (res *pb.GetSalesPaymentTermDetailResponse, err error)
	GetDivisionList(ctx context.Context, req *pb.GetDivisionListRequest) (res *pb.GetDivisionListResponse, err error)
	GetDivisionDetail(ctx context.Context, req *pb.GetDivisionDetailRequest) (res *pb.GetDivisionDetailResponse, err error)
	GetDistrictList(ctx context.Context, req *pb.GetDistrictListRequest) (res *pb.GetDistrictListResponse, err error)
	GetDistrictDetail(ctx context.Context, req *pb.GetDistrictDetailRequest) (res *pb.GetDistrictDetailResponse, err error)
	GetDistrictInIdsList(ctx context.Context, req *pb.GetDistrictInIdsListRequest) (res *pb.GetDistrictListResponse, err error)
	GetBankList(ctx context.Context, req *pb.GetBankListRequest) (res *pb.GetBankListResponse, err error)
	GetBankDetail(ctx context.Context, req *pb.GetBankDetailRequest) (res *pb.GetBankDetailResponse, err error)
	GetVendorList(ctx context.Context, req *pb.GetVendorListRequest) (res *pb.GetVendorListResponse, err error)
	GetVendorDetail(ctx context.Context, req *pb.GetVendorDetailRequest) (res *pb.GetVendorDetailResponse, err error)
	GetVendorOrganizationList(ctx context.Context, req *pb.GetVendorOrganizationListRequest) (res *pb.GetVendorOrganizationListResponse, err error)
	GetVendorOrganizationDetail(ctx context.Context, req *pb.GetVendorOrganizationDetailRequest) (res *pb.GetVendorOrganizationDetailResponse, err error)
	GetVendorOrganizationGPList(ctx context.Context, req *pb.GetVendorOrganizationGPListRequest) (res *pb.GetVendorOrganizationGPResponse, err error)
	GetVendorOrganizationGPDetail(ctx context.Context, req *pb.GetVendorOrganizationGPDetailRequest) (res *pb.GetVendorOrganizationGPResponse, err error)
	GetVendorClassificationList(ctx context.Context, req *pb.GetVendorClassificationListRequest) (res *pb.GetVendorClassificationListResponse, err error)
	GetVendorClassificationDetail(ctx context.Context, req *pb.GetVendorClassificationDetailRequest) (res *pb.GetVendorClassificationDetailResponse, err error)

	// purchase
	GetPurchasePlanList(ctx context.Context, req *pb.GetPurchasePlanListRequest) (res *pb.GetPurchasePlanListResponse, err error)
	GetPurchasePlanDetail(ctx context.Context, req *pb.GetPurchasePlanDetailRequest) (res *pb.GetPurchasePlanDetailResponse, err error)
	GetPurchasePlanItemList(ctx context.Context, req *pb.GetPurchasePlanItemListRequest) (res *pb.GetPurchasePlanItemListResponse, err error)
	GetPurchasePlanItemDetail(ctx context.Context, req *pb.GetPurchasePlanItemDetailRequest) (res *pb.GetPurchasePlanItemDetailResponse, err error)
	CreatePurchaseOrder(ctx context.Context, req *pb.CreatePurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error)
	CreatePurchaseOrderGP(ctx context.Context, req *pb.CreatePurchaseOrderGPRequest) (res *pb.CreatePurchaseOrderGPResponse, err error)
	GetPurchaseOrderList(ctx context.Context, req *pb.GetPurchaseOrderListRequest) (res *pb.GetPurchaseOrderListResponse, err error)
	GetPurchaseOrderDetail(ctx context.Context, req *pb.GetPurchaseOrderDetailRequest) (res *pb.GetPurchaseOrderDetailResponse, err error)
	GetPurchaseOrderItemList(ctx context.Context, req *pb.GetPurchaseOrderItemListRequest) (res *pb.GetPurchaseOrderItemListResponse, err error)
	GetPurchaseOrderItemDetail(ctx context.Context, req *pb.GetPurchaseOrderItemDetailRequest) (res *pb.GetPurchaseOrderItemDetailResponse, err error)
	CommitPurchaseOrder(ctx context.Context, req *pb.CommitPurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error)
	UpdatePurchaseOrder(ctx context.Context, req *pb.UpdatePurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error)
	UpdatePurchaseOrderGP(ctx context.Context, req *pb.UpdatePurchaseOrderGPRequest) (res *pb.CreatePurchaseOrderGPResponse, err error)
	UpdateProductPurchaseOrder(ctx context.Context, req *pb.UpdateProductPurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error)
	AssignPurchasePlanGP(ctx context.Context, req *pb.AssignPurchasePlanGPRequest) (res *pb.AssignPurchasePlanGPResponse, err error)
	CancelAssignPurchasePlan(ctx context.Context, req *pb.CancelAssignPurchasePlanRequest) (res *pb.CancelAssignPurchasePlanResponse, err error)
	GetPurchasePlanGPList(ctx context.Context, req *pb.GetPurchasePlanGPListRequest) (res *pb.GetPurchasePlanGPResponse, err error)
	GetPurchasePlanGPDetail(ctx context.Context, req *pb.GetPurchasePlanGPDetailRequest) (res *pb.GetPurchasePlanGPResponse, err error)
	GetPurchaseOrderGPList(ctx context.Context, req *pb.GetPurchaseOrderGPListRequest) (res *pb.GetPurchaseOrderGPResponse, err error)
	GetPurchaseOrderGPDetail(ctx context.Context, req *pb.GetPurchaseOrderGPDetailRequest) (res *pb.GetPurchaseOrderGPResponse, err error)
	CommitPurchaseOrderGP(ctx context.Context, req *pb.CommitPurchaseOrderGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)

	CancelPurchaseOrder(ctx context.Context, req *pb.CancelPurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error)
	CancelPurchaseOrderGP(ctx context.Context, req *pb.CancelPurchaseOrderGPRequest) (res *pb.CancelPurchaseOrderGPResponse, err error)
	GetDeliveryFeeList(ctx context.Context, req *pb.GetDeliveryFeeListRequest) (res *pb.GetDeliveryFeeListResponse, err error)
	GetDeliveryFeeDetail(ctx context.Context, req *pb.GetDeliveryFeeDetailRequest) (res *pb.GetDeliveryFeeDetailResponse, err error)
	GetDeliveryFeeGPList(ctx context.Context, req *pb.GetDeliveryFeeGPListRequest) (res *pb.GetDeliveryFeeGPListResponse, err error)

	// Consolidated Shipment
	CreateConsolidatedShipmentGP(ctx context.Context, req *pb.CreateConsolidatedShipmentGPRequest) (res *pb.CreateConsolidatedShipmentGPResponse, err error)
	UpdateConsolidatedShipmentGP(ctx context.Context, req *pb.UpdateConsolidatedShipmentGPRequest) (res *pb.UpdateConsolidatedShipmentGPResponse, err error)

	// Item transfer
	GetItemTransferList(ctx context.Context, req *pb.GetItemTransferListRequest) (res *pb.GetItemTransferListResponse, err error)
	GetItemTransferDetail(ctx context.Context, req *pb.GetItemTransferDetailRequest) (res *pb.GetItemTransferDetailResponse, err error)
	GetItemTransferItemDetail(ctx context.Context, req *pb.GetItemTransferItemDetailRequest) (res *pb.GetItemTransferItemDetailResponse, err error)
	CreateItemTransfer(ctx context.Context, req *pb.CreateItemTransferRequest) (res *pb.GetItemTransferDetailResponse, err error)
	UpdateItemTransfer(ctx context.Context, req *pb.UpdateItemTransferRequest) (res *pb.GetItemTransferDetailResponse, err error)
	CommitItemTransfer(ctx context.Context, req *pb.CommitItemTransferRequest) (res *pb.GetItemTransferDetailResponse, err error)
	GetInTransitTransferGPList(ctx context.Context, req *pb.GetInTransitTransferGPListRequest) (res *pb.GetInTransitTransferGPResponse, err error)
	GetInTransitTransferGPDetail(ctx context.Context, req *pb.GetInTransitTransferGPDetailRequest) (res *pb.GetInTransitTransferGPResponse, err error)
	GetTransferRequestGPList(ctx context.Context, req *pb.GetTransferRequestGPListRequest) (res *pb.GetTransferRequestGPResponse, err error)
	GetTransferRequestGPDetail(ctx context.Context, req *pb.GetTransferRequestGPDetailRequest) (res *pb.GetTransferRequestGPResponse, err error)
	CreateTransferRequestGP(ctx context.Context, req *pb.CreateTransferRequestGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)
	UpdateTransferRequestGP(ctx context.Context, req *pb.UpdateTransferRequestGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)
	UpdateInTransitTransferGP(ctx context.Context, req *pb.UpdateInTransitTransferGPRequest) (res *pb.UpdateInTransitTransferGPResponse, err error)
	CommitTransferRequestGP(ctx context.Context, req *pb.CommitTransferRequestGPRequest) (res *pb.CommitTransferRequestGPResponse, err error)

	// Receiving
	GetReceivingList(ctx context.Context, req *pb.GetReceivingListRequest) (res *pb.GetReceivingListResponse, err error)
	GetReceivingDetail(ctx context.Context, req *pb.GetReceivingDetailRequest) (res *pb.GetReceivingDetailResponse, err error)
	CreateReceiving(ctx context.Context, req *pb.CreateReceivingRequest) (res *pb.GetReceivingDetailResponse, err error)
	ConfirmReceiving(ctx context.Context, req *pb.ConfirmReceivingRequest) (res *pb.GetReceivingDetailResponse, err error)
	GetGoodsReceiptGPList(ctx context.Context, req *pb.GetGoodsReceiptGPListRequest) (res *pb.GetGoodsReceiptGPResponse, err error)
	GetGoodsReceiptGPDetail(ctx context.Context, req *pb.GetGoodsReceiptGPDetailRequest) (res *pb.GetGoodsReceiptGPResponse, err error)
	CreateGoodsReceiptGP(ctx context.Context, req *pb.CreateGoodsReceiptGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)
	UpdateGoodsReceiptGP(ctx context.Context, req *pb.UpdateGoodsReceiptGPRequest) (res *pb.CreateTransferRequestGPResponse, err error)
	// SalesInvoice & Payment
	GetSalesInvoiceList(ctx context.Context, req *pb.GetSalesInvoiceListRequest) (res *pb.GetSalesInvoiceListResponse, err error)
	GetSalesInvoiceGPList(ctx context.Context, req *pb.GetSalesInvoiceGPListRequest) (res *pb.GetSalesInvoiceGPListResponse, err error)
	GetSalesInvoiceGPDetail(ctx context.Context, req *pb.GetSalesInvoiceGPDetailRequest) (res *pb.GetSalesInvoiceGPDetailResponse, err error)
	GetSalesInvoiceDetail(ctx context.Context, req *pb.GetSalesInvoiceDetailRequest) (res *pb.GetSalesInvoiceDetailResponse, err error)
	CreateSalesInvoiceGP(ctx context.Context, req *pb.CreateSalesInvoiceGPRequest) (res *pb.CreateSalesInvoiceGPResponse, err error)
	GetSalesPaymentList(ctx context.Context, req *pb.GetSalesPaymentListRequest) (res *pb.GetSalesPaymentListResponse, err error)
	GetSalesPaymentDetail(ctx context.Context, req *pb.GetSalesPaymentDetailRequest) (res *pb.GetSalesPaymentDetailResponse, err error)
	GetSalesInvoiceItemList(ctx context.Context, req *pb.GetSalesInvoiceItemListRequest) (res *pb.GetSalesInvoiceItemListResponse, err error)
	GetDeliveryOrderDetail(ctx context.Context, req *pb.GetDeliveryOrderDetailRequest) (res *pb.GetDeliveryOrderDetailResponse, err error)
	CreateDeliveryOrder(ctx context.Context, req *pb.CreateDeliveryOrderRequest) (res *pb.CreateDeliveryOrderResponse, err error)
	CreateVendor(ctx context.Context, req *pb.CreateVendorRequest) (res *pb.CreateVendorResponse, err error)
	GetCashReceiptList(ctx context.Context, req *pb.GetCashReceiptListRequest) (res *pb.GetCashReceiptListResponse, err error)
	CreateCashReceipt(ctx context.Context, req *pb.CreateCashReceiptRequest) (res *pb.CreateCashReceiptResponse, err error)
	CreateSalesPaymentGP(ctx context.Context, req *pb.CreateSalesPaymentGPRequest) (res *pb.CreateSalesInvoiceGPResponse, err error)
	CreateSalesPaymentGPnonPBD(ctx context.Context, req *pb.CreateSalesPaymentGPnonPBDRequest) (res *pb.CreateSalesPaymentGPnonPBDResponse, err error)
	GetSalesPaymentGPList(ctx context.Context, req *pb.GetSalesPaymentGPListRequest) (res *pb.GetSalesPaymentGPResponse, err error)
	GetSalesPaymentGPDetail(ctx context.Context, req *pb.GetSalesPaymentGPDetailRequest) (res *pb.GetSalesPaymentGPResponse, err error)

	// Delivery Order GP
	GetDeliveryOrderListGP(ctx context.Context, req *pb.GetDeliveryOrderGPListRequest) (res *pb.GetDeliveryOrderGPListResponse, err error)

	// Customer Class
	GetCustomerClassList(ctx context.Context, req *pb.GetCustomerClassListRequest) (res *pb.GetCustomerClassResponse, err error)
	GetCustomerClassDetail(ctx context.Context, req *pb.GetCustomerClassDetailRequest) (res *pb.GetCustomerClassResponse, err error)

	// Sales Price Level
	GetSalesPriceLevelList(ctx context.Context, req *pb.GetSalesPriceLevelListRequest) (res *pb.GetSalesPriceLevelResponse, err error)
	GetSalesPriceLevelDetail(ctx context.Context, req *pb.GetSalesPriceLevelDetailRequest) (res *pb.GetSalesPriceLevelResponse, err error)

	// Shipping Method
	GetShippingMethodList(ctx context.Context, req *pb.GetShippingMethodListRequest) (res *pb.GetShippingMethodResponse, err error)
	GetShippingMethodDetail(ctx context.Context, req *pb.GetShippingMethodDetailRequest) (res *pb.GetShippingMethodResponse, err error)

	// Voucher
	CreateVoucherGP(ctx context.Context, req *pb.CreateVoucherGPRequest) (res *pb.CreateVoucherGPResponse, err error)
	GetVoucherGPList(ctx context.Context, req *pb.GetVoucherGPListRequest) (res *pb.GetVoucherGPResponse, err error)

	// Movement SO GP
	GetSalesMovementGP(ctx context.Context, req *pb.GetSalesMovementGPRequest) (res *pb.GetSalesMovementGPResponse, err error)
}

type BridgeServiceGrpcOption struct {
	Host                  string
	Port                  int
	Timeout               time.Duration
	MaxConcurrentRequests int
	ErrorPercentThreshold int
	Tls                   bool
	PemPath               string
	Secret                string
	Realtime              bool
}

type bridgeServiceGrpc struct {
	Option        BridgeServiceGrpcOption
	GrpcClient    pb.BridgeServiceClient
	HystrixClient *cirbreax.Client
}

func NewBridgeServiceGrpc(opt BridgeServiceGrpcOption) (iBridgeService IBridgeServiceGrpc, err error) {
	var opts []grpc.DialOption
	env, e := env.Env("env")
	if e != nil {
		return
	}
	serviceGrpcHTTPBackoffInterval := time.Duration(env.GetInt("client.serviceGrpcHTTPBackoffInterval")) * time.Millisecond
	serviceGrpcHTTPMaxJitterInterval := time.Duration(env.GetInt("client.serviceGrpcHTTPMaxJitterInterval")) * time.Millisecond
	serviceGrpcHTTPTimeout := time.Duration(env.GetInt("client.serviceGrpcHTTPTimeout")) * time.Millisecond
	serviceGrpcHTTPRetryCount := env.GetInt("client.serviceGrpcHTTPRetryCount")

	if opt.Tls {
		var pemServerCA []byte
		pemServerCA, err = ioutil.ReadFile(opt.PemPath)
		if err != nil {
			return
		}

		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(pemServerCA) {
			err = errors.New("failed to add server ca's certificate")
			return
		}

		// Create the credentials and return it
		config := &tls.Config{
			RootCAs: certPool,
		}

		tlsCredentials := credentials.NewTLS(config)

		opts = append(opts, grpc.WithTransportCredentials(tlsCredentials))
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}
	opts = append(opts, grpc.WithReturnConnectionError())
	opts = append(opts, grpc.FailOnNonTempDialError(true))

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(fmt.Sprintf("%s:%d", opt.Host, opt.Port),
		opts...,
	)
	if err != nil {
		return
	}

	backoff := cirbreax.NewConstantBackoff(serviceGrpcHTTPBackoffInterval, serviceGrpcHTTPMaxJitterInterval)
	retrier := cirbreax.NewRetrier(backoff)

	client := cirbreax.NewHttpClient(
		cirbreax.WithHTTPTimeout(serviceGrpcHTTPTimeout),
		cirbreax.WithCommandName(bridgeServiceGrpcCommandName),
		cirbreax.WithHystrixTimeout(opt.Timeout),
		cirbreax.WithMaxConcurrentRequests(opt.MaxConcurrentRequests),
		cirbreax.WithErrorPercentThreshold(opt.ErrorPercentThreshold),
		cirbreax.WithRetryCount(serviceGrpcHTTPRetryCount),
		cirbreax.WithRetrier(retrier),
	)

	gRPCClient := pb.NewBridgeServiceClient(conn)

	iBridgeService = bridgeServiceGrpc{
		Option:        opt,
		GrpcClient:    gRPCClient,
		HystrixClient: client,
	}
	return
}

func (o bridgeServiceGrpc) GetAddressList(ctx context.Context, req *pb.GetAddressListRequest) (res *pb.GetAddressListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAddressList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDeliveryOrderListGP(ctx context.Context, req *pb.GetDeliveryOrderGPListRequest) (res *pb.GetDeliveryOrderGPListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryOrderListGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesMovementGP(ctx context.Context, req *pb.GetSalesMovementGPRequest) (res *pb.GetSalesMovementGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesMovementGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAddressListWithExcludedIds(ctx context.Context, req *pb.GetAddressListWithExcludedIdsRequest) (res *pb.GetAddressListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAddressListWithExcludedIds(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAddressDetail(ctx context.Context, req *pb.GetAddressDetailRequest) (res *pb.GetAddressDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAddressDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) DeleteAddress(ctx context.Context, req *pb.DeleteAddressRequest) (res *pb.DeleteAddressResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.DeleteAddress(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAdmDivisionList(ctx context.Context, req *pb.GetAdmDivisionListRequest) (res *pb.GetAdmDivisionListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAdmDivisionList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAdmDivisionDetail(ctx context.Context, req *pb.GetAdmDivisionDetailRequest) (res *pb.GetAdmDivisionDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAdmDivisionDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetArchetypeList(ctx context.Context, req *pb.GetArchetypeListRequest) (res *pb.GetArchetypeListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetArchetypeList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetArchetypeDetail(ctx context.Context, req *pb.GetArchetypeDetailRequest) (res *pb.GetArchetypeDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetArchetypeDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerTypeList(ctx context.Context, req *pb.GetCustomerTypeListRequest) (res *pb.GetCustomerTypeListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerTypeList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerTypeDetail(ctx context.Context, req *pb.GetCustomerTypeDetailRequest) (res *pb.GetCustomerTypeDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerTypeDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetClassList(ctx context.Context, req *pb.GetClassListRequest) (res *pb.GetClassListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetClassList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetClassDetail(ctx context.Context, req *pb.GetClassDetailRequest) (res *pb.GetClassDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetClassDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemList(ctx context.Context, req *pb.GetItemListRequest) (res *pb.GetItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemDetail(ctx context.Context, req *pb.GetItemDetailRequest) (res *pb.GetItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetRegionList(ctx context.Context, req *pb.GetRegionListRequest) (res *pb.GetRegionListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetRegionList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetRegionDetail(ctx context.Context, req *pb.GetRegionDetailRequest) (res *pb.GetRegionDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetRegionDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalespersonList(ctx context.Context, req *pb.GetSalespersonListRequest) (res *pb.GetSalespersonListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalespersonList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalespersonDetail(ctx context.Context, req *pb.GetSalespersonDetailRequest) (res *pb.GetSalespersonDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalespersonDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSiteList(ctx context.Context, req *pb.GetSiteListRequest) (res *pb.GetSiteListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSiteList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSiteInIdsList(ctx context.Context, req *pb.GetSiteInIdsListRequest) (res *pb.GetSiteListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSiteInIdsList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSiteDetail(ctx context.Context, req *pb.GetSiteDetailRequest) (res *pb.GetSiteDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSiteDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSubDistrictList(ctx context.Context, req *pb.GetSubDistrictListRequest) (res *pb.GetSubDistrictListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSubDistrictList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSubDistrictDetail(ctx context.Context, req *pb.GetSubDistrictDetailRequest) (res *pb.GetSubDistrictDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSubDistrictDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetTerritoryList(ctx context.Context, req *pb.GetTerritoryListRequest) (res *pb.GetTerritoryListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetTerritoryList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetTerritoryDetail(ctx context.Context, req *pb.GetTerritoryDetailRequest) (res *pb.GetTerritoryDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetTerritoryDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetUomList(ctx context.Context, req *pb.GetUomListRequest) (res *pb.GetUomListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUomList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetUomDetail(ctx context.Context, req *pb.GetUomDetailRequest) (res *pb.GetUomDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUomDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesOrderList(ctx context.Context, req *pb.GetSalesOrderListRequest) (res *pb.GetSalesOrderListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesOrderDetail(ctx context.Context, req *pb.GetSalesOrderDetailRequest) (res *pb.GetSalesOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesOrderItemList(ctx context.Context, req *pb.GetSalesOrderItemListRequest) (res *pb.GetSalesOrderItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderItemList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesOrderItemDetail(ctx context.Context, req *pb.GetSalesOrderItemDetailRequest) (res *pb.GetSalesOrderItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateSalesOrder(ctx context.Context, req *pb.CreateSalesOrderRequest) (res *pb.CreateSalesOrderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesOrder(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateSalesOrderGP(ctx context.Context, req *pb.CreateSalesOrderGPRequest) (res *pb.CreateSalesOrderGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesOrderGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCourierList(ctx context.Context, req *pb.GetCourierListRequest) (res *pb.GetCourierListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCourierList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCourierDetail(ctx context.Context, req *pb.GetCourierDetailRequest) (res *pb.GetCourierDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCourierDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCourierGPList(ctx context.Context, req *pb.GetCourierGPListRequest) (res *pb.GetCourierGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCourierGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCourierGPDetail(ctx context.Context, req *pb.GetCourierGPDetailRequest) (res *pb.GetCourierGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCourierGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCourierVendorList(ctx context.Context, req *pb.GetCourierVendorListRequest) (res *pb.GetCourierVendorListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCourierVendorList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCourierVendorDetail(ctx context.Context, req *pb.GetCourierVendorDetailRequest) (res *pb.GetCourierVendorDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCourierVendorDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCourierVendorGPList(ctx context.Context, req *pb.GetCourierVendorGPListRequest) (res *pb.GetCourierVendorGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCourierVendorGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCourierVendorGPDetail(ctx context.Context, req *pb.GetCourierVendorGPDetailRequest) (res *pb.GetCourierVendorGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCourierVendorGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVehicleProfileList(ctx context.Context, req *pb.GetVehicleProfileListRequest) (res *pb.GetVehicleProfileListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVehicleProfileList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVehicleProfileDetail(ctx context.Context, req *pb.GetVehicleProfileDetailRequest) (res *pb.GetVehicleProfileDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVehicleProfileDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVehicleProfileGPList(ctx context.Context, req *pb.GetVehicleProfileGPListRequest) (res *pb.GetVehicleProfileGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVehicleProfileGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVehicleProfileGPDetail(ctx context.Context, req *pb.GetVehicleProfileGPDetailRequest) (res *pb.GetVehicleProfileGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVehicleProfileGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetWrtList(ctx context.Context, req *pb.GetWrtListRequest) (res *pb.GetWrtListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetWrtList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerList(ctx context.Context, req *pb.GetCustomerListRequest) (res *pb.GetCustomerListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerDetail(ctx context.Context, req *pb.GetCustomerDetailRequest) (res *pb.GetCustomerDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateCustomer(ctx context.Context, req *pb.UpdateCustomerRequest) (res *pb.UpdateCustomerResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateCustomer(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateCustomerGP(ctx context.Context, req *pb.CreateCustomerGPRequest) (res *pb.CreateCustomerGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateCustomerGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetWrtDetail(ctx context.Context, req *pb.GetWrtDetailRequest) (res *pb.GetWrtDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetWrtDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateItemPackable(ctx context.Context, req *pb.UpdateItemPackableRequest) (res *pb.UpdateItemPackableResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateItemPackable(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateItemFragile(ctx context.Context, req *pb.UpdateItemFragileRequest) (res *pb.UpdateItemFragileResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateItemFragile(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetOrderTypeList(ctx context.Context, req *pb.GetOrderTypeListRequest) (res *pb.GetOrderTypeListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetOrderTypeList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetOrderTypeDetail(ctx context.Context, req *pb.GetOrderTypeDetailRequest) (res *pb.GetOrderTypeDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetOrderTypeDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPaymentTermList(ctx context.Context, req *pb.GetSalesPaymentTermListRequest) (res *pb.GetSalesPaymentTermListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPaymentTermList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPaymentTermDetail(ctx context.Context, req *pb.GetSalesPaymentTermDetailRequest) (res *pb.GetSalesPaymentTermDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPaymentTermDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) ActivateEmergencyCourier(ctx context.Context, req *pb.EmergencyCourierRequest) (res *pb.EmergencyCourierResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.ActivateEmergencyCourier(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) DeactivateEmergencyCourier(ctx context.Context, req *pb.EmergencyCourierRequest) (res *pb.EmergencyCourierResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.DeactivateEmergencyCourier(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDivisionList(ctx context.Context, req *pb.GetDivisionListRequest) (res *pb.GetDivisionListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDivisionList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDivisionDetail(ctx context.Context, req *pb.GetDivisionDetailRequest) (res *pb.GetDivisionDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDivisionDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDistrictList(ctx context.Context, req *pb.GetDistrictListRequest) (res *pb.GetDistrictListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDistrictList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDistrictDetail(ctx context.Context, req *pb.GetDistrictDetailRequest) (res *pb.GetDistrictDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDistrictDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDistrictInIdsList(ctx context.Context, req *pb.GetDistrictInIdsListRequest) (res *pb.GetDistrictListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDistrictInIdsList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetBankDetail(ctx context.Context, req *pb.GetBankDetailRequest) (res *pb.GetBankDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetBankDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetBankList(ctx context.Context, req *pb.GetBankListRequest) (res *pb.GetBankListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetBankList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorList(ctx context.Context, req *pb.GetVendorListRequest) (res *pb.GetVendorListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorDetail(ctx context.Context, req *pb.GetVendorDetailRequest) (res *pb.GetVendorDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorOrganizationList(ctx context.Context, req *pb.GetVendorOrganizationListRequest) (res *pb.GetVendorOrganizationListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorOrganizationList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorOrganizationDetail(ctx context.Context, req *pb.GetVendorOrganizationDetailRequest) (res *pb.GetVendorOrganizationDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorOrganizationDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorClassificationList(ctx context.Context, req *pb.GetVendorClassificationListRequest) (res *pb.GetVendorClassificationListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorClassificationList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorClassificationDetail(ctx context.Context, req *pb.GetVendorClassificationDetailRequest) (res *pb.GetVendorClassificationDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorClassificationDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchasePlanList(ctx context.Context, req *pb.GetPurchasePlanListRequest) (res *pb.GetPurchasePlanListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchasePlanList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchasePlanDetail(ctx context.Context, req *pb.GetPurchasePlanDetailRequest) (res *pb.GetPurchasePlanDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchasePlanDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchasePlanItemList(ctx context.Context, req *pb.GetPurchasePlanItemListRequest) (res *pb.GetPurchasePlanItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchasePlanItemList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchasePlanItemDetail(ctx context.Context, req *pb.GetPurchasePlanItemDetailRequest) (res *pb.GetPurchasePlanItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchasePlanItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchaseOrderList(ctx context.Context, req *pb.GetPurchaseOrderListRequest) (res *pb.GetPurchaseOrderListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchaseOrderList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchaseOrderDetail(ctx context.Context, req *pb.GetPurchaseOrderDetailRequest) (res *pb.GetPurchaseOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchaseOrderDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchaseOrderItemList(ctx context.Context, req *pb.GetPurchaseOrderItemListRequest) (res *pb.GetPurchaseOrderItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchaseOrderItemList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchaseOrderItemDetail(ctx context.Context, req *pb.GetPurchaseOrderItemDetailRequest) (res *pb.GetPurchaseOrderItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchaseOrderItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDeliveryFeeList(ctx context.Context, req *pb.GetDeliveryFeeListRequest) (res *pb.GetDeliveryFeeListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryFeeList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDeliveryFeeDetail(ctx context.Context, req *pb.GetDeliveryFeeDetailRequest) (res *pb.GetDeliveryFeeDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryFeeDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDeliveryFeeGPList(ctx context.Context, req *pb.GetDeliveryFeeGPListRequest) (res *pb.GetDeliveryFeeGPListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryFeeGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateConsolidatedShipmentGP(ctx context.Context, req *pb.CreateConsolidatedShipmentGPRequest) (res *pb.CreateConsolidatedShipmentGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateConsolidatedShipmentGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateConsolidatedShipmentGP(ctx context.Context, req *pb.UpdateConsolidatedShipmentGPRequest) (res *pb.UpdateConsolidatedShipmentGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateConsolidatedShipmentGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPickingOrderGPHeader(ctx context.Context, req *pb.GetPickingOrderGPHeaderRequest) (res *pb.GetPickingOrderGPHeaderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPickingOrderGPHeader(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPickingOrderGPDetail(ctx context.Context, req *pb.GetPickingOrderGPDetailRequest) (res *pb.GetPickingOrderGPDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPickingOrderGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) SubmitPickingCheckingPickingOrder(ctx context.Context, req *pb.SubmitPickingCheckingRequest) (res *pb.SubmitPickingCheckingResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SubmitPickingCheckingPickingOrder(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) LoginHelper(ctx context.Context, req *pb.LoginHelperRequest) (res *pb.LoginHelperResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.LoginHelper(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreatePurchaseOrder(ctx context.Context, req *pb.CreatePurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreatePurchaseOrder(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CommitPurchaseOrder(ctx context.Context, req *pb.CommitPurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CommitPurchaseOrder(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CancelPurchaseOrder(ctx context.Context, req *pb.CancelPurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CancelPurchaseOrder(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetWrtGPList(ctx context.Context, req *pb.GetWrtGPListRequest) (res *pb.GetWrtGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetWrtGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetWrtGPDetail(ctx context.Context, req *pb.GetWrtGPDetailRequest) (res *pb.GetWrtGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetWrtGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemGPList(ctx context.Context, req *pb.GetItemGPListRequest) (res *pb.GetItemGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemMasterComplexGP(ctx context.Context, req *pb.GetItemMasterComplexGPListRequest) (res *pb.GetItemMasterComplexGPListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemMasterComplexGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemGPDetail(ctx context.Context, req *pb.GetItemGPDetailRequest) (res *pb.GetItemGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSiteGPList(ctx context.Context, req *pb.GetSiteGPListRequest) (res *pb.GetSiteGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSiteGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSiteGPDetail(ctx context.Context, req *pb.GetSiteGPDetailRequest) (res *pb.GetSiteGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSiteGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdatePurchaseOrder(ctx context.Context, req *pb.UpdatePurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdatePurchaseOrder(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateProductPurchaseOrder(ctx context.Context, req *pb.UpdateProductPurchaseOrderRequest) (res *pb.GetPurchaseOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateProductPurchaseOrder(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemTransferList(ctx context.Context, req *pb.GetItemTransferListRequest) (res *pb.GetItemTransferListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemTransferList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemTransferDetail(ctx context.Context, req *pb.GetItemTransferDetailRequest) (res *pb.GetItemTransferDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemTransferDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemTransferItemDetail(ctx context.Context, req *pb.GetItemTransferItemDetailRequest) (res *pb.GetItemTransferItemDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemTransferItemDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesInvoiceDetail(ctx context.Context, req *pb.GetSalesInvoiceDetailRequest) (res *pb.GetSalesInvoiceDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesInvoiceDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPaymentList(ctx context.Context, req *pb.GetSalesPaymentListRequest) (res *pb.GetSalesPaymentListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPaymentList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPaymentDetail(ctx context.Context, req *pb.GetSalesPaymentDetailRequest) (res *pb.GetSalesPaymentDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPaymentDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesInvoiceItemList(ctx context.Context, req *pb.GetSalesInvoiceItemListRequest) (res *pb.GetSalesInvoiceItemListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesInvoiceItemList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetDeliveryOrderDetail(ctx context.Context, req *pb.GetDeliveryOrderDetailRequest) (res *pb.GetDeliveryOrderDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetDeliveryOrderDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateDeliveryOrder(ctx context.Context, req *pb.CreateDeliveryOrderRequest) (res *pb.CreateDeliveryOrderResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateDeliveryOrder(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAdmDivisionGPDetail(ctx context.Context, req *pb.GetAdmDivisionGPDetailRequest) (res *pb.GetAdmDivisionGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAdmDivisionGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAdmDivisionGPList(ctx context.Context, req *pb.GetAdmDivisionGPListRequest) (res *pb.GetAdmDivisionGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAdmDivisionGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAdmDivisionCoverageGPDetail(ctx context.Context, req *pb.GetAdmDivisionCoverageGPDetailRequest) (res *pb.GetAdmDivisionCoverageGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAdmDivisionCoverageGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAdmDivisionCoverageGPList(ctx context.Context, req *pb.GetAdmDivisionCoverageGPListRequest) (res *pb.GetAdmDivisionCoverageGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAdmDivisionCoverageGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPaymentMethodGPDetail(ctx context.Context, req *pb.GetPaymentMethodGPDetailRequest) (res *pb.GetPaymentMethodGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPaymentMethodGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPaymentMethodGPList(ctx context.Context, req *pb.GetPaymentMethodGPListRequest) (res *pb.GetPaymentMethodGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPaymentMethodGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPaymentTermGPDetail(ctx context.Context, req *pb.GetPaymentTermGPDetailRequest) (res *pb.GetPaymentTermGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPaymentTermGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPaymentTermGPList(ctx context.Context, req *pb.GetPaymentTermGPListRequest) (res *pb.GetPaymentTermGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPaymentTermGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPersonGPDetail(ctx context.Context, req *pb.GetSalesPersonGPDetailRequest) (res *pb.GetSalesPersonGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPersonGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPersonGPList(ctx context.Context, req *pb.GetSalesPersonGPListRequest) (res *pb.GetSalesPersonGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPersonGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerTypeGPDetail(ctx context.Context, req *pb.GetCustomerTypeGPDetailRequest) (res *pb.GetCustomerTypeGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerTypeGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerTypeGPList(ctx context.Context, req *pb.GetCustomerTypeGPListRequest) (res *pb.GetCustomerTypeGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerTypeGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetArchetypeGPDetail(ctx context.Context, req *pb.GetArchetypeGPDetailRequest) (res *pb.GetArchetypeGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetArchetypeGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetArchetypeGPList(ctx context.Context, req *pb.GetArchetypeGPListRequest) (res *pb.GetArchetypeGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetArchetypeGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetOrderTypeGPDetail(ctx context.Context, req *pb.GetOrderTypeGPDetailRequest) (res *pb.GetOrderTypeGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetOrderTypeGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetOrderTypeGPList(ctx context.Context, req *pb.GetOrderTypeGPListRequest) (res *pb.GetOrderTypeGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetOrderTypeGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorClassificationGPDetail(ctx context.Context, req *pb.GetVendorClassificationGPDetailRequest) (res *pb.GetVendorClassificationGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorClassificationGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorClassificationGPList(ctx context.Context, req *pb.GetVendorClassificationGPListRequest) (res *pb.GetVendorClassificationGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorClassificationGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorOrganizationGPDetail(ctx context.Context, req *pb.GetVendorOrganizationGPDetailRequest) (res *pb.GetVendorOrganizationGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorOrganizationGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorOrganizationGPList(ctx context.Context, req *pb.GetVendorOrganizationGPListRequest) (res *pb.GetVendorOrganizationGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorOrganizationGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetUomGPDetail(ctx context.Context, req *pb.GetUomGPDetailRequest) (res *pb.GetUomGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUomGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetUomGPList(ctx context.Context, req *pb.GetUomGPListRequest) (res *pb.GetUomGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetUomGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemClassGPDetail(ctx context.Context, req *pb.GetItemClassGPDetailRequest) (res *pb.GetItemClassGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemClassGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetItemClassGPList(ctx context.Context, req *pb.GetItemClassGPListRequest) (res *pb.GetItemClassGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetItemClassGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesTerritoryGPDetail(ctx context.Context, req *pb.GetSalesTerritoryGPDetailRequest) (res *pb.GetSalesTerritoryGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesTerritoryGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesTerritoryGPList(ctx context.Context, req *pb.GetSalesTerritoryGPListRequest) (res *pb.GetSalesTerritoryGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesTerritoryGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetTransactionListGPDetail(ctx context.Context, req *pb.GetTransactionListGPDetailRequest) (res *pb.GetTransactionListGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetTransactionListGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetTransactionListGPList(ctx context.Context, req *pb.GetTransactionListGPListRequest) (res *pb.GetTransactionListGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetTransactionListGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetTransactionDetailGPDetail(ctx context.Context, req *pb.GetTransactionDetailGPDetailRequest) (res *pb.GetTransactionDetailGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetTransactionDetailGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetTransactionDetailGPList(ctx context.Context, req *pb.GetTransactionDetailGPListRequest) (res *pb.GetTransactionDetailGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetTransactionDetailGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAddressGPDetail(ctx context.Context, req *pb.GetAddressGPDetailRequest) (res *pb.GetAddressGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAddressGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetAddressGPList(ctx context.Context, req *pb.GetAddressGPListRequest) (res *pb.GetAddressGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetAddressGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateItemTransfer(ctx context.Context, req *pb.CreateItemTransferRequest) (res *pb.GetItemTransferDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateItemTransfer(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateItemTransfer(ctx context.Context, req *pb.UpdateItemTransferRequest) (res *pb.GetItemTransferDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateItemTransfer(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) AssignPurchasePlanGP(ctx context.Context, req *pb.AssignPurchasePlanGPRequest) (res *pb.AssignPurchasePlanGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.AssignPurchasePlanGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CommitItemTransfer(ctx context.Context, req *pb.CommitItemTransferRequest) (res *pb.GetItemTransferDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CommitItemTransfer(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CancelAssignPurchasePlan(ctx context.Context, req *pb.CancelAssignPurchasePlanRequest) (res *pb.CancelAssignPurchasePlanResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CancelAssignPurchasePlan(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetReceivingList(ctx context.Context, req *pb.GetReceivingListRequest) (res *pb.GetReceivingListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetReceivingList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateVendor(ctx context.Context, req *pb.CreateVendorRequest) (res *pb.CreateVendorResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateVendor(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetReceivingDetail(ctx context.Context, req *pb.GetReceivingDetailRequest) (res *pb.GetReceivingDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetReceivingDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateReceiving(ctx context.Context, req *pb.CreateReceivingRequest) (res *pb.GetReceivingDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateReceiving(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateAddress(ctx context.Context, req *pb.CreateAddressRequest) (res *pb.CreateAddressResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateAddress(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateAddress(ctx context.Context, req *pb.UpdateAddressRequest) (res *pb.UpdateAddressResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateAddress(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) SetDefaultAddress(ctx context.Context, req *pb.SetDefaultAddressRequest) (res *pb.SetDefaultAddressResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.SetDefaultAddress(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) ConfirmReceiving(ctx context.Context, req *pb.ConfirmReceivingRequest) (res *pb.GetReceivingDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.ConfirmReceiving(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorGPList(ctx context.Context, req *pb.GetVendorGPListRequest) (res *pb.GetVendorGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVendorGPDetail(ctx context.Context, req *pb.GetVendorGPDetailRequest) (res *pb.GetVendorGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVendorGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetHelperGPList(ctx context.Context, req *pb.GetHelperGPListRequest) (res *pb.GetHelperGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetHelperGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetHelperGPDetail(ctx context.Context, req *pb.GetHelperGPDetailRequest) (res *pb.GetHelperGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetHelperGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesInvoiceList(ctx context.Context, req *pb.GetSalesInvoiceListRequest) (res *pb.GetSalesInvoiceListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesInvoiceList(context.TODO(), req)
		return
	})
	return
}
func (o bridgeServiceGrpc) GetSalesInvoiceGPList(ctx context.Context, req *pb.GetSalesInvoiceGPListRequest) (res *pb.GetSalesInvoiceGPListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesInvoiceGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesInvoiceGPDetail(ctx context.Context, req *pb.GetSalesInvoiceGPDetailRequest) (res *pb.GetSalesInvoiceGPDetailResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesInvoiceGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerGPList(ctx context.Context, req *pb.GetCustomerGPListRequest) (res *pb.GetCustomerGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerGPDetail(ctx context.Context, req *pb.GetCustomerGPDetailRequest) (res *pb.GetCustomerGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesOrderListGPAll(ctx context.Context, req *pb.GetSalesOrderGPListRequest) (res *pb.GetSalesOrderGPListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderListGPAll(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesOrderListGPByID(ctx context.Context, req *pb.GetSalesOrderGPListByIDRequest) (res *pb.GetSalesOrderGPListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesOrderListGPByID(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchasePlanGPList(ctx context.Context, req *pb.GetPurchasePlanGPListRequest) (res *pb.GetPurchasePlanGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchasePlanGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchasePlanGPDetail(ctx context.Context, req *pb.GetPurchasePlanGPDetailRequest) (res *pb.GetPurchasePlanGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchasePlanGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchaseOrderGPList(ctx context.Context, req *pb.GetPurchaseOrderGPListRequest) (res *pb.GetPurchaseOrderGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchaseOrderGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetPurchaseOrderGPDetail(ctx context.Context, req *pb.GetPurchaseOrderGPDetailRequest) (res *pb.GetPurchaseOrderGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetPurchaseOrderGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreatePurchaseOrderGP(ctx context.Context, req *pb.CreatePurchaseOrderGPRequest) (res *pb.CreatePurchaseOrderGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreatePurchaseOrderGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdatePurchaseOrderGP(ctx context.Context, req *pb.UpdatePurchaseOrderGPRequest) (res *pb.CreatePurchaseOrderGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdatePurchaseOrderGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetInTransitTransferGPList(ctx context.Context, req *pb.GetInTransitTransferGPListRequest) (res *pb.GetInTransitTransferGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetInTransitTransferGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetInTransitTransferGPDetail(ctx context.Context, req *pb.GetInTransitTransferGPDetailRequest) (res *pb.GetInTransitTransferGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetInTransitTransferGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetTransferRequestGPList(ctx context.Context, req *pb.GetTransferRequestGPListRequest) (res *pb.GetTransferRequestGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetTransferRequestGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetTransferRequestGPDetail(ctx context.Context, req *pb.GetTransferRequestGPDetailRequest) (res *pb.GetTransferRequestGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetTransferRequestGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCashReceiptList(ctx context.Context, req *pb.GetCashReceiptListRequest) (res *pb.GetCashReceiptListResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCashReceiptList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateCashReceipt(ctx context.Context, req *pb.CreateCashReceiptRequest) (res *pb.CreateCashReceiptResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateCashReceipt(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerClassList(ctx context.Context, req *pb.GetCustomerClassListRequest) (res *pb.GetCustomerClassResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerClassList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetCustomerClassDetail(ctx context.Context, req *pb.GetCustomerClassDetailRequest) (res *pb.GetCustomerClassResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetCustomerClassDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPriceLevelList(ctx context.Context, req *pb.GetSalesPriceLevelListRequest) (res *pb.GetSalesPriceLevelResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPriceLevelList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPriceLevelDetail(ctx context.Context, req *pb.GetSalesPriceLevelDetailRequest) (res *pb.GetSalesPriceLevelResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPriceLevelDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetShippingMethodList(ctx context.Context, req *pb.GetShippingMethodListRequest) (res *pb.GetShippingMethodResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetShippingMethodList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetShippingMethodDetail(ctx context.Context, req *pb.GetShippingMethodDetailRequest) (res *pb.GetShippingMethodResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetShippingMethodDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateTransferRequestGP(ctx context.Context, req *pb.CreateTransferRequestGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateTransferRequestGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CommitPurchaseOrderGP(ctx context.Context, req *pb.CommitPurchaseOrderGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CommitPurchaseOrderGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetGoodsReceiptGPList(ctx context.Context, req *pb.GetGoodsReceiptGPListRequest) (res *pb.GetGoodsReceiptGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetGoodsReceiptGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetGoodsReceiptGPDetail(ctx context.Context, req *pb.GetGoodsReceiptGPDetailRequest) (res *pb.GetGoodsReceiptGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetGoodsReceiptGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateTransferRequestGP(ctx context.Context, req *pb.UpdateTransferRequestGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateTransferRequestGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateInTransitTransferGP(ctx context.Context, req *pb.UpdateInTransitTransferGPRequest) (res *pb.UpdateInTransitTransferGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateInTransitTransferGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CommitTransferRequestGP(ctx context.Context, req *pb.CommitTransferRequestGPRequest) (res *pb.CommitTransferRequestGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CommitTransferRequestGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateCustomerGP(ctx context.Context, req *pb.UpdateCustomerGPRequest) (res *pb.UpdateCustomerGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateCustomerGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateSalesInvoiceGP(ctx context.Context, req *pb.CreateSalesInvoiceGPRequest) (res *pb.CreateSalesInvoiceGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesInvoiceGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateSalesPaymentGP(ctx context.Context, req *pb.CreateSalesPaymentGPRequest) (res *pb.CreateSalesInvoiceGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesPaymentGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateSalesPaymentGPnonPBD(ctx context.Context, req *pb.CreateSalesPaymentGPnonPBDRequest) (res *pb.CreateSalesPaymentGPnonPBDResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateSalesPaymentGPnonPBD(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateGoodsReceiptGP(ctx context.Context, req *pb.CreateGoodsReceiptGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateGoodsReceiptGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateGoodsReceiptGP(ctx context.Context, req *pb.UpdateGoodsReceiptGPRequest) (res *pb.CreateTransferRequestGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateGoodsReceiptGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPaymentGPList(ctx context.Context, req *pb.GetSalesPaymentGPListRequest) (res *pb.GetSalesPaymentGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPaymentGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetSalesPaymentGPDetail(ctx context.Context, req *pb.GetSalesPaymentGPDetailRequest) (res *pb.GetSalesPaymentGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetSalesPaymentGPDetail(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CancelPurchaseOrderGP(ctx context.Context, req *pb.CancelPurchaseOrderGPRequest) (res *pb.CancelPurchaseOrderGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CancelPurchaseOrderGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) CreateVoucherGP(ctx context.Context, req *pb.CreateVoucherGPRequest) (res *pb.CreateVoucherGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.CreateVoucherGP(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) GetVoucherGPList(ctx context.Context, req *pb.GetVoucherGPListRequest) (res *pb.GetVoucherGPResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.GetVoucherGPList(context.TODO(), req)
		return
	})
	return
}

func (o bridgeServiceGrpc) UpdateFixedVa(ctx context.Context, req *pb.UpdateFixedVaRequest) (res *pb.UpdateFixedVaResponse, err error) {
	err = o.HystrixClient.Execute(func() (err error) {
		res, err = o.GrpcClient.UpdateFixedVa(context.TODO(), req)
		return
	})
	return
}