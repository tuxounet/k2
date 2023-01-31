export const separator = ".";

export function idNs(id: string): string {
  const segemnts = id.split(separator);

  return segemnts
    .filter((item, idx) => idx < segemnts.length - 1)
    .join(separator);
}

export function idName(id: string): string {
  const segemnts = id.split(separator);
  return segemnts.find((item, idx) => idx === segemnts.length - 1) ?? "unknow";
}

export function idVersion(id?: string): string {
  return id ?? "latest";
}

export function idRoot(id: string): string {
  if (id.includes(separator)) {
    const segemnts = id.split(separator);
    return segemnts.find((item, idx) => idx === 0) ?? "unknow";
  }
  return id;
}
