import {
  describe,
  it,
  expect,
  vi,
  beforeEach,
  type MockedFunction,
} from "vitest";
import { createItemsStore } from "./store";
import type { Item } from "../list";
import type { ActiveList, ListSummary } from "./store";

// Mock classifier for deterministic results
vi.mock("../classifier", () => ({
  classify: vi.fn().mockReturnValue({ id: "other", name: "Other", score: 0 }),
  departments: [{ id: "other", name: "Other", keywords: [] }],
  score: vi.fn().mockReturnValue([]),
  tokenize: vi.fn().mockReturnValue([]),
}));

let uuidCounter = 0;

const createStore = async () => {
  const persistence = {
    load: vi.fn().mockResolvedValue(undefined) as MockedFunction<
      () => Promise<boolean>
    >,
    save: vi.fn().mockResolvedValue(undefined) as MockedFunction<
      () => Promise<void>
    >,
  };

  const store = await createItemsStore({
    persistence,
    randomUUID: () => {
      uuidCounter += 1;
      return `uuid-${uuidCounter}`;
    },
  });

  return { persistence, store };
};

// Helper: collect current active-list items via subscribe
const collect = (store: {
  activeList: { subscribe: (fn: (v: ActiveList) => void) => () => void };
}): Promise<Item[]> =>
  new Promise((resolve) => {
    store.activeList.subscribe((v) => resolve(v.items));
  });

describe("listStore", () => {
  beforeEach(() => {
    uuidCounter = 0;
    vi.clearAllMocks();
  });

  it("starts with an empty list", async () => {
    const { store } = await createStore();
    expect(await collect(store)).toEqual([]);
  });

  describe("lists", () => {
    it("exposes the registered lists with default name and zero items", async () => {
      const { store } = await createStore();
      const lists = await new Promise<ListSummary[]>((resolve) =>
        store.lists.subscribe(resolve),
      );

      expect(lists).toHaveLength(1);
      expect(lists[0].name).toBe("Shopping List");
      expect(lists[0].itemCount).toBe(0);
    });
  });

  describe("addItem", () => {
    it("adds an item", async () => {
      const { store } = await createStore();
      store.addItem("milk");

      const items = await collect(store);
      expect(items).toHaveLength(1);
      expect(items[0].text).toBe("milk");
    });
  });

  describe("toggleItem", () => {
    it("toggles checkedAt on an item", async () => {
      const { store } = await createStore();
      store.addItem("eggs");

      const [item] = await collect(store);
      expect(item.checkedAt).toBeUndefined();

      store.toggleItem(item.id);
      const [checked] = await collect(store);
      expect(checked.checkedAt).toBeInstanceOf(Date);

      store.toggleItem(item.id);
      const [unchecked] = await collect(store);
      expect(unchecked.checkedAt).toBeUndefined();
    });
  });

  describe("deleteItem", () => {
    it("deletes an item", async () => {
      const { store } = await createStore();
      store.addItem("cheese");
      const [item] = await collect(store);

      store.deleteItem(item.id);
      expect(await collect(store)).toEqual([]);
    });
  });

  describe("deleteCheckedItems", () => {
    it("deletes only checked items", async () => {
      const { store } = await createStore();
      store.addItem("a");
      store.addItem("b");
      store.addItem("c");
      const all = await collect(store);
      store.toggleItem(all[0].id);
      store.toggleItem(all[2].id);

      store.deleteCheckedItems();
      const remaining = await collect(store);
      expect(remaining).toHaveLength(1);
      expect(remaining[0].text).toBe("b");
    });
  });

  it("calls persistence.save after mutations", async () => {
    const { persistence, store } = await createStore();
    // initial saves from registry + list doc creation
    persistence.save.mockClear();

    store.addItem("butter");
    expect(persistence.save).toHaveBeenCalled();

    store.toggleItem((await collect(store))[0].id);
    persistence.save.mockClear();
    store.deleteCheckedItems();
    expect(persistence.save).toHaveBeenCalled();
  });
});
