package edenlabs

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"git.edenfarm.id/edenlabs/edenlabs/timex"
	"github.com/labstack/echo/v4"
)

// Context is custom echo.Context
// has defined as middleware.
type Context struct {
	echo.Context
	ResponseFormat *FormatResponse
	ResponseData   interface{}
}

// NewContext new instances of context
func NewContext(c echo.Context) *Context {
	return &Context{c, NewResponse(), nil}
}

// GetParamID: get params id in path
func (c *Context) GetParamID() (id int64, err error) {
	paramValue := c.Param("id")
	var value int
	value, err = strconv.Atoi(paramValue)
	id = int64(value)
	return
}

func (c *Context) GetUriParamInt(param string) (value int) {
	paramValue := c.Param(param)
	value, _ = strconv.Atoi(paramValue)
	return
}

func (c *Context) GetUriParamString(param string) (value string) {
	value = c.Param(param)
	return
}

func (c *Context) GetParamInt(param string) (value int) {
	paramValue := c.Context.QueryParam(param)
	if paramValue != "" {
		value, _ = strconv.Atoi(paramValue)
	}
	return
}

func (c *Context) GetParamArrayInt(param string) (values []int) {
	paramValue := c.Context.QueryParam(param)
	if paramValue != "" {
		paramValues := strings.Split(paramValue, ",")
		for _, valueStr := range paramValues {
			var value int
			value, _ = strconv.Atoi(valueStr)
			values = append(values, value)
		}
	}
	return
}

func (c *Context) GetParamFloat64(param string) (value float64) {
	paramValue := c.Context.QueryParam(param)
	if paramValue != "" {
		value, _ = strconv.ParseFloat(paramValue, 64)
	}
	return
}

func (c *Context) GetParamString(param string) (value string) {
	value = c.Context.QueryParam(param)
	return
}

func (c *Context) GetParamArrayString(param string) (values []string) {
	paramValue := c.Context.QueryParam(param)
	if paramValue != "" {
		values = strings.Split(paramValue, ",")
	}
	return
}

func (c *Context) GetParamDate(param string) (values time.Time) {
	paramValue := c.Context.QueryParam(param)
	if paramValue != "" {
		values, _ = time.Parse(timex.InFormatDate, paramValue)
	}
	return
}

func (c *Context) GetParamDateTime(param string) (values time.Time) {
	paramValue := c.Context.QueryParam(param)
	if paramValue != "" {
		values, _ = time.Parse(timex.InFormatDateTime, paramValue)
	}
	return
}

func (c *Context) GetParamTime(param string) (values time.Time) {
	paramValue := c.Context.QueryParam(param)
	if paramValue != "" {
		values, _ = time.Parse(timex.InFormatTime, paramValue)
	}
	return
}

// Data set data and total into response format
func (c *Context) Data(data interface{}, total ...int64) {
	c.ResponseFormat.SetData(data)
}

// DataList set data list into response format
func (c *Context) DataList(data interface{}, total int64, page *Paginator) {
	c.ResponseFormat.SetDataList(data, total, page)
}

// Failure set response format errors
// its equal with validation errors.
func (c *Context) Failure(fail ...string) {
	c.ResponseFormat.Errors = map[string]string{fail[0]: fail[1]}
}

// Serve response json data with data that already collected
// if error is not nill will returning error responses.
func (c *Context) Serve(e error) (err error) {
	c.ResponseFormat.Code = http.StatusOK
	if e != nil {
		c.ResponseFormat.SetError(e)
	} else {
		if c.ResponseData != nil {
			c.ResponseFormat.SetData(c.ResponseData)
		}
	}

	if c.Request().Method == echo.HEAD || c.Request().Method == echo.OPTIONS {
		err = c.NoContent(http.StatusNoContent)
	} else {
		err = c.JSON(c.ResponseFormat.Code, c.ResponseFormat)
	}

	c.ResponseFormat.Reset()

	return
}

func (c *Context) Message(statusCode int, status string, message string) (err error) {
	c.ResponseFormat.Code = statusCode
	c.ResponseFormat.Status = status
	c.ResponseFormat.Message = message

	if c.Request().Method == echo.HEAD || c.Request().Method == echo.OPTIONS {
		err = c.NoContent(http.StatusNoContent)
	} else {
		err = c.JSON(c.ResponseFormat.Code, c.ResponseFormat)
	}
	c.ResponseFormat.Reset()
	return
}

const PerPage = 10

// GetPage get params page for pagination
func (c *Context) GetPage() int {
	p := c.QueryParam("page")

	if p == "" {
		return 1
	}
	page, err := strconv.Atoi(p)
	if err != nil {
		page = 1
	}
	return page
}

// GetPerPage get params per_page for pagination
func (c *Context) GetPerPage() int {
	p := c.QueryParam("per_page")
	perPage, err := strconv.Atoi(p)
	if err != nil {
	}
	return perPage
}
