---
layout: "alicloud"
page_title: "Alicloud: alicloud_ram_account_aliases"
sidebar_current: "docs-alicloud-datasource-ram-account-alias"
description: |-
    Provides an alias of the Alibaba Cloud account.
---

# alicloud\_ram\_account\_aliases

This data source provides an alias for the Alibaba Cloud account.

## Example

```
data "alicloud_ram_account_aliases" "alias_ds" {
  output_file = "alias.txt"
}

output "account_alias" {
  value = "${data.alicloud_ram_account_aliases.alias_ds.account_alias}"
}
```

## Argument Reference

The following arguments are supported:

* `output_file` - (Optional) Set the name of the file where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `account_alias` - Alias of the account.