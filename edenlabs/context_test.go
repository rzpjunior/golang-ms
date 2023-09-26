package edenlabs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"git.edenfarm.id/edenlabs/edenlabs/validation"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func fakeContext(method string, endpoint string, response string, rec *httptest.ResponseRecorder) (*Context, error) {
	e := New()
	e.HTTPErrorHandler = HTTPErrorHandler

	req, err := http.NewRequest(method, endpoint, strings.NewReader(response))
	ctx := e.NewContext(req, rec)
	c := NewContext(ctx)

	return c, err
}

func TestNewResponse(t *testing.T) {
	r := NewResponse()

	assert.Equal(t, HTTPResponseSuccess, r.Status)
}

func TestContextData(t *testing.T) {
	type user struct {
		ID   int    `json:"id" xml:"id" form:"id"`
		Name string `json:"name" xml:"name" form:"name"`
	}

	js := `{"code":200,"status":"success","data":{"id":1,"name":"Jon Snow"}}
`

	rec := httptest.NewRecorder()
	ctx, err := fakeContext(echo.POST, "/", js, rec)
	ctx.Data(user{1, "Jon Snow"}, 20)

	err = ctx.Serve(err)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, js, rec.Body.String())
	}

	rec = httptest.NewRecorder()
	ctx, err = fakeContext(echo.OPTIONS, "/", "", rec)
	err = ctx.Serve(err)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Equal(t, "", rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, "", string(rec.Body.Bytes()))
	}
}

func TestContext_DataList(t *testing.T) {
	type user struct {
		ID   int    `json:"id" xml:"id" form:"id"`
		Name string `json:"name" xml:"name" form:"name"`
	}

	js := `{"code":200,"status":"success","data":{"id":1,"name":"Jon Snow"},"total":20,"page":1,"per_page":10}
`

	rec := httptest.NewRecorder()
	ctx, err := fakeContext(echo.POST, "/", js, rec)
	ctx.SetParamNames("page", "per_page")
	ctx.SetParamValues("1", "100")
	page, _ := NewPaginator(ctx)

	ctx.DataList(user{1, "Jon Snow"}, 20, page)

	err = ctx.Serve(err)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, js, rec.Body.String())
	}
}

func TestContextResponseErrorType(t *testing.T) {
	js := `{"code":422,"status":"success","message":"Unprocessable Entity","errors":{"name":"The name field is required."}}
`

	// validation error
	rec := httptest.NewRecorder()
	ctx, _ := fakeContext(echo.POST, "/", js, rec)

	type user struct {
		Name string `json:"name" valid:"required"`
	}

	v := validation.New()
	o := v.Struct(user{""})

	err := ctx.Serve(o)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, js, rec.Body.String())
	}

	// http error
	rec = httptest.NewRecorder()
	ctx, _ = fakeContext(echo.POST, "/", "", rec)
	err = ctx.Serve(echo.ErrNotFound)
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSONCharsetUTF8, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, `{"code":404,"status":"failure","message":"Not Found","errors":{"error":"Not Found"}}
`, rec.Body.String())
	}

}

func TestContext_GetParamID(t *testing.T) {
	rec := httptest.NewRecorder()
	ctx, err := fakeContext(echo.GET, "/", ``, rec)
	ctx.SetParamValues("1")
	id, err := ctx.GetParamID()
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), id)

}
