package core_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dieklingel/doorpix/core"
	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/stretchr/testify/assert"
)

func TestHttpHandler_HandleEvent(t *testing.T) {
	config := doorpix.NewConfig()
	system := doorpix.System{
		Config: config,
		Bus:    core.NewEventEmitterWithConfig(config),
	}
	handler := core.HttpHandler{
		System: system,
	}

	t.Run("should handle events", func(t *testing.T) {
		events := []string{
			"api:ring",
			"api:unlock",
		}
		for _, event := range events {
			t.Run("should handle event "+event, func(t *testing.T) {

				body := strings.NewReader(`{"event": "` + event + `"}`)
				rr := httptest.NewRecorder()
				req, err := http.NewRequest("POST", "/api/events", body)
				assert.NoError(t, err)

				handler.HandleEmitEvent(rr, req)
				assert.Equal(t, http.StatusOK, rr.Code)
			})
		}
	})

	t.Run("should return 400 on invalid event", func(t *testing.T) {
		events := []string{
			"api:invalid",
			"api:*",
			"startup",
			"shutdown",
			"incoming-call",
			"1234",
			"0x00",
		}
		for _, event := range events {
			t.Run("should return 400 on invalid event "+event, func(t *testing.T) {

				body := strings.NewReader(`{"event": "` + event + `"}`)
				rr := httptest.NewRecorder()
				req, err := http.NewRequest("POST", "/api/events", body)
				assert.NoError(t, err)

				handler.HandleEmitEvent(rr, req)
				assert.Equal(t, http.StatusBadRequest, rr.Code)
			})
		}
	})
}
