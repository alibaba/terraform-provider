package alicloud

import (
	"fmt"
	"time"

	"github.com/denverdino/aliyungo/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudEssAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssAttachmentCreate,
		Read:   resourceAliyunEssAttachmentRead,
		Update: resourceAliyunEssAttachmentUpdate,
		Delete: resourceAliyunEssAttachmentDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"scaling_group_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},

			"instance_ids": &schema.Schema{
				Type:     schema.TypeSet,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
				MaxItems: 20,
				MinItems: 1,
			},

			"force": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceAliyunEssAttachmentCreate(d *schema.ResourceData, meta interface{}) error {

	d.SetId(d.Get("scaling_group_id").(string))

	return resourceAliyunEssAttachmentUpdate(d, meta)
}

func resourceAliyunEssAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	d.Partial(true)

	groupId := d.Id()
	if d.HasChange("instance_ids") {
		group, err := client.DescribeScalingGroupById(groupId)
		if err != nil {
			return fmt.Errorf("DescribeScalingGroupById %s error: %#v", groupId, err)
		}
		if group.LifecycleState == ess.Inacitve {
			return fmt.Errorf("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState)
		} else {
			if err := client.essconn.WaitForScalingGroup(getRegion(d, meta), group.ScalingGroupId, ess.Active, DefaultTimeout); err != nil {
				return fmt.Errorf("WaitForScalingGroup is %#v got an error: %#v.", ess.Active, err)
			}
		}
		o, n := d.GetChange("instance_ids")
		os := o.(*schema.Set)
		ns := n.(*schema.Set)
		remove := os.Difference(ns).List()
		add := ns.Difference(os).List()

		if len(add) > 0 {

			if err := resource.Retry(5*time.Minute, func() *resource.RetryError {

				if _, err := client.essconn.AttachInstances(&ess.AttachInstancesArgs{
					ScalingGroupId: groupId,
					InstanceId:     convertArrayInterfaceToArrayString(add),
				}); err != nil {
					if IsExceptedError(err, IncorrectCapacityMaxSize) {
						instances, _, err := client.essconn.DescribeScalingInstances(&ess.DescribeScalingInstancesArgs{
							RegionId:       getRegion(d, meta),
							ScalingGroupId: d.Id(),
						})
						if err != nil {
							return resource.NonRetryableError(fmt.Errorf("DescribeScalingInstances got an error: %#v", err))
						}
						var autoAdded, attached []string
						if len(instances) > 0 {
							for _, inst := range instances {
								if inst.CreationType == "Attached" {
									attached = append(attached, inst.InstanceId)
								} else {
									autoAdded = append(autoAdded, inst.InstanceId)
								}
							}
						}
						if len(add) > group.MaxSize {
							return resource.NonRetryableError(fmt.Errorf("To attach %d instances, the total capacity will be greater than the scaling group max size %d. "+
								"Please enlarge scaling group max size.", len(add), group.MaxSize))
						}

						if len(autoAdded) > 0 {
							if d.Get("force").(bool) {
								if err := client.EssRemoveInstances(groupId, autoAdded); err != nil {
									return resource.NonRetryableError(err)
								}
								time.Sleep(5)
								return resource.RetryableError(fmt.Errorf("Autocreated result in attaching instances got an error: %#v", err))
							} else {
								return resource.NonRetryableError(fmt.Errorf("To attach the instances, the total capacity will be greater than the scaling group max size %d."+
									"Please enlarge scaling group max size or set 'force' to true to remove autocreated instances: %#v.", group.MaxSize, autoAdded))
							}
						}

						if len(attached) > 0 {
							return resource.NonRetryableError(fmt.Errorf("To attach the instances, the total capacity will be greater than the scaling group max size %d. "+
								"Please enlarge scaling group max size or remove already attached instances: %#v.", group.MaxSize, attached))
						}
					}
					if IsExceptedError(err, ScalingActivityInProgress) {
						time.Sleep(5)
						return resource.RetryableError(fmt.Errorf("Progress results in Attaching instances got an error: %#v", err))
					}
					return resource.NonRetryableError(fmt.Errorf("Attaching instances got an error: %#v", err))
				}
				return nil
			}); err != nil {
				return err
			}

			if err := resource.Retry(3*time.Minute, func() *resource.RetryError {

				instances, _, err := client.essconn.DescribeScalingInstances(&ess.DescribeScalingInstancesArgs{
					RegionId:       getRegion(d, meta),
					ScalingGroupId: d.Id(),
					InstanceId:     convertArrayInterfaceToArrayString(add),
				})
				if err != nil {
					return resource.NonRetryableError(err)
				}
				if len(instances) < 0 {
					return resource.RetryableError(fmt.Errorf("There are no ECS instances have been attached."))
				}

				for _, inst := range instances {
					if inst.LifecycleState != ess.InService {
						return resource.RetryableError(fmt.Errorf("There are still ECS instances are not %s.", ess.InService))
					}
				}
				return nil
			}); err != nil {
				return err
			}
		}
		if len(remove) > 0 {
			if err := client.EssRemoveInstances(groupId, convertArrayInterfaceToArrayString(remove)); err != nil {
				return err
			}
		}

		d.SetPartial("instance_ids")
	}

	d.Partial(false)

	return resourceAliyunEssAttachmentRead(d, meta)
}

func resourceAliyunEssAttachmentRead(d *schema.ResourceData, meta interface{}) error {

	instances, _, err := meta.(*AliyunClient).essconn.DescribeScalingInstances(&ess.DescribeScalingInstancesArgs{
		RegionId:       getRegion(d, meta),
		ScalingGroupId: d.Id(),
		CreationType:   "Attached",
	})

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe ESS scaling instances: %#v", err)
	}

	if len(instances) < 1 {
		d.SetId("")
		return nil
	}

	var instanceIds []string
	for _, inst := range instances {
		instanceIds = append(instanceIds, inst.InstanceId)
	}

	d.Set("scaling_group_id", instances[0].ScalingGroupId)
	d.Set("instance_ids", instanceIds)

	return nil
}

func resourceAliyunEssAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	group, err := client.DescribeScalingGroupById(d.Id())
	if err != nil {
		return fmt.Errorf("DescribeScalingGroupById %s error: %#v", d.Id(), err)
	}
	if group.LifecycleState != ess.Active {
		return fmt.Errorf("Scaling group current status is %s, please active it before attaching or removing ECS instances.", group.LifecycleState)
	}

	return client.EssRemoveInstances(d.Id(), convertArrayInterfaceToArrayString(d.Get("instance_ids").(*schema.Set).List()))
}

func convertArrayInterfaceToArrayString(elm []interface{}) (arr []string) {
	if len(elm) < 1 {
		return
	}
	for _, e := range elm {
		arr = append(arr, e.(string))
	}
	return
}
