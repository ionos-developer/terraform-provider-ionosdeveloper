package ionosdeveloper

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviderFactories = map[string]func() (*schema.Provider, error){
	"ionosdeveloper": func() (*schema.Provider, error) {
		return Provider(), nil
	},
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv(apiKeyEnvVar) == "" {
		t.Fatalf("%s must be set for acceptance tests", apiKeyEnvVar)
	}
}
