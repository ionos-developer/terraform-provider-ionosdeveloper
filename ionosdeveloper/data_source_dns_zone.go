package ionosdeveloper

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDnsZone() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDnsZoneRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		}}
}

func dataSourceDnsZoneRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(SdkBundle).DnsApiClient
	var diags diag.Diagnostics

	resp, _, err := c.ZonesApi.GetZones(context.Background()).Execute()
	if err != nil {
		return appendError(diags, "Unable to get DNS zone", err)
	}

	zoneName := d.Get("name")
	for _, zone := range resp {
		if *zone.Name == zoneName {
			d.SetId(*zone.Id)
			return diags
		}
	}

	return append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  "DNS zone does not exist",
	})
}
