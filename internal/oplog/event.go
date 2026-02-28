package oplog

import (
	"bytes"
	"encoding/json"
)

func UnmarshalEvent(input map[string]any, v any) error {
	buffer := &bytes.Buffer{}

	err := json.NewEncoder(buffer).Encode(input)
	if err != nil {
		return err
	}

	err = json.NewDecoder(buffer).Decode(v)
	if err != nil {
		return err
	}

	return nil
}
