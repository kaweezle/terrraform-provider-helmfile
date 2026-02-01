// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// Package provider implements the Terraform Provider for Helmfile.
//
// This package contains the main provider implementation that integrates
// helmfile operations with Terraform's resource lifecycle. It includes:
//
//   - Provider registration and configuration
//   - Resource definitions (helmfile_release)
//   - Data source definitions (if any)
//   - Schema definitions and validation
//
// The provider enables declarative management of Helm releases through
// Helmfile configurations within Terraform infrastructure code.
package provider
