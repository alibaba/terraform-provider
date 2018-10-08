---
layout: "alicloud"
page_title: "Alicloud: alicloud_kms_keys"
sidebar_current: "docs-alicloud-datasource-kms-keys"
description: |-
    Provides a list of available KMS Keys.
---

# alicloud\_kms\_keys

This data source provides a list of KMS keys in an Alibaba Cloud account according to the specified filters.

## Example

```
# Declare the data source
data "alicloud_kms_keys" "kms_keys_ds" {
	description_regex = "Hello KMS"
	output_file = "kms_keys.json"
}

output "first_key_id" {
  value = "${data.alicloud_kms_keys.kms_keys_ds.keys.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) Filter by KMS key ID.
* `description_regex` - (Optional) Filter the results by KMS key description with a regex string.
* `status` - (Optional) Filter the results by status of the KMS keys. Valid values: `Enabled`, `Disabled`, `PendingDeletion`.
* `output_file` - (Optional) Set the name of the file where data source results will be saved after running `terraform plan`.

## Attributes Reference

The following attributes are returned in addition to the arguments listed above:

* `keys` - A list of KMS keys. Each element contains the following attributes:
  * `id` - ID of the key.
  * `arn` - The Alibaba Cloud Resource Name (ARN) of the key.
  * `description` - Description of the key.
  * `status` - Status of the key. Possible values: `Enabled`, `Disabled` and `PendingDeletion`.
  * `creation_date` - Creation date of key.
  * `delete_date` - Deletion date of key.
  * `creator` - The owner of the key.