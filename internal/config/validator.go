package config

import (
	"errors"
	"net/url"
)

func Validate(cfg *Config) error {

	if cfg.Server.Port <= 0 {
		return errors.New("invalid server port")
	}

	if len(cfg.Upstreams) == 0 {
		return errors.New("no upstreams configured")
	}

	for _, up := range cfg.Upstreams {

		_, err := url.Parse(up.URL)
		if err != nil {
			return err
		}

		if up.Weight <= 0 {
			return errors.New("weight must be > 0")
		}

		if up.HealthCheckPath == "" {
			return errors.New("health check path required")
		}
	}

	return nil
}