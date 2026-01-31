// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package resource_helmfile_release

// cSpell: words planmodifier basetypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func BuildDigestModifier(computedResource string) planmodifier.String {
	return buildDigestModifier{
		computedResource: computedResource,
	}
}

// buildDigestModifier implements the plan modifier.
type buildDigestModifier struct {
	computedResource string
}

// Description returns a human-readable description of the plan modifier.
func (m buildDigestModifier) Description(_ context.Context) string {
	return "Set the planned value to the digest of the Helmfile build output."
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (m buildDigestModifier) MarkdownDescription(_ context.Context) string {
	return "Set the planned value to the digest of the Helmfile build output."
}

// PlanModifyString implements the plan modification logic.
//
//nolint:gocritic // False positive: hugeParam (implementation of interface)
func (m buildDigestModifier) PlanModifyString(
	ctx context.Context,
	req planmodifier.StringRequest,
	resp *planmodifier.StringResponse,
) {
	tflog.Debug(ctx, "##BDM## Trying modifier", map[string]any{
		"path":              req.Path.String(),
		"computed_resource": m.computedResource,
		"config_value":      req.ConfigValue.ValueString(),
		"plan_value":        req.PlanValue.ValueString(),
		"state_value":       req.StateValue.ValueString(),
	})

	if !req.ConfigValue.IsUnknown() && !req.ConfigValue.IsNull() {
		tflog.Debug(
			ctx,
			"##BDM## Config value is set, skipping build digest modifier",
			map[string]any{
				"config_value": req.ConfigValue.String(),
			},
		)
		return
	}

	// Get the provider-specific private state data.
	computedPath := req.Path.ParentPath().AtName(m.computedResource)
	var computedValue basetypes.StringValue
	resp.Diagnostics.Append(req.Plan.GetAttribute(ctx, computedPath, &computedValue)...)
	if resp.Diagnostics.HasError() {
		tflog.Debug(ctx, "##BDM## Error getting computed resource attribute", map[string]any{
			"diagnostics": resp.Diagnostics,
		})
		return
	}

	if computedValue.IsNull() || computedValue.IsUnknown() {
		tflog.Debug(
			ctx,
			"##BDM## Computed value is null or unknown, skipping build digest modifier",
			map[string]any{
				"computed_value": computedValue,
			},
		)
		return
	}

	tflog.Debug(ctx, "##BDM## Setting computed digest value", map[string]any{
		"computed_value": computedValue.ValueString(),
	})
	resp.PlanValue = computedValue
}
