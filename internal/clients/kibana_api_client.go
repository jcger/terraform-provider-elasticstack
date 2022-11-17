package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

type ApiError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func (apiError *ApiError) Error() string {
	return apiError.Message
}

func NewKibanaApiClient(d *schema.ResourceData, meta interface{}) (*KibanaApiClient, error) {
	u, _ := url.Parse("http://127.0.0.1:5601")
	cfg := elastictransport.Config{
		URLs: []*url.URL{u},
	}
	transport, err := elastictransport.New(cfg)
	return &KibanaApiClient{transport}, err
}

func (t *KibanaApiClient) PostKibanaRule(ctx context.Context, rule *models.Rule) diag.Diagnostics {
	var diags diag.Diagnostics

	// this should be implemented outside this provider

	// creates the post request body
	var body = make(map[string]interface{})
	var params = make(map[string]interface{})
	var schedule = make(map[string]string)

	// id is ignored as it will be loaded from the response
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
	body["actions"] = rule.Actions

	// creates a JSON based on the body struct
	reqBodyJSON, err := json.Marshal(body)
	if err != nil {
		return diag.FromErr(err)
	}

	fmt.Printf("\nReqBodyJSON:\n%v\n\n", bytes.NewBuffer(reqBodyJSON))

	// creates the request
	req, _ := http.NewRequest("POST", "/api/alerting/rule/", bytes.NewBuffer(reqBodyJSON))

	// auth should be loaded from somewhere else
	req.SetBasicAuth("elastic", "changeme")
	req.Header.Add("Content-Type", "application/json")
	// required by the API
	req.Header.Add("kbn-xsrf", "true")

	// does the actual request
	res, err := t.est.Perform(req)
	if err != nil {
		return diag.FromErr(err)
	}

	// loads the body from the response as the id is in there
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	// if the request doesn't return a 200
	if res.StatusCode != http.StatusOK {
		if err != nil {
			return diag.FromErr(err)
		}

		var apiError ApiError
		if err := json.Unmarshal([]byte(responseBody), &apiError); err != nil {
			return diag.FromErr(err)
		}

		// we load the message prop to show it to the user
		return diag.FromErr(fmt.Errorf("api error response message \"%s\"", apiError.Message))
	}

	// fmt.Printf("\nResponseBody:\n%s\n\n", string(responseBody))
	// saves the response body into an auxiliar rule to retrieve the id
	var auxRule models.Rule
	if err := json.Unmarshal([]byte(responseBody), &auxRule); err != nil {
		return diag.FromErr(err)
	}

	// fmt.Printf("\nauxRule:\n%+v\n\n", auxRule)
	// assigns the id to the rule as the rule is passed by as it's passed by as reference
	// the value will be saved outside this function
	rule.Id = auxRule.Id

	return diags
}

func (t *KibanaApiClient) PutKibanaRule(ctx context.Context, rule *models.Rule) diag.Diagnostics {
	var diags diag.Diagnostics

	// creates the post request body
	var body = make(map[string]interface{})
	var params = make(map[string]interface{})
	var schedule = make(map[string]string)

	// id is ignored as it will be loaded from the response
	body["name"] = rule.Name
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

	reqBodyJSON, err := json.Marshal(body)
	if err != nil {
		return diag.FromErr(err)
	}
	fmt.Printf("\nReqBodyJSON:\n%v\n\nRule.Id:\n%s\n\n", bytes.NewBuffer(reqBodyJSON), rule.Id)

	path := fmt.Sprintf("/api/alerting/rule/%s", rule.Id)
	req, _ := http.NewRequest("PUT", path, bytes.NewBuffer(reqBodyJSON))

	req.SetBasicAuth("elastic", "changeme")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("kbn-xsrf", "true")

	res, err := t.est.Perform(req)
	if err != nil {
		return diag.FromErr(err)
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			return diag.FromErr(err)
		}

		var apiError ApiError
		if err := json.Unmarshal([]byte(bodyBytes), &apiError); err != nil {
			return diag.FromErr(err)
		}

		fmt.Printf("%#v", apiError)
		return diag.FromErr(fmt.Errorf("api error response message \"%s\"", apiError.Message))
	}

	return diags
}

func (t *KibanaApiClient) PostKibanaIndexConnector(ctx context.Context, indexConnector *models.IndexConnector) diag.Diagnostics {
	var diags diag.Diagnostics

	// this should be implemented outside this provider

	// creates the post request body
	var body = make(map[string]interface{})
	var config = make(map[string]interface{})

	// id is ignored as it will be loaded from the response

	config["index"] = indexConnector.Config.Index
	body["name"] = indexConnector.Name
	body["connector_type_id"] = indexConnector.ConnectorTypeId
	body["config"] = config

	// creates a JSON based on the body struct
	reqBodyJSON, err := json.Marshal(body)
	if err != nil {
		return diag.FromErr(err)
	}

	// fmt.Printf("\nReqBodyJSON:\n%v\n\n", bytes.NewBuffer(reqBodyJSON))

	// creates the request
	req, _ := http.NewRequest("POST", "/api/actions/connector", bytes.NewBuffer(reqBodyJSON))

	// auth should be loaded from somewhere else
	req.SetBasicAuth("elastic", "changeme")
	req.Header.Add("Content-Type", "application/json")
	// required by the API
	req.Header.Add("kbn-xsrf", "true")

	// does the actual request
	res, err := t.est.Perform(req)
	if err != nil {
		return diag.FromErr(err)
	}

	// loads the body from the response as the id is in there
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return diag.FromErr(err)
	}
	defer res.Body.Close()

	// if the request doesn't return a 200
	if res.StatusCode != http.StatusOK {
		if err != nil {
			return diag.FromErr(err)
		}

		var apiError ApiError
		if err := json.Unmarshal([]byte(responseBody), &apiError); err != nil {
			return diag.FromErr(err)
		}

		// we load the message prop to show it to the user
		return diag.FromErr(fmt.Errorf("api error response message \"%s\"", apiError.Message))
	}

	// fmt.Printf("\nResponseBody:\n%s\n\n", string(responseBody))
	// saves the response body into an auxiliar index connector to retrieve the id
	var auxIndexConnector models.IndexConnector
	if err := json.Unmarshal([]byte(responseBody), &auxIndexConnector); err != nil {
		return diag.FromErr(err)
	}

	// fmt.Printf("\nauxIndexConnector:\n%+v\n\n", auxIndexConnector)
	// assigns the id to the index connector as the index connector is passed by as it's passed by as reference
	// the value will be saved outside this function
	indexConnector.Id = auxIndexConnector.Id

	return diags
}
