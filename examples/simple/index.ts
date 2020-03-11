import * as xyz from "../../sdk/nodejs";

const random = new xyz.Random("my-random", { length: 24 });

export const output = random.result;