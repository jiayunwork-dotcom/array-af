package web

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()
	fsys := fstest.MapFS{
		"index.html": &fstest.MapFile{Data: []byte("<html><body>array-af console</body></html>")},
		"app.js":     &fstest.MapFile{Data: []byte("// test asset")},
		"style.css":  &fstest.MapFile{Data: []byte("body{}")},
	}
	return NewServer(Assets{
		WebFS: fsys,
		Examples: map[string][]byte{
			"broadside-8": []byte(`{"N":8,"d":0.5,"lambda":1.0,"beta":0.0,"element":"iso"}`),
		},
	})
}

func doReq(t *testing.T, s *Server, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	var rdr *strings.Reader
	if body == "" {
		rdr = strings.NewReader("")
	} else {
		rdr = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, path, rdr)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	s.ServeHTTP(rec, req)
	return rec
}

func TestAFEndpointBroadside(t *testing.T) {
	s := newTestServer(t)
	rec := doReq(t, s, http.MethodPost, "/api/af",
		`{"N":8,"d":0.5,"lambda":1.0,"beta":0.0,"element":"iso"}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}
	var resp afResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Mainlobe.Visible != true || resp.Mainlobe.AngleDeg == nil {
		t.Errorf("mainlobe: got visible=%v angle=%v, want visible true angle 90", resp.Mainlobe.Visible, resp.Mainlobe.AngleDeg)
	} else if *resp.Mainlobe.AngleDeg != 90 {
		t.Errorf("mainlobe angle: got %v, want 90", *resp.Mainlobe.AngleDeg)
	}
	if resp.AfPeak != 8 {
		t.Errorf("af_peak: got %v, want 8", resp.AfPeak)
	}
	if !resp.AfPeakMatchesN {
		t.Errorf("af_peak_matches_n: got false, want true")
	}
	if resp.Grating.Present {
		t.Errorf("grating: expected none for d=lambda/2")
	}
	if len(resp.Points) == 0 {
		t.Errorf("points: expected non-empty AF sample list")
	}
}

func TestAFEndpointInvalidInput(t *testing.T) {
	s := newTestServer(t)
	cases := []struct {
		name string
		body string
		want string
	}{
		{"N=1", `{"N":1,"d":0.5,"lambda":1.0}`, "N"},
		{"d=0", `{"N":8,"d":0,"lambda":1.0}`, "d"},
		{"lambda=-1", `{"N":8,"d":0.5,"lambda":-1.0}`, "lambda"},
		{"missing N", `{"d":0.5,"lambda":1.0}`, "N"},
		{"missing d", `{"N":8,"lambda":1.0}`, "d"},
		{"missing lambda", `{"N":8,"d":0.5}`, "lambda"},
		{"bad element", `{"N":8,"d":0.5,"lambda":1.0,"element":"magic"}`, "element"},
	}
	for _, c := range cases {
		rec := doReq(t, s, http.MethodPost, "/api/af", c.body)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%s: status got %d, want 400", c.name, rec.Code)
			continue
		}
		var eb errorBody
		if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
			t.Errorf("%s: error body not JSON: %v (%s)", c.name, err, rec.Body.String())
			continue
		}
		if eb.Error == "" {
			t.Errorf("%s: error body has empty error field", c.name)
			continue
		}
		if !strings.Contains(eb.Error, c.want) {
			t.Errorf("%s: error %q does not mention %q", c.name, eb.Error, c.want)
		}
	}
}

func TestAFEndpointBadJSON(t *testing.T) {
	s := newTestServer(t)
	rec := doReq(t, s, http.MethodPost, "/api/af", `{"N":8,`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", rec.Code)
	}
	var eb errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
		t.Errorf("error body not JSON: %v", err)
	}
	if !strings.Contains(eb.Error, "invalid JSON") {
		t.Errorf("error %q should mention invalid JSON", eb.Error)
	}
}

func TestScanEndpoint(t *testing.T) {
	s := newTestServer(t)
	rec := doReq(t, s, http.MethodPost, "/api/scan",
		`{"N":8,"d":0.5,"lambda":1.0,"steps":8}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200 (body=%s)", rec.Code, rec.Body.String())
	}
	var resp scanResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Rows) != 9 {
		t.Errorf("rows: got %d, want 9", len(resp.Rows))
	}
	if resp.FirstMainlobeAngleDeg == nil || *resp.FirstMainlobeAngleDeg != 90 {
		t.Errorf("first mainlobe: got %v, want 90", resp.FirstMainlobeAngleDeg)
	}
	if resp.LastMainlobeAngleDeg == nil || *resp.LastMainlobeAngleDeg != 0 {
		t.Errorf("last mainlobe: got %v, want 0", resp.LastMainlobeAngleDeg)
	}
	if !resp.Summary.MainlobeMovesTowardEndfire {
		t.Errorf("expected mainlobe to move toward endfire in summary")
	}
}

func TestScanEndpointInvalidInput(t *testing.T) {
	s := newTestServer(t)
	rec := doReq(t, s, http.MethodPost, "/api/scan",
		`{"N":1,"d":0.5,"lambda":1.0}`)
	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want 400", rec.Code)
	}
	var eb errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
		t.Fatalf("error body not JSON: %v", err)
	}
	if eb.Error == "" {
		t.Errorf("expected non-empty error field")
	}
}

func TestExamplesEndpoint(t *testing.T) {
	s := newTestServer(t)
	rec := doReq(t, s, http.MethodGet, "/api/examples", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	var resp map[string]json.RawMessage
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	var payloads map[string]json.RawMessage
	if err := json.Unmarshal(resp["payloads"], &payloads); err != nil {
		t.Fatalf("decode payloads: %v", err)
	}
	if _, ok := payloads["broadside-8"]; !ok {
		t.Errorf("expected example broadside-8 in payloads")
	}
}

func TestIndexPageServed(t *testing.T) {
	s := newTestServer(t)
	rec := doReq(t, s, http.MethodGet, "/", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	if !strings.Contains(rec.Body.String(), "array-af console") {
		t.Errorf("index page body missing expected marker")
	}
}

func TestStaticAssetServed(t *testing.T) {
	s := newTestServer(t)
	rec := doReq(t, s, http.MethodGet, "/static/app.js", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/javascript") {
		t.Errorf("content type: got %q, want text/javascript", ct)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	s := newTestServer(t)
	rec := doReq(t, s, http.MethodGet, "/api/af", "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status: got %d, want 405", rec.Code)
	}
	var eb errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
		t.Fatalf("error body not JSON: %v", err)
	}
	if eb.Error == "" {
		t.Errorf("expected non-empty error field")
	}
}

func TestUnknownPageNotFound(t *testing.T) {
	s := newTestServer(t)
	rec := doReq(t, s, http.MethodGet, "/no/such/page", "")
	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want 404", rec.Code)
	}
	var eb errorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &eb); err != nil {
		t.Fatalf("error body not JSON: %v", err)
	}
	if eb.Error == "" {
		t.Errorf("expected non-empty error field")
	}
}
