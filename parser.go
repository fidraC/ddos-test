package main

import "encoding/json"

type Config struct {
	Host      string     `json:"host"`
	Endpoints []Endpoint `json:"endpoints"`
	Threads   int        `json:"threads"`
}

type Endpoint struct {
	Path    string            `json:"path"`
	Method  string            `json:"method"`
	Payload *Payload          `json:"payload"`
	Headers map[string]string `json:"headers"`
}

type Payload struct {
	Params string `json:"params"`
	Data   string `json:"data"`
}

func validateConfig(c *Config) error {
	if c.Host == "" {
		return ErrMissingHost
	}
	for i, e := range c.Endpoints {
		if e.Path == "" {
			return ErrMissingPath
		}
		if e.Method == "" {
			c.Endpoints[i].Method = "GET"
		}
	}
	return nil
}

func NewConfig(s string) (*Config, error) {
	var c Config
	err := json.Unmarshal([]byte(s), &c)
	if err != nil {
		return nil, err
	}
	if err := validateConfig(&c); err != nil {
		return nil, err
	}
	return &c, nil
}
