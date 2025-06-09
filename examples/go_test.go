//go:build go || all
// +build go all

package examples

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/pulumi/providertest/pulumitest"
	"github.com/pulumi/providertest/pulumitest/opttest"
	"github.com/stretchr/testify/require"
)

func TestGoExampleLifecycle(t *testing.T) {
	cwd, err := os.Getwd()
	require.NoError(t, err)

	module := filepath.Join(cwd, "../sdk/go/pulumi-provider-boilerplate")
	pt := pulumitest.NewPulumiTest(t, "go",
		opttest.GoModReplacement("github.com/pulumi/pulumi-provider-boilerplate/sdk/go/pulumi-provider-boilerplate", module),
		opttest.AttachProviderServer("provider-boilerplate", providerFactory),
		opttest.SkipInstall(),
	)

	pt.Preview(t)
	pt.Up(t)
	pt.Destroy(t)
}
