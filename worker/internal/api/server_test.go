package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/0xmigzy/capstone-worker/internal/cache"
)

func TestSubmissionRejectsUnauthorizedAndInvalidRequests(t *testing.T) {
	server := NewServer(nil, nil, nil, "secret-token", "")
	for _, test := range []struct {
		name, auth, body string
		status           int
	}{
		{"missing token", "", `{"pair":"ETH/USD"}`, http.StatusUnauthorized},
		{"wrong token", "Bearer wrong", `{"pair":"ETH/USD"}`, http.StatusUnauthorized},
		{"unsupported pair", "Bearer secret-token", `{"pair":"DOGE/USD"}`, http.StatusBadRequest},
		{"unknown field", "Bearer secret-token", `{"pair":"ETH/USD","extra":1}`, http.StatusBadRequest},
		{"trailing JSON", "Bearer secret-token", `{"pair":"ETH/USD"}{}`, http.StatusBadRequest},
		{"oversized body", "Bearer secret-token", `{"pair":"` + strings.Repeat("x", maxBodyBytes) + `"}`, http.StatusBadRequest},
	} {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/oracle/request", strings.NewReader(test.body))
			request.Header.Set("Authorization", test.auth)
			response := httptest.NewRecorder()
			server.handleOracleRequest(response, request)
			if response.Code != test.status {
				t.Fatalf("status=%d want=%d", response.Code, test.status)
			}
		})
	}
}

func TestJobResponseNeverExposesSignedTransaction(t *testing.T) {
	job := cache.Job{ID: "job", Pair: "ETH/USD", PriceDecimals: 8, SignedTransaction: []byte("private-signed-bytes")}
	job.Transition("processing")
	response := httptest.NewRecorder()
	jsonResponse(response, http.StatusOK, job)
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body.Bytes(), &fields); err != nil {
		t.Fatal(err)
	}
	for _, forbidden := range []string{"SignedTransaction", "signedTransaction"} {
		if _, exists := fields[forbidden]; exists {
			t.Fatalf("leaked field %s", forbidden)
		}
	}
	if response.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("job response must not be cached")
	}
}
