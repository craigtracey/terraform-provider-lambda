package provider

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"

	"github.com/h2non/gock"
)

func setupMocks(t *testing.T) {

	if os.Getenv("DEBUG_MOCKS") != "" {
		gock.Observe(func(req *http.Request, mock gock.Mock) {
			if req != nil {
				t.Log("===")
				t.Logf("%s %s\n", req.Method, req.URL)
				t.Logf("Headers: %+v\n", req.Header)
				if req.Body != nil {
					body, _ := ioutil.ReadAll(req.Body)
					t.Logf("Body: %s\n", body)
				}
				t.Log("===")
			}
		})
	}

	gock.New("https://example.com/").
		MatchHeaders(map[string]string{
			"Content-Type":  "application/json",
			"Authorization": "Bearer BAADF00D",
		}).
		Post("/instance-operations/launch").
		BodyString(`{
			"instance_type_name": "gpu_1x_a100",
			"name": "training-node-1",
			"region_name": "us-tx-1",
			"ssh_key_names": ["macbook-pro"],
			"file_system_names": ["shared-fs"]
		      }`).
		Reply(200).
		BodyString(`{
		"data": {
		  "instance_ids": [
		    "0920582c7ff041399e34823a0be62549"
		  ]
		}
		}`)

	gock.New("https://example.com").
		MatchHeaders(map[string]string{
			"Authorization": "Bearer BAADF00D",
		}).
		Get("/instances/0920582c7ff041399e34823a0be62549").
		Reply(200).
		BodyString(`{
	        "data": {
		  "file_system_names": [
		    "shared-fs"
		  ],
		  "hostname": "10-0-8-196.cloud.lambdalabs.com",
		  "id": "0920582c7ff041399e34823a0be62549",
		  "instance_type": {
		    "description": "1x RTX A100 (24 GB)",
		    "name": "gpu_1x_a100",
		    "price_cents_per_hour": 110,
		    "specs": {
		      "memory_gib": 800,
		      "storage_gib": 512,
		      "vcpus": 24
		    }
		  },
		  "ip": "10.10.10.1",
		  "jupyter_token": "53968f128c4a4489b688c2c0a181d083",
		  "jupyter_url": "https://jupyter-3ac4c5c6-9026-47d2-9a33-71efccbcd0ee.lambdaspaces.com/?token=53968f128c4a4489b688c2c0a181d083",
		  "name": "training-node-1",
		  "region": {
		    "description": "Austin, Texas",
		    "name": "us-tx-1"
		  },
		  "ssh_key_names": [
		    "macbook-pro"
		  ],
		  "status": "active"
		}
	      }`)

	gock.New("https://example.com").
		MatchHeaders(map[string]string{
			"Authorization": "Bearer BAADF00D",
			"Content-Type":  "application/json",
		}).
		Post("/instance-operations/terminate").
		BodyString(`{
		"instance_ids": [
		  "0920582c7ff041399e34823a0be62549"
		]
		}`).
		Reply(200).
		BodyString(`{
	      "data": {
		"file_system_names": [
		  "shared-fs"
		],
		"hostname": "10-0-8-196.cloud.lambdalabs.com",
		"id": "0920582c7ff041399e34823a0be62549",
		"instance_type": {
		  "description": "1x RTX A100 (24 GB)",
		  "name": "gpu_1x_a100",
		  "price_cents_per_hour": 110,
		  "specs": {
		    "memory_gib": 800,
		    "storage_gib": 512,
		    "vcpus": 24
		  }
		},
		"ip": "10.10.10.1",
		"jupyter_token": "53968f128c4a4489b688c2c0a181d083",
		"jupyter_url": "https://jupyter-3ac4c5c6-9026-47d2-9a33-71efccbcd0ee.lambdaspaces.com/?token=53968f128c4a4489b688c2c0a181d083",
		"name": "training-node-1",
		"region": {
		  "description": "Austin, Texas",
		  "name": "us-tx-1"
		},
		"ssh_key_names": [
		  "macbook-pro"
		],
		"status": "active"
	      }
	    }`)

}
