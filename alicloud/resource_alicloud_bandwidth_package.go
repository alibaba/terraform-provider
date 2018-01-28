package alicloud

import (
	"fmt"
	"log"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceAliyunBandwidthPackage() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunBandwidthPackageCreate,
		Read:   resourceAliyunBandwidthPackageRead,
		Update: resourceAliyunBandwidthPackageUpdate,
		Delete: resourceAliyunBandwidthPackageDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"nat_gateway_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ip_count": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"bandwidth": &schema.Schema{
				Type:     schema.TypeInt,
				Required: true,
			},
			"zone": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_ip_addresses": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}
func resourceAliyunBandwidthPackageCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := &ecs.CreateBandwidthPackageArgs{
		RegionId:     getRegion(d, meta),
		NatGatewayId: d.Get("nat_gateway_id").(string),
		IpCount:      d.Get("ip_count").(int),
		Bandwidth:    d.Get("bandwidth").(int),
	}

	var zone string
	if v, ok := d.GetOk("zone"); ok {
		zone = v.(string)
	}
	args.Zone = zone

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}

	args.Name = name

	if v, ok := d.GetOk("description"); ok {
		args.Description = v.(string)
	}

	resp, err := conn.CreateBandwidthPackage(args)
	if err != nil {
		return fmt.Errorf("CreateBandwidthPackage got error: %#v", err)
	}

	d.SetId(resp.BandwidthPackageId)

	return resourceAliyunBandwidthPackageRead(d, meta)
}

func resourceAliyunBandwidthPackageRead(d *schema.ResourceData, meta interface{}) error {
	var bandwidthPackagesIds []string
	bandwidthPackagesIds = append(bandwidthPackagesIds, d.Id())
	bandwidthPackages, err := flattenBandWidthPackages(bandwidthPackagesIds, meta, d)
	if err != nil {
		log.Printf("[ERROR] bandWidthPackages flattenBandWidthPackages failed. bandwidth id is %#v", d.Id())
	} else {
		var bandWidthPackage = bandwidthPackages[0]
		d.Set("bandwidth", bandWidthPackage["bandwidth"])
		d.Set("ip_count", bandWidthPackage["ip_count"])
		d.Set("zone", bandWidthPackage["zone"])
		d.Set("name", bandWidthPackage["name"])
		d.Set("description", bandWidthPackage["description"])
		d.Set("nat_gateway_id", bandWidthPackage["nat_gateway_id"])
		d.Set("public_ip_addresses", bandWidthPackage["public_ip_addresses"])
	}

	return nil
}

func resourceAliyunBandwidthPackageUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	conn := client.vpcconn

	d.Partial(true)
	attributeUpdate := false
	args := &ecs.ModifyBandwithPackageAttributeArgs{
		RegionId:           getRegion(d, meta),
		BandwidthPackageId: d.Id(),
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		} else {
			return fmt.Errorf("can't change name to empty string")
		}
		args.Name = name

		attributeUpdate = true
	}

	if d.HasChange("description") {
		d.SetPartial("description")
		var description string
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		} else {
			return fmt.Errorf("can to change description to empty string")
		}

		args.Description = description

		attributeUpdate = true
	}

	if attributeUpdate {
		if err := conn.ModifyBandwithPackageAttribute(args); err != nil {
			return err
		}
	}

	if d.HasChange("bandwidth") {
		d.SetPartial("bandwidth")
		var bandwidth string
		if v, ok := d.GetOk("bandwidth"); ok {
			bandwidth = v.(string)
		} else {
			return fmt.Errorf("can to change bandwidth to empty string")
		}

		args := &ecs.ModifyBandwithPackageSpecArgs{
			RegionId:           getRegion(d, meta),
			BandwidthPackageId: d.Id(),
			Bandwidth:          bandwidth,
		}

		err := conn.ModifyBandwithPackageSpec(args)
		if err != nil {
			return fmt.Errorf("%#v %#v", err, *args)
		}

	}
	d.Partial(false)

	return resourceAliyunBandwidthPackageRead(d, meta)
}

func resourceAliyunBandwidthPackageDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
