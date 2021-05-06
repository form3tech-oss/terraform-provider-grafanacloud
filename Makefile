VERSION := 0.0.1
INSTALL_DIR := ~/.terraform.d/plugins/github.com/form3tech-oss/grafanacloud/$(VERSION)/linux_amd64
BINARY := terraform-provider-grafanacloud_v$(VERSION)
SHELL := /bin/bash

# Default values used by tests
GRAFANA_CLOUD_MOCK ?= 1
GRAFANA_CLOUD_API_KEY ?= very-secret
GRAFANA_CLOUD_ORGANISATION ?= dummy-org
GRAFANA_CLOUD_STACK ?= dummy-stack

.PHONY: build
build: testacc
	mkdir -p bin
	go build -o bin/$(BINARY) main.go

.PHONY: test
test: vet
	GRAFANA_CLOUD_MOCK=$(GRAFANA_CLOUD_MOCK) \
	go test -count 1 -v ./...

.PHONY: testacc
testacc: vet
	TF_ACC=1 \
	GRAFANA_CLOUD_API_KEY=$(GRAFANA_CLOUD_API_KEY) \
	GRAFANA_CLOUD_ORGANISATION=$(GRAFANA_CLOUD_ORGANISATION) \
	GRAFANA_CLOUD_STACK=$(GRAFANA_CLOUD_STACK) \
	GRAFANA_CLOUD_MOCK=$(GRAFANA_CLOUD_MOCK) \
	go test -count=1 ./... -v $(TESTARGS) -timeout 120m

.PHONY: vet
vet:
	go vet ./...

.PHONY: install
install: test build
	mkdir -p $(INSTALL_DIR)
	cp bin/$(BINARY) $(INSTALL_DIR)/

.PHONY: release
release: bin/goreleaser
	./bin/goreleaser --skip-publish --snapshot

bin/goreleaser:
	wget https://github.com/goreleaser/goreleaser/releases/download/v0.164.0/goreleaser_Linux_x86_64.tar.gz
	tar -C bin/ -xzf goreleaser*tar.gz goreleaser
	rm goreleaser*tar.gz*

.PHONY: docs
docs: bin/tfplugindocs
	./bin/tfplugindocs generate

bin/tfplugindocs:
	@if [[ "$$OSTYPE" == "linux-gnu"* ]]; then \
		wget https://github.com/hashicorp/terraform-plugin-docs/releases/download/v0.4.0/tfplugindocs_0.4.0_linux_amd64.zip; \
	elif [[ "$$OSTYPE" == "darwin"* ]]; then \
		wget https://github.com/hashicorp/terraform-plugin-docs/releases/download/v0.4.0/tfplugindocs_0.4.0_darwin_amd64.zip; \
	fi
	unzip -d bin/ tfplugindocs*zip tfplugindocs
	rm tfplugindocs*zip*

.PHONY: tf-plan
tf-plan: install
	cd examples/full && rm -f .terraform.lock.hcl && terraform init && terraform plan

.PHONY: tf-apply
tf-apply: install
	cd examples/full && rm -f .terraform.lock.hcl && terraform init && terraform apply

.PHONY: tf-destroy
tf-destroy: install
	cd examples/full && rm -f .terraform.lock.hcl && terraform init && terraform destroy
