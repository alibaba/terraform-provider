package alicloud

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudRKVInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudRKVInstanceCreate,
		Read:   resourceAlicloudRKVInstanceRead,
		Update: resourceAlicloudRKVInstanceUpdate,
		Delete: resourceAlicloudRKVInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"instance_name": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateRKVInstanceName,
			},
			"password": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Sensitive:    true,
				ValidateFunc: validateRKVPassword,
			},
			"instance_class": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"charge_type": &schema.Schema{
				Type:         schema.TypeString,
				ValidateFunc: validateAllowedStringValue([]string{string(PostPaid), string(PrePaid)}),
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
			},
			"period": &schema.Schema{
				Type:             schema.TypeInt,
				ValidateFunc:     validateAllowedIntValue([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 12, 24, 36}),
				Optional:         true,
				Default:          1,
				DiffSuppressFunc: rkvPostPaidDiffSuppressFunc,
			},

			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "CLASSIC",
				ValidateFunc: validateAllowedStringValue([]string{
					string(CLASSIC),
					string(VPC),
				}),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return true
				},
				Deprecated: "Field 'instance_network_type' has been deprecated from provider version 1.5.0.",
			},
			"instance_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "Redis",
				ValidateFunc: validateAllowedStringValue([]string{
					"Memcache",
					"Redis",
				}),
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
			},
			"engine_version": &schema.Schema{
				Type:     schema.TypeString,
				ForceNew: true,
				Optional: true,
				Default:  "2.8",
				ValidateFunc: validateAllowedStringValue([]string{
					"2.8",
					"4.0",
				}),
			},
			"connection_domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"private_ip": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudRKVInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn

	request, err := buildRKVCreateRequest(d, meta)
	if err != nil {
		return err
	}

	resp, err := conn.CreateInstance(request)

	if err != nil {
		return fmt.Errorf("Error creating Alicloud db instance: %#v", err)
	}

	d.SetId(resp.InstanceId)

	// wait instance status change from Creating to Normal

	if err := client.WaitForRKVInstance(d.Id(), Normall, DefaultLongTimeout); err != nil {
		return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
	}

	return resourceAlicloudRKVInstanceRead(d, meta)
}

func resourceAlicloudRKVInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	conn := client.rkvconn
	d.Partial(true)

	if d.HasChange("instance_class") {
		request := r_kvstore.CreateModifyInstanceSpecRequest()
		request.InstanceId = d.Id()
		request.InstanceClass = d.Get("instance_class").(string)
		if _, err := conn.ModifyInstanceSpec(request); err != nil {
			return err
		}
		// wait instance status is Normal after modifying
		if err := client.WaitForRKVInstance(d.Id(), Normall, 500); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}

		d.SetPartial("instance_class")
	}

	if d.HasChange("instance_name") {
		request := r_kvstore.CreateModifyInstanceAttributeRequest()
		request.InstanceId = d.Id()
		request.InstanceName = d.Get("instance_name").(string)
		// wait instance status is Normal before modifying
		if err := client.WaitForRKVInstance(d.Id(), Normall, 500); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}
		if _, err := conn.ModifyInstanceAttribute(request); err != nil {
			return fmt.Errorf("ModifyRKVInstanceDescription got an error: %#v", err)
		}
		// wait instance status is Normal after modifying
		if err := client.WaitForRKVInstance(d.Id(), Normall, 500); err != nil {
			return fmt.Errorf("WaitForInstance %s got error: %#v", Running, err)
		}

		d.SetPartial("instance_name")
	}

	d.Partial(false)
	return resourceAlicloudDBInstanceRead(d, meta)
}

func resourceAlicloudRKVInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)
	instance, err := client.DescribeRKVInstanceById(d.Id())
	if err != nil {
		if NotFoundRKVInstance(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe rKV InstanceAttribute: %#v", err)
	}
	d.SetId(instance.InstanceId)
	d.Set("instance_name", instance.InstanceName)
	d.Set("instance_class", instance.InstanceClass)
	d.Set("zone_id", instance.ZoneId)
	d.Set("charge_type", instance.ChargeType)
	d.Set("instance_type", instance.InstanceType)
	d.Set("vswitch_id", instance.VSwitchId)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("connection_domain", instance.ConnectionDomain)
	d.Set("port", instance.Port)
	d.Set("private_ip", instance.PrivateIp)

	return nil
}

func resourceAlicloudRKVInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	instance, err := client.DescribeRKVInstanceById(d.Id())
	if err != nil {
		if NotFoundRKVInstance(err) {
			return nil
		}
		return fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err)
	}
	if PayType(instance.ChargeType) == Prepaid {
		return fmt.Errorf("At present, 'Prepaid' instance cannot be deleted and must wait it to be expired and release it automatically")
	}
	request := r_kvstore.CreateDeleteInstanceRequest()
	request.InstanceId = d.Id()

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.rkvconn.DeleteInstance(request)

		if err != nil {
			if NotFoundRKVInstance(err) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Delete DB instance timeout and got an error: %#v", err))
		}

		instance, err := client.DescribeRKVInstanceById(d.Id())
		if err != nil {
			if NotFoundRKVInstance(err) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err))
		}
		if instance == nil {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Delete DB instance timeout and got an error: %#v", err))
	})
}

func buildRKVCreateRequest(d *schema.ResourceData, meta interface{}) (*r_kvstore.CreateInstanceRequest, error) {
	client := meta.(*AliyunClient)
	request := r_kvstore.CreateCreateInstanceRequest()
	request.InstanceName = Trim(d.Get("instance_name").(string))
	request.RegionId = string(getRegion(d, meta))
	request.EngineVersion = Trim(d.Get("engine_version").(string))
	request.InstanceClass = Trim(d.Get("instance_class").(string))
	request.NetworkType = Trim(d.Get("network_type").(string))
	request.ChargeType = Trim(d.Get("charge_type").(string))
	request.Password = Trim(d.Get("password").(string))
	if PayType(request.ChargeType) == PrePaid {
		request.Period = d.Get("Period").(string)
	}

	if zone, ok := d.GetOk("zone_id"); ok && Trim(zone.(string)) != "" {
		request.ZoneId = Trim(zone.(string))
	}

	vswitchID := Trim(d.Get("vswitch_id").(string))
	if vswitchID != "" {
		request.VSwitchId = vswitchID
		request.NetworkType = strings.ToUpper(string(Vpc))
		d.Set("network_type", string(VPC))

		// check vswitchId in zone
		vsw, err := client.DescribeVswitch(vswitchID)
		if err != nil {
			return nil, fmt.Errorf("DescribeVSwitch got an error: %#v", err)
		}

		if request.ZoneId == "" {
			request.ZoneId = vsw.ZoneId
		} else if strings.Contains(request.ZoneId, MULTI_IZ_SYMBOL) {
			zonestr := strings.Split(strings.SplitAfter(request.ZoneId, "(")[1], ")")[0]
			if !strings.Contains(zonestr, string([]byte(vsw.ZoneId)[len(vsw.ZoneId)-1])) {
				return nil, fmt.Errorf("The specified vswitch %s isn't in the multi zone %s", vsw.VSwitchId, request.ZoneId)
			}
		} else if request.ZoneId != vsw.ZoneId {
			return nil, fmt.Errorf("The specified vswitch %s isn't in the zone %s", vsw.VSwitchId, request.ZoneId)
		}

		request.VpcId = vsw.VpcId
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		uuid = resource.UniqueId()
	}
	request.Token = fmt.Sprintf("%s", uuid)

	return request, nil
}
