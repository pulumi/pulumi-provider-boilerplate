using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Boilerplate = Mynamespace.ProviderBoilerplate;

return await Deployment.RunAsync(() => 
{
    var myRandomResource = new Boilerplate.Random("myRandomResource", new()
    {
        Length = 24,
    });

    var myRandomComponent = new Boilerplate.RandomComponent("myRandomComponent", new()
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

