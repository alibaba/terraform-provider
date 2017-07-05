package alicloud

import (
	"encoding/json"
	"fmt"
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
	"strings"
)

func (client *AliyunClient) DescribeImage(imageId string) (*ecs.ImageType, error) {

	pagination := common.Pagination{
		PageNumber: 1,
	}
	args := ecs.DescribeImagesArgs{
		Pagination: pagination,
		RegionId:   client.Region,
		Status:     ecs.ImageStatusAvailable,
	}

	var allImages []ecs.ImageType

	for {
		images, _, err := client.ecsconn.DescribeImages(&args)
		if err != nil {
			break
		}

		if len(images) == 0 {
			break
		}

		allImages = append(allImages, images...)

		args.Pagination.PageNumber++
	}

	if len(allImages) == 0 {
		return nil, common.GetClientErrorFromString("Not found")
	}

	var image *ecs.ImageType
	imageIds := []string{}
	for _, im := range allImages {
		if im.ImageId == imageId {
			image = &im
		}
		imageIds = append(imageIds, im.ImageId)
	}

	if image == nil {
		return nil, fmt.Errorf("image_id %s not exists in range %s, all images are %s", imageId, client.Region, imageIds)
	}

	return image, nil
}

// DescribeZone validate zoneId is valid in region
func (client *AliyunClient) DescribeZone(zoneID string) (*ecs.ZoneType, error) {
	zones, err := client.ecsconn.DescribeZones(client.Region)
	if err != nil {
		return nil, fmt.Errorf("error to list zones not found")
	}

	var zone *ecs.ZoneType
	zoneIds := []string{}
	for _, z := range zones {
		if z.ZoneId == zoneID {
			zone = &ecs.ZoneType{
				ZoneId:                    z.ZoneId,
				LocalName:                 z.LocalName,
				AvailableResourceCreation: z.AvailableResourceCreation,
				AvailableDiskCategories:   z.AvailableDiskCategories,
			}
		}
		zoneIds = append(zoneIds, z.ZoneId)
	}

	if zone == nil {
		return nil, fmt.Errorf("availability_zone not exists in range %s, all zones are %s", client.Region, zoneIds)
	}

	return zone, nil
}

// return multiIZ list of current region
func (client *AliyunClient) DescribeMultiIZByRegion() (izs []string, err error) {
	resp, err := client.rdsconn.DescribeRegions()
	if err != nil {
		return nil, fmt.Errorf("error to list regions not found")
	}
	regions := resp.Regions.RDSRegion

	zoneIds := []string{}
	for _, r := range regions {
		if r.RegionId == string(client.Region) && strings.Contains(r.ZoneId, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, r.ZoneId)
		}
	}

	return zoneIds, nil
}

func (client *AliyunClient) QueryInstancesByIds(ids []string) (instances []ecs.InstanceAttributesType, err error) {
	idsStr, jerr := json.Marshal(ids)
	if jerr != nil {
		return nil, jerr
	}

	args := ecs.DescribeInstancesArgs{
		RegionId:    client.Region,
		InstanceIds: string(idsStr),
	}

	instances, _, errs := client.ecsconn.DescribeInstances(&args)

	if errs != nil {
		return nil, errs
	}

	return instances, nil
}

func (client *AliyunClient) QueryInstancesById(id string) (instance *ecs.InstanceAttributesType, err error) {
	ids := []string{id}

	instances, errs := client.QueryInstancesByIds(ids)
	if errs != nil {
		return nil, errs
	}

	if len(instances) == 0 {
		return nil, GetNotFoundErrorFromString(InstanceNotFound)
	}

	return &instances[0], nil
}

func (client *AliyunClient) QueryInstanceSystemDisk(id string) (disk *ecs.DiskItemType, err error) {
	args := ecs.DescribeDisksArgs{
		RegionId:   client.Region,
		InstanceId: string(id),
		DiskType:   ecs.DiskTypeAllSystem,
	}
	disks, _, err := client.ecsconn.DescribeDisks(&args)
	if err != nil {
		return nil, err
	}
	if len(disks) == 0 {
		return nil, common.GetClientErrorFromString(SystemDiskNotFound)
	}

	return &disks[0], nil
}

