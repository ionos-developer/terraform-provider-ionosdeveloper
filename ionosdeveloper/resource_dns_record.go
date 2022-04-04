package ionosdeveloper

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	dnsSdk "github.com/ionos-developer/dns-sdk-go"
)

func resourceDnsRecord() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDnsRecordCreate,
		ReadContext:   resourceDnsRecordRead,
		UpdateContext: resourceDnsRecordUpdate,
		DeleteContext: resourceDnsRecordDelete,
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.EqualFold(old, new) {
						return true
					}

					record := dnsSdk.NewRecord()
					record.SetName(new)
					record.SetType(getRecordType(d.Get("type")))
					record.SetContent(d.Get("content").(string))
					normalized, _, err := sdkBundle.DnsApiClient.RecordsApi.NormalizeRecord(context.TODO()).Record(*record).Execute()

					if err != nil {
						fmt.Printf("Error getting the normalized content from API")
						return false
					}

					return normalized.GetName() == old
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				StateFunc: func(v interface{}) string {
					value := strings.ToUpper(v.(string))
					return value
				},
				ValidateFunc: validation.StringInSlice([]string{"A", "AAAA", "CNAME", "MX", "NS", "SOA", "SRV", "TXT", "CAA"}, true),
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if strings.TrimSpace(old) == strings.TrimSpace(new) {
						return true
					}

					recordType := getRecordType(d.Get("type"))

					record := dnsSdk.NewRecord()
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
				Default:          0,
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceDnsRecordCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(SdkBundle).DnsApiClient
	var diags diag.Diagnostics

	record := createRecord(d)
	zoneId := d.Get("zone_id").(string)

	createdRecords, _, err := client.RecordsApi.CreateRecords(context.Background(), zoneId).Record([]dnsSdk.Record{*record}).Execute()
	if err != nil {
		return appendError(diags, "Unable to create zone record", err)
	}

	d.SetId(*createdRecords[0].Id)

	return resourceDnsRecordRead(ctx, d, m)
}

func createRecord(d *schema.ResourceData) *dnsSdk.Record {
	record := dnsSdk.NewRecord()

	record.SetName(d.Get("name").(string))
	record.SetType(getRecordType(d.Get("type")))
	record.SetContent(d.Get("content").(string))
	record.SetTtl(int32(d.Get("ttl").(int)))
	record.SetPrio(int32(d.Get("prio").(int)))
	record.SetDisabled(d.Get("disabled").(bool))

	return record
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

	d.Set("name", *record.Name)
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

	recordUpdate := *dnsSdk.NewRecordUpdate()

	if d.HasChange("content") {
		recordUpdate.SetContent(d.Get("content").(string))
	}

	if d.HasChange("ttl") {
		recordUpdate.SetTtl(int32(d.Get("ttl").(int)))
	}

	if d.HasChange("prio") {
		recordUpdate.SetPrio(int32(d.Get("prio").(int)))
	}

	if d.HasChange("disabled") {
		recordUpdate.SetDisabled(d.Get("disabled").(bool))
	}

	updatedRecord, _, err := c.RecordsApi.UpdateRecord(context.Background(), zoneId, recordId).RecordUpdate(recordUpdate).Execute()
	if err != nil {
		return appendError(diags, "Unable to update record", err)
	}

	d.SetId(*updatedRecord.Id)

	return resourceDnsRecordRead(ctx, d, m)
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

func getRecordType(value interface{}) dnsSdk.RecordTypes {
	return dnsSdk.RecordTypes(strings.ToUpper(value.(string)))
}
