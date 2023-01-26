export interface IK2Body extends Record<string, unknown> {
  vars?: Record<string, string>;
}

export interface IK2<TBody extends IK2Body = Record<string, unknown>> {
  k2: {
    metadata: {
      id: string;
      kind: string;
      path: string;
      folder: string;
    };
    body: TBody;
  };
}
