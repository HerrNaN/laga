import { openDB } from "idb";
import type { Item } from "./lib/item";

const DB_NAME = "laga";
const DB_VERSION = 1;
const STORE_NAME = "items";

async function initDB() {
  return openDB(DB_NAME, DB_VERSION, {
    upgrade(db) {
      if (!db.objectStoreNames.contains(STORE_NAME)) {
        db.createObjectStore(STORE_NAME, { keyPath: "id" });
      }
    },
  });
}

export async function getAllItems(): Promise<Item[]> {
  const db = await initDB();
  return db.getAll(STORE_NAME);
}

export async function setItems(items: Item[]) {
  const db = await initDB();
  const tx = db.transaction(STORE_NAME, "readwrite");
  const store = tx.objectStore(STORE_NAME);
  await store.clear();
  for (const item of items) {
    await store.put(item);
  }
  await tx.done;
}
