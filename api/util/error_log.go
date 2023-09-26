package util

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

type ErrorLog struct {
	ErrorCode    int    `orm:"column(error_code);" json:"error_code"`
	Name         string `orm:"column(name)" json:"name"`
	Email        string `orm:"column(email)" json:"email"`
	ErrorMessage string `orm:"column(error_message)" json:"error_message"`
	Function     string `orm:"column(function)" json:"function"`
	Platform     string `orm:"column(platform)" json:"platform"`
}

// function for post error log to service log
func PostToServiceErrorLog(r ErrorLog) error {
	tr := &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	}
	var client = &http.Client{Transport: tr}
	jsonReq, _ := json.Marshal(r)
	request, _ := http.NewRequest("POST", ErrorLogs, bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if resp != nil {
		defer resp.Body.Close() // MUST CLOSED THIS
	}
	if err != nil {
		return err
	}
	//=================================================
	// READ THE BODY EVEN THE DATA IS NOT IMPORTANT
	// THIS MUST TO DO, TO AVOID MEMORY LEAK WHEN REUSING HTTP
	// CONNECTION
	//=================================================
	_, err = io.Copy(ioutil.Discard, resp.Body) // WE READ THE BODY
	if err != nil {
		return err
	}
	return nil
}
