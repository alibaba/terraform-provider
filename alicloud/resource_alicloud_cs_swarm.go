package alicloud

import (
	"fmt"
	"time"

	"github.com/alibaba/terraform-provider/alicloud/connectivity"
	newsdk "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/denverdino/aliyungo/cs"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAlicloudCSSwarm() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlicloudCSSwarmCreate,
		Read:   resourceAlicloudCSSwarmRead,
		Update: resourceAlicloudCSSwarmUpdate,
		Delete: resourceAlicloudCSSwarmDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Computed:      true,
				ValidateFunc:  validateContainerName,
				ConflictsWith: []string{"name_prefix"},
			},
			"name_prefix": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "Terraform-Creation",
				ValidateFunc:  validateContainerNamePrefix,
				ConflictsWith: []string{"name"},
			},
			"size": &schema.Schema{
				Type:       schema.TypeInt,
				Optional:   true,
				Deprecated: "Field 'size' has been deprecated from provider version 1.9.1. New field 'node_number' replaces it.",
			},
			"node_number": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validateIntegerInRange(0, 50),
			},
			"cidr_block": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateInstanceType,
			},
			"vswitch_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": &schema.Schema{
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				Sensitive: true,
			},
			"disk_size": &schema.Schema{
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      20,
				ValidateFunc: validateIntegerInRange(20, 32768),
			},
			"disk_category": &schema.Schema{
				Type:         schema.TypeString,
				Optional:     true,
				Default:      ecs.DiskCategoryCloudEfficiency,
				ForceNew:     true,
				ValidateFunc: validateDiskCategory,
			},
			"image_id": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"release_eip": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return old != ""
				},
			},
			"is_outdated": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"need_slb": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
			"nodes": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"eip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"slb_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"security_group_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"agent_version": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAlicloudCSSwarmCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	ecsService := EcsService{client}
	vpcService := VpcService{client}

	// Ensure instance_type is valid
	zoneId, validZones, err := ecsService.DescribeAvailableResources(d, meta, InstanceTypeResource)
	if err != nil {
		return err
	}
	if err := ecsService.InstanceTypeValidation(d.Get("instance_type").(string), zoneId, validZones); err != nil {
		return err
	}

	var clusterName string
	if v, ok := d.GetOk("name"); ok {
		clusterName = v.(string)
	} else {
		clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
	}

	args := &cs.ClusterCreationArgs{
		Name:             clusterName,
		InstanceType:     d.Get("instance_type").(string),
		Password:         d.Get("password").(string),
		Size:             int64(d.Get("node_number").(int)),
		IOOptimized:      ecs.IoOptimized("true"),
		DataDiskCategory: ecs.DiskCategory(d.Get("disk_category").(string)),
		DataDiskSize:     int64(d.Get("disk_size").(int)),
		NetworkMode:      cs.VPCNetwork,
		VSwitchID:        d.Get("vswitch_id").(string),
		SubnetCIDR:       d.Get("cidr_block").(string),
		ReleaseEipFlag:   d.Get("release_eip").(bool),
		NeedSLB:          d.Get("need_slb").(bool),
	}

	vsw, err := vpcService.DescribeVswitch(args.VSwitchID)
	if err != nil {
		return fmt.Errorf("Error DescribeVSwitches: %#v", err)
	}

	if vsw.CidrBlock == args.SubnetCIDR {
		return fmt.Errorf("Container cluster's cidr_block only accepts 192.168.X.0/24 or 172.18.X.0/24 ~ 172.31.X.0/24. " +
			"And it cannot be equal to vswitch's cidr_block and sub cidr block.")
	}
	args.VPCID = vsw.VpcId

	if imageId, ok := d.GetOk("image_id"); ok {
		if _, err := ecsService.DescribeImageById(imageId.(string)); err != nil {
			return err
		}

		args.ECSImageID = imageId.(string)
	}

	region := client.Region
	raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
		return csClient.CreateCluster(region, args)
	})

	if err != nil {
		return fmt.Errorf("Creating container Cluster got an error: %#v", err)
	}
	cluster, _ := raw.(cs.ClusterCreationResponse)
	d.SetId(cluster.ClusterID)

	_, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
		return nil, csClient.WaitForClusterAsyn(cluster.ClusterID, cs.Running, 500)
	})

	if err != nil {
		return fmt.Errorf("Waitting for container Cluster %#v got an error: %#v", cs.Running, err)
	}

	return resourceAlicloudCSSwarmUpdate(d, meta)
}

func resourceAlicloudCSSwarmUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	d.Partial(true)
	if d.HasChange("node_number") && !d.IsNewResource() {
		o, n := d.GetChange("node_number")
		oi := o.(int)
		ni := n.(int)
		if ni <= oi {
			return fmt.Errorf("The node number must greater than the current. The cluster's current node number is %d.", oi)
		}
		d.SetPartial("node_number")
		_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.ResizeCluster(d.Id(), &cs.ClusterResizeArgs{
				Size:             int64(ni),
				InstanceType:     d.Get("instance_type").(string),
				Password:         d.Get("password").(string),
				DataDiskCategory: ecs.DiskCategory(d.Get("disk_category").(string)),
				DataDiskSize:     int64(d.Get("disk_size").(int)),
				ECSImageID:       d.Get("image_id").(string),
				IOOptimized:      ecs.IoOptimized("true"),
			})
		})
		if err != nil {
			return fmt.Errorf("Resize Cluster got an error: %#v", err)
		}

		_, err = client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.WaitForClusterAsyn(d.Id(), cs.Running, 500)
		})

		if err != nil {
			return fmt.Errorf("Waitting for container Cluster %#v got an error: %#v", cs.Running, err)
		}
	}

	if !d.IsNewResource() && (d.HasChange("name") || d.HasChange("name_prefix")) {
		var clusterName string
		if v, ok := d.GetOk("name"); ok {
			clusterName = v.(string)
		} else {
			clusterName = resource.PrefixedUniqueId(d.Get("name_prefix").(string))
		}
		_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.ModifyClusterName(d.Id(), clusterName)
		})
		if err != nil && !IsExceptedError(err, ErrorClusterNameAlreadyExist) {
			return fmt.Errorf("Modify Cluster Name got an error: %#v", err)
		}
		d.SetPartial("name")
		d.SetPartial("name_prefix")
	}

	d.Partial(false)

	return resourceAlicloudCSSwarmRead(d, meta)
}

func resourceAlicloudCSSwarmRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)
	csService := CsService{client}
	ecsService := EcsService{client}

	raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
		return csClient.DescribeCluster(d.Id())
	})

	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	cluster, _ := raw.(cs.ClusterType)
	d.Set("name", cluster.Name)
	d.Set("node_number", cluster.Size)
	d.Set("vpc_id", cluster.VPCID)
	d.Set("vswitch_id", cluster.VSwitchID)
	d.Set("security_group_id", cluster.SecurityGroupID)
	d.Set("slb_id", cluster.ExternalLoadbalancerID)
	d.Set("agent_version", cluster.AgentVersion)

	pcluster, certs, err := csService.GetContainerClusterAndCertsByName(cluster.Name)
	if err != nil {
		return err
	}
	raw, err = client.WithCsProjectClient(pcluster.ClusterID, pcluster.MasterURL, *certs, func(csProjectClient *cs.ProjectClient) (interface{}, error) {
		return csProjectClient.GetSwarmClusterNodes()
	})
	if err != nil {
		return err
	}
	resp, _ := raw.(cs.GetSwarmClusterNodesResponse)
	var nodes []map[string]interface{}
	var oneNode newsdk.Instance

	for _, node := range resp {
		mapping := map[string]interface{}{
			"id":         node.InstanceId,
			"name":       node.Name,
			"private_ip": node.IP,
			"status":     node.Status,
		}
		if inst, err := ecsService.DescribeInstanceById(node.InstanceId); err != nil {
			return fmt.Errorf("[ERROR] QueryInstancesById %s: %#v.", node.InstanceId, err)
		} else {
			mapping["eip"] = inst.EipAddress.IpAddress
			oneNode = inst
		}

		nodes = append(nodes, mapping)
	}

	d.Set("nodes", nodes)

	d.Set("instance_type", oneNode.InstanceType)
	if disks, err := ecsService.DescribeDisksByType(oneNode.InstanceId, DiskTypeData); err != nil {
		return fmt.Errorf("[ERROR] DescribeDisks By Id %s: %#v.", resp[0].InstanceId, err)
	} else {
		for _, disk := range disks {
			d.Set("disk_size", disk.Size)
			d.Set("disk_category", disk.Category)
		}
	}

	return nil
}

func resourceAlicloudCSSwarmDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AliyunClient)

	return resource.Retry(3*time.Minute, func() *resource.RetryError {
		_, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return nil, csClient.DeleteCluster(d.Id())
		})
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return resource.RetryableError(fmt.Errorf("Deleting container cluster got an error: %#v", err))
		}

		raw, err := client.WithCsClient(func(csClient *cs.Client) (interface{}, error) {
			return csClient.DescribeCluster(d.Id())
		})
		if err != nil {
			if NotFoundError(err) || IsExceptedError(err, ErrorClusterNotFound) {
				return nil
			}
			return resource.NonRetryableError(fmt.Errorf("Describe container cluster got an error: %#v", err))
		}
		resp, _ := raw.(cs.ClusterType)
		if resp.ClusterID == "" {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("Deleting container cluster got an error: %#v", err))
	})
}