// ResourceAvailable check resource available for zone
func (client *AliyunClient) ResourceAvailable(zone *ecs.ZoneType, resourceType ecs.ResourceType) error {
	available := false
	for _, res := range zone.AvailableResourceCreation.ResourceTypes {
		if res == resourceType {
			available = true
		}
	}
	if !available {
		return fmt.Errorf("%s is not available in %s zone of %s region", resourceType, zone.ZoneId, client.Region)
	}

	return nil
}

func (client *AliyunClient) DiskAvailable(zone *ecs.ZoneType, diskCategory ecs.DiskCategory) error {
	available := false
	for _, dist := range zone.AvailableDiskCategories.DiskCategories {
		if dist == diskCategory {
			available = true
		}
	}
	if !available {
		return fmt.Errorf("%s is not available in %s zone of %s region", diskCategory, zone.ZoneId, client.Region)
	}
	return nil
}

// todo: support syc
func (client *AliyunClient) JoinSecurityGroups(instanceId string, securityGroupIds []string) error {
	for _, sid := range securityGroupIds {
		err := client.ecsconn.JoinSecurityGroup(instanceId, sid)
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code != InvalidInstanceIdAlreadyExists {
				return err
			}
		}
	}

	return nil
}

func (client *AliyunClient) LeaveSecurityGroups(instanceId string, securityGroupIds []string) error {
	for _, sid := range securityGroupIds {
		err := client.ecsconn.LeaveSecurityGroup(instanceId, sid)
		if err != nil {
			e, _ := err.(*common.Error)
			if e.ErrorResponse.Code != InvalidSecurityGroupIdNotFound {
				return err
			}
		}
	}

	return nil
}

func (client *AliyunClient) DescribeSecurity(securityGroupId string) (*ecs.DescribeSecurityGroupAttributeResponse, error) {

	args := &ecs.DescribeSecurityGroupAttributeArgs{
		RegionId:        client.Region,
		SecurityGroupId: securityGroupId,
	}

	return client.ecsconn.DescribeSecurityGroupAttribute(args)
}

func (client *AliyunClient) DescribeSecurityByAttr(securityGroupId, direction, nicType string) (*ecs.DescribeSecurityGroupAttributeResponse, error) {

	args := &ecs.DescribeSecurityGroupAttributeArgs{
		RegionId:        client.Region,
		SecurityGroupId: securityGroupId,
		Direction:       direction,
		NicType:         ecs.NicType(nicType),
	}

	return client.ecsconn.DescribeSecurityGroupAttribute(args)
}

func (client *AliyunClient) DescribeSecurityGroupRule(securityGroupId, direction, nicType, ipProtocol, portRange string) (*ecs.PermissionType, error) {
	sg, err := client.DescribeSecurityByAttr(securityGroupId, direction, nicType)
	if err != nil {
		return nil, err
	}

	for _, p := range sg.Permissions.Permission {
		if strings.ToLower(string(p.IpProtocol)) == ipProtocol && p.PortRange == portRange {
			return &p, nil
		}
	}
	return nil, GetNotFoundErrorFromString("Security group rule not found")

}

func (client *AliyunClient) RevokeSecurityGroup(args *ecs.RevokeSecurityGroupArgs) error {
	//when the rule is not exist, api will return success(200)
	return client.ecsconn.RevokeSecurityGroup(args)
}

func (client *AliyunClient) RevokeSecurityGroupEgress(args *ecs.RevokeSecurityGroupEgressArgs) error {
	//when the rule is not exist, api will return success(200)
	return client.ecsconn.RevokeSecurityGroupEgress(args)
}

