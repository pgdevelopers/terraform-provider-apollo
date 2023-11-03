// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGraphResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccGraphResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollo_graph.test", "configurable_attribute", "one"),
					resource.TestCheckResourceAttr("apollo_graph.test", "defaulted", "example value when not configured"),
					resource.TestCheckResourceAttr("apollo_graph.test", "id", "example-id"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "test",
				ImportState:       true,
				ImportStateVerify: true,
				// This is not normally necessary, but is here because this
				// example code does not have an actual upstream service.
				// Once the Read method is able to refresh information from
				// the upstream service, this can be removed.
				ImportStateVerifyIgnore: []string{"configurable_attribute", "defaulted"},
			},
			// Update and Read testing
			{
				Config: testAccGraphResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollo_graph.test", "configurable_attribute", "two"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccGraphResourceConfig(configurableAttribute string) string {
	return fmt.Sprintf(`
resource "apollo_graph" "test" {
  configurable_attribute = %[1]q
}
`, configurableAttribute)
}