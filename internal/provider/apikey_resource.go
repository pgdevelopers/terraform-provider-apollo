// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/client"
	"github.com/machinebox/graphql"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &ExampleResource{}
var _ resource.ResourceWithImportState = &ExampleResource{}

func NewApiKeyResource() resource.Resource {
	return &ApiKeyResource{}
}

// ApiKeyResource defines the resource implementation.
type ApiKeyResource struct {
	client *client.Client
}

type ApiKeyResponse struct {
	Service struct {
		NewKey struct {
			KeyName string `json:"keyName"`
			KeyId   string `json:"id"`
			Token   string `json:"token"`
		} `json:"newKey"`
	} `json:"service"`
}

// ApiKeyResourceModel describes the resource data model.
type ApiKeyResourceModel struct {
	GraphId types.String `tfsdk:"graph_id"`
	KeyName types.String `tfsdk:"key_name"`
	//Role    types.String `tfsdk:"role"`
	Token types.String `tfsdk:"token"`
}

func (r *ApiKeyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_apikey"
}

func (r *ApiKeyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "API key resource",

		Attributes: map[string]schema.Attribute{
			"graph_id": schema.StringAttribute{
				MarkdownDescription: "the id of the graph that the key is for",
				Required:            true,
			},
			"key_name": schema.StringAttribute{
				MarkdownDescription: "the name of the api key",
				Required:            true,
			},
			// "role": schema.StringAttribute{
			// 	MarkdownDescription: "the role assigned to the key",
			// 	Required:            true,
			// },
			"token": schema.StringAttribute{
				MarkdownDescription: "the token of the key",
				Computed:            true,
			},
		},
	}
}

func (r *ApiKeyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ApiKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApiKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//Initialize Apollo client
	apollo := client.Client{
		ApiKey:            r.client.ApiKey,
		EnterPriseEnabled: false,
		GraphClient:       graphql.NewClient("https://graphql.api.apollographql.com/api/graphql"),
	}
	apollo.Init()
	graphId := data.GraphId.ValueString()
	keyName := data.KeyName.ValueString()
	//role := data.Role.ValueString()

	query := `
	mutation Service {
		service(id: "` + graphId + `") {
			newKey(keyName: "` + keyName + `") {
				keyName
				id
				token
			}
		}
	}
`

	var result ApiKeyResponse
	err := apollo.Query(ctx, query, result)
	if err != nil {
		resp.Diagnostics.AddError("create api key error", fmt.Sprintf("Unable to create api key, got error: %s", err))
		return
	}
	//resp.Diagnostics.AddError("check it", fmt.Sprintf("check it: %s", result))
	data.Token = basetypes.NewStringValue(result.Service.NewKey.Token)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created an apikey")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApiKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ApiKey, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApiKeyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update ApiKey, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApiKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApiKeyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ApiKey, got error: %s", err))
	//     return
	// }
}

func (r *ApiKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
