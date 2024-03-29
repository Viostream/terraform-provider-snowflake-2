package resources_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccWarehouseGrant(t *testing.T) {
	wName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	roleName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: warehouseGrantConfig(wName, roleName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("snowflake_warehouse_grant.test", "warehouse_name", wName),
					resource.TestCheckResourceAttr("snowflake_warehouse_grant.test", "privilege", "USAGE"),
				),
			},
			// IMPORT
			{
				ResourceName:      "snowflake_warehouse_grant.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func warehouseGrantConfig(n, role string) string {
	return fmt.Sprintf(`

resource "snowflake_warehouse" "test" {
  name      = "%v"
}

resource "snowflake_role" "test" {
  name = "%v"
}

resource "snowflake_warehouse_grant" "test" {
  warehouse_name = snowflake_warehouse.test.name
  roles          = [snowflake_role.test.name]
}
`, n, role)
}
