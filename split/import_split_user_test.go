package split

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"strings"
	"testing"
)

func TestAccSplitUser_importBasic(t *testing.T) {
	email := testAccConfig.GetUserEmailorSkip(t)
	emailSplit := strings.Split(email, "@")
	emailFormatted := fmt.Sprintf("%s+%s@%s", emailSplit[0], acctest.RandString(8), emailSplit[1])

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		ProviderFactories: testAccProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckSplitUser_basic(emailFormatted),
			},
			{
				ResourceName:      "split_user.foobar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}
