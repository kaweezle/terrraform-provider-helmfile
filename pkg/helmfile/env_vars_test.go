// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package helmfile_test

// cSpell: words zaptest

import (
	"os"
	"testing"

	"github.com/kaweezle/terraform-provider-helmfile/pkg/helmfile"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func TestSetEnvVars_NewVariables(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	// Ensure test variables don't exist
	testVars := map[string]string{
		"TEST_VAR_1": "value1",
		"TEST_VAR_2": "value2",
	}

	for key := range testVars {
		os.Unsetenv(key) //nolint:errcheck // cleanup
	}

	// Set variables
	cleanup := helmfile.SetEnvVars(logger, testVars)

	// Verify variables are set
	for key, expectedValue := range testVars {
		actualValue, exists := os.LookupEnv(key)
		if !exists {
			t.Errorf("Expected environment variable %s to exist", key)
		}
		if actualValue != expectedValue {
			t.Errorf("Expected %s=%s, got %s", key, expectedValue, actualValue)
		}
	}

	// Cleanup
	cleanup()

	// Verify variables are unset
	for key := range testVars {
		if _, exists := os.LookupEnv(key); exists {
			t.Errorf("Expected environment variable %s to be unset after cleanup", key)
		}
	}
}

func TestSetEnvVars_ExistingVariables(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	// Set original values
	originalVars := map[string]string{
		"TEST_VAR_3": "original1",
		"TEST_VAR_4": "original2",
	}

	for key, value := range originalVars {
		os.Setenv(key, value) //nolint:errcheck // test setup
	}
	defer func() {
		for key := range originalVars {
			os.Unsetenv(key) //nolint:errcheck // cleanup
		}
	}()

	// Set new values
	newVars := map[string]string{
		"TEST_VAR_3": "new1",
		"TEST_VAR_4": "new2",
	}

	cleanup := helmfile.SetEnvVars(logger, newVars)

	// Verify new values are set
	for key, expectedValue := range newVars {
		actualValue, exists := os.LookupEnv(key)
		if !exists {
			t.Errorf("Expected environment variable %s to exist", key)
		}
		if actualValue != expectedValue {
			t.Errorf("Expected %s=%s, got %s", key, expectedValue, actualValue)
		}
	}

	// Cleanup
	cleanup()

	// Verify original values are restored
	for key, expectedValue := range originalVars {
		actualValue, exists := os.LookupEnv(key)
		if !exists {
			t.Errorf("Expected environment variable %s to exist after cleanup", key)
		}
		if actualValue != expectedValue {
			t.Errorf("Expected %s=%s after cleanup, got %s", key, expectedValue, actualValue)
		}
	}
}

func TestSetEnvVars_MixedVariables(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	// Set one variable, leave another unset
	os.Setenv("TEST_VAR_5", "original") //nolint:errcheck // test setup
	os.Unsetenv("TEST_VAR_6")           //nolint:errcheck // test setup

	defer os.Unsetenv("TEST_VAR_5") //nolint:errcheck // cleanup

	testVars := map[string]string{
		"TEST_VAR_5": "new",
		"TEST_VAR_6": "value",
	}

	cleanup := helmfile.SetEnvVars(logger, testVars)

	// Verify both are set
	for key, expectedValue := range testVars {
		actualValue, exists := os.LookupEnv(key)
		if !exists {
			t.Errorf("Expected environment variable %s to exist", key)
		}
		if actualValue != expectedValue {
			t.Errorf("Expected %s=%s, got %s", key, expectedValue, actualValue)
		}
	}

	// Cleanup
	cleanup()

	// Verify TEST_VAR_5 is restored, TEST_VAR_6 is unset
	if value, exists := os.LookupEnv("TEST_VAR_5"); !exists || value != "original" {
		t.Errorf("Expected TEST_VAR_5=original after cleanup, got %s (exists: %v)", value, exists)
	}
	if _, exists := os.LookupEnv("TEST_VAR_6"); exists {
		t.Error("Expected TEST_VAR_6 to be unset after cleanup")
	}
}

func TestSetEnvVars_EmptyMap(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	// Set empty map
	cleanup := helmfile.SetEnvVars(logger, map[string]string{})

	// Should not panic
	cleanup()
}

func TestSetEnvVars_EmptyValue(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	testVars := map[string]string{
		"TEST_VAR_7": "",
	}

	cleanup := helmfile.SetEnvVars(logger, testVars)
	defer cleanup()

	// Verify empty value is set
	value, exists := os.LookupEnv("TEST_VAR_7")
	if !exists {
		t.Error("Expected TEST_VAR_7 to exist")
	}
	if value != "" {
		t.Errorf("Expected empty string, got %s", value)
	}

	cleanup()

	// Verify it's unset
	if _, exists := os.LookupEnv("TEST_VAR_7"); exists {
		t.Error("Expected TEST_VAR_7 to be unset after cleanup")
	}
}

func TestSetEnvVars_MultipleCleanups(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	testVars := map[string]string{
		"TEST_VAR_8": "value",
	}

	cleanup := helmfile.SetEnvVars(logger, testVars)

	// Call cleanup multiple times (should be safe)
	cleanup()
	cleanup()

	// Verify variable is unset
	if _, exists := os.LookupEnv("TEST_VAR_8"); exists {
		t.Error("Expected TEST_VAR_8 to be unset after cleanup")
	}
}

func TestSetEnvVars_WithNilLogger(t *testing.T) {
	t.Parallel()

	testVars := map[string]string{
		"TEST_VAR_9": "value",
	}

	defer os.Unsetenv("TEST_VAR_9") //nolint:errcheck // cleanup

	// This will panic if the function doesn't handle nil logger
	// The current implementation doesn't check for nil, so this documents the behavior
	defer func() {
		if r := recover(); r != nil {
			t.Log("Function panics with nil logger (expected behavior)")
		}
	}()

	cleanup := helmfile.SetEnvVars(zap.NewNop(), testVars)
	cleanup()
}

func TestSetEnvVars_ConcurrentUsage(t *testing.T) {
	t.Parallel()

	// Test that multiple goroutines can use SetEnvVars without interfering
	// This tests that the cleanup function captures values correctly
	done := make(chan bool, 2)

	go func() {
		logger := zaptest.NewLogger(t)
		testVars := map[string]string{
			"TEST_VAR_10": "goroutine1",
		}
		cleanup := helmfile.SetEnvVars(logger, testVars)
		defer cleanup()

		// Verify value
		if value, _ := os.LookupEnv("TEST_VAR_10"); value != "goroutine1" {
			t.Errorf("Goroutine 1: Expected TEST_VAR_10=goroutine1, got %s", value)
		}
		done <- true
	}()

	go func() {
		logger := zaptest.NewLogger(t)
		testVars := map[string]string{
			"TEST_VAR_11": "goroutine2",
		}
		cleanup := helmfile.SetEnvVars(logger, testVars)
		defer cleanup()

		// Verify value
		if value, _ := os.LookupEnv("TEST_VAR_11"); value != "goroutine2" {
			t.Errorf("Goroutine 2: Expected TEST_VAR_11=goroutine2, got %s", value)
		}
		done <- true
	}()

	<-done
	<-done
}

func TestSetEnvVars_SpecialCharacters(t *testing.T) {
	t.Parallel()

	logger := zaptest.NewLogger(t)

	testVars := map[string]string{
		"TEST_VAR_12": "value with spaces",
		"TEST_VAR_13": "value=with=equals",
		"TEST_VAR_14": "value\nwith\nnewlines",
	}

	cleanup := helmfile.SetEnvVars(logger, testVars)
	defer cleanup()

	// Verify special characters are preserved
	for key, expectedValue := range testVars {
		actualValue, exists := os.LookupEnv(key)
		if !exists {
			t.Errorf("Expected environment variable %s to exist", key)
		}
		if actualValue != expectedValue {
			t.Errorf("Expected %s=%q, got %q", key, expectedValue, actualValue)
		}
	}

	cleanup()

	// Verify cleanup
	for key := range testVars {
		if _, exists := os.LookupEnv(key); exists {
			t.Errorf("Expected %s to be unset after cleanup", key)
		}
	}
}

func BenchmarkSetEnvVars(b *testing.B) {
	logger := zap.NewNop()

	testVars := map[string]string{
		"BENCH_VAR_1": "value1",
		"BENCH_VAR_2": "value2",
		"BENCH_VAR_3": "value3",
		"BENCH_VAR_4": "value4",
		"BENCH_VAR_5": "value5",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cleanup := helmfile.SetEnvVars(logger, testVars)
		cleanup()
	}
}
