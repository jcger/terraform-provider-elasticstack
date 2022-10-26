package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/terraform-provider-elasticstack/internal/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type KibanaApiClient struct {
	est *elastictransport.Client
	// version string
}

func NewKibanaApiClient(d *schema.ResourceData, meta interface{}) (*KibanaApiClient, error) {
	u, _ := url.Parse("http://127.0.0.1:5601")
	cfg := elastictransport.Config{
		URLs: []*url.URL{u},
	}
	transport, err := elastictransport.New(cfg)
	return &KibanaApiClient{transport}, err
}

func (t *KibanaApiClient) PutKibanaRule(ctx context.Context, rule *models.Rule) diag.Diagnostics {
	var diags diag.Diagnostics

	// type RequestBodyStruct struct {
	// 	name         string
	// 	consumer     string
	// 	rule_type_id string
	// 	notify_when  string
	// 	schedule     map[string]interface{}
	// 	params       map[string]interface{}
	// }

	// var reqBody RequestBodyStruct

	// reqBody.consumer = rule.Consumer
	// reqBody.name = rule.Name
	// reqBody.rule_type_id = rule.RuleTypeId
	// reqBody.notify_when = rule.NotifyWhen
	// reqBody.schedule = rule.Schedule
	// reqBody.params = rule.Params

	fmt.Printf("++++++++++++++++++++++ %#v", rule)

	var body = make(map[string]interface{})
	var params = make(map[string]interface{})
	var schedule = make(map[string]string)

	body["consumer"] = rule.Consumer
	body["name"] = rule.Name
	body["rule_type_id"] = rule.RuleTypeId
	body["notify_when"] = rule.NotifyWhen
	schedule["interval"] = rule.Schedule.Interval

	params["aggType"] = rule.Params.AggType
	params["termSize"] = rule.Params.TermSize
	params["thresholdComparator"] = rule.Params.ThresholdComparator
	params["timeWindowSize"] = rule.Params.TimeWindowSize
	params["timeWindowUnit"] = rule.Params.TimeWindowUnit
	params["groupBy"] = rule.Params.GroupBy
	params["threshold"] = rule.Params.Threshold
	params["index"] = rule.Params.Index
	params["timeField"] = rule.Params.TimeField
	params["aggField"] = rule.Params.AggField
	params["termField"] = rule.Params.TermField

	body["params"] = params
	body["schedule"] = schedule

	reqBodyJSON, err := json.Marshal(rule)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Println("-----------------------")
	fmt.Printf("%v", bytes.NewBuffer(reqBodyJSON))
	fmt.Println("-----------------------")

	req, _ := http.NewRequest("POST", "/api/alerting/rule/", bytes.NewBuffer(reqBodyJSON))

	req.SetBasicAuth("elastic", "changeme")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("kbn-xsrf", "true")

	res, err := t.est.Perform(req)
	fmt.Printf("called it! %#v", res)
	if err != nil {
		diag.FromErr(err)
	}
	defer res.Body.Close()

	return diags
}

// func main() {
// 	u, _ := url.Parse("http://127.0.0.1:9200")

// 	cfg := elastictransport.Config{
// 		URLs: []*url.URL{u},
// 	}
// 	transport, err := elastictransport.New(cfg)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	req, _ := http.NewRequest("GET", "/", nil)

// 	res, err := transport.Perform(req)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer res.Body.Close()

// 	log.Println(res)
// }
