package grafanacloud_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceStacks_Basic(t *testing.T) {
	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStacksConfig(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.grafanacloud_stacks.test", "stacks.0.name", name),
					resource.TestCheckResourceAttr("data.grafanacloud_stacks.test", "stacks.0.slug", name+"slug"),
					resource.TestCheckResourceAttrSet("data.grafanacloud_stacks.test", "stacks.0.prometheus_url"),
					resource.TestCheckResourceAttrSet("data.grafanacloud_stacks.test", "stacks.0.prometheus_user_id"),
					resource.TestCheckResourceAttrSet("data.grafanacloud_stacks.test", "stacks.0.alertmanager_user_id"),
				),
			},
		},
	})
}

func TestAccDataSourceStacks_Multiple(t *testing.T) {
	name1 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	name2 := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceStacksMultipleConfig(name1, name2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.grafanacloud_stacks.test", "stacks.#", "2"),
				),
			},
		},
	})
}

func testAccDataSourceStacksConfig(name string) string {
	return fmt.Sprintf(`
resource "grafanacloud_stack" "foo" {
  name = "%s"
  slug = "%sslug"
}

data "grafanacloud_stacks" "test" {
	depends_on = [
		grafanacloud_stack.foo
	]
}
`, name, name)
}

func testAccDataSourceStacksMultipleConfig(name1, name2 string) string {
	return fmt.Sprintf(`
resource "grafanacloud_stack" "foo" {
  name = "%s"
  slug = "%sslug"
}

resource "grafanacloud_stack" "bar" {
  name = "%s"
  slug = "%sslug"
}

data "grafanacloud_stacks" "test" {
	depends_on = [
		grafanacloud_stack.foo,
		grafanacloud_stack.bar
	]
}
`, name1, name1, name2, name2)
}
