import { IK2, IK2Body } from "./IK2";

export interface IK2TemplateBody extends IK2Body {
  name: string;
}

export interface IK2Template extends IK2<IK2TemplateBody> {}
