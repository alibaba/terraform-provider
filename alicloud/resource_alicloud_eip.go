package alicloud

import (
	"fmt"
	"strconv"
	"time"

	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunEip() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunEipCreate,
		Read:   resourceAliyunEipRead,
		Update: resourceAliyunEipUpdate,
		Delete: resourceAliyunEipDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"bandwidth": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"internet_charge_type": &schema.Schema{
				Type:         schema.TypeString,
				Default:      "PayByBandwidth",
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInternetChargeType,
			},

			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"instance": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAliyunEipCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	args, err := buildAliyunEipArgs(d, meta)
	if err != nil {
		return err
	}

	_, allocationID, err := conn.AllocateEipAddress(args)
	if err != nil {
		return err
	}

	err = conn.WaitForEip(getRegion(d, meta), allocationID, ecs.EipStatusAvailable, 60)
	if err != nil {
		return fmt.Errorf("Error Waitting for EIP available: %#v", err)
	}

	d.SetId(allocationID)

	return resourceAliyunEipUpdate(d, meta)
}

func resourceAliyunEipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	eip, err := client.DescribeEipAddress(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe Eip Attribute: %#v", err)
	}

	// Output parameter 'instance' would be deprecated in the next version.
	if eip.InstanceId != "" {
		d.Set("instance", eip.InstanceId)
	} else {
		d.Set("instance", "")
	}

	bandwidth, _ := strconv.Atoi(eip.Bandwidth)
	d.Set("bandwidth", bandwidth)
	d.Set("internet_charge_type", eip.InternetChargeType)
	d.Set("ip_address", eip.IpAddress)
	d.Set("status", eip.Status)

	return nil
}

func resourceAliyunEipUpdate(d *schema.ResourceData, meta interface{}) error {

	conn := meta.(*AliyunClient).ecsconn

	d.Partial(true)

	if d.HasChange("bandwidth") && !d.IsNewResource() {
		err := conn.ModifyEipAddressAttribute(d.Id(), d.Get("bandwidth").(int))
		if err != nil {
			return err
		}

		d.SetPartial("bandwidth")
	}

	d.Partial(false)

	return resourceAliyunEipRead(d, meta)
}

func resourceAliyunEipDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).ecsconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.ReleaseEipAddress(d.Id())

		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code == EipIncorrectStatus {
				return resource.RetryableError(fmt.Errorf("Delete EIP timeout and got an error:%#v.", err))
			}
		}

		args := &ecs.DescribeEipAddressesArgs{
			RegionId:     getRegion(d, meta),
			AllocationId: d.Id(),
		}

		eips, _, descErr := conn.DescribeEipAddresses(args)
		if descErr != nil {
			return resource.NonRetryableError(descErr)
		} else if eips == nil || len(eips) < 1 {
			return nil
		}
		return resource.RetryableError(fmt.Errorf("Delete EIP timeout and got an error:%#v.", err))
	})
}

func buildAliyunEipArgs(d *schema.ResourceData, meta interface{}) (*ecs.AllocateEipAddressArgs, error) {

	args := &ecs.AllocateEipAddressArgs{
		RegionId:           getRegion(d, meta),
		Bandwidth:          d.Get("bandwidth").(int),
		InternetChargeType: common.InternetChargeType(d.Get("internet_charge_type").(string)),
	}

	return args, nil
}
