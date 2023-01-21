import { OptionItem } from "./OptionItem"

export type ItemKind = {
  id: string;
  name: string;
  priority: number;
  optionItemIds: string[];
};

export type AggregatedItemKind = {
  id: string;
  name: string;
  priority: number;
  options: OptionItem[];
}