package aliyunclient

import (
	"github.com/denverdino/aliyungo/common"
	"fmt"
)

// Config of aliyun
type Config struct {
	AccessKey       string
	SecretKey       string
	Region          common.Region
	RegionId        string
	SecurityToken   string
	OtsInstanceName string
	LogEndpoint     string
	AccountId       string
	FcEndpoint      string
}

func (c *Config) loadAndValidate() error {
	err := c.validateRegion()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) validateRegion() error {

	for _, valid := range common.ValidRegions {
		if c.Region == valid {
			return nil
		}
	}

	return fmt.Errorf("Not a valid region: %s", c.Region)
}