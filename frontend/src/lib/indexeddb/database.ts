import { core, type ZodObject, type ZodType } from "zod";

type ZShape = Readonly<Record<string, ZodType>>;

export class Database<T extends ZShape> {
  constructor(
    private schema: ZodObject<T>,
    private name: string,
    private version: number,
    private store: string,
    private idbFactory: IDBFactory = self.indexedDB,
  ) {}

  open = async () =>
    new Promise<IDBDatabase>((resolve, reject) => {
      const req = this.idbFactory.open(this.name, this.version);
      req.onupgradeneeded = () => {
        req.result.createObjectStore(this.store);
      };
      req.onerror = () => reject(req.error);
      req.onsuccess = () => resolve(req.result);
    });

  get = async <K extends keyof T & string>(
    key: K,
  ): Promise<core.output<T[K]> | undefined> => {
    const db = await this.open();

    try {
      const value = await new Promise<unknown | undefined>(
        (resolve, reject) => {
          const tx = db.transaction(this.store, "readonly");
          const req = tx.objectStore(this.store).get(key);
          req.onerror = () => reject(req.error);
          req.onsuccess = () => resolve(req.result ?? undefined);
        },
      );

      if (value === undefined) return undefined;
      return await this.schema.shape[key].parseAsync(value);
    } finally {
      db.close();
    }
  };

  put = async <K extends keyof T & string>(
    key: K,
    value: core.output<T[K]>,
  ): Promise<core.output<T[K]>> => {
    const db = await this.open();

    try {
      const encoded = await this.schema.shape[key].encodeAsync(value);
      return new Promise<core.output<T[K]>>((resolve, reject) => {
        const tx = db.transaction(this.store, "readwrite");
        tx.objectStore(this.store).put(encoded, key);
        tx.oncomplete = () => resolve(value);
        tx.onerror = () => reject(tx.error);
      });
    } finally {
      db.close();
    }
  };

  add = async <K extends keyof T & string>(
    key: K,
    value: core.output<T[K]>,
  ): Promise<core.output<T[K]>> => {
    const db = await this.open();

    try {
      const encoded = await this.schema.shape[key].encodeAsync(value);
      return new Promise<core.output<T[K]>>((resolve, reject) => {
        const tx = db.transaction(this.store, "readwrite");
        tx.objectStore(this.store).add(encoded, key);
        tx.oncomplete = () => resolve(value);
        tx.onerror = () => reject(tx.error);
      });
    } finally {
      db.close();
    }
  };

  getOrAdd = async <K extends keyof T & string>(
    key: K,
    value: core.output<T[K]>,
  ): Promise<core.output<T[K]>> =>
    this.get(key).then((old) => old || this.add(key, value));
}
