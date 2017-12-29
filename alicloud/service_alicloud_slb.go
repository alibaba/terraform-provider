package alicloud

import (
	"fmt"
	"github.com/denverdino/aliyungo/slb"
)

func (client *AliyunClient) DescribeLoadBalancerAttribute(slbId string) (*slb.LoadBalancerType, error) {

	loadBalancers, err := client.slbconn.DescribeLoadBalancers(&slb.DescribeLoadBalancersArgs{
		RegionId:       client.Region,
		LoadBalancerId: slbId,
	})
	if err != nil {
		if IsExceptedError(err, LoadBalancerNotFound) {
			return nil, fmt.Errorf("Special SLB Id not found: %#v", err)
		}

		return nil, fmt.Errorf("DescribeLoadBalancers got an error: %#v", err)
	}
	if len(loadBalancers) < 1 {
		return nil, fmt.Errorf("Special SLB Id %s is not found in %#v.", slbId, client.Region)
	}
	return &loadBalancers[0], nil
}
