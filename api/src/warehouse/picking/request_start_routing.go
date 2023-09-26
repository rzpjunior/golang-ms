// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package picking

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/cuxs/common"
	"git.edenfarm.id/cuxs/orm"
	"git.edenfarm.id/cuxs/validation"
	"git.edenfarm.id/project-version2/api/datastore/repository"
	"git.edenfarm.id/project-version2/api/src/auth"
	"git.edenfarm.id/project-version2/api/util"
	"git.edenfarm.id/project-version2/datamodel/model"
)

// startRoutingAssignment : struct to hold picking assign routing request data
type startRoutingAssignment struct {
	ID int64 `json:"-"`

	PickerID        []string  `json:"picker_id"`
	PickerCapacity  int64     `json:"picker_max_weight"`
	AssignTimeStamp time.Time `json:"-"`
	PickerString    string    `json:"-"`

	PickingOrderAssign []*model.PickingOrderAssign `json:"-"`
	PickingOrder       *model.PickingOrder         `json:"-"`
	PickingList        *model.PickingList          `json:"-"`
	PickingRoutingStep []*model.PickingRoutingStep `json:"-"`
	PickerArr          []int64                     `json:"-"`
	VroomRequest       *model.VroomRequest         `json:"-"`

	Warehouse *model.Warehouse  `json:"-"`
	Picker    *model.Staff      `json:"-"`
	Session   *auth.SessionData `json:"-"`
}

