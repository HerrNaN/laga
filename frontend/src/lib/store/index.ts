import { createItemsStore } from "./store";
import { LoroPersistence } from "./persistence";

const persistence = new LoroPersistence();
export const items = await createItemsStore({
  persistence,
  randomUUID: crypto.randomUUID,
});
