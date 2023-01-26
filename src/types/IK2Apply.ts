import { IK2, IK2Body } from "./IK2";

export interface IK2ApplyBody extends IK2Body {
  template: string;
}

export interface IK2Apply extends IK2<IK2ApplyBody> {}
