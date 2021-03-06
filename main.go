package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	ionosdeveloper "github.com/ionos-developer/terraform-provider-ionosdeveloper/ionosdeveloper"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return ionosdeveloper.Provider()
		},
	})
}
