package edenlabs

import (
	"fmt"
	"net/http"

	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"git.edenfarm.id/edenlabs/edenlabs/validation"
	"github.com/labstack/echo/v4"
)

type FormatResponse struct {
	Code    int               `json:"code"`
	Status  string            `json:"status,omitempty"`
	Message interface{}       `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
	Total   int64             `json:"total,omitempty"`
	Page    int               `json:"page,omitempty"`
	PerPage int               `json:"per_page,omitempty"`
}

type SuccessResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status,omitempty"`
	Message interface{} `json:"message,omitempty"`
}

type ErrorResponse struct {
	Code    int               `json:"code"`
	Status  string            `json:"status,omitempty"`
	Message interface{}       `json:"message,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

const (
	// HTTPResponseSuccess default status for success responses
	HTTPResponseSuccess = "success"

	// HTTPResponseFailure default status for failure responses
	HTTPResponseFailure = "failure"
)

// NewResponse return new instances of response formater.
func NewResponse() *FormatResponse {
	return &FormatResponse{
		Code:   http.StatusOK,
		Status: HTTPResponseSuccess,
	}
}

// SetMessage fill message into response formater.
func (r *FormatResponse) SetMessage(message string) *FormatResponse {
	r.Status = HTTPResponseSuccess
	r.Code = http.StatusOK
	r.Message = message

	return r
}

// SetData fill data and total into response formater.
func (r *FormatResponse) SetData(d interface{}) *FormatResponse {
	r.Status = HTTPResponseSuccess
	r.Code = http.StatusOK
	r.Data = &d

	return r
}

// SetDataList fill data list into response formater.
func (r *FormatResponse) SetDataList(d interface{}, total int64, page *Paginator) *FormatResponse {
	r.Status = HTTPResponseSuccess
	r.Code = http.StatusOK
	r.Data = &d

	if total > 0 {
		r.Total = total
		if page.PerPage == 0 {
			r.PerPage = int(total)
			r.Page = 1
		} else {
			r.PerPage = page.PerPage
			r.Page = page.Page
		}
	}

	return r
}

// func (r *FormatResponse) ChangeTimezone(timezone string, object interface{}) interface{} {
// 	sv := reflect.ValueOf(object)
// 	fmt.Println("sv.Kind() ", sv.Kind())
// 	// if sv.Kind() == reflect.Ptr && !sv.IsNil() {
// 	// 	fmt.Println("Loop")
// 	// 	return r.ChangeTimezone(timezone, sv.Elem().Interface())
// 	// }

// 	// if sv.Kind() == reflect.Slice && sv.Len() > 0 {
// 	// 	fmt.Println("Sv Is Slice")
// 	// }

// 	// if sv.Kind() != reflect.Struct && sv.Kind() != reflect.Interface {
// 	// 	fmt.Println("Not Struct Not Interface")
// 	// 	return nil
// 	// }
// 	svi := reflect.ValueOf(sv.Elem().Interface())
// 	nf := svi.NumField()
// 	for i := 0; i < nf; i++ {
// 		f := svi.Field(i)
// 		if f.Type() == reflect.TypeOf(time.Time{}) {
// 			fmt.Println("F ", i, "Is Time")
// 			fmt.Println("F ", i, " Value ", f)
// 			locTime := timex.ToLocTime(f.Interface().(time.Time), timezone)
// 			fmt.Println("locTime ", locTime)
// 			fmt.Println("f.CanSet() ", reflect.ValueOf(object).Field(i).CanSet())
// 			// reflect.ValueOf(&object).Field(i).SetString(locTime.Format(time.RFC3339))
// 		}

// 		if f.Kind() == reflect.Ptr && !f.IsNil() {
// 			fmt.Println("F ", i, "Is Ptr")
// 		}

// 		for f.Kind() == reflect.Ptr && !f.IsNil() {
// 			fmt.Println("F ", i, "Change Elm")
// 			f = f.Elem()
// 		}

// 		if (f.Kind() == reflect.Struct || f.Kind() == reflect.Interface) && f.Type() != reflect.TypeOf(time.Time{}) {
// 			fmt.Println("F ", i, "Is Struct or Interface")
// 		}

// 		if f.Kind() == reflect.Slice && f.Len() > 0 {
// 			fmt.Println("F ", i, "Is Slice")
// 			if f.Index(0).Kind() == reflect.Struct || f.Index(0).Kind() == reflect.Ptr {
// 				for i := 0; i < f.Len(); i++ {
// 					fmt.Println("f.Index(i).Kind() : ", f.Index(i).Kind())
// 				}
// 			}
// 		}
// 	}

// 	return object
// }

// SetError set an error into response formater.
func (r *FormatResponse) SetError(err error) *FormatResponse {
	// Check error based on type
	if he, ok := err.(*echo.HTTPError); ok {
		// http error
		r.Code = he.Code
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = map[string]string{
			"error": fmt.Sprintf("%v", he.Message),
		}
	} else if o, ok := err.(*validation.Output); ok {
		// validation error
		r.Code = http.StatusUnprocessableEntity
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = o.Messages()

	} else if err == orm.ErrTxDone {
		r.Code = http.StatusUnprocessableEntity
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = map[string]string{
			"error": "Transaction already done",
		}
	} else if err == orm.ErrMultiRows {
		r.Code = http.StatusUnprocessableEntity
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = map[string]string{
			"error": "Return multi rows",
		}
	} else if err == orm.ErrNoRows {
		r.Code = http.StatusUnprocessableEntity
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = map[string]string{
			"error": "No row Found",
		}
	} else if err == orm.ErrStmtClosed {
		r.Code = http.StatusUnprocessableEntity
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = map[string]string{
			"error": "Stmt already closed",
		}
	} else if err == orm.ErrArgs {
		r.Code = http.StatusUnprocessableEntity
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = map[string]string{
			"error": "Args error may be empty",
		}
	} else if err == orm.ErrNotImplement {
		r.Code = http.StatusUnprocessableEntity
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = map[string]string{
			"error": "Have not implement",
		}
	} else if err == orm.ErrLastInsertIdUnavailable {
		r.Code = http.StatusUnprocessableEntity
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = map[string]string{
			"error": "Last insert id is unavailable",
		}
	} else {
		message := err.Error()
		// other error
		r.Code = http.StatusInternalServerError
		r.Status = HTTPResponseFailure
		r.Message = http.StatusText(r.Code)
		r.Errors = map[string]string{
			"error": message,
		}
	}
	return r
}

// Reset all data in response formater
func (r *FormatResponse) Reset() {
	r.Data = nil
	r.Errors = nil
	r.Message = nil
}
