package alicloud

import (
	"fmt"
	"strings"

	"time"

	"strconv"

	"github.com/alibaba/terraform-provider/alicloud/aliyunclient"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func JudgeRegionValidation(key, region string, client *aliyunclient.AliyunClient) error {
	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeRegions(ecs.CreateDescribeRegionsRequest())
	})
	if err != nil {
		return fmt.Errorf("DescribeRegions got an error: %#v", err)
	}
	resp := raw.(*ecs.DescribeRegionsResponse)
	if resp == nil || len(resp.Regions.Region) < 1 {
		return GetNotFoundErrorFromString("There is no any available region.")
	}

	var rs []string
	for _, v := range resp.Regions.Region {
		if v.RegionId == region {
			return nil
		}
		rs = append(rs, v.RegionId)
	}
	return fmt.Errorf("'%s' is invalid. Expected on %v.", key, strings.Join(rs, ", "))
}

// DescribeZone validate zoneId is valid in region
func DescribeZone(zoneID string, client *aliyunclient.AliyunClient) (zone ecs.Zone, err error) {
	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeZones(ecs.CreateDescribeZonesRequest())
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeZonesResponse)
	if resp == nil || len(resp.Zones.Zone) < 1 {
		return zone, fmt.Errorf("There is no any availability zone in region %s.", client.RegionId)
	}

	zoneIds := []string{}
	for _, z := range resp.Zones.Zone {
		if z.ZoneId == zoneID {
			return z, nil
		}
		zoneIds = append(zoneIds, z.ZoneId)
	}
	return zone, fmt.Errorf("availability_zone not exists in range %s, all zones are %s", client.RegionId, zoneIds)
}

func DescribeInstanceById(id string, client *aliyunclient.AliyunClient) (instance ecs.Instance, err error) {
	req := ecs.CreateDescribeInstancesRequest()
	req.InstanceIds = convertListToJsonString([]interface{}{id})

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstances(req)
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeInstancesResponse)
	if resp == nil || len(resp.Instances.Instance) < 1 {
		return instance, GetNotFoundErrorFromString(GetNotFoundMessage("Instance", id))
	}

	return resp.Instances.Instance[0], nil
}

func DescribeInstanceAttribute(id string, client *aliyunclient.AliyunClient) (instance ecs.DescribeInstanceAttributeResponse, err error) {
	req := ecs.CreateDescribeInstanceAttributeRequest()
	req.InstanceId = id

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeInstanceAttribute(req)
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeInstanceAttributeResponse)
	if resp == nil {
		return instance, GetNotFoundErrorFromString(GetNotFoundMessage("Instance", id))
	}

	return *resp, nil
}

func QueryInstanceSystemDisk(id string, client *aliyunclient.AliyunClient) (disk ecs.Disk, err error) {
	args := ecs.CreateDescribeDisksRequest()
	args.InstanceId = id
	args.DiskType = string(DiskTypeSystem)

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(args)
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeDisksResponse)
	if resp != nil && len(resp.Disks.Disk) < 1 {
		return disk, GetNotFoundErrorFromString(fmt.Sprintf("The specified system disk is not found by instance id %s.", id))
	}

	return resp.Disks.Disk[0], nil
}

// ResourceAvailable check resource available for zone
func ResourceAvailable(zone ecs.Zone, resourceType ResourceType, client *aliyunclient.AliyunClient) error {
	for _, res := range zone.AvailableResourceCreation.ResourceTypes {
		if res == string(resourceType) {
			return nil
		}
	}
	return fmt.Errorf("%s is not available in %s zone of %s region", resourceType, zone.ZoneId, client.Region)
}

func DiskAvailable(zone ecs.Zone, diskCategory DiskCategory, client *aliyunclient.AliyunClient) error {
	for _, disk := range zone.AvailableDiskCategories.DiskCategories {
		if disk == string(diskCategory) {
			return nil
		}
	}
	return fmt.Errorf("%s is not available in %s zone of %s region", diskCategory, zone.ZoneId, client.Region)
}

func JoinSecurityGroups(instanceId string, securityGroupIds []string, client *aliyunclient.AliyunClient) error {
	req := ecs.CreateJoinSecurityGroupRequest()
	req.InstanceId = instanceId
	for _, sid := range securityGroupIds {
		req.SecurityGroupId = sid
		_, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.JoinSecurityGroup(req)
		})
		if err != nil && IsExceptedErrors(err, []string{InvalidInstanceIdAlreadyExists}) {
			return err
		}
	}

	return nil
}

