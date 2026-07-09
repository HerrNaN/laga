import { createRouter } from "sv-router";
import List from "./lib/list/List.svelte";
import Lists from "./lib/list/Lists.svelte";
import Layout from "./Layout.svelte";
import CatchAll from "./CatchAll.svelte";

export const { p, navigate, isActive, route } = createRouter({
  "/lists": {
    "/": Lists,
    "/:id": List,
  },
  "*": CatchAll,
  layout: Layout,
});
