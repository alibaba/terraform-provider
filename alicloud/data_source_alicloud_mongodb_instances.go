package alicloud

import (
	"regexp"
	"strconv"

	"github.com/alibaba/terraform-provider/alicloud/connectivity"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudMongoDBInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudMongoDBInstancesRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateNameRegex,
			},
			"db_type": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validateAllowedStringValue([]string{
					"sharding",
					"replicate",
				}),
				Default: "replicate",
			},
			"tags": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateJsonString,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"replication_factor": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"last_downgrade_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudMongoDBInstancesRead(d *schema.ResourceData, meta interface{}) error {
	aliyunClient := meta.(*connectivity.AliyunClient)
	client := MongoDBService{aliyunClient}
	request := CommonRequestInit(aliyunClient.RegionId, MONGODBCode, MongoDBDomain)
	request.QueryParams["DBInstanceType"] = d.Get("db_type").(string)
	request.QueryParams["PageSize"] = strconv.Itoa(PageSizeLarge)
	request.QueryParams["PageNumber"] = "1"
	request.RegionId = aliyunClient.RegionId

	var mdb []MongoDBInstance
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		if r, err := regexp.Compile(v.(string)); err == nil {
			nameRegex = r
		}
	}

	for {
		resp, err := client.DescribeMongoDBInstances(request, aliyunClient)

		if err != nil {
			return err
		}

		if resp == nil || len(resp.Items.DBInstances) < 1 {
			break
		}

		for _, item := range resp.Items.DBInstances {
			if nameRegex != nil {
				if !nameRegex.MatchString(item.DBInstanceDescription) {
					continue
				}
			}
			mdb = append(mdb, item)
		}

		if len(resp.Items.DBInstances) < PageSizeLarge {
			break
		}
		pageNum := request.QueryParams["PageNumber"]
		i, _ := strconv.Atoi(pageNum)
		i++
		request.QueryParams["PageNumber"] = strconv.Itoa(i)
	}

	return mongoDBInstancesDescription(d, mdb)
}

func mongoDBInstancesDescription(d *schema.ResourceData, mdb []MongoDBInstance) error {
	var ids []string
	var s []map[string]interface{}

	for _, item := range mdb {
		mapping := map[string]interface{}{
			"id":                 item.DBInstanceID,
			"replication_factor": item.ReplicationFactor,
			"description":        item.DBInstanceDescription,
			"region_id":          item.RegionID,
			"zone_id":            item.ZoneID,
			"engine":             item.Engine,
			"engine_version":     item.EngineVersion,
			"instance_class":     item.DBInstanceClass,
			"instance_storage":   strconv.Itoa(item.DBInstanceStorage),
			"status":             item.DBInstanceStatus,
			"charge_type":        item.ChargeType,
			"network_type":       item.NetworkType,
			"creation_time":      item.CreationTime,
			"expire_time":        item.ExpireTime,
			"instance_type":      item.DBInstanceType,
		}

		ids = append(ids, item.DBInstanceID)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
