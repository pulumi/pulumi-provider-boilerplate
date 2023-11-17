import pulumi
import pulumi_xyz as xyz

my_random_resource = xyz.Random("myRandomResource", length=24)
pulumi.export("output", {
    "value": my_random_resource.result,
})
