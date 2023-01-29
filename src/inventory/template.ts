import { IK2Template } from "../types/IK2Template";
import { IK2TemplateRef } from "../types/templates/IK2TemplateRef";
import { IK2TemplateRefGitParams } from "../types/templates/IK2TemplateRefGitParams";
import { IK2TemplateRefInventoryParams } from "../types/templates/IK2TemplateRefInventoryParams";
import { Inventory } from "./Inventory";
import fs from "fs";
import path from "path";
import childProc from "child_process";
export function resolveTemplate(
  inventory: Inventory,
  ref: IK2TemplateRef
): IK2Template {
  switch (ref.source) {
    case "inventory": {
      const param = ref.params as IK2TemplateRefInventoryParams;
      const template = inventory.templates.get(param.id);
      if (template == null) {
        throw new Error("template introuvable dans l'inventaire " + param.id);
      }
      return template;
    }
    case "git": {
      const param = ref.params as IK2TemplateRefGitParams;
      const template = resolveGitTemplate(inventory, param);
      if (template == null) {
        throw new Error("template non résolu  " + param.repository);
      }
      return template;
    }
    default:
      throw new Error("source de template non trouvé " + ref.source);
  }
}

function resolveGitTemplate(
  inventory: Inventory,
  refParams: IK2TemplateRefGitParams
) {
  throw new Error("template non résolu  " + refParams.repository);
}
