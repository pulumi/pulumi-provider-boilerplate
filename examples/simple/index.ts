import * as xyz from "@pulumi/xyz";

const random = new xyz.Random("my-random", { length: 24 });

export const output = random.result;