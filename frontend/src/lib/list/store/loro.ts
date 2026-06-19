import { LoroMap } from "loro-crdt";
import * as z from "zod";

type Out<T extends Record<string, z.ZodTypeAny>> = {
  [K in keyof T]: z.output<T[K]>;
};

export const wrapMap = <T extends Record<string, z.ZodTypeAny>>(
  shape: T,
  map: LoroMap,
) => ({
  inner: map,
  get: <K extends keyof T & string>(key: K): Out<T>[K] | undefined => {
    const raw = map.get(key);
    if (raw === undefined) return undefined;
    return shape[key].parse(raw);
  },
  set: <K extends keyof T & string>(key: K, value: Out<T>[K]): void => {
    if (value === undefined) {
      map.delete(key as string);
      return;
    }
    map.set(key as string, shape[key].encode(value));
  },
  delete: <K extends keyof T & string>(key: K) => map.delete(key),
  assign: (data: Partial<Out<T>>): void => {
    for (const [key, value] of Object.entries(data)) {
      if (value !== undefined) {
        map.set(key, shape[key as keyof T]!.encode(value));
      } else {
        map.delete(key);
      }
    }
  },
  toJSON: () => z.object(shape).parse(map.toJSON()) as Out<T>,
  raw: () => map,
});
