# Resource: ionosdeveloper_dns_record

Provides a DNS record.

## Example usage

The following example shows how to create a record:

```hcl
resource "ionosdeveloper_dns_record" "example" {
  zone_id = "${data.ionosdeveloper_dns_zone.selected.id}"
  name    = "www.${data.ionosdeveloper_dns_zone.selected.name}"
  type    = "CNAME"
  content = "www.cname.com"
  ttl     = 3600
}
```

## Argument Reference

The following arguments are required:

- `zone_id` - The ID of the zone that contains the record.
- `name` - The DNS record name. Must be absolute. No trailing dot needed.
- `type` - The DNS record type. Valid values are `A`,` AAAA`,` CNAME`, `MX`, `NS`, `SOA`, `SRV`, `TXT` and `CAA`.
- `content` - The string data for the record whose meaning depends on the DNS type. For `MX` and `SRV` records, it must be set to the exchange field of the record content.
- `ttl` - The time-to-live of this record (seconds).

The following arguments are optional:

- `prio` - The preference field of the record data for MX and SRV records.
- `disabled` - If `false`, not visible in DNS.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` - The ID of the record.
