package ionosdeveloper

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	sdk "github.com/ionos-developer/dns-sdk-go"
)

func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDnsRecordCreate,
		ReadContext:   resourceDnsRecordRead,
		UpdateContext: resourceDnsRecordUpdate,
		DeleteContext: resourceDnsRecordDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(old, new)
				},
				ValidateDiagFunc: func(v interface{}, p cty.Path) diag.Diagnostics {
					var diags diag.Diagnostics
					switch strings.ToUpper(v.(string)) {
					case
						"A",
						"AAAA",
						"CNAME",
						"MX",
						"NS",
						"SOA",
						"SRV",
						"TXT",
						"CAA":
						return diags
					}
					return append(diags, diag.Diagnostic{
						Severity: diag.Error,
						Summary:  "Invalid record type",
						Detail:   fmt.Sprintf("%q is not one of A, AAAA, CNAME, MX, NS, SOA, SRV, TXT, CAA", v),
					})
				},
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						return true
					}

					recordType := getRecordType(d.Get("type"))

					record := sdk.NewRecord()
					record.SetType(recordType)
					record.SetContent(new)
					normalized, _, err := sdkBundle.DnsApiClient.RecordsApi.NormalizeRecord(context.TODO()).Record(*record).Execute()

					if err != nil {
						fmt.Printf("Error getting the normalized content from API")
						return false
					}

					return normalized.GetContent() == old
				},
			},
			"ttl": {
				Type:             schema.TypeInt,
				Required:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(60)),
			},
			"prio": {
				Type:             schema.TypeInt,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 65535)),
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceDnsRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(SdkBundle).DnsApiClient
	var diags diag.Diagnostics

	zoneId := d.Get("zone_id").(string)
	recordName := d.Get("name").(string)
	recordType := getRecordType(d.Get("type"))
	recordContent := d.Get("content").(string)
	prio := d.Get("prio").(int)

	record := sdk.NewRecord()
	record.SetName(recordName)
	record.SetType(recordType)
	record.SetContent(recordContent)
	if d.Get("ttl") != 0 {
		record.SetTtl(int32(d.Get("ttl").(int)))
	}
	if recordType == sdk.MX || recordType == sdk.SRV {
		record.SetPrio(int32(prio))
	} else if prio != 0 {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Only MX and SRV records can set prio",
		})
	}
	record.SetDisabled(d.Get("disabled").(bool))

	createdRecords, _, err := c.RecordsApi.CreateRecords(context.Background(), zoneId).Record([]sdk.Record{*record}).Execute()
	if err != nil {
		return appendError(diags, "Unable to create zone record", err)
	}

	updateRecord(d, &createdRecords[0])

	return diags
}

func resourceDnsRecordRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(SdkBundle).DnsApiClient
	var diags diag.Diagnostics

	zoneId := d.Get("zone_id").(string)
	recordId := d.Id()
	record, _, err := c.RecordsApi.GetRecord(context.Background(), zoneId, recordId).Execute()
	if err != nil {
		return appendError(diags, "Unable to read record", err)
	}

	d.Set("type", *record.Type)
	d.Set("content", *record.Content)
	d.Set("ttl", *record.Ttl)
	if record.Prio != nil {
		d.Set("prio", *record.Prio)
	}
	d.Set("disabled", *record.Disabled)

	return diags
}

func resourceDnsRecordUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(SdkBundle).DnsApiClient
	var diags diag.Diagnostics

	zoneId := d.Get("zone_id").(string)
	recordId := d.Id()
	recordUpdate := *sdk.NewRecordUpdate()
	recordType := getRecordType(d.Get("type"))
	prio := d.Get("prio").(int)

	if d.HasChange("content") {
		recordUpdate.SetContent(d.Get("content").(string))
	}

	if d.HasChange("ttl") {
		recordUpdate.SetTtl(int32(d.Get("ttl").(int)))
	}

	if d.HasChange("prio") && (recordType == sdk.MX || recordType == sdk.SRV) {
		recordUpdate.SetPrio(int32(prio))
	} else if prio != 0 {
		return append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Only MX and SRV records can set prio",
		})
	}

	if d.HasChange("disabled") {
		recordUpdate.SetDisabled(d.Get("disabled").(bool))
	}

	updatedRecord, _, err := c.RecordsApi.UpdateRecord(context.Background(), zoneId, recordId).RecordUpdate(recordUpdate).Execute()
	if err != nil {
		return appendError(diags, "Unable to update record", err)
	}

	updateRecord(d, updatedRecord)

	return diags
}

func resourceDnsRecordDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(SdkBundle).DnsApiClient
	var diags diag.Diagnostics

	zoneId := d.Get("zone_id").(string)
	recordId := d.Id()

	if _, err := c.RecordsApi.DeleteRecord(context.Background(), zoneId, recordId).Execute(); err != nil {
		return appendError(diags, "Unable to delete record", err)
	}

	return diags
}

func getRecordType(value interface{}) sdk.RecordTypes {
	return sdk.RecordTypes(strings.ToUpper(value.(string)))
}

func updateRecord(d *schema.ResourceData, record *sdk.RecordResponse) {
	d.SetId(*record.Id)
	d.Set("name", *record.Name)
	d.Set("type", *record.Type)
	d.Set("content", *record.Content)
	d.Set("ttl", *record.Ttl)
	if record.Prio != nil {
		d.Set("prio", *record.Prio)
	}
	d.Set("disabled", *record.Disabled)
}
