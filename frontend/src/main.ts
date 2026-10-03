import { mount } from "svelte";
import App from "./App.svelte";
import { unregisterServiceWorkers } from "./lib/service-worker";

import "@awesome.me/webawesome/dist/styles/themes/default.css";

unregisterServiceWorkers();

const app = mount(App, {
  target: document.getElementById("app")!,
});

export default app;
