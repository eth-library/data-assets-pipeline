package env

import (
	"encoding/json"
	"runtime"
	"strings"

	"github.com/eth-library/dap/cli/internal/exec"
	"github.com/eth-library/dap/cli/internal/ui"
	"github.com/spf13/cobra"
)

var versionsAll bool

// VersionsCmd displays tool versions.
var VersionsCmd = &cobra.Command{
	Use:     "versions",
	Short:   "Display tool versions",
	Long:    "Show versions of development tools. Use --all to show all Nix flake dependencies.",
	GroupID: GroupID,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Subtitle("Versions")
		if versionsAll {
			ShowVersionsFull()
		} else {
			ShowVersionsCompact()
		}
	},
}

func init() {
	VersionsCmd.Flags().BoolVarP(&versionsAll, "all", "a", false, "Show all Nix flake dependencies")
}

// ShowVersionsCompact displays essential tool versions for the welcome screen.
func ShowVersionsCompact() {
	showVersion("python", getPythonVersion())
	showVersion("uv", getUVVersion())
	showVersion("dagster", getDagsterVersion())
}

// ShowVersionsFull displays all Nix flake dependencies.
func ShowVersionsFull() {
	showVersion("python", getPythonVersion())
	showVersion("uv", getUVVersion())
	showVersion("go", runtime.Version())
	showVersion("kubectl", getKubectlVersion())
	showVersion("helm", getHelmVersion())
	showVersion("jq", getJqVersion())
	showVersion("curl", getCurlVersion())
	showVersion("openssl", getOpensslVersion())
}

func showVersion(name, version string) {
	if version != "" {
		ui.KeyValue(name, version)
	} else {
		ui.KeyValue(name, ui.Styles.Dim.Render("not found"))
	}
}

func getPythonVersion() string {
	out, err := exec.Run("python", "--version")
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(out, "Python ")
}

func getUVVersion() string {
	out, err := exec.Run("uv", "--version")
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(out, "uv ")
}

func getDagsterVersion() string {
	out, err := exec.Run("python", "-c", "import dagster; print(dagster.__version__)")
	if err != nil {
		return ""
	}
	return out
}

func getRuffVersion() string {
	out, err := exec.Run("ruff", "--version")
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(out, "ruff ")
}

func getMypyVersion() string {
	out, err := exec.Run("mypy", "--version")
	if err != nil {
		return ""
	}
	// Output is "mypy X.Y.Z (compiled: yes)"
	parts := strings.Fields(out)
	if len(parts) >= 2 {
		return parts[1]
	}
	return out
}

func getKubectlVersion() string {
	out, err := exec.Run("kubectl", "version", "--client", "-o", "json")
	if err != nil {
		return ""
	}

	var v struct {
		ClientVersion struct {
			GitVersion string `json:"gitVersion"`
		} `json:"clientVersion"`
	}
	if err := json.Unmarshal([]byte(out), &v); err != nil {
		return ""
	}
	return v.ClientVersion.GitVersion
}

func getHelmVersion() string {
	out, err := exec.Run("helm", "version", "--short")
	if err != nil {
		return ""
	}
	if idx := strings.Index(out, "+"); idx > 0 {
		out = out[:idx]
	}
	return out
}

func getJqVersion() string {
	out, err := exec.Run("jq", "--version")
	if err != nil {
		return ""
	}
	return strings.TrimPrefix(out, "jq-")
}

func getCurlVersion() string {
	out, err := exec.Run("curl", "--version")
	if err != nil {
		return ""
	}
	// First line is "curl X.Y.Z ..."
	if idx := strings.Index(out, " "); idx > 0 {
		rest := out[idx+1:]
		if end := strings.Index(rest, " "); end > 0 {
			return rest[:end]
		}
	}
	return ""
}

func getOpensslVersion() string {
	out, err := exec.Run("openssl", "version")
	if err != nil {
		return ""
	}
	// Output is "OpenSSL X.Y.Z ..."
	parts := strings.Fields(out)
	if len(parts) >= 2 {
		return parts[1]
	}
	return ""
}