// Validate : function to validate picking assign routing request data
func (r *startRoutingAssignment) Validate() *validation.Output {
	o := &validation.Output{Valid: true}
	o1 := orm.NewOrm()
	o1.Using("read_only")
	var (
		err                  error
		filter, exclude      map[string]interface{}
		total                int64
		pickingOrderItem     []*model.PickingOrderItem
		arrVroomPicker       []*model.Vehicle
		arrShipment          []*model.Shipment
		task                 int64
		weight               float64
		capacityInsufficient bool

		locationCounter        int64
		locationIndexWarehouse int64
		locationIndexMap       map[string]int64
		locationsURLParameter  string
	)
	locationIndexMap = make(map[string]int64)
	locationsURLParameter = "/table/v1/driving/"

	// Picking list has to be requested for routed type
	r.PickingList = &model.PickingList{ID: r.ID}
	err = r.PickingList.Read("id")
	if err != nil {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
		return o
	}
	if r.PickingList.Status != 1 && r.PickingList.Status != 3 && r.PickingList.Status != 4 {
		o.Failure("picking_list_id.invalid", util.ErrorType("picking list", "new"))
		return o
	}

	// check if there's a routing on going for this picking list
	filter = map[string]interface{}{"picking_list_id": r.PickingList.ID, "status_step__in": []int{2, 3}}
	_, total, err = repository.CheckPickingRoutingStepData(filter, exclude)
	if err != nil {
		o.Failure("picking_routing_step.invalid", util.ErrorInvalidData("picking routing step"))
		return o
	}
	if total != 0 {
		o.Failure("picking_routing_step.invalid", util.ErrorRoutingOnGoing())
		return o
	}

	// warehouse bin info is required
	r.Warehouse = &model.Warehouse{ID: r.PickingList.Warehouse.ID}
	if err = r.Warehouse.Read("id"); err != nil {
		o.Failure("warehouse.invalid", util.ErrorInvalidData("warehouse"))
	}
	if r.Warehouse.BinInfo == nil {
		o.Failure("bin_info.invalid", util.ErrorInvalidData("bin info"))
		return o
	}
	if err = r.Warehouse.BinInfo.Read("ID"); err != nil {
		o.Failure("bin_info.invalid", util.ErrorInvalidData("bin info"))
		return o
	}

	locationIndexMap[longLatConcatenator(r.Warehouse.Longitude, r.Warehouse.Latitude)] = locationCounter
	locationsURLParameter += longLatConcatenator(r.Warehouse.Longitude, r.Warehouse.Latitude) + ";"
	locationCounter++

	filter = map[string]interface{}{"picking_list_id": r.ID, "status__in": []int64{1, 3, 8}}
	r.PickingOrderAssign, total, err = repository.CheckPickingOrderAssignData(filter, exclude)
	if err != nil {
		o.Failure("picking_list_id.invalid", util.ErrorInvalidData("picking list"))
	}
	if total == 0 {
		o.Failure("picking_order_assign.invalid", util.ErrorInvalidData("picking order assign"))
	}

	// change kg to gram
	if r.PickerCapacity <= 0 {
		o.Failure("picker_capacity_invalid", util.ErrorGreater("picker capacity", "0"))
	}
	r.PickerCapacity *= 1000

	for _, v1 := range r.PickingOrderAssign {
		if v1.Status == int8(2) {
			o.Failure("picking_order_assign.active", util.ErrorPickingStatus("new", "on progress"))
		}

		err = v1.SalesOrder.Read("ID")
		if err != nil {
			o.Failure("sales_order_id.invalid", util.ErrorInvalidData("sales order"))
		}

		if v1.SalesOrder.Status != 1 && v1.SalesOrder.Status != 3 && v1.SalesOrder.Status != 9 && v1.SalesOrder.Status != 12 {
			o.Failure("status.invalid", util.ErrorSalesOrderOnPicking())
		}

		if v1.Helper.ID != r.Session.Staff.ID {
			o.Failure("staff_id.invalid", util.ErrorPickingListStaff())
		}

		filter = map[string]interface{}{"picking_order_assign_id": v1.ID, "picking_flag__in": []int64{1, 3, 4}}
		pickingOrderItem, _, err = repository.CheckPickingOrderItemData(filter, exclude)
		if err != nil {
			o.Failure("picking_order_item.invalid", util.ErrorNotFound("picking order item"))
		}

		for _, v2 := range pickingOrderItem {
			salesOrderItem := &model.SalesOrderItem{SalesOrder: v1.SalesOrder, Product: v2.Product}
			if err = salesOrderItem.Read("salesorder", "product"); err != nil {
				weight = v2.OrderQuantity * 1000
			} else {
				weight = salesOrderItem.Weight * 1000
			}

			if int64(weight) > r.PickerCapacity {
				r.PickerCapacity = int64(weight)
				capacityInsufficient = true
			}

			product := &model.Product{ID: v2.Product.ID}
			err = product.Read("ID")
			if err != nil {
				o.Failure("product.not_found", util.ErrorNotFound("product"))
			}

			stock := &model.Stock{
				Product:   product,
				Warehouse: r.Warehouse,
			}
			if err = o1.Read(stock, "Product", "Warehouse"); err != nil {
				o.Failure("stock.invalid", util.ErrorInvalidData("stock"))
			} else {
				if stock.Bin == nil {
					o.Failure("bin.invalid", util.ErrorInvalidData("bin"))
				} else {
					if stock.Bin.ID == 0 {
						o.Failure("bin.invalid", util.ErrorInvalidData("bin"))
					} else {
						bin := &model.Bin{
							ID:        stock.Bin.ID,
							Warehouse: r.Warehouse,
							Product:   product,
						}
						if err = bin.Read("ID", "Warehouse", "Product"); err != nil {
							o.Failure("bin.invalid", util.ErrorInvalidData("bin"))
						} else {

							if locationIndexMap[longLatConcatenator(*bin.Longitude, *bin.Latitude)+bin.Code] == 0 {
								locationIndexMap[longLatConcatenator(*bin.Longitude, *bin.Latitude)+bin.Code] = locationCounter
								locationsURLParameter += longLatConcatenator(*bin.Longitude, *bin.Latitude) + ";"
								locationCounter++
							}
							locationIndex := locationIndexMap[longLatConcatenator(*bin.Longitude, *bin.Latitude)+bin.Code]

							arrShipment = append(arrShipment, &model.Shipment{
								Pickup: &model.PickupDelivery{
									SalesOrderID:  v2.ID,
									Description:   product.Name,
									Location:      []float64{*bin.Longitude, *bin.Latitude},
									LocationIndex: &locationIndex,
									Service:       bin.ServiceTime,
								},
								Delivery: &model.PickupDelivery{
									SalesOrderID:  v2.ID,
									Description:   product.Name,
									Location:      []float64{*r.Warehouse.BinInfo.Longitude, *r.Warehouse.BinInfo.Latitude},
									LocationIndex: &locationIndexWarehouse,
									Service:       0,
								},
								Amount: []int64{int64(weight)},
							})
							task++
						}
					}
				}
			}
		}
	}

	if len(arrShipment) == 0 {
		o.Failure("sales_order_item.not_found", util.ErrorNotFound("sales order"))
		return o
	}

	if len(r.PickerID) != 0 {
		for _, v := range r.PickerID {
			pickerID, err := common.Decrypt(v)
			if err != nil {
				o.Failure("picker_id.invalid", util.ErrorInvalidData("picker"))
			}

			if r.Picker, err = repository.GetStaff("id", pickerID); err != nil {
				o.Failure("picker_id.invalid", util.ErrorInvalidData("picker"))
			}

			if r.Picker.Status != 1 {
				o.Failure("picker_id.invalid", util.ErrorActiveInd("picker"))
			}

			r.PickerArr = append(r.PickerArr, pickerID)
			pickerString := strconv.Itoa(int(pickerID))
			r.PickerString += pickerString + ","

			arrVroomPicker = append(arrVroomPicker, &model.Vehicle{
				ID:          int64(pickerID),
				Profile:     "car",
				Start:       []float64{*r.Warehouse.BinInfo.Longitude, *r.Warehouse.BinInfo.Latitude},
				End:         []float64{*r.Warehouse.BinInfo.Longitude, *r.Warehouse.BinInfo.Latitude},
				StartIndex:  &locationIndexWarehouse,
				EndIndex:    &locationIndexWarehouse,
				Capacity:    []int64{r.PickerCapacity},
				SpeedFactor: 0.1,
			})

		}
		r.PickerString = strings.TrimSuffix(r.PickerString, ",")
	} else {
		o.Failure("pickers.invalid", util.ErrorInputRequired("pickers"))
	}

	// check if the pickers assigned have on going routing
	filter = map[string]interface{}{"staff_id__in": r.PickerArr, "status_step__in": []int{2, 3}}
	_, total, err = repository.CheckPickingRoutingStepData(filter, exclude)
	if err != nil {
		o.Failure("picking_routing_step.invalid", util.ErrorInvalidData("picking routing step"))
		return o
	}
	if total != 0 {
		o.Failure("picking_routing_step.invalid", util.ErrorStaffBusy())
		return o
	}

	var maxTasks float64
	maxTasks = float64(len(arrShipment)*2) / float64(len(arrVroomPicker))
	// plus two is to make sure that all vroom shipments can be fullfiled because of the division
	// example 1 SOI = 2 shipments and we have 2 pickers
	// 2 shipments / 2 = 1 shipment per picker that will cause the routing to be fail
	maxTasks = math.Ceil(maxTasks) + 1

	for _, v := range arrVroomPicker {
		v.MaxTasks = int64(maxTasks)
		if capacityInsufficient == true {
			v.Capacity = []int64{r.PickerCapacity}
		}
	}

	mapDurationsCar, err := osrmHandler(locationsURLParameter)
	if err != nil {
		o.Failure("osrm_failed", util.ErrorInvalidData("osrm"))
		return o
	}

	r.VroomRequest = &model.VroomRequest{
		Code:      r.PickingList.Code,
		Vehicles:  arrVroomPicker,
		Shipments: arrShipment,
		Options: &model.Option{
			G: false,
			C: false,
		},
		Matrices: mapDurationsCar,
	}

	return o
}

