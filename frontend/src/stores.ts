import { writable } from "svelte/store";
import type { Item } from "./lib/item";
import { setItems } from "./db";

const createItemsStore = () => {
  const { subscribe, set, update } = writable<Item[]>([]);
  let loaded = false;

  subscribe(async (items) => {
    if (loaded) {
      await setItems(items);
    }
  });

  return {
    subscribe,
    set,
    hydrate: (initialItems: Item[]) => {
      loaded = false;
      set(initialItems);
      loaded = true;
    },
    addItem: (text: string) => {
      const item: Item = {
        id: crypto.randomUUID(),
        text,
        checked: false,
        createdAt: new Date(),
      };
      update((items) => [...items, item]);
    },
    toggleItem: (id: string) => {
      update((items) =>
        items.map((item) =>
          item.id === id ? { ...item, checked: !item.checked } : item,
        ),
      );
    },
    deleteItem: (id: string) => {
      update((items) => items.filter((item) => item.id !== id));
    },
  };
};

export const items = createItemsStore();
