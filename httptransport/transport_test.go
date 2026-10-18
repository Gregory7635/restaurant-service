package httptransport

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-kit/log"

	"example.com/restaurant-service/implementation"
	"example.com/restaurant-service/repository"
	"example.com/restaurant-service/transport"
)

func newTestServer() *httptest.Server {
	svc := implementation.NewService(repository.NewMemoryRepository(), log.NewNopLogger())
	return httptest.NewServer(NewHandler(transport.MakeEndpoints(svc), log.NewNopLogger()))
}

func do(t *testing.T, method, url, body string) (*http.Response, map[string]interface{}) {
	t.Helper()
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var out map[string]interface{}
	_ = json.NewDecoder(resp.Body).Decode(&out)
	return resp, out
}

func TestHTTPFlow(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, out := do(t, "POST", srv.URL+"/restaurants", `{"name":"Пицца-Хаус","address":"ул. Ленина, 1","rating":4.2,"open":true}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("POST status = %d", resp.StatusCode)
	}
	id, _ := out["id"].(string)
	if id == "" {
		t.Fatal("POST: empty id")
	}

	resp, _ = do(t, "PATCH", srv.URL+"/restaurants/"+id+"/rating", `{"rating":4.9}`)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("PATCH status = %d", resp.StatusCode)
	}

	resp, out = do(t, "GET", srv.URL+"/restaurants/"+id, "")
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("GET status = %d", resp.StatusCode)
	}
	r, _ := out["restaurant"].(map[string]interface{})
	if r["rating"] != 4.9 {
		t.Errorf("rating = %v, want 4.9", r["rating"])
	}
}

func TestHTTPErrors(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	tests := []struct {
		name, method, path, body string
		want                     int
	}{
		{"not found", "GET", "/restaurants/missing", "", http.StatusNotFound},
		{"invalid input", "POST", "/restaurants", `{"rating":3}`, http.StatusBadRequest},
		{"bad json", "POST", "/restaurants", `{bad`, http.StatusBadRequest},
		{"invalid rating", "PATCH", "/restaurants/missing/rating", `{"rating":10}`, http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := do(t, tt.method, srv.URL+tt.path, tt.body)
			if resp.StatusCode != tt.want {
				t.Errorf("status = %d, want %d", resp.StatusCode, tt.want)
			}
		})
	}
}