func LeaveSecurityGroups(instanceId string, securityGroupIds []string, client *aliyunclient.AliyunClient) error {
	req := ecs.CreateLeaveSecurityGroupRequest()
	req.InstanceId = instanceId
	for _, sid := range securityGroupIds {
		req.SecurityGroupId = sid
		_, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.LeaveSecurityGroup(req)
		})
		if err != nil && IsExceptedErrors(err, []string{InvalidSecurityGroupIdNotFound}) {
			return err
		}
	}

	return nil
}

func DescribeSecurityGroupAttribute(securityGroupId string, client *aliyunclient.AliyunClient) (group ecs.DescribeSecurityGroupAttributeResponse, err error) {
	args := ecs.CreateDescribeSecurityGroupAttributeRequest()
	args.SecurityGroupId = securityGroupId

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(args)
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if resp == nil {
		return group, GetNotFoundErrorFromString(GetNotFoundMessage("Security Group", securityGroupId))
	}

	return *resp, nil
}

func DescribeSecurityGroupRule(groupId, direction, ipProtocol, portRange, nicType, cidr_ip, policy string, priority int, client *aliyunclient.AliyunClient) (rule ecs.Permission, err error) {
	args := ecs.CreateDescribeSecurityGroupAttributeRequest()
	args.SecurityGroupId = groupId
	args.Direction = direction
	args.NicType = nicType

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeSecurityGroupAttribute(args)
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeSecurityGroupAttributeResponse)
	if resp == nil {
		return rule, GetNotFoundErrorFromString(GetNotFoundMessage("Security Group", groupId))
	}

	for _, ru := range resp.Permissions.Permission {
		if strings.ToLower(string(ru.IpProtocol)) == ipProtocol && ru.PortRange == portRange {
			cidr := ru.SourceCidrIp
			if direction == string(DirectionIngress) && cidr == "" {
				cidr = ru.SourceGroupId
			}
			if direction == string(DirectionEgress) {
				if cidr = ru.DestCidrIp; cidr == "" {
					cidr = ru.DestGroupId
				}
			}

			if cidr == cidr_ip && strings.ToLower(string(ru.Policy)) == policy && ru.Priority == strconv.Itoa(priority) {
				return ru, nil
			}
		}
	}

	return rule, GetNotFoundErrorFromString(fmt.Sprintf("Security group rule not found by group id %s.", groupId))

}

func DescribeAvailableResources(d *schema.ResourceData, meta interface{}, destination DestinationResource, client *aliyunclient.AliyunClient) (zoneId string, validZones []ecs.AvailableZone, err error) {
	// Before creating resources, check input parameters validity according available zone.
	// If availability zone is nil, it will return all of supported resources in the current.
	args := ecs.CreateDescribeAvailableResourceRequest()
	args.DestinationResource = string(destination)
	args.IoOptimized = string(IOOptimized)

	if v, ok := d.GetOk("availability_zone"); ok && strings.TrimSpace(v.(string)) != "" {
		zoneId = strings.TrimSpace(v.(string))
	} else if v, ok := d.GetOk("vswitch_id"); ok && strings.TrimSpace(v.(string)) != "" {
		if vsw, err := DescribeVswitch(strings.TrimSpace(v.(string)), client); err == nil {
			zoneId = vsw.ZoneId
		}
	}

	if v, ok := d.GetOk("instance_charge_type"); ok && strings.TrimSpace(v.(string)) != "" {
		args.InstanceChargeType = strings.TrimSpace(v.(string))
	}

	if v, ok := d.GetOk("spot_strategy"); ok && strings.TrimSpace(v.(string)) != "" {
		args.SpotStrategy = strings.TrimSpace(v.(string))
	}

	if v, ok := d.GetOk("network_type"); ok && strings.TrimSpace(v.(string)) != "" {
		args.NetworkCategory = strings.TrimSpace(v.(string))
	}

	if v, ok := d.GetOk("is_outdated"); ok && v.(bool) == true {
		args.IoOptimized = string(NoneOptimized)
	}

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeAvailableResource(args)
	})
	if err != nil {
		return "", nil, fmt.Errorf("Error DescribeAvailableResource: %#v", err)
	}
	resources := raw.(*ecs.DescribeAvailableResourceResponse)

	if resources == nil || len(resources.AvailableZones.AvailableZone) < 1 {
		err = fmt.Errorf("There are no availability resources in the region: %s.", meta.(*aliyunclient.AliyunClient).RegionId)
		return
	}

	valid := false
	soldout := false
	var expectedZones []string
	for _, zone := range resources.AvailableZones.AvailableZone {
		if zone.Status == string(SoldOut) {
			if zone.ZoneId == zoneId {
				soldout = true
			}
			continue
		}
		if zoneId != "" && zone.ZoneId == zoneId {
			valid = true
			validZones = append(make([]ecs.AvailableZone, 1), zone)
			break
		}
		expectedZones = append(expectedZones, zone.ZoneId)
		validZones = append(validZones, zone)
	}
	if zoneId != "" {
		if !valid {
			err = fmt.Errorf("Availability zone %s status is not available in the region %s. Expected availability zones: %s.",
				zoneId, meta.(*aliyunclient.AliyunClient).RegionId, strings.Join(expectedZones, ", "))
			return
		}
		if soldout {
			err = fmt.Errorf("Availability zone %s status is sold out in the region %s. Expected availability zones: %s.",
				zoneId, meta.(*aliyunclient.AliyunClient).RegionId, strings.Join(expectedZones, ", "))
			return
		}
	}

	if len(validZones) <= 0 {
		err = fmt.Errorf("There is no availability resources in the region %s. Please choose another region.", meta.(*aliyunclient.AliyunClient).RegionId)
		return
	}

	return
}

