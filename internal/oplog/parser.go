package oplog

import "errors"

var ErrPropertyDoesNotExist = errors.New("event: property does not exist")
var ErrPropertyIsNotAString = errors.New("event: property is not of type string")

func ParseString(evt Event, key string) (string, error) {
	raw, exists := evt.Properties[key]
	if !exists {
		return "", ErrPropertyDoesNotExist
	}

	stringValue, ok := raw.(string)
	if !ok {
		return "", ErrPropertyIsNotAString
	}

	return stringValue, nil
}

func ParseBoolOrDefault(evt Event, key string, defaultValue bool) bool {
	raw, exists := evt.Properties[key]
	if !exists {
		return defaultValue
	}

	boolValue, ok := raw.(bool)
	if !ok {
		return defaultValue
	}

	return boolValue
}
