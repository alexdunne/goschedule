package http_test

import (
	"goschedule/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPOSTAuthLogin(t *testing.T) {
	server := http.NewServer()

	t.Run("Create's a new account using Github OAuth", func(t *testing.T) {
		request := httptest.NewRequest("POST", "/auth/login", nil)
		rr := httptest.NewRecorder()

		server.ServeHTTP(rr, request)

		assertStatus(t, rr, 200)
		assertResponseBody(t, rr, `{"user":{"id":"abc"}}`)
	})
}

func assertStatus(t testing.TB, rr *httptest.ResponseRecorder, want int) {
	t.Helper()

	got := rr.Code

	if got != want {
		t.Errorf("response status is wrong, got %q want %q", got, want)
	}
}

func assertResponseBody(t testing.TB, rr *httptest.ResponseRecorder, want string) {
	t.Helper()

	got := strings.TrimSpace(rr.Body.String())

	if got != want {
		t.Errorf("response body is wrong, got %q want %q", got, want)
	}
}
