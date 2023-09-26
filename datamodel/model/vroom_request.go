// Copyright 2020 PT. Eden Pangan Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package model

type VroomRequest struct {
	ID        int64       `json:"id"`
	Code      string      `json:"code"`
	Vehicles  []*Vehicle  `json:"vehicles"`
	Jobs      []*Job      `json:"jobs,omitempty"`
	Shipments []*Shipment `json:"shipments,omitempty"`
	Options   *Option     `json:"options"`
	Matrices  interface{} `json:"matrices"`
}

type Vehicle struct {
	ID          int64     `json:"id"`
	Profile     string    `json:"profile"`
	Start       []float64 `json:"start"`
	StartIndex  *int64    `json:"start_index,omitempty"`
	End         []float64 `json:"end"`
	EndIndex    *int64    `json:"end_index,omitempty"`
	Capacity    []int64   `json:"capacity"`
	Skills      []int64   `json:"skills,omitempty"`
	SpeedFactor float64   `json:"speed_factor"`
	MaxTasks    int64     `json:"max_tasks,omitempty"`
}
type Job struct {
	SalesOrderID  int64     `json:"id"`
	Description   string    `json:"description"` //Description filled by Sales Order Code
	Location      []float64 `json:"location"`
	LocationIndex *int64    `json:"location_index,omitempty"`
	Setup         int64     `json:"setup"`
	Service       int64     `json:"service"`
	Delivery      []int64   `json:"delivery"` //Delivery filled with weight and koli, index zero is the total weight, index one is the total koli
	Skills        []int64   `json:"skills,omitempty"`
	TimeWindows   [][]int64 `json:"time_windows,omitempty"`
}
type Shipment struct {
	Pickup   *PickupDelivery `json:"pickup"`   //Pickup filled with a shipment_step object describing pickup
	Delivery *PickupDelivery `json:"delivery"` //Delivery filled with a shipment_step object describing delivery
	Amount   []int64         `json:"amount"`
	Skills   []int64         `json:"skills,omitempty"`
	Priority int64           `json:"priority"`
}
type PickupDelivery struct {
	SalesOrderID  int64     `json:"id"`
	Description   string    `json:"description"` //Description filled by Sales Order Code
	Location      []float64 `json:"location"`
	LocationIndex *int64    `json:"location_index,omitempty"`
	Setup         int64     `json:"setup"`
	Service       int64     `json:"service"`
	TimeWindows   [][]int64 `json:"time_windows,omitempty"`
}

type Option struct {
	G bool `json:"g"`
	C bool `json:"c"`
}

type Matrice struct {
	Durations [][]int64 `json:"durations"`
	Costs     [][]int64 `json:"costs,omitempty"`
}

type RoutingMongoModel struct {
	ID            int64          `json:"id" bson:"id,omitempty"`
	Code          string         `json:"code" bson:"code,omitempty"`
	VroomResponse *VroomResponse `json:"vroom_response" bson:"vroom_response,omitempty"`
}

type VroomResponse struct {
	Code       int8                       `json:"code" bson:"code,omitempty"`
	Error      string                     `json:"error" bson:"error,omitempty"`
	Summary    *vroomResponseSummary      `json:"summary" bson:"summary,omitempty"`
	Routes     []*vroomResponseRoutes     `json:"routes" bson:"routes,omitempty"`
	Unassigned []*vroomResponseUnassigned `json:"unassigned" bson:"unassigned,omitempty"`
}

type vroomResponseSummary struct {
	Cost           int64                              `json:"cost" bson:"cost,omitempty"`
	Delivery       []int64                            `json:"delivery" bson:"delivery,omitempty"`
	Amount         []int64                            `json:"amount" bson:"amount,omitempty"`
	Pickup         []int64                            `json:"pickup" bson:"pickup,omitempty"`
	Setup          int64                              `json:"setup" bson:"setup,omitempty"`
	Service        int64                              `json:"service" bson:"service,omitempty"`
	Duration       int64                              `json:"duration" bson:"duration,omitempty"`
	WaitingTime    int64                              `json:"waiting_time" bson:"waiting_time,omitempty"`
	Priority       int64                              `json:"priority" bson:"priority,omitempty"`
	Distance       int64                              `json:"distance" bson:"distance,omitempty"`
	ComputingTimes vroomResponseSummaryComputingTimes `json:"computing_times" bson:"computing_times,omitempty"`
	Routes         int64                              `json:"routes" bson:"routes,omitempty"`
	Unassigned     int64                              `json:"unassigned" bson:"unassigned,omitempty"`
}

type vroomResponseRoutes struct {
	Vehicle     int64                       `json:"vehicle" bson:"vehicle,omitempty"`
	Cost        int64                       `json:"cost" bson:"cost,omitempty"`
	Delivery    []int64                     `json:"delivery" bson:"delivery,omitempty"`
	Amount      []int64                     `json:"amount" bson:"amount,omitempty"`
	Pickup      []int64                     `json:"pickup" bson:"pickup,omitempty"`
	Setup       int64                       `json:"setup" bson:"setup,omitempty"`
	Duration    int64                       `json:"duration" bson:"duration,omitempty"`
	Service     int64                       `json:"service" bson:"service,omitempty"`
	WaitingTime int64                       `json:"waiting_time" bson:"waiting_time,omitempty"`
	Priority    int64                       `json:"priority" bson:"priority,omitempty"`
	Distance    int64                       `json:"distance" bson:"distance,omitempty"`
	Steps       []*vroomResponseRoutesSteps `json:"steps" bson:"steps,omitempty"`
	Geometry    string                      `json:"geometry" bson:"geometry,omitempty"`
}

type vroomResponseRoutesSteps struct {
	ID          int64     `json:"id" bson:"id,omitempty"`
	Type        string    `json:"type" bson:"type,omitempty"`
	Description string    `json:"description" bson:"description,omitempty"`
	Location    []float64 `json:"location" bson:"location,omitempty"`
	Setup       int64     `json:"setup" bson:"setup,omitempty"`
	Service     int64     `json:"service" bson:"service,omitempty"`
	WaitingTime int64     `json:"waiting_time" bson:"waiting_time,omitempty"`
	Job         int64     `json:"job" bson:"job,omitempty"`
	Load        []int64   `json:"load" bson:"load,omitempty"`
	Arrival     int64     `json:"arrival" bson:"arrival,omitempty"`
	Duration    int64     `json:"duration" bson:"duration,omitempty"`
	Distance    int64     `json:"distance" bson:"distance,omitempty"`
}

type vroomResponseUnassigned struct {
	ID       int64     `json:"id" bson:"id,omitempty"`
	Location []float64 `json:"location" bson:"location,omitempty"`
	Type     string    `json:"type" bson:"type,omitempty"`
}

type vroomResponseSummaryComputingTimes struct {
	Loading int64 `json:"loading" bson:"loading,omitempty"`
	Solving int64 `json:"solving" bson:"solving,omitempty"`
	Routing int64 `json:"routing" bson:"routing,omitempty"`
}

type OSRMResponse struct {
	Code         string      `json:"code"`
	Durations    [][]float64 `json:"durations"`
	Destinations []struct {
		Hint     string    `json:"hint"`
		Distance float64   `json:"distance"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"destinations"`
	Sources []struct {
		Hint     string    `json:"hint"`
		Distance float64   `json:"distance"`
		Name     string    `json:"name"`
		Location []float64 `json:"location"`
	} `json:"sources"`
}
