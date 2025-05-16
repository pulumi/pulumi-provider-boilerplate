using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Xyz = Mynamespace.Xyz;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new Xyz.Random("myRandomResource", new()
    {
        Length = 24,
    });

    var myRandomComponent = new Xyz.RandomComponent("myRandomComponent", new()
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

