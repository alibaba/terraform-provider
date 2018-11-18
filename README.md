# [NOTES]: This repo has been deprecated and please refer to the latest repo: https://github.com/terraform-providers/terraform-provider-alicloud

## Alicloud ([Alibaba Cloud](http://www.aliyun.com)) Terraform provider

This is the official repository of the Alicloud Terraform provider.
It currently supports Terraform version ≥ v0.8.2.

If you are not planning to contribute to this repo, you can download the compiled binaries from [https://github.com/alibaba/terraform-provider/releases](https://github.com/alibaba/terraform-provider/releases) and move the binaries (bin/terraform-provider-alicloud) into the folder under the Terraform **PATH** such as **/usr/local/terraform**.

Alternatively, the provider can be installed as described in the [developer notes](#developer-notes).  
This way you will be able to sync the repo as active development is going on.

-> **Note:** When you use Terraform on a `Windows` computer, please install [golang](https://golang.org/dl/) first,
otherwise you might get [this issue](https://github.com/alibaba/terraform-provider/issues/469)
(the bug appeared in the version 1.8.1 and was fixed in the version 1.11.0).

#### Example

Example modules can be found in the [terraform/examples](examples) directory.

### Developer notes

#### Setting up
* Install Terraform: https://www.terraform.io/intro/getting-started/install.html
* Install golang:    https://golang.org/doc/install
* Install goimports: https://godoc.org/golang.org/x/tools/cmd/goimports
    ```
    go get golang.org/x/tools/cmd/goimports
    ```
* Finally:

```
mkdir -p $GOPATH/src/github.com/alibaba
cd $GOPATH/src/github.com/alibaba
git clone https://github.com/alibaba/terraform-provider.git

# switch to project
cd $GOPATH/src/github.com/alibaba/terraform-provider

# build provider

## Mac
make dev
## Linux
make devlinux
## Windows
make devwin

# install modules
terraform get

# set the creds
export ALICLOUD_ACCESS_KEY="***"
export ALICLOUD_SECRET_KEY="***"
export ALICLOUD_REGION="***"
export ALICLOUD_ACCOUNT_ID="***"

# you're good to start rocking
# alicloud.tf contains the default example
terraform plan
# terraform apply
# terraform destroy
```

#### Regions
```
cn-qingdao
cn-beijing
cn-zhangjiakou
cn-huhehaote
cn-hangzhou
cn-shanghai
cn-shenzhen
cn-hongkong

ap-northeast-1
ap-southeast-1
ap-southeast-2
ap-southeast-3
ap-southeast-5
ap-south-1

us-east-1
us-west-1

me-east-1

eu-central-1
```
For more information about the regions and availability zones, please use the data sources `alicloud_regions` and `alicloud_zones`.

#### Supported products
* [ECS](https://www.aliyun.com/product/ecs)
* [Block Storage](https://www.aliyun.com/product/disk)
* [SLB](https://www.aliyun.com/product/slb)
* [VPC](https://www.aliyun.com/product/vpc)
* [NAT Gateway](https://www.aliyun.com/product/nat)
* [RDS](https://www.aliyun.com/product/rds)
* [ESS](https://www.aliyun.com/product/ess)
* [OSS](https://www.aliyun.com/product/oss)
* [RAM](https://www.aliyun.com/product/ram)
* [CDN](https://www.aliyun.com/product/cdn)
* [DNS](https://wanwang.aliyun.com/domain/dns)
* [Container Service](https://www.aliyun.com/product/containerservice)
* [Log Service](https://www.aliyun.com/product/sls)
* [Function Compute](https://www.aliyun.com/product/fc)
* [TableStore](https://www.aliyun.com/product/ots)

#### Documents
The most recent documentation is available here:
* [Terraform Docs](https://www.terraform.io/docs/providers/alicloud/index.html)
* [Github](https://github.com/alibaba/terraform-provider-docs)

#### Acceptance Testing
Before making a release, the resources and data sources are tested automatically with acceptance tests (the tests are located in the alicloud/*_test.go files).
You can run them by entering the following instructions in a terminal:
```
cd $GOPATH/src/github.com/alibaba/terraform-provider
export ALICLOUD_ACCESS_KEY=xxx
export ALICLOUD_SECRET_KEY=xxx
export ALICLOUD_REGION=xxx
export ALICLOUD_ACCOUNT_ID=xxx
export outfile=gotest.out
TF_ACC=1 TF_LOG=INFO go test ./alicloud -v -run=TestAccAlicloud -timeout=1440m | tee $outfile
go2xunit -input $outfile -output $GOPATH/tests.xml
```

-> **Note:** The last line is optional, it allows to convert test results into a XML format compatible with xUnit.

Because some features are not available in all regions, the following environment variables can be set in order to
skip tests that use these features:
* ALICLOUD_SKIP_TESTS_FOR_SLB_SPECIFICATION=true    - Server Load Balancer with guaranteed performance specifications (old implementation has only shared performance)
* ALICLOUD_SKIP_TESTS_FOR_SLB_PAY_BY_BANDWIDTH=true - Server Load Balancer with a "pay by bandwidth" billing method (mostly available in China)
* ALICLOUD_SKIP_TESTS_FOR_FUNCTION_COMPUTE=true     - Function Compute
* ALICLOUD_SKIP_TESTS_FOR_PVTZ_ZONE=true            - Private Zone
* ALICLOUD_SKIP_TESTS_FOR_RDS_MULTIAZ=true          - Apsara RDS with multiple availability zones
* ALICLOUD_SKIP_TESTS_FOR_CLASSIC_NETWORK=true      - Classic network configuration

#### Common problems

1.
```
Error configuring: 1 error(s) occurred:
* Incompatible API version with plugin. Plugin version: 2, Ours: 1

# fix by manually setting the branch in the sources
cd src/github.com/hashicorp/terraform/
git checkout v<YOUR_TF_VERSION_HERE>
cd -

# rebuild
sudo -E "PATH=$PATH" make all
```


### How to contribute
* If you are not sure or have any doubt, feel free to ask and/or submit an issue or PR.
  We appreciate all contributions and try to make the process as smooth as possible.
* Contributions are welcome and will be merged via PRs.

### Contributors
* demonwy
* heww(heww0205@gmail.com)
* ShuWei
* WangYuelucky(wangyuelucky@126.com)
* GuiMin(heguimin36@163.com)

### License
* This project is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/alibaba/terraform-provider/blob/master/LICENSE) for the full license text.

### Reference
* Terraform Document: https://www.terraform.io/intro/
* Terraform Registry Of Alicloud: https://registry.terraform.io/modules/alibaba
* Terraform Alicloud Examples: [Official Examples](https://github.com/terraform-providers/terraform-provider-alicloud/tree/master/examples) , [Other Examples](https://github.com/mosuke5/terraform_examples_for_alibabacloud)
