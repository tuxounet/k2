import { IK2, IK2Body } from "./IK2";

export interface IK2TemplateBody extends IK2Body {
  name: string;
  parameters: Record<string, unknown>;
  scripts: {
    bootstrap?: string[];
    pre?: string[];
    post?: string[];
  };
}

export interface IK2Template extends IK2<IK2TemplateBody> {}
