import crypto from "crypto";
export function md5(message: string): string {
  return crypto.createHash("md5").update(message).digest("hex");
}
