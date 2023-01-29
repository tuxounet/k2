import crypto from "crypto";
export function md5(message: string) {
  return crypto.createHash("md5").update(message).digest("hex");
}
