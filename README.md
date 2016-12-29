## Aliyun (Alibaba Cloud) terraform provider


This is the official repository for the Aliyun terraform provider. Currently it supports terraform version ≥ v0.8.2 .

If just use the terraform-provide without develop the project, you can download the following binary package and move it into special forlder where download the Terraform PATH, such as "/usr/local/terraform".

* Mac OS X  [64-bit](http://tf-mac.oss-cn-shanghai.aliyuncs.com/terraform-provider-alicloud.zip )
* Linux  [64-bit](http://tf-linux.oss-cn-shanghai.aliyuncs.com/terraform-provider-alicloud.zip )
* Windows  [64-bit](http://tf-windows.oss-cn-shanghai.aliyuncs.com/terraform-provider-alicloud.exe)

Currently it is under active development and must be installed as described in the [developer notes](#developer-notes).

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
export TF_VAR_ali_access_key="***"
export TF_VAR_ali_secret_key="***"

# you're good to start rocking
# alicloud.tf contains the default example
terraform plan
# terraform apply
# terraform destroy
```

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
* If you're unsure or afraid of anything, just ask or submit the issue or pull request anyways. We appreciate any sort of contributions, and don't want a wall of rules to get in the way of that.
* Contributions are welcome and will be merged via PRs.

### Contributors
* heww(heww0205@gmail.com)
* ShuWei(shuwei.yin@alibaba-inc.com)
* WangYuelucky(wangyuelucky@126.com)
* GuiMin(guimin.hgm@alibaba-inc.com)

### License
* This project is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/alibaba/terraform-provider/blob/master/LICENSE) for the full license text.

### Refrence
* Terraform document: https://www.terraform.io/intro/
