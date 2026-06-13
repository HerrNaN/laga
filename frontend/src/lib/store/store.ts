import { writable } from "svelte/store";
import { LoroDoc, LoroMap } from "loro-crdt";
import { ItemSchema, type Item } from "../list";
import { classify, departments } from "../classifier";
import { wrapMap } from "./loro";
import type z from "zod";

// ── Data model ──────────────────────────────────────────

type Registry = {
  lists: LoroMap<{ [listId: string]: 0 }>;
};

type Items = {
  [itemId: string]: LoroMap<EncodedItem>;
};

type List = {
  meta: LoroMap<{ name: string }>;
  items: LoroMap<Items>;
};

type EncodedItem = z.input<typeof ItemSchema>;

// ── Factory ─────────────────────────────────────────────

export const createItemsStore = async (deps: {
  persistence: {
    load: (doc: LoroDoc, key: string) => Promise<boolean>;
    save: (doc: LoroDoc, key: string) => Promise<void>;
  };
  randomUUID: () => string;
}) => {
  const { persistence, randomUUID } = deps;

  // ── Registry doc ────────────────────────────────────────

  const registryDoc = new LoroDoc<Registry>();
  await persistence.load(registryDoc, "registry");
  const registryLists = registryDoc.getMap("lists");

  registryDoc.subscribe(() => persistence.save(registryDoc, "registry"));

  // ── Per-list docs ───────────────────────────────────────

  const listDocs = new Map<string, LoroDoc<List>>();

  const setupAutoSave = (doc: LoroDoc<List>, id: string) =>
    doc.subscribe(() => persistence.save(doc, `list-${id}`));

  const loadListDoc = async (id: string): Promise<LoroDoc<List>> => {
    let doc = listDocs.get(id);
    if (doc) return doc;

    doc = new LoroDoc<List>();
    await persistence.load(doc, `list-${id}`);
    listDocs.set(id, doc);
    setupAutoSave(doc, id);

    return doc;
  };

  const createListDoc = () => {
    const id = randomUUID();
    const doc = new LoroDoc<List>();
    setupAutoSave(doc, id);
    doc.getMap("meta").set("name", "");
    doc.getMap("items");
    listDocs.set(id, doc);

    registryLists.set(id, 0);
    registryDoc.commit();

    return id;
  };

  // ── Active list ─────────────────────────────────────────

  const ensureDefaultList = (): string =>
    registryLists.entries()[0]?.[0] ?? createListDoc();

  let activeListId = ensureDefaultList();
  await loadListDoc(activeListId);

  // Load remaining list docs in background
  Promise.all(
    registryLists
      .keys()
      .filter((id) => id !== activeListId)
      .map((id) => loadListDoc(id).catch(console.warn)),
  );

  // ── Active doc accessors ────────────────────────────────

  const getActiveDoc = (): LoroDoc<List> => {
    const doc = listDocs.get(activeListId);
    if (!doc) throw new Error(`List doc ${activeListId} not loaded`);
    return doc;
  };

  const getActiveItems = (): LoroMap<Items> => getActiveDoc().getMap("items");

  // ── Store ───────────────────────────────────────────────

  const { subscribe, set } = writable<Item[]>([]);

  const rebuildItems = () => {
    const itemsMap = getActiveItems();
    set(
      itemsMap.keys().map((id) => ItemSchema.parse(itemsMap.get(id).toJSON())),
    );
  };

  rebuildItems();

  getActiveDoc().subscribe(rebuildItems);

  return {
    subscribe,
    addItem: (text: string) => {
      const { id: department } = classify(text, departments);
      const id = randomUUID();
      const itemMap = new LoroMap();
      wrapMap(ItemSchema.shape, itemMap).assign({
        id,
        text,
        createdAt: new Date(),
        department,
      });
      getActiveItems().setContainer(id, itemMap);
      getActiveDoc().commit();
    },
    toggleItem: (id: string) => {
      const child = getActiveItems().get(id);
      if (child instanceof LoroMap) {
        const item = wrapMap(ItemSchema.shape, child);
        if (item.get("checkedAt") !== undefined) {
          item.delete("checkedAt");
        } else {
          item.set("checkedAt", new Date());
        }
        getActiveDoc().commit();
      }
    },
    deleteItem: (id: string) => {
      getActiveItems().delete(id);
      getActiveDoc().commit();
    },
    deleteCheckedItems: () => {
      const items = getActiveItems();
      const checked: string[] = [];
      for (const key of items.keys()) {
        const child = items.get(key);
        if (child instanceof LoroMap && child.get("checkedAt") !== undefined) {
          checked.push(key);
        }
      }
      for (const id of checked) {
        items.delete(id);
      }
      getActiveDoc().commit();
    },
  };
};
