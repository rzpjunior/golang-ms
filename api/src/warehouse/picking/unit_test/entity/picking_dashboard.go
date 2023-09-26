package entity

type PickingDashboard struct {
	WrtId                     string
	NewPickingStatus          int
	OnProgressPickingStatus   int
	NeedApprovalPickingStatus int
	PickedPickingStatus       int
	CheckingPickingStatus     int
	FinishedPickingStatus     int
	TotalSO                   int
}
