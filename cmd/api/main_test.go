package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMainApi(t *testing.T) {
	server := NewApplication()
	t.Run("tests if the API server is running", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/", nil)
		require.Nil(t, err)

		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		require.Equal(t, 200, res.Code)
		require.Equal(t, "Hello World", res.Body.String())
	})

	t.Run("tests if the ping-pong handler is working as expected", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, "/ping", nil)
		require.Nil(t, err)

		res := httptest.NewRecorder()

		server.ServeHTTP(res, req)

		require.Equal(t, 200, res.Code)
		require.Equal(t, "PONG", res.Body.String())
	})
}
