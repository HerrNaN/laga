import { LoroDoc } from "loro-crdt";

const DB_NAME = "laga-crdt";
const STORE = "snapshots";
const VERSION = 1;

export class LoroPersistence {
  constructor(private idbFactory: IDBFactory = indexedDB) {}

  async load(doc: LoroDoc, key: string): Promise<boolean> {
    const db = await this.openDB();
    try {
      const blob = await new Promise<Blob | undefined>((resolve, reject) => {
        const tx = db.transaction(STORE, "readonly");
        const req = tx.objectStore(STORE).get(key);
        req.onerror = () => reject(req.error);
        req.onsuccess = () => resolve(req.result ?? undefined);
      });
      if (blob) {
        doc.import(new Uint8Array(await blob.arrayBuffer()));
        return true;
      }
      return false;
    } finally {
      db.close();
    }
  }

  async save(doc: LoroDoc, key: string): Promise<void> {
    const snapshot = doc.export({ mode: "snapshot" });
    const db = await this.openDB();
    try {
      await new Promise<void>((resolve, reject) => {
        const tx = db.transaction(STORE, "readwrite");
        tx.objectStore(STORE).put(new Blob([snapshot as BlobPart]), key);
        tx.oncomplete = () => resolve();
        tx.onerror = () => reject(tx.error);
      });
    } finally {
      db.close();
    }
  }

  private async openDB(): Promise<IDBDatabase> {
    return new Promise((resolve, reject) => {
      const req = this.idbFactory.open(DB_NAME, VERSION);
      req.onupgradeneeded = () => {
        req.result.createObjectStore(STORE);
      };
      req.onerror = () => reject(req.error);
      req.onsuccess = () => resolve(req.result);
    });
  }
}
