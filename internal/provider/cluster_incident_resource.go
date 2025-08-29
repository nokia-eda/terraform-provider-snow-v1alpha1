package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/nokia/eda/apps/terraform-provider-snow/internal/eda/apiclient"
	"github.com/nokia/eda/apps/terraform-provider-snow/internal/resource_cluster_incident"
	"github.com/nokia/eda/apps/terraform-provider-snow/internal/tfutils"
)

const (
	create_rs_clusterIncident = "/apps/snow.eda.nokia.com/v1alpha1/clusterincidents"
	read_rs_clusterIncident   = "/apps/snow.eda.nokia.com/v1alpha1/clusterincidents/{name}"
	update_rs_clusterIncident = "/apps/snow.eda.nokia.com/v1alpha1/clusterincidents/{name}"
	delete_rs_clusterIncident = "/apps/snow.eda.nokia.com/v1alpha1/clusterincidents/{name}"
)

var (
	_ resource.Resource                = (*clusterIncidentResource)(nil)
	_ resource.ResourceWithConfigure   = (*clusterIncidentResource)(nil)
	_ resource.ResourceWithImportState = (*clusterIncidentResource)(nil)
)

func NewClusterIncidentResource() resource.Resource {
	return &clusterIncidentResource{}
}

type clusterIncidentResource struct {
	client *apiclient.EdaApiClient
}

func (r *clusterIncidentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster_incident"
}

func (r *clusterIncidentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_cluster_incident.ClusterIncidentResourceSchema(ctx)
}

func (r *clusterIncidentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resource_cluster_incident.ClusterIncidentModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize unknown values with null defaults
	err := tfutils.FillMissingValues(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error filling missing values", err.Error())
		return
	}

	// Convert Terraform model to API request body
	reqBody, err := tfutils.ModelToAnyMap(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error building request", err.Error())
		return
	}

	// Create API call logic
	tflog.Info(ctx, "Create()::API request", map[string]any{
		"path": create_rs_clusterIncident,
		"body": spew.Sdump(reqBody),
	})

	t0 := time.Now()
	result := map[string]any{}

	err = r.client.Create(ctx, create_rs_clusterIncident, nil, reqBody, &result)

	tflog.Info(ctx, "Create()::API returned", map[string]any{
		"path":      create_rs_clusterIncident,
		"result":    spew.Sdump(result),
		"timeTaken": time.Since(t0).String(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error creating resource", err.Error())
		return
	}

	// Read the resource again to populate any values not available in the response from Create()
	t0 = time.Now()

	err = r.client.Get(ctx, read_rs_clusterIncident, map[string]string{
		"name": tfutils.StringValue(data.Metadata.Name),
	}, &result)

	tflog.Info(ctx, "Read()::API returned", map[string]any{
		"path":      read_rs_clusterIncident,
		"result":    spew.Sdump(result),
		"timeTaken": time.Since(t0).String(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
		return
	}

	// Convert API response to Terraform model
	err = tfutils.AnyMapToModel(ctx, result, &data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build response from API result", err.Error())
		return
	}
	// Save created data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)
}

func (r *clusterIncidentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resource_cluster_incident.ClusterIncidentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	tflog.Info(ctx, "Read()::API request", map[string]any{
		"path": read_rs_clusterIncident,
		"data": spew.Sdump(data),
	})

	t0 := time.Now()
	result := map[string]any{}

	err := r.client.Get(ctx, read_rs_clusterIncident, map[string]string{
		"name": tfutils.StringValue(data.Metadata.Name),
	}, &result)

	tflog.Info(ctx, "Read()::API returned", map[string]any{
		"path":      read_rs_clusterIncident,
		"result":    spew.Sdump(result),
		"timeTaken": time.Since(t0).String(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
		return
	}

	// Convert API response to Terraform model
	err = tfutils.AnyMapToModel(ctx, result, &data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build response from API result", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterIncidentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data resource_cluster_incident.ClusterIncidentModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := tfutils.FillMissingValues(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error filling missing values", err.Error())
		return
	}

	reqBody, err := tfutils.ModelToAnyMap(ctx, &data)
	if err != nil {
		resp.Diagnostics.AddError("Error building request", err.Error())
		return
	}

	// Update API call logic
	tflog.Info(ctx, "Update()::API request", map[string]any{
		"path": update_rs_clusterIncident,
		"body": spew.Sdump(reqBody),
	})

	t0 := time.Now()
	result := map[string]any{}

	err = r.client.Update(ctx, update_rs_clusterIncident, map[string]string{
		"name": tfutils.StringValue(data.Metadata.Name),
	}, reqBody, &result)

	tflog.Info(ctx, "Update()::API returned", map[string]any{
		"path":      update_rs_clusterIncident,
		"result":    spew.Sdump(result),
		"timeTaken": time.Since(t0).String(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error updating resource", err.Error())
		return
	}

	// Read the resource again to populate any values not available in the response from Update()
	t0 = time.Now()

	err = r.client.Get(ctx, read_rs_clusterIncident, map[string]string{
		"name": tfutils.StringValue(data.Metadata.Name),
	}, &result)

	tflog.Info(ctx, "Read()::API returned", map[string]any{
		"path":      read_rs_clusterIncident,
		"result":    spew.Sdump(result),
		"timeTaken": time.Since(t0).String(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error reading resource", err.Error())
		return
	}

	// Convert API response to Terraform model
	err = tfutils.AnyMapToModel(ctx, result, &data)
	if err != nil {
		resp.Diagnostics.AddError("Failed to build response from API result", err.Error())
		return
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *clusterIncidentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resource_cluster_incident.ClusterIncidentModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete API call logic
	tflog.Info(ctx, "Delete()::API request", map[string]any{
		"path": delete_rs_clusterIncident,
		"data": spew.Sdump(data),
	})

	t0 := time.Now()
	result := map[string]any{}

	err := r.client.Delete(ctx, delete_rs_clusterIncident, map[string]string{
		"name": tfutils.StringValue(data.Metadata.Name),
	}, &result)

	tflog.Info(ctx, "Delete()::API returned", map[string]any{
		"path":      delete_rs_clusterIncident,
		"result":    spew.Sdump(result),
		"timeTaken": time.Since(t0).String(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Error deleting resource", err.Error())
		return
	}
}

// Configure adds the provider configured client to the resource.
func (r *clusterIncidentResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Add a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.EdaApiClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *api.EdaApiClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}
	r.client = client
}

// ImportState implements resource.ResourceWithImportState.
func (r *clusterIncidentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.Split(req.ID, "/")
	if len(parts) < 1 {
		resp.Diagnostics.AddError("Invalid ID", fmt.Sprintf("Expected <namespace/name> format, got: %s", req.ID))
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("metadata").AtName("name"), parts[0])...)
}
