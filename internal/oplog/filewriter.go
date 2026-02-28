package oplog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type FileWriter struct {
	File string
}

func (f *FileWriter) Write(evt Event) error {
	payload, err := json.Marshal(evt)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOplogWriteFailed, err)
	}

	buffer := &bytes.Buffer{}
	err = json.Compact(buffer, payload)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOplogWriteFailed, err)
	}

	_, err = buffer.WriteString("\r\n")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOplogWriteFailed, err)
	}

	file, err := os.OpenFile(f.File, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOplogWriteFailed, err)
	}
	defer file.Close()

	_, err = file.Write(buffer.Bytes())
	if err != nil {
		return fmt.Errorf("%w: %w", ErrOplogWriteFailed, err)
	}

	return nil
}
