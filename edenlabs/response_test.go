package edenlabs

import (
	"git.edenfarm.id/edenlabs/edenlabs/orm"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewResponse(t *testing.T) {
	tests := []struct {
		expected *FormatResponse
	}{
		{
			expected: &FormatResponse{
				Code:   http.StatusOK,
				Status: HTTPResponseSuccess,
			},
		},
	}
	for _, tt := range tests {
		res := NewResponse()
		assert.Equal(t, tt.expected, res)
	}
}

func Test_SetMessage(t *testing.T) {
	tests := []struct {
		expected *FormatResponse
	}{
		{
			expected: &FormatResponse{
				Code:    http.StatusOK,
				Status:  HTTPResponseSuccess,
				Message: "any-message",
			},
		},
	}
	for _, tt := range tests {
		res := NewResponse()
		res.SetMessage("any-message")
		assert.Equal(t, tt.expected, res)
	}
}

func Test_SetData(t *testing.T) {
	tests := []struct {
		expected *FormatResponse
	}{
		{
			expected: &FormatResponse{
				Code:   http.StatusOK,
				Status: HTTPResponseSuccess,
				Data: struct {
					ID   int    "json:\"id\""
					Name string "json:\"name\""
				}{
					ID:   1,
					Name: "any-name",
				},
			},
		},
	}
	for _, tt := range tests {
		res := NewResponse()
		res.SetData(struct {
			ID   int    "json:\"id\""
			Name string "json:\"name\""
		}{
			ID:   1,
			Name: "any-name",
		})

		assert.Equal(t, tt.expected, res)
	}
}

func Test_SetDataList(t *testing.T) {
	tests := []struct {
		page     *Paginator
		total    int64
		expected *FormatResponse
	}{
		{
			page: &Paginator{
				Page:    1,
				PerPage: 100,
			},
			total: 2,
			expected: &FormatResponse{
				Code:   http.StatusOK,
				Status: HTTPResponseSuccess,
				Data: []struct {
					ID   int    "json:\"id\""
					Name string "json:\"name\""
				}{
					{
						ID:   1,
						Name: "any-name",
					},
					{
						ID:   2,
						Name: "any-name",
					},
				},
				Total:   2,
				Page:    1,
				PerPage: 100,
			},
		},
	}
	for _, tt := range tests {
		res := NewResponse()
		res.SetDataList([]struct {
			ID   int    "json:\"id\""
			Name string "json:\"name\""
		}{
			{
				ID:   1,
				Name: "any-name",
			},
			{
				ID:   2,
				Name: "any-name",
			},
		}, tt.total, tt.page)

		assert.Equal(t, tt.expected, res)
	}
}

func TestFormatResponse_SetError(t *testing.T) {
	var err error
	res := NewResponse()
	err = orm.ErrTxDone
	res.SetError(err)
	assert.Equal(t, map[string]string{
		"error": "Transaction already done",
	}, res.Errors)

	res = NewResponse()
	err = orm.ErrMultiRows
	res.SetError(err)
	assert.Equal(t, map[string]string{
		"error": "Return multi rows",
	}, res.Errors)

	err = orm.ErrNoRows
	res.SetError(err)
	assert.Equal(t, map[string]string{
		"error": "No row Found",
	}, res.Errors)
	err = orm.ErrStmtClosed
	res.SetError(err)
	assert.Equal(t, map[string]string{
		"error": "Stmt already closed",
	}, res.Errors)
	err = orm.ErrArgs
	res.SetError(err)
	assert.Equal(t, map[string]string{
		"error": "Args error may be empty",
	}, res.Errors)
	err = orm.ErrNotImplement
	res.SetError(err)
	assert.Equal(t, map[string]string{
		"error": "Have not implement",
	}, res.Errors)
	err = orm.ErrLastInsertIdUnavailable
	res.SetError(err)
	assert.Equal(t, map[string]string{
		"error": "Last insert id is unavailable",
	}, res.Errors)
}
