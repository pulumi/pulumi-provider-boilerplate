// *** WARNING: this file was generated by pulumi-language-nodejs. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

import * as pulumi from "@pulumi/pulumi";
import * as utilities from "./utilities";

export class Random extends pulumi.CustomResource {
    /**
     * Get an existing Random resource's state with the given name, ID, and optional extra
     * properties used to qualify the lookup.
     *
     * @param name The _unique_ name of the resulting resource.
     * @param id The _unique_ provider ID of the resource to lookup.
     * @param opts Optional settings to control the behavior of the CustomResource.
     */
    public static get(name: string, id: pulumi.Input<pulumi.ID>, opts?: pulumi.CustomResourceOptions): Random {
        return new Random(name, undefined as any, { ...opts, id: id });
    }

    /** @internal */
    public static readonly __pulumiType = 'provider-boilerplate:index:Random';

    /**
     * Returns true if the given object is an instance of Random.  This is designed to work even
     * when multiple copies of the Pulumi SDK have been loaded into the same process.
     */
    public static isInstance(obj: any): obj is Random {
        if (obj === undefined || obj === null) {
            return false;
        }
        return obj['__pulumiType'] === Random.__pulumiType;
    }

    public readonly length!: pulumi.Output<number>;
    public /*out*/ readonly result!: pulumi.Output<string>;

    /**
     * Create a Random resource with the given unique name, arguments, and options.
     *
     * @param name The _unique_ name of the resource.
     * @param args The arguments to use to populate this resource's properties.
     * @param opts A bag of options that control this resource's behavior.
     */
    constructor(name: string, args: RandomArgs, opts?: pulumi.CustomResourceOptions) {
        let resourceInputs: pulumi.Inputs = {};
        opts = opts || {};
        if (!opts.id) {
            if ((!args || args.length === undefined) && !opts.urn) {
                throw new Error("Missing required property 'length'");
            }
            resourceInputs["length"] = args ? args.length : undefined;
            resourceInputs["result"] = undefined /*out*/;
        } else {
            resourceInputs["length"] = undefined /*out*/;
            resourceInputs["result"] = undefined /*out*/;
        }
        opts = pulumi.mergeOptions(utilities.resourceOptsDefaults(), opts);
        super(Random.__pulumiType, name, resourceInputs, opts);
    }
}

/**
 * The set of arguments for constructing a Random resource.
 */
export interface RandomArgs {
    length: pulumi.Input<number>;
}
