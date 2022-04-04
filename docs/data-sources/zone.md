# Data Source: ionosdeveloper_dns_zone

`ionosdeveloper_dns_zone` provides details about a specific zone managed by Ionos DNS.

## Example usage

The following example shows how to create a data source for a zone and how to use it to create a record:

```hcl
data "ionosdeveloper_dns_zone" "selected" {
 name = "example.com"
}

resource "ionosdeveloper_dns_record" "cname" {
 zone_id = data.ionosdeveloper_dns_zone.selected.id
 name    = "cname.${data.ionosdeveloper_dns_zone.selected.name}"
 type    = "CNAME"
 content = "www.cname.com"
 ttl     = 3600
}
```

## Argument Reference

The following arguments are required:

- `name` - The name of the zone.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

- `id` - The ID of the zone.
