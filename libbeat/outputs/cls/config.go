package cls

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/elastic/beats/v7/libbeat/outputs/codec"
	"github.com/elastic/elastic-agent-libs/config"
)

type clsConfig struct {
	Endpoint    string           `config:"endpoint"`
	Topic       string           `config:"topic"`
	AccessKey   string           `config:"access_key"`
	SecretKey   string           `config:"secret_key"`
	BulkMaxSize int              `config:"bulk_max_size"`
	MaxRetries  int              `config:"max_retries"         validate:"min=-1,nonzero"`
	Codec       codec.Config     `config:"codec"`
	Queue       config.Namespace `config:"queue"`
}

func defaultConfig() clsConfig {
	return clsConfig{
		Topic:       "",
		AccessKey:   "",
		SecretKey:   "",
		BulkMaxSize: 4096,
	}
}

func readConfig(cfg *config.C) (*clsConfig, error) {
	c := defaultConfig()
	if err := cfg.Unpack(&c); err != nil {
		return nil, err
	}
	return &c, nil
}

func (c *clsConfig) Validate() error {
	if c == nil {
		return errors.New("nil config")
	}
	if c.Endpoint == "" {
		return errors.New("cls.endpoint is required")
	}
	urlobj, err := url.Parse(c.Endpoint)
	if err != nil {
		return fmt.Errorf("invalid cls.endpoint, %w", err)
	}
	if urlobj.Scheme != "http" && urlobj.Scheme != "https" {
		return errors.New("invalid cls.endpoint, it's scheme must be http or https")
	}
	if c.Topic == "" {
		return errors.New("cls.topic is required")
	}
	if c.AccessKey == "" {
		return errors.New("cls.access_key is required")
	}
	if c.SecretKey == "" {
		return errors.New("cls.secret_key is required")
	}
	if c.BulkMaxSize <= 0 {
		return errors.New("cls.bulk_max_size must be greater than zero")
	}
	return nil
}