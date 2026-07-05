import z from "zod";
import { Database } from "../indexeddb/database";

const DB_NAME = "laga_user";
const VERSION = 1;
const STORE = "user";

const CryptoKeySchema = z.custom<CryptoKey>(
  (val: unknown) =>
    z
      .object({
        type: z.union([
          z.literal("secret"),
          z.literal("private"),
          z.literal("public"),
        ]),
        extractable: z.boolean(),
        algorithm: z.unknown(),
        usages: z.array(
          z.union([
            z.literal("encrypt"),
            z.literal("decrypt"),
            z.literal("sign"),
            z.literal("verify"),
            z.literal("deriveKey"),
            z.literal("deriveBits"),
            z.literal("wrapKey"),
            z.literal("unwrapKey"),
          ]),
        ),
      })
      .safeParse(val).success,
);

const CryptoKeyPairSchema = z.custom<CryptoKeyPair>(
  (val: unknown) =>
    z
      .object({
        privateKey: CryptoKeySchema,
        publicKey: CryptoKeySchema,
      })
      .safeParse(val).success,
);

const UserStoreSchema = z.object({
  userId: z.string(),
  deviceKeyPair: CryptoKeyPairSchema,
});

type User = {
  userId: string;
  deviceKeyPair: CryptoKeyPair;
};

const setupUser = async (): Promise<User> => ({
  userId: await userDB.getOrAdd("userId", generateUserId()),
  deviceKeyPair: await userDB.getOrAdd(
    "deviceKeyPair",
    await generateDeviceKeyPair(),
  ),
});

const generateUserId = (): string => {
  return crypto.randomUUID();
};

const generateDeviceKeyPair = async (): Promise<CryptoKeyPair> =>
  await self.crypto.subtle.generateKey(
    {
      name: "ECDSA",
      namedCurve: "P-384",
    },
    false,
    ["sign"],
  );

const userDB = new Database(UserStoreSchema, DB_NAME, VERSION, STORE);

export const user = await setupUser();
