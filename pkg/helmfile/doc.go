// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// Package helmfile provides an abstraction layer for executing helmfile operations.
//
// This package defines interfaces and implementations for interacting with helmfile,
// a declarative spec for deploying Helm charts. It includes:
//
//   - HelmfileExecutor interface: Defines the contract for executing helmfile commands
//   - HelmfileLibraryExecutor: Implementation that uses the helmfile library directly
//   - Options types: Configuration structures for controlling helmfile behavior
//   - Output capture: Utilities for capturing and logging helmfile output
//   - TFLogEncoder: Integration with Terraform's logging system
//
// The package is designed to work with the Terraform Plugin Framework and provides
// seamless integration between helmfile operations and Terraform resource management.
package helmfile
