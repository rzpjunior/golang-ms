package entity

type PickingAssign struct {
	Id               string
	Status           int8
	StatusPrint      int8
	PrintState       string
	Tolerable        float64
	LockedSalesOrder int8
	ProductItem      []Item
}

type Item struct {
	ProductCode string
	FlagOrder   int8
}

type ProductCode struct {
	Code string
}
