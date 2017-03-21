package alicloud

import (
	"fmt"

	"github.com/denverdino/aliyungo/ecs"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"log"
	"strings"
	"time"
	"github.com/denverdino/aliyungo/common"
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

			"bandwidth_package_ids": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"snat_table_ids" : &schema.Schema{
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
						"public_ip_addresses": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Required: true,
				MaxItems: 4,
			},

			"public_ip_addresses": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAliyunNatGatewayCreate(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AliyunClient).vpcconn

	args := &ecs.CreateNatGatewayArgs{
		RegionId: getRegion(d, meta),
		VpcId:    d.Get("vpc_id").(string),
		Spec:     d.Get("spec").(string),
	}

	bandwidthPackages := d.Get("bandwidth_packages").([]interface{})

	bandwidthPackageTypes := []ecs.BandwidthPackageType{}

	for _, e := range bandwidthPackages {
		pack := e.(map[string]interface{})
		bandwidthPackage := ecs.BandwidthPackageType{
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

	args.Name = name

	if v, ok := d.GetOk("description"); ok {
		args.Description = v.(string)
	}
	resp, err := conn.CreateNatGateway(args)
	if err != nil {
		return fmt.Errorf("CreateNatGateway got error: %#v", err)
	}

	d.SetId(resp.NatGatewayId)

	return resourceAliyunNatGatewayRead(d, meta)
}

func resourceAliyunNatGatewayRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	//conn := client.vpcconn

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
	d.Set("bandwidth_package_ids", strings.Join(natGateway.BandwidthPackageIds.BandwidthPackageId, ","))
	d.Set("snat_table_ids", strings.Join(natGateway.SnatTableIds.SnatTableId, ","))
	d.Set("description", natGateway.Description)
	d.Set("vpc_id", natGateway.VpcId)
	bindWidthPackages, err := flattenBandWidthPackages(natGateway.BandwidthPackageIds.BandwidthPackageId, meta, d)
	if(err != nil){
		log.Printf("[ERROR] bindWidthPackages flattenBandWidthPackages failed. natgateway id is %#v", d.Id())
		fmt.Println("[ERROR] bindWidthPackages flattenBandWidthPackages failed. natgateway id is %#v", d.Id())
	}else {
		fmt.Println("bindWidthPackages %#v", bindWidthPackages)
		d.Set("bandwidth_packages", bindWidthPackages)
	}
	ips, _ := flattenBandWidthPackagesIp(natGateway.BandwidthPackageIds.BandwidthPackageId, meta, d)
	d.Set("public_ip_addresses", ips)

	return nil
}

func resourceAliyunNatGatewayUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	conn := client.vpcconn

	natGateway, err := client.DescribeNatGateway(d.Id())
	if err != nil {
		return err
	}

	d.Partial(true)
	attributeUpdate := false
	args := &ecs.ModifyNatGatewayAttributeArgs{
		RegionId:     natGateway.RegionId,
		NatGatewayId: natGateway.NatGatewayId,
	}

	if d.HasChange("name") {
		d.SetPartial("name")
		var name string
		if v, ok := d.GetOk("name"); ok {
			name = v.(string)
		} else {
			return fmt.Errorf("cann't change name to empty string")
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
		if err := conn.ModifyNatGatewayAttribute(args); err != nil {
			return err
		}
	}

	if d.HasChange("spec") {
		d.SetPartial("spec")
		var spec ecs.NatGatewaySpec
		if v, ok := d.GetOk("spec"); ok {
			spec = ecs.NatGatewaySpec(v.(string))
		} else {
			// set default to small spec
			spec = ecs.NatGatewaySmallSpec
		}

		args := &ecs.ModifyNatGatewaySpecArgs{
			RegionId:     natGateway.RegionId,
			NatGatewayId: natGateway.NatGatewayId,
			Spec:         spec,
		}

		err := conn.ModifyNatGatewaySpec(args)
		if err != nil {
			return fmt.Errorf("%#v %#v", err, *args)
		}

	}
	d.Partial(false)

	return resourceAliyunNatGatewayRead(d, meta)
}

func resourceAliyunNatGatewayDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*AliyunClient)
	conn := client.vpcconn

	return resource.Retry(5*time.Minute, func() *resource.RetryError {

		packages, err := conn.DescribeBandwidthPackages(&ecs.DescribeBandwidthPackagesArgs{
			RegionId:     getRegion(d, meta),
			NatGatewayId: d.Id(),
		})
		if err != nil {
			log.Printf("[ERROR] Describe bandwidth package is failed, natGateway Id: %s", d.Id())
			return resource.NonRetryableError(err)
		}

		retry := false
		for _, pack := range packages {
			err = conn.DeleteBandwidthPackage(&ecs.DeleteBandwidthPackageArgs{
				RegionId:           getRegion(d, meta),
				BandwidthPackageId: pack.BandwidthPackageId,
			})

			if err != nil {
				er, _ := err.(*common.Error)
				if er.ErrorResponse.Code == NatGatewayInvalidRegionId {
					log.Printf("[ERROR] Delete bandwidth package is failed, bandwidthPackageId: %#v", pack.BandwidthPackageId)
					return resource.NonRetryableError(err)
				}
				retry = true
			}
		}

		if retry {
			return resource.RetryableError(fmt.Errorf("Bandwidth package in use - trying again while it is deleted."))
		}

		args := &ecs.DeleteNatGatewayArgs{
			RegionId:     getRegion(d, meta),
			NatGatewayId: d.Id(),
		}

		err = conn.DeleteNatGateway(args)
		if err != nil {
			er, _ := err.(*common.Error)
			if er.ErrorResponse.Code == DependencyViolationBandwidthPackages {
				return resource.RetryableError(fmt.Errorf("NatGateway in use - trying again while it is deleted."))
			}
		}

		describeArgs := &ecs.DescribeNatGatewaysArgs{
			RegionId:     getRegion(d, meta),
			NatGatewayId: d.Id(),
		}
		gw, _, gwErr := conn.DescribeNatGateways(describeArgs)

		if gwErr != nil {
			log.Printf("[ERROR] Describe NatGateways failed.")
			return resource.NonRetryableError(gwErr)
		} else if gw == nil || len(gw) < 1 {
			return nil
		}

		return resource.RetryableError(fmt.Errorf("NatGateway in use - trying again while it is deleted."))
	})
}

