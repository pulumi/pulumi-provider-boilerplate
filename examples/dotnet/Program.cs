using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Xyz = Pulumi.Xyz;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new Xyz.Random("myRandomResource", new()
    {
        Length = 24,
    });

    return new Dictionary<string, object?>
    {
        ["output"] = 
        {
            { "value", myRandomResource.Result },
        },
    };
});

