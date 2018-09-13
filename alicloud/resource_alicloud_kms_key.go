package alicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/denverdino/aliyungo/kms"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudKmsKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKmsKeyCreate,
		Read:   resourceAlicloudKmsKeyRead,
		Update: resourceAlicloudKmsKeyUpdate,
		Delete: resourceAlicloudKmsKeyDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "From Terraform",
				ValidateFunc: validateStringLengthInRange(0, 8192),
			},
			"key_usage": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, es []error) {
					value := v.(string)
					if !(kms.KeyUsage(value) == kms.KEY_USAGE_ENCRYPT_DECRYPT) {
						es = append(es, fmt.Errorf(
							"%q must be %s", k, kms.KEY_USAGE_ENCRYPT_DECRYPT))
					}
					return
				},
				Default: kms.KEY_USAGE_ENCRYPT_DECRYPT,
			},
			"is_enabled": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"deletion_window_in_days": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validateIntegerInRange(7, 30),
				Default:      30,
			},
			"arn": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudKmsKeyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	args := kms.CreateKeyArgs{
		KeyUsage: kms.KeyUsage(d.Get("key_usage").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		args.Description = v.(string)
	}
	raw, err := client.RunSafelyWithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.CreateKey(&args)
	})
	if err != nil {
		return fmt.Errorf("CreateKey got an error: %#v.", err)
	}
	resp := raw.(*kms.CreateKeyResponse)
	d.SetId(resp.KeyMetadata.KeyId)

	return resourceAlicloudKmsKeyUpdate(d, meta)
}

func resourceAlicloudKmsKeyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	raw, err := client.RunSafelyWithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.DescribeKey(d.Id())
	})
	if err != nil {
		if IsExceptedError(err, ForbiddenKeyNotFound) {
			return nil
		}
		return fmt.Errorf("DescribeKey got an error: %#v.", err)
	}
	key := raw.(*kms.DescribeKeyResponse)
	if KeyState(key.KeyMetadata.KeyState) == PendingDeletion {
		log.Printf("[WARN] Removing KMS key %s because it's already gone", d.Id())
		d.SetId("")
		return nil
	}

	d.Set("description", key.KeyMetadata.Description)
	d.Set("key_usage", key.KeyMetadata.KeyUsage)
	d.Set("is_enabled", KeyState(key.KeyMetadata.KeyState) == Enabled)
	d.Set("deletion_window_in_days", d.Get("deletion_window_in_days").(int))
	d.Set("arn", key.KeyMetadata.Arn)

	return nil
}

func resourceAlicloudKmsKeyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	d.Partial(true)

	if d.HasChange("is_enabled") {
		raw, err := client.RunSafelyWithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.DescribeKey(d.Id())
		})
		if err != nil {
			return fmt.Errorf("DescribeKey got an error: %#v.", err)
		}
		key := raw.(*kms.DescribeKeyResponse)
		if d.Get("is_enabled").(bool) && KeyState(key.KeyMetadata.KeyState) == Disabled {
			_, err := client.RunSafelyWithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
				return kmsClient.EnableKey(d.Id())
			})
			if err != nil {
				return fmt.Errorf("Enable key got an error: %#v.", err)
			}
		}

		if !d.Get("is_enabled").(bool) && KeyState(key.KeyMetadata.KeyState) == Enabled {
			_, err := client.RunSafelyWithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
				return kmsClient.DisableKey(d.Id())
			})
			if err != nil {
				return fmt.Errorf("Disable key got an error: %#v.", err)
			}
		}
		d.SetPartial("is_enabled")
	}

	d.Partial(false)

	return resourceAlicloudKmsKeyRead(d, meta)
}

func resourceAlicloudKmsKeyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	_, err := client.RunSafelyWithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
		return kmsClient.ScheduleKeyDeletion(&kms.ScheduleKeyDeletionArgs{
			KeyId:               d.Id(),
			PendingWindowInDays: d.Get("deletion_window_in_days").(int),
		})
	})
	if err != nil {
		return err
	}

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.RunSafelyWithKmsClient(func(kmsClient *kms.Client) (interface{}, error) {
			return kmsClient.DescribeKey(d.Id())
		})
		if err != nil {
			if IsExceptedError(err, ForbiddenKeyNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("DescribeKey got an error: %#v.", err))
		}
		key := raw.(*kms.DescribeKeyResponse)

		if key == nil || KeyState(key.KeyMetadata.KeyState) == PendingDeletion {
			log.Printf("[WARN] Removing KMS key %s because it's already gone", d.Id())
			d.SetId("")
			return nil
		}
		return resource.RetryableError(fmt.Errorf("ScheduleKeyDeletion timeout."))
	})
}
