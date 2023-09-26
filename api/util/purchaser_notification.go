package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type ModelPurchaserNotification struct {
	SendTo    string `json:"send_to"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	RefID     int64  `json:"ref_id"`
	StaffID   int64  `json:"staff_id"`
	ServerKey string `json:"server_key"`
}

func PostPurchaserModelNotification(r *ModelPurchaserNotification) error {
	var err error
	var client = &http.Client{}
	jsonReq, _ := json.Marshal(r)
	request, err := http.NewRequest("POST", PostPurchaserNotifURL, bytes.NewBuffer(jsonReq))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()
	_, err = io.Copy(ioutil.Discard, response.Body)
	if err != nil {
		return err
	}
	return err
}
