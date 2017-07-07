package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ram"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRamAK() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRamAKCreate,
		Read:   resourceAlicloudRamAKRead,
		Update: resourceAlicloudRamAKUpdate,
		Delete: resourceAlicloudRamAKDelete,

		Schema: map[string]*schema.Schema{
			"user_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRamName,
			},
			"access_key_secret": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"create_date": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRamAKCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.UserQueryRequest{}
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		args.UserName = v.(string)
	}

	response, err := conn.CreateAccessKey(args)
	if err != nil {
		return fmt.Errorf("CreateAccessKey got an error: %#v", err)
	}

	d.SetId(response.AccessKey.AccessKeyId)
	return resourceAlicloudRamAKUpdate(d, meta)
}

func resourceAlicloudRamAKUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	d.Partial(true)

	args := ram.UpdateAccessKeyRequest{
		UserAccessKeyId: d.Id(),
		Status:          ram.State(d.Get("status").(string)),
	}
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		args.UserName = v.(string)
	}

	if d.HasChange("status") {
		d.SetPartial("status")
		if _, err := conn.UpdateAccessKey(args); err != nil {
			return fmt.Errorf("UpdateAccessKey got an error: %#v", err)
		}
	}

	d.Partial(false)
	return resourceAlicloudRamAKRead(d, meta)
}

func resourceAlicloudRamAKRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.UserQueryRequest{}
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		args.UserName = v.(string)
	}

	response, err := conn.ListAccessKeys(args)
	if err != nil {
		return fmt.Errorf("Get list access keys got an error: %#v", err)
	}

	accessKeys := response.AccessKeys.AccessKey
	if accessKeys == nil || len(accessKeys) <= 0 {
		return fmt.Errorf("No access keys found.")
	}

	ak := accessKeys[0]
	d.Set("access_key_secret", ak.AccessKeySecret)
	d.Set("create_date", ak.CreateDate)
	d.Set("status", ak.Status)
	return nil
}

func resourceAlicloudRamAKDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ramconn

	args := ram.UpdateAccessKeyRequest{
		UserAccessKeyId: d.Id(),
	}
	if v, ok := d.GetOk("user_name"); ok && v.(string) != "" {
		args.UserName = v.(string)
	}

	_, err := conn.DeleteAccessKey(args)
	if err != nil {
		return fmt.Errorf("DeleteAccessKey got an error: %#v", err)
	}
	return nil
}
