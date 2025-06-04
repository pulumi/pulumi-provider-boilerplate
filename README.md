# Pulumi Native Provider Boilerplate

This repository is a boilerplate showing how to create and locally test a native Pulumi provider (with examples of both CustomResource and ComponentResource [resource types](https://www.pulumi.com/docs/iac/concepts/resources/)). 

## Authoring a Pulumi Native Provider

This boilerplate creates a working Pulumi-owned provider named `provider-boilerplate`.
It implements a random number generator that you can [build and test out for yourself](#test-against-the-example) and then replace the Random code with code specific to your provider.


### Prerequisites

You will need to ensure the following tools are installed and present in your `$PATH`:

* [`pulumictl`](https://github.com/pulumi/pulumictl#installation)
* [Go 1.21](https://golang.org/dl/) or 1.latest
* [NodeJS](https://nodejs.org/en/) 14.x.  We recommend using [nvm](https://github.com/nvm-sh/nvm) to manage NodeJS installations.
* [Yarn](https://yarnpkg.com/)
* [TypeScript](https://www.typescriptlang.org/)
* [Python](https://www.python.org/downloads/) (called as `python3`).  For recent versions of MacOS, the system-installed version is fine.
* [.NET](https://dotnet.microsoft.com/download)


### Build & test the boilerplate provider

1. Run `make build install` to build and install the provider.
1. Run `make gen_examples` to generate the example programs in `examples/` off of the source `examples/yaml` example program.
1. Run `make up` to run the example program in `examples/yaml`.
1. Run `make down` to tear down the example program.

### Creating a new provider repository

Pulumi offers this repository as a [GitHub template repository](https://docs.github.com/en/repositories/creating-and-managing-repositories/creating-a-repository-from-a-template) for convenience.  From this repository:

1. Click "Use this template".
1. Set the following options:
   * Owner: pulumi 
   * Repository name: pulumi-provider-boilerplate (replace "provider-boilerplate" with the name of your provider)
   * Description: Pulumi provider for xyz
   * Repository type: Public
1. Clone the generated repository.

From the templated repository:

1. Run the following command to update files to use the name of your provider (third-party: use your GitHub organization/username):

    ```bash
    make prepare NAME=foo ORG=myorg REPOSITORY=github.com/myorg/pulumi-foo
    ```

   This will do the following:
   - rename folders in `provider/cmd` to `pulumi-resource-{NAME}`
   - replace dependencies in `provider/go.mod` to reflect your repository name
   - find and replace all instances of `provider-boilerplate` with the `NAME` of your provider.
   - find and replace all instances of the boilerplate `abc` with the `ORG` of your provider.
   - replace all instances of the `github.com/pulumi/pulumi-provider-boilerplate` repository with the `REPOSITORY` location

#### Build the provider and install the plugin

   ```bash
   $ make build install
   ```
   
This will:

1. Create the SDK codegen binary and place it in a `./bin` folder (gitignored)
2. Create the provider binary and place it in the `./bin` folder (gitignored)
3. Generate the dotnet, Go, Node, and Python SDKs and place them in the `./sdk` folder
4. Install the provider on your machine.

#### Test against the example
   
```bash
$ cd examples/simple
$ yarn link @pulumi/provider-boilerplate
$ yarn install
$ pulumi stack init test
$ pulumi up
```

Now that you have completed all of the above steps, you have a working provider that generates a random string for you.

#### A brief repository overview

You now have:

1. A `provider/` folder containing the building and implementation logic.
    1. `cmd/pulumi-resource-provider-boilerplate/main.go` - holds the provider's sample implementation logic.
2. `Makefile` - targets to help with building and publishing the provider. Run `make ci-mgmt` to regenerate CI workflows.
3. `sdk` - holds the generated code libraries created by `pulumi gen-sdk`.
4. `examples` a folder of Pulumi programs to try locally and/or use in CI.
5. A `Makefile` and this `README`.

#### Additional Details

This repository depends on the pulumi-go-provider library. For more details on building providers, please check
the [Pulumi Go Provider docs](https://github.com/pulumi/pulumi-go-provider).

### Build Examples

Create an example program using the resources defined in your provider, and place it in the `examples/` folder.

You can now repeat the steps for [build, install, and test](#test-against-the-example).

## Configuring CI and releases

1. Follow the instructions laid out in the [deployment templates](./deployment-templates/README-DEPLOYMENT.md).

## References

Other resources/examples for implementing providers:
* [Pulumi Command provider](https://github.com/pulumi/pulumi-command/blob/master/provider/pkg/provider/provider.go)
* [Pulumi Go Provider repository](https://github.com/pulumi/pulumi-go-provider)
