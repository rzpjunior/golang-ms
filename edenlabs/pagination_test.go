package edenlabs

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_Paginator(t *testing.T) {
	tests := []struct {
		page       int
		perPage    int
		expected   *Paginator
		totalItems int64
		err        error
	}{
		{
			page:    1,
			perPage: 10,
			expected: &Paginator{
				Page:     1,
				PerPage:  10,
				Limit:    10,
				Start:    0,
				End:      10,
				NumPages: 2,
				HasPrev:  false,
				HasNext:  true,
			},
			totalItems: 17,
			err:        nil,
		},
		{
			page:    2,
			perPage: 10,
			expected: &Paginator{
				Page:     2,
				PerPage:  10,
				Offset:   10,
				Limit:    10,
				Start:    10,
				End:      20,
				NumPages: 2,
				HasPrev:  true,
				HasNext:  false,
			},
			totalItems: 17,
			err:        nil,
		},
		{
			page:       -1,
			perPage:    -1,
			totalItems: 10,
			expected:   &Paginator{},
			err:        errors.New("Invalid page number"),
		},
		{
			page:       1,
			perPage:    1,
			totalItems: 1,
			expected: &Paginator{
				Page:     1,
				PerPage:  1,
				Offset:   0,
				Limit:    1,
				Start:    0,
				End:      1,
				NumPages: 1,
				HasPrev:  false,
				HasNext:  false,
			},
			err: nil,
		},
	}
	for _, tt := range tests {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/?page=%d&per_page=%d", tt.page, tt.perPage), nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		ctx := &Context{
			Context:        c,
			ResponseFormat: NewResponse(),
			ResponseData:   nil,
		}

		var pagination *Paginator
		var err error
		pagination, err = NewPaginator(ctx)

		if tt.err != nil {
			assert.NotNil(t, err)
			assert.Nil(t, pagination)
		} else {
			pagination.Json(tt.totalItems)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, pagination)
		}
	}
}
