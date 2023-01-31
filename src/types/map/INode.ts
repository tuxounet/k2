export interface INode {
  id: string;
  name: string;
  namespace: string;
  kind: string;
  version: string;
  childs?: INode[];
}
