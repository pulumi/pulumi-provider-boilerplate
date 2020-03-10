import * as pulumi from "@pulumi/pulumi";

interface RandomArgs {
    length: pulumi.Input<number>;
}

export class Random extends pulumi.CustomResource {
    result: pulumi.Output<string>;
    constructor(name: string, args: RandomArgs, opts?: pulumi.CustomResourceOptions) {
        super("xyz:random:Random", name, {result: undefined, ...args}, opts);
    }
}
