package grafanacloud_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/form3tech-oss/terraform-provider-grafanacloud/grafanacloud"
)

func TestProvider(t *testing.T) {
	if err := grafanacloud.NewProvider("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func getProvider(p *schema.Provider) *grafanacloud.Provider {
	return p.Meta().(*grafanacloud.Provider)
}
