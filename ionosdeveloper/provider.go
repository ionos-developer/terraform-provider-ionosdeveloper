package ionosdeveloper

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/hashicorp/terraform-plugin-sdk/meta"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	dns "github.com/ionos-developer/dns-sdk-go"
)

const apiKeyEnvVar = "IONOS_API_KEY"

type SdkBundle struct {
	DnsApiClient *dns.APIClient
}

// The provider meta is exported in order to be used by DiffSuppressFunc functions
var sdkBundle SdkBundle

func Provider() *schema.Provider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IONOS_API_URL", nil),
			},
			"auth_header": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("IONOS_AUTH_HEADER", "X-API-Key"),
			},
			"api_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc(apiKeyEnvVar, nil),
			},
		},
		ResourcesMap:   map[string]*schema.Resource{"ionosdeveloper_dns_record": resourceDnsRecord()},
		DataSourcesMap: map[string]*schema.Resource{"ionosdeveloper_dns_zone": dataSourceDnsZone()},
	}

	provider.ConfigureContextFunc = func(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
		terraformVersion := provider.TerraformVersion

		if terraformVersion == "" {
			// Terraform 0.12 introduced this field to the protocol
			// We can therefore assume that if it's missing it's 0.10 or 0.11
			terraformVersion = "0.11+compatible"
		}

		log.Printf("[DEBUG] Setting terraformVersion to %s", terraformVersion)

		tmp, diags := providerConfigure(d, terraformVersion)
		sdkBundle = tmp.(SdkBundle)
		return tmp, diags
	}

	return provider
}

func providerConfigure(d *schema.ResourceData, terraformVersion string) (interface{}, diag.Diagnostics) {
	authHeader := d.Get("auth_header").(string)
	url := d.Get("url").(string)
	apiKey := d.Get("api_key").(string)
	var diags diag.Diagnostics

	if apiKey == "" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create Ionos Developer client",
			Detail:   "api_key not provided",
		})

		return nil, diags
	}

	configuration := dns.NewConfiguration()
	if url != "" {
		configuration.Servers[0].URL = url
	}
	configuration.AddDefaultHeader(authHeader, apiKey)
	configuration.UserAgent = fmt.Sprintf(
		"terraform-provider/hashicorp-terraform/%s_terraform-plugin-sdk/%s_os/%s_arch/%s",
		terraformVersion, meta.SDKVersionString(), runtime.GOOS, runtime.GOARCH)

	if os.Getenv("IONOS_DEBUG") != "" {
		configuration.Debug = true
	}

	dnsApiClient := dns.NewAPIClient(configuration)

	return SdkBundle{
		DnsApiClient: dnsApiClient,
	}, diags
}
