package rule_test

import (
	"fmt"
	"testing"

	"github.com/elastic/terraform-provider-elasticstack/internal/acctest"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceRule(t *testing.T) {
	indexName := sdkacctest.RandStringFromCharSet(22, sdkacctest.CharSetAlphaNum)
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { acctest.PreCheck(t) },
		CheckDestroy:      checkResourceRuleDestroy,
		ProviderFactories: acctest.Providers,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceIndexCreate(indexName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test", "name", indexName),
					resource.TestCheckResourceAttr("elasticstack_elasticsearch_index.test", "schedule.interval.#", "1m"),
				),
			},
		},
	})
}

func testAccResourceIndexCreate(name string) string {
	return fmt.Sprintf(`
provider "elasticstack" {
  elasticsearch {}
}

resource "elasticstack_elasticsearch_rule" "test" {
  name = "%s"
	id = "first_alert"
	consumer = "alerts"
	notify_when = "onActionGroupChange"
	rule_type_id = ".index-threshold"
	
	schedule = jsonencode({
		interval: "1m"
	})

	params = jsonencode({
		aggType: "avg",
		termSize: 6,
		thresholdComparator: ">",
		timeWindowSize: 5,
		timeWindowUnit: "m",
		groupBy: "top",
		threshold: [ 1000 ],
		index: [ ".test-index" ],
		timeField: "@timestamp",
		aggField: "sheet.version",
		termField: "name.keyword"
	})
}
	`, name)
}

func checkResourceRuleDestroy(s *terraform.State) error {
	return nil
	// client := acctest.Provider.Meta().(*clients.ApiClient)

	// for _, rs := range s.RootModule().Resources {
	// 	if rs.Type != "elasticstack_elasticsearch_index" {
	// 		continue
	// 	}
	// 	compId, _ := clients.CompositeIdFromStr(rs.Primary.ID)

	// 	res, err := client.GetESClient().Indices.Get([]string{compId.ResourceId})
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if res.StatusCode != 404 {
	// 		return fmt.Errorf("Index (%s) still exists", compId.ResourceId)
	// 	}
	// }
	// return nil
}
