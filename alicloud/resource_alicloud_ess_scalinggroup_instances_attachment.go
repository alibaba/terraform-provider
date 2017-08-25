package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"time"
)

func resourceAlicloudEssScalingGroupInstancesAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssScalingGroupInstancesAttachmentCreate,
		Read:   resourceAliyunEssScalingGroupInstancesAttachmentRead,
		Delete: resourceAliyunEssScalingGroupInstancesAttachmentDelete,

		Schema: map[string]*schema.Schema{
			"scaling_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
				MaxItems: 20,
			},
		},
	}
}

func resourceAliyunEssScalingGroupInstancesAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	essconn := meta.(*AliyunClient).essconn

	ids := []string{}
	for _, v := range d.Get("instance_ids").(*schema.Set).List() {
		ids = append(ids, v.(string))
	}
	instanceId := common.FlattenArray(ids)
	args := &ess.AttachInstancesArgs{
		ScalingGroupId: d.Get("scaling_group_id").(string),
		InstanceId:     instanceId,
	}

	_, err := essconn.AttachInstances(args)
	if err != nil {
		return fmt.Errorf("AttachInstances got an error: %#v", err)
	}

	d.SetId(args.ScalingGroupId + COLON_SEPARATED + strings.Join(args.InstanceId, ","))
	return resourceAliyunEssScalingGroupInstancesAttachmentRead(d, meta)
}

func resourceAliyunEssScalingGroupInstancesAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	essconn := meta.(*AliyunClient).essconn

	parts := strings.Split(d.Id(), COLON_SEPARATED)
	scalingGroupId := parts[0]
	instanceId := common.FlattenArray(strings.Split(parts[1], ","))
	args := &ess.DescribeScalingInstancesArgs{
		RegionId:       getRegion(d, meta),
		ScalingGroupId: scalingGroupId,
		InstanceId:     instanceId,
	}
	instances, _, err := essconn.DescribeScalingInstances(args)
	if err != nil {
		return fmt.Errorf("Error Describe ESS scaling group instances Attribute: %#v", err)
	}

	if len(instances) < 1 {
		return fmt.Errorf("No instances found.")
	}

	instanceIds := []string{}
	for _, v := range instances {
		if v.ScalingGroupId != scalingGroupId {
			return fmt.Errorf("Error scaling group id for instances.")
		}
		instanceIds = append(instanceIds, v.InstanceId)
	}

	d.Set("scaling_group_id", instances[0].ScalingGroupId)
	d.Set("instance_ids", instanceIds)

	return nil
}

func resourceAliyunEssScalingGroupInstancesAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	essconn := meta.(*AliyunClient).essconn

	parts := strings.Split(d.Id(), COLON_SEPARATED)
	scalingGroupId := parts[0]
	instanceId := common.FlattenArray(strings.Split(parts[1], ","))
	args := &ess.AttachInstancesArgs{
		ScalingGroupId: scalingGroupId,
		InstanceId:     instanceId,
	}

	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		if _, err := essconn.RemoveInstances(args); err != nil {
			if IsExceptedError(err, IncorrectScalingGroupStatus) || IsExceptedError(err, ScalingActivityInProgress) {
				return resource.RetryableError(fmt.Errorf("Scaling group is in use or not active - trying again while it is deleted."))
			}
			if IsExceptedError(err, IncorrectLoadBalancerStatus) || IsExceptedError(err, IncorrectDBInstanceStatus) {
				return resource.RetryableError(fmt.Errorf("The action is not supported by the current status of the specified load balancer or the DB instance - trying again while it is deleted."))
			}
			return resource.NonRetryableError(err)
		}

		return nil
	})
}
