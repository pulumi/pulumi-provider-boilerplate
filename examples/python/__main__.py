import pulumi
import mynamespace_xyz as xyz

my_random_resource = xyz.Random("myRandomResource", length=24)
my_random_component = xyz.RandomComponent("myRandomComponent", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
