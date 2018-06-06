package slb

import "testing"

func testListeners(t *testing.T, client *Client, loadBalancerId string) {

	port := 1234

	creationArgs := CreateLoadBalancerTCPListenerArgs{
		LoadBalancerId:    loadBalancerId,
		ListenerPort:      port,
		BackendServerPort: 1234,
		Bandwidth:         -1,
	}

	err := client.CreateLoadBalancerTCPListener(&creationArgs)
	if err != nil {
		t.Errorf("Failed to CreateLoadBalancerTCPListener: %v", err)
	}

	response, err := client.DescribeLoadBalancerTCPListenerAttribute(loadBalancerId, port)
	if err != nil {
		t.Errorf("Failed to DescribeLoadBalancerTCPListenerAttribute: %v", err)
	}
	t.Logf("Listener: %++v", *response)

	err = client.StartLoadBalancerListener(loadBalancerId, port)
	if err != nil {
		t.Errorf("Failed to StartLoadBalancerListener: %v", err)
	}

	status, err := client.WaitForListener(loadBalancerId, port, TCP)
	if err != nil {
		t.Errorf("Failed to WaitForListener: %v", err)
	}
	t.Logf("Listener status: %s", status)

	response, err = client.DescribeLoadBalancerTCPListenerAttribute(loadBalancerId, port)
	if err != nil {
		t.Errorf("Failed to DescribeLoadBalancerTCPListenerAttribute: %v", err)
	}
	t.Logf("Listener: %++v", *response)
}
