//go:build nodejs || all
// +build nodejs all

package examples

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

func TestNodejsExampleLifecycle(t *testing.T) {
	t.Skip("linking isn't working correctly")

	pt := pulumitest.NewPulumiTest(t, "nodejs",
		opttest.YarnLink("@mynamespace/provider-boilerplate"),
		opttest.AttachProviderServer("provider-boilerplate", providerFactory),
	)

	pt.Preview(t)
	pt.Up(t)
	pt.Destroy(t)
}
