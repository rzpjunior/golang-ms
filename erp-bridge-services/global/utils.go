package global

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

// this function for unmarshal body response from https request
func unmarshalBody(r *http.Response, bodyResponse *interface{}) error {
	var err error
	if r != nil {

		err = json.NewDecoder(r.Body).Decode(&bodyResponse)
		if err != nil {
			log.Printf("error decoding response: %v", err)
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			log.Printf("response: %q", bodyResponse)
			return err
		}
	}

	return err
}

// this function for generate params before send using https request
func buildParams(params map[string]string, retry bool) string {
	var paramStr string
	var pageNumber, pageSize int

	// formula to calculate the page number based on limit and offset
	if params["PageSize"] != "" && params["PageNumber"] != "" {
		pageSize, _ = strconv.Atoi(params["PageSize"])
		params["PageSize"] = strconv.Itoa(pageSize)

		offset, _ := strconv.Atoi(params["PageNumber"])
		// if offset == 0 {
		// 	pageNumber = 1
		// } else {
		// 	pageNumber = pageSize * (offset - 1)
		// }
		if !retry {
			pageNumber = offset + 1
			params["PageNumber"] = strconv.Itoa(pageNumber)
		} else {
			params["PageNumber"] = strconv.Itoa(offset)
		}
	}

	// looping for make parameters in string url
	for key, value := range params {
		if len(paramStr) == 0 {
			paramStr += "?" + key + "=" + value
		} else {
			paramStr += "&" + key + "=" + value
		}
	}
	return paramStr
}

func Retry(fn func() error, maxAttempts int, sleep time.Duration) error {
	var err error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if err = fn(); err == nil {
			return nil // success
		}
		time.Sleep(sleep)
	}
	return err // failed after all attempts
}
