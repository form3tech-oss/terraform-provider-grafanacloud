package grafanacloud_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/naag/terraform-provider-grafanacloud/grafanacloud"
	"github.com/stretchr/testify/require"
)

func TestValidatePortalApiKeyRole(t *testing.T) {
	fn := grafanacloud.ValidatePortalApiKeyRole()

	var tests = []struct {
		role  string
		valid bool
	}{
		{"Viewer", true},
		{"Editor", true},
		{"Admin", true},
		{"MetricsPublisher", true},
		{"PluginPublisher", true},
		{"Invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			warn, err := fn(tt.role, "role")
			if tt.valid {
				require.Empty(t, warn)
				require.Empty(t, err)
			} else {
				require.Empty(t, warn)
				require.NotEmpty(t, err)
			}
		})
	}
}

func TestAccPortalApiKey_Basic(t *testing.T) {
	var tests = []struct {
		role string
	}{
		{"Viewer"},
		{"Editor"},
		{"Admin"},
		{"MetricsPublisher"},
		{"PluginPublisher"},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

			resource.Test(t, resource.TestCase{
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckPortalAPIKeyDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccPortalAPIKeyConfig(resourceName, tt.role),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckPortalAPIKeyExists("grafanacloud_portal_api_key.test"),
							resource.TestCheckResourceAttrSet("grafanacloud_portal_api_key.test", "id"),
							resource.TestCheckResourceAttrSet("grafanacloud_portal_api_key.test", "key"),
							resource.TestCheckResourceAttr("grafanacloud_portal_api_key.test", "name", resourceName),
							resource.TestCheckResourceAttr("grafanacloud_portal_api_key.test", "role", tt.role),
						),
					},
				},
			})
		})
	}
}

func testAccCheckPortalAPIKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource `%s` not found", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource `%s` has no ID set", resourceName)
		}

		p := getProvider(testAccProvider)
		res, err := p.Client.ListAPIKeys(p.Organisation)
		if err != nil {
			return err
		}

		apiKey := res.FindByName(rs.Primary.ID)
		if apiKey == nil {
			return fmt.Errorf("resource `%s` not found via API", resourceName)
		}

		return nil
	}
}

func testAccCheckPortalAPIKeyDestroy(s *terraform.State) error {
	p := getProvider(testAccProvider)

	for name, rs := range s.RootModule().Resources {
		if rs.Type != "grafanacloud_portal_api_key" {
			continue
		}

		res, err := p.Client.ListAPIKeys(p.Organisation)
		if err != nil {
			return err
		}

		apiKey := res.FindByName(rs.Primary.ID)
		if apiKey != nil {
			return fmt.Errorf("resource `%s` with ID `%s` still exists after destroy", name, rs.Primary.ID)
		}
	}

	return nil
}

func testAccPortalAPIKeyConfig(resourceName, role string) string {
	return fmt.Sprintf(`
resource "grafanacloud_portal_api_key" "test" {
  name = "%s"
  role = "%s"
}
`, resourceName, role)
}
