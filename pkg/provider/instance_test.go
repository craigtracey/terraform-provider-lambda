package provider

import (
	"os"
	"testing"

	"github.com/h2non/gock"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestInstanceCreateDelete(t *testing.T) {

	factories := map[string]func() (*schema.Provider, error){
		"lambda": func() (*schema.Provider, error) {
			os.Setenv("LAMBDA_API_ENDPOINT", "https://example.com")
			os.Setenv("LAMBDA_API_KEY", "BAADF00D")
			p := Provider()
			return p, nil
		},
	}
	defer gock.Off()


	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: factories,
		Steps: []resource.TestStep{
			{
				PreConfig: func() { setupMocks(t) },
				Config:    basicInstanceResource,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambda_instance.training_node_1", "name", "training-node-1"),
				),
			},
			{
				PreConfig: func() { setupMocks(t) },
				Config:    basicInstanceResource,
				Destroy:   true,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("lambda_instance.training_node_1", "name", "training-node-1"),
				),
			},
		},
	})
}

const basicInstanceResource = `resource "lambda_instance" "training_node_1" {
  name          = "training-node-1"
  instance_type = "gpu_1x_a100"
  region        = "us-tx-1"
  ssh_key       = "macbook-pro"
  filesystem	= "shared-fs"
}`
