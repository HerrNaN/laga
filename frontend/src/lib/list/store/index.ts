import { createItemsStore } from "./store";
import { LoroPersistence } from "../../loro/persistence";

const persistence = new LoroPersistence();
export const items = await createItemsStore({
  persistence,
  randomUUID: () => crypto.randomUUID(),
});
