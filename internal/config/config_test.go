package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}

	return path
}

func TestLoad(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		want      *Config
		wantError bool
	}{
		{
			name: "valid config",
			yaml: `
backends:
  - address: "127.0.0.1:8001"
    endpoints: ["api1", "api2"]
`,
			want: &Config{
				Backends: []BackendConfig{
					{
						Address:   "127.0.0.1:8001",
						Endpoints: []string{"api1", "api2"},
					},
				},
			},
			wantError: false,
		},
		{
			name: "invalid yaml",
			yaml: `
backends:
  - address: 127.0.0.1:8001
    endpoints: [api1, api2
`,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeTempConfig(t, tt.yaml)

			cfg, err := Load(path)

			if tt.wantError {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(cfg, tt.want) {
				t.Fatalf("config mismatch\n got: %+v\nwant: %+v", cfg, tt.want)
			}
		})
	}
}
