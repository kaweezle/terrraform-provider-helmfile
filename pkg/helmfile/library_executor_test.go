// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// cSpell: words notworking stretchr

package helmfile_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/helmfile/helmfile/pkg/config"
	"github.com/kaweezle/terraform-provider-helmfile/pkg/helmfile"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// resolveKubeconfig returns the path to the active kubeconfig file, or an
// empty string when none can be located.
func resolveKubeconfig() string {
	if kc := os.Getenv("KUBECONFIG"); kc != "" {
		// KUBECONFIG may be a colon-separated list; take the first entry.
		return strings.SplitN(kc, ":", 2)[0]
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	defaultPath := filepath.Join(home, ".kube", "config")
	if _, err := os.Stat(defaultPath); err == nil {
		return defaultPath
	}
	return ""
}

// TestHelmfileLibraryExecutor_Apply_MultiRelease runs helmfile apply against a
// helmfile with two releases:
//
//   - working:    a simple ConfigMap – expected to be applied successfully.
//   - notworking: a kubernetes.io/tls Secret with invalid key names (crt / key
//     instead of the required tls.crt / tls.key) – expected to be rejected by
//     the Kubernetes API server.
//
// notworking declares `needs: [test/working]` so helmfile always applies
// working first. The overall Apply call is expected to return a non-nil error
// because notworking fails, while working should have been applied
// successfully.
//
// This is an integration test. It is skipped unless:
//  1. The TEST_INTEGRATION environment variable is set to "true", and
//  2. A kubeconfig is available (KUBECONFIG env var or ~/.kube/config).
func TestHelmfileLibraryExecutor_Apply_MultiRelease(t *testing.T) {
	t.Parallel()
	if os.Getenv("TEST_INTEGRATION") != "true" {
		t.Skip("skipping integration test: set TEST_INTEGRATION=true to run")
	}

	kubeconfig := resolveKubeconfig()
	if kubeconfig == "" {
		t.Skip("skipping integration test: no kubeconfig found")
	}

	// During `go test`, the working directory is the package directory, so
	// testdata/ is accessible via a relative path.
	helmfilePath := filepath.Join("testdata", "multi-release", "helmfile.yaml")
	if _, err := os.Stat(helmfilePath); err != nil {
		t.Fatalf("testdata not found at %s: %v", helmfilePath, err)
	}

	// Build a minimal executor: no extra plugins, no extra global options
	// beyond what the defaults provide.
	globalOpts := &helmfile.GlobalOptions{}
	executor := helmfile.NewHelmfileLibraryExecutor(globalOpts, nil)

	// Configure resource-level options.
	opts := &helmfile.Options{}
	opts.WithFileOrDir(helmfilePath)
	opts.WithKubeconfig(kubeconfig)

	applyOptions := &config.ApplyOptions{
		// Skip the pre-apply diff so the test focuses on the apply itself.
		SkipDiffOnInstall: true,
		// Do not suppress secrets output so errors from notworking are visible.
		SuppressSecrets: false,
	}

	stdout, logs, err := executor.Apply(context.Background(), opts, applyOptions)

	combined := stdout + "\n" + logs

	// The overall Apply must fail because the notworking release is rejected
	// by Kubernetes (invalid key names in a kubernetes.io/tls Secret).
	require.Error(t, err, "Apply should return an error when a release fails:\n%s", combined)

	// The working release was processed before notworking (enforced by the
	// needs: directive), so it should appear somewhere in the output.
	assert.Contains(t, combined, "working",
		"output should reference the working release")

	// The notworking release should also appear in the output so the caller
	// can identify which release failed.
	assert.Contains(t, combined, "notworking",
		"output should reference the notworking release")
}