func InstanceTypeValidation(targetType, zoneId string, validZones []ecs.AvailableZone, client *aliyunclient.AliyunClient) error {

	mapInstanceTypeToZones := make(map[string]string)
	var expectedInstanceTypes []string
	for _, zone := range validZones {
		if zoneId != "" && zoneId != zone.ZoneId {
			continue
		}
		for _, r := range zone.AvailableResources.AvailableResource {
			if r.Type == string(InstanceTypeResource) {
				for _, t := range r.SupportedResources.SupportedResource {
					if t.Status == string(SoldOut) {
						continue
					}
					if targetType == t.Value {
						return nil
					}

					if _, ok := mapInstanceTypeToZones[t.Value]; !ok {
						expectedInstanceTypes = append(expectedInstanceTypes, t.Value)
						mapInstanceTypeToZones[t.Value] = t.Value
					}
				}
			}
		}
	}
	if zoneId != "" {
		return fmt.Errorf("The instance type %s is solded out or is not supported in the zone %s. Expected instance types: %s", targetType, zoneId, strings.Join(expectedInstanceTypes, ", "))
	}
	return fmt.Errorf("The instance type %s is solded out or is not supported in the region %s. Expected instance types: %s", targetType, client.RegionId, strings.Join(expectedInstanceTypes, ", "))
}

func QueryInstancesWithKeyPair(instanceIdsStr, keypair string, client *aliyunclient.AliyunClient) (instanceIds []string, instances []ecs.Instance, err error) {

	args := ecs.CreateDescribeInstancesRequest()
	args.PageSize = requests.NewInteger(PageSizeLarge)
	args.PageNumber = requests.NewInteger(1)
	args.InstanceIds = instanceIdsStr
	args.KeyPairName = keypair
	for true {
		raw, e := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.DescribeInstances(args)
		})
		if e != nil {
			err = e
			return
		}
		resp := raw.(*ecs.DescribeInstancesResponse)
		if resp == nil || len(resp.Instances.Instance) < 0 {
			return
		}
		for _, inst := range resp.Instances.Instance {
			instanceIds = append(instanceIds, inst.InstanceId)
			instances = append(instances, inst)
		}
		if len(instances) < PageSizeLarge {
			break
		}
		if page, e := getNextpageNumber(args.PageNumber); e != nil {
			err = e
			return
		} else {
			args.PageNumber = page
		}
	}
	return
}

func DescribeKeyPair(keyName string, client *aliyunclient.AliyunClient) (keypair ecs.KeyPair, err error) {
	req := ecs.CreateDescribeKeyPairsRequest()
	req.KeyPairName = keyName
	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeKeyPairs(req)
	})

	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeKeyPairsResponse)
	if resp == nil || len(resp.KeyPairs.KeyPair) < 1 {
		return keypair, GetNotFoundErrorFromString(GetNotFoundMessage("KeyPair", keyName))
	}
	return resp.KeyPairs.KeyPair[0], nil

}

