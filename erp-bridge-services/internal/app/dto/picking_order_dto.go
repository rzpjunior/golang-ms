package dto

type PickingDetails struct {
	Sopnumbe     string  `json:"sopnumbe"`
	Lnitmseq     int32   `json:"lnitmseq"`
	IvmQtyPickso float64 `json:"ivm_qty_pickso"`
}

type Picking struct {
	Docnumbr string            `json:"docnumbr"`
	Strttime string            `json:"strttime"`
	Endtime  string            `json:"endtime"`
	Details  []*PickingDetails `json:"details"`
}

type CheckingDetails struct {
	Sopnumbe     string  `json:"sopnumbe"`
	Lnitmseq     int32   `json:"lnitmseq"`
	IvmQtyPickso float64 `json:"ivm_qty_pickso"`
}

type Checking struct {
	Docnumbr     string             `json:"docnumbr"`
	Sopnumbe     string             `json:"sopnumbe"`
	Strttime     string             `json:"strttime"`
	Endtime      string             `json:"endtime"`
	WmsPickerId  string             `json:"wms_picker_id"`
	IvmKoli      int32              `json:"ivm_koli"`
	IvmJenisKoli string             `json:"ivm_jenis_koli"`
	Details      []*CheckingDetails `json:"details"`
}

type SubmitPickingCheckingRequest struct {
	Interid  string      `json:"interid"`
	Uniqueid string      `json:"uniqueid"`
	Bachnumb string      `json:"bachnumb"`
	Picking  *Picking    `json:"picking"`
	Checking []*Checking `json:"checking"`
}
