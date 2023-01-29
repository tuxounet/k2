import { IK2Template } from "../types/IK2Template";
import { IK2TemplateRef } from "../types/templates/IK2TemplateRef";
import { IK2TemplateRefInventoryParams } from "../types/templates/IK2TemplateRefInventoryParams";
import { Inventory } from "./Inventory";

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
    default:
      throw new Error("source de template non trouvé " + ref.source);
  }
}
