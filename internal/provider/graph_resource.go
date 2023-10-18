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
	"github.com/hashicorp/terraform-provider-scaffolding-framework/internal/helpers"
	"github.com/machinebox/graphql"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &GraphResource{}
var _ resource.ResourceWithImportState = &GraphResource{}

func NewGraphResource() resource.Resource {
	return &GraphResource{}
}

type graph struct {
	NewService struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Title string `json:"title"`
	} `json:"newService"`
}

// GraphResource defines the resource implementation.
type GraphResource struct {
	client *client.Client
}

// GraphResourceModel describes the resource data model.
type GraphResourceModel struct {
	OrgId     types.String `tfsdk:"org_id"`
	GraphName types.String `tfsdk:"graph_name"`
	GraphId   types.String `tfsdk:"graph_id"`
}

func (r *GraphResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_graph"
}

func (r *GraphResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Graph resource",

		Attributes: map[string]schema.Attribute{
			"org_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID for Apollo Studio",
				Required:            true,
			},
			"graph_name": schema.StringAttribute{
				MarkdownDescription: "Name of your graph",
				Required:            true,
			},
			"graph_id": schema.StringAttribute{
				MarkdownDescription: "ID of your graph",
				Computed:            true,
			},
		},
	}
}

func (r *GraphResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *GraphResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GraphResourceModel

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
	var result graph

	orgId := data.OrgId.ValueString()
	id := data.GraphName.ValueString() + helpers.RandomNumberString(5)
	data.GraphId = basetypes.NewStringValue(id)
	name := data.GraphName.ValueString()
	adminOnly := "false"

	query := `
		mutation Service {
			newService(accountId: "` + orgId + `", id:  "` + id + `", name: "` + name + `", title: "` + name + `", hiddenFromUninvitedNonAdminAccountMembers: ` + adminOnly + `) {
				id
				name
				title
			}
	  	}
		`

	err := apollo.Query(ctx, query, result)
	if err != nil {
		resp.Diagnostics.AddError("create graph error", fmt.Sprintf("Unable to create graph, got error: %s", err))
		return
	}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "created a resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GraphResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GraphResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read Graph, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GraphResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GraphResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	// httpResp, err := r.client.Do(httpReq)
	// if err != nil {
	//     resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update Graph, got error: %s", err))
	//     return
	// }

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GraphResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GraphResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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

	var result graph
	graphId := data.GraphId.ValueString()

	err := apollo.Query(ctx, `

		mutation Service
		{
			service( id:  "`+graphId+`") {
		 		delete
			}
	  }
		`,
		&result)

	if err != nil {
		resp.Diagnostics.AddError("delete graph error", fmt.Sprintf("Unable to delete graph, got error: %s", err))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *GraphResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
