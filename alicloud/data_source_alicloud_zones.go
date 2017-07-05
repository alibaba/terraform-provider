package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"reflect"
	"strings"
)

func dataSourceAlicloudZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlicloudZonesRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"available_instance_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"available_resource_creation": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"available_disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validateAllowedStringValue([]string{
					string(ecs.DiskCategoryCloudSSD),
					string(ecs.DiskCategoryCloudEfficiency),
				}),
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
	conn := meta.(*AliyunClient).ecsconn

	insType, _ := d.Get("available_instance_type").(string)
	resType, _ := d.Get("available_resource_creation").(string)
	diskType, _ := d.Get("available_disk_category").(string)

	resp, err := conn.DescribeZones(getRegion(d, meta))
	if err != nil {
		return err
	}

	familiesWithGeneration, err := meta.(*AliyunClient).FetchSpecifiedInstanceTypeFamily(getRegion(d, meta), "", GenerationThree)
	if err != nil {
		return err
	}
	var zoneTypes []ecs.ZoneType
	for _, zone := range resp {

		if len(zone.AvailableInstanceTypes.InstanceTypes) == 0 {
			continue
		}

		if insType != "" {
			if !constraints(zone.AvailableInstanceTypes.InstanceTypes, insType) {
				continue
			}
			// Ensure current instance type belong to series III
			instanceTypeSplit := strings.Split(insType, DOT_SEPARATED)
			prefix := string(instanceTypeSplit[0] + DOT_SEPARATED + instanceTypeSplit[1])
			if _, ok := familiesWithGeneration[prefix]; !ok {
				continue
			}
		}

		if len(zone.AvailableResourceCreation.ResourceTypes) == 0 || (resType != "" && !constraints(zone.AvailableResourceCreation.ResourceTypes, resType)) {
			continue
		}

		if len(zone.AvailableDiskCategories.DiskCategories) == 0 || (diskType != "" && !constraints(zone.AvailableDiskCategories.DiskCategories, diskType)) {
			continue
		}

		// Filter and find supported resource types after finding valid zones
		var vaildInstanceTypes []string
		var vaildDiskCategories []ecs.DiskCategory
		for _, typeItem := range zone.AvailableInstanceTypes.InstanceTypes {
			instanceTypeSplit := strings.Split(typeItem, DOT_SEPARATED)
			prefix := string(instanceTypeSplit[0] + DOT_SEPARATED + instanceTypeSplit[1])
			if _, ok := familiesWithGeneration[prefix]; ok {
				vaildInstanceTypes = append(vaildInstanceTypes, typeItem)
			}
		}
		for _, diskItem := range zone.AvailableDiskCategories.DiskCategories {
			if ecs.DiskCategory(diskItem) == ecs.DiskCategoryCloudEfficiency ||
				ecs.DiskCategory(diskItem) == ecs.DiskCategoryCloudSSD {
				vaildDiskCategories = append(vaildDiskCategories, ecs.DiskCategory(diskItem))
			}
		}
		zone.AvailableInstanceTypes.InstanceTypes = vaildInstanceTypes
		zone.AvailableDiskCategories.DiskCategories = vaildDiskCategories
		zoneTypes = append(zoneTypes, zone)
	}

	if len(zoneTypes) < 1 {
		return fmt.Errorf("Your query returned no results. Please change your search criteria and try again.")
	}

	log.Printf("[DEBUG] alicloud_zones - Zones found: %#v", zoneTypes)
	return zonesDescriptionAttributes(d, zoneTypes)
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

func zonesDescriptionAttributes(d *schema.ResourceData, types []ecs.ZoneType) error {
	var ids []string
	var s []map[string]interface{}
	for _, t := range types {
		mapping := map[string]interface{}{
			"id":                          t.ZoneId,
			"local_name":                  t.LocalName,
			"available_instance_types":    t.AvailableInstanceTypes.InstanceTypes,
			"available_resource_creation": t.AvailableResourceCreation.ResourceTypes,
			"available_disk_categories":   t.AvailableDiskCategories.DiskCategories,
		}

		log.Printf("[DEBUG] alicloud_zones - adding zone mapping: %v", mapping)
		ids = append(ids, t.ZoneId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("zones", s); err != nil {
		return err
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output != nil {
		writeToFile(output.(string), s)
	}

	return nil
}
