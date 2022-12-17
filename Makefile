VERSION = 1.2.0
PROVIDER_NAME = terraform-provider-lambda_v${VERSION}
API_SPEC = https://cloud.lambdalabs.com/static/api/v1/openapi.yaml

.PHONY: all
all: provider

.PHONY: gen
gen: deps
	curl "${API_SPEC}" -H "Accepts: application/yaml" -o openapi.yaml
	python scripts/annotate-openapi.py openapi.yaml
	oapi-codegen -old-config-style -package v1 -generate client,types openapi.yaml > pkg/api/v1/cli.gen.go

.PHONY: deps
deps:
	python -m pip install pyyaml
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.2

.PHONY: provider
provider: test gen
	go build -o ${PROVIDER_NAME} .

.PHONY: test
test: gen
	go test -v ./pkg/provider/