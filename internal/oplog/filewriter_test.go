package oplog_test

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/dieklingel/doorpix/internal/oplog"
	"github.com/stretchr/testify/assert"
)

func TestFileWriter_Write(t *testing.T) {
	t.Run("should create file", func(t *testing.T) {
		filename := fmt.Sprintf("filewriter_test.go.%s.jsonl", rand.Text())
		writer := &oplog.FileWriter{
			File: filename,
		}

		assert.NoFileExists(t, filename)
		err := writer.Write(oplog.Event{})
		assert.NoError(t, err)
		assert.FileExists(t, filename)
		os.Remove(filename)
	})

	t.Run("should append to file", func(t *testing.T) {
		filename := fmt.Sprintf("filewriter_test.go.%s.jsonl", rand.Text())
		os.WriteFile(filename, []byte("{}\r\n"), 0644)
		writer := &oplog.FileWriter{
			File: filename,
		}

		event := oplog.Event{}
		err := writer.Write(event)
		assert.NoError(t, err)

		c, err := os.ReadFile(filename)
		assert.NoError(t, err)
		content := string(c)

		lines := strings.Split(content, "\r\n")
		assert.Equal(t, "{}", lines[0])
		assert.Len(t, lines, 3)

		os.Remove(filename)
	})

	t.Run("should write oplog event as json", func(t *testing.T) {
		filename := fmt.Sprintf("filewriter_test.go.%s.jsonl", rand.Text())
		writer := &oplog.FileWriter{
			File: filename,
		}

		timestamp, err := time.Parse(time.DateOnly, "2020-01-02")
		assert.NoError(t, err)

		input := oplog.Event{
			Id:        "1234",
			Timestamp: timestamp,
			Path:      "system/test",
			Properties: map[string]any{
				"Hello": "World",
			},
		}

		err = writer.Write(input)
		assert.NoError(t, err)

		output := oplog.Event{}
		c, err := os.ReadFile(filename)
		assert.NoError(t, err)
		err = json.Unmarshal(c, &output)
		assert.NoError(t, err)

		assert.Equal(t, input, output)
		os.Remove(filename)
	})

	t.Run("should return erro on wrong file permission", func(t *testing.T) {
		filename := fmt.Sprintf("filewriter_test.go.%s.jsonl", rand.Text())
		os.WriteFile(filename, []byte("{}\r\n"), 0000)
		writer := &oplog.FileWriter{
			File: filename,
		}

		err := writer.Write(oplog.Event{})
		assert.ErrorIs(t, err, oplog.ErrOplogWriteFailed)

		os.Remove(filename)
	})
}
