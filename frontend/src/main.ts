import { mount } from "svelte";
import App from "./App.svelte";

import "@awesome.me/webawesome/dist/styles/themes/default.css";

const app = mount(App, {
  target: document.getElementById("app")!,
});

export default app;
