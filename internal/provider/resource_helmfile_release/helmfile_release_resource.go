// Copyright Antoine Martin 2026
// SPDX-License-Identifier: MIT

package resource_helmfile_release

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/kaweezle/terraform-provider-helmfile/internal/provider/provider_helmfile"
)

var _ resource.ResourceWithConfigure = &HelmfileReleaseResource{}

type HelmfileReleaseResource struct {
	provider *provider_helmfile.HelmfileProvider
}

func NewHelmfileReleaseResource() resource.Resource {
	return &HelmfileReleaseResource{}
}

func (r *HelmfileReleaseResource) Metadata(
	_ context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_release"
}

func (r *HelmfileReleaseResource) Schema(
	ctx context.Context,
	_ resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = HelmfileReleaseResourceSchema(ctx)
	resp.Schema.Description = "Manages a Helm release defined in a Helmfile."
}

func (r *HelmfileReleaseResource) Configure(
	_ context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}
	provider, ok := req.ProviderData.(*provider_helmfile.HelmfileProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf(
				"Expected *provider_helmfile.HelmfileProvider, got: %T. Please report this issue to the provider developers.",
				req.ProviderData,
			),
		)

		return
	}
	r.provider = provider
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	var data HelmfileReleaseModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: Implement the logic to create the Helmfile release using r.provider
	// Steps (To be completed):
	// 1. Make a build --embed-values call to helmfile with the relevant helmfile and environment and
	//    create a sha256 hash of the output
	// 2. Make a list in order to list releases
	// 3. Make an apply call to helmfile with the relevant helmfile and environment to ensure the release is
	//    created/updated
	// 4. Store the hash in the state to detect future changes

	// Set the state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	var data HelmfileReleaseModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: Implement the logic to create the Helmfile release using r.provider

	// Set the state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
	var data HelmfileReleaseModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// TODO: Implement the logic to create the Helmfile release using r.provider

	// Set the state
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

//nolint:gocritic // Interface implementation
func (r *HelmfileReleaseResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
	// Implementation of Delete operation
	var data HelmfileReleaseModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
