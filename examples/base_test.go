package examples

import (
	"github.com/pulumi/providertest/providers"
	goprovider "github.com/pulumi/pulumi-go-provider"
	pulumirpc "github.com/pulumi/pulumi/sdk/v3/proto/go"

	"github.com/pulumi/pulumi-provider-boilerplate/provider"
)

var providerFactory = func(_ providers.PulumiTest) (pulumirpc.ResourceProviderServer, error) { //nolint:unused
	return goprovider.RawServer("provider-boilerplate", "1.0.0", provider.Provider())(nil)
}
