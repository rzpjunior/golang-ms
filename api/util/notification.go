package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type ModelNotification struct {
	SendTo     string `json:"send_to"`
	Title      string `json:"title"`
	Message    string `json:"message"`
	Type       string `json:"type"` // isi nya untuk type 1 = sales_order
	RefID      int64  `json:"ref_id"`
	MerchantID int64  `json:"merchant_id"`
	StaffID    int64  `json:"staff_id"`
	ServerKey  string `json:"server_key"`
}

type MessageNotification struct {
	Title   string `orm:"column(title);null" json:"title"`
	Message string `orm:"column(message);null" json:"message"`
	Type    string `orm:"column(type)"`
}

type MessageNotificationCampaign struct {
	ID             string  `json:"id" valid:"required"`
	Code           string  `json:"code" valid:"required"`
	CampaignName   string  `json:"campaign_name" valid:"required"`
	Area           []int64 `json:"area" valid:"required"`
	Archetype      []int64 `json:"archetype" valid:"required"`
	RedirectTo     int8    `json:"redirect_to" valid:"required"`
	RedirectToName string  `json:"redirect_to_name"`
	RedirectValue  string  `json:"redirect_value" valid:"required"`
	Title          string  `json:"title" valid:"required"`
	Message        string  `json:"message" valid:"required"`
	ServerKey      string  `json:"server_key" valid:"required"`
}

func PostModelNotification(r *ModelNotification) error {
	var err error
	var client = &http.Client{}
	jsonReq, _ := json.Marshal(r)
	request, err := http.NewRequest("POST", PostNotifURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	_, err = io.Copy(ioutil.Discard, response.Body) // WE READ THE BODY
	if err != nil {
		return err
	}
	return err
}

func PostModelNotificationFieldSales(r *ModelNotification) error {
	var err error
	var client = &http.Client{}
	jsonReq, _ := json.Marshal(r)
	request, err := http.NewRequest("POST", PostFieldSalesNotifURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	_, err = io.Copy(ioutil.Discard, response.Body) // WE READ THE BODY
	if err != nil {
		return err
	}
	return err
}

func PostModelNotificationCampaign(r *MessageNotificationCampaign) error {
	var err error
	var client = &http.Client{}
	jsonReq, _ := json.Marshal(r)
	request, err := http.NewRequest("POST", PostCampaignNotifURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	_, err = io.Copy(ioutil.Discard, response.Body) // WE READ THE BODY
	if err != nil {
		return err
	}
	return err
}
