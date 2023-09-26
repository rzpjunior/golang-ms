package client

// Clients all client object injected here
type Clients struct {
	AccountServiceGrpc        IAccountServiceGrpc
	AuditServiceGrpc          IAuditServiceGrpc
	BridgeServiceGrpc         IBridgeServiceGrpc
	ConfigurationServiceGrpc  IConfigurationServiceGrpc
	CatalogServiceGrpc        ICatalogServiceGrpc
	CrmServiceGrpc            ICrmServiceGrpc
	LogisticServiceGrpc       ILogisticServiceGrpc
	CampaignServiceGrpc       ICampaignServiceGrpc
	SalesServiceGrpc          ISalesServiceGrpc
	PromotionServiceGrpc      IPromotionServiceGrpc
	SiteServiceGrpc           ISiteServiceGrpc
	SettlementGrpc            ISettlementServiceGrpc
	NotificationServiceGrpc   INotificationServiceGrpc
	CustomerMobileServiceGrpc ICustomerMobileServiceGrpc
	StorageServiceGrpc        IStorageServiceGrpc
}
