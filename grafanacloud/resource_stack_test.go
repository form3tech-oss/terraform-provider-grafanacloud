package grafanacloud_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccStack_Basic(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	urlRegexpString := fmt.Sprintf("http://.+/grafana/%s-slug", resourceName)

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStackConfig(resourceName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStackExists("grafanacloud_stack.test"),
					resource.TestCheckResourceAttrSet("grafanacloud_stack.test", "id"),
					resource.TestCheckResourceAttr("grafanacloud_stack.test", "name", resourceName),
					resource.TestCheckResourceAttr("grafanacloud_stack.test", "slug", resourceName+"-slug"),
					resource.TestMatchResourceAttr("grafanacloud_stack.test", "url", regexp.MustCompile(urlRegexpString)),
				),
			},
		},
	})
}

func TestAccStack_URL(t *testing.T) {
	resourceName := acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum)
	url := "https://my.grafana.instance"

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckStackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccStackConfigURL(resourceName, url),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckStackExists("grafanacloud_stack.test"),
					resource.TestCheckResourceAttrSet("grafanacloud_stack.test", "id"),
					resource.TestCheckResourceAttr("grafanacloud_stack.test", "name", resourceName),
					resource.TestCheckResourceAttr("grafanacloud_stack.test", "slug", resourceName+"-slug"),
					resource.TestCheckResourceAttr("grafanacloud_stack.test", "url", url),
				),
			},
		},
	})
}

func testAccCheckStackExists(resourceName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource `%s` not found", resourceName)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource `%s` has no ID set", resourceName)
		}

		p := getProvider(testAccProvider)
		stack, err := p.Client.GetStack(p.Organisation, rs.Primary.Attributes["slug"])
		if err != nil {
			return err
		}

		if stack == nil {
			return fmt.Errorf("resource `%s` not found via API", resourceName)
		}

		return nil
	}
}

func testAccCheckStackDestroy(s *terraform.State) error {
	p := getProvider(testAccProvider)

	for name, rs := range s.RootModule().Resources {
		if rs.Type != "grafanacloud_stack" {
			continue
		}

		stack, err := p.Client.GetStack(p.Organisation, rs.Primary.Attributes["slug"])
		if err != nil {
			return err
		}

		if stack != nil {
			return fmt.Errorf("resource `%s` with ID `%s` still exists after destroy", name, rs.Primary.ID)
		}
	}

	return nil
}

func testAccStackConfig(resourceName string) string {
	return fmt.Sprintf(`
resource "grafanacloud_stack" "test" {
  name = "%s"
  slug = "%s-slug"
}
`, resourceName, resourceName)
}

func testAccStackConfigURL(resourceName, url string) string {
	return fmt.Sprintf(`
resource "grafanacloud_stack" "test" {
  name = "%s"
  slug = "%s-slug"
	url  = "%s"
}
`, resourceName, resourceName, url)
}
