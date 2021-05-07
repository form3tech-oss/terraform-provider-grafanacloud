package grafanacloud_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceStack_Basic(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStackConfig(resourceName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.grafanacloud_stack.test", "id"),
					resource.TestCheckResourceAttr("data.grafanacloud_stack.test", "name", resourceName),
					resource.TestCheckResourceAttr("data.grafanacloud_stack.test", "slug", resourceName+"slug"),
					resource.TestCheckResourceAttrSet("data.grafanacloud_stack.test", "prometheus_url"),
					resource.TestCheckResourceAttrSet("data.grafanacloud_stack.test", "prometheus_user_id"),
					resource.TestCheckResourceAttrSet("data.grafanacloud_stack.test", "alertmanager_user_id"),
				),
			},
		},
	})
}

func testAccDataSourceStackConfig(resourceName string) string {
	return fmt.Sprintf(`
resource "grafanacloud_stack" "test" {
  name = "%s"
  slug = "%sslug"
}

data "grafanacloud_stack" "test" {
  slug = grafanacloud_stack.test.slug
}
`, resourceName, resourceName)
}
