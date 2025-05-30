import pulumi
import provider_boilerplate as boilerplate

my_random_resource = boilerplate.Random("myRandomResource", length=24)
my_random_component = boilerplate.RandomComponent("myRandomComponent", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
