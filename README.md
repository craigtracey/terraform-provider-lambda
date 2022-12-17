Lambda GPU Cloud Terraform Provider
===
A Terraform provider for managing Lambda GPU Cloud resources.

Requirements
---
- [Terraform](https://www.terraform.io/downloads.html)
- [Go](https://golang.org/doc/install) 1.18 (to build the provider)

Example Usage
---
```
provider "lambda" {
  endpoint = "https://example.com"
  apikey   = "BAADF00D"
}

resource "lambda_instance" "training_node_1" {
  name          = "training-node-1"
  instance_type = "gpu_1x_a100"
  region        = "us-tx-1"
  ssh_key       = "macbook-pro"
  filesystem	= "shared-fs"
}
```

Building the Provider
---
Clone this repository and simply execute:
```bash
make
```

Testing the Provider
---
Testing the provider makes use of a mocked API endpoint. It will not create an resources against a Lambda API endpoint.

Executing the tests is as simple as:
```bash
export TF_ACC=1
make test
go test -v ./pkg/provider/
=== RUN   TestInstanceCreateDelete
=== PAUSE TestInstanceCreateDelete
=== RUN   TestProvider
--- PASS: TestProvider (0.00s)
=== CONT  TestInstanceCreateDelete
--- PASS: TestInstanceCreateDelete (1.39s)
PASS
ok      github.com/craigtracey/terraform-provider-lambda/pkg/provider   (cached)
```

Optionally, you may obtain additional details about the mocked requests by setting the `DEBUG_MOCKS` environment variable:

```bash
export TF_ACC=1
export DEBUG_MOCKS=1
make test
go test -v ./pkg/provider/
=== RUN   TestInstanceCreateDelete
=== PAUSE TestInstanceCreateDelete
=== RUN   TestProvider
--- PASS: TestProvider (0.00s)
=== CONT  TestInstanceCreateDelete
    mocks.go:17: ===
    mocks.go:18: POST https://example.com/instance-operations/launch
    mocks.go:19: Headers: map[Authorization:[Bearer BAADF00D] Content-Type:[application/json]]
    mocks.go:22: Body: {"file_system_names":["shared-fs"],"instance_type_name":"gpu_1x_a100","name":"training-node-1","region_name":"us-tx-1","ssh_key_names":["macbook-pro"]}
    mocks.go:24: ===
    mocks.go:17: ===
    mocks.go:18: GET https://example.com/instances/0920582c7ff041399e34823a0be62549
    mocks.go:19: Headers: map[Authorization:[Bearer BAADF00D]]
    mocks.go:24: ===
    mocks.go:17: ===
    mocks.go:18: GET https://example.com/instances/0920582c7ff041399e34823a0be62549
    mocks.go:19: Headers: map[Authorization:[Bearer BAADF00D]]
    mocks.go:24: ===
    mocks.go:17: ===
    mocks.go:18: POST https://example.com/instance-operations/terminate
    mocks.go:19: Headers: map[Authorization:[Bearer BAADF00D] Content-Type:[application/json]]
    mocks.go:22: Body: {"instance_ids":["0920582c7ff041399e34823a0be62549"]}
    mocks.go:24: ===
--- PASS: TestInstanceCreateDelete (1.39s)
PASS
ok      github.com/craigtracey/terraform-provider-lambda/pkg/provider   (cached)
```
