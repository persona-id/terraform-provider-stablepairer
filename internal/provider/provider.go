// Copyright (c) Persona
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure StablePairer satisfies various provider interfaces.
var _ provider.Provider = &StablePairer{}

// StablePairer defines the provider implementation.
type StablePairer struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &StablePairer{
			version: version,
		}
	}
}

func (p *StablePairer) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *StablePairer) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *StablePairer) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "stablepairer"
	resp.Version = p.version
}

func (p *StablePairer) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewPairResource,
	}
}

func (p *StablePairer) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "This Terraform provider provides a resource that keeps a stable mapping between a set of keys and a set of values to prevent churn like using module / division indexing.",
	}
}
