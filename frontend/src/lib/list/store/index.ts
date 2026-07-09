import { createItemsStore } from "./store";
import { LoroPersistence } from "../../loro/persistence";

const persistence = new LoroPersistence();
export const listStore = await createItemsStore({
  persistence,
  randomUUID: () => crypto.randomUUID(),
});
