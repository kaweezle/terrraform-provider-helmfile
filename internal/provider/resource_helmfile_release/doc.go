// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// Package resource_helmfile_release implements the helmfile_release Terraform resource.
//
// This package provides the resource implementation for managing Helm releases
// through Helmfile configurations. It includes:
//
//   - Resource CRUD operations (Create, Read, Update, Delete)
//   - Plan modification logic for detecting changes
//   - Build digest calculation for change detection
//   - Integration with helmfile executor
//   - State management and serialization
//
// The resource allows Terraform users to declaratively manage Helm releases
// defined in Helmfile configurations, with full support for the Terraform
// resource lifecycle.
package resource_helmfile_release
