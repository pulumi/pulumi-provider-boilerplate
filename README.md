# xyz Pulumi Provider

This repo is a boilerplate showing how to create a native Pulumi provider.  You can search-replace `xyz` with the name of your desired provider as a starting point for creating a provider that manages resources in the target cloud.

Most of the code for the provider implementation is in `pkg/provider/provider.go`.  

An example of using the single resource defined in this example is in `examples/simple`.

This repo does not yet contain a code-generator for generating SDKs in each Pulumi language (TypeScript/Python/Go/.NET).  Instead, the SDK for TypeScript for this example is included in the example in `examples/simple/xyz.ts`.

Note that the generated provider plugin (`pulumi-resource-xyz`) must be on your `PATH` to be used by Pulumi deployments.  If creating a provider for distribution to other users, you should ensure they install this plugin to their `PATH`.


## Build and Test

```bash
# build
$ go install ./cmd/pulumi-resource-xyz

# test
$ cd examples/simple
$ npm install
$ pulumi stack init test
$ pulumi up
```

## References

Other resoruces for learning about the Pulumi resource model:
* [Pulumi Kubernetes provider](https://github.com/pulumi/pulumi-kubernetes/blob/master/pkg/provider/provider.go)
* [Pulumi Terraform Remote State provider](https://github.com/pulumi/pulumi-terraform/blob/master/pkg/provider/provider.go)
* [Dynamic Providers](https://www.pulumi.com/docs/intro/concepts/programming-model/#dynamicproviders)
