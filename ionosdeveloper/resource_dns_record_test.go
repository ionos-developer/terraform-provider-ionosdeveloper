//go:build all || dns

package ionosdeveloper

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDnsRecord_Validations(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      invalidTtl,
				ExpectError: regexp.MustCompile("ttl"),
			},
			{
				Config:      invalidType,
				ExpectError: regexp.MustCompile("invalid"),
			},
			{
				Config:      invalidPrio,
				ExpectError: regexp.MustCompile("prio"),
			},
		},
	})
}

func TestAccDnsRecord_MX(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: mx,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "name", "test-acc."+testZoneName),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "type", "MX"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "content", "a.de"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "prio", "0"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "ttl", "1000"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "disabled", "false"),
				),
			},
			{
				Config: mx2,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "name", "test-acc."+testZoneName),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "type", "MX"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "content", "new.de"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "prio", "20"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "ttl", "2000"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "disabled", "true"),
				),
			},
		},
	})
}

func TestAccDnsRecord_TXT(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: txt,
				Check:  resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "content", "\"text\""),
			},
		},
	})
}

func TestAccDnsRecord_UpdateType(t *testing.T) {
	var initialId string
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: mx,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "name", "test-acc."+testZoneName),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "type", "MX"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "content", "a.de"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "prio", "0"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "ttl", "1000"),
					getCurrentId("ionosdeveloper_dns_record.r", &initialId),
				),
			},
			{
				Config: txt,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkDifferentId("ionosdeveloper_dns_record.r", &initialId),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "name", "test-acc."+testZoneName),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "type", "TXT"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "content", "\"text\""),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "ttl", "100"),
				),
			},
		},
	})
}

func TestAccDnsRecord_UpdateName(t *testing.T) {
	var initialId string
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: mx,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "name", "test-acc."+testZoneName),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "type", "MX"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "content", "a.de"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "ttl", "1000"),
					getCurrentId("ionosdeveloper_dns_record.r", &initialId),
				),
			},
			{
				Config: mxUpdateName,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkDifferentId("ionosdeveloper_dns_record.r", &initialId),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "name", "test-acc2."+testZoneName),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "type", "MX"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "content", "a.de"),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "ttl", "1000"),
				),
			},
		},
	})
}

func TestAccDnsRecord_NormalizedName(t *testing.T) {
	var initialId string
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: punycodeName,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "name", "ss.test-acc."+testZoneName),
					getCurrentId("ionosdeveloper_dns_record.r", &initialId),
				),
			},
			{
				Config: unicodeName,
				Check: resource.ComposeAggregateTestCheckFunc(
					checkSameId("ionosdeveloper_dns_record.r", &initialId),
					resource.TestCheckResourceAttr("ionosdeveloper_dns_record.r", "name", "ss.test-acc."+testZoneName),
				),
			},
		},
	})
}

func getCurrentId(n string, id *string) resource.TestCheckFunc {
	return resource.TestCheckFunc(func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		*id = rs.Primary.ID

		return nil
	})
}

func checkDifferentId(n string, id *string) resource.TestCheckFunc {
	return resource.TestCheckFunc(func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == *id {
			return fmt.Errorf("The resource id didn't change")
		}

		return nil
	})
}

func checkSameId(n string, id *string) resource.TestCheckFunc {
	return resource.TestCheckFunc(func(s *terraform.State) error {
		// find the corresponding state object
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID != *id {
			return fmt.Errorf("The resource id changed")
		}

		return nil
	})
}

var mx = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "test-acc.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "mx"
  content  = "a.de"
  ttl      = 1000
  disabled = false
}`

var mx2 = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "test-acc.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "MX"
  content  = "new.de"
  prio     = "20"
  ttl      = 2000
  disabled = true
}`

var mxUpdateName = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "test-acc2.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "mx"
  content  = "a.de"
  ttl      = 1000
  disabled = false
}`

var txt = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "test-acc.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "TXT"
  content  = "text"
  ttl      = 100
}`

var a = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "test-acc.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "A"
  content  = "1.1.1.1"
  ttl      = 100
}`

var invalidTtl = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "test-acc.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "MX"
  content  = "a.de"
  ttl      = 1
}`

var invalidType = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "test-acc.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "invalid"
  content  = "a.de"
  ttl      = 1000
}`

var invalidPrio = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "test-acc.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "mx"
  content  = "a.de"
  ttl      = 1000
  prio     = 65536
}`

var punycodeName = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "ss.test-acc.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "a"
  content  = "1.1.1.1"
  ttl      = 1000
}`

var unicodeName = zoneConfig(testZoneName) + `
resource ionosdeveloper_dns_record r {
  zone_id  = data.ionosdeveloper_dns_zone.z.id
  name     = "ß.test-acc.${data.ionosdeveloper_dns_zone.z.name}"
  type     = "a"
  content  = "1.1.1.1"
  ttl      = 1000
}`
