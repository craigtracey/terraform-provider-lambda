package provider

import (
	"context"
	"fmt"
	"net/http"

	v1 "github.com/craigtracey/terraform-provider-lambda/pkg/api/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func Provider() *schema.Provider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("LAMBDA_API_ENDPOINT", nil),
				Description:  "URL of Lambda GPU Cloud API endpoint.",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"apikey": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("LAMBDA_API_KEY", nil),
				Description:  "Lambda GPU Cloud API access key",
				ValidateFunc: validation.NoZeroValues,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"lambda_instance": resourceInstance(),
		},
	}
	p.ConfigureFunc = providerConfigure(p)
	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		endpoint := d.Get("endpoint").(string)
		apikey := d.Get("apikey").(string)
		client, err := v1.NewClient(endpoint, v1.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apikey))
			return nil
		}))
		if err != nil {
			return nil, err
		}
		return client, nil
	}
}
