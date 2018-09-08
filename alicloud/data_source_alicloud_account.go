package alicloud

import (
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudAccount() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudAccountRead,

		Schema: map[string]*schema.Schema{
			// Computed values
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAlicloudAccountRead(d *schema.ResourceData, meta interface{}) error {
	accountId, err := meta.(*aliyunclient.AliyunClient).AccountId()

	if err != nil {
		return err
	}

	log.Printf("[DEBUG] alicloud_account - account ID found: %#v", accountId)

	d.SetId(accountId)

	return nil
}