// Messages : function to return error validation messages
func (r *startRoutingAssignment) Messages() map[string]string {
	messages := map[string]string{}

	return messages
}

// function to combine longitude and latitude and change the data from float to string
func longLatConcatenator(long, lat float64) string {
	longitudeStr := fmt.Sprintf("%f", long)
	latitudeStr := fmt.Sprintf("%f", lat)
	longLat := longitudeStr + "," + latitudeStr

	return longLat
}

//  function to send location url parameter to OSRM and return the durations
func osrmHandler(locationsURLParameter string) (matricesMap map[string]*model.Matrice, err error) {
	matricesMap = make(map[string]*model.Matrice)
	var durationsCar [][]int64
	var osrmCarResponse model.OSRMResponse
	var client = http.Client{}
	osrmCarURL := util.OsrmCar

	locationsURLParameter = strings.TrimSuffix(locationsURLParameter, ";")

	request, err := http.NewRequest("GET", osrmCarURL+locationsURLParameter, nil)
	if err != nil {
		return
	}

	responseOSRMCar, err := client.Do(request)
	if err != nil {
		return
	}

	defer responseOSRMCar.Body.Close()

	if err = json.NewDecoder(responseOSRMCar.Body).Decode(&osrmCarResponse); err != nil {
		return
	}

	for _, v1 := range osrmCarResponse.Durations {
		var durationCar []int64
		for _, v2 := range v1 {
			durationCar = append(durationCar, int64(v2))
		}
		durationsCar = append(durationsCar, durationCar)
	}

	for a := 0; a < len(durationsCar); a++ {
		for i := 0; i < len(durationsCar[a]); i++ {
			if a != i {
				durationsCar[a][i]++
			}
		}
	}

	matricesMap["car"] = &model.Matrice{
		Durations: durationsCar,
	}

	return
}
