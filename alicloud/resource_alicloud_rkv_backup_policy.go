package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRKVBackupPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRKVBackupPolicyCreate,
		Read:   resourceAlicloudRKVBackupPolicyRead,
		Update: resourceAlicloudRKVBackupPolicyUpdate,
		Delete: resourceAlicloudRKVBackupPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"instance_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"preferred_backup_time": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"preferred_backup_period": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"Monday",
					"Tuesday",
					"Wednesday",
					"Thursday",
					"Friday",
					"Saturday",
					"Sunday",
				}),
			},
		},
	}
}

func resourceAlicloudRKVBackupPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn

	request := r_kvstore.CreateModifyBackupPolicyRequest()
	request.InstanceId = d.Get("instance_id").(string)
	request.PreferredBackupTime = d.Get("preferred_backup_time").(string)
	request.PreferredBackupPeriod = d.Get("preferred_backup_period").(string)

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if _, err := conn.ModifyBackupPolicy(request); err != nil {
			return resource.NonRetryableError(fmt.Errorf("Create security whitelist ips got an error: %#v", err))
		}
		return nil
	})

	if err != nil {
		return err
	}

	// A security ip whitelist does not have a native IP.
	d.SetId(fmt.Sprintf("%s%s%s", request.InstanceId, COLON_SEPARATED, resource.UniqueId()))

	return resourceAlicloudRKVBackupPolicyRead(d, meta)
}

func resourceAlicloudRKVBackupPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn
	instanceID := strings.Split(d.Id(), COLON_SEPARATED)[0]

	request := r_kvstore.CreateDescribeBackupPolicyRequest()
	request.InstanceId = instanceID
	policy, err := conn.DescribeBackupPolicy(request)
	if err != nil {
		if NotFoundRKVInstance(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe RKV Security IPs: %#v", err)
	}
	if policy == nil {
		d.SetId("")
		return nil
	}

	d.Set("instance_id", instanceID)
	d.Set("preferred_backup_time", policy.PreferredBackupTime)
	d.Set("preferred_backup_period", policy.PreferredBackupPeriod)

	return nil
}

func resourceAlicloudRKVBackupPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn
	instanceID := strings.Split(d.Id(), COLON_SEPARATED)[0]

	if d.HasChange("preferred_backup_time") && d.HasChange("preferred_backup_period") {
		request := r_kvstore.CreateModifyBackupPolicyRequest()
		request.InstanceId = instanceID
		request.PreferredBackupTime = d.Get("preferred_backup_time").(string)
		request.PreferredBackupPeriod = d.Get("preferred_backup_period").(string)
		if _, err := conn.ModifyBackupPolicy(request); err != nil {
			return err
		}
	}

	return resourceAlicloudRKVBackupPolicyRead(d, meta)
}

func resourceAlicloudRKVBackupPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	// There is no explicit delete, only update with modified security ips
	return nil
}
