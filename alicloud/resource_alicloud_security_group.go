package alicloud

import (
	"fmt"
	"time"

	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunSecurityGroupCreate,
		Read:   resourceAliyunSecurityGroupRead,
		Update: resourceAliyunSecurityGroupUpdate,
		Delete: resourceAliyunSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityGroupName,
			},

			"description": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateSecurityGroupDescription,
			},

			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"inner_access": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceAliyunSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.CreateSecurityGroup(buildAliyunSecurityGroupArgs(d, meta))
	})
	if err != nil {
		return err
	}
	resp := raw.(*ecs.CreateSecurityGroupResponse)
	if resp == nil {
		return fmt.Errorf("Creating security group got a nil response.")
	}
	d.SetId(resp.SecurityGroupId)
	return resourceAliyunSecurityGroupUpdate(d, meta)
}

func resourceAliyunSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	var sg *ecs.DescribeSecurityGroupAttributeResponse
	err := resource.Retry(1*time.Minute, func() *resource.RetryError {
		raw, e := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeSecurityGroupAttribute(d.Id())
		})
		if e != nil {
			if NotFoundError(e) || IsExceptedErrors(e, []string{InvalidSecurityGroupIdNotFound}) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error DescribeSecurityGroupAttribute: %#v", e))
		}
		group := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
		if group.SecurityGroupId != "" {
			sg = &group
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Create security group timeout and got an error: %#v", e))
	})

	if err != nil {
		return err
	}
	if sg == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", sg.SecurityGroupName)
	d.Set("description", sg.Description)
	d.Set("vpc_id", sg.VpcId)
	d.Set("inner_access", sg.InnerAccessPolicy == string(GroupInnerAccept))

	return nil
}

func resourceAliyunSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)

	d.Partial(true)
	attributeUpdate := false
	args := ecs.CreateModifySecurityGroupAttributeRequest()
	args.SecurityGroupId = d.Id()

	if d.HasChange("name") && !d.IsNewResource() {
		d.SetPartial("name")
		args.SecurityGroupName = d.Get("name").(string)

		attributeUpdate = true
	}

	if d.HasChange("description") && !d.IsNewResource() {
		d.SetPartial("description")
		args.Description = d.Get("description").(string)

		attributeUpdate = true
	}
	if attributeUpdate {
		_, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupAttribute(args)
		})
		if err != nil {
			return err
		}
	}

	if d.HasChange("inner_access") {
		policy := GroupInnerAccept
		if !d.Get("inner_access").(bool) {
			policy = GroupInnerDrop
		}
		args := ecs.CreateModifySecurityGroupPolicyRequest()
		args.SecurityGroupId = d.Id()
		args.InnerAccessPolicy = string(policy)

		_, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ModifySecurityGroupPolicy(args)
		})
		if err != nil {
			return fmt.Errorf("ModifySecurityGroupPolicy got an error: %#v.", err)
		}
	}

	d.Partial(false)

	return resourceAliyunSecurityGroupRead(d, meta)
}

func resourceAliyunSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)
	ecsService := EcsService{client}
	req := ecs.CreateDeleteSecurityGroupRequest()
	req.SecurityGroupId = d.Id()

	return resource.Retry(6*time.Minute, func() *resource.RetryError {
		_, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DeleteSecurityGroup(req)
		})

		if err != nil {
			if IsExceptedError(err, SgDependencyViolation) {
				return resource.RetryableError(fmt.Errorf("Delete security group timeout and got an error: %#v", err))
			}
		}

		sg, err := ecsService.DescribeSecurityGroupAttribute(d.Id())

		if err != nil {
			if IsExceptedError(err, InvalidSecurityGroupIdNotFound) {
				return nil
			}
			return resource.NonRetryableError(err)
		} else if sg.SecurityGroupId == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete security group timeout and got an error: %#v", err))
	})

}

func buildAliyunSecurityGroupArgs(d *schema.ResourceData, meta interface{}) *ecs.CreateSecurityGroupRequest {

	args := ecs.CreateCreateSecurityGroupRequest()

	if v := d.Get("name").(string); v != "" {
		args.SecurityGroupName = v
	}

	if v := d.Get("description").(string); v != "" {
		args.Description = v
	}

	if v := d.Get("vpc_id").(string); v != "" {
		args.VpcId = v
	}
	args.ClientToken = buildClientToken("TF-CreateSecurityGroup")

	return args
}
