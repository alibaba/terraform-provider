package alicloud

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"

	//"github.com/denverdino/aliyungo/ecs"
	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceAlicloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudZonesRead,

		Schema: map[string]*schema.Schema{
			"available_instance_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceType,
			},
			"available_resource_creation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ResourceTypeInstance),
					string(ResourceTypeRds),
					string(ResourceTypeRkv),
					string(ResourceTypeVSwitch),
					string(ResourceTypeDisk),
					string(IoOptimized),
				}),
			},
			"available_disk_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateDiskCategory,
			},

			"multi": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validateInstanceChargeType,
			},
			"network_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateAllowedStringValue([]string{string(Vpc), string(Classic)}),
			},
			"spot_strategy": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      NoSpot,
				ValidateFunc: validateInstanceSpotStrategy,
			},

			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// Computed values.
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"local_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available_resource_creation": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"available_disk_categories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlicloudZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*aliyunclient.AliyunClient)
	ecsService := EcsService{client}

	resType, _ := d.Get("available_resource_creation").(string)
	multi := d.Get("multi").(bool)
	var zoneIds []string
	rdsZones := make(map[string]string)
	rkvZones := make(map[string]string)
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeRds)) {
		request := rds.CreateDescribeRegionsRequest()
		raw, err := client.RunSafelyWithRdsClient(func(rdsClient *rds.Client) (interface{}, error) {
			return rdsClient.DescribeRegions(request)
		})
		if err != nil {
			return fmt.Errorf("[ERROR] DescribeRegions got an error: %#v", err)
		}
		regions, _ := raw.(*rds.DescribeRegionsResponse)
		if len(regions.Regions.RDSRegion) <= 0 {
			return fmt.Errorf("[ERROR] There is no available region for RDS.")
		}
		for _, r := range regions.Regions.RDSRegion {
			if multi && strings.Contains(r.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(meta.(*aliyunclient.AliyunClient).Region) {
				zoneIds = append(zoneIds, r.ZoneId)
				continue
			}
			rdsZones[r.ZoneId] = r.RegionId
		}
	}
	if strings.ToLower(Trim(resType)) == strings.ToLower(string(ResourceTypeRkv)) {
		request := r_kvstore.CreateDescribeRegionsRequest()
		raw, err := client.RunSafelyWithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
			return rkvClient.DescribeRegions(request)
		})
		if err != nil {
			return fmt.Errorf("[ERROR] DescribeRegions got an error: %#v", err)
		}
		regions, _ := raw.(*r_kvstore.DescribeRegionsResponse)
		if len(regions.RegionIds.KVStoreRegion) <= 0 {
			return fmt.Errorf("[ERROR] There is no available region for KVStore")
		}
		for _, r := range regions.RegionIds.KVStoreRegion {
			for _, zoneID := range r.ZoneIdList.ZoneId {
				if multi && strings.Contains(zoneID, MULTI_IZ_SYMBOL) && r.RegionId == string(meta.(*aliyunclient.AliyunClient).Region) {
					zoneIds = append(zoneIds, zoneID)
					continue
				}
				rkvZones[zoneID] = r.RegionId
			}
		}
	}
	if len(zoneIds) > 0 {
		sort.Strings(zoneIds)
		return multiZonesDescriptionAttributes(d, zoneIds)
	} else if multi {
		return fmt.Errorf("There is no multi zones in the current region %s. Please change region and try again.", meta.(*aliyunclient.AliyunClient).Region)
	}

	_, validZones, err := ecsService.DescribeAvailableResources(d, meta, ZoneResource)
	if err != nil {
		return err
	}

	req := ecs.CreateDescribeZonesRequest()
	if v, ok := d.GetOk("instance_charge_type"); ok && v.(string) != "" {
		req.InstanceChargeType = v.(string)
	}
	if v, ok := d.GetOk("spot_strategy"); ok && v.(string) != "" {
		req.SpotStrategy = v.(string)
	}

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(req)
	})
	if err != nil {
		return fmt.Errorf("DescribeZones got an error: %#v", err)
	}
	resp, _ := raw.(*ecs.DescribeZonesResponse)
	if resp == nil || len(resp.Zones.Zone) < 1 {
		return fmt.Errorf("There are no availability zones in the region: %#v.", meta.(*aliyunclient.AliyunClient).Region)
	}

	mapZones := make(map[string]ecs.Zone)
	insType, _ := d.Get("available_instance_type").(string)
	diskType, _ := d.Get("available_disk_category").(string)

	for _, zone := range resp.Zones.Zone {
		for _, v := range validZones {
			if zone.ZoneId != v.ZoneId {
				continue
			}
			if len(zone.AvailableInstanceTypes.InstanceTypes) <= 0 ||
				(insType != "" && !constraints(zone.AvailableInstanceTypes.InstanceTypes, insType)) {
				continue
			}
			if len(zone.AvailableDiskCategories.DiskCategories) <= 0 ||
				(diskType != "" && !constraints(zone.AvailableDiskCategories.DiskCategories, diskType)) {
				continue
			}
			if len(rdsZones) > 0 {
				if _, ok := rdsZones[zone.ZoneId]; !ok {
					continue
				}
			}
			if len(rkvZones) > 0 {
				if _, ok := rkvZones[zone.ZoneId]; !ok {
					continue
				}
			}
			zoneIds = append(zoneIds, zone.ZoneId)
			mapZones[zone.ZoneId] = zone
		}
	}

	if len(zoneIds) <= 0 {
		return fmt.Errorf("Your query zones returned no results. Please change your search criteria and try again.")
	}

	// Sort zones before reading
	sort.Strings(zoneIds)

	var s []map[string]interface{}
	for _, zoneId := range zoneIds {
		mapping := map[string]interface{}{
			"id":                          zoneId,
			"local_name":                  mapZones[zoneId].LocalName,
			"available_instance_types":    mapZones[zoneId].AvailableInstanceTypes.InstanceTypes,
			"available_resource_creation": mapZones[zoneId].AvailableResourceCreation.ResourceTypes,
			"available_disk_categories":   mapZones[zoneId].AvailableDiskCategories.DiskCategories,
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(zoneIds))
	if err := d.Set("zones", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

// check array constraints str
func constraints(arr interface{}, v string) bool {
	arrs := reflect.ValueOf(arr)
	len := arrs.Len()
	for i := 0; i < len; i++ {
		if arrs.Index(i).String() == v {
			return true
		}
	}
	return false
}

func multiZonesDescriptionAttributes(d *schema.ResourceData, zones []string) error {
	var s []map[string]interface{}
	for _, t := range zones {
		mapping := map[string]interface{}{
			"id": t,
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(zones))
	if err := d.Set("zones", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
