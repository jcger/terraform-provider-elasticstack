package rule

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

func ResourceRule() *schema.Resource {
	rulesSchema := map[string]*schema.Schema{
		// "id": {
		// 	Description: "Internal identifier of the resource",
		// 	Type:        schema.TypeString,
		// 	Required:    true,
		// },
		"name": {
			Description: "A name to reference and search",
			Type:        schema.TypeString,
			Required:    true,
		},
		"consumer": {
			Description: "The name of the application or feature that owns the rule",
			Type:        schema.TypeString,
			Required:    true,
		},
		"notify_when": {
			Description: "The condition for throttling the notification",
			Type:        schema.TypeString,
			Required:    true,
		},
		"rule_type_id": {
			Description: "The ID of the rule type that you want to call when the rule is scheduled to run",
			Type:        schema.TypeString,
			Required:    true,
		},
		// "schedule": {
		// 	Description: "The check interval, which specifies how frequently the rule conditions are checked",
		// 	Type:        schema.TypeMap,
		// 	Required:    true,
		// 	Elem: &schema.Resource{
		// 		Schema: map[string]*schema.Schema{
		// 			"interval": {
		// 				Type:     schema.TypeString,
		// 				Required: true,
		// 			},
		// 		},
		// 	},
		// },
		"schedule": {
			Description:      "The check interval, which specifies how frequently the rule conditions are checked",
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: utils.DiffJsonSuppress,
			ValidateFunc:     validation.StringIsJSON,
			Default:          "{}",
		},
		"params": {
			Description:      "Rule parameters",
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: utils.DiffJsonSuppress,
			ValidateFunc:     validation.StringIsJSON,
			Default:          "{}",
		},
		// params depend on each rule type
		// we don't have a specific definition
		// "params": {
		// 	Description: "Rule parameters",
		// 	Required:    true,
		// 	Type:        schema.TypeMap,

		// 	Elem: &schema.Resource{
		// 		Schema: map[string]*schema.Schema{
		// 			"agg_type": {
		// 				Type:     schema.TypeString,
		// 				Required: true,
		// 			},
		// 			"term_size": {
		// 				Type:     schema.TypeInt,
		// 				Required: true,
		// 			},
		// 			"threshold_comparator": {
		// 				Type:     schema.TypeString,
		// 				Required: true,
		// 			},
		// 			"time_window_size": {
		// 				Type:     schema.TypeInt,
		// 				Required: true,
		// 			},
		// 			"time_window_unit": {
		// 				Type:     schema.TypeString,
		// 				Required: true,
		// 			},
		// 			"group_by": {
		// 				Type:     schema.TypeString,
		// 				Required: true,
		// 			},
		// 			"threshold": {
		// 				Type: schema.TypeSet,
		// 				Elem: &schema.Schema{
		// 					Type: schema.TypeInt,
		// 				},
		// 				Required: true,
		// 			},
		// 			"index": {
		// 				Type: schema.TypeSet,
		// 				Elem: &schema.Schema{
		// 					Type: schema.TypeString,
		// 				},
		// 				Required: true,
		// 			},
		// 			"time_field": {
		// 				Type:     schema.TypeString,
		// 				Required: true,
		// 			},
		// 			"agg_field": {
		// 				Type:     schema.TypeString,
		// 				Required: true,
		// 			},
		// 			"term_field": {
		// 				Type:     schema.TypeString,
		// 				Required: true,
		// 			},
		// 		},
		// 	},
		// },
	}
	utils.AddConnectionSchema(rulesSchema)
	return &schema.Resource{
		Description:   "Creates Elasticsearch rules. See https://www.elastic.co/guide/en/kibana/current/create-rule-api.html",
		CreateContext: resourceRuleCreate,
		ReadContext:   resourceRuleRead,
		UpdateContext: resourceRuleUpdate,
		DeleteContext: resourceRuleDelete,
		// Importer: &schema.ResourceImporter{
		// 	StateContext: func(ctx context.Context, d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
		// 		// first populate what we can with Read
		// 		diags := resourceRuleRead(ctx, d, m)
		// 		if diags.HasError() {
		// 			return nil, fmt.Errorf("Unable to import requested index")
		// 		}
		// 		return diags
		// 	},
		// },
		// CustomizeDiff: customdiff.ForceNewIfChange("mappings", func(ctx context.Context, old, new, meta interface{}) bool {
		// 	fmt.Println("************************************* calling CustomizeDiff")

		// 	return false
		// }),
		Schema: rulesSchema,
	}
}

// The provider uses the Create method to create a new resource based on the schema data.
// https://learn.hashicorp.com/tutorials/terraform/providers-plugin-framework-resource-create?in=terraform/providers-plugin-framework
func resourceRuleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Info(ctx, "********************************* resourceRuleCreate")
	var diags diag.Diagnostics
	client, err := clients.NewKibanaApiClient(d, meta)
	if err != nil {
		return diag.FromErr(err)
	}
	var rule models.Rule

	rule.Name = d.Get("name").(string)
	rule.Consumer = d.Get("consumer").(string)
	rule.NotifyWhen = d.Get("notify_when").(string)
	rule.RuleTypeId = d.Get("rule_type_id").(string)

	if v, ok := d.GetOk("schedule"); ok {
		var schedule models.AlertRuleSchedule
		if v.(string) != "" {
			if err := json.Unmarshal([]byte(v.(string)), &schedule); err != nil {
				return diag.FromErr(err)
			}
		}
		rule.Schedule = schedule
	}

	if v, ok := d.GetOk("params"); ok {
		var params models.AlertRuleParams
		if v.(string) != "" {
			if err := json.Unmarshal([]byte(v.(string)), &params); err != nil {
				return diag.FromErr(err)
			}
		}
		rule.Params = params
	}

	if diags := client.PutKibanaRule(ctx, &rule); diags.HasError() {
		return diags
	}

	// d.SetId(id.String())
	return diags
}

func resourceRuleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Info(ctx, "********************************* resourceRuleRead")

	var diags diag.Diagnostics
	return diags
}

func resourceRuleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Debug(ctx, "********************************* resourceRuleDelete")

	var diags diag.Diagnostics
	return diags
}

func resourceRuleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	tflog.Warn(ctx, "********************************* resourceRuleUpdate")

	var diags diag.Diagnostics
	return diags
}
