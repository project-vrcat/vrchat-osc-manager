import { Manager } from "https://deno.land/x/vrchat_osc_manager@v0.1.0/mod.ts";

const manager = new Manager();
await manager.connect();

const options = await manager.getOptions();
const parameters = options.parameters as string[];
manager.listenParameters(parameters);
manager.on("parameters", (name, value) =>
  console.log("Parameter:", name, value)
);
