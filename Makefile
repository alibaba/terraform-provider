GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)
VETARGS?=-all
TEST?=$$(go list ./...)


all: build

build: mac windows linux

dev: clean fmt mac copy

copy:
	tar -xvf bin/terraform-provider-alicloud_darwin-amd64.tgz && mv bin/terraform-provider-alicloud /usr/local/Cellar/terraform/0.11.5/bin

test: vet fmtcheck errcheck
	TF_ACC=1 go test -v ./alicloud -run=TestAccAlicloud -timeout=180m -parallel=4

vet:
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

fmt:
	gofmt -w $(GOFMT_FILES)
	goimports -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

clean:
	rm -rf bin/*

mac:
	GOOS=darwin GOARCH=amd64 go build -o bin/terraform-provider-alicloud
	mv bin/terraform-provider-alicloud /usr/local/bin

windows:
	GOOS=windows GOARCH=amd64 go build -o bin/terraform-provider-alicloud.exe
	tar czvf bin/terraform-provider-alicloud_windows-amd64.tgz bin/terraform-provider-alicloud.exe
	

linux:
	GOOS=linux GOARCH=amd64 go build -o bin/terraform-provider-alicloud
	tar czvf bin/terraform-provider-alicloud_linux-amd64.tgz bin/terraform-provider-alicloud
	rm -rf bin/terraform-provider-alicloud

