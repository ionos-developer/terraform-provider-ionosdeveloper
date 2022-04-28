# IonosDeveloper Provider

Use the IONOS Developer provider to manage your DNS configuration.

Use the navigation to the left to read about the available data sources and resources.

## Authentication and Configuration

The provider needs to be configured with proper credentials before it can be used.

One way to configure the provider is to set the `IONOS_API_KEY` environment variable as shown in the below example:

```hcl
$ export IONOS_API_KEY="x-api-key"
```

Another way of configuring it, is by providing your credentials in the .tf configuration file.
You can either explicitly write them or use `var.name`:

```hcl
provider "ionosdeveloper" {
    api_key = var.api_key
}
```

**Important notes**

- If you use var.name, the environment variables must be in the format `TF_VAR_name` and this will be checked last for a value. For example:

```hcl
$ export TF_VAR_api_key="x-api-key"
```

## Configuration Reference

The following arguments are required:

- `api_key` - If omitted, the IONOS_API_KEY environment variable is used.

## Example usage

```hcl
terraform {
  required_providers {
    ionosdeveloper = {
      source  = "ionos-developer/ionosdeveloper"
      version = ">= 0.1"
    }
  }
}

# Configure the IonosDeveloper Provider
provider "ionosdeveloper" {
  api_key = "X-API-Key"
}

# Create a DNS record
resource "ionosdeveloper_dns_record" "example" {
#...
}
```

**Important notes**

- The `required_providers` section must be specified in order for Terraform to be able to find and download the ionosdeveloper provider.
- The `credentials` provided in a .tf file will override the credentials from environment variables.

## Debugging

Setting up the environment variable `IONOS_DEBUG` will enable logging the HTTP traffic with the customer API in the DNS SDK. This alone will not display any logs in the console.
To see the logs in the console, a Terraform environment variable must be set as shown in the below example:

```hcl
$ export TF_LOG=info
$ export IONOS_DEBUG=1
$ terraform apply
```

**Important notes**

- We recommend you only use IONOS_DEBUG for debugging purposes. Disable it in your production environments because it can log sensitive data.
