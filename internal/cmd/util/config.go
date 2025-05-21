package util

import (
	"bytes"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

const (
	envAddress      = "TODOIST_ADDRESS"
	envApiToken     = "TODOIST_API_TOKEN"
	envApiTokenFile = "TODOIST_API_TOKEN_FILE"
)

func (f *Factory) loadConfigFromFile() error {
	data, err := os.ReadFile(f.ConfigPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}

		buf := &bytes.Buffer{}
		encoder := toml.NewEncoder(buf)
		encoder.Indent = ""
		if err := encoder.Encode(f); err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(f.ConfigPath), 0o755); err != nil {
			return err
		}

		return os.WriteFile(f.ConfigPath, buf.Bytes(), 0o644)
	}

	return toml.Unmarshal(data, f)
}

func (f *Factory) loadConfigFromEnv() error {
	address := os.Getenv(envAddress)
	if address != "" {
		f.DeamonConfig.Address = address
	}

	apiToken := os.Getenv(envApiToken)
	if apiToken != "" {
		f.DeamonConfig.ApiToken = apiToken
	}

	apiTokenFile := os.Getenv(envApiTokenFile)
	if apiTokenFile != "" {
		f.DeamonConfig.ApiTokenFile = apiTokenFile
	}

	return nil
}

func (f *Factory) loadApiToken() error {
	if f.DeamonConfig.ApiToken != "" {
		return nil
	}

	if f.DeamonConfig.ApiTokenFile != "" {
		token, err := os.ReadFile(f.DeamonConfig.ApiTokenFile)
		if err != nil {
			return err
		}
		f.DeamonConfig.ApiToken = string(token)
	}

	return nil
}

func (f *Factory) LoadConfig() error {
	if err := f.loadConfigFromFile(); err != nil {
		return err
	}

	if err := f.loadConfigFromEnv(); err != nil {
		return err
	}

	if err := f.loadApiToken(); err != nil {
		return err
	}

	return nil
}
