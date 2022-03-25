package ionosdeveloper

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var testZoneName = os.Getenv("TEST_DNS_ZONE_NAME")
var testZoneId = os.Getenv("TEST_DNS_ZONE_ID")

func testDnsAccPreCheck(t *testing.T) {
	testAccPreCheck(t)

	if os.Getenv(apiKeyEnvVar) == "" {
		t.Fatalf("%s must be set for acceptance tests", apiKeyEnvVar)
	}
}

func TestAccZoneOk(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testDnsAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: zoneConfig(testZoneName),
				Check:  resource.TestCheckResourceAttr("data.ionosdeveloper_dns_zone.z", "id", testZoneId),
			},
			{
				Config:      zoneConfig("inexistent-zone.de"),
				ExpectError: regexp.MustCompile("DNS zone does not exist"),
			},
			{
				Config:      zoneConfig("-1"),
				ExpectError: regexp.MustCompile("DNS zone does not exist"),
			},
		},
	})
}

func zoneConfig(name string) string {
	return `
data ionosdeveloper_dns_zone z {
  name = "` + name + `"
}`
}
