
import {
  generateKeyPairSync,
  sign,
  verify,
} from "node:crypto";

const { publicKey, privateKey } = generateKeyPairSync(
  "ec",
  {
    namedCurve: "secp256k1",
  },
);

const message = Buffer.from(
  "Alice authorizes a payment of 5000",
);

const signature = sign(
  "sha256",
  message,
  privateKey,
);

const isValid = verify(
  "sha256",
  message,
  publicKey,
  signature,
);

console.log("Signature valid:", isValid);

const modifiedMessage = Buffer.from(
  "Alice authorizes a payment of 50000",
);

const modifiedMessageIsValid = verify(
  "sha256",
  modifiedMessage,
  publicKey,
  signature,
);

console.log(
  "Modified message valid:",
  modifiedMessageIsValid,
);
