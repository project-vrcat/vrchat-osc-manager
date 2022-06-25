import { Manager } from "vrchat-osc-manager";

const manager = new Manager();

(async () => {
  await manager.connect();
  const options = await manager.getOptions();
  const parameters = options.parameters as string[];
  manager.listenParameters(parameters);
  manager.on("parameters", (name, value) =>
    console.log("Parameter:", name, value)
  );
})();
