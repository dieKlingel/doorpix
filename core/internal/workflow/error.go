package workflow

import (
	"errors"
)

var ErrProviderAlreadyRegistered = errors.New("provider is already registered")
var ErrRegistryIsNil = errors.New("registry is nil")
var ErrPipelineIsNil = errors.New("pipeline is nil")
var ErrProviderNotFound = errors.New("provider is not found")
var ErrSourceIsNil = errors.New("source is nil")
