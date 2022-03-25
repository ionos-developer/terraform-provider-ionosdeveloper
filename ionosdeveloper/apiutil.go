package ionosdeveloper

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	sdk "github.com/ionos-developer/dns-sdk-go"
)

func getIndentedBody(err error) string {
	if e, ok := interface{}(err).(*sdk.GenericOpenAPIError); ok {
		var unmarshalled []interface{}
		if json.Unmarshal(e.Body(), &unmarshalled) == nil {
			if body, jsonErr := json.MarshalIndent(unmarshalled, "", "    "); jsonErr == nil {
				return string(body)
			}
		}
		return string(e.Body())
	}
	return ""
}

func appendError(diags diag.Diagnostics, summary string, err error) diag.Diagnostics {
	return append(diags, diag.Diagnostic{
		Severity: diag.Error,
		Summary:  summary,
		Detail:   fmt.Sprintf("%v\n%v\n", err, getIndentedBody(err)),
	})
}
