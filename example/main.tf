provider "lambda" {
  endpoint = "https://example.com"
  apikey   = "BAADF00D"
}

resource "lambda_instance" "training_node_1" {
  name          = "training-node-1"
  instance_type = "gpu_1x_a100"
  region        = "us-tx-1"
  ssh_key       = "macbook-pro"
  filesystem    = "shared-fs"
}
