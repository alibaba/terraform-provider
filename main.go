package main

import (
	"github.com/alibaba/terraform-alicloud/alicloud"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: alicloud.Provider,
	})
}
