// SPDX-License-Identifier: BSD-3-Clause
// Copyright (c) 2022 Intel Corporation

package cfg

import (
    "bytes"
    "fmt"
    "os"
	"encoding/json"

    "github.com/tidwall/jsonc"
)
type Config struct {
}

type System struct {
    cfg *Config
}

func (cfg *Config) validateConfig() error {
    return nil
}

// OpenWithConfig by passing in a initialized Config structure
func OpenWithConfig(c *Config) (*System, error) {

    if err := c.validateConfig(); err != nil {
        return nil, fmt.Errorf("failed to validate configuration: %w", err)
    }

    sys := &System{cfg: c}

    return sys, nil
}

func OpenWithText(b []byte) (*System, error) {

    text := jsonc.ToJSON(bytes.TrimSpace(b))

	if len(text) == 0 {
		return nil, fmt.Errorf("empty json text string")
	}

	// test for JSON string, which must start with a '{'
	if text[0] != '{' {
		return nil, fmt.Errorf("string does not appear to be a valid JSON text missing starting '{'")
	}

	cfg := &Config{}

	// Unmarshal json text into the Config structure
	if err := json.Unmarshal(text, cfg); err != nil {
		return nil, err
	}
    return OpenWithConfig(cfg)
}

// OpenWithFile by passing in a filename or path to a JSON-C or JSON configuration
func OpenWithFile(path string) (*System, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    return OpenWithText(b)
}

func (c *Config) String() string {

	if data, err := json.MarshalIndent(c, "", "  "); err != nil {
		return fmt.Sprintf("error marshalling JSON: %v", err)
	} else {
		return string(data)
	}
}