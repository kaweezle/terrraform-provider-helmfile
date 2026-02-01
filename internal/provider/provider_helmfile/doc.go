// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// Package provider_helmfile contains the core provider implementation.
//
// This package defines the HelmfileProvider type which manages the lifecycle
// of helmfile operations within the Terraform provider context. It includes:
//
//   - Provider configuration and initialization
//   - Global options management
//   - Helmfile executor integration
//   - Plugin management
//
// The provider acts as a bridge between Terraform's resource management
// system and the helmfile execution layer.
package provider_helmfile
