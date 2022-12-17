package provider

import (
	"encoding/json"
	"net/http"

	v1 "github.com/craigtracey/terraform-provider-lambda/pkg/api/v1"
)

func decodeErrorResponse(resp *http.Response) (*v1.Error, error) {
	var errorResponse v1.Error
	err := json.NewDecoder(resp.Body).Decode(&errorResponse)
	return &errorResponse, err
}