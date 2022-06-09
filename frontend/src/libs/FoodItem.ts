import { ItemKind } from "./ItemKind";

export type FoodItem = {
  id: string;
  name: string;
  description: string;
  priority: number;
  price: number;
  maxOrder: number;
  maxOrderPerDay: number;
  enabled: boolean;
  kind: ItemKind | null;
  scheduleIds: string[];
};
