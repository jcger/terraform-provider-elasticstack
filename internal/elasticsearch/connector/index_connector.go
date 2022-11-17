package indexConnector

import (
	"context"
	"encoding/json"

	"github.com/elastic/terraform-provider-elasticstack/internal/clients"
	"github.com/elastic/terraform-provider-elasticstack/internal/models"
	"github.com/elastic/terraform-provider-elasticstack/internal/utils"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func ResourceIndexConnector() *schema.Resource {
	indexConnectorSchema := map[string]*schema.Schema{
		"id": {
			Description: "Internal identifier of the resource",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"name": {
			Description: "The display name for the connector",
			Type:        schema.TypeString,
			Required:    true,
		},
		"connector_type_id": {
			Description: "The connector type ID for the connector. For example, .cases-webhook, .index, .jira, or .server-log.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"config": {
			Description:      "The configuration for the connector. Configuration properties vary depending on the connector type",
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: utils.DiffJsonSuppress,
			ValidateFunc:     validation.StringIsJSON,
			Default:          "{}",
		},
	}

	utils.AddConnectionSchema(indexConnectorSchema)
	return &schema.Resource{
		Description:   "Creates a connector. See https://www.elastic.co/guide/en/kibana/master/create-connector-api.html",
		CreateContext: resourceIndexConnectorCreate,
		ReadContext:   resourceIndexConnectorRead,
		UpdateContext: resourceIndexConnectorUpdate,
		DeleteContext: resourceIndexConnectorDelete,
		Schema:        indexConnectorSchema,
	}
}

// The provider uses the Create method to create a new resource based on the schema data.
// https://learn.hashicorp.com/tutorials/terraform/providers-plugin-framework-resource-create?in=terraform/providers-plugin-framework
func resourceIndexConnectorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Info(ctx, "********************************* resourceIndexConnectorCreate")
	var diags diag.Diagnostics
	client, err := clients.NewKibanaApiClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}

	indexConnector := models.IndexConnector{
		Name:            d.Get("name").(string),
		ConnectorTypeId: d.Get("connector_type_id").(string),
	}

	if v, ok := d.GetOk("config"); ok {
		var config models.IndexConnectorConfig
		if v.(string) != "" {
			if err := json.Unmarshal([]byte(v.(string)), &config); err != nil {
				return diag.FromErr(err)
			}
		}
		indexConnector.Config = config
	}

	if diags := client.PostKibanaIndexConnector(ctx, &indexConnector); diags.HasError() {
		return diags
	}

	// needs to be set, otherwise terraform won't realize the object was created and launch an error
	// from https://developer.hashicorp.com/terraform/tutorials/providers/provider-setup
	// "The existence of a non-blank ID tells Terraform that a resource was created. This ID can be
	// any string value, but should be a value that Terraform can use to read the resource again. Since this
	// data resource doesn't have a unique ID, you set the ID to the current UNIX time, which will force this
	// resource to refresh during every Terraform apply."
	// fmt.Printf("the rule is %s\n", rule.Id)
	d.SetId(indexConnector.Id)
	return diags
}

func resourceIndexConnectorRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Info(ctx, "********************************* resourceIndexConnectorRead")

	var diags diag.Diagnostics
	return diags
}

func resourceIndexConnectorDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "********************************* resourceIndexConnectorDelete")

	var diags diag.Diagnostics
	return diags
}

func resourceIndexConnectorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Warn(ctx, "********************************* resourceIndexConnectorUpdate")
	var diags diag.Diagnostics

	return diags
}
