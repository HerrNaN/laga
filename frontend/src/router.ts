import { createRouter } from "sv-router";
import List from "./lib/list/List.svelte";
import Layout from "./Layout.svelte";
import CatchAll from "./CatchAll.svelte";

export const { p, navigate, isActive, route } = createRouter({
  "/lists": List,
  "*": CatchAll,
  layout: Layout,
});
