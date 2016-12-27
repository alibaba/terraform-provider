package alicloud

import (
	"fmt"

	"github.com/denverdino/aliyungo/common"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"time"
)

func resourceAliyunNatGateway() *schema.Resource {
	return &schema.Resource{
		Create: resourceAliyunNatGatewayCreate,
		Read:   resourceAliyunNatGatewayRead,
		Update: resourceAliyunNatGatewayUpdate,
		Delete: resourceAliyunNatGatewayDelete,

		Schema: map[string]*schema.Schema{
			"vpc_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"spec": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"bandwidth_package_id": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"bandwidth_packages": &schema.Schema{
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						},
					},
				},
				Required: true,
			},
		},
	}
}

func resourceAliyunNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := &CreateNatGatewayArgs{
		RegionId: getRegion(d, meta),
		VpcId:    d.Get("vpc_id").(string),
		Spec:     d.Get("spec").(string),
	}

	bandwidthPackages := d.Get("bandwidth_packages").([]interface{})

	if len(bandwidthPackages) > 4 {
		return fmt.Errorf("Only less than 4 bandwidth packages form per NatGateway is supported")
	}

	bandwidthPackageTypes := []BandwidthPackageType{}

	for _, e := range bandwidthPackages {
		pack := e.(map[string]interface{})
		bandwidthPackage := BandwidthPackageType{
			IpCount:   pack["ip_count"].(int),
			Bandwidth: pack["bandwidth"].(int),
		}
		if pack["zone"].(string) != "" {
			bandwidthPackage.Zone = pack["zone"].(string)
		}

		bandwidthPackageTypes = append(bandwidthPackageTypes, bandwidthPackage)
	}

	args.BandwidthPackage = bandwidthPackageTypes

	var name string
	if v, ok := d.GetOk("name"); ok {
		name = v.(string)
	}
	// } else {
	// 	name = resource.PrefixedUniqueId("tf-ngw-")
	// 	d.Set("name", name)
	// }

	args.Name = name

	if v, ok := d.GetOk("description"); ok {
		args.Description = v.(string)
	}

	resp, err := CreateNatGateway(conn, args)
	if err != nil {
		return fmt.Errorf("CreateNatGateway got error: %#v", err)
	}

	d.SetId(resp.NatGatewayId)
	d.Partial(true)
	d.SetPartial("name")
	d.SetPartial("description")
	d.SetPartial("spec")
	d.SetPartial("vpc_id")

	// for i, packageId := range resp.BandwidthPackageIds.BandwidthPackageId {
	// 	packages[i].(map[string]interface{})["id"] = packageId
	// }

	d.SetPartial("bandwidth_packages")

	//d.Set("bandwidth_package_id", resp.BandwidthPackageIds.BandwidthPackageId[0])
	d.SetPartial("bandwidth_package_id")

	return resourceAliyunNatGatewayRead(d, meta)
}

func resourceAliyunNatGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	natGateway, err := client.DescribeNatGateway(d.Id())
	if err != nil {
		if notFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("name", natGateway.Name)
	d.Set("spec", natGateway.Spec)
	d.Set("bandwidth_package_ids", natGateway.BandwidthPackageIds.BandwidthPackageId)
	d.Set("description", natGateway.Description)

	return nil
}

func resourceAliyunNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	natGateway, err := client.DescribeNatGateway(d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)

	if d.HasChange("name") {
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		} else {
			return fmt.Errorf("can to change name to empty string")
		}

		args := &ModifyNatGatewayAttributeArgs{
			RegionId:     natGateway.RegionId,
			NatGatewayId: natGateway.NatGatewayId,
			Name:         name,
		}

		err := ModifyNatGatewayAttribute(client.vpcconn, args)
		if err != nil {
			return err
		}

		d.SetPartial("name")
	}

	if d.HasChange("description") {
		var description string
		if v, ok := d.GetOk("description"); ok {
			description = v.(string)
		} else {
			return fmt.Errorf("can to change description to empty string")
		}

		args := &ModifyNatGatewayAttributeArgs{
			RegionId:     natGateway.RegionId,
			NatGatewayId: natGateway.NatGatewayId,
			Description:  description,
		}

		err := ModifyNatGatewayAttribute(client.vpcconn, args)
		if err != nil {
			return fmt.Errorf("%#v %#v", err, *args)
		}

		d.SetPartial("description")
	}

	if d.HasChange("spec") {
		var spec NatGatewaySpec
		if v, ok := d.GetOk("spec"); ok {
			spec = NatGatewaySpec(v.(string))
		} else {
			// set default to small spec
			spec = NatGatewaySmallSpec
		}

		args := &ModifyNatGatewaySpecArgs{
			RegionId:     natGateway.RegionId,
			NatGatewayId: natGateway.NatGatewayId,
			Spec:         spec,
		}

		err := ModifyNatGatewaySpec(client.vpcconn, args)
		if err != nil {
			return fmt.Errorf("%#v %#v", err, *args)
		}

		d.SetPartial("spec")
	}

	return nil
}

func resourceAliyunNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)

	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		packages, err := DescribeBandwidthPackages(client.vpcconn, &DescribeBandwidthPackagesArgs{
			RegionId:     getRegion(d, meta),
			NatGatewayId: d.Id(),
		})
		if err != nil {
			log.Printf("[ERROR] Describe bandwidth package is failed, natGateway Id: %s", d.Id())
			return resource.NonRetryableError(err)
		}

		retry := false
		for _, pack := range packages {
			err = DeleteBandwidthPackage(client.vpcconn, &DeleteBandwidthPackageArgs{
				RegionId:           getRegion(d, meta),
				BandwidthPackageId: pack.BandwidthPackageId,
			})

			if err != nil {
				er, _ := err.(*common.Error)
				if er.ErrorResponse.Code == "Invalid.RegionId" {
					log.Printf("[ERROR] Delete bandwidth package is failed, bandwidthPackageId: %#v", pack.BandwidthPackageId)
					return resource.NonRetryableError(err)
				}
				retry = true
			}
		}

		if retry {
			return resource.RetryableError(fmt.Errorf("Bandwidth package in use - trying again while it is deleted."))
		}

		args := &DeleteNatGatewayArgs{
			RegionId:     client.Region,
			NatGatewayId: d.Id(),
		}

		err = DeleteNatGateway(client.vpcconn, args)
		if err != nil {
			er, _ := err.(*common.Error)
			if er.ErrorResponse.Code == "DependencyViolation.BandwidthPackages" {
				return resource.RetryableError(fmt.Errorf("NatGateway in use - trying again while it is deleted."))
			}
		}

		describeArgs := &DescribeNatGatewaysArgs{
			RegionId:     client.Region,
			NatGatewayId: d.Id(),
		}
		gw, _, gwErr := DescribeNatGateways(client.vpcconn, describeArgs)

		if gwErr != nil {
			log.Printf("[ERROR] Describe NatGateways failed.")
			return resource.NonRetryableError(gwErr)
		} else if gw == nil || len(gw) < 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("NatGateway in use - trying again while it is deleted."))
	})
}
