import { LoroDoc } from "loro-crdt";

const DB_NAME = "laga-crdt";
const STORE = "snapshots";
const VERSION = 1;

const openDB = (): Promise<IDBDatabase> =>
  new Promise((resolve, reject) => {
    const req = indexedDB.open(DB_NAME, VERSION);
    req.onupgradeneeded = () => {
      req.result.createObjectStore(STORE);
    };
    req.onerror = () => reject(req.error);
    req.onsuccess = () => resolve(req.result);
  });

export const loadDoc = async (doc: LoroDoc, key: string): Promise<void> => {
  const db = await openDB();
  try {
    const blob = await new Promise<Blob | undefined>((resolve, reject) => {
      const tx = db.transaction(STORE, "readonly");
      const req = tx.objectStore(STORE).get(key);
      req.onerror = () => reject(req.error);
      req.onsuccess = () => resolve(req.result ?? undefined);
    });
    if (blob) {
      doc.import(new Uint8Array(await blob.arrayBuffer()));
    }
  } finally {
    db.close();
  }
};

export const saveDoc = async (doc: LoroDoc, key: string): Promise<void> => {
  const snapshot = doc.export({ mode: "snapshot" });
  const db = await openDB();
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
};
