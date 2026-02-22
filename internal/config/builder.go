package config

import (
	"errors"
	"os"
)

type Builder struct {
	filePaths []string
}

func NewBuilder() *Builder {
	return &Builder{
		filePaths: []string{},
	}
}

func (b *Builder) AddConfigFile(path string) *Builder {
	b.filePaths = append(b.filePaths, path)
	return b
}

func (b *Builder) Build() (*Config, error) {
	var content []byte
	var err error

	for _, path := range b.filePaths {
		content, err = os.ReadFile(path)
		if err == nil {
			break
		}

		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return Parse(content)
}
