package sip_test

import (
	"testing"

	"github.com/dieklingel/doorpix/internal/transport/sip"
	"github.com/stretchr/testify/assert"
)

func TestUri(t *testing.T) {

	t.Run("should return uri", func(t *testing.T) {
		inputs := []string{
			"\"test\" <sip:test@sip.example.org>",
			"\"test\" <sip:test@sip.example.org>",
			"<sip:test@sip.example.org>",
			"sip:test@sip.example.org",
		}

		for _, input := range inputs {
			uri := sip.Uri(input)
			assert.Equal(t, "sip:test@sip.example.org", uri)
		}
	})

	t.Run("should return error", func(t *testing.T) {
		inputs := []string{
			"example.org",
			"http://example.org",
			"\"Hallo Welt\"",
			"<>",
			"\"\" <sip:@example.org>",
			"sip:test@",
		}

		for _, input := range inputs {
			uri := sip.Uri(input)
			assert.Equal(t, "", uri)
		}
	})
}
