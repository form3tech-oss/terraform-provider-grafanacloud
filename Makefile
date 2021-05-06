include Makefile.deps

VERSION := 0.0.1
INSTALL_DIR := ~/.terraform.d/plugins/github.com/form3tech-oss/grafanacloud/$(VERSION)/linux_amd64
BINARY := terraform-provider-grafanacloud_v$(VERSION)
SHELL := /bin/bash
PATH := $(PATH):$(PWD)/bin

# Default values used by tests
GRAFANA_CLOUD_MOCK ?= 1
GRAFANA_CLOUD_API_KEY ?= very-secret
GRAFANA_CLOUD_ORGANISATION ?= dummy-org
GRAFANA_CLOUD_STACK ?= dummy-stack

build: lint testacc
	mkdir -p bin
	go build -o bin/$(BINARY) main.go

test:
	GRAFANA_CLOUD_MOCK=$(GRAFANA_CLOUD_MOCK) \
	go test -count 1 -v ./...

testacc:
	TF_ACC=1 \
	GRAFANA_CLOUD_API_KEY=$(GRAFANA_CLOUD_API_KEY) \
	GRAFANA_CLOUD_ORGANISATION=$(GRAFANA_CLOUD_ORGANISATION) \
	GRAFANA_CLOUD_STACK=$(GRAFANA_CLOUD_STACK) \
	GRAFANA_CLOUD_MOCK=$(GRAFANA_CLOUD_MOCK) \
	go test -count=1 ./... -v $(TESTARGS) -timeout 120m

lint: vet tflint tffmtcheck

vet:
	go vet ./...

tflint:
	find ./examples/ -type d -exec tflint \{\} \;

tffmtcheck:
	terraform fmt -check -recursive ./examples/

fmt:
	go fmt ./...
	terraform fmt -recursive ./examples/

install: test build
	mkdir -p $(INSTALL_DIR)
	cp bin/$(BINARY) $(INSTALL_DIR)/

release:
	goreleaser

docs:
	tfplugindocs generate

tf-plan: install
	cd examples/full && rm -f .terraform.lock.hcl && terraform init && terraform plan

tf-apply: install
	cd examples/full && rm -f .terraform.lock.hcl && terraform init && terraform apply

tf-destroy: install
	cd examples/full && rm -f .terraform.lock.hcl && terraform init && terraform destroy

.PHONY: build test testacc lint vet tffmtcheck fmt install release docs tf-plan tf-apply tf-destroy
