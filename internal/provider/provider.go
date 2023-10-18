// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/client"
)

// Ensure ApolloProvider satisfies various provider interfaces.
var _ provider.Provider = &ApolloProvider{}

// ApolloProvider defines the provider implementation.
type ApolloProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ApolloProviderModel describes the provider data model. - Reflects the schema
type ApolloProviderModel struct {
	PersonalApiKey types.String `tfsdk:"personal_api_key"`
}

func (p *ApolloProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "apollo"
	resp.Version = p.version
}

func (p *ApolloProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"personal_api_key": schema.StringAttribute{
				MarkdownDescription: "User's personal Apollo API key",
				Required:            true,
			},
		},
	}
}

func (p *ApolloProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ApolloProviderModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }
	if data.PersonalApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_api_key"),
			"Unknown apollo api key",
			"The provider cannot connect to the Apollo client because there is an unknown configuration value for the personal api key.",
		)
	}

	if data.PersonalApiKey.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("personal_api_key"),
			"Missing apollo api key",
			"The provider cannot connect to the Apollo client because there is a missing configuration value for the personal api key.",
		)
	}

	Apollo := &client.Client{
		ApiKey: data.PersonalApiKey.ValueString(),
	}

	resp.DataSourceData = Apollo
	resp.ResourceData = Apollo
}

func (p *ApolloProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewGraphResource,
	}
}

func (p *ApolloProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ApolloProvider{
			version: version,
		}
	}
}

