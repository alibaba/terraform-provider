all: build copy

build:
	go build -o terraform-provider-alicloud

copy:
	cp terraform-provider-alicloud $(shell dirname `which terraform`)

test:
	TF_ACC=1 go test -v ./alicloud
