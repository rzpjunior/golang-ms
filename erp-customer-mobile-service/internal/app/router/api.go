package router

import (
	"git.edenfarm.id/project-version3/erp-services/erp-customer-mobile-service/internal/app/handler"
)

func init() {
	handlers["auth"] = &handler.AuthHandler{}
	handlers["config"] = &handler.ConfigHandler{}
	handlers["region"] = &handler.RegionPolicyHandler{}
	handlers["delivery"] = &handler.DeliveryDateHandler{}
	handlers["adm/division"] = &handler.AdmDivisionHandler{}
	handlers["customer/archetype"] = &handler.ArchetypeHandler{}
	handlers["term_condition"] = &handler.TermConditionHandler{}
	//handlers["wrt"] = &handler.WrtHandler{}
	handlers["item_category"] = &handler.ItemCategoryHandler{}
	handlers["sms_viro"] = &handler.SMSViroHandler{}
	handlers["campaign/banner"] = &handler.BannerHandler{}
	handlers["campaign/item-section"] = &handler.ItemSectionHandler{}
	handlers["address"] = &handler.AddressHandler{}
	handlers["term_condition"] = &handler.TermConditionHandler{}
	handlers["item"] = &handler.ItemHandler{}
	handlers["sales/payment"] = &handler.PaymentHandler{}
	handlers["search_suggestion"] = &handler.SearchSuggestionHandler{}
	handlers["promotion/voucher"] = &handler.VoucherHandler{}
	handlers["transaction_history"] = &handler.TransactionHistoryHandler{}
	handlers["sales/order"] = &handler.SalesOrderHandler{}
	handlers["eden_point"] = &handler.EdenPointHandler{}
	handlers["campaign/membership"] = &handler.MembershipHandler{}
	handlers["notification/transaction"] = &handler.NotificationTransactionHandler{}
	handlers["notification/campaign"] = &handler.NotificationCampaignHandler{}
	handlers["customer"] = &handler.CustomerHandler{}
	handlers["gmaps"] = &handler.GmapsHandler{}
	handlers["personal"] = &handler.PersonalHandler{}
	handlers["upload"] = &handler.StorageHandler{}
}
