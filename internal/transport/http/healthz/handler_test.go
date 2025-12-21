package healthz_test

import (
	"net/http/httptest"
	"testing"

	"github.com/dieklingel/doorpix/internal/transport/http/healthz"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	t.Run("serve /", func(t *testing.T) {
		handler := healthz.Handler()

		req := httptest.NewRequest("GET", "http://example.com/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)

		resp := w.Result()
		assert.Equal(t, 200, resp.StatusCode)

		body := w.Body.String()
		assert.Equal(t, "OK\r\n", body)
	})
}
