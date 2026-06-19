import { describe, it, expect, beforeEach } from "vitest";
import { LoroDoc } from "loro-crdt";
import { LoroPersistence } from "./persistence";
import { IDBFactory } from "fake-indexeddb";

describe("LoroPersistence", () => {
  let persistence: LoroPersistence;

  beforeEach(() => {
    persistence = new LoroPersistence(new IDBFactory());
  });

  it("returns true when data exists for a key", async () => {
    const doc = new LoroDoc();
    doc.getMap("meta").set("name", "test list");
    await persistence.save(doc, "list-1");

    const loaded = new LoroDoc();
    const result = await persistence.load(loaded, "list-1");
    expect(result).toBe(true);
    expect(loaded.getMap("meta").get("name")).toBe("test list");
  });

  it("returns false when key does not exist", async () => {
    const loaded = new LoroDoc();
    const result = await persistence.load(loaded, "nonexistent");
    expect(result).toBe(false);
    expect(loaded.getMap("meta").size).toBe(0);
  });

  it("returns true on overwrite", async () => {
    const doc1 = new LoroDoc();
    doc1.getMap("meta").set("name", "first");
    await persistence.save(doc1, "list-1");

    const doc2 = new LoroDoc();
    doc2.getMap("meta").set("name", "second");
    await persistence.save(doc2, "list-1");

    const loaded = new LoroDoc();
    const result = await persistence.load(loaded, "list-1");
    expect(result).toBe(true);
    expect(loaded.getMap("meta").get("name")).toBe("second");
  });

  it("handles multiple keys independently", async () => {
    const a = new LoroDoc();
    a.getMap("meta").set("name", "A");
    await persistence.save(a, "key-a");

    const b = new LoroDoc();
    b.getMap("meta").set("name", "B");
    await persistence.save(b, "key-b");

    const loadedB = new LoroDoc();
    expect(await persistence.load(loadedB, "key-b")).toBe(true);
    expect(loadedB.getMap("meta").get("name")).toBe("B");

    const loadedA = new LoroDoc();
    expect(await persistence.load(loadedA, "key-a")).toBe(true);
    expect(loadedA.getMap("meta").get("name")).toBe("A");
  });
});
