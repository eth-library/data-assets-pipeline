package exec

import (
	"testing"
)

func TestWhich(t *testing.T) {
	tests := []struct {
		name     string
		command  string
		expected bool
	}{
		{"go exists", "go", true},
		{"nonexistent command", "definitely-not-a-real-command-12345", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Which(tt.command)
			if got != tt.expected {
				t.Errorf("Which(%q) = %v, want %v", tt.command, got, tt.expected)
			}
		})
	}
}

func TestRun(t *testing.T) {
	tests := []struct {
		name        string
		command     string
		args        []string
		wantErr     bool
		wantContain string
	}{
		{
			name:        "echo command",
			command:     "echo",
			args:        []string{"hello"},
			wantErr:     false,
			wantContain: "hello",
		},
		{
			name:    "nonexistent command",
			command: "definitely-not-a-real-command-12345",
			args:    []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Run(tt.command, tt.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.wantContain != "" {
				if got != tt.wantContain {
					t.Errorf("Run() = %q, want to contain %q", got, tt.wantContain)
				}
			}
		})
	}
}
