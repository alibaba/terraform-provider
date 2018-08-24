package alicloud

import (
	"github.com/denverdino/aliyungo/sts"
)

func (client *AliyunClient) GetCallerIdentity() (*sts.GetCallerIdentityResponse, error) {
	invoker := NewInvoker()
	var identityResponse *sts.GetCallerIdentityResponse

	err := invoker.Run(func() error {
		identity, err := client.stsconn.GetCallerIdentity()
		if err != nil {
			return err
		}
		if identity == nil {
			return GetNotFoundErrorFromString("Caller identity not found.")
		}
		identityResponse = identity
		return nil
	})
	return identityResponse, err
}
