package statusx

const (
	Active          = "Active"
	Finished        = "Finished"
	Cancelled       = "Cancelled"
	Deleted         = "Deleted"
	Draft           = "Draft"
	New             = "New"
	Archived        = "Archived"
	Approved        = "Approved"
	Rejected        = "Rejected"
	Declined        = "Declined"
	Registered      = "Registered"
	Requested       = "Requested"
	Valid           = "Valid"
	Used            = "Used"
	Accepted        = "Accepted"
	Picked          = "Picked"
	Done            = "Done"
	FullyCompleted  = "Fully Completed"
	Sold            = "Sold"
	Checking        = "Checking"
	InProgress      = "In Progress"
	Processing      = "Processing"
	InCart          = "In Cart"
	Pending         = "Pending"
	Postponed       = "Postponed"
	NotActive       = "Not Active"
	Failed          = "Failed"
	Expired         = "Expired"
	Undeliverable   = "Undeliverable"
	Waste           = "Waste"
	CancelledIssued = "Cancelled Issued"
	CancelledRedeem = "Cancelled Redeem"
	Closed          = "Closed"
)

var StatusMap = map[int8]string{
	1:  Active,
	2:  Finished,
	3:  Cancelled,
	4:  Deleted,
	5:  Draft,
	6:  New,
	7:  Archived,
	8:  Approved,
	9:  Rejected,
	10: Declined,
	11: Registered,
	12: Requested,
	13: Valid,
	14: Used,
	15: Accepted,
	16: Picked,
	17: Done,
	18: FullyCompleted,
	19: Sold,
	20: Checking,
	21: InProgress,
	22: Processing,
	23: InCart,
	24: Pending,
	25: Postponed,
	26: NotActive,
	27: Failed,
	28: Expired,
	29: Undeliverable,
	30: Waste,
	31: CancelledIssued,
	32: CancelledRedeem,
	33: Closed,
}

func ConvertStatusValue(statusValue int8) (statusName string) {
	var exists bool
	statusName, exists = StatusMap[statusValue]
	if !exists {
		statusName = "Unknow"
		return
	}
	return
}

func ConvertStatusName(statusName string) (statusValue int8) {
	for key, value := range StatusMap {
		if value == statusName {
			statusValue = key
			return
		}
	}
	statusValue = 0
	return
}
