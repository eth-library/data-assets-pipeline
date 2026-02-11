package k8s

import (
	"testing"
)

func TestConstants(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"Namespace", Namespace, "dagster"},
		{"Release", Release, "dagster"},
		{"Image", Image, "da-pipeline:local"},
		{"HelmChart", HelmChart, "dagster/dagster"},
		{"K8sContext", K8sContext, "docker-desktop"},
	}

	for _, tt := range tests {
		if tt.value != tt.want {
			t.Errorf("%s = %q, want %q", tt.name, tt.value, tt.want)
		}
	}
}

func TestK8sCmdHasSubcommands(t *testing.T) {
	if !K8sCmd.HasSubCommands() {
		t.Error("K8sCmd should have subcommands")
	}
}

func TestK8sCmdUse(t *testing.T) {
	if K8sCmd.Use != "k8s" {
		t.Errorf("K8sCmd.Use = %q, want %q", K8sCmd.Use, "k8s")
	}
}
