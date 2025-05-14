import * as pulumi from "@pulumi/pulumi";
import * as xyz from "@mynamespace/xyz";

const myRandomResource = new xyz.Random("myRandomResource", {length: 24});
const myRandomComponent = new xyz.RandomComponent("myRandomComponent", {length: 24});
export const output = {
    value: myRandomResource.result,
};