func (client *AliyunClient) CheckParameterValidity(d *schema.ResourceData, meta interface{}) error {
	// before creating resources, check input parameters validity according available zone.
	conn := meta.(*AliyunClient).ecsconn
	zones, err := conn.DescribeZones(getRegion(d, meta))

	if err != nil {
		return fmt.Errorf("Error DescribeZone: %#v", err)
	}

	if zones == nil || len(zones) < 1 {
		return fmt.Errorf("There are no availability zones in the region: %#v.", getRegion(d, meta))
	}

	zoneId := d.Get("availability_zone").(string)
	valid := false
	var validZones []string
	for _, zone := range zones {
		if zoneId != "" && zone.ZoneId == zoneId {
			valid = true
			break
		}
		validZones = append(validZones, zone.ZoneId)
	}
	if zoneId != "" && !valid {
		return fmt.Errorf("Availablity zone %s is not supported in the region %s. Expected availablity zones: %s.",
			zoneId, getRegion(d, meta), strings.Join(validZones, ", "))
	}

	var instanceType string
	if insType, ok := d.GetOk("available_instance_type"); ok {
		instanceType = insType.(string)
	} else if insType, ok := d.GetOk("instance_type"); ok {
		instanceType = insType.(string)
	}
	if instanceType != "" {
		// Retrieve series III instance type family
		families, err := client.FetchSpecifiedInstanceTypeFamily(getRegion(d, meta), zoneId, GenerationThree)

		if err != nil {
			return err
		}

		instanceTypeSplit := strings.Split(instanceType, DOT_SEPARATED)
		prefix := string(instanceTypeSplit[0] + DOT_SEPARATED + instanceTypeSplit[1])

		if _, ok := families[prefix]; !ok {
			var validFamilies []string
			for key := range families {
				validFamilies = append(validFamilies, key)
			}

			if zoneId == "" {
				return fmt.Errorf("Instance type %s is not supported in the region %s. Expected instance types: %s.",
					instanceType, getRegion(d, meta), strings.Join(validFamilies, ", "))
			}
			return fmt.Errorf("Instance type %s is not supported in the availability zone %s. Expected instance types: %s.",
				instanceType, zoneId, strings.Join(validFamilies, ", "))
		}
	}

	if instanceTypeFamily, ok := d.GetOk("instance_type_family"); ok {
		// Retrieve series III instance type family
		families, err := client.FetchSpecifiedInstanceTypeFamily(getRegion(d, meta), zoneId, GenerationThree)

		if err != nil {
			return err
		}

		if _, ok := families[instanceTypeFamily.(string)]; !ok {
			var validFamilies []string
			for key := range families {
				validFamilies = append(validFamilies, key)
			}
			if zoneId == "" {
				return fmt.Errorf("Instance type family %s is not supported in the region %s. Expected instance type families: %s.",
					instanceTypeFamily, getRegion(d, meta), strings.Join(validFamilies, ", "))
			}
			return fmt.Errorf("Instance type family %s is not supported in the availability zone %s. Expected instance type families: %s.",
				instanceTypeFamily, zoneId, strings.Join(validFamilies, ", "))
		}
	}
	return nil
}

func (client *AliyunClient) FetchSpecifiedInstanceTypeFamily(regionId common.Region, zoneId, generation string) (map[string]ecs.InstanceTypeFamily, error) {
	// Describe specified series instance type families
	response, err := client.ecsconn.DescribeInstanceTypeFamilies(&ecs.DescribeInstanceTypeFamiliesArgs{
		RegionId:   regionId,
		Generation: generation,
	})

	if err != nil {
		return nil, fmt.Errorf("Error DescribeInstanceTypeFamilies: %#v.", err)
	}

	familiesWithGeneration := make(map[string]ecs.InstanceTypeFamily)
	for _, family := range response.InstanceTypeFamilies.InstanceTypeFamily {
		familiesWithGeneration[family.InstanceTypeFamilyId] = family
	}

	// Filter specified zone's instance type families, and make them fit for specified generation
	if zoneId != "" {
		validFamilies := make(map[string]ecs.InstanceTypeFamily)
		zones, err := client.ecsconn.DescribeZones(regionId)
		if err != nil {
			return nil, fmt.Errorf("Error DescribeZones: %#v", err)
		}
		for _, zone := range zones {
			if zone.ZoneId == zoneId {
				for _, resource := range zone.AvailableResources.ResourcesInfo {
					families, _ := resource.InstanceTypeFamilies[string(ecs.SupportedInstanceTypeFamily)]
					for _, familyId := range families {
						if val, ok := familiesWithGeneration[familyId]; ok {
							validFamilies[familyId] = val
						}
					}

				}
				return validFamilies, nil
			}
		}
	}
	return familiesWithGeneration, nil
}
