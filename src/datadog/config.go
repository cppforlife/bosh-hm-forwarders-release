package main

import (
	"fmt"
)

type Config struct {
	APIKey string
}

func (c Config) Validate() error {
	if len(c.APIKey) == 0 {
		return fmt.Errorf("Missing 'APIKey'")
	}

	return nil
}
