// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

// Package utils provides utility functions and validators for the Terraform provider.
//
// This package contains helper functions and custom validators used across
// the provider implementation, including:
//
//   - Custom attribute validators (file/directory existence)
//   - Common validation logic
//   - Utility functions for resource operations
//
// These utilities help ensure data integrity and provide meaningful error
// messages to Terraform users during plan and apply operations.
package utils
