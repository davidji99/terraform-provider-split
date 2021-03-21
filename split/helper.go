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

func parseCompositeID(id string, numOfSplits int) ([]string, error) {
	parts := strings.SplitN(id, ":", numOfSplits)

	if len(parts) != numOfSplits {
		return nil, fmt.Errorf("Error: import composite ID requires %d parts separated by a colon (x:y). "+
			"Please check resource documentation for more information.", numOfSplits)
	}
	return parts, nil
}
