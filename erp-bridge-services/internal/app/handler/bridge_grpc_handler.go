package handler

import (
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/global"
	"git.edenfarm.id/project-version3/erp-services/erp-bridge-services/internal/app/service"
)

type BridgeGrpcHandler struct {
	Option                       global.HandlerOptions
	ServicesAddress              service.IAddressService
	ServicesAdmDivision          service.IAdmDivisionService
	ServicesArchetype            service.IArchetypeService
	ServicesCashReceipt          service.ICashReceiptService
	ServicesCustomerType         service.ICustomerTypeService
	ServicesClass                service.IClassService
	ServicesItem                 service.IItemService
	ServicesRegion               service.IRegionService
	ServicesSalesperson          service.ISalespersonService
	ServicesSite                 service.ISiteService
	ServicesSubDistrict          service.ISubDistrictService
	ServicesTerritory            service.ITerritoryService
	ServicesUom                  service.IUomService
	ServicesSalesOrder           service.ISalesOrderService
	ServicesSalesOrderItem       service.ISalesOrderItemService
	ServicesCourier              service.ICourierService
	ServicesCourierVendor        service.ICourierVendorService
	ServicesVehicleProfile       service.IVehicleProfileService
	ServicesWrt                  service.IWrtService
	ServicesCustomer             service.ICustomerService
	ServicesOrderType            service.IOrderTypeService
	ServicesSalesPaymentTerm     service.ISalesPaymentTermService
	ServicesDivision             service.IDivisionService
	ServicesDistrict             service.IDistrictService
	ServicesBank                 service.IBankService
	ServicesDeliveryFee          service.IDeliveryFeeService
	ServicesVendor               service.IVendorService
	ServicesHelper               service.IHelperService
	ServicesVendorOrganization   service.IVendorOrganizationService
	ServicesVendorClassification service.IVendorClassificationService
	ServicesPurchasePlan         service.IPurchasePlanService
	ServicesPurchasePlanItem     service.IPurchasePlanItemService
	ServicesPurchaseOrder        service.IPurchaseOrderService
	ServicesPurchaseOrderItem    service.IPurchaseOrderItemService
	ServiceAdmDivisionCoverage   service.IAdmDivisionCoverageService
	ServicePickingOrder          service.IPickingOrderService
	ServicePurchaseOrder         service.IPurchaseOrderService
	ServiceItemTransfer          service.IItemTransferService
	ServiceSalesInvoice          service.ISalesInvoiceService
	ServiceSalesInvoiceItem      service.ISalesInvoiceItemService
	ServiceSalesPayment          service.ISalesPaymentService
	ServiceDeliveryOrder         service.IDeliveryOrderService
	ServicesPaymentMethod        service.IPaymentMethodService
	ServicesItemClass            service.IItemClassService
	ServicesPaymentTerm          service.IPaymentTermService
	ServicesSalesPerson          service.ISalesPersonService
	ServicesSalesTerritory       service.ISalesTerritoryService
	ServicesTransactionList      service.ITransactionListService
	ServicesTransactionDetail    service.ITransactionDetailService
	ServicesItemTransfer         service.IItemTransferService
	ServicesItemTransferItem     service.IItemTransferItemService
	ServiceReceiving             service.IReceivingService
	ServicesCustomerClass        service.ICustomerClassService
	ServicesSalesPriceLevel      service.ISalesPriceLevelService
	ServicesShippingMethod       service.IShippingMethodService
	ServicesVoucher              service.IVoucherService
}
