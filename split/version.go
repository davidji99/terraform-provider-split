// +build ignore

// This a hack to populate the version in the custom binary file as this provider is not official.

package main

import (
	"fmt"
	"github.com/davidji99/terraform-provider-split/version"
)

var ver = version.ProviderVersion

func main() {
	fmt.Println(ver)
}
