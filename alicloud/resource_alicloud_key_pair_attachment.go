package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"time"
)

func resourceAlicloudKeyPairAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudKeyPairAttachmentCreate,
		Read:   resourceAlicloudKeyPairAttachmentRead,
		Update: resourceAlicloudKeyPairAttachmentUpdate,
		Delete: resourceAlicloudKeyPairAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key_name": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateKeyPairName,
			},
			"instance_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceAlicloudKeyPairAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	instance_ids := convertListToJsonString(d.Get("instance_ids").(*schema.Set).List())

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		er := conn.AttachKeyPair(&ecs.AttachKeyPairArgs{
			RegionId:    getRegion(d, meta),
			KeyPairName: d.Get("key_name").(string),
			InstanceIds: instance_ids,
		})
		if er != nil {
			if IsExceptedError(er, KeyPairServiceUnavailable) {
				return resource.RetryableError(fmt.Errorf("Key Pair is attaching and gets an error: %#v -- try again...", er))
			}
			return resource.NonRetryableError(fmt.Errorf("Error Attach KeyPair: %#v", er))
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error Attachment KeyPair with InstanceIds: %s", err)
	}

	d.SetId(d.Get("key_name").(string) + ":" + instance_ids)

	return resourceAlicloudKeyPairAttachmentRead(d, meta)
}

func resourceAlicloudKeyPairAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	keyname := strings.Split(d.Id(), ":")[0]
	keypairs, _, err := conn.DescribeKeyPairs(&ecs.DescribeKeyPairsArgs{
		RegionId:    getRegion(d, meta),
		KeyPairName: keyname,
	})
	if err != nil {
		if IsExceptedError(err, KeyPairNotFound) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Retrieving KeyPair: %s", err)
	}

	if len(keypairs) > 0 {
		d.Set("key_name", keypairs[0].KeyPairName)
		d.Set("instance_ids", d.Get("instance_ids"))
		return nil
	}

	return fmt.Errorf("Unable to find key pair within: %#v", keypairs)
}

func resourceAlicloudKeyPairAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	d.Partial(true)
	newKey := d.Get("key_name").(string)
	var oldKey string
	if d.HasChange("key_name") {
		d.SetPartial("key_name")
		o, _ := d.GetChange("key_name")
		oldKey = o.(string)
	}

	instanceIds := convertListToJsonString(d.Get("instance_ids").(*schema.Set).List())
	var newInstanceIds, oldInstanceIds string
	if d.HasChange("instance_ids") {
		d.SetPartial("instance_ids")
		o, n := d.GetChange("instance_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		newInstanceIds = convertListToJsonString(ns.Difference(os).List())
		oldInstanceIds = convertListToJsonString(os.Difference(ns).List())
	}

	if oldInstanceIds != "" {
		keyname := newKey
		if oldKey != "" {
			keyname = oldKey
		}
		err := conn.DetachKeyPair(&ecs.DetachKeyPairArgs{
			RegionId:    getRegion(d, meta),
			KeyPairName: keyname,
			InstanceIds: oldInstanceIds,
		})
		if err != nil {
			return fmt.Errorf("Error Detach Key Pair: %#v", err)
		}
	}

	args := &ecs.AttachKeyPairArgs{
		RegionId:    getRegion(d, meta),
		KeyPairName: newKey,
	}

	if oldKey != "" {
		args.InstanceIds = instanceIds
	} else if newInstanceIds != "" {
		args.InstanceIds = newInstanceIds
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		if er := conn.AttachKeyPair(args); er != nil {
			if IsExceptedError(er, KeyPairServiceUnavailable) {
				return resource.RetryableError(fmt.Errorf("Key Pair is attaching and gets an error: %#v -- try again...", er))
			}
			return resource.NonRetryableError(fmt.Errorf("Error Attach KeyPair: %#v", er))
		}
		return nil
	})

	if err != nil {
		return err
	}

	return resourceAlicloudKeyPairAttachmentRead(d, meta)
}

func resourceAlicloudKeyPairAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn
	keyname := strings.Split(d.Id(), ":")[0]
	instanceIds := strings.Split(d.Id(), ":")[1]

	err := conn.DetachKeyPair(&ecs.DetachKeyPairArgs{
		RegionId:    getRegion(d, meta),
		KeyPairName: keyname,
		InstanceIds: instanceIds,
	})
	return err
}
