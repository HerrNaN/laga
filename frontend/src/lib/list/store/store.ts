import { writable } from "svelte/store";
import { LoroDoc, LoroMap } from "loro-crdt";
import { ItemSchema, type Item } from "../list";
import { classify, departments } from "../classifier";
import type z from "zod";
import { WrapMap } from "./loro";

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

export type ListSummary = {
  id: string;
  name: string;
  itemCount: number;
};

export type ActiveList = {
  id: string;
  name: string;
  items: Item[];
};

// ── Factory ─────────────────────────────────────────────

export const createItemsStore = async (deps: {
  persistence: {
    load: (doc: LoroDoc, key: string) => Promise<boolean>;
    save: (doc: LoroDoc, key: string) => Promise<void>;
  };
  randomUUID: () => string;
}) => {
  const { persistence, randomUUID } = deps;

  // ── Stores ─────────────────────────────────────────────

  const listsWritable = writable<ListSummary[]>([]);
  const activeListWritable = writable<ActiveList>({
    id: "",
    name: "",
    items: [],
  });

  // ── Registry doc ────────────────────────────────────────

  const registryDoc = new LoroDoc<Registry>();
  await persistence.load(registryDoc, "registry");
  const registryLists = registryDoc.getMap("lists");

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
    doc.subscribe(rebuildLists);

    if (!doc.getMap("meta").get("name")) {
      doc.getMap("meta").set("name", "Shopping List");
      doc.commit();
    }

    return doc;
  };

  const createListDoc = () => {
    const id = randomUUID();
    const doc = new LoroDoc<List>();
    setupAutoSave(doc, id);
    doc.getMap("meta").set("name", "Shopping List");
    doc.getMap("items");
    listDocs.set(id, doc);

    registryLists.set(id, 0);
    registryDoc.commit();

    return id;
  };

  // ── Active doc accessors ────────────────────────────────

  let activeListId = "";

  const getActiveDoc = (): LoroDoc<List> => {
    const doc = listDocs.get(activeListId);
    if (!doc) throw new Error(`List doc ${activeListId} not loaded`);
    return doc;
  };

  const getActiveItems = (): LoroMap<Items> => getActiveDoc().getMap("items");

  // ── Rebuilds ────────────────────────────────────────────

  const readActiveName = (): string =>
    (getActiveDoc().getMap("meta").get("name") as string | undefined) ?? "";

  const rebuildActiveList = () => {
    const itemsMap = getActiveItems();
    activeListWritable.set({
      id: activeListId,
      name: readActiveName(),
      items: itemsMap
        .keys()
        .map((id) => ItemSchema.parse(itemsMap.get(id).toJSON())),
    });
  };

  const rebuildLists = () => {
    listsWritable.set(
      registryLists.keys().map((id) => {
        const doc = listDocs.get(id);
        const meta = doc?.getMap("meta");
        const itemsMap = doc?.getMap("items");
        return {
          id,
          name: (meta?.get("name") as string | undefined) ?? "",
          itemCount: itemsMap ? itemsMap.size : 0,
        };
      }),
    );
    rebuildActiveList();
  };

  // ── Active list ─────────────────────────────────────────

  const ensureDefaultList = (): string =>
    registryLists.entries()[0]?.[0] ?? createListDoc();

  activeListId = ensureDefaultList();
  await loadListDoc(activeListId);

  registryDoc.subscribe(() => {
    persistence.save(registryDoc, "registry");
    rebuildLists();
  });

  // Load remaining list docs in background
  Promise.all(
    registryLists
      .keys()
      .filter((id) => id !== activeListId)
      .map((id) =>
        loadListDoc(id).catch((reason) =>
          console.warn("failed to load list doc:", reason),
        ),
      ),
  );

  // ── Store ───────────────────────────────────────────────

  let activeUnsub: (() => void) | undefined;

  rebuildActiveList();
  activeUnsub = getActiveDoc().subscribe(rebuildActiveList);
  rebuildLists();

  const setActiveList = async (id: string) => {
    await loadListDoc(id);
    activeListId = id;
    activeUnsub?.();
    rebuildActiveList();
    activeUnsub = getActiveDoc().subscribe(rebuildActiveList);
  };

  return {
    lists: { subscribe: listsWritable.subscribe },
    activeList: { subscribe: activeListWritable.subscribe },
    setActiveList,
    addItem: (text: string) => {
      const { id: department } = classify(text, departments);
      const id = randomUUID();
      const itemMap = new LoroMap();
      new WrapMap(ItemSchema.shape, itemMap).assign({
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
        const item = new WrapMap(ItemSchema.shape, child);
        if (item.get("checkedAt") !== undefined) {
          item.delete("checkedAt");
        } else {
          item.set("checkedAt", new Date());
        }
        getActiveDoc().commit();
      }
    },
    updateItem: (id: string, updates: Partial<Item>) => {
      const child = getActiveItems().get(id);
      const item = new WrapMap(ItemSchema.shape, child);
      item.update(updates);
      getActiveDoc().commit();
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
