import { IK2Template } from "../types/IK2Template";
import { IK2TemplateRef } from "../types/templates/IK2TemplateRef";
import { IK2TemplateRefGitParams } from "../types/templates/IK2TemplateRefGitParams";
import { IK2TemplateRefInventoryParams } from "../types/templates/IK2TemplateRefInventoryParams";
import { Inventory } from "./Inventory";
import fs from "fs";
import path from "path";
import { md5 } from "../helpers/hash";
import { exec } from "../helpers/exec";
export function resolveTemplate(
  inventory: Inventory,
  ref: IK2TemplateRef
): IK2Template {
  switch (ref.source) {
    case "inventory": {
      const param = ref.params as IK2TemplateRefInventoryParams;
      const template = inventory.templates.get(param.id);
      if (template == null) {
        throw new Error(`template introuvable dans l'inventaire ${param.id}`);
      }
      return template;
    }
    case "git": {
      const param = ref.params as IK2TemplateRefGitParams;
      const template = resolveGitTemplate(inventory, param);
      if (template == null) {
        throw new Error(`template non résolu ${param.repository}`);
      }
      return template;
    }
    default:
      throw new Error(`source de template non trouvé ${String(ref.source)}`);
  }
}

function resolveGitTemplate(
  inventory: Inventory,
  refParams: IK2TemplateRefGitParams
): IK2Template {
  console.info("resolve", refParams);

  const id = md5(
    JSON.stringify({
      repository: refParams.repository,
      branch: refParams.branch,
    })
  );

  const templateRefPath = path.join(inventory.inventoryFolder, "refs", id);
  if (!fs.existsSync(templateRefPath)) {
    exec(
      `git clone  ${
        refParams.branch !== undefined
          ? `--branch ${refParams.branch} --single-branch`
          : ""
      } ${refParams.repository} ${templateRefPath}`,
      inventory.inventoryFolder
    );
  } else {
    exec(`git pull`, templateRefPath);
  }

  const targetPath = templateRefPath + "/" + refParams.path;

  if (!fs.existsSync(targetPath)) {
    throw new Error(
      "Impossible de trouver le fichier de template " + targetPath
    );
  }

  const templateFile = inventory.loadK2File<IK2Template>(targetPath);
  return templateFile;
}
