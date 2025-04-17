package actions_test

import (
	"testing"

	"github.com/dieklingel/doorpix/core/internal/doorpix"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func Test_YamlStringOrList_UnmarshalYAML(t *testing.T) {
	type Fields struct {
		Value doorpix.YamlScalarOrList[string] `yaml:"value"`
	}

	t.Run("scalar", func(t *testing.T) {
		value := `value: "test"`
		var fields Fields

		err := yaml.Unmarshal([]byte(value), &fields)
		assert.NoError(t, err)

		assert.Len(t, fields.Value, 1)
		assert.Equal(t, "test", fields.Value[0])
	})

	t.Run("list", func(t *testing.T) {
		value := `
value:
  - test
  - test2
`
		var fields Fields

		err := yaml.Unmarshal([]byte(value), &fields)
		assert.NoError(t, err)

		assert.Len(t, fields.Value, 2)
		assert.Equal(t, "test", fields.Value[0])
		assert.Equal(t, "test2", fields.Value[1])
	})
}
