package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.edenfarm.id/project-version2/api/util"
)

func TestPublishWrongResponseStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer ts.Close()
	err := util.SendIDToOapi("296222731", "")
	if err != nil {
		t.Errorf("SendIDToOapi() returned an error: %s", err)
	}
}
