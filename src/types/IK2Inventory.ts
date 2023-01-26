import { IK2, IK2Body } from "./IK2";

export interface IK2InventoryBody extends IK2Body {
  folders: {
    ignore: string[];
    sources: string[];
    templates: string[];
  };
}

export interface IK2Inventory extends IK2<IK2InventoryBody> {}
