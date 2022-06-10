# Pulumi Native Provider Boilerplate

This repo is a boilerplate showing how to create a native Pulumi provider.

### Background

A Pulumi provider is a piece of the Pulumi toolchain that implements a gRPC server which allows for the Pulumi engine to send requests to the provider.

```go
type xyzProvider struct {
	host    *provider.HostClient
	name    string
	version string
	schema  []byte
}
```

There are several parts to implement to make this work.
We encourage you to read the links below for background.
Many, but not all, providers implement a cloud provider service of some kind, such as Azure, AWS, Civo, to name a few.
The important part is that providers allow for Pulumi to interface with the upstream tool's API in a programmatic way.
The key functionality a provider brings is the implementation of behavior for the cloud resources that are being created during a Pulumi program.
A note on jargon: when we speak of a "native" provider, we mean that all implementation is native to Pulumi, as opposed to Terraform based providers TODO: link
There are no wrappers around Terraform infrastucture providers here.





This repository is part of the [guide for authoring and publishing a Pulumi Package](https://www.pulumi.com/docs/guides/pulumi-packages/how-to-author).

Learn about the concepts behind [Pulumi Packages](https://www.pulumi.com/docs/guides/pulumi-packages/#pulumi-packages).

## Creating a Pulumi Native Provider

The following instructions cover providers maintained by Pulumi (denoted with a "Pulumi Official" checkmark on the Pulumi registry).
In the future, we will add instruction for providers published and maintained by the Pulumi community, referred to as "third-party" providers.

This boilerplate creates a Pulumi-owned provider named `xyz`. For a full example please see [random-native](TODO:link)

### Prerequisites

Ensure the following tools are installed and present in your `$PATH`:

* [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
* [Go 1.17](https://golang.org/dl/) or 1.latest
* [NodeJS](https://nodejs.org/en/) 14.x.  We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations.
* [Yarn](https://yarnpkg.com/)
* [TypeScript](https://www.typescriptlang.org/)
* [Python](https://www.python.org/downloads/) (called as `python3`).  For recent versions of MacOS, the system-installed version is fine.
* [.NET](https://dotnet.microsoft.com/download)


### Creating and Initializing the Repository

Pulumi offers this repository as a [GitHub template repository](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template) for convenience.  From this repository:

1. Click "Use this template".
1. Set the following options:
   * Owner: pulumi 
   * Repository name: pulumi-xyz-native (replace "xyz" with the name of your provider)
   * Description: Pulumi provider for xyz
   * Repository type: Public
1. Clone the generated repository.

From the templated repository:

1. Search-replace `xyz` with the name of your desired provider.

2. Build the provider and install the plugin

   ```bash
   $ make build install
   ```
   
   This will:
   1. create the SDK codegen binary and place it in a `./bin` folder (gitignored)
   2. create the provider binary and place it in the `./bin` folder
   3. generate the dotnet, Go, Node, and Python SDKs and place them in the `./sdk` folder
   4. install the provider on your machine.

3. Test against the example
   
   ```bash
   $ cd examples/simple
   $ yarn link @pulumi/xyz
   $ yarn install
   $ pulumi stack init test
   $ pulumi up
   ```

Now that you have completed all of the above steps, you have a working provider that generates a random string for you.
Here is an example of a more complete native random provider (TODO: link to random-native).

#### A brief repository overview

You now have:

1. A `provider/` folder containing the building and implementation logic
   1. `cmd/`
      1. `pulumi-gen-xyz/` - generates language SDKs from the schema
      2. `pulumi-resource-xyz/` - holds the package schema, injects the package version, and starts the gRPC server
   2. `pkg`
      1. `provider` - holds the gRPC methods (and for now, the sample implementation logic) required by the Pulumi engine
      2. `version` - semver package to be consumed by build processes
2. `deployment-templates` - a set of files to help you around deployment and publication
3. `sdk` - holds the generated code libraries created by `pulumi-gen-xyz/main.go`
4. `examples` a folder of Pulumi programs to try locally and/or use in CI.
5. A `Makefile` and this `README`.

#### About `schema.json`

The JSON schema file is what enables `pulumi-gen-xyz` to create language-specific SDKs.
We read in the schema file and use Pulumi's codegen package (TODO: link) to generate language-specific SDKs.
For a native provider, you must create (usually by hand) a schema that describes your provider's resources.
TODO: link to schema docs
TODO: provide rough outline of schema

### Implement the gRPC methods

Once you have a schema that describes your provider, you will need to implement the desired gRPC methods.
You will find a mostly blank implementation of these in `pkg/provider/provider.go`.
Note that these methods do not necessarily link 1:1 to the Pulumi CLI commands.

#### The Necessary Methods

You need to implement the following methods _for each resource_:

1. Check
2. Diff
3. Create
4. Read
5. Update
6. Delete
7. GetPluginInfo (already implemented)
8. GetSchema (already implemented)

TODO: link to a native provider implementation with multiple resources

#### Additional Methods

The provider interface (TODO: link to info in p/p) includes a few more gRPC methods that you may need to implement and can read more about.

You can now repeat the steps for build, install, and test

## Configuring CI and releases

We have centralized CI management for our providers. [Follow instructions here to implement them](TODO: add CI mgmt link with instructions)

### Third-party providers

1. Follow the instructions laid out in the [deployment templates](./deployment-templates/README-DEPLOYMENT.md).

## References

Other resources for learning about the Pulumi resource model:
* [Pulumi Kubernetes provider](https://github.com/pulumi/pulumi-kubernetes/blob/master/provider/pkg/provider/provider.go)
* [Pulumi Terraform Remote State provider](https://github.com/pulumi/pulumi-terraform/blob/master/provider/cmd/pulumi-resource-terraform/provider.go)
* [Dynamic Providers](https://www.pulumi.com/docs/intro/concepts/programming-model/#dynamicproviders)
