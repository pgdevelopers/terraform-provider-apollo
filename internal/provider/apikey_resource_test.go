// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccaApikeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApikeyResourceConfig("one"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollo_apikey.test", "graph_id", "test"),
					resource.TestCheckResourceAttr("apollo_apikey.test", "key_name", "test"),
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
				Config: testAccApikeyResourceConfig("two"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("apollo_apikey.test", "graph_id", "test"),
					resource.TestCheckResourceAttr("apollo_apikey.test", "key_name", "test"),
				),
			},
			// Delete testing automatically occurs in TestCase
		},
	})
}

func testAccApikeyResourceConfig(configurableAttribute string) string {
	return fmt.Sprintln(`
resource "apollo_apikey" "test" {
	key_name = "test-key"
	graph_id = "test-graph"
}
`, configurableAttribute)
}
