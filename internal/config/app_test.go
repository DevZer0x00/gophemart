package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAppOptions(t *testing.T) {
	tests := []struct {
		Name         string
		Environments []string
		Arguments    []string
		Config       *Config
		HasError     bool
	}{
		{
			Name:         "default values",
			Environments: []string{},
			Arguments:    []string{},
			Config: &Config{
				Address:              ":8080",
				DatabaseUri:          "user=user password=password dbname=pg sslmode=disable host=localhost port=5433",
				AccrualSystemAddress: "localhost",
			},
			HasError: false,
		},
		{
			Name:         "cli options",
			Environments: []string{},
			Arguments: []string{
				"-a",
				"8",
				"-d",
				"localhost",
				"-r",
				"9",
			},
			Config: &Config{
				Address:              "8",
				DatabaseUri:          "localhost",
				AccrualSystemAddress: "9",
			},
		},
		{
			Name: "env override",
			Environments: []string{
				"RUN_ADDRESS=127.0.0.1",
				"DATABASE_URI=2",
				"ACCURAL_SYSTEM_ADDRESS=3",
			},
			Arguments: []string{
				"-a",
				"8",
				"-d",
				"localhost",
				"-r",
				"9",
			},
			Config: &Config{
				Address:              "127.0.0.1",
				DatabaseUri:          "2",
				AccrualSystemAddress: "3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			cfg, err := ParseAppOptions(tt.Environments, tt.Arguments)
			require.Equal(t, tt.HasError, err != nil, err)
			assert.Equal(t, tt.Config, cfg)
		})
	}
}
