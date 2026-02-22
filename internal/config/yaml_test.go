package config_test

import (
	"testing"

	"github.com/dieklingel/doorpix/internal/config"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestBoolUnmarshalYAML(t *testing.T) {
	t.Run("should parse truthy boolean values", func(t *testing.T) {
		truhtyValues := []string{
			"value: 'true'",
			"value: true",
			"value: 1",
		}

		for _, input := range truhtyValues {
			parsed := struct {
				Value config.Bool `yaml:"value"`
			}{}
			err := yaml.Unmarshal([]byte(input), &parsed)

			assert.NoError(t, err)
			assert.True(t, bool(parsed.Value))
		}
	})

	t.Run("should parse falsy boolean values", func(t *testing.T) {
		falsyValues := []string{
			"value: 'false'",
			"value: false",
			"value: 0",
		}

		for _, input := range falsyValues {
			parsed := struct {
				Value config.Bool `yaml:"value"`
			}{}
			err := yaml.Unmarshal([]byte(input), &parsed)

			assert.NoError(t, err)
			assert.False(t, bool(parsed.Value))
		}
	})

	t.Run("should return error for invalid values", func(t *testing.T) {
		invalidValues := []string{
			"value: 'some string'",
			"value: falsy",
			"value: -1",
		}

		for _, input := range invalidValues {
			parsed := struct {
				Value config.Bool `yaml:"value"`
			}{}

			err := yaml.Unmarshal([]byte(input), &parsed)
			assert.ErrorIs(t, err, config.ErrUnexpectedInput)
			assert.False(t, bool(parsed.Value))
		}
	})
}
