import { LoroMap } from "loro-crdt";
import * as z from "zod";

type Out<T extends Record<string, z.ZodTypeAny>> = {
  [K in keyof T]: z.output<T[K]>;
};

export class WrapMap<T extends Record<string, z.ZodTypeAny>> {
  constructor(
    private shape: T,
    private map: LoroMap,
  ) {}
  get = <K extends keyof T & string>(key: K): Out<T>[K] | undefined => {
    const raw = this.map.get(key);
    if (raw === undefined) return undefined;
    return this.shape[key].parse(raw);
  };
  set = <K extends keyof T & string>(key: K, value: Out<T>[K]): void => {
    if (value === undefined) {
      this.map.delete(key as string);
      return;
    }
    this.map.set(key as string, this.shape[key].encode(value));
  };
  delete = <K extends keyof T & string>(key: K) => this.map.delete(key);
  assign = (data: Partial<Out<T>>): void => {
    for (const [key, value] of Object.entries(data)) {
      if (value !== undefined) {
        this.map.set(key, this.shape[key as keyof T]!.encode(value));
      } else {
        this.map.delete(key);
      }
    }
  };
  toJSON = () => z.object(this.shape).parse(this.map.toJSON()) as Out<T>;
  raw = () => this.map;
}
