package alicloud

import "github.com/denverdino/aliyungo/common"


func (client *AliyunClient) DescribeNatGateway(natGatewayId string) (*NatGatewaySetType, error) {

	args := &DescribeNetGatewaysArgs{
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