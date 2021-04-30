# Terraform Provider Grafana Cloud

A Terraform provider for managing Grafana Cloud resources.

## Use Cases

A few possible use-cases for `terraform-provider-grafanacloud` are:

- Managing Grafana Cloud stacks
- Managing API keys for both Grafana Cloud and Grafana instances inside stacks
- Rolling API keys by tainting TF resources
- Collecting information about configured stacks, such as Prometheus / Alertmanager endpoints or user IDs
- Reading Grafana data sources

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.15

## Installing the provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using `make build`
1. Install the provider locally using `make install`
1. Add the following code snippet to your Terraform code:
```tf
terraform {
  required_providers {
    grafanacloud = {
      source  = "github.com/form3tech-oss/grafanacloud"
      version = "0.0.1"
    }
  }
}
```

## Using the provider

### Configuration

The following provider block variables are available for configuration:

| Name | Description | Default |
| ---- | ----------- | ------- |
| `url` | The URL to Grafana Cloud API | `https://grafana.com/api` |
| `api_key` | The API key used to authenticate with Grafana Cloud. If you want to manage API keys using this provider, this needs to have the `Admin` role | - |
| `organisation` | Slug name of the organisation to manage | - |

For more detailed docs, please refer to the [generated docs](/docs/index.md).

## Developing the provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make build install`. This will build the provider and put the provider binary in the Terraform plugin directory.

To generate or update documentation, run `go generate`.

## Testing the provider

Both unit tests (`make test`) and acceptance tests (`make testacc`) are provided.

By default, the acceptance tests run against a mock implementation of the Grafana Cloud/Grafana APIs provided by the __internal/mock__ package. If you'd like to instead run acceptance tests against the real API endpoints, please run the following command with proper values for the placeholder variables:

```sh
GRAFANA_CLOUD_MOCK=0 \
GRAFANA_CLOUD_ORGANISATION=my-org-slug \
GRAFANA_CLOUD_STACK=my-stack-slug \
GRAFANA_CLOUD_API_KEY=a-secret-admin-api-key \
make testacc
```

Due to API rate limits, such a test run may not actually pass or might even impact your production API usage, so please be careful about that.

You can obtain debugging output from the `go-resty` HTTP client by setting `HTTP_DEBUG=1` when running tests.
