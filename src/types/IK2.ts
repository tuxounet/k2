export interface IK2Body {}

export interface IK2<TBody extends IK2Body = IK2Body> {
  k2: {
    metadata: {
      id: string;
      kind: string;
      version?: string;
      path: string;
      folder: string;
    };
    body: TBody;
  };
}
