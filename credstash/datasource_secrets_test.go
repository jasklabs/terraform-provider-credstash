package credstash

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataCredstashSecret_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCredstashSecret,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("credstash_secret.terraform", "value", "test"),
					resource.TestCheckResourceAttr("data.credstash_secret.terraform", "value", "test"),
					resource.TestCheckResourceAttr("data.credstash_secret.terraform-fake", "value", "fake"),
				),
			},
		},
	})
}

const testAccCheckCredstashSecret = `
resource "credstash_secret" "terraform" {
	name = "terraform-data-acc"
	value = "test"
}

data "credstash_secret" "terraform" {
	name    = "${credstash_secret.terraform.name}"
}

data "credstash_secret" "terraform-fake" {
	name    = "terraform-data-acc-fake"
	default = "fake"
}
`
