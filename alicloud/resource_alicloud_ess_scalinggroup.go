package alicloud

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/denverdino/aliyungo/ess"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
	"time"
)

func resourceAlicloudEss() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEssCreate,
		Read:   resourceAliyunEssRead,
		Update: resourceAliyunEssUpdate,
		Delete: resourceAliyunEssDelete,

		Schema: map[string]*schema.Schema{
			"min_size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"max_size": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"scaling_group_name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_cooldown": &schema.Schema{
				Type:     schema.TypeInt,
				Default:  300,
				Optional: true,
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"removal_policys": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MaxItems: 4,
			},
			"db_instance_ids": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
				MaxItems: 3,
			},
			"loadbalancer_ids": &schema.Schema{
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Optional: true,
			},
		},
	}
}

func resourceAliyunEssCreate(d *schema.ResourceData, meta interface{}) error {

	args, err := buildAliyunVpcArgs(d, meta)
	if err != nil {
		return err
	}

	ecsconn := meta.(*AliyunClient).ecsconn

	vpc, err := ecsconn.CreateVpc(args)
	if err != nil {
		return err
	}

	d.SetId(vpc.VpcId)
	d.Set("router_table_id", vpc.RouteTableId)

	err = ecsconn.WaitForVpcAvailable(args.RegionId, vpc.VpcId, 60)
	if err != nil {
		return fmt.Errorf("Timeout when WaitForVpcAvailable")
	}

	return resourceAliyunVpcRead(d, meta)
}

func resourceAliyunEssRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	vpc, err := client.DescribeVpc(d.Id())
	if err != nil {
		return err
	}

	if vpc == nil {
		d.SetId("")
		return nil
	}

	d.Set("cidr_block", vpc.CidrBlock)
	d.Set("name", vpc.VpcName)
	d.Set("description", vpc.Description)
	d.Set("router_id", vpc.VRouterId)

	return nil
}

func resourceAliyunEssUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).ecsconn

	d.Partial(true)

	attributeUpdate := false
	args := &ecs.ModifyVpcAttributeArgs{
		VpcId: d.Id(),
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		args.VpcName = d.Get("name").(string)

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		args.Description = d.Get("description").(string)

		attributeUpdate = true
	}

	if attributeUpdate {
		if err := conn.ModifyVpcAttribute(args); err != nil {
			return err
		}
	}

	d.Partial(false)

	return nil
}

func resourceAliyunEssDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DeleteVpc(d.Id())

		if err != nil {
			return resource.RetryableError(fmt.Errorf("Vpc in use - trying again while it is deleted."))
		}

		args := &ecs.DescribeVpcsArgs{
			RegionId: getRegion(d, meta),
			VpcId:    d.Id(),
		}
		vpc, _, descErr := conn.DescribeVpcs(args)
		if descErr != nil {
			return resource.NonRetryableError(err)
		} else if vpc == nil || len(vpc) < 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Vpc in use - trying again while it is deleted."))
	})
}

func buildAlicloudEssArgs(d *schema.ResourceData, meta interface{}) (*ess.CreateScalingGroupRequest, error) {
	client := meta.(*AliyunClient)
	args := &ess.CreateScalingGroupRequest{
		RegionId:        getRegion(d, meta),
		MinSize:         d.Get("min_size").(int),
		MaxSize:         d.Get("max_size").(int),
		DefaultCooldown: d.Get("default_cooldown").(int),
	}

	if v := d.Get("scaling_group_name").(string); v != "" {
		args.ScalingGroupName = v
	}

	if v := d.Get("vswitch_id").(string); v != "" {
		args.VSwitchId = v

		// get vpcId
		vpcId, err := client.GetVpcIdByVSwitchId(v)

		if err != nil {
			return nil, fmt.Errorf("VswitchId %s is not valid of current region", v)
		}
		// fill vpcId by vswitchId
		args.VpcId = vpcId

	}

	return args, nil
}
