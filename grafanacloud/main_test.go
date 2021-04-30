package grafanacloud_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/naag/terraform-provider-grafanacloud/grafanacloud"
	"github.com/naag/terraform-provider-grafanacloud/internal/mock"
)

var (
	testAccProviders map[string]*schema.Provider
	testAccProvider  *schema.Provider
	grafanaCloudMock *mock.GrafanaCloud
)

const (
	EnvMock = "GRAFANA_CLOUD_MOCK"
)

func TestMain(m *testing.M) {
	startMock()
	if grafanaCloudMock != nil {
		defer grafanaCloudMock.Close()
	}

	testAccProvider = grafanacloud.NewProvider("0.0.1")()
	testAccProviders = map[string]*schema.Provider{
		"grafanacloud": testAccProvider,
	}

	os.Exit(m.Run())
}

func startMock() {
	if os.Getenv(EnvMock) == "1" {
		org := os.Getenv(grafanacloud.EnvOrganisation)
		grafanaCloudMock = mock.NewGrafanaCloud(org).
			Start()

		os.Setenv(grafanacloud.EnvURL, grafanaCloudMock.URL())
	}
}
