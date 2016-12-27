all: build copy

build:
	go build -o terraform-provider-alicloud

copy:
	cp terraform-provider-alicloud $(shell dirname `which terraform`)

test:
	TF_ACC=1 go test -v ./alicloud -timeout 120m

vet:
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi