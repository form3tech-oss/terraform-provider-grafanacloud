package grafanacloud_test

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/naag/terraform-provider-grafanacloud/grafanacloud"
	"github.com/stretchr/testify/require"
)

func TestProvider(t *testing.T) {
	if err := grafanacloud.NewProvider("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProviderConfigure(t *testing.T) {
	resourceSchema := map[string]*schema.Schema{
		"url": {
			Type: schema.TypeString,
		},
		"api_key": {
			Type: schema.TypeString,
		},
		"organisation": {
			Type: schema.TypeString,
		},
	}

	resourceDataMap := map[string]interface{}{
		"url":          os.Getenv(grafanacloud.EnvURL),
		"api_key":      os.Getenv(grafanacloud.EnvAPIKey),
		"organisation": os.Getenv(grafanacloud.EnvOrganisation),
	}
	resourceLocalData := schema.TestResourceDataRaw(t, resourceSchema, resourceDataMap)

	configureFunc := grafanacloud.ConfigureProvider("0.0.1", &schema.Provider{TerraformVersion: "0.15"})
	provider, err := configureFunc(context.TODO(), resourceLocalData)
	require.Nil(t, err)

	_, ok := provider.(*grafanacloud.Provider)
	require.True(t, ok)
}

func getProvider(p *schema.Provider) *grafanacloud.Provider {
	return p.Meta().(*grafanacloud.Provider)
}
