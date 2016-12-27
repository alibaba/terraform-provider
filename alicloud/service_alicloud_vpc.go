package alicloud

import (
	"github.com/denverdino/aliyungo/common"
	"github.com/denverdino/aliyungo/ecs"
)

func (client *AliyunClient) DescribeEipAddress(allocationId string) (*ecs.EipAddressSetType, error) {

	args := ecs.DescribeEipAddressesArgs{
		RegionId:     client.Region,
		AllocationId: allocationId,
	}

	eips, _, err := client.ecsconn.DescribeEipAddresses(&args)
	if err != nil {
		return nil, err
	}
	if len(eips) == 0 {
		return nil, common.GetClientErrorFromString("Not found")
	}

	return &eips[0], nil
}

func (client *AliyunClient) DescribeNatGateway(natGatewayId string) (*NatGatewaySetType, error) {

	args := &DescribeNatGatewaysArgs{
		RegionId:     client.Region,
		NatGatewayId: natGatewayId,
	}

	natGateways, _, err := DescribeNatGateways(client.ecsconn, args)
	if err != nil {
		return nil, err
	}

	if len(natGateways) == 0 {
		return nil, common.GetClientErrorFromString("Not found")
	}

	return &natGateways[0], nil
}

func (client *AliyunClient) DescribeVpc(vpcId string) (*ecs.VpcSetType, error) {
	args := ecs.DescribeVpcsArgs{
		RegionId: client.Region,
		VpcId:    vpcId,
	}

	vpcs, _, err := client.ecsconn.DescribeVpcs(&args)
	if err != nil {
		if notFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	if len(vpcs) == 0 {
		return nil, nil
	}

	return &vpcs[0], nil
}

// describe vswitch by param filters
func (client *AliyunClient) QueryVswitches(args *ecs.DescribeVSwitchesArgs) (vswitches []ecs.VSwitchSetType, err error) {
	vsws, _, err := client.ecsconn.DescribeVSwitches(args)
	if err != nil {
		if notFoundError(err) {
			return nil, nil
		}
		return nil, err
	}

	return vsws, nil
}

func (client *AliyunClient) QueryVswitchById(vpcId, vswitchId string) (vsw *ecs.VSwitchSetType, err error) {
	args := &ecs.DescribeVSwitchesArgs{
		VpcId:     vpcId,
		VSwitchId: vswitchId,
	}
	vsws, err := client.QueryVswitches(args)
	if err != nil {
		return nil, err
	}

	if len(vsws) == 0 {
		return nil, nil
	}

	return &vsws[0], nil
}
