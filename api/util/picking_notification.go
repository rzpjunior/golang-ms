package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

type ModelPickingNotification struct {
	SendTo    string `json:"send_to"`
	Title     string `json:"title"`
	Message   string `json:"message"`
	Type      string `json:"type"`
	RefID     int64  `json:"ref_id"`
	StaffID   int64  `json:"staff_id"`
	ServerKey string `json:"server_key"`
}

func PostPickingModelNotification(r *ModelPickingNotification) error {
	var err error
	var client = &http.Client{}
	jsonReq, _ := json.Marshal(r)
	request, err := http.NewRequest("POST", PostPickingNotifURL, bytes.NewBuffer(jsonReq))
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