func flattenBandWidthPackagesIp(bandWidthPackageIds []string,  meta interface{}, d *schema.ResourceData) (string, error)  {
	var result string

	result = "ipaddresstest"

	for _, packageId := range bandWidthPackageIds {
		packages, err := getPackages(packageId, meta, d)
		if(err != nil){
			log.Printf("[ERROR] NatGateways getPackages failed. packageId is %#v", packageId)
			return result, err
		}
		for _, pack := range packages {
			ipAddress := flattenPackPublicIp(pack.PublicIpAddresses.PublicIpAddresse)

			result = ipAddress
		}

	}
	return result, nil
}


func flattenBandWidthPackages(bandWidthPackageIds []string,  meta interface{}, d *schema.ResourceData) ([]map[string]interface{}, error)  {

	//result := make([]map[string]interface{}, 0, len(bandWidthPackageIds))
	//
	//for _, packageId := range bandWidthPackageIds {
	//	hasPublicIpErr := waitGetPublicIp(packageId, meta, d)
	//	if(hasPublicIpErr == nil){
	//		packages, err := getPackages(packageId, meta, d)
	//		if(err != nil){
	//			log.Printf("[ERROR] NatGateways getPackages failed. packageId is %#v", packageId)
	//			return result, err
	//		}
	//		for _, pack := range packages {
	//			ipAddress := flattenPackPublicIp(pack.PublicIpAddresses.PublicIpAddresse)
	//			println("ipAddress:%#v", ipAddress)
	//			l := map[string]interface{}{
	//				"ip_count":            pack.IpCount,
	//				"bandwidth": 8,
	//				"zone":     pack.ZoneId,
	//				"public_ip_addresses": "ipaddresstest",
	//			}
	//			result = append(result, l)
	//		}
	//	}
	//
	//}
	//return result, nil

	result := make([]map[string]interface{}, 0, len(bandWidthPackageIds))

	for _, packageId := range bandWidthPackageIds {
			packages, err := getPackages(packageId, meta, d)
			if(err != nil){
				log.Printf("[ERROR] NatGateways getPackages failed. packageId is %#v", packageId)
				return result, err
			}
			for _, pack := range packages {
				ipAddress := flattenPackPublicIp(pack.PublicIpAddresses.PublicIpAddresse)
				println("ipAddress:%#v", ipAddress)
				l := map[string]interface{}{
					"ip_count":            pack.IpCount,
					"bandwidth": 8,
					"zone":     pack.ZoneId,
					"public_ip_addresses": "ipaddresstest",
				}
				result = append(result, l)
			}

	}
	return result, nil

}

func waitGetPublicIp(packageId string, meta interface{}, d *schema.ResourceData)  error{
	return resource.Retry(1*time.Minute, func() *resource.RetryError {
		packages, err := getPackages(packageId, meta, d)
		if(err != nil){
			log.Printf("[ERROR] NatGateways getPackages failed. packageId is %#v", packageId)
			return resource.NonRetryableError(err)
		}
		for _, pack := range packages {
			println("PublicIpAddress %#v",pack.PublicIpAddresses.PublicIpAddresse)
			ipAddrLen := len(pack.PublicIpAddresses.PublicIpAddresse)
			println("ipAddrLen: %#v", ipAddrLen)
			if ipAddrLen > 0{
				return nil
			}else {
				return resource.RetryableError(fmt.Errorf(NoPublicIpAddressInPackage))
			}
		}
		return resource.RetryableError(fmt.Errorf("DescribeBandwidthPackages no packages."))
	})
}

func getPackages(packageId string, meta interface{}, d *schema.ResourceData) ([] ecs.DescribeBandwidthPackageType, error) {
	client := meta.(*AliyunClient)
	conn := client.vpcconn
	packages, err := conn.DescribeBandwidthPackages(&ecs.DescribeBandwidthPackagesArgs{
		RegionId:     getRegion(d, meta),
		BandwidthPackageId: packageId,
	})

	if err != nil {
		log.Printf("[ERROR] Describe bandwidth package is failed, BandwidthPackageId Id: %s", packageId)
		return nil,err
	}

	return packages, nil

}

func flattenPackPublicIp(publicIpAddressList []ecs.PublicIpAddresseType) string {
	var result []string

	for _, publicIpAddresses := range publicIpAddressList {
		ipAddress := publicIpAddresses.IpAddress
		result = append(result, ipAddress)
	}

	return strings.Join(result, ",")
}