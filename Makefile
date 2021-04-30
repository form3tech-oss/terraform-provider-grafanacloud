VERSION := 0.0.1
INSTALL_DIR := ~/.terraform.d/plugins/github.com/form3tech-oss/grafanacloud/$(VERSION)/linux_amd64
BINARY := terraform-provider-grafanacloud_v$(VERSION)

# Default values used by acceptance tests (testacc target)
GRAFANA_CLOUD_API_KEY ?= very-secret
GRAFANA_CLOUD_ORGANISATION ?= dummy-org
GRAFANA_CLOUD_STACK ?= dummy-stack
GRAFANA_CLOUD_MOCK ?= 1

.PHONY: build
build:
	mkdir -p bin
	go build -o bin/$(BINARY) main.go

.PHONY: test
test:
	GRAFANA_CLOUD_MOCK=$(GRAFANA_CLOUD_MOCK) \
	go test -count 1 -v ./...

.PHONY: testacc
testacc:
	TF_ACC=1 \
	GRAFANA_CLOUD_API_KEY=$(GRAFANA_CLOUD_API_KEY) \
	GRAFANA_CLOUD_ORGANISATION=$(GRAFANA_CLOUD_ORGANISATION) \
	GRAFANA_CLOUD_STACK=$(GRAFANA_CLOUD_STACK) \
	GRAFANA_CLOUD_MOCK=$(GRAFANA_CLOUD_MOCK) \
	go test -count=1 ./... -v $(TESTARGS) -timeout 120m

.PHONY: install
install: test build
	mkdir -p $(INSTALL_DIR)
	cp bin/$(BINARY) $(INSTALL_DIR)/

.PHONY: docs
docs:
	tfplugindocs generate

.PHONY: tf-plan
tf-plan: install
	cd examples/full && rm -f .terraform.lock.hcl && terraform init && terraform plan

.PHONY: tf-apply
tf-apply: install
	cd examples/full && rm -f .terraform.lock.hcl && terraform init && terraform apply

.PHONY: tf-destroy
tf-destroy: install
	cd examples/full && rm -f .terraform.lock.hcl && terraform init && terraform destroy
