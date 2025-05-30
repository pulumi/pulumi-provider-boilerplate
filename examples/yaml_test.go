//go:build yaml || all
// +build yaml all

package examples

import (
	"testing"

	"github.com/pulumi/providertest"
	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/assertpreview"
	"github.com/pulumi/providertest/pulumitest/opttest"
)

func TestYAMLExampleLifecycle(t *testing.T) {
	pt := pulumitest.NewPulumiTest(t, "yaml",
		opttest.AttachProviderServer("provider-boilerplate", providerFactory),
		opttest.SkipInstall(),
	)

	pt.Preview(t)
	pt.Up(t)
	pt.Destroy(t)
}

func TestYAMLExampleUpgrade(t *testing.T) {
	pt := pulumitest.NewPulumiTest(t, "yaml",
		opttest.AttachProviderServer("provider-boilerplate", providerFactory),
		opttest.SkipInstall(),
	)
	previewResult := providertest.PreviewProviderUpgrade(t, pt, "provider-boilerplate", "0.0.1")

	assertpreview.HasNoChanges(t, previewResult)
}
