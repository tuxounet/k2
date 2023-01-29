import { IK2TemplateRefGitParams } from "./IK2TemplateRefGitParams";
import { IK2TemplateRefInventoryParams } from "./IK2TemplateRefInventoryParams";

export interface IK2TemplateRef {
  source: "inventory" | "git";
  params: IK2TemplateRefInventoryParams | IK2TemplateRefGitParams;
}
