import { IK2, IK2Body } from "./IK2";
import { IK2TemplateRef } from "./templates/IK2TemplateRef";

export interface IK2ApplyBody extends IK2Body {
  template: IK2TemplateRef;
}

export interface IK2Apply extends IK2<IK2ApplyBody> {}
