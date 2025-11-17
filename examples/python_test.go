//go:build python || all
// +build python all

package examples

import (
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

func TestPython(t *testing.T) {
	pt := pulumitest.NewPulumiTest(t, "python",
		opttest.YarnLink("@pulumi/provider-boilerplate"),
		opttest.AttachProviderServer("provider-boilerplate", providerFactory),
	)

	pt.Preview(t)
	pt.Up(t)
	pt.Destroy(t)
}
