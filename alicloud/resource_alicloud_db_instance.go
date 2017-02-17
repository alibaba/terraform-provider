package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/rds"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strconv"
	"time"
)

func resourceAliyunDBInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunDBInstanceCreate,
		Read:   resourceAliyunDBInstanceRead,
		Update: resourceAliyunDBInstanceUpdate,
		Delete: resourceAliyunDBInstanceDelete,

		Schema: map[string]*schema.Schema{
			"commodity_code": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"engine": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_version": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_class": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"db_instance_storage": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},

			"instance_charge_type": &schema.Schema{
				Type:     schema.TypeString, // rds.DBPayType
				Optional: true,
				Default:  rds.Postpaid,
			},
			"period_type": &schema.Schema{
				Type:     schema.TypeString, // common.TimeType
				Optional: true,
			},
			"period": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"auto_pay": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			//"order_id": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Computed: true,
			//},
			//"price": &schema.Schema{
			//	Type:     schema.TypeFloat,
			//	Computed: true,
			//},

			"zone_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"multi_az": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"db_instance_net_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"allocate_public_connection": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			//"connection_mode": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			//"master_user_name": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},
			//"master_user_password": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},

			//"preferred_backup_period": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},
			//"preferred_backup_time": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},
			//"backup_retention_period": &schema.Schema{
			//	Type:     schema.TypeString,
			//	Optional: true,
			//},

			"security_ip_list": &schema.Schema{
				Type:     schema.TypeString,
				Required: false,
				Default:  "127.0.0.1",
				Optional: true,
			},

			//"db_mappings": &schema.Schema{
			//	Type: schema.TypeList,
			//	Elem: &schema.Resource{
			//		Schema: map[string]*schema.Schema{
			//			"db_description": &schema.Schema{
			//				Type:     schema.TypeString,
			//				Optional: true,
			//			},
			//			"db_name": &schema.Schema{
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//			"character_set_name": &schema.Schema{
			//				Type:     schema.TypeString,
			//				Required: true,
			//			},
			//		},
			//	},
			//	Optional: true,
			//},
		},
	}
}

func resourceAliyunDBInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).rdsconn

	args, err := buildDBCreateOrderArgs(d, meta)
	if err != nil {
		return err
	}

	resp, err := conn.CreateOrder(args)

	if err != nil {
		return fmt.Errorf("Error creating Aliyun db instance: %#v", err)
	}

	instanceId := resp.DBInstanceId
	if instanceId == "" {
		return fmt.Errorf("Error get Aliyun db instance id")
	}

	d.SetId(instanceId)
	// we can't get this attr from DescribeDBInstance, so do it here
	d.Set("commodity_code", d.Get("commodity_code"))
	d.Set("instance_charge_type", d.Get("instance_charge_type"))
	d.Set("period", d.Get("period"))
	d.Set("period_type", d.Get("period_type"))
	d.Set("auto_pay", args.AutoPay)

	// after instance created, its status change from Creating to running
	if err := conn.WaitForInstance(d.Id(), rds.Running, defaultLongTimeout); err != nil {
		log.Printf("[DEBUG] WaitForInstance %s got error: %#v", rds.Running, err)
	}

	return resourceAliyunDBInstanceUpdate(d, meta)
}

func resourceAliyunDBInstanceUpdate(d *schema.ResourceData, meta interface{}) error {

	// todo: security_ip_list
	return resourceAliyunDBInstanceRead(d, meta)
}

func resourceAliyunDBInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*AliyunClient)

	instance, err := client.DescribeDBInstanceById(d.Id())
	if err != nil {
		if notFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error Describe DB InstanceAttribute: %#v", err)
	}

	d.Set("engine", instance.Engine)
	d.Set("engine_version", instance.EngineVersion)
	d.Set("db_instance_class", instance.DBInstanceClass)
	d.Set("db_instance_storage", instance.DBInstanceStorage)
	d.Set("zone_id", instance.ZoneId)
	d.Set("db_instance_net_type", instance.DBInstanceNetType)
	d.Set("security_ip_list", instance.SecurityIPList)

	return nil
}

func resourceAliyunDBInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).rdsconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := conn.DeleteInstance(d.Id())

		if err != nil {
			return resource.RetryableError(fmt.Errorf("DB Instance in use - trying again while it is deleted."))
		}

		args := &rds.DescribeDBInstancesArgs{
			DBInstanceId: d.Id(),
		}
		resp, err := conn.DescribeDBInstanceAttribute(args)
		if err != nil {
			return resource.NonRetryableError(err)
		} else if len(resp.Items.DBInstanceAttribute) < 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Vpc in use - trying again while it is deleted."))
	})
}

func buildDBCreateOrderArgs(d *schema.ResourceData, meta interface{}) (*rds.CreateOrderArgs, error) {
	args := &rds.CreateOrderArgs{
		RegionId:          getRegion(d, meta),
		EngineVersion:     d.Get("engine_version").(string),
		Engine:            rds.Engine(d.Get("engine").(string)),
		DBInstanceStorage: d.Get("db_instance_storage").(int),
		DBInstanceClass:   d.Get("db_instance_class").(string),
		Quantity:          DEFAULT_INSTANCE_COUNT,
		Resource:          rds.DefaultResource,
	}

	bussStr, err := json.Marshal(DefaultBusinessInfo)
	if err != nil {
		log.Printf("Failed to translate bussiness info %#v from json to string", DefaultBusinessInfo)
	}

	args.BusinessInfo = string(bussStr)

	chargeType := d.Get("instance_charge_type").(string)
	if chargeType != "" {
		args.PayType = rds.DBPayType(chargeType)
	}

	commodityCode := d.Get("commodity_code").(string)
	// if charge type is postpaid, the commodity code must set to bards
	if commodityCode == string(rds.Rds) && chargeType == string(rds.Postpaid) {
		args.CommodityCode = rds.Bards
	} else {
		args.CommodityCode = rds.CommodityCode(commodityCode)
	}

	//zoneId := d.Get("zone_id").(string)
	//multiAZ := d.Get("multi_az").(string)
	//allocatePublicCon := d.Get("allocate_public_connection").(string)
	//connectionMode := d.Get("connection_mode").(string)
	//vswitchId := d.Get("vswitch_id").(string)
	//
	//masterUserName := d.Get("master_user_name").(string)
	//masterUserPwd := d.Get("master_user_password").(string)
	//
	//backupPeriod := d.Get("preferred_backup_period").(string)
	//backupTime := d.Get("preferred_backup_time").(string)
	//retentionPeriod := d.Get("backup_retention_period").(string)
	//
	//securityIpList := d.Get("security_ip_list").(string)
	//
	//dbMapping := d.Get("db_mappings").(string)

	if v := d.Get("db_instance_net_type").(string); v != "" {
		args.DBInstanceNetType = common.NetType(v)
	}

	//period := d.Get("period").(string)
	// if charge_type == postpaid, then auto_pay = true
	// else charge_type == prepaid, then auto_pay default value is false
	autoPay := strconv.FormatBool(d.Get("auto_pay").(bool))
	if chargeType == string(rds.Prepaid) && autoPay == "" {
		args.AutoPay = strconv.FormatBool(false)
	} else {
		args.AutoPay = autoPay
	}

	// todo: deal commodity_code by charge_type

	return args, nil
}
