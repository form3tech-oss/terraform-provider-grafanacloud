package grafanacloud_test

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/form3tech-oss/terraform-provider-grafanacloud/grafanacloud"
	"github.com/stretchr/testify/require"
)

func TestValidateGrafanaApiKeyRole(t *testing.T) {
	fn := grafanacloud.ValidateGrafanaApiKeyRole()

	var tests = []struct {
		role  string
		valid bool
	}{
		{"Viewer", true},
		{"Editor", true},
		{"Admin", true},
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

func TestAccGrafanaApiKey_Basic(t *testing.T) {
	var tests = []struct {
		role string
	}{
		{"Viewer"},
		{"Editor"},
		{"Admin"},
	}

	for _, tt := range tests {
		t.Run(tt.role, func(t *testing.T) {
			resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

			resource.Test(t, resource.TestCase{
				Providers:    testAccProviders,
				CheckDestroy: testAccCheckGrafanaAPIKeyDestroy,
				Steps: []resource.TestStep{
					{
						Config: testAccGrafanaAPIKeyConfig(resourceName, tt.role),
						Check: resource.ComposeTestCheckFunc(
							testAccCheckGrafanaAPIKeyExists("grafanacloud_grafana_api_key.test"),
							resource.TestCheckResourceAttrSet("grafanacloud_grafana_api_key.test", "id"),
							resource.TestCheckResourceAttrSet("grafanacloud_grafana_api_key.test", "key"),
							resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "name", resourceName),
							resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "stack", "dummystack"),
							resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "role", tt.role),
							resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "expiration", ""),
							resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "is_expired", "false"),
							resource.TestCheckNoResourceAttr("grafanacloud_grafana_api_key.test", "seconds_to_live"),
						),
					},
				},
			})
		})
	}
}

func TestAccGrafanaApiKey_Expiring(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGrafanaAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGrafanaAPIKeyConfigExpiring(resourceName, 10),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaAPIKeyExists("grafanacloud_grafana_api_key.test"),
					resource.TestCheckResourceAttrSet("grafanacloud_grafana_api_key.test", "id"),
					resource.TestCheckResourceAttrSet("grafanacloud_grafana_api_key.test", "key"),
					resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "name", resourceName),
					resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "stack", "dummystack"),
					resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "role", "Viewer"),
					resource.TestMatchResourceAttr("grafanacloud_grafana_api_key.test", "expiration", regexp.MustCompile("^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}.*$")),
					resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "is_expired", "false"),
					resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "seconds_to_live", "10"),
				),
			},
		},
	})
}

func TestAccGrafanaApiKey_Expired(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGrafanaAPIKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGrafanaAPIKeyConfigExpiring(resourceName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaAPIKeyExists("grafanacloud_grafana_api_key.test"),
					resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "is_expired", "false"),
				),
			},
			{
				Config: testAccGrafanaAPIKeyConfigExpiring(resourceName, 2),
				// This is supposed to recreate the now expired API key
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeTestCheckFunc(
					testAccSleep(3*time.Second),
					resource.TestCheckResourceAttr("grafanacloud_grafana_api_key.test", "is_expired", "false"),
				),
			},
			{
				Config: testAccGrafanaAPIKeyConfigExpiring(resourceName, 2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGrafanaAPIKeyExists("grafanacloud_grafana_api_key.test"),
				),
			},
		},
	})
}

func testAccCheckGrafanaAPIKeyExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ctx := context.Background()

		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource `%s` not found", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource `%s` has no ID set", resourceName)
		}

		p := getProvider(testAccProvider)
		gc, cleanup, err := p.Client.GetAuthedGrafanaClient(ctx, p.Organisation, rs.Primary.Attributes["stack"])
		if err != nil {
			return err
		}

		if cleanup != nil {
			defer cleanup()
		}

		res, err := gc.ListAPIKeys(ctx, true)
		if err != nil {
			return err
		}

		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		apiKey := res.FindByID(id)
		if apiKey == nil {
			return fmt.Errorf("resource `%s` not found via API", resourceName)
		}

		return nil
	}
}

func testAccCheckGrafanaAPIKeyDestroy(s *terraform.State) error {
	ctx := context.Background()
	p := getProvider(testAccProvider)

	for name, rs := range s.RootModule().Resources {
		if rs.Type != "grafanacloud_grafana_api_key" {
			continue
		}

		res, err := p.Client.ListAPIKeys(ctx, p.Organisation)
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

func testAccGrafanaAPIKeyConfig(resourceName, role string) string {
	return fmt.Sprintf(`
resource "grafanacloud_grafana_api_key" "test" {
  name = "%s"
  role = "%s"
	stack = grafanacloud_stack.test.slug
}

resource "grafanacloud_stack" "test" {
  name = "dummy-stack"
	slug = "dummystack"
}
`, resourceName, role)
}

func testAccGrafanaAPIKeyConfigExpiring(resourceName string, secondsToLive int) string {
	return fmt.Sprintf(`
resource "grafanacloud_grafana_api_key" "test" {
  name = "%s"
  role = "Viewer"
	stack = grafanacloud_stack.test.slug
	seconds_to_live = %d
}

resource "grafanacloud_stack" "test" {
  name = "dummy-stack"
	slug = "dummystack"
}
`, resourceName, secondsToLive)
}

func testAccSleep(d time.Duration) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(d)
		return nil
	}
}
