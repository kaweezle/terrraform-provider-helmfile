// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package utils

// cSpell: words validatordiag
import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = (*ExistingFileOrDirectoryValidator)(nil)

type ExistingFileOrDirectoryValidator struct {
	AllowDirectory bool
}

func (v ExistingFileOrDirectoryValidator) Description(_ context.Context) string {
	if v.AllowDirectory {
		return "value must be an existing file or directory path"
	}
	return "value must be an existing file path"
}

func (v ExistingFileOrDirectoryValidator) MarkdownDescription(ctx context.Context) string {
	return v.Description(ctx)
}

//nolint:gocritic // validator interface
func (v ExistingFileOrDirectoryValidator) ValidateString(
	_ context.Context,
	request validator.StringRequest,
	response *validator.StringResponse,
) {
	if request.ConfigValue.IsNull() || request.ConfigValue.IsUnknown() {
		return
	}
	path := request.ConfigValue.ValueString()
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"File or Directory does not exist",
			fmt.Sprintf("The path '%s' does not exist.", path),
		))
		return
	}
	if err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"Error accessing path",
			fmt.Sprintf("An error occurred while accessing the path '%s': %s", path, err.Error()),
		))
		return
	}
	if info.IsDir() && !v.AllowDirectory {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(
			request.Path,
			"Expected a file but found a directory",
			fmt.Sprintf("The path '%s' is a directory, but a file was expected.", path),
		))
		return
	}
}
