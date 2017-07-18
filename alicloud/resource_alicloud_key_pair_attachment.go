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

	return resourceAlicloudKeyPairAttachmentUpdate(d, meta)
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

	update := false
	if d.HasChange("key_name") {
		d.SetPartial("key_name")
		update = true
	}

	instanceIds := convertListToJsonString(d.Get("instance_ids").(*schema.Set).List())
	if d.HasChange("instance_ids") {
		d.SetPartial("instance_ids")
		update = true
	}

	if update {
		args := &ecs.AttachKeyPairArgs{
			RegionId:    getRegion(d, meta),
			KeyPairName: d.Get("key_name").(string),
			InstanceIds: instanceIds,
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
	}
	d.Partial(false)

	d.SetId(d.Get("key_name").(string) + ":" + instanceIds)

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
