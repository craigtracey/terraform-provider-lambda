package provider

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	v1 "github.com/craigtracey/terraform-provider-lambda/pkg/api/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

// FIXME: update these "constants"
var lambdaInstanceTypes = []string{"gpu_1x_a100", "gpu_4x_a6000"}
var lambdaRegions = []string{"us-az-1", "us-tx-1"}

func resourceInstance() *schema.Resource {

	return &schema.Resource{
		CreateContext: resourceInstanceCreate,
		ReadContext:   resourceInstanceRead,
		DeleteContext: resourceInstanceDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"region": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(lambdaRegions, true),
			},
			"instance_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice(lambdaInstanceTypes, true),
			},
			"ssh_key": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"filesystem": {
				Type:         schema.TypeString,
				Required:     false,
				ForceNew:     true,
				Optional:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func resourceInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*v1.Client)

	name := d.Get("name").(string)
	regionName := d.Get("region").(string)
	instanceType := d.Get("instance_type").(string)

	instanceBody := v1.LaunchInstanceJSONBody{
		Name:             &name,
		InstanceTypeName: instanceType,
		RegionName:       regionName,
	}

	// the API only accepts a single ssh key/filesystem,
	// so only expose these as strings for the moment
	if attr, ok := d.GetOk("ssh_key"); ok {
		instanceBody.SshKeyNames = []string{attr.(string)}
	}
	if attr, ok := d.GetOk("filesystem"); ok {
		instanceBody.FileSystemNames = &[]string{attr.(string)}
	}

	resp, err := client.LaunchInstance(ctx, v1.LaunchInstanceJSONRequestBody(instanceBody))
	if err != nil {
		return diag.Errorf("Failed to launch instance with error: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errorResponse, err := decodeErrorResponse(resp)
		if err != nil {
			return diag.Errorf("Failed to decode API error reponse: %s", err)
		}
		return diag.Errorf("Failed to launch instance with error: %s", errorResponse.Message)
	}

	var launchInstanceResponse v1.LaunchAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&launchInstanceResponse)
	if err != nil {
		return diag.Errorf("Failed to decode API reponse: %s", err)
	}

	// again, we only support creating a single instance per resource
	d.SetId(launchInstanceResponse.Data.InstanceIds[0])
	log.Printf("[INFO] Lambda Instance ID: %s", d.Id())

	return nil
}

func resourceInstanceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*v1.Client)

	resp, err := client.GetInstance(ctx, d.Id())
	if err != nil {
		return diag.Errorf("Failed to fetch instance data: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("[WARN] Lambda Instance not found: %s", d.Id())
			d.SetId("")
			return nil
		}

		errorResponse, err := decodeErrorResponse(resp)
		if err != nil {
			return diag.Errorf("Failed to decode API error reponse: %s", err)
		}
		return diag.Errorf("Failed to fetch instance data: %s", errorResponse.Message)
	}

	var instanceData v1.InstanceAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&instanceData)
	if err != nil {
		return diag.Errorf("Failed to decode instance data: %s", err)
	}
	d.Set("name", instanceData.Data.Name)
	d.Set("region", instanceData.Data.Region.Name)
	d.Set("instance_type", instanceData.Data.InstanceType.Name)

	// FIXME: these cases are not great because the API is currently limited to
	// single values on instance launch, but semantically allows for an array of values
	if len(instanceData.Data.FileSystemNames) == 1 {
		d.Set("filesystem", instanceData.Data.FileSystemNames[0])
	} else if len(instanceData.Data.FileSystemNames) > 1 {
		diag.Errorf("[ERROR] Multiple filesystems provisioned for instance")
	}
	if len(instanceData.Data.SshKeyNames) == 1 {
		d.Set("ssh_key", instanceData.Data.SshKeyNames[0])
	} else if len(instanceData.Data.SshKeyNames) > 1 {
		diag.Errorf("[ERROR] Multiple ssh keys provisioned for instance")
	}

	return nil
}

func resourceInstanceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*v1.Client)

	id := d.Id()
	body := v1.Terminate{
		InstanceIds: []string{
			id,
		},
	}
	resp, err := client.TerminateInstance(ctx, v1.TerminateInstanceJSONRequestBody(body))
	if err != nil {
		return diag.Errorf("Failed to delete instance: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("[WARN] Lambda Instance not found: %s", id)
			d.SetId("")
			return nil
		}

		errorResponse, err := decodeErrorResponse(resp)
		if err != nil {
			return diag.Errorf("Failed to decode API error reponse: %s", err)
		}
		return diag.Errorf("Failed to delete instance %s: %s", id, errorResponse.Message)
	}

	return nil
}
