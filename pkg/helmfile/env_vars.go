// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package helmfile

// cSpell: words nestif

import (
	"os"

	"go.uber.org/zap"
)

// SetEnvVars sets environment variables from the given map and returns a cleanup function
// that restores the original environment state. The cleanup function should be deferred.
//
// Example usage:
//
//	envVars := map[string]string{"FOO": "bar", "BAZ": "qux"}
//	cleanup := helmfile.SetEnvVars(ctx, logger, envVars)
//	defer cleanup()
func SetEnvVars(logger *zap.Logger, envVars map[string]string) func() {
	// Store original values for variables that exist
	originalValues := make(map[string]*string)

	for key := range envVars {
		if value, exists := os.LookupEnv(key); exists {
			// Create a copy of the value to avoid issues with loop variable capture
			valueCopy := value
			originalValues[key] = &valueCopy
			logger.Debug("Saving original environment variable",
				zap.String("key", key),
				zap.String("original_value", value),
			)
		} else {
			// Mark as non-existent with nil
			originalValues[key] = nil
			logger.Debug("Environment variable does not exist",
				zap.String("key", key),
			)
		}
	}

	// Set new values
	for key, value := range envVars {
		if err := os.Setenv(key, value); err != nil {
			logger.Error("Failed to set environment variable",
				zap.String("key", key),
				zap.String("value", value),
				zap.Error(err),
			)
		} else {
			logger.Debug("Set environment variable",
				zap.String("key", key),
				zap.String("value", value),
			)
		}
	}

	// Return cleanup function
	return func() {
		for key, originalValue := range originalValues {
			if originalValue != nil { //nolint:nestif // Restore original value
				// Restore original value
				if err := os.Setenv(key, *originalValue); err != nil {
					logger.Error("Failed to restore environment variable",
						zap.String("key", key),
						zap.String("original_value", *originalValue),
						zap.Error(err),
					)
				} else {
					logger.Debug("Restored environment variable",
						zap.String("key", key),
						zap.String("original_value", *originalValue),
					)
				}
			} else {
				// Unset variable (it didn't exist originally)
				if err := os.Unsetenv(key); err != nil {
					logger.Error("Failed to unset environment variable",
						zap.String("key", key),
						zap.Error(err),
					)
				} else {
					logger.Debug("Unset environment variable",
						zap.String("key", key),
					)
				}
			}
		}
	}
}
