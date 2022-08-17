package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
)

// getWorkspaceID extracts the workspace ID attribute generically from a Split resource.
func getWorkspaceID(d *schema.ResourceData) string {
	var workspaceID string
	if v, ok := d.GetOk("workspace_id"); ok {
		vs := v.(string)
		log.Printf("[DEBUG] workspace_id: %s", vs)
		workspaceID = vs
	}

	return workspaceID
}

// getTrafficTypeID extracts the traffic type ID attribute generically from a Split resource.
func getTrafficTypeID(d *schema.ResourceData) string {
	var trafficTypeID string
	if v, ok := d.GetOk("traffic_type_id"); ok {
		vs := v.(string)
		log.Printf("[DEBUG] traffic_type_id: %s", vs)
		trafficTypeID = vs
	}

	return trafficTypeID
}

// getEnvironmentID extracts the environment ID attribute generically from a Split resource.
func getEnvironmentID(d *schema.ResourceData) string {
	var envID string
	if v, ok := d.GetOk("environment_id"); ok {
		vs := v.(string)
		log.Printf("[DEBUG] environment_id: %s", vs)
		envID = vs
	}

	return envID
}

// getSplitName extracts the split name attribute generically from a Split resource.
func getSplitName(d *schema.ResourceData) string {
	var n string
	if v, ok := d.GetOk("split_name"); ok {
		vs := v.(string)
		log.Printf("[DEBUG] split_name: %s", vs)
		n = vs
	}

	return n
}

func parseCompositeID(id string, numOfSplits int) ([]string, error) {
	parts := strings.SplitN(id, ":", numOfSplits)

	if len(parts) != numOfSplits {
		return nil, fmt.Errorf("error: import composite ID requires %d parts separated by a colon (x:y). "+
			"Please check resource documentation for more information", numOfSplits)
	}
	return parts, nil
}
