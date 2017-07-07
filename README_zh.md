
# 利用Terraform创建/配置Aliyun资源

* Terraform是HashiCorp公司出品的，著名Vagrant、Packer工具也出自于该公司。Terraform是“基础设施即代码”的开源工具，通过模板定义“基础设施”，目前已支持AWS、Azure、OpenStack等主流云平台，这个项目是阿里云官方出品，利用Terraform创建阿里云的资源。

### 安装 Terraform
* Alicloud terraform provider 当前支持的 Terraform 的版本 ≥ v0.8.2
* 安装Terraform可参考：https://www.terraform.io/intro/getting-started/install.html
		
* 注意：设置环境变量时，需要指向terraform所在的父目录，如：terraform的路径是"~/work/terraform_0.8.2"，则指定环境时设定为export PATH=$PATH:~/work/terraform_0.8.2

* 如果想直接使用此项目中提供的模板，可以直接下载如下编译好的二进制文件，并拷贝到与Terraform同级的目录下，如~/work/terraform_0.8.2。

  * Mac OS X  [64-bit](http://tf-mac.oss-cn-shanghai.aliyuncs.com/terraform-provider-alicloud.zip )
  * Linux  [64-bit](http://tf-linux.oss-cn-shanghai.aliyuncs.com/terraform-provider-alicloud.zip )
  * Windows  [64-bit](http://tf-windows.oss-cn-shanghai.aliyuncs.com/terraform-provider-alicloud.exe)


### 安装GoLang运行环境
* 如果只想引用此项目的tf模板，不需要安装GoLang的运行环境，此运行环境是为了编译修改源代码所用。
* 定义Go的工作目录，如： ~/work/go
* 下载 GoLong SDK: http://www.golangtc.com/download 
* 这里有mac下安装 Golang SDK的参考： http://www.tuicool.com/articles/Fv6zUfE
* 设置环境变量，如：
		
		export PATH=$PATH:/usr/local/go/bin
		export GOPATH="~/work/go"

* 创建Go的基本目录结构，如 go get github.com/denverdino/aliyungo

### 如何下载编译使用本项目
1. 在GoPath中创建 "alibaba" 目录，如 "work/go/src/github.com/alibaba"，然后克隆此项目。
2. 在 "github.com/alibaba/terraform-provider" 目录下执行：
		
		 "go get ./..." //this will download depend package
		 "make all"
		 "terraform get"
		 "terraform plan" //input the tips parameters, such as access_key, secret_key , or zone, ecs instance name etc.
		 "terraform apply"
		 
3. 或者使用预设置的参数：:

		export ALICLOUD_ACCESS_KEY=*** 
		export ALICLOUD_SECRET_KEY=***
		terraform get
		terraform plan
		terraform apply
		terraform destroy
		...
		
提示：如果想直接使用此项目中提供的模板，而不需要编译或修改源码，可以直接编写Terraform的模板并运行"terraform plan"等命令，模板参考下面的“Terraform模板文件夹介绍”

### Example

* [这里](terraform/examples)提供了一些场景的Example，包括ECS、SLB、VPC、安全组等模板。


### 如何贡献代码
* 如果您在使用的过程中遇到任何问题，可直接通过 submit issue 或者 pull request进行询问，我们将尽快予以解决。
* 我们欢迎有更多的开发者贡献代码，包括GoLang的源码及Terraform的模板。Fork此项目，然后提交pull request，我们review代码后，即可合并进主干。

### 贡献者
* heww(heww0205@gmail.com)
* ShuWei
* WangYuelucky(wangyuelucky@126.com)
* GuiMin(heguimin36@163.com)

### License
* This project is licensed under the Apache License, Version 2.0. See [LICENSE](https://github.com/denverdino/aliyungo/blob/master/LICENSE.txt) for the full license text.

#### 参考
* Terraform官方文档：https://www.terraform.io/intro/