func DescribeDiskById(instanceId, diskId string, client *aliyunclient.AliyunClient) (disk ecs.Disk, err error) {
	req := ecs.CreateDescribeDisksRequest()
	if instanceId != "" {
		req.InstanceId = instanceId
	}
	req.DiskIds = convertListToJsonString([]interface{}{diskId})

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(req)
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeDisksResponse)
	if resp == nil || len(resp.Disks.Disk) < 1 {
		err = GetNotFoundErrorFromString(GetNotFoundMessage("ECS disk", diskId))
		return
	}
	return resp.Disks.Disk[0], nil
}

func DescribeDisksByType(instanceId string, diskType DiskType, client *aliyunclient.AliyunClient) (disk []ecs.Disk, err error) {
	req := ecs.CreateDescribeDisksRequest()
	if instanceId != "" {
		req.InstanceId = instanceId
	}
	req.DiskType = string(diskType)

	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeDisks(req)
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeDisksResponse)
	if resp == nil {
		return
	}
	return resp.Disks.Disk, nil
}

func DescribeTags(resourceId string, resourceType TagResourceType, client *aliyunclient.AliyunClient) (tags []ecs.Tag, err error) {
	req := ecs.CreateDescribeTagsRequest()
	req.ResourceType = string(resourceType)
	req.ResourceId = resourceId
	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeTags(req)
	})

	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeTagsResponse)
	if resp == nil || len(resp.Tags.Tag) < 1 {
		err = GetNotFoundErrorFromString(fmt.Sprintf("Describe %s tag by id %s got an error.", resourceType, resourceId))
		return
	}

	return resp.Tags.Tag, nil
}

func DescribeImageById(id string, client *aliyunclient.AliyunClient) (image ecs.Image, err error) {
	req := ecs.CreateDescribeImagesRequest()
	req.ImageId = id
	raw, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.DescribeImages(req)
	})
	if err != nil {
		return
	}
	resp := raw.(*ecs.DescribeImagesResponse)
	if resp == nil || len(resp.Images.Image) < 1 {
		return image, GetNotFoundErrorFromString(GetNotFoundMessage("Image", id))
	}
	return resp.Images.Image[0], nil
}

// WaitForInstance waits for instance to given status
func WaitForEcsInstance(instanceId string, status Status, timeout int, client *aliyunclient.AliyunClient) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := DescribeInstanceById(instanceId, client)
		if err != nil {
			return err
		}
		if instance.Status == string(status) {
			//Sleep one more time for timing issues
			time.Sleep(DefaultIntervalMedium * time.Second)
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("ECS Instance", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}

// WaitForInstance waits for instance to given status
func WaitForEcsDisk(diskId string, status Status, timeout int, client *aliyunclient.AliyunClient) error {
	if timeout <= 0 {
		timeout = DefaultTimeout
	}
	for {
		instance, err := DescribeDiskById("", diskId, client)
		if err != nil {
			return err
		}
		if instance.Status == string(status) {
			//Sleep one more time for timing issues
			time.Sleep(DefaultIntervalMedium * time.Second)
			break
		}
		timeout = timeout - DefaultIntervalShort
		if timeout <= 0 {
			return GetTimeErrorFromString(GetTimeoutMessage("ECS Disk", string(status)))
		}
		time.Sleep(DefaultIntervalShort * time.Second)

	}
	return nil
}

func AttachKeyPair(keyname string, instanceIds []interface{}, client *aliyunclient.AliyunClient) error {
	args := ecs.CreateAttachKeyPairRequest()
	args.KeyPairName = keyname
	args.InstanceIds = convertListToJsonString(instanceIds)
	return resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.RunSafelyWithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.AttachKeyPair(args)
		})

		if err != nil {
			if IsExceptedError(err, KeyPairServiceUnavailable) {
				return resource.RetryableError(fmt.Errorf("Attach Key Pair timeout and got an error: %#v.", err))
			}
			return resource.NonRetryableError(fmt.Errorf("Error Attach KeyPair: %#v", err))
		}
		return nil
	})
}
