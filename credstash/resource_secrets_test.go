package credstash

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccCredstashSecret_basic(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCredstashSecretBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("credstash_secret.terraform", "name", "terraform-acc"),
					resource.TestCheckResourceAttr("credstash_secret.terraform", "value", "test"),
				),
			},
		},
	})
}

func TestAccCredstashSecret_overwrite(t *testing.T) {
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckCredstashSecretOverwrite,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("credstash_secret.terraform1", "name", "terraform-acc-overwrite"),
					resource.TestCheckResourceAttr("credstash_secret.terraform1", "value", "test"),
					resource.TestCheckResourceAttr("credstash_secret.terraform2", "name", "terraform-acc-overwrite"),
					resource.TestCheckResourceAttr("credstash_secret.terraform2", "value", "test"),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

const testAccCheckCredstashSecretBasic = `
resource "credstash_secret" "terraform" {
	name = "terraform-acc"
	value = "test"
}
`

const testAccCheckCredstashSecretOverwrite = `
resource "credstash_secret" "terraform1" {
	name = "terraform-acc-overwrite"
	value = "test"
	overwrite = true
}

resource "credstash_secret" "terraform2" {
	name = "${credstash_secret.terraform1.name}"
	value = "overwrite"
	overwrite = false
}
`
