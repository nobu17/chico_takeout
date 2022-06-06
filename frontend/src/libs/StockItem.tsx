import { ItemKind } from "./ItemKind";

export type StockItem = {
  id: string;
  name: string;
  description: string;
  priority: number;
  price: number;
  maxOrder: number;
  enabled: boolean;
  remain: number;
  kind: ItemKind | null;
};
