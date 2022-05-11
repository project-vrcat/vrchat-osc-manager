import { Manager } from "https://cdn.jsdelivr.net/gh/project-vrcat/vrchat-osc-manager-plugins/module/deno/vrchat-osc-manager.ts";

const manager = new Manager();
await manager.connect();

const options = await manager.getOptions();
const parametes = options.parametes as string[];
manager.listenParameters(parametes);
manager.on(
  "parameters",
  (name, value) => console.log("Parameter:", name, value),
);
