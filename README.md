## Alicloud ([Alibaba Cloud](http://www.aliyun.com)) terraform provider

This is the official repository for the Alicloud terraform provider.
Currently it supports terraform version ≥ v0.8.2.

If you are not planning to contribute to this repo, you can download the compiled binaries from [https://github.com/alibaba/terraform-provider/releases](https://github.com/alibaba/terraform-provider/releases) and move the banaries (bin/terraform-provider-alicloud) into the folder under the Terraform **PATH** such as **/usr/local/terraform**.

Alternatively, the provider can be installed as described in the [developer notes](#developer-notes).  
This way you will be able to sync the repo as active development is going on.

#### Example

Example modules can be found in the [terraform/examples](terraform/examples) directory.

### Developer notes

#### Setting up
* install terraform: https://www.terraform.io/intro/getting-started/install.html
* install golang:    https://golang.org/doc/install
* install glide: https://github.com/Masterminds/glide
* finally:

```
cd $GOPATH
mkdir -p src/github.com/alibaba
cd $GOPATH/src/github.com/alibaba
git clone https://github.com/alibaba/terraform-provider.git

# switch to project
cd $GOPATH/src/github.com/alibaba/terraform-provider

# get all dependencies and install modules
go get ./...
glide up
sudo -E "PATH=$PATH" make all
terraform get

# set the creds
export ALICLOUD_ACCESS_KEY="***"
export ALICLOUD_SECRET_KEY="***"

# you're good to start rocking
# alicloud.tf contains the default example
terraform plan
# terraform apply
# terraform destroy
```

#### Regions
```
cn-shenzhen
ap-southeast-1
cn-qingdao
cn-beijing
cn-shanghai
us-east-1
cn-hongkong
me-east-1
ap-southeast-2
cn-hangzhou
eu-central-1
ap-northeast-1
us-west-1
```

#### Support products
* [ECS](https://www.aliyun.com/product/ecs)
* [Block Storage](https://www.aliyun.com/product/disk)
* [SLB](https://www.aliyun.com/product/slb)
* [VPC](https://www.aliyun.com/product/vpc)
* [NAT Gateway](https://www.aliyun.com/product/nat)
* [RDS](https://www.aliyun.com/product/rds)


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


### How to contribute code
* If you are not sure or have any doubts, feel free to ask and/or submit an issue or PR. We appreciate all contributions and don't want to create artificial obstacles that get in the way.
* Contributions are welcome and will be merged via PRs.

### Contributors
* demonwy(demon.wy@alibaba-inc.com)
* heww(heww0205@gmail.com)
* ShuWei(shuwei.yin@alibaba-inc.com)
* WangYuelucky(wangyuelucky@126.com)
* GuiMin(guimin.hgm@alibaba-inc.com)

### License
* This project is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/alibaba/terraform-provider/blob/master/LICENSE) for the full license text.

### Refrence
* Terraform document: https://www.terraform.io/intro/
